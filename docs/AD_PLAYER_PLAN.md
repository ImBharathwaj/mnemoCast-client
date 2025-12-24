# Ad Player System - Implementation Plan

## Overview

This document outlines the plan for implementing an ad player system that displays ads fetched from the server. The player will handle multiple ad types (image, video, HTML, text), scheduling, priorities, and continuous playback.

---

## Goals

1. **Play Ads**: Display ads fetched from the server in a continuous loop
2. **Support Multiple Formats**: Handle image, video, HTML, and text ad types
3. **Media Management**: Download and cache ad media files locally
4. **Scheduling**: Respect ad startTime/endTime windows
5. **Priority Handling**: Play ads based on priority levels
6. **Seamless Playback**: Smooth transitions between ads
7. **Offline Support**: Play cached ads when offline
8. **Error Handling**: Graceful handling of missing or corrupted media

---

## Architecture

### Components

```
internal/
├── player/
│   ├── player.go          # Main player orchestrator
│   ├── playlist.go        # Playlist management and filtering
│   ├── scheduler.go       # Ad scheduling logic
│   ├── downloader.go      # Media file downloader
│   ├── renderer.go        # Ad rendering/display interface
│   └── renderers/
│       ├── image.go       # Image ad renderer
│       ├── video.go       # Video ad renderer
│       ├── html.go        # HTML ad renderer
│       └── text.go        # Text ad renderer
```

---

## Phase 1: Core Player Infrastructure

### 1.1 Player Component (`internal/player/player.go`)

**Responsibilities:**
- Orchestrate ad playback
- Manage playback state (playing, paused, stopped)
- Coordinate between playlist, scheduler, and renderers
- Handle playback loop
- Provide player status and statistics

**Key Methods:**
```go
type Player struct {
    playlist    *Playlist
    scheduler   *Scheduler
    downloader  *Downloader
    renderer    Renderer
    currentAd   *models.Ad
    state       PlayerState
    stats       PlayerStats
}

func NewPlayer(storage *ads.Storage, config *models.ScreenConfig) *Player
func (p *Player) Start() error
func (p *Player) Stop() error
func (p *Player) Pause() error
func (p *Player) Resume() error
func (p *Player) GetState() PlayerState
func (p *Player) GetStats() PlayerStats
```

**Player States:**
- `Stopped` - Player is not running
- `Playing` - Currently playing an ad
- `Paused` - Playback paused
- `Loading` - Loading next ad
- `Error` - Error state

### 1.2 Playlist Manager (`internal/player/playlist.go`)

**Responsibilities:**
- Filter ads based on current time (startTime/endTime)
- Sort ads by priority
- Manage ad queue
- Handle playlist updates from fetcher
- Provide next ad to play

**Key Methods:**
```go
type Playlist struct {
    ads         []models.Ad
    currentIdx  int
    lastUpdate  time.Time
    mu          sync.RWMutex
}

func NewPlaylist() *Playlist
func (p *Playlist) UpdateAds(adResponse *models.AdDeliveryResponse)
func (p *Playlist) GetNextAd() *models.Ad
func (p *Playlist) GetActiveAds() []models.Ad
func (p *Playlist) FilterByTime(now time.Time) []models.Ad
func (p *Playlist) SortByPriority(ads []models.Ad) []models.Ad
func (p *Playlist) GetCount() int
```

**Filtering Logic:**
- Filter ads where `now >= startTime && now <= endTime`
- If no startTime/endTime, ad is always active
- Sort by priority (higher priority first)
- If priorities are equal, sort by creation/update time

### 1.3 Scheduler (`internal/player/scheduler.go`)

**Responsibilities:**
- Determine when to switch to next ad
- Handle ad duration timing
- Manage transition delays
- Schedule ad updates

