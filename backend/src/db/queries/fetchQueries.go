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

package queries

import (
	"time"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
)

func GetAllSystems() ([]models.System, error) {
	var systems []models.System
	err := db.DB.Find(&systems).Error
	return systems, err
}

func GetSystemByID(systemID int) (*models.System, error) {
	var system models.System
	err := db.DB.First(&system, systemID).Error
	if err != nil {
		return nil, err
	}
	return &system, nil
}

func GetSystemsByRegionID(regionID int) ([]models.System, error) {
	var systems []models.System
	err := db.DB.Where("region_id = ?", regionID).Find(&systems).Error
	return systems, err
}

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
		Where("systems.region_id = ?", regionID)

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

func GetLastKillTimeForCharacter(characterID int64) (time.Time, error) {
	var lastKill models.Kill
	err := db.DB.Where("character_id = ?", characterID).Order("kill_time DESC").First(&lastKill).Error
	if err != nil {
		return time.Time{}, err
	}
	return lastKill.KillmailTime, nil
}

func GetAllRegions() ([]models.Region, error) {
	var regions []models.Region
	err := db.DB.Find(&regions).Error
	if err != nil {
		return nil, err
	}
	return regions, nil
}

func GetKillsByRegion(regionID int, startDate, endDate string) ([]models.Kill, error) {
	var kills []models.Kill

	query := db.DB.Table("kills").
		Joins("JOIN systems ON kills.solar_system_id = systems.system_id").
		Joins("JOIN constellations ON systems.constellation_id = constellations.constellation_id").
		Where("constellations.region_id = ?", regionID)

	if startDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		query = query.Where("kills.killmail_time >= ?", startTime)
	}
	if endDate != "" {
		endTime, _ := time.Parse("2006-01-02", endDate)
		query = query.Where("kills.killmail_time <= ?", endTime)
	}

	err := query.Find(&kills).Error
	if err != nil {
		return nil, err
	}

	return kills, nil
}
