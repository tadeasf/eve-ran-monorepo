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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ESIItem struct {
	TypeID         int     `gorm:"primaryKey" json:"type_id"`
	GroupID        int     `gorm:"index" json:"group_id"`
	Name           string  `gorm:"type:text" json:"name"`
	Description    string  `gorm:"type:text" json:"description"`
	Mass           float64 `json:"mass"`
	Volume         float64 `json:"volume"`
	Capacity       float64 `json:"capacity"`
	PortionSize    int     `json:"portion_size"`
	PackagedVolume float64 `json:"packaged_volume"`
	Published      bool    `json:"published"`
	Radius         float64 `json:"radius"`
}

type ZKillboardItem struct {
	Flag              int              `json:"flag"`
	ItemTypeID        int              `json:"item_type_id"`
	QuantityDestroyed *int64           `json:"quantity_destroyed,omitempty"`
	QuantityDropped   *int64           `json:"quantity_dropped,omitempty"`
	Singleton         int              `json:"singleton"`
	Items             []ZKillboardItem `json:"items,omitempty"`
}

type Item struct {
	ItemTypeID        int   `json:"item_type_id"`
	QuantityDestroyed int64 `json:"quantity_destroyed,omitempty"`
	QuantityDropped   int64 `json:"quantity_dropped,omitempty"`
	Flag              int   `json:"flag"`
	Singleton         int   `json:"singleton"`
}

type ItemArray []Item

func (a ItemArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ItemArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var items []Item
	err := json.Unmarshal(bytes, &items)
	*a = ItemArray(items)
	return err
}
