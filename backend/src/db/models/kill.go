// Copyright (C) 2024 Tadeáš Fořt
// 
// This file is part of EVE Ran Services.
// 
// EVE Ran Services is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// EVE Ran Services is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with EVE Ran Services.  If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"encoding/json"
	"time"
)

type Kill struct {
	ID            uint  `gorm:"primaryKey"`
	KillmailID    int64 `gorm:"uniqueIndex"`
	KillmailTime  time.Time
	SolarSystemID int
	CharacterID   int64
	Victim        Victim `gorm:"embedded;embeddedPrefix:victim_"`
	Attackers     []byte `gorm:"type:jsonb"`
	ZkillData     Zkill  `gorm:"foreignKey:KillmailID;references:KillmailID"`
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

func (k *Kill) GetAttackers() ([]Attacker, error) {
	var attackers []Attacker
	err := json.Unmarshal(k.Attackers, &attackers)
	if err != nil {
		return nil, err
	}
	return attackers, nil
}
