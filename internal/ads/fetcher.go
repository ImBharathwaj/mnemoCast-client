package ads

import (
	"context"
	"encoding/json"
	"log"
	"mnemoCast-client/internal/client"
	"mnemoCast-client/internal/models"
	"sync"
	"time"
)

// Fetcher manages periodic ad fetching from the server
type Fetcher struct {
	client        *client.Client
	screenID      string
	interval      time.Duration
	retryAttempts int
	retryDelay    time.Duration
	storage       *Storage

	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	mu       sync.RWMutex
	lastAds  *models.AdDeliveryResponse
	lastFetch time.Time
	lastError error
}

// NewFetcher creates a new ad fetcher
func NewFetcher(
	adClient *client.Client,
	screenID string,
	configDir string,
	intervalSeconds int,
	retryAttempts int,
	retryDelaySeconds int,
) *Fetcher {
	ctx, cancel := context.WithCancel(context.Background())

	return &Fetcher{
		client:        adClient,
		screenID:      screenID,
		interval:      time.Duration(intervalSeconds) * time.Second,
		retryAttempts: retryAttempts,
		retryDelay:    time.Duration(retryDelaySeconds) * time.Second,
		storage:       NewStorage(configDir),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Start starts the ad fetcher in a background goroutine
func (f *Fetcher) Start() {
	f.wg.Add(1)
	go f.run()
	log.Printf("[ADS] Ad fetcher started (interval: %v)", f.interval)
}

// Stop gracefully stops the ad fetcher
func (f *Fetcher) Stop() {
	log.Println("[ADS] Stopping ad fetcher...")
	f.cancel()
	f.wg.Wait()
	log.Println("[ADS] Ad fetcher stopped")
}

// run is the main fetcher loop
func (f *Fetcher) run() {
	defer f.wg.Done()

	// Fetch ads immediately on start
	log.Printf("[%s] [INIT] Fetching initial ads...", time.Now().Format("15:04:05.000"))
	f.fetchAds()

	// Create ticker for periodic ad fetching
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()

	for {
		select {
		case <-f.ctx.Done():
			log.Printf("[%s] [SHUTDOWN] Ad fetcher stopping...", time.Now().Format("15:04:05.000"))
			return
		case <-ticker.C:
			tickTime := time.Now()
			log.Printf("[%s] [TIMER] Ad fetch interval reached (every %v), fetching ads...", 
				tickTime.Format("15:04:05.000"), f.interval)
			f.fetchAds()
		}
	}
}

// fetchAds fetches ads from the server with retry logic
func (f *Fetcher) fetchAds() {
	startTime := time.Now()
	log.Printf("[%s] [FETCH] Starting ad fetch cycle for screen: %s", startTime.Format("15:04:05.000"), f.screenID)

	var lastErr error

	// Retry logic
	for attempt := 0; attempt <= f.retryAttempts; attempt++ {
		if attempt > 0 {
			delay := f.retryDelay * time.Duration(attempt)
			retryTime := time.Now()
			log.Printf("[%s] [RETRY] Retrying ad fetch (attempt %d/%d) after %v...", 
				retryTime.Format("15:04:05.000"), attempt, f.retryAttempts, delay)
			time.Sleep(delay)
		}

		attemptTime := time.Now()
		log.Printf("[%s] [ATTEMPT] Attempting ad fetch (attempt %d/%d)...", 
			attemptTime.Format("15:04:05.000"), attempt+1, f.retryAttempts+1)

		ads, err := f.client.GetAds(f.screenID)
		if err == nil {
			// Success
			successTime := time.Now()
			totalDuration := successTime.Sub(startTime)

			f.mu.Lock()
			f.lastAds = ads
			f.lastFetch = time.Now()
			f.lastError = nil
			f.mu.Unlock()

			// Save ads to filesystem
			if err := f.storage.SaveAds(ads); err != nil {
				log.Printf("[%s] [WARN] Failed to save ads to filesystem: %v", successTime.Format("15:04:05.000"), err)
			} else {
				log.Printf("[%s] [OK] Ads saved to filesystem: %s", successTime.Format("15:04:05.000"), f.storage.GetAdsDir())
			}

			if attempt > 0 {
				log.Printf("[%s] [OK] Ad fetch succeeded after %d retries | Total duration: %v | Ads received: %d", 
					successTime.Format("15:04:05.000"), attempt, totalDuration, len(ads.Ads))
			} else {
				log.Printf("[%s] [OK] Ad fetch completed successfully | Total duration: %v | Ads received: %d", 
					successTime.Format("15:04:05.000"), totalDuration, len(ads.Ads))
			}

			// Log ad details and JSON response
			if len(ads.Ads) > 0 {
				log.Printf("[%s] [INFO] Received %d ads for display", successTime.Format("15:04:05.000"), len(ads.Ads))
				
				// Log full JSON response
				jsonData, err := json.MarshalIndent(ads, "", "  ")
				if err != nil {
					log.Printf("[%s] [WARN] Failed to marshal ads to JSON: %v", successTime.Format("15:04:05.000"), err)
				} else {
					log.Printf("[%s] [RESPONSE] Ads JSON Response:\n%s", successTime.Format("15:04:05.000"), string(jsonData))
				}
				
				// Log individual ad summary
				for i, ad := range ads.Ads {
					log.Printf("[%s] [INFO] Ad %d: ID=%s, Type=%s, URL=%s", 
						successTime.Format("15:04:05.000"), i+1, ad.ID, ad.Type, ad.ContentURL)
				}
			} else {
				log.Printf("[%s] [INFO] No ads available for display", successTime.Format("15:04:05.000"))
				
				// Log JSON response even when no ads
				jsonData, err := json.MarshalIndent(ads, "", "  ")
				if err != nil {
					log.Printf("[%s] [WARN] Failed to marshal ads to JSON: %v", successTime.Format("15:04:05.000"), err)
				} else {
					log.Printf("[%s] [RESPONSE] Ads JSON Response:\n%s", successTime.Format("15:04:05.000"), string(jsonData))
				}
			}

			return
		}

		lastErr = err
		errorTime := time.Now()
		log.Printf("[%s] [WARN] Ad fetch attempt %d/%d failed: %v", 
			errorTime.Format("15:04:05.000"), attempt+1, f.retryAttempts+1, err)
	}

	// All retries failed
	failureTime := time.Now()
	totalDuration := failureTime.Sub(startTime)

	f.mu.Lock()
	f.lastError = lastErr
	f.mu.Unlock()

	log.Printf("[%s] [ERROR] Ad fetch cycle failed after %d attempts | Total duration: %v | Error: %v", 
		failureTime.Format("15:04:05.000"), f.retryAttempts+1, totalDuration, lastErr)
}

// GetLastAds returns the last successfully fetched ads
func (f *Fetcher) GetLastAds() *models.AdDeliveryResponse {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.lastAds
}

// LoadAdsFromStorage loads ads from the filesystem
func (f *Fetcher) LoadAdsFromStorage() (*models.AdDeliveryResponse, error) {
	return f.storage.LoadAds()
}

// GetStorage returns the ad storage instance
func (f *Fetcher) GetStorage() *Storage {
	return f.storage
}

// GetLastFetch returns the time of the last successful fetch
func (f *Fetcher) GetLastFetch() time.Time {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.lastFetch
}

// GetLastError returns the last error encountered
func (f *Fetcher) GetLastError() error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.lastError
}

// GetStats returns ad fetcher statistics
func (f *Fetcher) GetStats() map[string]interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()

	stats := map[string]interface{}{
		"lastFetch":  f.lastFetch,
		"interval":   f.interval.String(),
		"adsCount":   0,
	}

	if f.lastAds != nil {
		stats["adsCount"] = len(f.lastAds.Ads)
	}

	if f.lastError != nil {
		stats["lastError"] = f.lastError.Error()
	}

	if !f.lastFetch.IsZero() {
		stats["timeSinceLastFetch"] = time.Since(f.lastFetch).String()
	}

	return stats
}