**Key Methods:**
```go
type Scheduler struct {
    defaultDuration time.Duration
    transitionDelay time.Duration
}

func NewScheduler(defaultDuration, transitionDelay time.Duration) *Scheduler
func (s *Scheduler) GetAdDuration(ad *models.Ad) time.Duration
func (s *Scheduler) ShouldTransition(ad *models.Ad, startTime time.Time) bool
```

**Duration Logic:**
- Use `ad.Duration` if specified (in seconds)
- Fall back to `defaultDuration` from config (e.g., 30 seconds)
- Minimum duration: 5 seconds
- Maximum duration: 300 seconds (5 minutes)

---

## Phase 2: Media Downloader

### 2.1 Downloader Component (`internal/player/downloader.go`)

**Responsibilities:**
- Download ad media files from URLs
- Cache media files locally
- Verify downloaded files
- Handle download retries
- Check for existing cached files

**Key Methods:**
```go
type Downloader struct {
    storage     *ads.Storage
    httpClient  *http.Client
    maxRetries  int
    retryDelay  time.Duration
}

func NewDownloader(storage *ads.Storage) *Downloader
func (d *Downloader) DownloadAdMedia(ad *models.Ad) (string, error)
func (d *Downloader) GetLocalPath(ad *models.Ad) (string, bool)
func (d *Downloader) IsCached(ad *models.Ad) bool
func (d *Downloader) DownloadFile(url, destPath string) error
```

**Download Strategy:**
- Check if media already cached locally
- If cached and valid, use cached version
- If not cached or invalid, download from URL
- Store in `~/.mnemocast/ads/media/{adId}/` directory
- Extract filename from URL or use ad ID + extension
- Verify file integrity (file size > 0, valid format)

**Supported Media Types:**
- Images: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`
- Videos: `.mp4`, `.webm`, `.mov`, `.avi`
- HTML: `.html` files or inline HTML content

**File Naming:**
- Images: `{adId}.{ext}` or `{adId}_image.{ext}`
- Videos: `{adId}.{ext}` or `{adId}_video.{ext}`
- HTML: `{adId}.html`

---

## Phase 3: Renderer System

### 3.1 Renderer Interface (`internal/player/renderer.go`)

**Responsibilities:**
- Define common interface for all renderers
- Handle renderer selection based on ad type
- Manage renderer lifecycle

**Interface:**
```go
type Renderer interface {
    CanRender(ad *models.Ad) bool
    Render(ad *models.Ad, localPath string) error
    Stop() error
    GetStatus() RendererStatus
}

type RendererStatus struct {
    IsPlaying bool
    Error     error
}
```

### 3.2 Image Renderer (`internal/player/renderers/image.go`)

**Responsibilities:**
- Display image ads
- Handle image formats (JPEG, PNG, GIF, WebP)
- Support fullscreen display
- Handle image transitions

**Implementation Options:**
1. **CLI Approach**: Use system commands (e.g., `feh`, `sxiv`, `imv` on Linux)
2. **Future GUI**: Prepare for GUI framework integration
3. **Web Server**: Start local HTTP server and open browser
4. **Terminal Display**: ASCII art or terminal image viewers

**Initial Approach: CLI with System Commands**
- Use `feh` or `imv` for Linux (fullscreen image viewer)
- Fallback to opening default image viewer
- For headless systems, prepare for future GUI integration

**Key Methods:**
```go
type ImageRenderer struct {
    currentProcess *os.Process
    displayCmd     string
}

func NewImageRenderer() *ImageRenderer
func (r *ImageRenderer) CanRender(ad *models.Ad) bool
func (r *ImageRenderer) Render(ad *models.Ad, localPath string) error
func (r *ImageRenderer) Stop() error
```

### 3.3 Video Renderer (`internal/player/renderers/video.go`)

**Responsibilities:**
- Play video ads
- Handle video formats (MP4, WebM, MOV, AVI)
- Support fullscreen playback
- Handle video loops and transitions

**Implementation Options:**
1. **CLI Approach**: Use `mpv`, `vlc`, or `ffplay`
2. **Future GUI**: Prepare for GUI framework integration
3. **Web Server**: Use HTML5 video player

**Initial Approach: CLI with mpv/vlc**
- Use `mpv` for Linux (lightweight, scriptable)
- Fallback to `vlc` or system default player
- Support fullscreen and loop options

**Key Methods:**
```go
type VideoRenderer struct {
    currentProcess *os.Process
    playerCmd      string
}

