package db

import (
	"fmt"
	"log"
)

// RunIndexMigrations creates database indexes for performance optimization
func RunIndexMigrations() error {
	log.Println("Running index migrations...")

	// Create indexes on kills table for query performance
	indexes := []struct {
		name  string
		table string
		sql   string
	}{
		{
			name:  "idx_kills_solar_system_id",
			table: "kills",
			sql:   "CREATE INDEX IF NOT EXISTS idx_kills_solar_system_id ON kills (solar_system_id)",
		},
		{
			name:  "idx_kills_killmail_time",
			table: "kills",
			sql:   "CREATE INDEX IF NOT EXISTS idx_kills_killmail_time ON kills (killmail_time)",
		},
		{
			name:  "idx_kills_character_id",
			table: "kills",
			sql:   "CREATE INDEX IF NOT EXISTS idx_kills_character_id ON kills (character_id)",
		},
		{
			name:  "idx_kills_system_time",
			table: "kills",
			sql:   "CREATE INDEX IF NOT EXISTS idx_kills_system_time ON kills (solar_system_id, killmail_time)",
		},
		{
			name:  "idx_kills_character_time",
			table: "kills",
			sql:   "CREATE INDEX IF NOT EXISTS idx_kills_character_time ON kills (character_id, killmail_time DESC)",
		},
		{
			name:  "idx_kills_victim_ship_type",
			table: "kills",
			sql:   "CREATE INDEX IF NOT EXISTS idx_kills_victim_ship_type ON kills (victim_ship_type_id)",
		},
		{
			name:  "idx_zkills_character_id",
			table: "zkills",
			sql:   "CREATE INDEX IF NOT EXISTS idx_zkills_character_id ON zkills (character_id)",
		},
		{
			name:  "idx_zkills_killmail_id",
			table: "zkills",
			sql:   "CREATE INDEX IF NOT EXISTS idx_zkills_killmail_id ON zkills (killmail_id)",
		},
		{
			name:  "idx_competition_results_year_month",
			table: "competition_results",
			sql:   "CREATE INDEX IF NOT EXISTS idx_competition_results_year_month ON competition_results (year, month)",
		},
		{
			name:  "idx_competition_results_character",
			table: "competition_results",
			sql:   "CREATE INDEX IF NOT EXISTS idx_competition_results_character ON competition_results (character_id)",
		},
		{
			name:  "idx_characters_name",
			table: "characters",
			sql:   "CREATE INDEX IF NOT EXISTS idx_characters_name ON characters (name)",
		},
		{
			name:  "idx_systems_region_id",
			table: "systems",
			sql:   "CREATE INDEX IF NOT EXISTS idx_systems_region_id ON systems (region_id)",
		},
		{
			name:  "idx_systems_constellation_id",
			table: "systems",
			sql:   "CREATE INDEX IF NOT EXISTS idx_systems_constellation_id ON systems (constellation_id)",
		},
	}

	for _, idx := range indexes {
		log.Printf("Creating index %s on table %s...", idx.name, idx.table)
		if err := DB.Exec(idx.sql).Error; err != nil {
			log.Printf("Warning: Failed to create index %s: %v", idx.name, err)
			// Don't fail if index already exists, just continue
		} else {
			log.Printf("Successfully created index %s", idx.name)
		}
	}

	// Analyze tables to update statistics for query planner
	log.Println("Analyzing kills table to update statistics...")
	if err := DB.Exec("ANALYZE kills").Error; err != nil {
		return fmt.Errorf("failed to analyze kills table: %v", err)
	}

	log.Println("Index migrations completed successfully")
	return nil
}

// RunDataMigrations handles one-time data migrations
func RunDataMigrations() error {
	log.Println("Running data migrations...")

	// List of tables that need migration from integer[] to jsonb
	tablesToMigrate := []struct {
		tableName  string
		columnName string
	}{
		{"competition_settings", "regions"},
		{"competition_results", "regions"},
	}

	for _, table := range tablesToMigrate {
		log.Printf("Checking %s.%s column type...", table.tableName, table.columnName)

		// Check if the table exists first
		var tableExists bool
		err := DB.Raw(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_name = ?
			)
		`, table.tableName).Scan(&tableExists).Error

		if err != nil || !tableExists {
			log.Printf("Table %s doesn't exist yet, will be created by AutoMigrate", table.tableName)
			continue
		}

		// Check if the column exists with integer[] type
		var columnType string
		err = DB.Raw(`
			SELECT data_type 
			FROM information_schema.columns 
			WHERE table_name = ? 
			AND column_name = ?
		`, table.tableName, table.columnName).Scan(&columnType).Error

		if err == nil && columnType == "ARRAY" {
			log.Printf("Found integer[] type for %s.%s, converting to jsonb...", table.tableName, table.columnName)

			// Check if there's any data
			var hasData bool
			DB.Raw(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s LIMIT 1)", table.tableName)).Scan(&hasData)

			if hasData {
				log.Printf("Warning: Found existing data in %s, it will be reset", table.tableName)
			}

			// Drop the table - CASCADE will handle any foreign keys
			err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table.tableName)).Error
			if err != nil {
				return fmt.Errorf("failed to drop %s table: %v", table.tableName, err)
			}

			log.Printf("Successfully dropped %s table, AutoMigrate will recreate it", table.tableName)
		}
	}

	log.Println("Data migrations completed successfully")
	return nil
}
