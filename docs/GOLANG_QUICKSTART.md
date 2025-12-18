# 🐹 Golang + Wails Quick Start Guide

**Step-by-step guide to get your Digital Display Client up and running with Golang.**

---

## 🎯 Prerequisites

### Required Software

```bash
# 1. Install Go (1.21 or later)
# Download from: https://go.dev/dl/
go version  # Should show 1.21+

# 2. Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 3. Verify Wails installation
wails version

# 4. Install Node.js (for frontend, 18+)
node --version  # Should show v18+
```

### Platform-Specific Requirements

**Linux:**
```bash
# Install build dependencies
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev
```

**Windows:**
- Install TDM-GCC or MinGW-w64
- Wails will guide you through setup

**macOS:**
```bash
# Install Xcode Command Line Tools
xcode-select --install
```

---

## 🚀 Step 1: Create Project

```bash
# Create new Wails project with vanilla JS template
wails init -n mnemoCast-client -t vanilla

# Navigate to project
cd mnemoCast-client

# Initialize Go module
go mod init mnemoCast-client
```

---

## 📦 Step 2: Install Dependencies

```bash
# Add required Go packages
go get github.com/wailsapp/wails/v2
go get github.com/go-resty/resty/v2
go get github.com/spf13/viper

# Or add to go.mod manually, then run:
go mod tidy
```

**`go.mod` should look like:**
```go
module mnemoCast-client

go 1.21

require (
    github.com/wailsapp/wails/v2 v2.8.0
    github.com/go-resty/resty/v2 v2.11.0
    github.com/spf13/viper v1.18.2
)
```

---

## 🏗️ Step 3: Project Structure

Create this structure:

```
mnemoCast-client/
├── app/
│   ├── main.go
│   └── internal/
│       ├── app.go
│       ├── models/
│       │   ├── screen.go
│       │   ├── playlist.go
│       │   └── events.go
│       ├── api/
│       │   ├── client.go
│       │   ├── screen.go
│       │   ├── playlist.go
│       │   └── events.go
│       └── storage/
│           ├── cache.go
│           ├── media_cache.go
│           └── event_queue.go
├── frontend/
│   ├── index.html
│   ├── main.js
│   ├── styles.css
│   └── components/
│       ├── Player.js
│       └── StatusBar.js
├── wails.json
└── go.mod
```

---

## 📝 Step 4: Core Files Setup

### 4.1 Models

**`app/internal/models/screen.go`**
```go
package models

type ScreenLocation struct {
    City      string `json:"city"`
    Area      string `json:"area"`
    VenueType string `json:"venueType"`
}

type Screen struct {
    ID            string         `json:"id"`
    Name          string         `json:"name"`
    Location      ScreenLocation `json:"location"`
    Classification int           `json:"classification"`
    IsOnline      bool           `json:"isOnline"`
    LastHeartbeat *string        `json:"lastHeartbeat,omitempty"`
}

type RegisterScreenRequest struct {
    ID            string         `json:"id"`
    Name          string         `json:"name"`
    Location      ScreenLocation `json:"location"`
    Classification int           `json:"classification"`
}
```

**`app/internal/models/playlist.go`**
```go
package models

type PlaylistItem struct {
    CreativeID     string `json:"creativeId"`
    CampaignID     string `json:"campaignId"`
    URL            string `json:"url"`
    DurationSeconds int   `json:"durationSeconds"`
    Order          int    `json:"order"`
    Type           string `json:"type"` // "image" or "video"
}

type Playlist struct {
    ScreenID       string        `json:"screenId"`
    GeneratedAt    string        `json:"generatedAt"`
    DurationMinutes int          `json:"durationMinutes"`
    Items          []PlaylistItem `json:"items"`
}
```

**`app/internal/models/events.go`**
```go
package models

type PlayEvent struct {
    ScreenID   string `json:"screenId"`
    CreativeID string `json:"creativeId"`
    CampaignID string `json:"campaignId"`
    Timestamp  string `json:"timestamp"`
}
```

### 4.2 API Client

