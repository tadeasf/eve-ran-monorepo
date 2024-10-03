package jobs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/services"
)

var (
	fetchQueue = make(chan int64, 100)
	jobMutex   sync.Mutex
	jobRunning bool
)

func StartKillFetcherJob() {
	c := cron.New()
	c.AddFunc("@every 10min", func() {
		if !isJobRunning() {
			log.Println("Starting to fetch kills for all characters")
			go FetchKillsForAllCharacters()
		} else {
			log.Println("Kill fetcher job is already running, skipping this run")
		}
	})
	c.Start()

	go killFetcherWorker()
}

// Add this new function to allow immediate execution
func TriggerImmediateKillFetch() {
	if !isJobRunning() {
		log.Println("Starting immediate kill fetch for all characters")
		go FetchKillsForAllCharacters()
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
				time.Sleep(5 * time.Minute) // Wait for 5 minutes before re-queuing
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

		characters, err := db.GetAllCharacters()
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

const (
	initialConcurrency = 2
	maxConcurrency     = 200
	minConcurrency     = 2
)

func fetchKillsForCharacter(characterID int64) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in fetchKillsForCharacter: %v", r)
		}
		setJobRunning(false)
	}()

	lastKillTime, err := db.GetLastKillTimeForCharacter(characterID)
	if err != nil {
		log.Printf("Error getting last kill time for character %d: %v", characterID, err)
		lastKillTime = time.Time{}
	}
	log.Printf("Last kill time for character %d: %v", characterID, lastKillTime)
	isNewCharacter := lastKillTime.IsZero()
	page := 1
	totalNewKills := 0

	concurrency := initialConcurrency
	semaphore := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	killChan := make(chan *models.Kill, concurrency)
	done := make(chan bool)
	stopProcessing := make(chan bool)

	// Batch processing goroutine
	go func() {
		batch := make([]*models.Kill, 0, concurrency)
		for kill := range killChan {
			batch = append(batch, kill)
			if len(batch) == concurrency {
				err := db.UpsertKillsBatch(batch)
				if err != nil {
					log.Printf("Error upserting kill batch: %v", err)
				} else {
					log.Printf("Successfully committed batch of %d kills to the database", len(batch))
					totalNewKills += len(batch)
				}
				batch = batch[:0] // Clear the batch
			}
		}
		// Process any remaining kills in the batch
		if len(batch) > 0 {
			err := db.UpsertKillsBatch(batch)
			if err != nil {
				log.Printf("Error upserting final kill batch: %v", err)
			} else {
				log.Printf("Successfully committed final batch of %d kills to the database", len(batch))
				totalNewKills += len(batch)
			}
		}
		done <- true
	}()

outerLoop:
	for {
		log.Printf("Fetching page %d for character %d with concurrency %d", page, characterID, concurrency)
		kills, err := services.FetchKillsFromZKillboard(characterID, page)
		if err != nil {
			log.Printf("Error fetching kills for character %d: %v", characterID, err)
			concurrency = max(minConcurrency, concurrency/2)
			semaphore = make(chan struct{}, concurrency)
			continue
		}

		if len(kills) == 0 {
			log.Printf("No more kills found for character %d", characterID)
			break
		}

		var newKills int32
		var errors int32
		for _, kill := range kills {
			select {
			case <-stopProcessing:
				break outerLoop
			default:
				wg.Add(1)
				go func(k models.Kill) {
					defer wg.Done()
					semaphore <- struct{}{}
					defer func() { <-semaphore }()

					esiKill, err := services.FetchKillmailFromESI(k.KillmailID, k.Hash)
					if err != nil {
						log.Printf("Error fetching ESI killmail %d: %v", k.KillmailID, err)
						atomic.AddInt32(&errors, 1)
						return
					}

					k.KillTime = esiKill.KillTime
					k.SolarSystemID = esiKill.SolarSystemID
					k.Victim = esiKill.Victim
					k.Attackers = esiKill.Attackers

					if isNewCharacter || k.KillTime.After(lastKillTime) {
						atomic.AddInt32(&newKills, 1)
						killChan <- &k
					} else {
						log.Printf("Reached already processed kills for character %d", characterID)
						stopProcessing <- true
					}
				}(kill)
			}
		}

		wg.Wait()

		log.Printf("Processed %d new kills for character %d on page %d", newKills, characterID, page)

		if errors > 0 {
			concurrency = max(minConcurrency, concurrency/2)
		} else {
			concurrency = min(maxConcurrency, concurrency*2)
		}
		semaphore = make(chan struct{}, concurrency)

		if newKills == 0 && !isNewCharacter {
			log.Printf("No new kills on page %d for character %d, stopping", page, characterID)
			break
		}

		page++
		time.Sleep(1 * time.Second)
	}

	close(killChan)
	<-done

	log.Printf("Finished fetching kills for character %d. Total new kills: %d", characterID, totalNewKills)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func FetchAllKillsForCharacter(characterID int64) {
	log.Printf("Starting full kill fetch for character %d", characterID)
	fetchKillsForCharacter(characterID)
	log.Printf("Finished full kill fetch for character %d", characterID)
}

