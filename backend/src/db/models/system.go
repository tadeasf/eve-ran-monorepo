package models

import "encoding/json"

// System model
type System struct {
	SystemID        int             `gorm:"primaryKey" json:"system_id"`
	ConstellationID int             `json:"constellation_id"`
	RegionID        int             `json:"region_id"`
	Name            string          `json:"name"`
	SecurityClass   string          `json:"security_class"`
	SecurityStatus  float64         `json:"security_status"`
	StarID          int             `json:"star_id"`
	Planets         json.RawMessage `gorm:"type:jsonb" json:"planets" swaggertype:"array,integer"`
	Stargates       json.RawMessage `gorm:"type:jsonb" json:"stargates" swaggertype:"array,integer"`
	Stations        json.RawMessage `gorm:"type:jsonb" json:"stations" swaggertype:"array,integer"`
	Position        json.RawMessage `gorm:"type:jsonb" json:"position" swaggertype:"object"`
}
