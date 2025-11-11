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

	// Create cache middleware with 5-minute TTL
	cacheMiddleware := middleware.CacheMiddleware(5 * time.Minute)

	// zKillboard routes
	r.POST("/characters", routes.AddCharacter)
	r.DELETE("/characters/:id", routes.RemoveCharacter)
	r.GET("/characters/:id/kills/db", routes.GetCharacterKillsFromDB)

	// Region routes
	r.POST("/regions/fetch", routes.FetchAndStoreRegions)
	r.GET("/regions", routes.GetAllRegions)

	// System routes (cached)
	r.POST("/systems/fetch", routes.FetchAndStoreSystems)
	r.GET("/systems", cacheMiddleware, routes.GetAllSystems)
	r.GET("/systems/:id", routes.GetSystemByID)
	r.GET("/systems/region/:regionID", routes.GetSystemsByRegion)

	// Constellation routes
	r.POST("/constellations/fetch", routes.FetchAndStoreConstellations)
	r.GET("/constellations", routes.GetAllConstellations)
	r.GET("/constellations/:id", routes.GetConstellationByID)
	r.GET("/constellations/region/:regionID", routes.GetConstellationsByRegion)

	// Item routes
	r.POST("/items/fetch", routes.FetchAndStoreItems)
	r.GET("/items", routes.GetAllItems)
	r.GET("/items/:typeID", routes.GetItemByTypeID)

	// New routes
	r.GET("/characters/:id/killmails", routes.GetCharacterKillmails)
	r.GET("/characters/stats", routes.GetAllCharacterStats)

	// New data routes (cached)
	r.GET("/characters", routes.GetAllCharacters)
	r.GET("/kills", cacheMiddleware, routes.GetAllKills)

	// Character search and batch management
	r.GET("/characters/search", routes.SearchCharacters)
	r.POST("/characters/batch", routes.BatchAddCharacters)

	// Admin dashboard routes (cached)
	r.GET("/admin/stats", cacheMiddleware, routes.GetDashboardStats)

	// Kill by region route (cached)
	r.GET("/kills/region/:regionID", cacheMiddleware, routes.GetKillsByRegion)

	// Competition routes (no caching - always fetch fresh data based on current settings)
	r.GET("/competition/settings", routes.GetCompetitionSettings)
	r.POST("/competition/settings", routes.UpdateCompetitionSettings)
	r.GET("/competition/current", routes.GetCurrentCompetitionStandings)
	r.GET("/competition/history", routes.GetAllCompetitionHistory)
	r.GET("/competition/history/:month/:year", routes.GetCompetitionHistory)
	r.GET("/competition/recent-winners", routes.GetRecentCompetitionWinners)
	r.GET("/competition/ytd-winners", routes.GetYearToDateWinners)
	r.POST("/competition/save/:month/:year", routes.SaveMonthlyResults)

	// Setup Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
