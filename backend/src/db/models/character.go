package models

// CharacterStats model
type CharacterStats struct {
	CharacterID int64   `json:"character_id"`
	KillCount   int     `json:"kill_count"`
	TotalISK    float64 `json:"total_isk"`
}
