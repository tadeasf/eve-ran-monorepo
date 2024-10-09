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

package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db/queries"
)

func GetCharacterKillmails(c *gin.Context) {
	characterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	startTime, _ := time.Parse(time.RFC3339, c.Query("start_time"))
	endTime, _ := time.Parse(time.RFC3339, c.Query("end_time"))
	systemID, _ := strconv.ParseInt(c.Query("system_id"), 10, 64)
	regionID, _ := strconv.ParseInt(c.Query("region_id"), 10, 64)

	kills, err := queries.GetCharacterKillmails(characterID, startTime, endTime, systemID, regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kills)
}
