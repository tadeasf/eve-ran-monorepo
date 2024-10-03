package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
)

const esiBaseURL = "https://esi.evetech.net/latest"

func EnhanceKills() {
	var zkills []models.Zkill
	result := db.DB.Find(&zkills)
	if result.Error != nil {
		fmt.Printf("Error fetching Zkills: %v\n", result.Error)
		return
	}

	for _, zkill := range zkills {
		enhancedKill, err := fetchEnhancedKillData(zkill)
		if err != nil {
			fmt.Printf("Error enhancing kill %d: %v\n", zkill.KillmailID, err)
			continue
		}

		// Log the enhancedKill data for debugging
		fmt.Printf("Enhanced kill data: %+v\n", enhancedKill)

		// Skip invalid killmail IDs
		if enhancedKill.KillmailID == 0 {
			fmt.Printf("Skipping invalid killmail ID: %d\n", zkill.KillmailID)
			continue
		}

		// Check if a Kill entry already exists
		var existingKill models.Kill
		if err := db.DB.Where("killmail_id = ?", zkill.KillmailID).First(&existingKill).Error; err == nil {
			// Update existing Kill entry
			existingKill.KillmailTime = enhancedKill.KillmailTime
			existingKill.SolarSystemID = enhancedKill.SolarSystemID
			existingKill.Victim = enhancedKill.Victim
			existingKill.Attackers = enhancedKill.Attackers
			existingKill.ZkillData = zkill

			if err := db.DB.Save(&existingKill).Error; err != nil {
				fmt.Printf("Error updating enhanced kill %d: %v\n", zkill.KillmailID, err)
			}
		} else {
			// Create new Kill entry
			if err := db.DB.Create(enhancedKill).Error; err != nil {
				fmt.Printf("Error storing enhanced kill %d: %v\n", zkill.KillmailID, err)
			}
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

	// Convert Attackers to JSONB
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
