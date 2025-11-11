package services

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ESIErrorManager struct {
	mu             sync.Mutex
	errorRemaining int
	resetTime      time.Time
	lastUpdate     time.Time
}

var globalESIManager = &ESIErrorManager{
	errorRemaining: 100, // Default conservative value
	resetTime:      time.Now(),
}

// GetGlobalESIManager returns the singleton ESI error manager
func GetGlobalESIManager() *ESIErrorManager {
	return globalESIManager
}

// UpdateLimitsFromHeaders updates the rate limit state from ESI response headers
func (em *ESIErrorManager) UpdateLimitsFromHeaders(headers http.Header) {
	em.mu.Lock()
	defer em.mu.Unlock()

	// Parse X-ESI-Error-Limit-Remain header
	if remainStr := headers.Get("X-ESI-Error-Limit-Remain"); remainStr != "" {
		if remaining, err := strconv.Atoi(remainStr); err == nil {
			em.errorRemaining = remaining
		}
	}

	// Parse X-ESI-Error-Limit-Reset header
	if resetStr := headers.Get("X-ESI-Error-Limit-Reset"); resetStr != "" {
		if resetSeconds, err := strconv.Atoi(resetStr); err == nil {
			em.resetTime = time.Now().Add(time.Duration(resetSeconds) * time.Second)
		}
	}

	em.lastUpdate = time.Now()
}

// UpdateLimits updates rate limits manually (for backwards compatibility)
func (em *ESIErrorManager) UpdateLimits(remaining int, resetSeconds int) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.errorRemaining = remaining
	em.resetTime = time.Now().Add(time.Duration(resetSeconds) * time.Second)
	em.lastUpdate = time.Now()
}

// CanMakeRequest checks if we can make a request without hitting the error limit
func (em *ESIErrorManager) CanMakeRequest() bool {
	em.mu.Lock()
	defer em.mu.Unlock()

	// If the reset time has passed, we can make requests
	if time.Now().After(em.resetTime) {
		return true
	}

	// Check if we have errors remaining
	return em.errorRemaining > 0
}

// DecrementErrorCount decrements the error count when an error occurs
func (em *ESIErrorManager) DecrementErrorCount() {
	em.mu.Lock()
	defer em.mu.Unlock()
	if em.errorRemaining > 0 {
		em.errorRemaining--
	}
}

// WaitForReset waits until the error limit resets
func (em *ESIErrorManager) WaitForReset() {
	em.mu.Lock()
	resetTime := em.resetTime
	em.mu.Unlock()

	sleepDuration := time.Until(resetTime)
	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
}

// GetStatus returns the current rate limit status
func (em *ESIErrorManager) GetStatus() (remaining int, resetTime time.Time) {
	em.mu.Lock()
	defer em.mu.Unlock()
	return em.errorRemaining, em.resetTime
}
