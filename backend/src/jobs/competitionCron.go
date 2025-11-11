package jobs

import (
	"fmt"
	"time"

	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/utils"
)

// StartCompetitionCron starts a cron job that saves monthly competition results
// It runs daily and checks if it's the first day of a new month
func StartCompetitionCron() {
	ticker := time.NewTicker(24 * time.Hour) // Check once per day
	defer ticker.Stop()

	utils.InfoLogger.Println("Competition cron job started")

	for range ticker.C {
		checkAndSaveMonthlyResults()
	}
}

func checkAndSaveMonthlyResults() {
	now := time.Now()

	// Check if it's the first day of the month
	if now.Day() != 1 {
		return
	}

	// Get last month's data
	lastMonth := now.AddDate(0, -1, 0)
	month := int(lastMonth.Month())
	year := lastMonth.Year()

	utils.InfoLogger.Printf("Saving competition results for %d/%d", month, year)

	// Save the results
	err := queries.SaveMonthlyCompetitionResults(month, year)
	if err != nil {
		utils.ErrorLogger.Printf("Error saving monthly competition results for %d/%d: %v", month, year, err)
		return
	}

	utils.InfoLogger.Printf("Successfully saved competition results for %d/%d", month, year)
}

// SavePreviousMonthResults manually saves the previous month's results (for initial setup)
func SavePreviousMonthResults() error {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	month := int(lastMonth.Month())
	year := lastMonth.Year()

	utils.InfoLogger.Printf("Manually saving competition results for %d/%d", month, year)

	err := queries.SaveMonthlyCompetitionResults(month, year)
	if err != nil {
		return fmt.Errorf("failed to save monthly competition results: %v", err)
	}

	utils.InfoLogger.Printf("Successfully saved competition results for %d/%d", month, year)
	return nil
}
