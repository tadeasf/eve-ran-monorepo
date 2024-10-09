// Copyright (C) 2024 Tadeáš Fořt
// 
// This file is part of EVE Ran Services.
// 
// EVE Ran Services is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// EVE Ran Services is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with EVE Ran Services.  If not, see <https://www.gnu.org/licenses/>.

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
			zkills, err := FetchKillsFromZKillboard(character.ID, page)
			if err != nil {
				fmt.Printf("Error fetching kills for character %d: %v\n", character.ID, err)
				break
			}

			if len(zkills) == 0 {
				break
			}

			newZKills := filterNewZKills(zkills)
			if len(newZKills) == 0 {
				break
			}

			err = StoreZKills(newZKills)
			if err != nil {
				fmt.Printf("Error storing new zkills for character %d: %v\n", character.ID, err)
				break
			}

			for _, zkill := range newZKills {
				err = EnhanceAndStoreKill(zkill)
				if err != nil {
					fmt.Printf("Error enhancing and storing kill %d: %v\n", zkill.KillmailID, err)
				}
			}

			page++
		}
	}
}

func filterNewZKills(zkills []models.Zkill) []models.Zkill {
	var newZKills []models.Zkill
	for _, zkill := range zkills {
		exists, err := queries.ZKillExists(zkill.KillmailID)
		if err != nil {
			fmt.Printf("Error checking if zkill %d exists: %v\n", zkill.KillmailID, err)
			continue
		}
		if !exists {
			newZKills = append(newZKills, zkill)
		}
	}
	return newZKills
}
