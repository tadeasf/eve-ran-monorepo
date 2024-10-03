package models

import "encoding/json"

// System model
type System struct {
	SystemID        int             `gorm:"primaryKey" json:"system_id"`
	ConstellationID int             `json:"constellation_id"`
	Name            string          `json:"name"`
	SecurityClass   string          `json:"security_class"`
	SecurityStatus  float64         `json:"security_status"`
	StarID          int             `json:"star_id"`
	Planets         json.RawMessage `gorm:"type:jsonb" json:"planets"`
	Stargates       json.RawMessage `gorm:"type:jsonb" json:"stargates"`
	Stations        json.RawMessage `gorm:"type:jsonb" json:"stations"`
	Position        json.RawMessage `gorm:"type:jsonb" json:"position"`
}
