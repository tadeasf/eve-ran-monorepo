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

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreRegions(c *gin.Context) {
	regions, err := services.FetchAllRegions(10) // Use 10 concurrent requests
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, region := range regions {
		err = queries.UpsertRegion(region)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Regions fetched and stored successfully", "count": len(regions)})
}

// GetAllRegions retrieves all regions from the database
// @Summary Get all regions
// @Description Fetch all regions from the database
// @Tags regions
// @Accept json
// @Produce json
// @Success 200 {array} models.Region
// @Failure 500 {object} models.ErrorResponse
// @Router /regions [get]
func GetAllRegions(c *gin.Context) {
	regions, err := queries.GetAllRegions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, regions)
}
