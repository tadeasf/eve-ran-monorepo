package jobs

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
	"gorm.io/gorm"
)

var (
	fetchQueue = make(chan int64, 100)
	jobMutex   sync.Mutex
	jobRunning bool
)

func StartKillFetcherJob() {
	go killFetcherWorker()
}

func QueueCharacterForKillFetch(characterID int64) {
	select {
	case fetchQueue <- characterID:
		log.Printf("Queued character %d for kill fetching", characterID)
	default:
		log.Printf("Queue is full, character %d will not be processed this time", characterID)
	}
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

func FetchKillsForAllCharacters() {
	characters, err := queries.GetAllCharacters()
	if err != nil {
		log.Printf("Error fetching characters: %v", err)
		return
	}

	for _, character := range characters {
		QueueCharacterForKillFetch(character.ID)
	}
}

func FetchAllKillsForCharacter(characterID int64) {
	log.Printf("Starting to fetch all kills for character %d", characterID)
	lastKillTime := time.Time{} // This will fetch all kills
	newKills, err := FetchNewKillsForCharacter(characterID, lastKillTime)
	if err != nil {
		log.Printf("Error fetching kills for character %d: %v", characterID, err)
	} else {
		log.Printf("Fetched %d new kills for character %d", newKills, characterID)
	}
}

func FetchNewKillsForCharacter(characterID int64, lastKillTime time.Time) (int, error) {
	log.Printf("Starting to fetch new kills for character %d since %v", characterID, lastKillTime)
	totalNewKills := int32(0)
	maxPageConcurrency := 10
	maxKillConcurrency := 200
	batchSize := 1

	// If lastKillTime is zero, set it to a very old date to fetch all kills
	if lastKillTime.IsZero() {
		lastKillTime = time.Date(2003, 5, 6, 0, 0, 0, 0, time.UTC) // EVE Online release date
		log.Printf("No last kill time provided, fetching all kills since EVE Online release")
	}

	var pageWg sync.WaitGroup
	pageSem := make(chan struct{}, maxPageConcurrency)
	killSem := make(chan struct{}, maxKillConcurrency)

	for page := 1; page <= maxPageConcurrency; page++ {
		pageWg.Add(1)
		pageSem <- struct{}{}
		go func(currentPage int) {
			defer pageWg.Done()
			defer func() { <-pageSem }()

			newKills, err := fetchKillsForPage(characterID, currentPage, lastKillTime, killSem, batchSize)
			if err != nil {
				if err == ErrNoMoreKills {
					log.Printf("No more kills found for character %d on page %d", characterID, currentPage)
					return
				}
				log.Printf("Error fetching kills for character %d on page %d: %v", characterID, currentPage, err)
				return
			}
			atomic.AddInt32(&totalNewKills, int32(newKills))
		}(page)
	}

	pageWg.Wait()

	log.Printf("Finished fetching new kills for character %d. Total new kills: %d", characterID, totalNewKills)
	return int(totalNewKills), nil
}

var ErrNoMoreKills = errors.New("no more kills")

func fetchKillsForPage(characterID int64, page int, lastKillTime time.Time, sem chan struct{}, batchSize int) (int, error) {
	log.Printf("Fetching page %d of kills for character %d", page, characterID)
	zkillKills, err := services.FetchKillsFromZKillboard(characterID, page)
	if err != nil {
		return 0, err
	}

	if len(zkillKills) == 0 {
		return 0, ErrNoMoreKills
	}

	log.Printf("Processing %d kills for character %d on page %d", len(zkillKills), characterID, page)

	var wg sync.WaitGroup
	newKillsInPage := int32(0)

	for i := 0; i < len(zkillKills); i += batchSize {
		end := i + batchSize
		if end > len(zkillKills) {
			end = len(zkillKills)
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(batch []models.ZKillKill) {
			defer wg.Done()
			defer func() { <-sem }()

			newKillsInBatch, err := processNewKillsBatch(batch, characterID, lastKillTime)
			if err != nil {
				log.Printf("Error processing batch for character %d on page %d: %v", characterID, page, err)
			} else {
				atomic.AddInt32(&newKillsInPage, int32(newKillsInBatch))
			}
		}(zkillKills[i:end])
	}

	wg.Wait()
	log.Printf("Processed %d new kills for character %d on page %d", newKillsInPage, characterID, page)

	return int(newKillsInPage), nil
}

func processNewKillsBatch(batch []models.ZKillKill, characterID int64, lastKillTime time.Time) (int, error) {
	var kills []models.Kill
	newKills := 0

	for _, zkillKill := range batch {
		existingKill, err := queries.GetKillByKillmailID(zkillKill.KillmailID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error checking existing kill %d: %v", zkillKill.KillmailID, err)
			continue
		}

		if existingKill != nil {
			continue
		}

		esiKill, err := services.FetchKillmailFromESI(zkillKill.KillmailID, zkillKill.ZKB.Hash)
		if err != nil {
			log.Printf("Error fetching killmail %d from ESI: %v", zkillKill.KillmailID, err)
			continue
		}

		// We'll always process the kill if lastKillTime is set to the EVE Online release date
		if esiKill.KillTime.After(lastKillTime) {
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
				Attackers:      models.AttackersJSON(esiKill.Attackers),
			}

			kills = append(kills, kill)
			newKills++
		}
	}

	if len(kills) > 0 {
		err := queries.BulkUpsertKills(kills)
		if err != nil {
			log.Printf("Error bulk upserting kills: %v", err)
			return 0, err
		}
		log.Printf("Successfully upserted %d kills for character %d", len(kills), characterID)
	} else {
		log.Printf("No new kills to upsert for character %d", characterID)
	}

	return newKills, nil
}
