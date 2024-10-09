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

// Constellation model
type Constellation struct {
	ConstellationID int             `gorm:"primaryKey" json:"constellation_id"`
	Name            string          `json:"name"`
	RegionID        int             `json:"region_id"`
	Systems         json.RawMessage `gorm:"type:jsonb" json:"systems"`
	Position        json.RawMessage `gorm:"type:jsonb" json:"position"`
}