**`app/internal/api/client.go`**
```go
package api

import (
    "net/http"
    "time"
)

type Client struct {
    httpClient *http.Client
    baseURL    string
}

func NewClient(baseURL string) *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        baseURL: baseURL,
    }
}

func (c *Client) GetHTTPClient() *http.Client {
    return c.httpClient
}

func (c *Client) GetBaseURL() string {
    return c.baseURL
}
```

**`app/internal/api/screen.go`**
```go
package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "mnemoCast-client/app/internal/models"
    "net/http"
)

type ScreenAPI struct {
    client *Client
}

func NewScreenAPI(baseURL string) *ScreenAPI {
    return &ScreenAPI{
        client: NewClient(baseURL),
    }
}

func (s *ScreenAPI) Register(request models.RegisterScreenRequest) (*models.Screen, error) {
    url := fmt.Sprintf("%s/api/v1/screens/register", s.client.GetBaseURL())
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, err
    }
    
    resp, err := s.client.GetHTTPClient().Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var screen models.Screen
    if err := json.NewDecoder(resp.Body).Decode(&screen); err != nil {
        return nil, err
    }
    
    return &screen, nil
}

func (s *ScreenAPI) Heartbeat(screenID string) error {
    url := fmt.Sprintf("%s/api/v1/screens/%s/heartbeat", s.client.GetBaseURL(), screenID)
    
    resp, err := s.client.GetHTTPClient().Post(url, "application/json", nil)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("heartbeat failed with status: %d", resp.StatusCode)
    }
    
    return nil
}
```

**`app/internal/api/playlist.go`**
```go
package api

import (
    "encoding/json"
    "fmt"
    "mnemoCast-client/app/internal/models"
    "net/http"
    "net/url"
)

type PlaylistAPI struct {
    client *Client
}

func NewPlaylistAPI(baseURL string) *PlaylistAPI {
    return &PlaylistAPI{
        client: NewClient(baseURL),
    }
}

func (p *PlaylistAPI) Fetch(screenID string, durationMinutes int) (*models.Playlist, error) {
    baseURL := fmt.Sprintf("%s/api/v1/screens/%s/playlist", p.client.GetBaseURL(), screenID)
    u, err := url.Parse(baseURL)
    if err != nil {
        return nil, err
    }
    
    q := u.Query()
    q.Set("durationMinutes", fmt.Sprintf("%d", durationMinutes))
    u.RawQuery = q.Encode()
    
    resp, err := p.client.GetHTTPClient().Get(u.String())
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var playlist models.Playlist
    if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
        return nil, err
    }
    
    return &playlist, nil
}
```

**`app/internal/api/events.go`**
```go
package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "mnemoCast-client/app/internal/models"
    "net/http"
)

type EventsAPI struct {
    client *Client
}

func NewEventsAPI(baseURL string) *EventsAPI {
    return &EventsAPI{
        client: NewClient(baseURL),
    }
}

func (e *EventsAPI) RecordImpression(event models.PlayEvent) error {
    url := fmt.Sprintf("%s/api/v1/events/impression", e.client.GetBaseURL())
    return e.sendEvent(url, event)
}

func (e *EventsAPI) RecordPlay(event models.PlayEvent) error {
    url := fmt.Sprintf("%s/api/v1/events/play", e.client.GetBaseURL())
    return e.sendEvent(url, event)
}

func (e *EventsAPI) sendEvent(url string, event models.PlayEvent) error {
    jsonData, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    resp, err := e.client.GetHTTPClient().Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("event send failed with status: %d", resp.StatusCode)
    }
    
    return nil
}
```

### 4.3 App Struct

