package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/jobs"
)

func GetCharacterKillmails(c *gin.Context) {
	characterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	startTime, _ := time.Parse(time.RFC3339, c.Query("start_time"))
	endTime, _ := time.Parse(time.RFC3339, c.Query("end_time"))
	systemID, _ := strconv.ParseInt(c.Query("system_id"), 10, 64)
	regionID, _ := strconv.ParseInt(c.Query("region_id"), 10, 64)

	kills, err := queries.GetCharacterKillmails(characterID, startTime, endTime, systemID, regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kills)
}

// TriggerKillFetchForCharacter triggers the kill fetcher job for a specific character
func TriggerKillFetchForCharacter(c *gin.Context) {
	characterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	go jobs.FetchAllKillsForCharacter(characterID)

	c.JSON(http.StatusOK, gin.H{"message": "Kill fetch job triggered for character " + c.Param("id")})
}

// TriggerKillFetchForAllCharacters triggers the kill fetcher job for all characters
func TriggerKillFetchForAllCharacters(c *gin.Context) {
	go jobs.FetchKillsForAllCharacters()

	c.JSON(http.StatusOK, gin.H{"message": "Kill fetch job triggered for all characters"})
}
