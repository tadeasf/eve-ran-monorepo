package jobs

import (
	"fmt"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
)

func StartKillCron() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		checkNewKills()
	}
}

func checkNewKills() {
	characters, err := queries.GetAllCharacters()
	if err != nil {
		fmt.Printf("Error fetching characters: %v\n", err)
		return
	}

	for _, character := range characters {
		page := 1
		for {
			kills, err := FetchKillsFromZKillboard(character.ID, page)
			if err != nil {
				fmt.Printf("Error fetching kills for character %d: %v\n", character.ID, err)
				break
			}

			if len(kills) == 0 {
				break
			}

			newKills := filterNewKills(kills)
			if len(newKills) == 0 {
				break
			}

			err = StoreKills(character.ID, newKills)
			if err != nil {
				fmt.Printf("Error storing new kills for character %d: %v\n", character.ID, err)
				break
			}

			page++
		}
	}
}

func filterNewKills(kills []models.Zkill) []models.Zkill {
	var newKills []models.Zkill
	for _, kill := range kills {
		exists, err := queries.KillExists(kill.KillmailID)
		if err != nil {
			fmt.Printf("Error checking if kill %d exists: %v\n", kill.KillmailID, err)
			continue
		}
		if !exists {
			newKills = append(newKills, kill)
		}
	}
	return newKills
}
