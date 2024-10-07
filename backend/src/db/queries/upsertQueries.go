package queries

import (
	"encoding/json"

	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/gorm"
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
		Columns: []clause.Column{{Name: "killmail_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"killmail_time",
			"solar_system_id",
			"victim_alliance_id",
			"victim_character_id",
			"victim_corporation_id",
			"victim_damage_taken",
			"victim_ship_type_id",
			"victim_position_x",
			"victim_position_y",
			"victim_position_z",
			"victim_items",
			"attackers",
		}),
	}).Create(kill).Error
}

func BatchUpsertSystems(systems []*models.System) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, system := range systems {
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "system_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"constellation_id", "name", "security_class", "security_status", "star_id", "planets", "stargates", "stations", "position"}),
			}).Create(system).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func UpsertRegion(region *models.Region) error {
	constellationsJSON, err := json.Marshal(region.Constellations)
	if err != nil {
		return err
	}

	return db.DB.Exec(`
        INSERT INTO regions (region_id, name, description, constellations)
        VALUES (?, ?, ?, ?)
        ON CONFLICT (region_id) DO UPDATE
        SET name = EXCLUDED.name,
            description = EXCLUDED.description,
            constellations = EXCLUDED.constellations
    `, region.RegionID, region.Name, region.Description, constellationsJSON).Error
}

func BatchUpsertConstellations(constellations []*models.Constellation) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, constellation := range constellations {
			systemsJSON, err := json.Marshal(constellation.Systems)
			if err != nil {
				return err
			}

			err = tx.Exec(`
				INSERT INTO constellations (constellation_id, name, region_id, systems, position)
				VALUES (?, ?, ?, ?, ?)
				ON CONFLICT (constellation_id) DO UPDATE
				SET name = EXCLUDED.name,
					region_id = EXCLUDED.region_id,
					systems = EXCLUDED.systems,
					position = EXCLUDED.position
			`, constellation.ConstellationID, constellation.Name, constellation.RegionID, systemsJSON, constellation.Position).Error

			if err != nil {
				return err
			}
		}
		return nil
	})
}

func UpsertESIItem(item *models.ESIItem) error {
	return db.DB.Save(item).Error
}

func BulkUpsertKills(kills []models.Kill) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, kill := range kills {
			result := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "killmail_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"character_id", "kill_time", "solar_system_id", "location_id", "hash", "fitted_value", "dropped_value", "destroyed_value", "total_value", "points", "npc", "solo", "awox", "victim_alliance_id", "victim_character_id", "victim_corporation_id", "victim_faction_id", "victim_damage_taken", "victim_ship_type_id", "victim_items", "victim_position", "attackers"}),
			}).Create(&kill)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}
