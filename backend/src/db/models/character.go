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

// CharacterStats model
type CharacterStats struct {
	CharacterID int64   `json:"character_id"`
	KillCount   int     `json:"kill_count"`
	TotalISK    float64 `json:"total_isk"`
}

// Character model
type Character struct {
	ID             int64   `gorm:"primaryKey" json:"id"`
	Name           string  `json:"name"`
	SecurityStatus float64 `json:"security_status"`
	Title          string  `json:"title"`
	RaceID         int     `json:"race_id"`
}
