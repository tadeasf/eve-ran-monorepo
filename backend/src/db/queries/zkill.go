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
