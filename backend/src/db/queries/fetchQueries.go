package queries

import (
	"time"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
)

func FetchKillsForCharacterWithFilters(characterID int64, page, pageSize, regionID int, startDate, endDate string) ([]models.Kill, error) {
	var kills []models.Kill
	query := db.DB.Where("character_id = ?", characterID)

	if regionID != 0 {
		query = query.Where("solar_system_id IN (SELECT system_id FROM systems WHERE region_id = ?)", regionID)
	}

	if startDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		query = query.Where("kill_time >= ?", startTime)
	}

	if endDate != "" {
		endTime, _ := time.Parse("2006-01-02", endDate)
		query = query.Where("kill_time <= ?", endTime)
	}

	result := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&kills)
	if result.Error != nil {
		return nil, result.Error
	}
	return kills, nil
}

func FetchTotalKillsForCharacterWithFilters(characterID int64, regionID int, startDate, endDate string) (int64, error) {
	var count int64
	query := db.DB.Model(&models.Kill{}).Where("character_id = ?", characterID)

	if regionID != 0 {
		query = query.Where("solar_system_id IN (SELECT system_id FROM systems WHERE region_id = ?)", regionID)
	}

	if startDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		query = query.Where("kill_time >= ?", startTime)
	}

	if endDate != "" {
		endTime, _ := time.Parse("2006-01-02", endDate)
		query = query.Where("kill_time <= ?", endTime)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func FetchKillsByRegion(regionID int, page, pageSize int, startDate, endDate string) ([]models.Kill, int64, error) {
	var kills []models.Kill
	var totalCount int64

	query := db.DB.Table("kills").
		Joins("JOIN systems ON kills.solar_system_id = systems.system_id").
		Joins("JOIN constellations ON systems.constellation_id = constellations.constellation_id").
		Where("constellations.region_id = ?", regionID)

	if startDate != "" {
		query = query.Where("kills.kill_time >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("kills.kill_time <= ?", endDate)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&kills).Error

	if err != nil {
		return nil, 0, err
	}

	return kills, totalCount, nil
}

// Add other fetch queries as needed
