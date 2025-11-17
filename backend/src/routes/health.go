package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/cache"
	"github.com/tadeasf/eve-ran/src/db"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Services  map[string]ServiceInfo `json:"services,omitempty"`
}

// ServiceInfo represents the status of a service
type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Health godoc
// @Summary Basic health check
// @Description Returns the health status of the API
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadyCheck godoc
// @Summary Readiness check
// @Description Returns the readiness status of the API and its dependencies
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /ready [get]
func ReadyCheck(c *gin.Context) {
	services := make(map[string]ServiceInfo)
	allHealthy := true

	// Check database connection
	sqlDB, err := db.DB.DB()
	if err != nil {
		services["database"] = ServiceInfo{
			Status:  "unhealthy",
			Message: "failed to get database instance: " + err.Error(),
		}
		allHealthy = false
	} else {
		err = sqlDB.Ping()
		if err != nil {
			services["database"] = ServiceInfo{
				Status:  "unhealthy",
				Message: "database ping failed: " + err.Error(),
			}
			allHealthy = false
		} else {
			services["database"] = ServiceInfo{
				Status: "healthy",
			}
		}
	}

	// Check Redis connection
	redisClient := cache.GetRedisClient()
	if redisClient == nil {
		services["redis"] = ServiceInfo{
			Status:  "unhealthy",
			Message: "redis client is nil",
		}
		allHealthy = false
	} else {
		err = redisClient.Ping(c).Err()
		if err != nil {
			services["redis"] = ServiceInfo{
				Status:  "unhealthy",
				Message: "redis ping failed: " + err.Error(),
			}
			allHealthy = false
		} else {
			services["redis"] = ServiceInfo{
				Status: "healthy",
			}
		}
	}

	status := "ready"
	statusCode := http.StatusOK
	if !allHealthy {
		status = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services:  services,
	})
}

// LiveCheck godoc
// @Summary Liveness check
// @Description Returns whether the API is alive (for Kubernetes liveness probe)
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /live [get]
func LiveCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "alive",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