**`app/internal/app.go`**
```go
package internal

import (
    "context"
    "encoding/json"
    "mnemoCast-client/app/internal/api"
    "mnemoCast-client/app/internal/models"
    "sync"
    "time"
)

type App struct {
    ctx           context.Context
    screenAPI     *api.ScreenAPI
    playlistAPI   *api.PlaylistAPI
    eventsAPI     *api.EventsAPI
    config        *Config
    configMutex   sync.RWMutex
    heartbeatStop chan bool
}

type Config struct {
    ScreenID                string              `json:"screenId"`
    Name                    string              `json:"name"`
    Location                models.ScreenLocation `json:"location"`
    Classification          int                 `json:"classification"`
    BackendURL              string              `json:"backendUrl"`
    PlaylistRefreshInterval int                 `json:"playlistRefreshInterval"`
    HeartbeatInterval       int                 `json:"heartbeatInterval"`
}

func NewApp() *App {
    defaultConfig := &Config{
        ScreenID:                "screen-1",
        Name:                    "Default Screen",
        Location:                models.ScreenLocation{City: "Chennai", Area: "Airport", VenueType: "airport"},
        Classification:          1,
        BackendURL:              "http://localhost:8080",
        PlaylistRefreshInterval: 3,
        HeartbeatInterval:       30,
    }
    
    return &App{
        screenAPI:     api.NewScreenAPI(defaultConfig.BackendURL),
        playlistAPI:   api.NewPlaylistAPI(defaultConfig.BackendURL),
        eventsAPI:     api.NewEventsAPI(defaultConfig.BackendURL),
        config:        defaultConfig,
        heartbeatStop: make(chan bool),
    }
}

func (a *App) OnStartup(ctx context.Context) {
    a.ctx = ctx
    a.startHeartbeat()
    a.registerScreen()
}

func (a *App) OnDomReady(ctx context.Context) {
    // Frontend is ready
}

func (a *App) OnBeforeClose(ctx context.Context) (prevent bool) {
    a.stopHeartbeat()
    return false
}

// Wails exposed methods

func (a *App) RegisterScreen(request models.RegisterScreenRequest) (string, error) {
    screen, err := a.screenAPI.Register(request)
    if err != nil {
        return "", err
    }
    return toJSON(screen), nil
}

func (a *App) FetchPlaylist(screenID string, durationMinutes int) (string, error) {
    playlist, err := a.playlistAPI.Fetch(screenID, durationMinutes)
    if err != nil {
        return "", err
    }
    return toJSON(playlist), nil
}

func (a *App) RecordImpression(eventJSON string) error {
    var event models.PlayEvent
    if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
        return err
    }
    return a.eventsAPI.RecordImpression(event)
}

func (a *App) RecordPlay(eventJSON string) error {
    var event models.PlayEvent
    if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
        return err
    }
    return a.eventsAPI.RecordPlay(event)
}

func (a *App) GetConfig() (string, error) {
    a.configMutex.RLock()
    defer a.configMutex.RUnlock()
    return toJSON(a.config), nil
}

func (a *App) SetConfig(configJSON string) error {
    var config Config
    if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
        return err
    }
    
    a.configMutex.Lock()
    defer a.configMutex.Unlock()
    a.config = &config
    
    // Update API clients with new base URL
    a.screenAPI = api.NewScreenAPI(config.BackendURL)
    a.playlistAPI = api.NewPlaylistAPI(config.BackendURL)
    a.eventsAPI = api.NewEventsAPI(config.BackendURL)
    
    return nil
}

// Internal methods

func (a *App) registerScreen() {
    go func() {
        request := models.RegisterScreenRequest{
            ID:            a.config.ScreenID,
            Name:          a.config.Name,
            Location:      a.config.Location,
            Classification: a.config.Classification,
        }
        
        if _, err := a.screenAPI.Register(request); err != nil {
            println("Screen registration failed:", err.Error())
        } else {
            println("Screen registered successfully")
        }
    }()
}

func (a *App) startHeartbeat() {
    go func() {
        ticker := time.NewTicker(time.Duration(a.config.HeartbeatInterval) * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                if err := a.screenAPI.Heartbeat(a.config.ScreenID); err != nil {
                    println("Heartbeat failed:", err.Error())
                } else {
                    println("Heartbeat sent")
                }
            case <-a.heartbeatStop:
                return
            }
        }
    }()
}

func (a *App) stopHeartbeat() {
    close(a.heartbeatStop)
}

func toJSON(v interface{}) (string, error) {
    data, err := json.Marshal(v)
    if err != nil {
        return "", err
    }
    return string(data), nil
}
```

### 4.4 Main Entry Point

**`app/main.go`**
```go
package main

import (
    "context"
    "embed"
    "mnemoCast-client/app/internal"
    
    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := internal.NewApp()
    
    err := wails.Run(&options.App{
        Title:  "MnemoCast Display Client",
        Width:  1920,
        Height: 1080,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
        OnStartup:        app.OnStartup,
        OnDomReady:       app.OnDomReady,
        OnBeforeClose:    app.OnBeforeClose,
        Context:          context.Background(),
    })
    
    if err != nil {
        println("Error:", err.Error())
    }
}
```

