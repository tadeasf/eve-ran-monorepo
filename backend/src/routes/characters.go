package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
)

// GetAllCharacters retrieves all characters from the database
// @Summary Get all characters
// @Description Fetch all characters from the database
// @Tags characters
// @Accept json
// @Produce json
// @Success 200 {array} models.Character
// @Failure 500 {object} models.ErrorResponse
// @Router /characters [get]
func GetAllCharacters(c *gin.Context) {
	characters, err := queries.GetAllCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, characters)
}

// GetAllKills retrieves kills from the database with optional filters and pagination
// @Summary Get kills with filters
// @Description Fetch kills from the database, optionally filtered by character_id, region, and date range with pagination
// @Tags kills
// @Accept json
// @Produce json
// @Param character_id query int false "Filter by character ID"
// @Param region_id query int false "Filter by region ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 100, max: 1000)"
// @Success 200 {object} map[string]interface{} "Returns kills array, total count, page, and page_size"
// @Failure 500 {object} models.ErrorResponse
// @Router /kills [get]
func GetAllKills(c *gin.Context) {
	// Get optional query parameters
	characterIDStr := c.Query("character_id")
	regionIDStr := c.Query("region_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Parse pagination parameters
	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 100
	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 {
			pageSize = ps
			// Cap at 1000 to prevent excessive memory usage
			if pageSize > 1000 {
				pageSize = 1000
			}
		}
	}

	// Build base query
	baseQuery := db.DB.Model(&models.Kill{})

	if characterIDStr != "" {
		characterID, err := strconv.ParseInt(characterIDStr, 10, 64)
		if err == nil {
			baseQuery = baseQuery.Where("character_id = ?", characterID)
		}
	}

	if regionIDStr != "" {
		regionID, err := strconv.Atoi(regionIDStr)
		if err == nil {
			// Join with systems table to filter by region
			baseQuery = baseQuery.Joins("JOIN systems ON kills.solar_system_id = systems.system_id").
				Where("systems.region_id = ?", regionID)
		}
	}

	if startDate != "" && endDate != "" {
		baseQuery = baseQuery.Where("killmail_time BETWEEN ? AND ?", startDate, endDate)
	}

	// Get total count
	var totalCount int64
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch paginated results
	var kills []models.Kill
	offset := (page - 1) * pageSize
	query := baseQuery.Preload("ZkillData").
		Order("killmail_time DESC").
		Limit(pageSize).
		Offset(offset)

	if err := query.Find(&kills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return paginated response with metadata
	response := gin.H{
		"kills":       kills,
		"total_count": totalCount,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (totalCount + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, response)
}

// GetAllCharacterStats retrieves stats for all characters with filters
// @Summary Get all character stats
// @Description Fetch stats for all characters from the database with optional filters. Returns aggregated kill counts and ISK destroyed per character.
// @Tags characters
// @Accept json
// @Produce json
// @Param regionID query []int false "Region IDs (can specify multiple: ?regionID=10000001&regionID=10000002)"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} models.CharacterStats
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/stats [get]
func GetAllCharacterStats(c *gin.Context) {
	regionIDs := c.QueryArray("regionID")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	var regionIDInts []int64
	for _, id := range regionIDs {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
			return
		}
		regionIDInts = append(regionIDInts, intID)
	}

	var startTime, endTime time.Time
	var err error
	if startDate != "" {
		startTime, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
	} else {
		// Default to 30 days ago if not specified
		startTime = time.Now().AddDate(0, 0, -30)
	}

	if endDate != "" {
		endTime, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	} else {
		// Default to today if not specified
		endTime = time.Now()
	}

	// Get stats from optimized query
	stats, err := queries.GetCharacterStats(startTime, endTime, 0, regionIDInts...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enrich with character names
	characters, err := queries.GetAllCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a map for quick character name lookup
	charMap := make(map[int64]string)
	for _, char := range characters {
		charMap[char.ID] = char.Name
	}

	// Add names to stats
	type EnrichedCharacterStats struct {
		CharacterID int64   `json:"character_id"`
		Name        string  `json:"name"`
		KillCount   int     `json:"kill_count"`
		TotalISK    float64 `json:"total_isk"`
	}

	enrichedStats := make([]EnrichedCharacterStats, 0, len(stats))
	for _, stat := range stats {
		enrichedStats = append(enrichedStats, EnrichedCharacterStats{
			CharacterID: stat.CharacterID,
			Name:        charMap[stat.CharacterID],
			KillCount:   stat.KillCount,
			TotalISK:    stat.TotalISK,
		})
	}

	c.JSON(http.StatusOK, enrichedStats)
}

// GetCharacterAnalytics retrieves aggregated analytics for a specific character
// @Summary Get character analytics
// @Description Fetch aggregated analytics data for a character (hourly activity, daily stats, top systems)
// @Tags characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id}/analytics [get]
func GetCharacterAnalytics(c *gin.Context) {
	characterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Build query for character's kills
	query := db.DB.Model(&models.Kill{}).
		Preload("ZkillData").
		Where("character_id = ?", characterID)

	if startDate != "" && endDate != "" {
		query = query.Where("killmail_time BETWEEN ? AND ?", startDate, endDate)
	}

	var kills []models.Kill
	if err := query.Find(&kills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Aggregate analytics data
	type HourlyActivity struct {
		Hour  int `json:"hour"`
		Kills int `json:"kills"`
	}

	type DailyActivity struct {
		Date  string  `json:"date"`
		Kills int     `json:"kills"`
		ISK   float64 `json:"isk"`
	}

	type SystemActivity struct {
		SystemID   int     `json:"system_id"`
		SystemName string  `json:"system_name"`
		Kills      int     `json:"kills"`
		ISK        float64 `json:"isk"`
	}

	hourlyMap := make(map[int]int)
	dailyMap := make(map[string]*DailyActivity)
	systemMap := make(map[int]*SystemActivity)

	var totalKills int
	var totalISK float64
	var totalPoints int

	for _, kill := range kills {
		totalKills++
		totalISK += kill.ZkillData.TotalValue
		totalPoints += kill.ZkillData.Points

		// Hourly activity
		hour := kill.KillmailTime.Hour()
		hourlyMap[hour]++

		// Daily activity
		date := kill.KillmailTime.Format("2006-01-02")
		if _, exists := dailyMap[date]; !exists {
			dailyMap[date] = &DailyActivity{Date: date, Kills: 0, ISK: 0}
		}
		dailyMap[date].Kills++
		dailyMap[date].ISK += kill.ZkillData.TotalValue

		// System activity
		if _, exists := systemMap[kill.SolarSystemID]; !exists {
			systemMap[kill.SolarSystemID] = &SystemActivity{
				SystemID: kill.SolarSystemID,
				Kills:    0,
				ISK:      0,
			}
		}
		systemMap[kill.SolarSystemID].Kills++
		systemMap[kill.SolarSystemID].ISK += kill.ZkillData.TotalValue
	}

	// Convert maps to arrays
	hourlyActivity := make([]HourlyActivity, 0, len(hourlyMap))
	for hour, kills := range hourlyMap {
		hourlyActivity = append(hourlyActivity, HourlyActivity{Hour: hour, Kills: kills})
	}

	dailyActivity := make([]DailyActivity, 0, len(dailyMap))
	for _, activity := range dailyMap {
		dailyActivity = append(dailyActivity, *activity)
	}

	systemActivity := make([]SystemActivity, 0, len(systemMap))
	for _, activity := range systemMap {
		systemActivity = append(systemActivity, *activity)
	}

	// Get system names
	if len(systemActivity) > 0 {
		systemIDs := make([]int, 0, len(systemActivity))
		for _, sys := range systemActivity {
			systemIDs = append(systemIDs, sys.SystemID)
		}

		var systems []models.System
		db.DB.Where("system_id IN ?", systemIDs).Find(&systems)

		systemNameMap := make(map[int]string)
		for _, sys := range systems {
			systemNameMap[sys.SystemID] = sys.Name
		}

		for i := range systemActivity {
			systemActivity[i].SystemName = systemNameMap[systemActivity[i].SystemID]
		}
	}

	response := gin.H{
		"total_kills":     totalKills,
		"total_isk":       totalISK,
		"total_points":    totalPoints,
		"hourly_activity": hourlyActivity,
		"daily_activity":  dailyActivity,
		"system_activity": systemActivity,
	}

	c.JSON(http.StatusOK, response)
}

// SearchCharacters searches for characters by name using EVE ESI API
// @Summary Search characters by name
// @Description Search for EVE Online characters by name using ESI API
// @Tags characters
// @Accept json
// @Produce json
// @Param search query string true "Character name to search"
// @Success 200 {array} models.Character
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/search [get]
func SearchCharacters(c *gin.Context) {
	searchTerm := c.Query("search")
	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term is required"})
		return
	}

	// Search for character IDs (no need to encode here, the service will handle it)
	characterIDs, err := services.SearchCharactersByName(searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(characterIDs) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Fetch detailed character information
	var characters []interface{}
	for _, characterID := range characterIDs {
		character, err := services.FetchCharacterInfo(characterID)
		if err != nil {
			// Skip characters that couldn't be fetched but continue with others
			continue
		}
		characters = append(characters, character)
	}

	c.JSON(http.StatusOK, characters)
}

// BatchAddCharacters adds multiple characters by ID
// @Summary Batch add characters
// @Description Add multiple characters by their IDs
// @Tags characters
// @Accept json
// @Produce json
// @Param characters body []int64 true "Array of character IDs"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/batch [post]
func BatchAddCharacters(c *gin.Context) {
	var characterIDs []int64
	if err := c.ShouldBindJSON(&characterIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := make(map[string]interface{})
	successCount := 0
	failureCount := 0
	duplicateCount := 0
	var errors []string
	var duplicates []int64

	for _, characterID := range characterIDs {
		// Check if character already exists
		existingCharacter, err := queries.GetCharacterByID(characterID)
		if err == nil && existingCharacter != nil {
			duplicateCount++
			duplicates = append(duplicates, characterID)
			continue
		}

		// Fetch character info from ESI
		character, err := services.FetchCharacterInfo(characterID)
		if err != nil {
			failureCount++
			errors = append(errors, err.Error())
			continue
		}

		// Store character in database
		err = queries.UpsertCharacter(character)
		if err != nil {
			failureCount++
			errors = append(errors, err.Error())
			continue
		}

		successCount++
	}

	results["success_count"] = successCount
	results["failure_count"] = failureCount
	results["duplicate_count"] = duplicateCount
	results["total_processed"] = len(characterIDs)
	if len(errors) > 0 {
		results["errors"] = errors
	}
	if len(duplicates) > 0 {
		results["duplicates"] = duplicates
	}

	c.JSON(http.StatusOK, results)
}

// RemoveCharacter removes a character from the database
// @Summary Remove a character
// @Description Remove a character from the database
// @Tags characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id} [delete]
func RemoveCharacter(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	// Check if character exists
	existingCharacter, err := queries.GetCharacterByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingCharacter == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	// Remove character
	err = queries.DeleteCharacter(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDashboardStats returns statistics for the admin dashboard
// @Summary Get dashboard statistics
// @Description Get various statistics for the admin dashboard
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/stats [get]
func GetDashboardStats(c *gin.Context) {
	stats := make(map[string]interface{})

	// Get total characters
	characters, err := queries.GetAllCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stats["total_characters"] = len(characters)

	// Get total kills (optimized with COUNT query instead of loading all records)
	totalKills, err := queries.GetKillCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stats["total_kills"] = totalKills

	// Get total regions
	regions, err := queries.GetAllRegions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stats["total_regions"] = len(regions)

	// Get recent activity (kills from last 7 days) - optimized with COUNT query
	recentActivity, err := queries.GetRecentKillCount(7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stats["recent_activity"] = recentActivity

	c.JSON(http.StatusOK, stats)
}
