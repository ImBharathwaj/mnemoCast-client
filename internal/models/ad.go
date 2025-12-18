package models

import "time"

// Ad represents an advertisement delivered to the screen
type Ad struct {
	ID          string    `json:"id"`                    // Ad ID
	Title       string    `json:"title,omitempty"`       // Ad title
	Type        string    `json:"type"`                  // Ad type (image, video, html, etc.)
	ContentURL  string    `json:"contentUrl"`            // URL to ad content
	Duration    int       `json:"duration,omitempty"`    // Display duration in seconds
	StartTime   time.Time `json:"startTime,omitempty"`   // Scheduled start time
	EndTime     time.Time `json:"endTime,omitempty"`     // Scheduled end time
	Priority    int       `json:"priority,omitempty"`    // Display priority
	Metadata    map[string]interface{} `json:"metadata,omitempty"` // Additional metadata
}

// AdDeliveryResponse represents the response from the ad delivery endpoint
type AdDeliveryResponse struct {
	Ads       []Ad       `json:"ads"`                    // List of ads to display
	PlaylistID string    `json:"playlistId,omitempty"`   // Associated playlist ID
	UpdatedAt time.Time  `json:"updatedAt"`              // Last update timestamp
}

