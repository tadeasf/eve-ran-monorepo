package models

import (
	"database/sql/driver"
	"encoding/json"
)

// Region model
type Region struct {
	RegionID       int    `gorm:"primaryKey" json:"region_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Constellations []int  `gorm:"type:jsonb" json:"constellations"`
}

func (a IntArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
