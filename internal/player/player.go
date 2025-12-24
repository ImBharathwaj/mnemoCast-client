package player

import (
	"context"
	"log"
	"mnemoCast-client/internal/ads"
	"mnemoCast-client/internal/models"
	"sync"
	"time"
)

// PlayerState represents the current state of the player
type PlayerState string

const (
	PlayerStateStopped PlayerState = "stopped"
	PlayerStatePlaying PlayerState = "playing"
	PlayerStatePaused  PlayerState = "paused"
	PlayerStateLoading PlayerState = "loading"
	PlayerStateError   PlayerState = "error"
)

// PlayerStats contains statistics about the player
type PlayerStats struct {
	TotalAdsPlayed    int
	CurrentAdID       string
	CurrentAdType     string
	PlaybackStartTime time.Time
	LastError         error
	State             PlayerState
}

// Player orchestrates ad playback
type Player struct {
	playlist   *Playlist
	scheduler  *Scheduler
	downloader *Downloader
	renderer   *RendererManager
	storage    *ads.Storage
	config     *models.ScreenConfig
	
	currentAd  *models.Ad
	state      PlayerState
	stats      PlayerStats
	
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex
	
	// Callback for ad updates
	onAdsUpdated func(*models.AdDeliveryResponse)
}

// NewPlayer creates a new player instance
func NewPlayer(storage *ads.Storage, config *models.ScreenConfig) *Player {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Default scheduler: 30 seconds default duration, 1 second transition delay
	scheduler := NewScheduler(30, 1)
	
	// Create downloader with retry settings from config
	maxRetries := 3
	retryDelay := 5
	if config != nil {
		maxRetries = config.RetryAttempts
		retryDelay = config.RetryDelay
	}
	downloader := NewDownloader(storage, maxRetries, retryDelay)
	
	// Create renderer manager
	renderer := NewRendererManager()
	
	return &Player{
		playlist:   NewPlaylist(),
		scheduler:  scheduler,
		downloader: downloader,
		renderer:   renderer,
		storage:    storage,
		config:     config,
		state:      PlayerStateStopped,
		ctx:        ctx,
		cancel:     cancel,
		stats: PlayerStats{
			State: PlayerStateStopped,
		},
	}
}

// Start starts the player
func (p *Player) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.state == PlayerStatePlaying {
		return nil // Already playing
	}
	
	// Load ads from storage if available
	if ads, err := p.storage.LoadAds(); err == nil {
		p.playlist.UpdateAds(ads)
		log.Printf("[%s] [PLAYER] Loaded %d ads from storage", time.Now().Format("15:04:05.000"), p.playlist.GetCount())
	}
	
	p.state = PlayerStatePlaying
	p.stats.State = PlayerStatePlaying
	
	// Start playback loop
	p.wg.Add(1)
	go p.playbackLoop()
	
	log.Printf("[%s] [PLAYER] Player started", time.Now().Format("15:04:05.000"))
	return nil
}

// Stop stops the player
func (p *Player) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.state == PlayerStateStopped {
		return nil // Already stopped
	}
	
	log.Printf("[%s] [PLAYER] Stopping player...", time.Now().Format("15:04:05.000"))
	
	// Stop renderer
	if p.renderer != nil {
		p.renderer.Stop()
	}
	
	p.cancel()
	p.state = PlayerStateStopped
	p.stats.State = PlayerStateStopped
	p.currentAd = nil
	
	p.mu.Unlock()
	p.wg.Wait()
	p.mu.Lock()
	
	log.Printf("[%s] [PLAYER] Player stopped", time.Now().Format("15:04:05.000"))
	return nil
}

// Pause pauses the player
func (p *Player) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.state != PlayerStatePlaying {
		return nil
	}
	
	p.state = PlayerStatePaused
	p.stats.State = PlayerStatePaused
	log.Printf("[%s] [PLAYER] Player paused", time.Now().Format("15:04:05.000"))
	return nil
}

// Resume resumes the player
func (p *Player) Resume() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.state != PlayerStatePaused {
		return nil
	}
	
	p.state = PlayerStatePlaying
	p.stats.State = PlayerStatePlaying
	log.Printf("[%s] [PLAYER] Player resumed", time.Now().Format("15:04:05.000"))
	return nil
}

// GetState returns the current player state
func (p *Player) GetState() PlayerState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// GetStats returns player statistics
func (p *Player) GetStats() PlayerStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stats
}

// UpdateAds updates the playlist with new ads
func (p *Player) UpdateAds(adResponse *models.AdDeliveryResponse) {
	p.mu.Lock()
	p.playlist.UpdateAds(adResponse)
	p.mu.Unlock()
	
	log.Printf("[%s] [PLAYER] Playlist updated: %d total ads, %d active ads", 
		time.Now().Format("15:04:05.000"), p.playlist.GetCount(), p.playlist.GetActiveCount())
	
	// Call callback if set
	if p.onAdsUpdated != nil {
		p.onAdsUpdated(adResponse)
	}
}

