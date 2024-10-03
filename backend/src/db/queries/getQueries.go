package queries

import (
	"errors"
	"time"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found")

func GetCharacterByID(id int64) (*models.Character, error) {
	var character models.Character
	result := db.DB.First(&character, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &character, nil
}

func GetAllCharacters() ([]models.Character, error) {
	var characters []models.Character
	err := db.DB.Find(&characters).Error
	return characters, err
}

func GetKillByKillmailID(killmailID int64) (*models.Kill, error) {
	var kill models.Kill
	err := db.DB.First(&kill, killmailID).Error
	return &kill, err
}

func GetAllKills() ([]models.Kill, error) {
	var kills []models.Kill
	err := db.DB.Find(&kills).Error
	return kills, err
}

func GetKillsForCharacter(characterID int64, page, pageSize int) ([]models.Kill, error) {
	var kills []models.Kill
	offset := (page - 1) * pageSize
	err := db.DB.Where("character_id = ?", characterID).Order("kill_time DESC").Offset(offset).Limit(pageSize).Find(&kills).Error
	return kills, err
}

func GetTotalKillsForCharacter(characterID int64) (int64, error) {
	var count int64
	err := db.DB.Model(&models.Kill{}).Where("character_id = ?", characterID).Count(&count).Error
	return count, err
}

func GetCharacterStats(startTime, endTime time.Time, systemID int64, regionIDs ...int64) ([]models.CharacterStats, error) {
	query := db.DB.Table("kills").
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

	var stats []models.CharacterStats
	err := query.Find(&stats).Error
	return stats, err
}

func IsInitialFetchForCharacter(characterID int64) (bool, error) {
	var count int64
	err := db.DB.Model(&models.Kill{}).Where("character_id = ?", characterID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func GetKillByID(killmailID int64) (*models.Kill, error) {
	var kill models.Kill
	err := db.DB.Where("killmail_id = ?", killmailID).First(&kill).Error
	return &kill, err
}

func GetUnenrichedKillsForCharacter(characterID int64) ([]models.Kill, error) {
	var kills []models.Kill
	err := db.DB.Where("character_id = ? AND kill_time IS NULL", characterID).Find(&kills).Error
	return kills, err
}

func GetAllESIItems() ([]models.ESIItem, error) {
	var items []models.ESIItem
	err := db.DB.Find(&items).Error
	return items, err
}

func GetESIItemByTypeID(typeID int) (*models.ESIItem, error) {
	var item models.ESIItem
	err := db.DB.First(&item, typeID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func GetAllConstellations() ([]models.Constellation, error) {
	var constellations []models.Constellation
	err := db.DB.Find(&constellations).Error
	return constellations, err
}
