package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	esiManager = GetGlobalESIManager()
)

// makeESIRequest is a helper function that makes ESI requests with proper rate limiting and header handling
func makeESIRequest(url string, method string) (*http.Response, error) {
	return makeESIRequestWithRetry(url, method, 3)
}

// makeESIRequestWithRetry makes ESI requests with retry logic for 429 errors
func makeESIRequestWithRetry(url string, method string, maxRetries int) (*http.Response, error) {
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retry attempt %d/%d for %s", attempt, maxRetries, url)
		}

		// Check if we can make a request
		if !esiManager.CanMakeRequest() {
			log.Printf("ESI rate limit reached, waiting for reset...")
			esiManager.WaitForReset()
		}

		// Apply backoff if we're approaching limits
		if shouldBackoff, backoffDuration := esiManager.ShouldBackoff(); shouldBackoff {
			log.Printf("ESI rate limit approaching, backing off for %v", backoffDuration)
			time.Sleep(backoffDuration)
		}

		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		// Set required headers as per ESI guidelines
		req.Header.Set("User-Agent", "EVE Ran Application - GitHub: tadeasf/eve-ran - Contact: github.com/tadeasf")
		req.Header.Set("Accept", "application/json")
		// Support caching with If-Modified-Since
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := esiClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}

		// Always update rate limiting information from response headers
		esiManager.UpdateLimitsFromHeaders(resp.Header)

		// Handle 429 (rate limited)
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt < maxRetries {
				log.Printf("ESI returned 429 for %s, waiting before retry...", url)
				esiManager.WaitForReset()
				continue
			}
			return nil, fmt.Errorf("ESI returned 429 after %d retries", maxRetries)
		}

		// Handle 420 (old rate limit) and 520 (errors)
		if resp.StatusCode == 420 || resp.StatusCode == 520 {
			resp.Body.Close()
			if attempt < maxRetries {
				log.Printf("ESI returned %d, waiting before retry...", resp.StatusCode)
				esiManager.WaitForReset()
				continue
			}
			return nil, fmt.Errorf("ESI rate limit exceeded: %d after %d retries", resp.StatusCode, maxRetries)
		}

		// Log rate limit status if approaching limits
		remaining, resetTime := esiManager.GetStatus()
		if remaining < 20 {
			log.Printf("ESI Error Limit LOW - Remaining: %d, Reset: %s", remaining, resetTime.Format(time.RFC3339))
		}

		// Decrement error count on 4xx client errors (but not 429)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
			esiManager.DecrementErrorCount()
		}

		// 5xx errors don't cost tokens, but we might want to retry
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			resp.Body.Close()
			if attempt < maxRetries {
				log.Printf("ESI returned %d (server error), retrying in 2s...", resp.StatusCode)
				time.Sleep(2 * time.Second)
				continue
			}
		}

		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded")
}

func FetchRegionIDs() ([]int, error) {
	url := fmt.Sprintf("%s/universe/regions/?datasource=tranquility", esiBaseURL)
	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
	}

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
	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
	}

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
	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
	}

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
	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
	}

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
	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
	}

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
	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
	}

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
		resp, err := makeESIRequest(url, "GET")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
		}

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
	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
	}

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

	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d for killmail %d", resp.StatusCode, killmailID)
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
	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, fmt.Errorf("error fetching constellation: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
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

// SearchCharactersByName searches for characters by name using EVE ESI API
func SearchCharactersByName(searchTerm string) ([]int64, error) {
	// Use the /universe/ids/ endpoint which resolves names to IDs
	url := fmt.Sprintf("%s/universe/ids/?datasource=tranquility", esiBaseURL)

	// Create the request body with the search term
	requestBody := []string{searchTerm}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "EVE Ran Application - GitHub: tadeasf/eve-ran - Contact: github.com/tadeasf")
	req.Header.Set("Accept", "application/json")

	resp, err := esiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error searching characters: %v", err)
	}
	defer resp.Body.Close()

	// Update rate limiting information from response headers
	esiManager.UpdateLimitsFromHeaders(resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d, body: %s", resp.StatusCode, string(body))
	}

	var searchResult struct {
		Characters []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"characters"`
	}

	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling search result: %v", err)
	}

	// Extract character IDs
	var characterIDs []int64
	for _, character := range searchResult.Characters {
		characterIDs = append(characterIDs, character.ID)
	}

	return characterIDs, nil
}

// FetchCharacterInfo fetches character information by ID
func FetchCharacterInfo(characterID int64) (*models.Character, error) {
	url := fmt.Sprintf("%s/characters/%d/?datasource=tranquility", esiBaseURL, characterID)

	resp, err := makeESIRequest(url, "GET")
	if err != nil {
		return nil, fmt.Errorf("error fetching character info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var esiCharacter struct {
		CharacterID    int64   `json:"character_id"`
		Name           string  `json:"name"`
		SecurityStatus float64 `json:"security_status"`
		Title          string  `json:"title"`
		RaceID         int     `json:"race_id"`
	}

	err = json.Unmarshal(body, &esiCharacter)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling character data: %v", err)
	}

	return &models.Character{
		ID:             characterID,
		Name:           esiCharacter.Name,
		SecurityStatus: esiCharacter.SecurityStatus,
		Title:          esiCharacter.Title,
		RaceID:         esiCharacter.RaceID,
	}, nil
}
