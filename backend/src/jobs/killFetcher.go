package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/services"
	"gorm.io/gorm"
)

var (
	fetchQueue = make(chan int64, 100)
	jobMutex   sync.Mutex
	jobRunning bool
)

type zKillKill struct {
	KillmailID int64 `json:"killmail_id"`
	ZKB        struct {
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
	page := 1
	totalNewKills := 0

	for {
		log.Printf("Fetching page %d for character %d", page, characterID)
		zkillKills, err := services.FetchKillsFromZKillboard(characterID, page)
		if err != nil {
			log.Printf("Error fetching kills from zKillboard for character %d: %v", characterID, err)
			break
		}

		if len(zkillKills) == 0 {
			log.Printf("No more kills found for character %d", characterID)
			break
		}

		for _, zkill := range zkillKills {
			// Check if kill already exists in database
			existingKill, err := db.GetKillByID(zkill.KillmailID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Error checking for existing kill %d: %v", zkill.KillmailID, err)
				continue
			}

			if existingKill != nil {
				// Kill already exists, skip
				continue
			}

			// Create new kill with basic zKillboard data
			kill := &models.Kill{
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

			// Insert basic kill data
			err = db.InsertKill(kill)
			if err != nil {
				log.Printf("Error inserting basic kill data for killmail %d: %v", zkill.KillmailID, err)
				continue
			}

			// Queue kill for ESI enrichment
			go enrichKillWithESIData(kill)

			totalNewKills++
		}

		page++
		time.Sleep(1 * time.Second) // Rate limiting
	}

	log.Printf("Finished fetching kills for character %d. Total new kills: %d", characterID, totalNewKills)
}

func enrichKillWithESIData(kill *models.Kill) {
	esiKill, err := services.FetchKillmailFromESI(kill.KillmailID, kill.Hash)
	if err != nil {
		log.Printf("Error fetching ESI data for killmail %d: %v", kill.KillmailID, err)
		return
	}

	// Update kill with ESI data
	kill.KillTime = esiKill.KillTime
	kill.SolarSystemID = esiKill.SolarSystemID
	kill.Victim = esiKill.Victim
	kill.Attackers = esiKill.Attackers

	// Update kill in database
	err = db.UpdateKill(kill)
	if err != nil {
		log.Printf("Error updating kill %d with ESI data: %v", kill.KillmailID, err)
	}
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
			killTime := time.Time{}
			if zkill.KillmailTime != "" {
				parsedTime, err := time.Parse("2006-01-02T15:04:05Z", zkill.KillmailTime)
				if err != nil {
					log.Printf("Error parsing kill time for killmail %d: %v", zkill.KillmailID, err)
					continue
				}
				killTime = parsedTime
			} else {
				log.Printf("Warning: Empty killmail time for killmail %d", zkill.KillmailID)
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
