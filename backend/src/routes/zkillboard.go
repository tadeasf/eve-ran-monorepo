package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/jobs"
	"github.com/tadeasf/eve-ran/src/utils"
	"gorm.io/gorm"
)

var addCharacterMutex sync.Mutex

// AddCharacter adds a new character ID
// @Summary Add a new character ID
// @Description Add a new character ID to the database and fetch all kills
// @Tags characters
// @Accept json
// @Produce json
// @Param character body models.Character true "Character ID"
// @Success 201 {object} models.Character
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters [post]
func AddCharacter(c *gin.Context) {
	addCharacterMutex.Lock()
	defer addCharacterMutex.Unlock()

	var character models.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingCharacter, err := queries.GetCharacterByID(character.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing character"})
		return
	}

	if existingCharacter != nil {
		c.JSON(http.StatusOK, existingCharacter)
		return
	}

	// Fetch character data from ESI API
	esiURL := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/?datasource=tranquility", character.ID)
	req, _ := http.NewRequest("GET", esiURL, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Cache-Control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch character data from ESI"})
		return
	}
	defer resp.Body.Close()

	var esiData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&esiData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ESI response"})
		return
	}

	// Update character with ESI data
	character.Name = esiData["name"].(string)
	character.SecurityStatus = esiData["security_status"].(float64)
	if title, ok := esiData["title"].(string); ok {
		character.Title = title
	}
	character.RaceID = int(esiData["race_id"].(float64))

	// Insert the character into the database
	if err := queries.UpsertCharacter(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add character"})
		return
	}

	// Trigger kill initialization
	go func() {
		if err := jobs.InitializeCharacterKills(character.ID); err != nil {
			utils.LogError(fmt.Sprintf("Error initializing kills for character %d: %v", character.ID, err))
		}
	}()

	utils.LogToConsole(fmt.Sprintf("Added character: %s (ID: %d)", character.Name, character.ID))
	c.JSON(http.StatusCreated, character)

}

// RemoveCharacter removes a character
// @Summary Remove a character
// @Description Remove a character from the database
// @Tags characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id} [delete]
func RemoveCharacter(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	err = db.DB.Delete(&models.Character{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCharacterKillsFromDB retrieves character kills from the database
// @Summary Get character kills from database
// @Description Fetch kills for a character from the database
// @Tags characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id}/kills/db [get]
func GetCharacterKillsFromDB(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	kills, err := queries.GetKillsForCharacter(id, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalItems, err := queries.GetTotalKillsForCharacter(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	response := models.PaginatedResponse{
		Data:       kills,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetKillsByRegion retrieves kills by region
// @Summary Get kills by region
// @Description Fetch all kills for a region from the database
// @Tags kills
// @Accept json
// @Produce json
// @Param regionID path int true "Region ID"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} models.Kill
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /kills/region/{regionID} [get]
func GetKillsByRegion(c *gin.Context) {
	regionID, err := strconv.Atoi(c.Param("regionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	kills, err := queries.GetKillsByRegion(regionID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kills)
}