func FetchKillsForCharacter(characterID int64) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in FetchKillsForCharacter: %v", r)
		}
		setJobRunning(false)
	}()

	lastKillTime, err := db.GetLastKillTimeForCharacter(characterID)
	if err != nil {
		log.Printf("Error getting last kill time for character %d: %v", characterID, err)
		lastKillTime = time.Time{}
	}
	log.Printf("Last kill time for character %d: %v", characterID, lastKillTime)

	page := 1
	totalNewKills := 0
	isNewCharacter := lastKillTime.IsZero()

	for {
		log.Printf("Fetching page %d for character %d", page, characterID)
		url := fmt.Sprintf("https://zkillboard.com/api/characterID/%d/page/%d/", characterID, page)

		var zkillResponse []struct {
			KillmailID   int64  `json:"killmail_id"`
			KillmailTime string `json:"killmail_time"`
			ZKB          struct {
				LocationID     int64   `json:"locationID"`
				Hash           string  `json:"hash"`
				FittedValue    float64 `json:"fittedValue"`
				DroppedValue   float64 `json:"droppedValue"`
				DestroyedValue float64 `json:"destroyedValue"`
				TotalValue     float64 `json:"totalValue"`
				Points         int     `json:"points"`
				NPC            bool    `json:"npc"`
				Solo           bool    `json:"solo"`
				Awox           bool    `json:"awox"`
			} `json:"zkb"`
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching kills for character %d: %v", characterID, err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}

		err = json.Unmarshal(body, &zkillResponse)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			return
		}

		newKills := 0
		for _, zkill := range zkillResponse {
			killTime, err := time.Parse("2006-01-02T15:04:05Z", zkill.KillmailTime)
			if err != nil {
				log.Printf("Error parsing kill time for killmail %d: %v", zkill.KillmailID, err)
				continue
			}

			if !isNewCharacter && killTime.Before(lastKillTime) || killTime.Equal(lastKillTime) {
				log.Printf("Reached already processed kills for character %d, stopping", characterID)
				return
			}

			kill := models.Kill{
				KillmailID:     zkill.KillmailID,
				CharacterID:    characterID,
				KillTime:       killTime,
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

			// Fetch additional data from ESI
			esiKill, err := services.FetchKillmailFromESI(kill.KillmailID, kill.Hash)
			if err != nil {
				log.Printf("Error fetching killmail %d from ESI: %v", kill.KillmailID, err)
				continue
			}

			// Combine zKillboard and ESI data
			kill.SolarSystemID = esiKill.SolarSystemID
			kill.Victim = esiKill.Victim
			kill.Attackers = esiKill.Attackers

			// Insert the kill
			err = db.InsertKill(&kill)
			if err != nil {
				log.Printf("Error inserting kill %d: %v", kill.KillmailID, err)
			} else {
				newKills++
				totalNewKills++
			}
		}

		log.Printf("Inserted %d new kills for character %d on page %d", newKills, characterID, page)

		if newKills == 0 {
			log.Printf("No new kills on page %d for character %d, stopping", page, characterID)
			break
		}

		page++
		time.Sleep(1 * time.Second) // Add a delay to avoid hitting rate limits
	}

	log.Printf("Finished fetching kills for character %d. Total new kills: %d", characterID, totalNewKills)
}
