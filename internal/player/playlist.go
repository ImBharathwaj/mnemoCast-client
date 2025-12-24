package player

import (
	"mnemoCast-client/internal/models"
	"sort"
	"sync"
	"time"
)

// Playlist manages the list of ads and provides filtering/sorting
type Playlist struct {
	ads        []models.Ad
	currentIdx int
	lastUpdate time.Time
	mu         sync.RWMutex
}

// NewPlaylist creates a new playlist
func NewPlaylist() *Playlist {
	return &Playlist{
		ads:        []models.Ad{},
		currentIdx: 0,
		lastUpdate: time.Time{},
	}
}

// UpdateAds updates the playlist with new ads from the server
func (p *Playlist) UpdateAds(adResponse *models.AdDeliveryResponse) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.ads = adResponse.Ads
	p.lastUpdate = time.Now()
	
	// Reset current index if it's out of bounds
	if p.currentIdx >= len(p.ads) {
		p.currentIdx = 0
	}
}

// GetActiveAds returns ads that are currently active (within their time window)
func (p *Playlist) GetActiveAds() []models.Ad {
	return p.FilterByTime(time.Now())
}

// FilterByTime filters ads based on current time and their startTime/endTime
// Ads without startTime/endTime are always included
func (p *Playlist) FilterByTime(now time.Time) []models.Ad {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	var activeAds []models.Ad
	
	for _, ad := range p.ads {
		// If ad has no time constraints, it's always active
		if ad.StartTime.IsZero() && ad.EndTime.IsZero() {
			activeAds = append(activeAds, ad)
			continue
		}
		
		// Check if current time is within ad's time window
		startOK := ad.StartTime.IsZero() || now.After(ad.StartTime) || now.Equal(ad.StartTime)
		endOK := ad.EndTime.IsZero() || now.Before(ad.EndTime) || now.Equal(ad.EndTime)
		
		if startOK && endOK {
			activeAds = append(activeAds, ad)
		}
	}
	
	return activeAds
}

// SortByPriority sorts ads by priority (higher priority first)
// If priorities are equal, maintains original order
func (p *Playlist) SortByPriority(ads []models.Ad) []models.Ad {
	sorted := make([]models.Ad, len(ads))
	copy(sorted, ads)
	
	sort.Slice(sorted, func(i, j int) bool {
		// Higher priority comes first
		if sorted[i].Priority != sorted[j].Priority {
			return sorted[i].Priority > sorted[j].Priority
		}
		// If priorities are equal, maintain original order (by ID)
		return sorted[i].ID < sorted[j].ID
	})
	
	return sorted
}

// GetNextAd returns the next ad to play from active ads
// Returns nil if no active ads are available
func (p *Playlist) GetNextAd() *models.Ad {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Get active ads filtered by time
	now := time.Now()
	activeAds := p.FilterByTime(now)
	
	if len(activeAds) == 0 {
		return nil
	}
	
	// Sort by priority
	sortedAds := p.SortByPriority(activeAds)
	
	// If we have ads, cycle through them
	if len(sortedAds) > 0 {
		// Find current ad in sorted list, or start from beginning
		ad := sortedAds[p.currentIdx%len(sortedAds)]
		p.currentIdx++
		return &ad
	}
	
	return nil
}

// GetCount returns the total number of ads in the playlist
func (p *Playlist) GetCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.ads)
}

// GetActiveCount returns the number of currently active ads
func (p *Playlist) GetActiveCount() int {
	return len(p.GetActiveAds())
}

// GetLastUpdate returns the time of the last playlist update
func (p *Playlist) GetLastUpdate() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastUpdate
}

// Reset resets the playlist index to start from the beginning
func (p *Playlist) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.currentIdx = 0
}

