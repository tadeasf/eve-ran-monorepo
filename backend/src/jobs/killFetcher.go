package jobs

import (
	"log"
	"sync"
	"sync/atomic"
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

func FetchKillsForCharacter(characterID int64) (int, error) {
	totalNewKills := int32(0)
	page := 1
	batchSize := 100
	maxConcurrent := 10

	for {
		zkillKills, err := services.FetchKillsFromZKillboard(characterID, page)
		if err != nil {
			return int(totalNewKills), err
		}

		if len(zkillKills) == 0 {
			break
		}

		sem := make(chan bool, maxConcurrent)
		var wg sync.WaitGroup

		for i := 0; i < len(zkillKills); i += batchSize {
			end := i + batchSize
			if end > len(zkillKills) {
				end = len(zkillKills)
			}

			wg.Add(1)
			sem <- true
			go func(batch []models.ZKillKill) {
				defer wg.Done()
				defer func() { <-sem }()

				processBatch(batch, characterID, &totalNewKills)
			}(zkillKills[i:end])
		}

		wg.Wait()
		log.Printf("Processed %d new kills for character %d on page %d", atomic.LoadInt32(&totalNewKills), characterID, page)

		if len(zkillKills) < 200 {
			log.Printf("Less than 200 new kills on page %d for character %d, stopping", page, characterID)
			break
		}

		page++
	}

	return int(totalNewKills), nil
}

func processBatch(batch []models.ZKillKill, characterID int64, newKills *int32) {
	var kills []models.Kill
	for _, zkillKill := range batch {
		exists, err := queries.KillExists(zkillKill.KillmailID)
		if err != nil {
			log.Printf("Error checking if kill %d exists: %v", zkillKill.KillmailID, err)
			continue
		}

		if !exists {
			esiKill, err := services.FetchKillmailFromESI(zkillKill.KillmailID, zkillKill.ZKB.Hash)
			if err != nil {
				log.Printf("Error fetching killmail %d from ESI: %v", zkillKill.KillmailID, err)
				continue
			}

			kill := models.Kill{
				KillmailID:     zkillKill.KillmailID,
				CharacterID:    characterID,
				KillTime:       esiKill.KillTime,
				SolarSystemID:  esiKill.SolarSystemID,
				LocationID:     zkillKill.ZKB.LocationID,
				Hash:           zkillKill.ZKB.Hash,
				FittedValue:    zkillKill.ZKB.FittedValue,
				DroppedValue:   zkillKill.ZKB.DroppedValue,
				DestroyedValue: zkillKill.ZKB.DestroyedValue,
				TotalValue:     zkillKill.ZKB.TotalValue,
				Points:         zkillKill.ZKB.Points,
				NPC:            zkillKill.ZKB.NPC,
				Solo:           zkillKill.ZKB.Solo,
				Awox:           zkillKill.ZKB.Awox,
				Victim:         esiKill.Victim,
				Attackers:      esiKill.Attackers,
			}

			kills = append(kills, kill)
			atomic.AddInt32(newKills, 1)
		}
	}

	if len(kills) > 0 {
		err := queries.BulkUpsertKills(kills)
		if err != nil {
			log.Printf("Error bulk upserting kills: %v", err)
		}
	}
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

		err = queries.UpsertKill(&kill)
		if err != nil {
			log.Printf("Error inserting kill %d: %v", kill.KillmailID, err)
			return 0
		}
		newKills++
	}
	return newKills
}

func enrichKillsForCharacter(characterID int64) {
	log.Printf("Starting to enrich kills for character %d", characterID)
	kills, err := queries.GetUnenrichedKillsForCharacter(characterID)
	if err != nil {
		log.Printf("Error fetching unenriched kills for character %d: %v", characterID, err)
		return
	}

	for _, kill := range kills {
		existingKill, err := queries.GetKillByKillmailID(kill.KillmailID)
		if err != nil {
			log.Printf("Error checking existing kill %d: %v", kill.KillmailID, err)
			continue
		}
		if existingKill == nil {
			log.Printf("Kill %d not found in database, skipping", kill.KillmailID)
			continue
		}

		esiKill, err := services.FetchKillmailFromESI(kill.KillmailID, kill.Hash)
		if err != nil {
			log.Printf("Error fetching killmail %d from ESI: %v", kill.KillmailID, err)
			continue
		}

		existingKill.KillTime = esiKill.KillTime
		existingKill.SolarSystemID = esiKill.SolarSystemID
		existingKill.Victim = esiKill.Victim
		existingKill.Attackers = esiKill.Attackers

		err = queries.UpsertKill(existingKill)
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
