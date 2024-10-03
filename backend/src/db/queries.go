package db

import (
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/gorm/clause"
)

func GetConstellation(id int) (*models.Constellation, error) {
	var constellation models.Constellation
	result := DB.First(&constellation, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &constellation, nil
}

func GetConstellationsByRegionID(regionID int) ([]models.Constellation, error) {
	var constellations []models.Constellation
	result := DB.Where("region_id = ?", regionID).Find(&constellations)
	if result.Error != nil {
		return nil, result.Error
	}
	return constellations, nil
}

func GetSystem(id int) (*models.System, error) {
	var system models.System
	result := DB.First(&system, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &system, nil
}

func GetSystemsByRegionID(regionID int) ([]models.System, error) {
	var systems []models.System
	result := DB.Where("region_id = ?", regionID).Find(&systems)
	if result.Error != nil {
		return nil, result.Error
	}
	return systems, nil
}

func GetKillsForCharacterWithFilters(characterID int64, page, pageSize, regionID int, startDate, endDate string) ([]models.Kill, error) {
	var kills []models.Kill
	query := DB.Where("character_id = ?", characterID)

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

func GetTotalKillsForCharacterWithFilters(characterID int64, regionID int, startDate, endDate string) (int64, error) {
	var count int64
	query := DB.Model(&models.Kill{}).Where("character_id = ?", characterID)

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

func GetKillsByRegion(regionID int, page, pageSize int, startDate, endDate string) ([]models.Kill, int64, error) {
	var kills []models.Kill
	var totalCount int64

	query := DB.Table("kills").
		Joins("JOIN systems ON kills.solar_system_id = systems.system_id").
		Joins("JOIN constellations ON systems.constellation_id = constellations.constellation_id").
		Where("constellations.region_id = ?", regionID)

	if startDate != "" {
		query = query.Where("kills.kill_time >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("kills.kill_time <= ?", endDate)
	}

	// Count total matching kills
	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	err = query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&kills).Error

	if err != nil {
		return nil, 0, err
	}

	return kills, totalCount, nil
}

func GetTotalKillsByRegion(regionID int, startDate, endDate string) (int64, error) {
	var count int64
	query := DB.Model(&models.Kill{}).Where("solar_system_id IN (SELECT system_id FROM systems WHERE region_id = ?)", regionID)

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

func GetCharacterKillmails(characterID int64, startTime, endTime time.Time, systemID, regionID int64) ([]models.Kill, error) {
	query := DB.Where("character_id = ? AND kill_time BETWEEN ? AND ?", characterID, startTime, endTime)

	if systemID != 0 {
		query = query.Where("solar_system_id = ?", systemID)
	}

	if regionID != 0 {
		query = query.Joins("JOIN systems ON kills.solar_system_id = systems.system_id").
			Where("systems.region_id = ?", regionID)
	}

	var kills []models.Kill
	err := query.Find(&kills).Error
	return kills, err
}

type CharacterStats struct {
	CharacterID int64   `json:"character_id"`
	KillCount   int     `json:"kill_count"`
	TotalISK    float64 `json:"total_isk"`
}

func GetCharacterStats(startTime, endTime time.Time, systemID int64, regionIDs ...int64) ([]CharacterStats, error) {
	query := DB.Table("kills").
		Select("character_id, COUNT(*) as kill_count, SUM(total_value) as total_isk").
		Where("kill_time BETWEEN ? AND ?", startTime, endTime).
		Group("character_id")

	if systemID != 0 {
		query = query.Where("solar_system_id = ?", systemID)
	}

	if len(regionIDs) > 0 {
		query = query.Joins("JOIN systems ON kills.solar_system_id = systems.system_id").
			Where("systems.region_id IN ?", regionIDs)
	}

	var stats []CharacterStats
	err := query.Find(&stats).Error
	return stats, err
}

func UpsertKillsBatch(kills []*models.Kill) error {
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "killmail_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"character_id", "kill_time", "solar_system_id", "location_id", "hash", "fitted_value", "dropped_value", "destroyed_value", "total_value", "points", "npc", "solo", "awox", "victim_alliance_id", "victim_character_id", "victim_corporation_id", "victim_faction_id", "victim_damage_taken", "victim_ship_type_id", "victim_items", "victim_position", "attackers"}),
	}).Create(kills).Error
}

func GetKillsForCharacter(characterID int64, page, pageSize int) ([]models.Kill, error) {
	var kills []models.Kill
	offset := (page - 1) * pageSize
	err := DB.Where("character_id = ?", characterID).Order("kill_time DESC").Offset(offset).Limit(pageSize).Find(&kills).Error
	return kills, err
}

func GetTotalKillsForCharacter(characterID int64) (int64, error) {
	var count int64
	err := DB.Model(&models.Kill{}).Where("character_id = ?", characterID).Count(&count).Error
	return count, err
}
