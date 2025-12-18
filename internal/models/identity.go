package models

import "time"

// ScreenIdentity represents the complete identity of a screen
type ScreenIdentity struct {
	ID            string    `json:"id"`              // Unique screen ID (UUID) - PRIMARY KEY
	Name          string    `json:"name"`            // Human-readable name - NOT NULL
	Country       string    `json:"country,omitempty"` // Country
	City          string    `json:"city,omitempty"`    // City
	Area          string    `json:"area,omitempty"`    // Area
	VenueType     string    `json:"venueType,omitempty"` // Venue type
	Timezone      string    `json:"timezone,omitempty"`   // Timezone
	Width         int       `json:"width,omitempty"`      // Screen width in pixels
	Height        int       `json:"height,omitempty"`     // Screen height in pixels
	IsAudible     bool      `json:"isAudible"`           // Audio capability - DEFAULT false
	IsOnline      bool      `json:"isOnline"`            // Online status - DEFAULT false
	LastSeen      *time.Time `json:"lastSeen,omitempty"`  // Last heartbeat time (TIMESTAMPTZ)
	Classification int      `json:"classification"`       // Screen classification - DEFAULT 1
	CreatedAt     time.Time `json:"createdAt"`           // First registration time - DEFAULT now()
	UpdatedAt     time.Time `json:"updatedAt"`           // Last update time - DEFAULT now()
}

// Validate checks if the identity is valid
func (si *ScreenIdentity) Validate() error {
	if si.ID == "" {
		return ErrInvalidIdentity
	}
	if si.Name == "" {
		return ErrInvalidIdentity
	}
	// Note: Other fields are optional, only ID and Name are required
	return nil
}

// GetLocation returns location information as a map (for backward compatibility)
func (si *ScreenIdentity) GetLocation() map[string]string {
	return map[string]string{
		"country":   si.Country,
		"city":      si.City,
		"area":      si.Area,
		"venueType": si.VenueType,
	}
}

