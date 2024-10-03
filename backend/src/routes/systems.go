package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreSystems(c *gin.Context) {
	batchSize := 100
	totalSystems := 0

	for {
		systems, err := services.FetchAllSystems(batchSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(systems) == 0 {
			break
		}

		err = queries.BatchUpsertSystems(systems)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalSystems += len(systems)

		if len(systems) < batchSize {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Systems fetched and stored successfully", "count": totalSystems})
}

func GetAllSystems(c *gin.Context) {
	systems, err := queries.GetAllSystems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, systems)
}

func GetSystemByID(c *gin.Context) {
	systemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system ID"})
		return
	}

	system, err := queries.GetSystemByID(systemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if system == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		return
	}

	c.JSON(http.StatusOK, system)
}

func GetSystemsByRegion(c *gin.Context) {
	regionID, err := strconv.Atoi(c.Param("regionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	systems, err := queries.GetSystemsByRegionID(regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, systems)
}
