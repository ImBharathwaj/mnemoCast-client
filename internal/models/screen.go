package models

import "time"

// Screen represents a registered screen (matches database schema)
type Screen struct {
	ID            string     `json:"id"`              // PRIMARY KEY
	Name          string     `json:"name"`            // NOT NULL
	Country       string     `json:"country,omitempty"`
	City          string     `json:"city,omitempty"`
	Area          string     `json:"area,omitempty"`
	VenueType     string     `json:"venueType,omitempty"`
	Timezone      string     `json:"timezone,omitempty"`
	Width         int        `json:"width,omitempty"`
	Height        int        `json:"height,omitempty"`
	IsAudible     bool       `json:"isAudible"`       // DEFAULT false
	IsOnline      bool       `json:"isOnline"`       // DEFAULT false
	LastSeen      *time.Time `json:"lastSeen,omitempty"` // TIMESTAMPTZ
	Classification int      `json:"classification"`   // DEFAULT 1
	CreatedAt     time.Time `json:"createdAt"`       // DEFAULT now()
	UpdatedAt     time.Time `json:"updatedAt"`       // DEFAULT now()
}

// RegisterScreenRequest represents a screen registration request
type RegisterScreenRequest struct {
	ID            string    `json:"id"`              // PRIMARY KEY
	Name          string    `json:"name"`            // NOT NULL
	Country       string    `json:"country,omitempty"`
	City          string    `json:"city,omitempty"`
	Area          string    `json:"area,omitempty"`
	VenueType     string    `json:"venueType,omitempty"`
	Timezone      string    `json:"timezone,omitempty"`
	Width         int       `json:"width,omitempty"`
	Height        int       `json:"height,omitempty"`
	IsAudible     bool      `json:"isAudible"`      // DEFAULT false
	Classification int      `json:"classification"` // DEFAULT 1
}

// HeartbeatRequest represents a heartbeat request
type HeartbeatRequest struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// HeartbeatResponse represents a heartbeat response
type HeartbeatResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

