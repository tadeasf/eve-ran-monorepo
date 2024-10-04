package queries

import (
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/gorm/clause"
)

func ZKillExists(killmailID int64) (bool, error) {
	var count int64
	result := db.DB.Model(&models.Zkill{}).Where("killmail_id = ?", killmailID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func UpsertZKills(zkills []models.Zkill) error {
	return db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "killmail_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"character_id", "location_id", "hash", "fitted_value", "dropped_value", "destroyed_value", "total_value", "points", "npc", "solo", "awox", "labels"}),
	}).Create(&zkills).Error
}

func GetZKillByID(killmailID int64) (*models.Zkill, error) {
	var zkill models.Zkill
	result := db.DB.Where("killmail_id = ?", killmailID).First(&zkill)
	if result.Error != nil {
		return nil, result.Error
	}
	return &zkill, nil
}