func NewVideoRenderer() *VideoRenderer
func (r *VideoRenderer) CanRender(ad *models.Ad) bool
func (r *VideoRenderer) Render(ad *models.Ad, localPath string) error
func (r *VideoRenderer) Stop() error
```

### 3.4 HTML Renderer (`internal/player/renderers/html.go`)

**Responsibilities:**
- Display HTML ads
- Handle inline HTML or HTML files
- Support embedded media
- Fullscreen browser display

**Implementation Options:**
1. **Local Web Server**: Start HTTP server, open browser
2. **Future GUI**: Use embedded webview
3. **CLI**: Use `w3m`, `lynx` (limited)

**Initial Approach: Local Web Server**
- Start local HTTP server on random port
- Serve HTML content
- Open default browser in fullscreen/kiosk mode
- Clean up on ad change

**Key Methods:**
```go
type HTMLRenderer struct {
    server     *http.Server
    serverPort int
    browserCmd string
}

func NewHTMLRenderer() *HTMLRenderer
func (r *HTMLRenderer) CanRender(ad *models.Ad) bool
func (r *HTMLRenderer) Render(ad *models.Ad, localPath string) error
func (r *HTMLRenderer) Stop() error
```

### 3.5 Text Renderer (`internal/player/renderers/text.go`)

**Responsibilities:**
- Display text ads
- Format text nicely
- Support scrolling text
- Terminal or GUI display

**Implementation Options:**
1. **Terminal Display**: Use ANSI colors and formatting
2. **Future GUI**: Text widget with styling
3. **Fullscreen Terminal**: Use `tput` for positioning

**Initial Approach: Terminal Display**
- Display in terminal with formatting
- Support scrolling for long text
- Use colors and borders for visibility

**Key Methods:**
```go
type TextRenderer struct {
    terminalWidth  int
    terminalHeight int
}

func NewTextRenderer() *TextRenderer
func (r *TextRenderer) CanRender(ad *models.Ad) bool
func (r *TextRenderer) Render(ad *models.Ad, localPath string) error
func (r *TextRenderer) Stop() error
```

---

## Phase 4: Integration

### 4.1 Player Integration with Fetcher

**Integration Points:**
- Player subscribes to ad updates from fetcher
- When new ads arrive, update playlist
- Handle ad removal (ads no longer in response)
- Maintain playback continuity

**Implementation:**
```go
// In fetcher.go
func (f *Fetcher) SetPlayerCallback(callback func(*models.AdDeliveryResponse))

