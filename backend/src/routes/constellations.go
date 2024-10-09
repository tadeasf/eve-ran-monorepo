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

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreConstellations(c *gin.Context) {
	batchSize := 250
	totalConstellations := 0

	for {
		// Fetch a batch of constellations
		constellations, err := services.FetchAllConstellations(batchSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(constellations) == 0 {
			break // No more constellations to fetch
		}

		// Batch upsert constellations
		err = queries.BatchUpsertConstellations(constellations)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalConstellations += len(constellations)

		if len(constellations) < batchSize {
			break // Last batch
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Constellations fetched and stored successfully", "count": totalConstellations})
}

func GetAllConstellations(c *gin.Context) {
	constellations, err := queries.GetAllConstellations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, constellations)
}

func GetConstellationByID(c *gin.Context) {
	constellationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constellation ID"})
		return
	}

	constellation, err := queries.GetConstellationByID(constellationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if constellation == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Constellation not found"})
		return
	}

	c.JSON(http.StatusOK, constellation)
}

func GetConstellationsByRegion(c *gin.Context) {
	regionID, err := strconv.Atoi(c.Param("regionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	constellations, err := queries.GetConstellationsByRegionID(regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, constellations)
}
