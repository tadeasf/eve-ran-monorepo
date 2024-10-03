package jobs

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

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
	QueueCharacterForKillFetch(characterID)
}

func FetchNewKillsForCharacter(characterID int64, lastKillTime time.Time) (int, error) {
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

				processNewKillsBatch(batch, characterID, lastKillTime, &totalNewKills)
			}(zkillKills[i:end])
		}

		wg.Wait()
		log.Printf("Processed %d new kills for character %d on page %d", atomic.LoadInt32(&totalNewKills), characterID, page)

		if len(zkillKills) < 200 {
			break
		}

		page++
	}

	return int(totalNewKills), nil
}

func processNewKillsBatch(batch []models.ZKillKill, characterID int64, lastKillTime time.Time, newKills *int32) {
	var kills []models.Kill
	for _, zkillKill := range batch {
		existingKill, err := queries.GetKillByKillmailID(zkillKill.KillmailID)
		if err != nil {
			log.Printf("Error checking existing kill %d: %v", zkillKill.KillmailID, err)
			continue
		}

		if existingKill != nil && !existingKill.KillTime.IsZero() {
			continue
		}

		esiKill, err := services.FetchKillmailFromESI(zkillKill.KillmailID, zkillKill.ZKB.Hash)
		if err != nil {
			log.Printf("Error fetching killmail %d from ESI: %v", zkillKill.KillmailID, err)
			continue
		}

		if esiKill.KillTime.Before(lastKillTime) {
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

	if len(kills) > 0 {
		err := queries.BulkUpsertKills(kills)
		if err != nil {
			log.Printf("Error bulk upserting kills: %v", err)
		}
	}
}
