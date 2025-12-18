# 🐹 Digital Display Client - Golang + Wails Implementation Plan

**Version:** 1.0.0  
**Last Updated:** December 2024  
**Status:** Ready for Implementation

---

## 🎯 Overview

This document provides a comprehensive implementation plan for building the Digital Display Client using **Golang** (backend) and **Wails** (frontend framework) to create a lightweight, efficient native application.

### Why Golang + Wails?

- ✅ **Lightweight:** ~15-25MB binary (vs 100-200MB for Electron)
- ✅ **Fast development:** Simple syntax, quick iteration
- ✅ **Excellent concurrency:** Goroutines for async operations
- ✅ **Memory efficient:** ~20-50MB RAM usage
- ✅ **Cross-platform:** Linux, Windows, macOS
- ✅ **Easy deployment:** Single binary, no dependencies

---

## 🏗️ Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                    Wails Application                         │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐    │
│  │              Frontend (HTML/CSS/JS)                  │    │
│  │  - Vanilla JS or React/Vue                           │    │
│  │  - HTML5 Video/Image for media                       │    │
│  │  - UI Components                                     │    │
│  └──────────────────┬───────────────────────────────────┘    │
│                     │ Context (Go <-> JS)                    │
│  ┌──────────────────▼───────────────────────────────────┐    │
│  │              Go Backend                              │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────┐    │    │
│  │  │ API Client   │  │ Event        │  │ Config   │    │    │
│  │  │ (net/http)   │  │ Tracker      │  │ Manager  │    │    │
│  │  └──────────────┘  └──────────────┘  └──────────┘    │    │
│  │  ┌──────────────┐  ┌──────────────┐                  │    │
│  │  │ Playlist     │  │ Heartbeat    │                  │    │
│  │  │ Manager      │  │ Scheduler    │                  │    │
│  │  └──────────────┘  └──────────────┘                  │    │
│  └──────────────────────────────────────────────────────┘    │
└───────────────────────────┬──────────────────────────────────┘
                            │ HTTP/REST API
┌───────────────────────────▼──────────────────────────────────┐
│              Scala/Pekko Backend                             │
│         (http://localhost:8080/api/v1/...)                   │
└──────────────────────────────────────────────────────────────┘
```

---

## 🛠️ Technology Stack

### Backend (Golang)
- **Language:** Go 1.21+
- **HTTP Client:** `net/http` (standard library) or `resty`
- **JSON:** `encoding/json` (standard library)
- **Concurrency:** Goroutines + Channels
- **Error Handling:** `errors`, `fmt`
- **Logging:** `log` or `logrus`
- **Config:** `viper` or `config`

### Frontend (Wails)
- **Framework:** Wails v2 (latest)
- **UI:** Vanilla JS/HTML/CSS or React/Vue
- **Media:** HTML5 `<video>` and `<img>` elements
- **Styling:** Tailwind CSS (optional)

### Build Tools
- **Package Manager:** Go Modules
- **Build System:** Wails CLI
- **Target Platforms:** x86_64, ARM64 (Linux, Windows, macOS)

---

## 📦 Project Structure

```
mnemoCast-client-go/
├── app/                          # Go backend
│   ├── main.go                   # Wails entry point
│   ├── internal/
│   │   ├── api/                  # API client
│   │   │   ├── client.go         # HTTP client wrapper
│   │   │   ├── screen.go         # Screen API
│   │   │   ├── playlist.go       # Playlist API
│   │   │   └── events.go         # Events API
│   │   ├── models/               # Data models
│   │   │   ├── screen.go
│   │   │   ├── playlist.go
│   │   │   └── events.go
│   │   ├── scheduler/            # Background tasks
│   │   │   ├── heartbeat.go     # Heartbeat scheduler
│   │   │   └── refresh.go        # Playlist refresh
│   │   ├── storage/              # Local storage
│   │   │   ├── cache.go          # Playlist cache
│   │   │   └── config.go         # Config storage
│   │   └── app.go                # App struct
│   └── go.mod
├── frontend/                      # Frontend (web)
│   ├── index.html
│   ├── main.js                   # Entry point
│   ├── styles.css
│   ├── components/
│   │   ├── Player.js
│   │   └── StatusBar.js
│   └── app.js
├── build/                         # Build artifacts
├── wails.json                     # Wails configuration
├── go.mod
└── README.md
```

---

## 🚀 Implementation Steps

### Phase 1: Project Setup (Day 1-2)

#### Step 1.1: Install Prerequisites

```bash
# Install Go 1.21+
# Download from https://go.dev/dl/

# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify installation
wails version
```

#### Step 1.2: Create Wails Project

```bash
# Create new Wails project
wails init -n mnemoCast-client -t vanilla

# Or with React/Vue template
# wails init -n mnemoCast-client -t react
# wails init -n mnemoCast-client -t vue
```

#### Step 1.3: Configure go.mod

```bash
cd mnemoCast-client
go mod init mnemoCast-client
go mod tidy
```

**`go.mod`**
```go
module mnemoCast-client

go 1.21

require (
    github.com/wailsapp/wails/v2 v2.8.0
    github.com/go-resty/resty/v2 v2.11.0
    github.com/spf13/viper v1.18.2
    github.com/sirupsen/logrus v1.9.3
)
```

---

### Phase 2: Core Backend Implementation (Day 3-5)

#### Step 2.1: Define Data Models

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

#### Step 2.2: Create API Client

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

#### Step 2.3: Create App Struct with Wails Methods

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
    PlaylistRefreshInterval int                 `json:"playlistRefreshInterval"` // minutes
    HeartbeatInterval       int                 `json:"heartbeatInterval"`       // seconds
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

// Wails lifecycle methods

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
            // Log error, but don't block
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

// Helper functions

func toJSON(v interface{}) (string, error) {
    data, err := json.Marshal(v)
    if err != nil {
        return "", err
    }
    return string(data), nil
}
```

#### Step 2.4: Main Entry Point

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

### Phase 3: Frontend Implementation (Day 6-8)

#### Step 3.1: Basic HTML Structure

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

#### Step 3.2: Main Application Logic

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
            // Load configuration
            const configJSON = await window.go.main.App.GetConfig();
            this.config = JSON.parse(configJSON);
            
            // Register screen
            await this.registerScreen();
            
            // Fetch and play playlist
            await this.fetchAndPlay();
            
            // Status bar will be updated by heartbeat
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
                3 // durationMinutes
            );
            
            this.playlist = JSON.parse(playlistJSON);
            this.currentIndex = 0;
            
            if (this.playlist.items.length === 0) {
                this.showError('Playlist is empty');
                return;
            }
            
            this.playNext();
            
            // Refresh playlist periodically
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
        
        // Send impression event
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

        // Display media
        this.player.display(item);

        // Schedule next item
        setTimeout(async () => {
            // Send play event
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

            // Move to next item
            this.currentIndex = (this.currentIndex + 1) % this.playlist.items.length;
            this.playNext();
        }, item.durationSeconds * 1000);
    }

    updateStatusBar() {
        // Status bar updates based on connection
        // In a real implementation, you'd check connection status
        this.statusBar.setConnected(true);
    }

    showError(message) {
        document.getElementById('player-container').innerHTML = 
            `<div class="error">${message}</div>`;
    }
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    const app = new DisplayClient();
    app.init();
});
```

#### Step 3.3: Player Component

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

#### Step 3.4: Status Bar Component

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

#### Step 3.5: Styles

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

### Phase 4: Configuration & Build (Day 9-10)

#### Step 4.1: Wails Configuration

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
  },
  "nsisType": "multiple",
  "obfuscated": false,
  "garbleargs": ""
}
```

