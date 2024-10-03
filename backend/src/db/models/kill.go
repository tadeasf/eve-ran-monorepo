package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Character model
type Character struct {
	ID             int64   `gorm:"primaryKey" json:"id"`
	Name           string  `json:"name"`
	SecurityStatus float64 `json:"security_status"`
	Title          string  `json:"title"`
	RaceID         int     `json:"race_id"`
}

// Kill model
type Kill struct {
	ID             int64         `gorm:"primaryKey" json:"id"`
	KillmailID     int64         `gorm:"uniqueIndex" json:"killmail_id"`
	CharacterID    int64         `json:"character_id"`
	KillTime       time.Time     `json:"killmail_time"`
	SolarSystemID  int           `json:"solar_system_id"`
	LocationID     int64         `json:"locationID"`
	Hash           string        `json:"hash"`
	FittedValue    float64       `json:"fitted_value"`
	DroppedValue   float64       `json:"dropped_value"`
	DestroyedValue float64       `json:"destroyed_value"`
	TotalValue     float64       `json:"total_value"`
	Points         int           `json:"points"`
	NPC            bool          `json:"npc"`
	Solo           bool          `json:"solo"`
	Awox           bool          `json:"awox"`
	Victim         Victim        `gorm:"embedded;embeddedPrefix:victim_" json:"victim"`
	Attackers      AttackersJSON `gorm:"type:jsonb" json:"attackers"`
}

// AttackersJSON type
type AttackersJSON []Attacker

// Value implementation for AttackersJSON
func (a AttackersJSON) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implementation for AttackersJSON
func (a *AttackersJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(bytes, &a)
}

// Victim model
type Victim struct {
	AllianceID    *int      `json:"alliance_id,omitempty"`
	CharacterID   *int      `json:"character_id,omitempty"`
	CorporationID *int      `json:"corporation_id,omitempty"`
	FactionID     *int      `json:"faction_id,omitempty"`
	DamageTaken   int       `json:"damage_taken"`
	ShipTypeID    int       `json:"ship_type_id"`
	Items         ItemArray `json:"items" gorm:"type:jsonb"`
	Position      *Position `json:"position" gorm:"type:jsonb"`
}

// Attacker model
type Attacker struct {
	AllianceID     *int    `json:"alliance_id,omitempty"`
	CharacterID    *int    `json:"character_id,omitempty"`
	CorporationID  *int    `json:"corporation_id,omitempty"`
	FactionID      *int    `json:"faction_id,omitempty"`
	DamageDone     int     `json:"damage_done"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float64 `json:"security_status"`
	ShipTypeID     int     `json:"ship_type_id"`
	WeaponTypeID   int     `json:"weapon_type_id"`
}

// ZKillKill model
type ZKillKill struct {
	KillmailID int64 `json:"killmail_id"`
	ZKB        struct {
		LocationID     int64   `json:"locationID"`
		Hash           string  `json:"hash"`
		FittedValue    float64 `json:"fittedValue"`
		DroppedValue   float64 `json:"droppedValue"`
		DestroyedValue float64 `json:"destroyedValue"`
		TotalValue     float64 `json:"totalValue"`
		Points         int     `json:"points"`
		NPC            bool    `json:"npc"`
		Solo           bool    `json:"solo"`
		Awox           bool    `json:"awox"`
	} `json:"zkb"`
}