---

## 🎨 Step 5: Frontend Setup

### 5.1 HTML

**`frontend/index.html`**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MnemoCast Display Client</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <div id="app">
        <div id="status-bar"></div>
        <div id="player-container">
            <div id="loading">Loading playlist...</div>
        </div>
    </div>
    <script src="wailsjs/runtime/runtime.js"></script>
    <script type="module" src="main.js"></script>
</body>
</html>
```

### 5.2 Main JavaScript

**`frontend/main.js`**
```javascript
import { Player } from './components/Player.js';
import { StatusBar } from './components/StatusBar.js';

class DisplayClient {
    constructor() {
        this.config = null;
        this.playlist = null;
        this.currentIndex = 0;
        this.player = new Player(document.getElementById('player-container'));
        this.statusBar = new StatusBar(document.getElementById('status-bar'));
    }

    async init() {
        try {
            const configJSON = await window.go.main.App.GetConfig();
            this.config = JSON.parse(configJSON);
            
            await this.registerScreen();
            await this.fetchAndPlay();
            this.updateStatusBar();
        } catch (error) {
            console.error('Initialization error:', error);
            this.showError('Failed to initialize: ' + error);
        }
    }

    async registerScreen() {
        try {
            const request = {
                id: this.config.screenId,
                name: this.config.name,
                location: this.config.location,
                classification: this.config.classification,
            };
            
            await window.go.main.App.RegisterScreen(JSON.stringify(request));
            console.log('Screen registered successfully');
        } catch (error) {
            console.error('Registration error:', error);
        }
    }

    async fetchAndPlay() {
        try {
            const playlistJSON = await window.go.main.App.FetchPlaylist(
                this.config.screenId,
                3
            );
            
            this.playlist = JSON.parse(playlistJSON);
            this.currentIndex = 0;
            
            if (this.playlist.items.length === 0) {
                this.showError('Playlist is empty');
                return;
            }
            
            this.playNext();
            
            setInterval(() => {
                this.fetchAndPlay();
            }, this.config.playlistRefreshInterval * 60 * 1000);
            
        } catch (error) {
            console.error('Playlist fetch error:', error);
            this.showError('Failed to fetch playlist: ' + error);
        }
    }

    async playNext() {
        if (!this.playlist || this.playlist.items.length === 0) {
            return;
        }

        const item = this.playlist.items[this.currentIndex];
        
        try {
            const impressionEvent = {
                screenId: this.config.screenId,
                creativeId: item.creativeId,
                campaignId: item.campaignId,
                timestamp: new Date().toISOString(),
            };
            await window.go.main.App.RecordImpression(JSON.stringify(impressionEvent));
        } catch (error) {
            console.error('Impression event error:', error);
        }

        this.player.display(item);

        setTimeout(async () => {
            try {
                const playEvent = {
                    screenId: this.config.screenId,
                    creativeId: item.creativeId,
                    campaignId: item.campaignId,
                    timestamp: new Date().toISOString(),
                };
                await window.go.main.App.RecordPlay(JSON.stringify(playEvent));
            } catch (error) {
                console.error('Play event error:', error);
            }

            this.currentIndex = (this.currentIndex + 1) % this.playlist.items.length;
            this.playNext();
        }, item.durationSeconds * 1000);
    }

    updateStatusBar() {
        this.statusBar.setConnected(true);
    }

    showError(message) {
        document.getElementById('player-container').innerHTML = 
            `<div class="error">${message}</div>`;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const app = new DisplayClient();
    app.init();
});
```

### 5.3 Components

**`frontend/components/Player.js`**
```javascript
export class Player {
    constructor(container) {
        this.container = container;
        this.currentVideo = null;
    }

