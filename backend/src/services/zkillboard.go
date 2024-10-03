package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tadeasf/eve-ran/src/db/models"
)

func FetchKillsFromZKillboard(characterID int64, page int) ([]models.ZKillKill, error) {
	url := fmt.Sprintf("https://zkillboard.com/api/kills/characterID/%d/page/%d/", characterID, page)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "EVE Ran Application - GitHub: tadeasf/eve-ran")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var zkillKills []models.ZKillKill
	err = json.NewDecoder(resp.Body).Decode(&zkillKills)
	if err != nil {
		return nil, err
	}

	return zkillKills, nil
}
