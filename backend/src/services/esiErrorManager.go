package services

import (
	"sync"
	"time"
)

type ESIErrorManager struct {
	mu             sync.Mutex
	errorRemaining int
	resetTime      time.Time
}

var (
	esiErrorManager = &ESIErrorManager{
		errorRemaining: 60,                              // Default to 60, will be updated with actual value from ESI
		resetTime:      time.Now().Add(1 * time.Minute), // Default to 1 minute, will be updated with actual value from ESI
	}
)

func (em *ESIErrorManager) UpdateLimits(remaining int, resetSeconds int) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.errorRemaining = remaining
	em.resetTime = time.Now().Add(time.Duration(resetSeconds) * time.Second)
}

func (em *ESIErrorManager) CanMakeRequest() bool {
	em.mu.Lock()
	defer em.mu.Unlock()
	if time.Now().After(em.resetTime) {
		// If we've passed the reset time, assume we can make a request
		return true
	}
	return em.errorRemaining > 0
}

func (em *ESIErrorManager) DecrementErrorCount() {
	em.mu.Lock()
	defer em.mu.Unlock()
	if em.errorRemaining > 0 {
		em.errorRemaining--
	}
}

func (em *ESIErrorManager) WaitForReset() {
	em.mu.Lock()
	resetTime := em.resetTime
	em.mu.Unlock()

	sleepDuration := time.Until(resetTime)
	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
}
