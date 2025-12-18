package models

import "time"

// ScreenConfig represents the complete configuration for a screen
type ScreenConfig struct {
	Identity         ScreenIdentity `json:"identity"`
	AdServerURL      string         `json:"adServerUrl"`      // Backend URL
	HeartbeatInterval int          `json:"heartbeatInterval"` // Seconds between heartbeats
	AdFetchInterval   int          `json:"adFetchInterval"`   // Seconds between ad fetches
	RetryAttempts    int           `json:"retryAttempts"`     // Max retry attempts
	RetryDelay       int           `json:"retryDelay"`        // Seconds between retries
}

// DefaultConfig returns a default configuration
func DefaultConfig() *ScreenConfig {
	now := time.Now()
	return &ScreenConfig{
		Identity: ScreenIdentity{
			Name:          "Default Screen",
			Country:       "India",
			City:          "Chennai",
			Area:          "Airport",
			VenueType:     "airport",
			Timezone:      "Asia/Kolkata",
			Width:         1920,
			Height:        1080,
			IsAudible:     false,
			IsOnline:      false,
			Classification: 1,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		AdServerURL:      "http://10.42.0.1:8080",
		HeartbeatInterval: 30,
		AdFetchInterval:   60, // Fetch ads every 60 seconds (1 minute)
		RetryAttempts:    3,
		RetryDelay:       5,
	}
}

