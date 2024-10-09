package models

// CharacterStats model
type CharacterStats struct {
	CharacterID int64   `json:"character_id"`
	KillCount   int     `json:"kill_count"`
	TotalISK    float64 `json:"total_isk"`
}

// Character model
type Character struct {
	ID             int64   `gorm:"primaryKey" json:"id"`
	Name           string  `json:"name"`
	SecurityStatus float64 `json:"security_status"`
	Title          string  `json:"title"`
	RaceID         int     `json:"race_id"`
}
