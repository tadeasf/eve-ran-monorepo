package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db/models"
	"github.com/tadeasf/eve-ran/src/db/queries"
)

// GetCompetitionSettings retrieves the current competition settings
// @Summary Get competition settings
// @Description Get the currently active competition settings (returns defaults if none exist)
// @Tags competition
// @Produce json
// @Success 200 {object} models.CompetitionSettings
// @Failure 500 {object} map[string]string
// @Router /competition/settings [get]
func GetCompetitionSettings(c *gin.Context) {
	settings, err := queries.GetCompetitionSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get competition settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateCompetitionSettings creates or updates competition settings
// @Summary Update competition settings
// @Description Create or update the competition settings (admin only)
// @Tags competition
// @Accept json
// @Produce json
// @Param settings body models.CompetitionSettings true "Competition settings"
// @Success 200 {object} models.CompetitionSettings
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /competition/settings [post]
func UpdateCompetitionSettings(c *gin.Context) {
	var settings models.CompetitionSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate metric
	if settings.Metric != "isk_destroyed" && settings.Metric != "kill_count" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric. Must be 'isk_destroyed' or 'kill_count'"})
		return
	}

	if err := queries.UpsertCompetitionSettings(&settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// GetCurrentCompetitionStandings returns current month competition standings
// @Summary Get current competition standings
// @Description Get the current month's competition standings based on active settings
// @Tags competition
// @Produce json
// @Success 200 {array} models.CompetitionStanding
// @Failure 500 {object} map[string]string
// @Router /competition/current [get]
func GetCurrentCompetitionStandings(c *gin.Context) {
	standings, err := queries.GetCurrentCompetitionStandings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, standings)
}

// GetCompetitionHistory returns historical competition results for a specific month
// @Summary Get competition history
// @Description Get historical competition results for a specific month and year
// @Tags competition
// @Produce json
// @Param month path int true "Month (1-12)"
// @Param year path int true "Year"
// @Success 200 {array} models.CompetitionWinner
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /competition/history/{month}/{year} [get]
func GetCompetitionHistory(c *gin.Context) {
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil || year < 2020 || year > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	history, err := queries.GetCompetitionHistory(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetRecentCompetitionWinners returns the most recent competition winners
// @Summary Get recent competition winners
// @Description Get the winners from the most recent completed competition (last month)
// @Tags competition
// @Produce json
// @Success 200 {array} models.CompetitionWinner
// @Failure 500 {object} map[string]string
// @Router /competition/recent-winners [get]
func GetRecentCompetitionWinners(c *gin.Context) {
	winners, err := queries.GetRecentCompetitionWinners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, winners)
}

// GetAllCompetitionHistory returns all historical competition results
// @Summary Get all competition history
// @Description Get all historical competition results across all months
// @Tags competition
// @Produce json
// @Success 200 {array} models.CompetitionWinner
// @Failure 500 {object} map[string]string
// @Router /competition/history [get]
func GetAllCompetitionHistory(c *gin.Context) {
	history, err := queries.GetAllCompetitionHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// SaveMonthlyResults manually triggers saving the monthly competition results
// @Summary Save monthly competition results
// @Description Manually save competition results for a specific month (admin only)
// @Tags competition
// @Produce json
// @Param month path int true "Month (1-12)"
// @Param year path int true "Year"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /competition/save/{month}/{year} [post]
func SaveMonthlyResults(c *gin.Context) {
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil || year < 2020 || year > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	if err := queries.SaveMonthlyCompetitionResults(month, year); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Competition results saved successfully"})
}
