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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
	"github.com/tadeasf/eve-ran/src/utils"
)

const (
	baseURL = "https://esi.evetech.net/latest"
)

func FetchAndUpdateTypes() {
	utils.LogToConsole("Starting FetchAndUpdateTypes job")
	fetchAndUpdateRegions()
	fetchAndUpdateConstellations()
	fetchAndUpdateSystems()
	fetchAndUpdateItems()
	utils.LogToConsole("Finished FetchAndUpdateTypes job")
}

func fetchAndUpdateRegions() {
	log.Println("Fetching and updating regions")
	regions, err := services.FetchAllRegions(10)
	if err != nil {
		log.Printf("Error fetching regions: %v", err)
		return
	}

	for _, region := range regions {
		err := queries.UpsertRegion(region)
		if err != nil {
			log.Printf("Error upserting region %d: %v", region.RegionID, err)
		}
	}
	log.Println("Finished fetching and updating regions")
}

func fetchAndUpdateConstellations() {
	log.Println("Fetching and updating constellations")
	url := baseURL + "/universe/constellations/"
	ids := fetchIDs(url)

	existingConstellations, _ := queries.GetAllConstellations()
	existingMap := make(map[int]bool)
	for _, constellation := range existingConstellations {
		existingMap[constellation.ConstellationID] = true
	}

	constellationsChan := make(chan *models.Constellation, len(ids))
	var wg sync.WaitGroup
	rateLimiter := time.Tick(10 * time.Millisecond)

	for _, id := range ids {
		if !existingMap[id] {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				<-rateLimiter // Wait for rate limiter
				constellation := fetchConstellation(id)
				if constellation != nil {
					constellationsChan <- constellation
				}
			}(id)
		}
	}

	go func() {
		wg.Wait()
		close(constellationsChan)
	}()

	batchSize := 250
	var constellationsBatch []*models.Constellation

	for constellation := range constellationsChan {
		constellationsBatch = append(constellationsBatch, constellation)

		if len(constellationsBatch) >= batchSize {
			err := queries.BatchUpsertConstellations(constellationsBatch)
			if err != nil {
				log.Printf("Error batch upserting constellations: %v", err)
			}
			constellationsBatch = []*models.Constellation{}
		}
	}

	// Upsert any remaining constellations
	if len(constellationsBatch) > 0 {
		err := queries.BatchUpsertConstellations(constellationsBatch)
		if err != nil {
			log.Printf("Error batch upserting remaining constellations: %v", err)
		}
	}

	log.Println("Finished fetching and updating constellations")
}

func fetchConstellation(id int) *models.Constellation {
	url := baseURL + "/universe/constellations/" + strconv.Itoa(id) + "/"
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error fetching constellation %d: %v", id, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching constellation %d: Status code %d", id, resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body for constellation %d: %v", id, err)
		return nil
	}

	var constellation models.Constellation
	err = json.Unmarshal(body, &constellation)
	if err != nil {
		log.Printf("Error unmarshaling constellation %d: %v", id, err)
		return nil
	}

	return &constellation
}

func fetchAndUpdateSystems() {
	log.Println("Fetching and updating systems")
	url := baseURL + "/universe/systems/"
	ids := fetchIDs(url)

	existingSystems, _ := queries.GetAllSystems()
	existingMap := make(map[int]bool)
	for _, system := range existingSystems {
		existingMap[system.SystemID] = true
	}

	systemsChan := make(chan *models.System, len(ids))
	var wg sync.WaitGroup
	rateLimiter := time.Tick(5 * time.Millisecond)

	for _, id := range ids {
		if !existingMap[id] {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				<-rateLimiter // Wait for rate limiter
				system := fetchSystem(id)
				if system != nil {
					// Fetch constellation to get region_id
					constellation, err := services.FetchConstellation(system.ConstellationID)
					if err == nil {
						system.RegionID = constellation.RegionID
					}
					systemsChan <- system
				}
			}(id)
		}
	}

	go func() {
		wg.Wait()
		close(systemsChan)
	}()

	batchSize := 1000
	var systemsBatch []*models.System

	for system := range systemsChan {
		systemsBatch = append(systemsBatch, system)

		if len(systemsBatch) >= batchSize {
			err := queries.BatchUpsertSystems(systemsBatch)
			if err != nil {
				log.Printf("Error batch upserting systems: %v", err)
			}
			systemsBatch = []*models.System{}
		}
	}

	// Upsert any remaining systems
	if len(systemsBatch) > 0 {
		err := queries.BatchUpsertSystems(systemsBatch)
		if err != nil {
			log.Printf("Error batch upserting remaining systems: %v", err)
		}
	}

	log.Println("Finished fetching and updating systems")
}

func fetchSystem(id int) *models.System {
	url := baseURL + "/universe/systems/" + strconv.Itoa(id) + "/"
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error fetching system %d: %v", id, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching system %d: Status code %d", id, resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body for system %d: %v", id, err)
		return nil
	}

	var system models.System
	err = json.Unmarshal(body, &system)
	if err != nil {
		log.Printf("Error unmarshaling system %d: %v", id, err)
		return nil
	}

	// Fetch constellation to get region_id
	constellation, err := services.FetchConstellation(system.ConstellationID)
	if err == nil {
		system.RegionID = constellation.RegionID
	} else {
		log.Printf("Error fetching constellation for system %d: %v", id, err)
	}

	return &system
}

func fetchAndUpdateItems() {
	log.Println("Fetching and updating items")
	baseURL := baseURL + "/universe/types/"

	existingItems, _ := queries.GetAllESIItems()
	existingMap := make(map[int]bool)
	for _, item := range existingItems {
		existingMap[item.TypeID] = true
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)
	itemIDsChan := make(chan int, 100)

	// Start worker goroutines
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range itemIDsChan {
				semaphore <- struct{}{}
				fetchAndSaveItem(id)
				<-semaphore
			}
		}()
	}

	page := 1
	for {
		itemIDs, err := fetchItemIDsWithPagination(baseURL, page)
		if err != nil {
			if err.Error() == "requested page does not exist" {
				log.Println("Reached the end of item pages")
				break
			}
			log.Printf("Error fetching item IDs for page %d: %v", page, err)
			break
		}

		for _, id := range itemIDs {
			if !existingMap[id] {
				itemIDsChan <- id
			}
		}

		page++
		time.Sleep(100 * time.Millisecond) // Small delay to avoid hitting rate limits
	}

	close(itemIDsChan)
	wg.Wait()

	log.Println("Finished fetching and updating items")
}

func fetchItemIDsWithPagination(baseURL string, page int) ([]int, error) {
	url := fmt.Sprintf("%s?datasource=tranquility&page=%d", baseURL, page)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("requested page does not exist")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ids []int
	err = json.Unmarshal(body, &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func fetchAndSaveItem(id int) {
	if id == 0 {
		log.Printf("Skipping item with ID 0")
		return
	}
	url := fmt.Sprintf("%s/universe/types/%d/?datasource=tranquility&language=en", baseURL, id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for item %d: %v", id, err)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching item %d: %v", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var item models.ESIItem
	err = json.Unmarshal(body, &item)
	if err != nil {
		log.Printf("Error unmarshaling item %d: %v", id, err)
		return
	}

	err = queries.UpsertESIItem(&item)
	if err != nil {
		log.Printf("Error upserting item %d: %v", id, err)
	}
}

func fetchIDs(url string) []int {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching IDs from %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ids []int
	json.Unmarshal(body, &ids)

	return ids
}
