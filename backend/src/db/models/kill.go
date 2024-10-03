package models

import (
	"time"
)

type Kill struct {
	ID            uint  `gorm:"primaryKey"`
	KillmailID    int64 `gorm:"uniqueIndex"`
	KillmailTime  time.Time
	SolarSystemID int
	Victim        Victim     `gorm:"embedded;embeddedPrefix:victim_"`
	Attackers     []Attacker `gorm:"type:jsonb"`
	ZkillData     Zkill      `gorm:"foreignKey:KillmailID;references:KillmailID"`
}

type Victim struct {
	AllianceID    int64
	CharacterID   int64
	CorporationID int64
	DamageTaken   int
	ShipTypeID    int
	Position      Position  `gorm:"embedded;embeddedPrefix:position_"`
	Items         ItemArray `gorm:"type:jsonb"`
}

type Attacker struct {
	AllianceID     int64   `json:"alliance_id,omitempty"`
	CharacterID    int64   `json:"character_id,omitempty"`
	CorporationID  int64   `json:"corporation_id,omitempty"`
	DamageDone     int     `json:"damage_done"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float64 `json:"security_status"`
	ShipTypeID     int     `json:"ship_type_id"`
	WeaponTypeID   int     `json:"weapon_type_id"`
}

type Zkill struct {
	ID             uint  `gorm:"primaryKey"`
	KillmailID     int64 `gorm:"uniqueIndex"`
	CharacterID    int64
	LocationID     int64
	Hash           string
	FittedValue    float64
	DroppedValue   float64
	DestroyedValue float64
	TotalValue     float64
	Points         int
	NPC            bool
	Solo           bool
	Awox           bool
	Labels         []string `gorm:"type:text[]"`
}