    display(item) {
        this.container.innerHTML = '';

        if (item.type === 'image') {
            const img = document.createElement('img');
            img.src = item.url;
            img.className = 'media-content';
            img.onerror = () => this.showError('Failed to load image');
            this.container.appendChild(img);
        } else if (item.type === 'video') {
            const video = document.createElement('video');
            video.src = item.url;
            video.className = 'media-content';
            video.autoplay = true;
            video.muted = true;
            video.playsInline = true;
            video.onerror = () => this.showError('Failed to load video');
            this.container.appendChild(video);
            this.currentVideo = video;
        }
    }

    showError(message) {
        this.container.innerHTML = `<div class="error">${message}</div>`;
    }
}
```

**`frontend/components/StatusBar.js`**
```javascript
export class StatusBar {
    constructor(container) {
        this.container = container;
        this.connected = false;
        this.render();
    }

    setConnected(connected) {
        this.connected = connected;
        this.render();
    }

    render() {
        this.container.innerHTML = `
            <div class="status-bar">
                <span class="status-indicator ${this.connected ? 'connected' : 'disconnected'}"></span>
                <span>${this.connected ? 'Connected' : 'Disconnected'}</span>
            </div>
        `;
    }
}
```

### 5.4 Styles

**`frontend/styles.css`**
```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    width: 100vw;
    height: 100vh;
    overflow: hidden;
    background: #000;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

#app {
    width: 100%;
    height: 100%;
    position: relative;
}

#status-bar {
    position: absolute;
    top: 10px;
    right: 10px;
    z-index: 1000;
}

.status-bar {
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 8px 12px;
    border-radius: 4px;
    font-size: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
}

.status-indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    display: inline-block;
}

.status-indicator.connected {
    background: #10b981;
}

.status-indicator.disconnected {
    background: #ef4444;
}

#player-container {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.media-content {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
}

#loading, .error {
    color: white;
    font-size: 24px;
    text-align: center;
}

.error {
    color: #ef4444;
}
```

---

## ⚙️ Step 6: Configuration

**`wails.json`**
```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "mnemoCast-client",
  "outputfilename": "mnemoCast-client",
  "frontend": {
    "dir": "./frontend",
    "install": "",
    "build": "",
    "dev": "",
    "bridge": "",
    "serve": ""
  },
  "author": {
    "name": "Your Name",
    "email": "your.email@example.com"
  },
  "info": {
    "companyName": "MnemoCast",
    "productName": "MnemoCast Display Client",
    "productVersion": "0.1.0",
    "copyright": "Copyright...",
    "comments": "Digital Display Client for MnemoCast"
  }
}
```

---

## 🚀 Step 7: Run & Test

### Development Mode

```bash
# Start backend first (in another terminal)
cd ../backend
sbt run

# Run client in dev mode
wails dev
```

### Build for Production

```bash
# Build for current platform
wails build

# Build for specific platform
wails build -platform linux/amd64
wails build -platform windows/amd64
wails build -platform darwin/amd64
```

---

## ✅ Next Steps

1. ✅ **Test basic functionality** - Register screen, fetch playlist
2. ✅ **Add caching** - Implement filesystem cache (see `OFFLINE_CACHING_STRATEGY.md`)
3. ✅ **Add error handling** - Retry logic, offline detection
4. ✅ **Add preloading** - Preload next items for smooth playback
5. ✅ **Add event queue** - Queue events when offline
6. ✅ **Polish UI** - Better transitions, loading states

---

## 📚 Reference Documents

- **Full Implementation Plan:** `GOLANG_WAILS_IMPLEMENTATION_PLAN.md`
- **Offline Caching:** `OFFLINE_CACHING_STRATEGY.md`
- **Storage Analysis:** `CACHE_STORAGE_ANALYSIS.md`
- **API Endpoints:** See backend documentation

---

## 🐛 Troubleshooting

### Issue: `wails dev` fails
- **Check:** Node.js is installed (`node --version`)
- **Check:** Frontend directory exists
- **Fix:** Run `wails init` again if needed

### Issue: Cannot connect to backend
- **Check:** Backend is running on port 8080
- **Check:** CORS is enabled in backend
- **Check:** BackendURL in config is correct

### Issue: Build fails
- **Check:** Go version is 1.21+
- **Check:** All dependencies installed (`go mod tidy`)
- **Check:** Platform-specific build tools installed

---

**Status:** Ready to Start  
**Estimated Setup Time:** 30-60 minutes

