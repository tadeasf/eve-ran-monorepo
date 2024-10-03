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
	_, err := killCron.AddFunc("@hourly", fetchNewKills)
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

func fetchNewKills() {
	log.Println("Starting periodic kill fetch")

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

		newKills, err := FetchNewKillsForCharacter(character.ID, lastKillTime)
		if err != nil {
			log.Printf("Error fetching new kills for character %d: %v", character.ID, err)
			continue
		}

		log.Printf("Fetched %d new kills for character %d", newKills, character.ID)
	}

	log.Println("Periodic kill fetch completed")
}
