package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"

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

		utils.LogToFile(fmt.Sprintf("Enhanced kill data: %+v", enhancedKill))

		// Create new Kill entry
		if err := db.DB.Create(enhancedKill).Error; err != nil {
			utils.LogError(fmt.Sprintf("Error storing enhanced kill %d: %v", zkill.KillmailID, err))
		} else {
			utils.LogToConsole(fmt.Sprintf("Added new kill: %d", zkill.KillmailID))
		}
	}
}

func fetchEnhancedKillData(zkill models.Zkill) (*models.Kill, error) {
	url := fmt.Sprintf("%s/killmails/%d/%s/?datasource=tranquility", esiBaseURL, zkill.KillmailID, zkill.Hash)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var esiKill struct {
		KillmailID    int64     `json:"killmail_id"`
		KillmailTime  time.Time `json:"killmail_time"`
		SolarSystemID int       `json:"solar_system_id"`
		Victim        struct {
			AllianceID    int64 `json:"alliance_id"`
			CharacterID   int64 `json:"character_id"`
			CorporationID int64 `json:"corporation_id"`
			DamageTaken   int   `json:"damage_taken"`
			ShipTypeID    int   `json:"ship_type_id"`
			Position      struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"position"`
		} `json:"victim"`
		Attackers []json.RawMessage `json:"attackers"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&esiKill); err != nil {
		return nil, err
	}

	// Convert Attackers to JSON byte array
	attackersJSON, err := json.Marshal(esiKill.Attackers)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal attackers: %v", err)
	}

	enhancedKill := &models.Kill{
		KillmailID:    esiKill.KillmailID,
		KillmailTime:  esiKill.KillmailTime,
		SolarSystemID: esiKill.SolarSystemID,
		Victim: models.Victim{
			AllianceID:    esiKill.Victim.AllianceID,
			CharacterID:   esiKill.Victim.CharacterID,
			CorporationID: esiKill.Victim.CorporationID,
			DamageTaken:   esiKill.Victim.DamageTaken,
			ShipTypeID:    esiKill.Victim.ShipTypeID,
			Position: models.Position{
				X: esiKill.Victim.Position.X,
				Y: esiKill.Victim.Position.Y,
				Z: esiKill.Victim.Position.Z,
			},
		},
		Attackers: attackersJSON,
		ZkillData: zkill,
	}

	return enhancedKill, nil
}

func EnhanceKill(killmailID int64) (*models.Kill, error) {
	// First, get the zKill data
	zkill, err := queries.GetZKillByID(killmailID)
	if err != nil {
		return nil, fmt.Errorf("failed to get zkill data: %v", err)
	}

	// Then fetch the killmail data from ESI
	url := fmt.Sprintf("%s/killmails/%d/%s/?datasource=tranquility", esiBaseURL, killmailID, zkill.Hash)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch killmail from ESI: %v", err)
	}
	defer resp.Body.Close()

	var esiKill struct {
		KillmailID    int64     `json:"killmail_id"`
		KillmailTime  time.Time `json:"killmail_time"`
		SolarSystemID int       `json:"solar_system_id"`
		Victim        struct {
			AllianceID    int64 `json:"alliance_id"`
			CharacterID   int64 `json:"character_id"`
			CorporationID int64 `json:"corporation_id"`
			DamageTaken   int   `json:"damage_taken"`
			ShipTypeID    int   `json:"ship_type_id"`
			Position      struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"position"`
		} `json:"victim"`
		Attackers []json.RawMessage `json:"attackers"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&esiKill); err != nil {
		return nil, fmt.Errorf("failed to decode ESI response: %v", err)
	}

	// Marshal the entire Attackers slice into a single JSON byte array
	attackersJSON, err := json.Marshal(esiKill.Attackers)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal attackers: %v", err)
	}

	enhancedKill := &models.Kill{
		KillmailID:    killmailID,
		KillmailTime:  esiKill.KillmailTime,
		SolarSystemID: esiKill.SolarSystemID,
		CharacterID:   zkill.CharacterID, // Use CharacterID from zKill data
		Victim: models.Victim{
			AllianceID:    esiKill.Victim.AllianceID,
			CharacterID:   esiKill.Victim.CharacterID,
			CorporationID: esiKill.Victim.CorporationID,
			DamageTaken:   esiKill.Victim.DamageTaken,
			ShipTypeID:    esiKill.Victim.ShipTypeID,
			Position: models.Position{
				X: esiKill.Victim.Position.X,
				Y: esiKill.Victim.Position.Y,
				Z: esiKill.Victim.Position.Z,
			},
		},
		Attackers: attackersJSON, // Now this is a []byte
		ZkillData: *zkill,
	}

	return enhancedKill, nil
}
