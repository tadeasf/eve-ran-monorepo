package jobs

import (
	"fmt"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
	"github.com/tadeasf/eve-ran/src/utils"
)

const esiBaseURL = "https://esi.evetech.net/latest"

func EnhanceKills() {
	// Get all killmail IDs from the kills table
	var existingKillmailIDs []int64
	if err := db.DB.Model(&models.Kill{}).Pluck("killmail_id", &existingKillmailIDs).Error; err != nil {
		utils.LogError(fmt.Sprintf("Error fetching existing killmail IDs: %v", err))
		return
	}

	utils.LogToConsole(fmt.Sprintf("Number of existing killmail IDs: %d", len(existingKillmailIDs)))

	// Create a map for faster lookup
	existingKillmailIDMap := make(map[int64]bool)
	for _, id := range existingKillmailIDs {
		existingKillmailIDMap[id] = true
	}

	// Fetch zkills that are not in the kills table
	var zkillsToEnhance []models.Zkill
	if err := db.DB.Where("killmail_id NOT IN (?)", existingKillmailIDs).Find(&zkillsToEnhance).Error; err != nil {
		utils.LogError(fmt.Sprintf("Error fetching Zkills to enhance: %v", err))
		return
	}

	utils.LogToConsole(fmt.Sprintf("Enhancing %d new kills", len(zkillsToEnhance)))

	for _, zkill := range zkillsToEnhance {
		enhancedKill, err := fetchEnhancedKillData(zkill)
		if err != nil {
			utils.LogError(fmt.Sprintf("Error enhancing kill %d: %v", zkill.KillmailID, err))
			continue
		}

		// Validate that we got real data before storing
		if enhancedKill.KillmailID == 0 {
			utils.LogError(fmt.Sprintf("Skipping kill %d: received empty data from ESI", zkill.KillmailID))
			continue
		}

		// Use UpsertKill instead of Create to handle duplicates
		if err := queries.UpsertKill(enhancedKill); err != nil {
			utils.LogError(fmt.Sprintf("Error upserting kill %d: %v", zkill.KillmailID, err))
		} else {
			utils.LogToConsole(fmt.Sprintf("Upserted kill: %d", zkill.KillmailID))
		}
	}
}

func fetchEnhancedKillData(zkill models.Zkill) (*models.Kill, error) {
	// Use the centralized ESI service which handles rate limiting
	kill, err := services.FetchKillmailFromESI(zkill.KillmailID, zkill.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance kill %d: %v", zkill.KillmailID, err)
	}
	return kill, nil
}

func EnhanceKill(killmailID int64) (*models.Kill, error) {
	// First, get the zKill data
	zkill, err := queries.GetZKillByID(killmailID)
	if err != nil {
		return nil, fmt.Errorf("failed to get zkill data: %v", err)
	}

	// Use the centralized ESI service which handles rate limiting
	kill, err := services.FetchKillmailFromESI(zkill.KillmailID, zkill.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance kill %d: %v", zkill.KillmailID, err)
	}

	return kill, nil
}
