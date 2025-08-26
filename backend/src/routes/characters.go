package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// GetAllKills retrieves all kills from the database
// @Summary Get all kills
// @Description Fetch all kills from the database
// @Tags kills
// @Accept json
// @Produce json
// @Success 200 {array} models.Kill
// @Failure 500 {object} models.ErrorResponse
// @Router /kills [get]
func GetAllKills(c *gin.Context) {
	kills, err := queries.GetAllKills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kills)
}

// GetAllCharacterStats retrieves stats for all characters with filters
// @Summary Get all character stats
// @Description Fetch stats for all characters from the database with optional filters
// @Tags characters
// @Accept json
// @Produce json
// @Param regionID query []int false "Region IDs"
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
	}
	if endDate != "" {
		endTime, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	}

	stats, err := queries.GetCharacterStats(startTime, endTime, 0, regionIDInts...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
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

	// Get total kills
	kills, err := queries.GetAllKills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stats["total_kills"] = len(kills)

	// Get total regions
	regions, err := queries.GetAllRegions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stats["total_regions"] = len(regions)

	// Get recent activity (kills from last 7 days)
	// This is a simplified version - you might want to make this more sophisticated
	recentCount := 0
	for _, kill := range kills {
		if time.Since(kill.KillmailTime).Hours() < 168 { // 7 days * 24 hours
			recentCount++
		}
	}
	stats["recent_activity"] = recentCount

	c.JSON(http.StatusOK, stats)
}
