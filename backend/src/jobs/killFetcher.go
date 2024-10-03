package jobs

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
	"gorm.io/gorm"
)

var (
	killFetchQueue       = make(chan killFetchJob, 1000) // Increase the queue size
	backoffMutex         sync.Mutex
	lastBackoffTime      time.Time
	workerRunning        bool
	workerMutex          sync.Mutex
	esiErrorLimitBackoff time.Time
	esiErrorLimitMutex   sync.Mutex
)

type killFetchJob struct {
	characterID  int64
	lastKillTime time.Time
	isInitial    bool
}

const (
	maxRetries = 3
	retryDelay = 5 * time.Second
)

func StartKillFetcherWorker() {
	workerMutex.Lock()
	if !workerRunning {
		workerRunning = true
		workerMutex.Unlock()
		go killFetcherWorker()
	} else {
		workerMutex.Unlock()
	}
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
		esiErrorLimitMutex.Lock()
		if time.Now().Before(esiErrorLimitBackoff) {
			sleepTime := time.Until(esiErrorLimitBackoff)
			esiErrorLimitMutex.Unlock()
			log.Printf("Waiting for ESI error limit backoff: %v", sleepTime)
			time.Sleep(sleepTime)
		} else {
			esiErrorLimitMutex.Unlock()
		}

		var newKills int
		var err error
		if job.isInitial {
			newKills, err = FetchKillsForCharacter(job.characterID)
		} else {
			newKills, err = FetchNewKillsForCharacter(job.characterID, job.lastKillTime)
		}

		if err != nil {
			log.Printf("Error fetching kills for character %d: %v", job.characterID, err)
		} else {
			log.Printf("Fetched %d kills for character %d (initial: %v)", newKills, job.characterID, job.isInitial)
		}
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
			// Continue to the next page instead of stopping
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
	var kills []models.Kill
	newKills := 0
	allExist := true

	for _, zkillKill := range zkillKills {
		existingKill, err := queries.GetKillByKillmailID(zkillKill.KillmailID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error checking existing kill %d: %v", zkillKill.KillmailID, err)
			continue
		}

		if existingKill != nil {
			log.Printf("Kill %d already exists, skipping", zkillKill.KillmailID)
			continue
		}

		allExist = false

		esiKill, err := fetchKillmailWithRetry(zkillKill.KillmailID, zkillKill.ZKB.Hash)
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

		kills = append(kills, kill)
		newKills++
	}

	if len(kills) > 0 {
		err := queries.BulkUpsertKills(kills)
		if err != nil {
			return newKills, false, fmt.Errorf("error bulk upserting kills: %v", err)
		}
		log.Printf("Successfully upserted %d new kills for character %d", len(kills), characterID)
	} else {
		log.Printf("No new kills found for character %d in this batch", characterID)
	}

	return newKills, allExist, nil
}

func fetchKillmailWithBackoff(killmailID int64, hash string) (*models.Kill, error) {
	esiErrorLimitMutex.Lock()
	if time.Now().Before(esiErrorLimitBackoff) {
		sleepTime := time.Until(esiErrorLimitBackoff)
		esiErrorLimitMutex.Unlock()
		log.Printf("Waiting for ESI error limit backoff: %v", sleepTime)
		time.Sleep(sleepTime)
	} else {
		esiErrorLimitMutex.Unlock()
	}

	backoffMutex.Lock()
	if time.Since(lastBackoffTime) < 15*time.Second {
		sleepTime := 15*time.Second - time.Since(lastBackoffTime)
		backoffMutex.Unlock()
		time.Sleep(sleepTime)
	} else {
		backoffMutex.Unlock()
	}

	esiKill, err := services.FetchKillmailFromESI(killmailID, hash)
	if err != nil {
		if services.IsESITimeout(err) {
			backoffMutex.Lock()
			lastBackoffTime = time.Now()
			backoffMutex.Unlock()
			log.Printf("ESI timeout encountered, backing off for 15 seconds")
			time.Sleep(15 * time.Second)
			return fetchKillmailWithBackoff(killmailID, hash) // Retry after backoff
		} else if services.IsESIErrorLimit(err) {
			esiErrorLimitMutex.Lock()
			esiErrorLimitBackoff = time.Now().Add(1 * time.Minute)
			esiErrorLimitMutex.Unlock()
			log.Printf("ESI error limit reached, backing off for 1 minute")
			time.Sleep(1 * time.Minute)
			return fetchKillmailWithBackoff(killmailID, hash) // Retry after waiting
		}
		return nil, err
	}
	return esiKill, nil
}

func fetchKillmailWithRetry(killmailID int64, hash string) (*models.Kill, error) {
	var esiKill *models.Kill
	var err error

	for i := 0; i < maxRetries; i++ {
		esiKill, err = services.FetchKillmailFromESI(killmailID, hash)
		if err == nil {
			return esiKill, nil
		}

		if services.IsESIErrorLimit(err) {
			log.Printf("ESI error limit reached, backing off for 1 minute")
			time.Sleep(1 * time.Minute)
		} else {
			log.Printf("Error fetching killmail %d from ESI (attempt %d/%d): %v", killmailID, i+1, maxRetries, err)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to fetch killmail after %d attempts: %v", maxRetries, err)
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
