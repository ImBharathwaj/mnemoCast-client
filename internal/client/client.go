package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mnemoCast-client/internal/models"
	"net/http"
	"time"
)

// Client handles communication with the ad server
type Client struct {
	baseURL    string
	screenID   string
	passkey    string
	httpClient *http.Client
}

// NewClient creates a new ad server client with screen ID and passkey
func NewClient(baseURL string, screenID string, passkey string) *Client {
	return &Client{
		baseURL:  baseURL,
		screenID: screenID,
		passkey:  passkey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetCredentials updates the screen ID and passkey
func (c *Client) SetCredentials(screenID string, passkey string) {
	c.screenID = screenID
	c.passkey = passkey
}

// createRequest creates an HTTP request with authentication headers
func (c *Client) createRequest(method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	// Set authentication headers (server expects X-Screen-Id and X-Screen-Passkey)
	if c.screenID != "" {
		req.Header.Set("X-Screen-Id", c.screenID)  // Server expects X-Screen-Id (not X-Screen-ID)
	}
	if c.passkey != "" {
		req.Header.Set("X-Screen-Passkey", c.passkey)  // Server expects X-Screen-Passkey (not X-Passkey)
	}

	return req, nil
}

// doRequest executes an HTTP request with retry logic
func (c *Client) doRequest(req *http.Request, maxRetries int, retryDelay time.Duration) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay * time.Duration(attempt)) // Exponential backoff
		}

		resp, err := c.httpClient.Do(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries+1, lastErr)
}

// Connect authenticates with the ad server using screen ID and passkey
// This replaces the registration flow - screen is pre-registered on server
func (c *Client) Connect() (*models.Screen, error) {
	url := fmt.Sprintf("%s/api/v1/screens/%s/connect", c.baseURL, c.screenID)

	// Create connection request (empty body, auth via headers)
	req, err := c.createRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	// Execute request with retry
	resp, err := c.doRequest(req, 3, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("authentication failed: invalid screen ID or passkey - %s", string(body))
	}
	if resp.StatusCode == http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("screen not found: screen ID %s does not exist - %s", c.screenID, string(body))
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("connection failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var screen models.Screen
	if err := json.NewDecoder(resp.Body).Decode(&screen); err != nil {
		return nil, fmt.Errorf("failed to parse connection response: %w", err)
	}

	return &screen, nil
}

// Heartbeat sends a heartbeat to the ad server using PUT method
func (c *Client) Heartbeat(screenID string) error {
	url := fmt.Sprintf("%s/api/v1/screens/%s/heartbeat", c.baseURL, screenID)
	requestTime := time.Now()

	// Create heartbeat request
	request := models.HeartbeatRequest{
		Status:    "online",
		Timestamp: requestTime.UTC().Format(time.RFC3339),
	}

	log.Printf("[%s] [REQUEST] Sending heartbeat request to: %s", requestTime.Format("15:04:05.000"), url)
	log.Printf("[%s] [REQUEST] Method: PUT | Screen ID: %s", requestTime.Format("15:04:05.000"), screenID)

	req, err := c.createRequest("PUT", url, request)
	if err != nil {
		log.Printf("[%s] [ERROR] Failed to create heartbeat request: %v", time.Now().Format("15:04:05.000"), err)
		return err
	}

	// Execute request with retry
	resp, err := c.doRequest(req, 3, 2*time.Second)
	responseTime := time.Now()
	duration := responseTime.Sub(requestTime)

	if err != nil {
		log.Printf("[%s] [ERROR] Heartbeat request failed after %v: %v", responseTime.Format("15:04:05.000"), duration, err)
		return fmt.Errorf("heartbeat failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[%s] [ERROR] Heartbeat failed: Status %d | Duration: %v | Response: %s", 
			responseTime.Format("15:04:05.000"), resp.StatusCode, duration, string(body))
		return fmt.Errorf("heartbeat failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[%s] [OK] Heartbeat successful: Status %d | Duration: %v | Path: %s", 
		responseTime.Format("15:04:05.000"), resp.StatusCode, duration, url)
	return nil
}

// GetAds fetches ads from the ad server for the screen
func (c *Client) GetAds(screenID string) (*models.AdDeliveryResponse, error) {
	url := fmt.Sprintf("%s/api/v1/screens/%s/ads/deliver", c.baseURL, screenID)
	requestTime := time.Now()

	log.Printf("[%s] [REQUEST] Fetching ads from: %s", requestTime.Format("15:04:05.000"), url)
	log.Printf("[%s] [REQUEST] Method: GET | Screen ID: %s", requestTime.Format("15:04:05.000"), screenID)

	req, err := c.createRequest("GET", url, nil)
	if err != nil {
		log.Printf("[%s] [ERROR] Failed to create ads request: %v", time.Now().Format("15:04:05.000"), err)
		return nil, err
	}

	// Execute request with retry
	resp, err := c.doRequest(req, 3, 2*time.Second)
	responseTime := time.Now()
	duration := responseTime.Sub(requestTime)

	if err != nil {
		log.Printf("[%s] [ERROR] Ads request failed after %v: %v", responseTime.Format("15:04:05.000"), duration, err)
		return nil, fmt.Errorf("failed to fetch ads: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusNoContent {
		// No ads available
		log.Printf("[%s] [INFO] No ads available (Status 204) | Duration: %v", responseTime.Format("15:04:05.000"), duration)
		return &models.AdDeliveryResponse{
			Ads:      []models.Ad{},
			UpdatedAt: time.Now(),
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[%s] [ERROR] Ads request failed: Status %d | Duration: %v | Response: %s", 
			responseTime.Format("15:04:05.000"), resp.StatusCode, duration, string(body))
		return nil, fmt.Errorf("failed to fetch ads with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var adResponse models.AdDeliveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&adResponse); err != nil {
		log.Printf("[%s] [ERROR] Failed to parse ads response: %v", responseTime.Format("15:04:05.000"), err)
		return nil, fmt.Errorf("failed to parse ads response: %w", err)
	}

	log.Printf("[%s] [OK] Ads fetched successfully: %d ads | Duration: %v | Path: %s", 
		responseTime.Format("15:04:05.000"), len(adResponse.Ads), duration, url)
	return &adResponse, nil
}

