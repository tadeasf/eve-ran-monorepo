package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
)

const esiBaseURL = "https://esi.evetech.net/latest"

var (
	esiClient = &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
)

func FetchRegionIDs() ([]int, error) {
	url := fmt.Sprintf("%s/universe/regions/?datasource=tranquility", esiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var regionIDs []int
	err = json.Unmarshal(body, &regionIDs)
	return regionIDs, err
}

func FetchRegionInfo(regionID int) (*models.Region, error) {
	url := fmt.Sprintf("%s/universe/regions/%d/?datasource=tranquility&language=en", esiBaseURL, regionID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var region models.Region
	err = json.Unmarshal(body, &region)
	if err != nil {
		return nil, err
	}

	// Ensure Constellations is initialized as an empty slice if it's null
	if region.Constellations == nil {
		region.Constellations = json.RawMessage("[]")
	}

	return &region, nil
}

func FetchSystemIDs() ([]int, error) {
	url := fmt.Sprintf("%s/universe/systems/?datasource=tranquility", esiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var systemIDs []int
	err = json.Unmarshal(body, &systemIDs)
	return systemIDs, err
}

func FetchSystemInfo(systemID int) (*models.System, error) {
	url := fmt.Sprintf("%s/universe/systems/%d/?datasource=tranquility&language=en", esiBaseURL, systemID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var system models.System
	err = json.Unmarshal(body, &system)
	return &system, err
}

func FetchConstellationIDs() ([]int, error) {
	url := fmt.Sprintf("%s/universe/constellations/?datasource=tranquility", esiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var constellationIDs []int
	err = json.Unmarshal(body, &constellationIDs)
	return constellationIDs, err
}

func FetchConstellationInfo(constellationID int) (*models.Constellation, error) {
	url := fmt.Sprintf("%s/universe/constellations/%d/?datasource=tranquility&language=en", esiBaseURL, constellationID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var constellation models.Constellation
	err = json.Unmarshal(body, &constellation)
	return &constellation, err
}

func FetchItemIDs() ([]int, error) {
	var allItemIDs []int
	page := 1
	for {
		url := fmt.Sprintf("%s/universe/types/?datasource=tranquility&page=%d", esiBaseURL, page)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var itemIDs []int
		err = json.Unmarshal(body, &itemIDs)
		if err != nil {
			return nil, err
		}

		if len(itemIDs) == 0 {
			break
		}

		allItemIDs = append(allItemIDs, itemIDs...)
		page++
	}

	return allItemIDs, nil
}

func FetchItemInfo(itemID int) (*models.ESIItem, error) {
	url := fmt.Sprintf("%s/universe/types/%d/?datasource=tranquility&language=en", esiBaseURL, itemID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item models.ESIItem
	err = json.Unmarshal(body, &item)
	return &item, err
}

func FetchAllItems(concurrency int) ([]*models.ESIItem, error) {
	itemIDs, err := FetchItemIDs()
	if err != nil {
		return nil, err
	}

	items := make([]*models.ESIItem, 0, len(itemIDs))
	itemChan := make(chan *models.ESIItem, len(itemIDs))
	errChan := make(chan error, len(itemIDs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, itemID := range itemIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			item, err := FetchItemInfo(id)
			if err != nil {
				errChan <- err
				return
			}
			itemChan <- item
		}(itemID)
	}

	go func() {
		wg.Wait()
		close(itemChan)
		close(errChan)
	}()

	for item := range itemChan {
		items = append(items, item)
	}

	if len(errChan) > 0 {
		return items, <-errChan
	}

	return items, nil
}

func FetchAllRegions(concurrency int) ([]*models.Region, error) {
	regionIDs, err := FetchRegionIDs()
	if err != nil {
		return nil, err
	}

	regions := make([]*models.Region, 0, len(regionIDs))
	regionChan := make(chan *models.Region, len(regionIDs))
	errChan := make(chan error, len(regionIDs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, regionID := range regionIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			region, err := FetchRegionInfo(id)
			if err != nil {
				errChan <- err
				return
			}
			regionChan <- region
		}(regionID)
	}

	go func() {
		wg.Wait()
		close(regionChan)
		close(errChan)
	}()

	for region := range regionChan {
		regions = append(regions, region)
	}

	if len(errChan) > 0 {
		return regions, <-errChan
	}

	return regions, nil
}

func FetchAllConstellations(concurrency int) ([]*models.Constellation, error) {
	constellationIDs, err := FetchConstellationIDs()
	if err != nil {
		return nil, err
	}

	constellations := make([]*models.Constellation, 0, len(constellationIDs))
	constellationChan := make(chan *models.Constellation, len(constellationIDs))
	errChan := make(chan error, len(constellationIDs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, constellationID := range constellationIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			constellation, err := FetchConstellationInfo(id)
			if err != nil {
				errChan <- err
				return
			}
			constellationChan <- constellation
		}(constellationID)
	}

	go func() {
		wg.Wait()
		close(constellationChan)
		close(errChan)
	}()

	for constellation := range constellationChan {
		constellations = append(constellations, constellation)
	}

	if len(errChan) > 0 {
		return constellations, <-errChan
	}

	return constellations, nil
}

func FetchAllSystems(concurrency int) ([]*models.System, error) {
	systemIDs, err := FetchSystemIDs()
	if err != nil {
		return nil, err
	}

	systems := make([]*models.System, 0, len(systemIDs))
	systemChan := make(chan *models.System, len(systemIDs))
	errChan := make(chan error, len(systemIDs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for _, systemID := range systemIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			system, err := FetchSystemInfo(id)
			if err != nil {
				errChan <- err
				return
			}
			systemChan <- system
		}(systemID)
	}

	go func() {
		wg.Wait()
		close(systemChan)
		close(errChan)
	}()

	for system := range systemChan {
		systems = append(systems, system)
	}

	if len(errChan) > 0 {
		return systems, <-errChan
	}

	return systems, nil
}

func FetchKillmailFromESI(killmailID int64, hash string) (*models.Kill, error) {
	url := fmt.Sprintf("%s/killmails/%d/%s/?datasource=tranquility", esiBaseURL, killmailID, hash)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("User-Agent", "EVE Ran Application - GitHub: tadeasf/eve-ran")

	resp, err := esiClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return nil, fmt.Errorf("ESI timeout: %v", err)
		}
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var esiKill struct {
		KillmailID    int64             `json:"killmail_id"`
		KillmailTime  time.Time         `json:"killmail_time"`
		SolarSystemID int               `json:"solar_system_id"`
		Victim        models.Victim     `json:"victim"`
		Attackers     []models.Attacker `json:"attackers"`
	}
	err = json.Unmarshal(body, &esiKill)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling killmail: %v", err)
	}

	// Marshal the Attackers slice into JSON
	attackersJSON, err := json.Marshal(esiKill.Attackers)
	if err != nil {
		return nil, fmt.Errorf("error marshaling attackers: %v", err)
	}

	return &models.Kill{
		KillmailID:    esiKill.KillmailID,
		KillmailTime:  esiKill.KillmailTime,
		SolarSystemID: esiKill.SolarSystemID,
		Victim:        esiKill.Victim,
		Attackers:     attackersJSON,
	}, nil
}

func IsESITimeout(err error) bool {
	return strings.Contains(err.Error(), "Timeout contacting tranquility")
}

func IsESIErrorLimit(err error) bool {
	return strings.Contains(err.Error(), "ESI error limit reached") ||
		strings.Contains(err.Error(), "This software has exceeded the error limit for ESI")
}

// Add this function to the existing esi.go file

func FetchConstellation(constellationID int) (*models.Constellation, error) {
	url := fmt.Sprintf("%s/universe/constellations/%d/?datasource=tranquility&language=en", esiBaseURL, constellationID)
	resp, err := esiClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching constellation: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned non-OK status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var constellation models.Constellation
	err = json.Unmarshal(body, &constellation)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling constellation data: %v", err)
	}

	return &constellation, nil
}
