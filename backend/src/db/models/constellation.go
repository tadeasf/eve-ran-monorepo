package models

// Constellation model
type Constellation struct {
	ConstellationID int       `gorm:"primaryKey" json:"constellation_id"`
	Name            string    `json:"name"`
	RegionID        int       `json:"region_id"`
	Systems         []int     `gorm:"type:jsonb" json:"systems"`
	Position        *Position `gorm:"type:jsonb" json:"position"`
}