// In player.go
func (p *Player) OnAdsUpdated(adResponse *models.AdDeliveryResponse) {
    p.playlist.UpdateAds(adResponse)
}
```

### 4.2 Main Application Integration

**Integration in `cmd/screen/main.go`:**
- Initialize player after fetcher starts
- Start player when ads are available
- Handle player lifecycle (start/stop)
- Display player status in main loop

**Flow:**
1. Start ad fetcher
2. Wait for initial ads
3. Initialize player with ads
4. Start player
5. Update player when new ads arrive
6. Handle graceful shutdown

---

## Phase 5: Configuration

### 5.1 Player Configuration

**Add to `internal/models/config.go`:**
```go
type ScreenConfig struct {
    // ... existing fields ...
    
    // Player settings
    PlayerEnabled        bool   `json:"playerEnabled"`        // Enable/disable player
    DefaultAdDuration    int    `json:"defaultAdDuration"`    // Default duration in seconds (default: 30)
    TransitionDelay      int    `json:"transitionDelay"`      // Delay between ads in seconds (default: 1)
    MediaCacheEnabled    bool   `json:"mediaCacheEnabled"`    // Enable media caching (default: true)
    MediaCacheMaxSize    int64  `json:"mediaCacheMaxSize"`    // Max cache size in MB (default: 500)
    
    // Renderer settings
    ImageViewer          string `json:"imageViewer"`          // Image viewer command (default: "feh")
    VideoPlayer          string `json:"videoPlayer"`          // Video player command (default: "mpv")
    Browser              string `json:"browser"`              // Browser command (default: "firefox")
    Fullscreen           bool   `json:"fullscreen"`           // Fullscreen mode (default: true)
}
```

---

## Phase 6: Error Handling & Edge Cases

### 6.1 Error Scenarios

**Media Download Failures:**
- Network errors → Use cached version if available
- Invalid URL → Skip ad, log error
- Unsupported format → Skip ad, log warning
- File corruption → Re-download, retry

**Rendering Failures:**
- Renderer not available → Skip ad, try next
- Process crash → Restart renderer, continue
- Display unavailable → Log error, pause player

**Playlist Scenarios:**
- No ads available → Show placeholder or wait
- All ads expired → Wait for next fetch
- Single ad → Loop single ad
- Empty playlist → Pause player, wait for ads

### 6.2 Recovery Mechanisms

- **Automatic Retry**: Retry failed downloads with exponential backoff
- **Fallback Rendering**: Try alternative renderers if primary fails
- **Cache Validation**: Verify cached files before use
- **Health Checks**: Periodic player health checks

---

## Phase 7: Logging & Monitoring

### 7.1 Player Logging

**Log Events:**
- Player start/stop
- Ad playback start/end
- Media download start/success/failure
- Renderer selection and status
- Playlist updates
- Error conditions

**Log Format:**
```
[HH:MM:SS.mmm] [PLAYER] Player started
[HH:MM:SS.mmm] [PLAYER] Playing ad: ad-12345 (type: image, duration: 30s)
[HH:MM:SS.mmm] [DOWNLOAD] Downloading media: https://example.com/ad.jpg
[HH:MM:SS.mmm] [DOWNLOAD] Media cached: ~/.mnemocast/ads/media/ad-12345/ad-12345.jpg
[HH:MM:SS.mmm] [RENDER] Rendering image ad: ad-12345
[HH:MM:SS.mmm] [PLAYER] Ad completed: ad-12345
[HH:MM:SS.mmm] [PLAYER] Transitioning to next ad...
```

### 7.2 Player Statistics

**Track:**
- Total ads played
- Current ad being played
- Playback duration
- Download statistics
- Error counts
- Cache hit/miss ratio

---

## Implementation Phases

### Phase 1: Core Infrastructure (Priority: High) ✅ COMPLETED
- [x] Create `internal/player/player.go`
- [x] Create `internal/player/playlist.go`
- [x] Create `internal/player/scheduler.go`
- [x] Basic player lifecycle (start/stop)
- [x] Playlist filtering and sorting

**Status:** ✅ Complete
- Player orchestrator with state management
- Playlist manager with time-based filtering and priority sorting
- Scheduler for ad duration and transitions

### Phase 2: Media Downloader (Priority: High) ✅ COMPLETED
- [x] Create `internal/player/downloader.go`
- [x] Implement download logic
- [x] Implement caching mechanism
- [x] Add retry logic
- [x] File validation

**Status:** ✅ Complete
- HTTP downloader with retry mechanism
- Local file caching in `~/.mnemocast/ads/media/{adId}/`
- File extension detection from URL or ad type
- Cleanup of old media files

### Phase 3: Renderer System (Priority: High) ✅ COMPLETED
- [x] Create renderer interface
- [x] Implement image renderer
- [x] Implement video renderer
- [x] Implement HTML renderer
- [x] Implement text renderer
- [x] Renderer selection logic

**Status:** ✅ Complete
- Renderer interface and manager
- Image renderer: Uses `feh`, `imv`, `sxiv`, or `xdg-open`
- Video renderer: Uses `mpv`, `vlc`, `ffplay`, or `xdg-open`
- HTML renderer: Local HTTP server + browser
- Text renderer: Terminal display with formatting

### Phase 4: Integration (Priority: Medium) ✅ COMPLETED
- [x] Integrate player with fetcher
- [x] Integrate player with main application
- [x] Handle ad updates
- [x] Graceful shutdown

**Status:** ✅ Complete
- Player integrated with ad fetcher via callback
- Player starts automatically when ads are available
- Updates playlist when new ads arrive
- Graceful shutdown on Ctrl+C

### Phase 5: Configuration (Priority: Medium)
- [ ] Add player config to models
- [ ] Load player config
- [ ] Apply config defaults
- [ ] Config validation

**Status:** ⏳ Pending
- Currently using hardcoded defaults (30s duration, 1s transition)
- Configuration can be added to `ScreenConfig` model

### Phase 6: Error Handling (Priority: Medium)
- [x] Basic error handling implemented
- [ ] Enhanced error recovery
- [ ] Handle edge cases (no ads, all expired, etc.)
- [ ] Add health checks
- [ ] Fallback mechanisms

**Status:** 🟡 Partial
- Basic error handling in place
- Player handles missing ads gracefully
- Media download failures use cached versions
- Enhanced recovery and edge cases pending

### Phase 7: Logging & Monitoring (Priority: Low)
- [x] Basic logging implemented
- [ ] Enhanced player statistics
- [ ] Status reporting in main loop
- [ ] Debug mode

**Status:** 🟡 Partial
- Basic logging for player events
- Statistics tracking (total ads played, current ad)
- Status reporting can be enhanced in main loop

---

## Technical Considerations

### Display Options

**Current (CLI):**
- Use system commands for image/video display
- Terminal-based text rendering
- Local web server for HTML

**Future (GUI):**
- Prepare for GUI framework integration (Wails, Fyne, etc.)
- Abstract renderer interface for easy migration
- Keep business logic separate from display

### Dependencies

**New Dependencies:**
- None required initially (use standard library + system commands)
- Future: May need GUI framework or webview library

**System Requirements:**
- Image viewer: `feh`, `imv`, or system default
- Video player: `mpv`, `vlc`, or system default
- Browser: `firefox`, `chromium`, or system default

### Performance

- **Memory**: Cache media files, not entire files in memory
- **CPU**: Lightweight rendering, efficient downloads
- **Disk**: Manage cache size, cleanup old files
- **Network**: Download in background, don't block playback

---

## Testing Strategy

### Unit Tests
- Playlist filtering and sorting
- Scheduler duration calculation
- Downloader caching logic
- Renderer selection

### Integration Tests
- Player with fetcher
- Player with storage
- End-to-end ad playback

### Manual Testing
- Different ad types
- Scheduling scenarios
- Error conditions
- Offline playback

---

## Future Enhancements

1. **GUI Integration**: Full GUI player with Wails or Fyne
2. **Advanced Scheduling**: Time-based playlists, day-of-week rules
3. **Analytics**: Track ad impressions, play duration
4. **Remote Control**: Control player via API or web interface
5. **Multi-Monitor**: Support multiple displays
6. **Transitions**: Fade, slide, and other transition effects
7. **Audio Support**: Audio ads, background music
8. **Interactive Ads**: Clickable ads, touch support

---

## Summary

This plan outlines a comprehensive ad player system that:
- Plays ads in a continuous loop
- Supports multiple ad formats (image, video, HTML, text)
- Downloads and caches media files
- Handles scheduling and priorities
- Provides robust error handling
- Integrates seamlessly with existing fetcher

The implementation will start with CLI-based rendering using system commands, with a clear path to future GUI integration.

