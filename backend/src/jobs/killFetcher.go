package jobs

import (
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
)

var (
	fetchQueue = make(chan int64, 100)
	jobMutex   sync.Mutex
	jobRunning bool
)

func StartKillFetcherJob() {
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		if !isJobRunning() {
			log.Println("Starting scheduled kill fetch for all characters")
			go FetchKillsForAllCharacters()
		} else {
			log.Println("Kill fetcher job is already running, skipping this scheduled run")
		}
	})
	c.Start()

	go TriggerImmediateKillFetch()
	go killFetcherWorker()
}

func TriggerImmediateKillFetch() {
	if !isJobRunning() {
		log.Println("Starting immediate kill fetch for all characters")
		FetchKillsForAllCharacters()
	} else {
		log.Println("Kill fetcher job is already running, skipping immediate fetch")
	}
}

func isJobRunning() bool {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	return jobRunning
}

func setJobRunning(running bool) {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	jobRunning = running
}

func killFetcherWorker() {
	for characterID := range fetchQueue {
		if !isJobRunning() {
			setJobRunning(true)
			FetchKillsForCharacter(characterID)
			setJobRunning(false)
		} else {
			log.Printf("Kill fetcher job is already running, queuing character %d for later", characterID)
			go func(id int64) {
				time.Sleep(5 * time.Minute)
				QueueCharacterForKillFetch(id)
			}(characterID)
		}
	}
}

func QueueCharacterForKillFetch(characterID int64) {
	select {
	case fetchQueue <- characterID:
		log.Printf("Queued character %d for kill fetching", characterID)
	default:
		log.Printf("Queue is full, character %d will not be processed this time", characterID)
	}
}

func FetchKillsForAllCharacters() {
	if !isJobRunning() {
		setJobRunning(true)
		defer setJobRunning(false)

		characters, err := queries.GetAllCharacters()
		if err != nil {
			log.Printf("Error fetching characters: %v", err)
			return
		}

		for _, character := range characters {
			QueueCharacterForKillFetch(character.ID)
		}
	} else {
		log.Println("Kill fetcher job is already running, skipping this run")
	}
}

func FetchKillsForCharacter(characterID int64) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in FetchKillsForCharacter: %v", r)
		}
		setJobRunning(false)
	}()

	log.Printf("Fetching kills for character %d", characterID)

	isInitialFetch, err := queries.IsInitialFetchForCharacter(characterID)
	if err != nil {
		log.Printf("Error checking initial fetch status for character %d: %v", characterID, err)
		return
	}

	if isInitialFetch {
		fetchAllKillsForCharacter(characterID)
	} else {
		fetchRecentKillsForCharacter(characterID)
	}

	enrichKillsForCharacter(characterID)
}

func fetchAllKillsForCharacter(characterID int64) {
	page := 1
	totalNewKills := 0

	for {
		zkillKills, err := services.FetchKillsFromZKillboard(characterID, page)
		if err != nil {
			log.Printf("Error fetching kills from zKillboard for character %d, page %d: %v", characterID, page, err)
			break
		}

		if len(zkillKills) == 0 {
			log.Printf("No more kills found for character %d", characterID)
			break
		}

		newKills := insertNewKills(characterID, zkillKills)
		totalNewKills += newKills

		log.Printf("Inserted %d new kills for character %d on page %d", newKills, characterID, page)

		if newKills < 200 {
			log.Printf("Less than 200 new kills on page %d for character %d, stopping", page, characterID)
			break
		}

		page++
		time.Sleep(1 * time.Second)
	}

	log.Printf("Finished initial kill fetch for character %d. Total new kills: %d", characterID, totalNewKills)
}

func fetchRecentKillsForCharacter(characterID int64) {
	zkillKills, err := services.FetchKillsFromZKillboard(characterID, 1)
	if err != nil {
		log.Printf("Error fetching recent kills from zKillboard for character %d: %v", characterID, err)
		return
	}

	newKills := insertNewKills(characterID, zkillKills)
	log.Printf("Inserted %d new kills for character %d from recent fetch", newKills, characterID)
}

func insertNewKills(characterID int64, zkillKills []models.ZKillKill) int {
	newKills := 0
	for _, zkill := range zkillKills {
		existingKill, err := queries.GetKillByID(zkill.KillmailID)
		if err == nil && existingKill != nil {
			continue
		}

		kill := models.Kill{
			KillmailID:     zkill.KillmailID,
			CharacterID:    characterID,
			LocationID:     zkill.ZKB.LocationID,
			Hash:           zkill.ZKB.Hash,
			FittedValue:    zkill.ZKB.FittedValue,
			DroppedValue:   zkill.ZKB.DroppedValue,
			DestroyedValue: zkill.ZKB.DestroyedValue,
			TotalValue:     zkill.ZKB.TotalValue,
			Points:         zkill.ZKB.Points,
			NPC:            zkill.ZKB.NPC,
			Solo:           zkill.ZKB.Solo,
			Awox:           zkill.ZKB.Awox,
		}

		err = queries.InsertKill(&kill)
		if err != nil {
			log.Printf("Error inserting kill %d: %v", kill.KillmailID, err)
		} else {
			newKills++
		}
	}
	return newKills
}

func enrichKillsForCharacter(characterID int64) {
	kills, err := queries.GetUnenrichedKillsForCharacter(characterID)
	if err != nil {
		log.Printf("Error fetching unenriched kills for character %d: %v", characterID, err)
		return
	}

	for _, kill := range kills {
		esiKill, err := services.FetchKillmailFromESI(kill.KillmailID, kill.Hash)
		if err != nil {
			log.Printf("Error fetching killmail %d from ESI: %v", kill.KillmailID, err)
			continue
		}

		kill.KillTime = esiKill.KillTime
		kill.SolarSystemID = esiKill.SolarSystemID
		kill.Victim = esiKill.Victim
		kill.Attackers = esiKill.Attackers

		err = queries.UpdateKill(&kill)
		if err != nil {
			log.Printf("Error updating kill %d with ESI data: %v", kill.KillmailID, err)
		} else {
			log.Printf("Successfully enriched kill %d", kill.KillmailID)
		}

		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("Finished enriching kills for character %d", characterID)
}

func FetchAllKillsForCharacter(characterID int64) {
	log.Printf("Starting full kill fetch for character %d", characterID)
	FetchKillsForCharacter(characterID)
	log.Printf("Finished full kill fetch for character %d", characterID)
}
