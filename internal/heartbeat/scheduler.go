package heartbeat

import (
	"context"
	"log"
	"mnemoCast-client/internal/client"
	"mnemoCast-client/internal/identity"
	"sync"
	"time"
)

// Status represents the connection status
type Status int

const (
	StatusUnknown Status = iota
	StatusConnected
	StatusDisconnected
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusConnected:
		return "Connected"
	case StatusDisconnected:
		return "Disconnected"
	case StatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// Scheduler manages periodic heartbeat sending
type Scheduler struct {
	client          *client.Client
	identityManager *identity.Manager
	screenID        string
	interval        time.Duration
	retryAttempts   int
	retryDelay      time.Duration

	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	mu       sync.RWMutex
	status   Status
	lastSent time.Time
	lastError error
}

// NewScheduler creates a new heartbeat scheduler
func NewScheduler(
	adClient *client.Client,
	identityManager *identity.Manager,
	screenID string,
	intervalSeconds int,
	retryAttempts int,
	retryDelaySeconds int,
) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		client:          adClient,
		identityManager: identityManager,
		screenID:        screenID,
		interval:        time.Duration(intervalSeconds) * time.Second,
		retryAttempts:   retryAttempts,
		retryDelay:      time.Duration(retryDelaySeconds) * time.Second,
		ctx:             ctx,
		cancel:          cancel,
		status:          StatusUnknown,
	}
}

// Start starts the heartbeat scheduler
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go s.run()
	log.Printf("Heartbeat scheduler started (interval: %v)", s.interval)
}

// Stop stops the heartbeat scheduler
func (s *Scheduler) Stop() {
	log.Println("Stopping heartbeat scheduler...")
	s.cancel()
	s.wg.Wait()
	log.Println("Heartbeat scheduler stopped")
}

// run is the main scheduler loop
func (s *Scheduler) run() {
	defer s.wg.Done()

	// Send initial heartbeat immediately
	log.Printf("[%s] [INIT] Sending initial heartbeat...", time.Now().Format("15:04:05.000"))
	s.sendHeartbeat()

	// Create ticker for periodic heartbeats
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("[%s] [SHUTDOWN] Heartbeat scheduler stopping...", time.Now().Format("15:04:05.000"))
			return
		case <-ticker.C:
			tickTime := time.Now()
			log.Printf("[%s] [TIMER] Heartbeat interval reached (every %v), sending heartbeat...", 
				tickTime.Format("15:04:05.000"), s.interval)
			s.sendHeartbeat()
		}
	}
}

// sendHeartbeat sends a heartbeat with retry logic
func (s *Scheduler) sendHeartbeat() {
	startTime := time.Now()
	log.Printf("[%s] [HB] Starting heartbeat cycle for screen: %s", startTime.Format("15:04:05.000"), s.screenID)
	
	var lastErr error

	// Retry logic
	for attempt := 0; attempt <= s.retryAttempts; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := s.retryDelay * time.Duration(attempt)
			retryTime := time.Now()
			log.Printf("[%s] [RETRY] Retrying heartbeat (attempt %d/%d) after %v...", 
				retryTime.Format("15:04:05.000"), attempt, s.retryAttempts, delay)
			time.Sleep(delay)
		}

		attemptTime := time.Now()
		log.Printf("[%s] [ATTEMPT] Attempting heartbeat (attempt %d/%d)...", 
			attemptTime.Format("15:04:05.000"), attempt+1, s.retryAttempts+1)
		
		err := s.client.Heartbeat(s.screenID)
		if err == nil {
			// Success
			successTime := time.Now()
			totalDuration := successTime.Sub(startTime)
			
			s.mu.Lock()
			s.status = StatusConnected
			s.lastSent = time.Now()
			s.lastError = nil
			s.mu.Unlock()

			// Update last seen in identity
			if identity, err := s.identityManager.LoadIdentity(); err == nil {
				_ = s.identityManager.UpdateLastSeen(identity)
			}

			if attempt > 0 {
				log.Printf("[%s] [OK] Heartbeat succeeded after %d retries | Total duration: %v", 
					successTime.Format("15:04:05.000"), attempt, totalDuration)
			} else {
				log.Printf("[%s] [OK] Heartbeat cycle completed successfully | Total duration: %v", 
					successTime.Format("15:04:05.000"), totalDuration)
			}
			return
		}

		lastErr = err
		errorTime := time.Now()
		log.Printf("[%s] [WARN] Heartbeat attempt %d/%d failed: %v", 
			errorTime.Format("15:04:05.000"), attempt+1, s.retryAttempts+1, err)
	}

	// All retries failed
	failureTime := time.Now()
	totalDuration := failureTime.Sub(startTime)
	
	s.mu.Lock()
	s.status = StatusError
	s.lastError = lastErr
	s.mu.Unlock()

	log.Printf("[%s] [ERROR] Heartbeat cycle failed after %d attempts | Total duration: %v | Error: %v", 
		failureTime.Format("15:04:05.000"), s.retryAttempts+1, totalDuration, lastErr)
}

// GetStatus returns the current connection status
func (s *Scheduler) GetStatus() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// GetLastSent returns the time of the last successful heartbeat
func (s *Scheduler) GetLastSent() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSent
}

// GetLastError returns the last error encountered
func (s *Scheduler) GetLastError() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastError
}

// IsConnected returns true if the scheduler is connected
func (s *Scheduler) IsConnected() bool {
	return s.GetStatus() == StatusConnected
}

// GetStats returns heartbeat statistics
func (s *Scheduler) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"status":    s.status.String(),
		"lastSent":  s.lastSent,
		"interval":  s.interval.String(),
		"connected": s.status == StatusConnected,
	}

	if s.lastError != nil {
		stats["lastError"] = s.lastError.Error()
	}

	if !s.lastSent.IsZero() {
		stats["timeSinceLastSent"] = time.Since(s.lastSent).String()
	}

	return stats
}

