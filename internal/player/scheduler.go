package player

import (
	"mnemoCast-client/internal/models"
	"time"
)

// Scheduler handles ad duration and transition timing
type Scheduler struct {
	defaultDuration time.Duration
	transitionDelay time.Duration
	minDuration     time.Duration
	maxDuration     time.Duration
}

// NewScheduler creates a new scheduler with default values
func NewScheduler(defaultDurationSeconds, transitionDelaySeconds int) *Scheduler {
	defaultDur := time.Duration(defaultDurationSeconds) * time.Second
	transitionDur := time.Duration(transitionDelaySeconds) * time.Second
	
	return &Scheduler{
		defaultDuration: defaultDur,
		transitionDelay: transitionDur,
		minDuration:     5 * time.Second,  // Minimum 5 seconds
		maxDuration:     300 * time.Second, // Maximum 5 minutes
	}
}

// GetAdDuration returns the duration for an ad
// Uses ad.Duration if specified, otherwise falls back to defaultDuration
// Enforces min/max duration limits
func (s *Scheduler) GetAdDuration(ad *models.Ad) time.Duration {
	var duration time.Duration
	
	if ad.Duration > 0 {
		duration = time.Duration(ad.Duration) * time.Second
	} else {
		duration = s.defaultDuration
	}
	
	// Enforce minimum duration
	if duration < s.minDuration {
		duration = s.minDuration
	}
	
	// Enforce maximum duration
	if duration > s.maxDuration {
		duration = s.maxDuration
	}
	
	return duration
}

// GetTransitionDelay returns the delay between ads
func (s *Scheduler) GetTransitionDelay() time.Duration {
	return s.transitionDelay
}

// ShouldTransition checks if an ad should transition to the next one
// Returns true if the ad has been playing longer than its duration
func (s *Scheduler) ShouldTransition(ad *models.Ad, startTime time.Time) bool {
	duration := s.GetAdDuration(ad)
	elapsed := time.Since(startTime)
	return elapsed >= duration
}

