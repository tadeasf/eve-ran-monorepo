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
	Planets         json.RawMessage `gorm:"type:jsonb" json:"planets"`
	Stargates       json.RawMessage `gorm:"type:jsonb" json:"stargates"`
	Stations        json.RawMessage `gorm:"type:jsonb" json:"stations"`
	Position        json.RawMessage `gorm:"type:jsonb" json:"position"`
}
