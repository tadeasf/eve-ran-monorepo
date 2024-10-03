package models

import (
	"encoding/json"
)

// Region model
type Region struct {
	RegionID       int             `gorm:"primaryKey" json:"region_id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Constellations json.RawMessage `gorm:"type:jsonb" json:"constellations"`
}
