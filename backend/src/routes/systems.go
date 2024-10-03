package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreSystems(c *gin.Context) {
	batchSize := 100
	totalSystems := 0

	for {
		// Fetch a batch of systems
		systems, err := services.FetchAllSystems(batchSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(systems) == 0 {
			break // No more systems to fetch
		}

		// Batch upsert systems
		err = db.BatchUpsertSystems(systems)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalSystems += len(systems)

		if len(systems) < batchSize {
			break // Last batch
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Systems fetched and stored successfully", "count": totalSystems})
}

func GetAllSystems(c *gin.Context) {
	systems, err := db.GetAllSystems()
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

	system, err := db.GetSystem(systemID)
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

	systems, err := db.GetSystemsByRegionID(regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, systems)
}