#### Step 4.2: Build Commands

```bash
# Development
wails dev

# Production build
wails build

# Build for specific platform
wails build -platform windows/amd64
wails build -platform linux/amd64
wails build -platform darwin/amd64
wails build -platform darwin/arm64
```

---

## 🧪 Testing

### Unit Tests

**`app/internal/api/screen_test.go`**
```go
package api_test

import (
    "mnemoCast-client/app/internal/api"
    "mnemoCast-client/app/internal/models"
    "testing"
)

func TestScreenRegistration(t *testing.T) {
    // Test implementation
    screenAPI := api.NewScreenAPI("http://localhost:8080")
    
    request := models.RegisterScreenRequest{
        ID:            "test-screen",
        Name:          "Test Screen",
        Location:      models.ScreenLocation{City: "Test", Area: "Test", VenueType: "test"},
        Classification: 1,
    }
    
    screen, err := screenAPI.Register(request)
    if err != nil {
        t.Fatalf("Registration failed: %v", err)
    }
    
    if screen.ID != "test-screen" {
        t.Errorf("Expected screen ID 'test-screen', got '%s'", screen.ID)
    }
}
```

### Integration Tests

```bash
# Start backend first
cd ../backend
sbt run

# In another terminal, run client
wails dev
```

---

## 📊 Performance Targets

- **Binary Size:** < 25MB
- **Memory Usage:** < 50MB
- **Startup Time:** < 300ms
- **API Latency:** < 100ms
- **Playlist Fetch:** < 500ms
- **Event Send:** < 50ms (async, non-blocking)

---

## 🚀 Deployment

### Linux
```bash
wails build -platform linux/amd64
# Output: build/bin/mnemoCast-client
```

### Windows
```bash
wails build -platform windows/amd64
# Output: build/bin/mnemoCast-client.exe
```

### macOS
```bash
wails build -platform darwin/amd64
# Output: build/bin/mnemoCast-client.app
```

### Creating Installers

Wails can generate installers using:
- **Windows:** NSIS installer
- **macOS:** DMG or PKG
- **Linux:** AppImage, DEB, or RPM

---

## 🔄 Async Operations

Golang's goroutines make async operations simple:

```go
// Example: Non-blocking event sending
go func() {
    if err := a.eventsAPI.RecordImpression(event); err != nil {
        log.Printf("Failed to send impression: %v", err)
    }
}()

// Example: Background playlist refresh
go func() {
    ticker := time.NewTicker(3 * time.Minute)
    for range ticker.C {
        playlist, err := a.playlistAPI.Fetch(screenID, 3)
        // Handle playlist update
    }
}()
```

---

## 📝 Next Steps

1. ✅ Set up project structure
2. ✅ Implement API client
3. ✅ Create Wails app methods
4. ✅ Build frontend UI
5. ✅ Add background schedulers
6. ✅ Test with Scala/Pekko backend
7. ✅ Optimize performance
8. ✅ Create installers

---

## 🆚 Comparison: Golang vs Rust

| Aspect | Golang | Rust |
|--------|--------|------|
| Development Speed | ⭐⭐⭐⭐⭐ Faster | ⭐⭐⭐ Slower (learning curve) |
| Binary Size | ~15-25MB | ~3-8MB |
| Memory Usage | ~20-50MB | ~10-30MB |
| Latency | Good (GC pauses possible) | Excellent (no GC) |
| Concurrency | ⭐⭐⭐⭐⭐ Excellent (goroutines) | ⭐⭐⭐⭐ Excellent (async) |
| Code Simplicity | ⭐⭐⭐⭐⭐ Very simple | ⭐⭐⭐ More complex |

---

**Status:** Ready for Implementation  
**Estimated Time:** 8-10 days for MVP

