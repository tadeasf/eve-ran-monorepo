package queries

import (
	"fmt"
	"time"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
)

// GetCompetitionSettings retrieves the current competition settings
// Returns default settings if none exist
func GetCompetitionSettings() (*models.CompetitionSettings, error) {
	var settings models.CompetitionSettings
	err := db.DB.Where("active = ?", true).First(&settings).Error
	if err != nil {
		// If no settings exist, return default settings
		if err.Error() == "record not found" {
			return &models.CompetitionSettings{
				Metric:  "kill_count",
				Regions: []int{},
				Active:  true,
			}, nil
		}
		return nil, err
	}

	// Ensure regions is never nil
	if settings.Regions == nil {
		settings.Regions = []int{}
	}

	return &settings, nil
}

// UpsertCompetitionSettings creates or updates competition settings
func UpsertCompetitionSettings(settings *models.CompetitionSettings) error {
	// Ensure regions is never nil
	if settings.Regions == nil {
		settings.Regions = []int{}
	}

	// Deactivate all existing settings first
	if err := db.DB.Model(&models.CompetitionSettings{}).Where("active = ?", true).Update("active", false).Error; err != nil {
		return err
	}

	// Create new settings
	now := time.Now()
	settings.Active = true
	settings.CreatedAt = now
	settings.UpdatedAt = now
	return db.DB.Create(settings).Error
}

// GetCurrentCompetitionStandings returns current month standings based on active settings
// Uses existing kills and zkills data from the database
func GetCurrentCompetitionStandings() ([]models.CompetitionStanding, error) {
	settings, err := GetCompetitionSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get competition settings: %v", err)
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var standings []models.CompetitionStanding

	// Base query joining kills with characters
	query := db.DB.Table("kills").
		Joins("LEFT JOIN characters ON kills.character_id = characters.id").
		Where("kills.killmail_time >= ? AND kills.killmail_time < ?", startOfMonth, endOfMonth)

	// Filter by regions if specified (non-empty array means specific regions)
	if len(settings.Regions) > 0 {
		query = query.Joins("LEFT JOIN systems ON kills.solar_system_id = systems.system_id").
			Where("systems.region_id IN ?", settings.Regions)
	}

	// Add appropriate aggregation based on metric
	if settings.Metric == "kill_count" {
		query = query.Select("kills.character_id, characters.name as character_name, COUNT(*) as value").
			Group("kills.character_id, characters.name").
			Order("value DESC")
	} else { // isk_destroyed
		query = query.Select("kills.character_id, characters.name as character_name, COALESCE(SUM(zkills.total_value), 0) as value").
			Joins("LEFT JOIN zkills ON kills.killmail_id = zkills.killmail_id").
			Group("kills.character_id, characters.name").
			Having("COALESCE(SUM(zkills.total_value), 0) > 0").
			Order("value DESC")
	}

	if err := query.Scan(&standings).Error; err != nil {
		return nil, err
	}

	// Assign ranks
	for i := range standings {
		standings[i].Rank = i + 1
	}

	return standings, nil
}

// SaveMonthlyCompetitionResults saves the final results for a given month
// Uses existing kills and zkills data from the database
func SaveMonthlyCompetitionResults(month, year int) error {
	settings, err := GetCompetitionSettings()
	if err != nil {
		return fmt.Errorf("failed to get competition settings: %v", err)
	}

	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var standings []models.CompetitionStanding

	// Base query joining kills with characters
	query := db.DB.Table("kills").
		Joins("LEFT JOIN characters ON kills.character_id = characters.id").
		Where("kills.killmail_time >= ? AND kills.killmail_time < ?", startOfMonth, endOfMonth)

	// Filter by regions if specified
	if len(settings.Regions) > 0 {
		query = query.Joins("LEFT JOIN systems ON kills.solar_system_id = systems.system_id").
			Where("systems.region_id IN ?", settings.Regions)
	}

	// Add appropriate aggregation based on metric
	if settings.Metric == "kill_count" {
		query = query.Select("kills.character_id, characters.name as character_name, COUNT(*) as value").
			Group("kills.character_id, characters.name").
			Order("value DESC").
			Limit(10) // Save top 10
	} else { // isk_destroyed
		query = query.Select("kills.character_id, characters.name as character_name, COALESCE(SUM(zkills.total_value), 0) as value").
			Joins("LEFT JOIN zkills ON kills.killmail_id = zkills.killmail_id").
			Group("kills.character_id, characters.name").
			Having("COALESCE(SUM(zkills.total_value), 0) > 0").
			Order("value DESC").
			Limit(10) // Save top 10
	}

	if err := query.Scan(&standings).Error; err != nil {
		return err
	}

	// Delete existing results for this month/year to avoid duplicates
	if err := db.DB.Where("month = ? AND year = ?", month, year).Delete(&models.CompetitionResult{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing results: %v", err)
	}

	// Save results to database
	for i, standing := range standings {
		result := models.CompetitionResult{
			Month:       month,
			Year:        year,
			CharacterID: standing.CharacterID,
			Metric:      settings.Metric,
			Value:       standing.Value,
			Rank:        i + 1,
			Regions:     settings.Regions,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.DB.Create(&result).Error; err != nil {
			return fmt.Errorf("failed to save competition result: %v", err)
		}
	}

	return nil
}

// GetCompetitionHistory returns historical competition results for a specific month/year
func GetCompetitionHistory(month, year int) ([]models.CompetitionWinner, error) {
	var winners []models.CompetitionWinner

	err := db.DB.Table("competition_results").
		Select("competition_results.character_id, characters.name as character_name, competition_results.month, competition_results.year, competition_results.metric, competition_results.value, competition_results.rank").
		Joins("LEFT JOIN characters ON competition_results.character_id = characters.id").
		Where("competition_results.month = ? AND competition_results.year = ?", month, year).
		Order("competition_results.rank ASC").
		Scan(&winners).Error

	if err != nil {
		return nil, err
	}

	return winners, nil
}

// GetRecentCompetitionWinners returns the most recent competition winners (last month)
func GetRecentCompetitionWinners() ([]models.CompetitionWinner, error) {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	return GetCompetitionHistory(int(lastMonth.Month()), lastMonth.Year())
}

// GetAllCompetitionHistory returns all historical competition results
func GetAllCompetitionHistory() ([]models.CompetitionWinner, error) {
	var winners []models.CompetitionWinner

	err := db.DB.Table("competition_results").
		Select("competition_results.character_id, characters.name as character_name, competition_results.month, competition_results.year, competition_results.metric, competition_results.value, competition_results.rank").
		Joins("LEFT JOIN characters ON competition_results.character_id = characters.id").
		Order("competition_results.year DESC, competition_results.month DESC, competition_results.rank ASC").
		Scan(&winners).Error

	if err != nil {
		return nil, err
	}

	return winners, nil
}
