package models

// Credentials represents authentication credentials for the screen
// Screen ID and Passkey are server-assigned and must be configured manually
type Credentials struct {
	ScreenID string `json:"screenId"`  // Server-assigned screen ID (NOT NULL)
	Passkey  string `json:"passkey"`    // Server-assigned passkey (NOT NULL)
}

// HasCredentials checks if both screen ID and passkey are present
func (c *Credentials) HasCredentials() bool {
	return c.ScreenID != "" && c.Passkey != ""
}

// IsValid checks if credentials are valid (both screen ID and passkey must be present)
func (c *Credentials) IsValid() bool {
	return c.HasCredentials()
}

