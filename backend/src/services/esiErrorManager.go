package services

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ESIRateLimitInfo tracks rate limiting for a specific route group
type ESIRateLimitInfo struct {
	group      string
	limit      int    // Total tokens per window
	remaining  int    // Available tokens
	used       int    // Tokens used by last request
	windowSize string // e.g., "15m", "1h"
	lastUpdate time.Time
	retryAfter time.Time // When we can retry after 429
}

// ESIErrorManager manages both new rate limiting and legacy error limiting
type ESIErrorManager struct {
	mu sync.RWMutex

	// New rate limiting (per route group)
	rateLimits map[string]*ESIRateLimitInfo

	// Legacy error limiting (100 non-2xx/3xx per minute)
	errorRemaining int
	errorResetTime time.Time

	// Backoff for 429 errors
	globalRetryAfter time.Time
	lastUpdate       time.Time
}

var globalESIManager = &ESIErrorManager{
	rateLimits:     make(map[string]*ESIRateLimitInfo),
	errorRemaining: 100, // Legacy error limit
	errorResetTime: time.Now().Add(time.Minute),
}

// GetGlobalESIManager returns the singleton ESI error manager
func GetGlobalESIManager() *ESIErrorManager {
	return globalESIManager
}

// UpdateLimitsFromHeaders updates rate limit state from ESI response headers
func (em *ESIErrorManager) UpdateLimitsFromHeaders(headers http.Header) {
	em.mu.Lock()
	defer em.mu.Unlock()

	// Update new rate limiting headers (X-Ratelimit-*)
	if group := headers.Get("X-Ratelimit-Group"); group != "" {
		if em.rateLimits[group] == nil {
			em.rateLimits[group] = &ESIRateLimitInfo{group: group}
		}

		info := em.rateLimits[group]

		// Parse X-Ratelimit-Limit (format: "150/15m")
		if limitStr := headers.Get("X-Ratelimit-Limit"); limitStr != "" {
			var limit int
			var window string
			if n, _ := fmt.Sscanf(limitStr, "%d/%s", &limit, &window); n == 2 {
				info.limit = limit
				info.windowSize = window
			}
		}

		// Parse X-Ratelimit-Remaining
		if remainStr := headers.Get("X-Ratelimit-Remaining"); remainStr != "" {
			if remaining, err := strconv.Atoi(remainStr); err == nil {
				info.remaining = remaining
			}
		}

		// Parse X-Ratelimit-Used
		if usedStr := headers.Get("X-Ratelimit-Used"); usedStr != "" {
			if used, err := strconv.Atoi(usedStr); err == nil {
				info.used = used
			}
		}

		info.lastUpdate = time.Now()
	}

	// Update legacy error limiting headers (X-ESI-Error-Limit-*)
	if remainStr := headers.Get("X-ESI-Error-Limit-Remain"); remainStr != "" {
		if remaining, err := strconv.Atoi(remainStr); err == nil {
			em.errorRemaining = remaining
		}
	}

	if resetStr := headers.Get("X-ESI-Error-Limit-Reset"); resetStr != "" {
		if resetSeconds, err := strconv.Atoi(resetStr); err == nil {
			em.errorResetTime = time.Now().Add(time.Duration(resetSeconds) * time.Second)
		}
	}

	// Handle Retry-After header (for 429 responses)
	if retryStr := headers.Get("Retry-After"); retryStr != "" {
		if retrySeconds, err := strconv.Atoi(retryStr); err == nil {
			em.globalRetryAfter = time.Now().Add(time.Duration(retrySeconds) * time.Second)
			log.Printf("ESI Rate Limited: Retry after %d seconds (at %s)", retrySeconds, em.globalRetryAfter.Format(time.RFC3339))
		}
	}

	em.lastUpdate = time.Now()
}

// CanMakeRequest checks if we can make a request without hitting limits
func (em *ESIErrorManager) CanMakeRequest() bool {
	em.mu.RLock()
	defer em.mu.RUnlock()

	now := time.Now()

	// Check global retry-after from 429 responses
	if now.Before(em.globalRetryAfter) {
		return false
	}

	// Check legacy error limit
	if now.Before(em.errorResetTime) && em.errorRemaining <= 5 {
		// Keep buffer of 5 errors
		return false
	}

	// Check if any route group is approaching limits
	for _, info := range em.rateLimits {
		if info.remaining < 10 && now.Before(info.lastUpdate.Add(time.Minute)) {
			// Approaching limit, slow down
			return false
		}
	}

	return true
}

// WaitForReset waits until we can make requests again
func (em *ESIErrorManager) WaitForReset() {
	em.mu.RLock()
	globalRetry := em.globalRetryAfter
	errorReset := em.errorResetTime
	em.mu.RUnlock()

	now := time.Now()
	var waitUntil time.Time

	// Find the nearest time we need to wait for
	if now.Before(globalRetry) {
		waitUntil = globalRetry
	} else if now.Before(errorReset) {
		waitUntil = errorReset
	} else {
		// No need to wait
		return
	}

	sleepDuration := time.Until(waitUntil)
	if sleepDuration > 0 {
		log.Printf("ESI rate limit: sleeping for %v until %s", sleepDuration, waitUntil.Format(time.RFC3339))
		time.Sleep(sleepDuration)
		// Add small buffer
		time.Sleep(500 * time.Millisecond)
	}
}

// ShouldBackoff determines if we should slow down requests
func (em *ESIErrorManager) ShouldBackoff() (bool, time.Duration) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	// Check if error remaining is low
	if em.errorRemaining < 20 {
		// Slow down significantly
		return true, 500 * time.Millisecond
	}

	// Check rate limit remaining for any group
	for _, info := range em.rateLimits {
		if info.limit > 0 {
			// Calculate percentage remaining
			pctRemaining := float64(info.remaining) / float64(info.limit) * 100

			if pctRemaining < 10 {
				// Less than 10% remaining, slow down significantly
				return true, 1 * time.Second
			} else if pctRemaining < 30 {
				// Less than 30% remaining, moderate slowdown
				return true, 300 * time.Millisecond
			}
		}
	}

	return false, 0
}

// DecrementErrorCount decrements the legacy error count
func (em *ESIErrorManager) DecrementErrorCount() {
	em.mu.Lock()
	defer em.mu.Unlock()
	if em.errorRemaining > 0 {
		em.errorRemaining--
	}
}

// GetStatus returns the current rate limit status
func (em *ESIErrorManager) GetStatus() (remaining int, resetTime time.Time) {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.errorRemaining, em.errorResetTime
}

// GetRateLimitInfo returns info for a specific route group
func (em *ESIErrorManager) GetRateLimitInfo(group string) *ESIRateLimitInfo {
	em.mu.RLock()
	defer em.mu.RUnlock()
	if info, ok := em.rateLimits[group]; ok {
		// Return a copy to avoid race conditions
		copy := *info
		return &copy
	}
	return nil
}

// LogStatus logs current rate limit status
func (em *ESIErrorManager) LogStatus() {
	em.mu.RLock()
	defer em.mu.RUnlock()

	log.Printf("ESI Rate Limit Status - Error Remaining: %d, Error Reset: %s",
		em.errorRemaining, em.errorResetTime.Format(time.RFC3339))

	for group, info := range em.rateLimits {
		log.Printf("  Group %s: %d/%d tokens remaining (%s window)",
			group, info.remaining, info.limit, info.windowSize)
	}
}
