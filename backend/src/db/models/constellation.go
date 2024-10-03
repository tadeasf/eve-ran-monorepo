package models

import "encoding/json"

// Constellation model
type Constellation struct {
	ConstellationID int             `gorm:"primaryKey" json:"constellation_id"`
	Name            string          `json:"name"`
	RegionID        int             `json:"region_id"`
	Systems         json.RawMessage `gorm:"type:jsonb" json:"systems"`
	Position        json.RawMessage `gorm:"type:jsonb" json:"position"`
}
