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
)

var (
	killFetchQueue       = make(chan killFetchJob, 100)
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
	totalNewKills := int32(0)
	maxPageConcurrency := 10
	maxKillConcurrency := 1
	batchSize := 1
	pageStaggerInterval := 500 * time.Millisecond

	if lastKillTime.IsZero() {
		lastKillTime = time.Date(2003, 5, 6, 0, 0, 0, 0, time.UTC)
		log.Printf("No last kill time provided, fetching all kills since EVE Online release")
	}

	var pageWg sync.WaitGroup
	pageSem := make(chan struct{}, maxPageConcurrency)
	killSem := make(chan struct{}, maxKillConcurrency)
	errorChan := make(chan error, maxPageConcurrency)

	page := 1
	continueFetching := true
	oldestKillTime := time.Now()

	for continueFetching {
		pageWg.Add(1)
		pageSem <- struct{}{}
		go func(currentPage int) {
			defer pageWg.Done()
			defer func() { <-pageSem }()

			time.Sleep(time.Duration(currentPage-1) * pageStaggerInterval)

			newKills, pageOldestKillTime, reachedOldKills, err := fetchKillsForPage(characterID, currentPage, lastKillTime, killSem, batchSize)
			if err != nil {
				if err == ErrNoMoreKills {
					log.Printf("No more kills found for character %d on page %d", characterID, currentPage)
					continueFetching = false
					return
				}
				log.Printf("Error fetching kills for character %d on page %d: %v", characterID, currentPage, err)
				errorChan <- err
				return
			}
			atomic.AddInt32(&totalNewKills, int32(newKills))

			if pageOldestKillTime.Before(oldestKillTime) {
				oldestKillTime = pageOldestKillTime
			}

			if newKills == 0 {
				log.Printf("No new kills found for character %d on page %d, stopping fetch", characterID, currentPage)
				continueFetching = false
				return
			}

			if reachedOldKills {
				log.Printf("Reached kills older than last known kill for character %d, stopping fetch", characterID)
				continueFetching = false
				return
			}
		}(page)

		pageWg.Wait()

		if len(errorChan) > 0 {
			return int(totalNewKills), <-errorChan
		}

		if continueFetching {
			page++
		}
	}

	log.Printf("Finished fetching new kills for character %d. Total new kills: %d", characterID, totalNewKills)
	return int(totalNewKills), nil
}

var ErrNoMoreKills = errors.New("no more kills")

func fetchKillsForPage(characterID int64, page int, lastKillTime time.Time, sem chan struct{}, batchSize int) (int, time.Time, bool, error) {
	log.Printf("Fetching page %d of kills for character %d", page, characterID)
	zkillKills, err := services.FetchKillsFromZKillboard(characterID, page)
	if err != nil {
		return 0, time.Time{}, false, err
	}

	if len(zkillKills) == 0 {
		return 0, time.Time{}, false, ErrNoMoreKills
	}

	log.Printf("Processing %d kills for character %d on page %d", len(zkillKills), characterID, page)

	var wg sync.WaitGroup
	newKillsInPage := int32(0)
	oldestKillTime := time.Now()
	reachedOldKills := false

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

			newKillsInBatch, batchOldestKillTime, batchReachedOldKills := processNewKillsBatch(batch, characterID, lastKillTime)
			atomic.AddInt32(&newKillsInPage, int32(newKillsInBatch))
			if batchOldestKillTime.Before(oldestKillTime) {
				oldestKillTime = batchOldestKillTime
			}
			if batchReachedOldKills {
				reachedOldKills = true
			}
		}(zkillKills[i:end])
	}

	wg.Wait()
	log.Printf("Processed %d new kills for character %d on page %d", newKillsInPage, characterID, page)

	return int(newKillsInPage), oldestKillTime, reachedOldKills, nil
}

func processNewKillsBatch(batch []models.ZKillKill, characterID int64, lastKillTime time.Time) (int, time.Time, bool) {
	var kills []models.Kill
	newKills := 0
	oldestKillTime := time.Now()
	reachedOldKills := false

	for _, zkillKill := range batch {
		existingKill, err := queries.GetKillByKillmailID(zkillKill.KillmailID)
		if err != nil {
			log.Printf("Error checking existing kill %d: %v", zkillKill.KillmailID, err)
			continue
		}

		if existingKill != nil {
			continue
		}

		esiKill, err := fetchKillmailWithBackoff(zkillKill.KillmailID, zkillKill.ZKB.Hash)
		if err != nil {
			log.Printf("Error fetching killmail %d from ESI: %v", zkillKill.KillmailID, err)
			continue
		}

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

			if esiKill.KillTime.Before(oldestKillTime) {
				oldestKillTime = esiKill.KillTime
			}
		} else {
			reachedOldKills = true
			log.Printf("Reached kill %d which is older than or equal to last known kill, will stop after processing current batch", zkillKill.KillmailID)
			break
		}
	}

	if len(kills) > 0 {
		err := queries.BulkUpsertKills(kills)
		if err != nil {
			log.Printf("Error bulk upserting kills: %v", err)
			return 0, oldestKillTime, reachedOldKills
		}
		log.Printf("Successfully upserted %d new kills for character %d", len(kills), characterID)
	} else {
		log.Printf("No new kills found for character %d in this batch", characterID)
	}

	return newKills, oldestKillTime, reachedOldKills
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
