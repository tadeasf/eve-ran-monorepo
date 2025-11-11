package models

import (
	"time"
)

// CompetitionSettings stores the configuration for competitions
type CompetitionSettings struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Metric    string    `gorm:"type:varchar(50);not null" json:"metric"` // "isk_destroyed" or "kill_count"
	Regions   []int     `gorm:"type:integer[];not null" json:"regions"`  // Array of region IDs, empty means all regions
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CompetitionResult stores monthly competition results
type CompetitionResult struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Month       int       `gorm:"not null;index:idx_month_year,priority:1" json:"month"` // 1-12
	Year        int       `gorm:"not null;index:idx_month_year,priority:2" json:"year"`
	CharacterID int64     `gorm:"not null;index:idx_character_competition" json:"character_id"`
	Metric      string    `gorm:"type:varchar(50);not null" json:"metric"` // "isk_destroyed" or "kill_count"
	Value       float64   `gorm:"not null" json:"value"`                   // Either ISK or kill count
	Rank        int       `gorm:"not null" json:"rank"`                    // 1 = winner, 2 = runner-up, etc.
	Regions     []int     `gorm:"type:integer[]" json:"regions"`           // Which regions were included
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CompetitionStanding represents current (live) competition standings
type CompetitionStanding struct {
	CharacterID   int64   `json:"character_id"`
	CharacterName string  `json:"character_name"`
	Value         float64 `json:"value"` // ISK or kill count
	Rank          int     `json:"rank"`
}

// CompetitionWinner represents a monthly winner with additional info
type CompetitionWinner struct {
	CharacterID   int64   `json:"character_id"`
	CharacterName string  `json:"character_name"`
	Month         int     `json:"month"`
	Year          int     `json:"year"`
	Metric        string  `json:"metric"`
	Value         float64 `json:"value"`
	Rank          int     `json:"rank"`
}