// SetOnAdsUpdated sets a callback for when ads are updated
func (p *Player) SetOnAdsUpdated(callback func(*models.AdDeliveryResponse)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.onAdsUpdated = callback
}

// playbackLoop is the main playback loop
func (p *Player) playbackLoop() {
	defer p.wg.Done()
	
	ticker := time.NewTicker(1 * time.Second) // Check every second
	defer ticker.Stop()
	
	var currentAdStartTime time.Time
	var currentAdDuration time.Duration
	
	for {
		select {
		case <-p.ctx.Done():
			log.Printf("[%s] [PLAYER] Playback loop stopping...", time.Now().Format("15:04:05.000"))
			return
		case <-ticker.C:
			p.mu.Lock()
			state := p.state
			p.mu.Unlock()
			
			if state != PlayerStatePlaying {
				continue
			}
			
			// Check if we need to load a new ad
			if p.currentAd == nil || p.scheduler.ShouldTransition(p.currentAd, currentAdStartTime) {
				p.loadNextAd()
				if p.currentAd != nil {
					currentAdStartTime = time.Now()
					currentAdDuration = p.scheduler.GetAdDuration(p.currentAd)
					
					p.mu.Lock()
					p.stats.CurrentAdID = p.currentAd.ID
					p.stats.CurrentAdType = p.currentAd.Type
					p.stats.PlaybackStartTime = currentAdStartTime
					p.mu.Unlock()
					
					log.Printf("[%s] [PLAYER] Playing ad: %s (type: %s, duration: %v)", 
						time.Now().Format("15:04:05.000"), p.currentAd.ID, p.currentAd.Type, currentAdDuration)
					
					// Wait for transition delay before starting next ad
					time.Sleep(p.scheduler.GetTransitionDelay())
				} else {
					// No ads available, wait a bit before checking again
					time.Sleep(5 * time.Second)
				}
			}
		}
	}
}

// loadNextAd loads the next ad from the playlist
func (p *Player) loadNextAd() {
	p.mu.Lock()
	p.state = PlayerStateLoading
	p.stats.State = PlayerStateLoading
	p.mu.Unlock()
	
	ad := p.playlist.GetNextAd()
	if ad == nil {
		p.mu.Lock()
		p.currentAd = nil
		p.state = PlayerStatePlaying
		p.stats.State = PlayerStatePlaying
		p.mu.Unlock()
		log.Printf("[%s] [PLAYER] No active ads available", time.Now().Format("15:04:05.000"))
		return
	}
	
	// Download media if needed
	localPath, err := p.downloader.DownloadAdMedia(ad)
	if err != nil {
		log.Printf("[%s] [PLAYER] Failed to download media for ad %s: %v", 
			time.Now().Format("15:04:05.000"), ad.ID, err)
		// Try to use cached version if available
		if cachedPath, exists := p.downloader.GetLocalPath(ad); exists {
			localPath = cachedPath
			log.Printf("[%s] [PLAYER] Using cached media: %s", time.Now().Format("15:04:05.000"), localPath)
		} else {
			log.Printf("[%s] [PLAYER] Skipping ad %s: no media available", 
				time.Now().Format("15:04:05.000"), ad.ID)
			return
		}
	}
	
	// Render the ad
	log.Printf("[%s] [PLAYER] Attempting to render ad %s (type: %s) from: %s", 
		time.Now().Format("15:04:05.000"), ad.ID, ad.Type, localPath)
	
	if err := p.renderer.Render(ad, localPath); err != nil {
		log.Printf("[%s] [PLAYER] [ERROR] Failed to render ad %s: %v", 
			time.Now().Format("15:04:05.000"), ad.ID, err)
		log.Printf("[%s] [PLAYER] [ERROR] Ad type: %s, Local path: %s", 
			time.Now().Format("15:04:05.000"), ad.Type, localPath)
		p.mu.Lock()
		p.state = PlayerStateError
		p.stats.State = PlayerStateError
		p.stats.LastError = err
		p.mu.Unlock()
		// Continue to next ad instead of returning
		return
	}
	
	log.Printf("[%s] [PLAYER] [OK] Successfully rendered ad %s", 
		time.Now().Format("15:04:05.000"), ad.ID)
	
	p.mu.Lock()
	p.currentAd = ad
	p.state = PlayerStatePlaying
	p.stats.State = PlayerStatePlaying
	p.stats.TotalAdsPlayed++
	p.mu.Unlock()
}

// GetCurrentAd returns the currently playing ad
func (p *Player) GetCurrentAd() *models.Ad {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.currentAd
}

// GetPlaylist returns the playlist instance
func (p *Player) GetPlaylist() *Playlist {
	return p.playlist
}

