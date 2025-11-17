package main

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tadeasf/eve-ran/docs"
	"github.com/tadeasf/eve-ran/src/cache"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/jobs"
	"github.com/tadeasf/eve-ran/src/middleware"
	"github.com/tadeasf/eve-ran/src/routes"
	"github.com/tadeasf/eve-ran/src/utils"
)

// @title EVE Ran API
// @version 1.0
// @description This is the API for EVE Ran application.
// @host api.tundragon.space
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key for authentication. Add your API key in the X-API-Key header.

func main() {
	utils.InitLogger()
	gin.SetMode(gin.ReleaseMode)

	db.InitDB()

	// Initialize Redis cache
	if err := cache.InitRedis(); err != nil {
		utils.ErrorLogger.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer cache.CloseRedis()

	// Add a delay to allow initial data to be stored
	time.Sleep(1 * time.Minute)

	// Run the type fetcher job
	jobs.FetchAndUpdateTypes()

	// Start the kill cron job
	go jobs.StartKillCron()

	// Start the kill enhancement job
	go func() {
		for {
			jobs.EnhanceKills()
			time.Sleep(1 * time.Minute) // Run every 15 minutes
		}
	}()

	// Start the competition cron job
	go jobs.StartCompetitionCron()

	r := gin.Default()

	// CORS middleware - configure allowed origins
	r.Use(func(c *gin.Context) {
		allowedOrigins := []string{
			"https://tundragon.space",
			"https://api.tundragon.space",
			"http://localhost:12921",
			"http://localhost:3000",
		}

		origin := c.Request.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Create cache middleware with 5-minute TTL
	cacheMiddleware := middleware.CacheMiddleware(5 * time.Minute)

	// Health check endpoints (no authentication required)
	r.GET("/health", routes.Health)
	r.GET("/ready", routes.ReadyCheck)
	r.GET("/live", routes.LiveCheck)

	// API key middleware for protected routes
	apiKeyMiddleware := middleware.APIKeyAuth()

	// Apply API key middleware to all routes except health checks and swagger
	api := r.Group("/")
	api.Use(apiKeyMiddleware)

	// Character routes
	api.POST("/characters", routes.AddCharacter)
	api.DELETE("/characters/:id", routes.RemoveCharacter)
	api.DELETE("/characters/batch", routes.BatchDeleteCharacters)
	api.GET("/characters/:id/kills/db", routes.GetCharacterKillsFromDB)

	// Region routes
	api.POST("/regions/fetch", routes.FetchAndStoreRegions)
	api.GET("/regions", routes.GetAllRegions)

	// System routes (cached)
	api.POST("/systems/fetch", routes.FetchAndStoreSystems)
	api.GET("/systems", cacheMiddleware, routes.GetAllSystems)
	api.GET("/systems/:id", routes.GetSystemByID)
	api.GET("/systems/region/:regionID", routes.GetSystemsByRegion)

	// Constellation routes
	api.POST("/constellations/fetch", routes.FetchAndStoreConstellations)
	api.GET("/constellations", routes.GetAllConstellations)
	api.GET("/constellations/:id", routes.GetConstellationByID)
	api.GET("/constellations/region/:regionID", routes.GetConstellationsByRegion)

	// Item routes
	api.POST("/items/fetch", routes.FetchAndStoreItems)
	api.GET("/items", routes.GetAllItems)
	api.GET("/items/:typeID", routes.GetItemByTypeID)

	// Character killmail and stats routes
	api.GET("/characters/:id/killmails", routes.GetCharacterKillmails)
	api.GET("/characters/stats", routes.GetAllCharacterStats)

	// Data routes (cached)
	api.GET("/characters", routes.GetAllCharacters)
	api.GET("/kills", cacheMiddleware, routes.GetAllKills)

	// Character search and batch management
	api.GET("/characters/search", routes.SearchCharacters)
	api.POST("/characters/batch", routes.BatchAddCharacters)

	// Admin dashboard routes (cached)
	api.GET("/admin/stats", cacheMiddleware, routes.GetDashboardStats)

	// Kill by region route (cached)
	api.GET("/kills/region/:regionID", cacheMiddleware, routes.GetKillsByRegion)

	// Kill comments routes
	api.GET("/kills/:killmail_id/comments", routes.GetKillComments)
	api.POST("/kills/:killmail_id/comments", routes.CreateKillComment)
	api.PUT("/comments/:id", routes.UpdateKillComment)
	api.DELETE("/comments/:id", routes.DeleteKillComment)

	// Competition routes (no caching - always fetch fresh data based on current settings)
	api.GET("/competition/settings", routes.GetCompetitionSettings)
	api.POST("/competition/settings", routes.UpdateCompetitionSettings)
	api.GET("/competition/current", routes.GetCurrentCompetitionStandings)
	api.GET("/competition/history", routes.GetAllCompetitionHistory)
	api.GET("/competition/history/:month/:year", routes.GetCompetitionHistory)
	api.GET("/competition/recent-winners", routes.GetRecentCompetitionWinners)
	api.GET("/competition/ytd-winners", routes.GetYearToDateWinners)
	api.POST("/competition/save/:month/:year", routes.SaveMonthlyResults)

	// Setup Swagger (no authentication required for documentation)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
