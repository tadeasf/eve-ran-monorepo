package jobs

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchKillsForCharacter(characterID int64) (int, error) {
	totalNewKills := int32(0)
	page := 1
	batchSize := 1
	maxConcurrent := 1000

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

				processBatch(batch, characterID, &totalNewKills)
			}(zkillKills[i:end])
		}

		wg.Wait()
		log.Printf("Processed %d new kills for character %d on page %d", atomic.LoadInt32(&totalNewKills), characterID, page)

		if len(zkillKills) < 200 {
			log.Printf("Less than 200 new kills on page %d for character %d, stopping", page, characterID)
			break
		}

		page++
	}

	return int(totalNewKills), nil
}

func processBatch(batch []models.ZKillKill, characterID int64, newKills *int32) {
	var kills []models.Kill
	for _, zkillKill := range batch {
		existingKill, err := queries.GetKillByKillmailID(zkillKill.KillmailID)
		if err != nil {
			log.Printf("Error checking existing kill %d: %v", zkillKill.KillmailID, err)
			continue
		}

		if existingKill != nil && !existingKill.KillTime.IsZero() {
			// Kill already exists and has been enriched, skip it
			continue
		}

		esiKill, err := services.FetchKillmailFromESI(zkillKill.KillmailID, zkillKill.ZKB.Hash)
		if err != nil {
			log.Printf("Error fetching killmail %d from ESI: %v", zkillKill.KillmailID, err)
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
