package jobs

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
)

var (
	killFetchQueue       = make(chan killFetchJob, 1000)
	failedKillFetchQueue = make(chan killFetchJob, 1000)
	workerPool           = make(chan struct{}, 10) // Limit to 10 concurrent workers
	esiErrorLimitMutex   sync.Mutex
	esiErrorLimitBackoff time.Time
)

type killFetchJob struct {
	characterID  int64
	lastKillTime time.Time
	isInitial    bool
}

func StartKillFetcherWorker() {
	go killFetcherWorker()
	go failedKillFetchWorker()
}

func QueueKillFetch(characterID int64, lastKillTime time.Time, isInitial bool) {
	select {
	case killFetchQueue <- killFetchJob{characterID: characterID, lastKillTime: lastKillTime, isInitial: isInitial}:
		log.Printf("Queued character %d for kill fetching (initial: %v)", characterID, isInitial)
	default:
		log.Printf("Kill fetch queue is full, character %d will not be processed this time", characterID)
	}
}

func killFetcherWorker() {
	for job := range killFetchQueue {
		workerPool <- struct{}{} // Acquire a worker
		go func(j killFetchJob) {
			defer func() { <-workerPool }() // Release the worker when done

			var newKills int
			var err error
			if j.isInitial {
				newKills, err = FetchKillsForCharacter(j.characterID)
			} else {
				newKills, err = FetchNewKillsForCharacter(j.characterID, j.lastKillTime)
			}

			if err != nil {
				log.Printf("Error fetching kills for character %d: %v", j.characterID, err)
				failedKillFetchQueue <- j
			} else {
				log.Printf("Fetched %d kills for character %d (initial: %v)", newKills, j.characterID, j.isInitial)
			}
		}(job)
	}
}

func failedKillFetchWorker() {
	for job := range failedKillFetchQueue {
		time.Sleep(5 * time.Minute) // Wait before retrying
		QueueKillFetch(job.characterID, job.lastKillTime, job.isInitial)
	}
}

func FetchKillsForCharacter(characterID int64) (int, error) {
	return FetchNewKillsForCharacter(characterID, time.Time{})
}

func FetchNewKillsForCharacter(characterID int64, lastKillTime time.Time) (int, error) {
	log.Printf("Starting to fetch new kills for character %d since %v", characterID, lastKillTime)
	totalNewKills := 0
	page := 1
	pageStaggerInterval := 500 * time.Millisecond

	for {
		log.Printf("Fetching page %d for character %d", page, characterID)
		zkillKills, err := services.FetchKillsFromZKillboard(characterID, page)
		if err != nil {
			return totalNewKills, fmt.Errorf("error fetching kills from zKillboard: %v", err)
		}

		if len(zkillKills) == 0 {
			log.Printf("No more kills found for character %d", characterID)
			break
		}

		newKills, allExist, err := processKillsPage(characterID, zkillKills)
		if err != nil {
			log.Printf("Error processing kills page for character %d: %v", characterID, err)
			page++
			time.Sleep(pageStaggerInterval)
			continue
		}

		totalNewKills += newKills
		log.Printf("Processed %d new kills on page %d for character %d", newKills, page, characterID)

		if allExist {
			log.Printf("All kills on page %d already exist for character %d, stopping fetch", page, characterID)
			break
		}

		page++
		time.Sleep(pageStaggerInterval)
	}

	log.Printf("Finished fetching new kills for character %d. Total new kills: %d", characterID, totalNewKills)
	return totalNewKills, nil
}

func processKillsPage(characterID int64, zkillKills []models.ZKillKill) (int, bool, error) {
	killmailIDs := make([]int64, len(zkillKills))
	for i, kill := range zkillKills {
		killmailIDs[i] = kill.KillmailID
	}

	existingKills, err := queries.GetKillsByKillmailIDs(killmailIDs)
	if err != nil {
		return 0, false, fmt.Errorf("error checking existing kills: %v", err)
	}

	existingKillMap := make(map[int64]bool)
	for _, kill := range existingKills {
		existingKillMap[kill.KillmailID] = true
	}

	var newKills []models.Kill
	for _, zkillKill := range zkillKills {
		if existingKillMap[zkillKill.KillmailID] {
			continue
		}

		esiKill, err := fetchKillmailWithExponentialBackoff(zkillKill.KillmailID, zkillKill.ZKB.Hash)
		if err != nil {
			log.Printf("Error fetching killmail %d from ESI after retries: %v", zkillKill.KillmailID, err)
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
			Attackers:      models.AttackersJSON(esiKill.Attackers),
		}

		newKills = append(newKills, kill)
	}

	if len(newKills) > 0 {
		err := queries.BulkUpsertKills(newKills)
		if err != nil {
			return 0, false, fmt.Errorf("error bulk upserting kills: %v", err)
		}
		log.Printf("Successfully upserted %d new kills for character %d", len(newKills), characterID)
	} else {
		log.Printf("No new kills found for character %d in this batch", characterID)
	}

	return len(newKills), len(newKills) == 0, nil
}

func fetchKillmailWithExponentialBackoff(killmailID int64, hash string) (*models.Kill, error) {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 5 * time.Minute

	var esiKill *models.Kill
	var err error

	operation := func() error {
		esiKill, err = services.FetchKillmailFromESI(killmailID, hash)
		if err != nil {
			if services.IsESIErrorLimit(err) {
				log.Printf("ESI error limit reached, backing off")
				return err
			}
			log.Printf("Error fetching killmail %d from ESI: %v", killmailID, err)
			return err
		}
		return nil
	}

	err = backoff.Retry(operation, b)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch killmail after multiple attempts: %v", err)
	}

	return esiKill, nil
}

// Add these functions at the end of the file

func FetchAllKillsForCharacter(characterID int64) {
	QueueKillFetch(characterID, time.Time{}, true)
}

func FetchKillsForAllCharacters() {
	characters, err := queries.GetAllCharacters()
	if err != nil {
		log.Printf("Error fetching characters: %v", err)
		return
	}

	for _, character := range characters {
		QueueKillFetch(character.ID, time.Time{}, true)
	}
}
