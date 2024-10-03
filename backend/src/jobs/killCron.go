package jobs

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tadeasf/eve-ran/src/db/queries"
)

var killCron *cron.Cron

func StartKillCron() {
	killCron = cron.New()

	// Run the job every hour
	_, err := killCron.AddFunc("@every 30s", queueNewKillFetches)
	if err != nil {
		log.Printf("Error adding kill fetch cron job: %v", err)
		return
	}

	killCron.Start()
	log.Println("Kill fetch cron job started")
}

func StopKillCron() {
	if killCron != nil {
		killCron.Stop()
		log.Println("Kill fetch cron job stopped")
	}
}

func queueNewKillFetches() {
	log.Println("Queueing periodic kill fetches")

	esiErrorLimitMutex.Lock()
	if time.Now().Before(esiErrorLimitBackoff) {
		sleepTime := time.Until(esiErrorLimitBackoff)
		esiErrorLimitMutex.Unlock()
		log.Printf("Waiting for ESI error limit backoff: %v", sleepTime)
		time.Sleep(sleepTime)
	} else {
		esiErrorLimitMutex.Unlock()
	}

	characters, err := queries.GetAllCharacters()
	if err != nil {
		log.Printf("Error fetching characters: %v", err)
		return
	}

	for _, character := range characters {
		lastKillTime, err := queries.GetLastKillTimeForCharacter(character.ID)
		if err != nil {
			log.Printf("Error getting last kill time for character %d: %v", character.ID, err)
			continue
		}

		// If no kills found, fetch kills from the last 24 hours
		if lastKillTime.IsZero() {
			lastKillTime = time.Now().Add(-24 * time.Hour)
		}

		QueueKillFetch(character.ID, lastKillTime, false)
	}

	log.Println("Periodic kill fetch queuing completed")
}
