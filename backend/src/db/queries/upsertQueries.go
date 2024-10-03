package queries

import (
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/gorm/clause"
)

func UpsertCharacter(character *models.Character) error {
	return db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "security_status", "title", "race_id"}),
	}).Create(character).Error
}

func UpsertKill(kill *models.Kill) error {
	return db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "killmail_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"character_id", "kill_time", "solar_system_id", "location_id", "hash", "fitted_value", "dropped_value", "destroyed_value", "total_value", "points", "npc", "solo", "awox", "victim_alliance_id", "victim_character_id", "victim_corporation_id", "victim_faction_id", "victim_damage_taken", "victim_ship_type_id", "victim_items", "victim_position", "attackers"}),
	}).Create(kill).Error
}

func UpsertKillsBatch(kills []*models.Kill) error {
	return db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "killmail_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"character_id", "kill_time", "solar_system_id", "location_id", "hash", "fitted_value", "dropped_value", "destroyed_value", "total_value", "points", "npc", "solo", "awox", "victim_alliance_id", "victim_character_id", "victim_corporation_id", "victim_faction_id", "victim_damage_taken", "victim_ship_type_id", "victim_items", "victim_position", "attackers"}),
	}).Create(kills).Error
}

func InsertKill(kill *models.Kill) error {
	return db.DB.Create(kill).Error
}

func UpdateKill(kill *models.Kill) error {
	return db.DB.Save(kill).Error
}
