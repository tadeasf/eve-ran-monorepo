package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
)

func InitializeCharacterKills(characterID int64) error {
	page := 1
	for {
		kills, err := FetchKillsFromZKillboard(characterID, page)
		if err != nil {
			return err
		}

		if len(kills) == 0 {
			break
		}

		err = StoreKills(characterID, kills)
		if err != nil {
			return err
		}

		page++
	}

	return nil
}

func FetchKillsFromZKillboard(characterID int64, page int) ([]models.Zkill, error) {
	url := fmt.Sprintf("https://zkillboard.com/api/kills/characterID/%d/page/%d/", characterID, page)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawKills []struct {
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

	err = json.NewDecoder(resp.Body).Decode(&rawKills)
	if err != nil {
		return nil, err
	}

	var kills []models.Zkill
	for _, rawKill := range rawKills {
		kill := models.Zkill{
			KillmailID:     rawKill.KillmailID,
			CharacterID:    characterID,
			LocationID:     rawKill.ZKB.LocationID,
			Hash:           rawKill.ZKB.Hash,
			FittedValue:    rawKill.ZKB.FittedValue,
			DroppedValue:   rawKill.ZKB.DroppedValue,
			DestroyedValue: rawKill.ZKB.DestroyedValue,
			TotalValue:     rawKill.ZKB.TotalValue,
			Points:         rawKill.ZKB.Points,
			NPC:            rawKill.ZKB.NPC,
			Solo:           rawKill.ZKB.Solo,
			Awox:           rawKill.ZKB.Awox,
		}
		kills = append(kills, kill)
	}

	return kills, nil
}

func StoreKills(characterID int64, kills []models.Zkill) error {
	for _, kill := range kills {
		if err := db.DB.Create(&kill).Error; err != nil {
			return err
		}
	}
	return nil
}
