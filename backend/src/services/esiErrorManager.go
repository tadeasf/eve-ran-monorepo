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
