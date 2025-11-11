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
