# 💾 Offline Caching Strategy for Digital Display Client

**Comprehensive plan for caching ads and playlists when internet connection is unavailable.**

---

## 🎯 Overview

The client must continue playing content seamlessly even when offline. This document outlines a multi-layered caching strategy for playlists, media files, and event queuing.

---

## 📋 Caching Requirements

### What to Cache

1. **Playlist Metadata**
   - Playlist JSON structure
   - Creative metadata (IDs, URLs, durations)
   - Timestamp of last successful fetch

2. **Media Files**
   - Images (JPG, PNG, WebP)
   - Videos (MP4, WebM)
   - Preload next items for smooth playback

3. **Events Queue**
   - Impression events
   - Play events
   - Heartbeat status
   - Sync when connection restored

4. **Configuration**
   - Screen configuration
   - Last known good state

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Application                         │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Playback Engine                          │   │
│  │  - Current playlist                                  │   │
│  │  - Media player                                      │   │
│  └──────────────────┬───────────────────────────────────┘   │
│                     │                                         │
│  ┌──────────────────▼───────────────────────────────────┐   │
│  │         Cache Manager (Multi-Layer)                  │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────┐  │   │
│  │  │ Playlist     │  │ Media        │  │ Event    │  │   │
│  │  │ Cache        │  │ Cache        │  │ Queue    │  │   │
│  │  └──────────────┘  └──────────────┘  └──────────┘  │   │
│  └──────────────────┬───────────────────────────────────┘   │
│                     │                                         │
│  ┌──────────────────▼───────────────────────────────────┐   │
│  │         Storage Layer                                 │   │
│  │  - File System (media files + JSON files)            │   │
│  │  - No database needed (filesystem sufficient)        │   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

---

## 💾 Storage Strategy

### Layer 1: Playlist Cache

**Purpose:** Store playlist metadata for quick access

**Storage Options:**

#### For Web/Electron (TypeScript):
- **IndexedDB** (preferred) - Structured storage, larger capacity
- **localStorage** (fallback) - Simple, limited to ~5-10MB

#### For Rust (Tauri):
- **JSON files** (recommended) - Simple, easy to read/write, sufficient for use case
- **SQLite** (optional) - Only if complex queries needed

#### For Golang (Wails):
- **JSON files** (recommended) - Simple, sufficient for use case
- **SQLite/BoltDB** (optional) - Only if complex queries needed

**Recommendation:** Use JSON files. Filesystem is sufficient for this use case. See `CACHE_STORAGE_ANALYSIS.md` for detailed analysis.

**Cache Structure:**
```json
{
  "playlist": {
    "screenId": "screen-1",
    "generatedAt": "2024-12-18T15:00:00Z",
    "durationMinutes": 3,
    "items": [...]
  },
  "metadata": {
    "cachedAt": "2024-12-18T15:00:00Z",
    "expiresAt": "2024-12-18T15:03:00Z",
    "version": 1
  }
}
```

### Layer 2: Media Cache

**Purpose:** Store actual image/video files locally

**Storage:**
- **File System** - Dedicated cache directory
- **Naming Convention:** `{creativeId}_{hash}.{ext}`
- **Size Limits:** Configurable (default: 500MB-2GB)

**Cache Directory Structure:**
```
cache/
├── media/
│   ├── images/
│   │   ├── creative-1_a1b2c3.jpg
│   │   └── creative-2_d4e5f6.png
│   └── videos/
│       ├── creative-3_g7h8i9.mp4
│       └── creative-4_j0k1l2.webm
├── playlists/
│   ├── playlist_20241218_150000.json
│   └── playlist_20241218_150300.json
└── events/
    └── events_queue.json
```

### Layer 3: Event Queue

**Purpose:** Store events when offline, sync when online

**Storage:**
- **In-Memory Queue** (primary) - Fast access
- **Persistent Queue** (backup) - Survive restarts
- **Max Queue Size:** 1000 events (configurable)

**Queue Structure:**
```json
{
  "events": [
    {
      "type": "impression",
      "data": {
        "screenId": "screen-1",
        "creativeId": "creative-1",
        "campaignId": "campaign-1",
        "timestamp": "2024-12-18T15:00:00Z"
      },
      "queuedAt": "2024-12-18T15:00:05Z",
      "retries": 0
    }
  ],
  "metadata": {
    "lastSyncAttempt": "2024-12-18T15:01:00Z",
    "totalQueued": 15
  }
}
```

---

## 🔄 Caching Workflow

### Online Mode

```
1. Fetch playlist from backend
   ↓
2. Save playlist to cache (with timestamp)
   ↓
3. For each creative in playlist:
   ├─ Check if media exists in cache
   ├─ If not cached: Download and save
   └─ Preload next 2-3 items
   ↓
4. Play from cache (faster, offline-capable)
   ↓
5. Send events immediately
   ↓
6. Refresh playlist every N minutes
```

### Offline Mode Detection

```javascript
// Pseudo-code
function isOnline() {
    return navigator.onLine && canReachBackend();
}

async function canReachBackend() {
    try {
        const response = await fetch('/api/v1/health', {
            timeout: 5000
        });
        return response.ok;
    } catch {
        return false;
    }
}
```

### Offline Mode Workflow

```
1. Detect offline status
   ↓
2. Load last cached playlist
   ↓
3. Check cache expiry
   ├─ If expired: Show "Offline - Using cached content"
   └─ If valid: Continue normally
   ↓
4. Play from cache
   ├─ Check if media file exists locally
   ├─ If missing: Skip item or show placeholder
   └─ Continue to next item
   ↓
5. Queue all events (impression, play)
   ↓
6. Periodically retry connection
   ├─ When online: Sync event queue
   └─ Fetch fresh playlist
```

---

## 🛠️ Implementation Details

### Rust (Tauri) Implementation

#### Playlist Cache

**`src-tauri/src/storage/cache.rs`**
```rust
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CachedPlaylist {
    pub playlist: Playlist,
    pub cached_at: DateTime<Utc>,
    pub expires_at: DateTime<Utc>,
}

pub struct PlaylistCache {
    cache_dir: PathBuf,
}

impl PlaylistCache {
    pub fn new(cache_dir: PathBuf) -> Self {
        fs::create_dir_all(&cache_dir).ok();
        Self { cache_dir }
    }

    pub fn save(&self, playlist: &Playlist, ttl_minutes: i64) -> Result<(), Box<dyn std::error::Error>> {
        let cached = CachedPlaylist {
            playlist: playlist.clone(),
            cached_at: Utc::now(),
            expires_at: Utc::now() + chrono::Duration::minutes(ttl_minutes),
        };

        let file_path = self.cache_dir.join("playlist.json");
        let json = serde_json::to_string_pretty(&cached)?;
        fs::write(file_path, json)?;
        
        Ok(())
    }

    pub fn load(&self) -> Option<CachedPlaylist> {
        let file_path = self.cache_dir.join("playlist.json");
        
        if !file_path.exists() {
            return None;
        }

        match fs::read_to_string(&file_path) {
            Ok(content) => {
                match serde_json::from_str::<CachedPlaylist>(&content) {
                    Ok(cached) => {
                        // Check if expired
                        if Utc::now() > cached.expires_at {
                            None
                        } else {
                            Some(cached)
                        }
                    }
                    Err(_) => None,
                }
            }
            Err(_) => None,
        }
    }

    pub fn is_expired(&self) -> bool {
        match self.load() {
            Some(cached) => Utc::now() > cached.expires_at,
            None => true,
        }
    }
}
```

#### Media Cache

**`src-tauri/src/storage/media_cache.rs`**
```rust
use std::fs;
use std::path::PathBuf;
use std::collections::HashMap;
use reqwest::Client;
use sha2::{Sha256, Digest};

pub struct MediaCache {
    cache_dir: PathBuf,
    max_size_mb: u64,
}

impl MediaCache {
    pub fn new(cache_dir: PathBuf, max_size_mb: u64) -> Self {
        fs::create_dir_all(&cache_dir.join("images")).ok();
        fs::create_dir_all(&cache_dir.join("videos")).ok();
        Self { cache_dir, max_size_mb }
    }

    pub async fn get_or_download(&self, url: &str, creative_id: &str) -> Result<PathBuf, Box<dyn std::error::Error>> {
        // Generate cache key
        let cache_key = self.generate_cache_key(url);
        let ext = self.get_extension(url);
        let media_type = if url.contains(".mp4") || url.contains(".webm") {
            "videos"
        } else {
            "images"
        };
        
        let file_path = self.cache_dir.join(media_type).join(format!("{}_{}.{}", creative_id, cache_key, ext));

        // Check if already cached
        if file_path.exists() {
            return Ok(file_path);
        }

        // Download and cache
        self.download_and_cache(url, &file_path).await?;
        
        // Check cache size and cleanup if needed
        self.cleanup_if_needed().await?;

        Ok(file_path)
    }

    async fn download_and_cache(&self, url: &str, file_path: &PathBuf) -> Result<(), Box<dyn std::error::Error>> {
        let client = Client::new();
        let response = client.get(url).send().await?;
        let bytes = response.bytes().await?;
        fs::write(file_path, bytes)?;
        Ok(())
    }

    fn generate_cache_key(&self, url: &str) -> String {
        let mut hasher = Sha256::new();
        hasher.update(url.as_bytes());
        format!("{:x}", hasher.finalize())[..8].to_string()
    }

    fn get_extension(&self, url: &str) -> &str {
        url.split('.').last().unwrap_or("bin")
    }

    async fn cleanup_if_needed(&self) -> Result<(), Box<dyn std::error::Error>> {
        let total_size = self.get_cache_size().await?;
        let max_size_bytes = self.max_size_mb * 1024 * 1024;

        if total_size > max_size_bytes {
            // Remove oldest files (LRU strategy)
            self.remove_oldest_files(total_size - max_size_bytes).await?;
        }

        Ok(())
    }

    async fn get_cache_size(&self) -> Result<u64, Box<dyn std::error::Error>> {
        let mut total = 0u64;
        for entry in fs::read_dir(&self.cache_dir)? {
            let entry = entry?;
            if entry.file_type()?.is_file() {
                total += entry.metadata()?.len();
            }
        }
        Ok(total)
    }

    async fn remove_oldest_files(&self, bytes_to_free: u64) -> Result<(), Box<dyn std::error::Error>> {
        // Implementation: Sort files by modification time, remove oldest
        // until bytes_to_free is freed
        Ok(())
    }
}
```

#### Event Queue

**`src-tauri/src/storage/event_queue.rs`**
```rust
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;
use std::sync::Mutex;
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QueuedEvent {
    pub event_type: String, // "impression" or "play"
    pub data: serde_json::Value,
    pub queued_at: DateTime<Utc>,
    pub retries: u32,
}

pub struct EventQueue {
    queue: Mutex<Vec<QueuedEvent>>,
    queue_file: PathBuf,
    max_size: usize,
}

impl EventQueue {
    pub fn new(queue_file: PathBuf, max_size: usize) -> Self {
        let queue = Self::load_from_file(&queue_file).unwrap_or_default();
        Self {
            queue: Mutex::new(queue),
            queue_file,
            max_size,
        }
    }

    pub fn enqueue(&self, event_type: String, data: serde_json::Value) -> Result<(), Box<dyn std::error::Error>> {
        let mut queue = self.queue.lock().unwrap();
        
        if queue.len() >= self.max_size {
            return Err("Event queue is full".into());
        }

        let queued_event = QueuedEvent {
            event_type,
            data,
            queued_at: Utc::now(),
            retries: 0,
        };

        queue.push(queued_event);
        self.save_to_file(&queue)?;
        Ok(())
    }

    pub fn dequeue_all(&self) -> Vec<QueuedEvent> {
        let mut queue = self.queue.lock().unwrap();
        let events = queue.drain(..).collect();
        self.save_to_file(&queue).ok();
        events
    }

    pub fn sync_events(&self, api_client: &api::EventsAPI) -> Result<usize, Box<dyn std::error::Error>> {
        let events = self.dequeue_all();
        let mut synced = 0;

        for event in events {
            match self.send_event(&event, api_client).await {
                Ok(_) => synced += 1,
                Err(e) => {
                    // Re-queue with retry count
                    if event.retries < 3 {
                        let mut queue = self.queue.lock().unwrap();
                        let mut retry_event = event;
                        retry_event.retries += 1;
                        queue.push(retry_event);
                    }
                    // Otherwise drop the event after max retries
                }
            }
        }

        Ok(synced)
    }

    async fn send_event(&self, event: &QueuedEvent, api_client: &api::EventsAPI) -> Result<(), Box<dyn std::error::Error>> {
        match event.event_type.as_str() {
            "impression" => {
                let play_event: models::PlayEvent = serde_json::from_value(event.data.clone())?;
                api_client.record_impression(play_event).await?;
            }
            "play" => {
                let play_event: models::PlayEvent = serde_json::from_value(event.data.clone())?;
                api_client.record_play(play_event).await?;
            }
            _ => return Err("Unknown event type".into()),
        }
        Ok(())
    }

    fn load_from_file(file_path: &PathBuf) -> Result<Vec<QueuedEvent>, Box<dyn std::error::Error>> {
        if !file_path.exists() {
            return Ok(Vec::new());
        }
        let content = fs::read_to_string(file_path)?;
        let events: Vec<QueuedEvent> = serde_json::from_str(&content)?;
        Ok(events)
    }

    fn save_to_file(&self, queue: &Vec<QueuedEvent>) -> Result<(), Box<dyn std::error::Error>> {
        let json = serde_json::to_string_pretty(queue)?;
        fs::write(&self.queue_file, json)?;
        Ok(())
    }
}
```

---

### Golang (Wails) Implementation

#### Playlist Cache

**`app/internal/storage/cache.go`**
```go
package storage

import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"
    "mnemoCast-client/app/internal/models"
)

type CachedPlaylist struct {
    Playlist   models.Playlist `json:"playlist"`
    CachedAt   time.Time       `json:"cachedAt"`
    ExpiresAt  time.Time       `json:"expiresAt"`
}

type PlaylistCache struct {
    cacheDir string
}

func NewPlaylistCache(cacheDir string) *PlaylistCache {
    os.MkdirAll(cacheDir, 0755)
    return &PlaylistCache{cacheDir: cacheDir}
}

func (c *PlaylistCache) Save(playlist *models.Playlist, ttlMinutes int) error {
    cached := CachedPlaylist{
        Playlist:  *playlist,
        CachedAt:  time.Now(),
        ExpiresAt: time.Now().Add(time.Duration(ttlMinutes) * time.Minute),
    }

    filePath := filepath.Join(c.cacheDir, "playlist.json")
    data, err := json.MarshalIndent(cached, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(filePath, data, 0644)
}

func (c *PlaylistCache) Load() (*CachedPlaylist, error) {
    filePath := filepath.Join(c.cacheDir, "playlist.json")
    
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return nil, nil
    }

    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var cached CachedPlaylist
    if err := json.Unmarshal(data, &cached); err != nil {
        return nil, err
    }

    // Check if expired
    if time.Now().After(cached.ExpiresAt) {
        return nil, nil
    }

    return &cached, nil
}

func (c *PlaylistCache) IsExpired() bool {
    cached, err := c.Load()
    if err != nil || cached == nil {
        return true
    }
    return time.Now().After(cached.ExpiresAt)
}
```

#### Media Cache

**`app/internal/storage/media_cache.go`**
```go
package storage

import (
    "crypto/sha256"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

type MediaCache struct {
    cacheDir  string
    maxSizeMB int64
}

func NewMediaCache(cacheDir string, maxSizeMB int64) *MediaCache {
    os.MkdirAll(filepath.Join(cacheDir, "images"), 0755)
    os.MkdirAll(filepath.Join(cacheDir, "videos"), 0755)
    return &MediaCache{
        cacheDir:  cacheDir,
        maxSizeMB: maxSizeMB,
    }
}

func (m *MediaCache) GetOrDownload(url, creativeID string) (string, error) {
    cacheKey := m.generateCacheKey(url)
    ext := m.getExtension(url)
    
    var mediaType string
    if strings.Contains(url, ".mp4") || strings.Contains(url, ".webm") {
        mediaType = "videos"
    } else {
        mediaType = "images"
    }

    filePath := filepath.Join(m.cacheDir, mediaType, fmt.Sprintf("%s_%s.%s", creativeID, cacheKey, ext))

    // Check if already cached
    if _, err := os.Stat(filePath); err == nil {
        return filePath, nil
    }

    // Download and cache
    if err := m.downloadAndCache(url, filePath); err != nil {
        return "", err
    }

    // Check cache size and cleanup if needed
    m.cleanupIfNeeded()

    return filePath, nil
}

func (m *MediaCache) downloadAndCache(url, filePath string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}

func (m *MediaCache) generateCacheKey(url string) string {
    h := sha256.New()
    h.Write([]byte(url))
    return fmt.Sprintf("%x", h.Sum(nil))[:8]
}

func (m *MediaCache) getExtension(url string) string {
    parts := strings.Split(url, ".")
    if len(parts) > 0 {
        return parts[len(parts)-1]
    }
    return "bin"
}

func (m *MediaCache) cleanupIfNeeded() {
    // Implementation: Check cache size, remove oldest files if over limit
}
```

#### Event Queue

**`app/internal/storage/event_queue.go`**
```go
package storage

import (
    "encoding/json"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type QueuedEvent struct {
    EventType string          `json:"eventType"`
    Data      json.RawMessage `json:"data"`
    QueuedAt  time.Time       `json:"queuedAt"`
    Retries   int             `json:"retries"`
}

type EventQueue struct {
    queue     []QueuedEvent
    mutex     sync.Mutex
    queueFile string
    maxSize   int
}

func NewEventQueue(queueFile string, maxSize int) *EventQueue {
    eq := &EventQueue{
        queue:     make([]QueuedEvent, 0),
        queueFile: queueFile,
        maxSize:   maxSize,
    }
    eq.loadFromFile()
    return eq
}

func (eq *EventQueue) Enqueue(eventType string, data json.RawMessage) error {
    eq.mutex.Lock()
    defer eq.mutex.Unlock()

    if len(eq.queue) >= eq.maxSize {
        return fmt.Errorf("event queue is full")
    }

    event := QueuedEvent{
        EventType: eventType,
        Data:      data,
        QueuedAt:  time.Now(),
        Retries:   0,
    }

    eq.queue = append(eq.queue, event)
    return eq.saveToFile()
}

func (eq *EventQueue) DequeueAll() []QueuedEvent {
    eq.mutex.Lock()
    defer eq.mutex.Unlock()

    events := eq.queue
    eq.queue = make([]QueuedEvent, 0)
    eq.saveToFile()
    return events
}

func (eq *EventQueue) SyncEvents(apiClient interface{}) (int, error) {
    events := eq.DequeueAll()
    synced := 0

    for _, event := range events {
        if err := eq.sendEvent(event, apiClient); err != nil {
            // Re-queue with retry count
            if event.Retries < 3 {
                event.Retries++
                eq.Enqueue(event.EventType, event.Data)
            }
            continue
        }
        synced++
    }

    return synced, nil
}

func (eq *EventQueue) sendEvent(event QueuedEvent, apiClient interface{}) error {
    // Implementation: Send event via API client
    return nil
}

func (eq *EventQueue) loadFromFile() {
    if _, err := os.Stat(eq.queueFile); os.IsNotExist(err) {
        return
    }

    data, err := os.ReadFile(eq.queueFile)
    if err != nil {
        return
    }

    json.Unmarshal(data, &eq.queue)
}

func (eq *EventQueue) saveToFile() error {
    data, err := json.MarshalIndent(eq.queue, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(eq.queueFile, data, 0644)
}
```

---

## 🔄 Cache Management Strategies

### 1. Cache Invalidation

**Time-Based (TTL):**
- Playlist expires after refresh interval (default: 3 minutes)
- Media files expire after 7 days of non-use
- Events expire after 30 days (should never happen)

**Size-Based (LRU):**
- When cache exceeds max size, remove least recently used files
- Keep at least 2 playlists worth of media

**Version-Based:**
- If playlist version changes, invalidate old media
- Download new creatives automatically

### 2. Preloading Strategy

```javascript
// Preload next 2-3 items while current item plays
async function preloadNextItems(playlist, currentIndex) {
    const preloadCount = 3;
    for (let i = 1; i <= preloadCount; i++) {
        const nextIndex = (currentIndex + i) % playlist.items.length;
        const item = playlist.items[nextIndex];
        
        // Download media if not cached
        await mediaCache.getOrDownload(item.url, item.creativeId);
    }
}
```

### 3. Cache Size Limits

**Recommended Limits:**
- **Playlist Cache:** 10 playlists max (~1MB)
- **Media Cache:** 500MB - 2GB (configurable)
- **Event Queue:** 1000 events max (~500KB)

**Cleanup Rules:**
- Remove media files older than 7 days
- Remove playlists older than 1 hour
- Remove events after successful sync

---

## 🚨 Offline Mode Indicators

### UI Indicators

```css
.offline-indicator {
    position: fixed;
    top: 10px;
    left: 10px;
    background: rgba(255, 193, 7, 0.9);
    color: #000;
    padding: 8px 12px;
    border-radius: 4px;
    font-size: 12px;
    z-index: 1000;
}

.offline-indicator::before {
    content: "⚠️ ";
}
```

### Status Messages

- **"Offline - Using cached content"** - Playing from cache
- **"Reconnecting..."** - Attempting to reconnect
- **"Syncing events..."** - Uploading queued events
- **"Cache expired"** - Cached content is too old

---

## 📊 Cache Statistics

Track and display:
- Cache hit rate
- Cache size (current/max)
- Number of queued events
- Last successful sync time
- Offline duration

---

## 🧪 Testing Offline Mode

### Test Scenarios

1. **Start offline:** Launch app without internet
2. **Go offline during playback:** Disconnect mid-session
3. **Cache expiry:** Wait for cache to expire
4. **Event queue overflow:** Generate >1000 events offline
5. **Cache size limit:** Fill cache to capacity
6. **Reconnection:** Reconnect after extended offline period

### Test Checklist

- [ ] Playlist loads from cache when offline
- [ ] Media files play from cache
- [ ] Events are queued when offline
- [ ] Events sync when connection restored
- [ ] Cache cleanup works correctly
- [ ] Preloading works as expected
- [ ] UI shows offline status
- [ ] App handles cache expiry gracefully

---

## ⚙️ Configuration

```json
{
  "cache": {
    "enabled": true,
    "playlistTTLMinutes": 3,
    "mediaTTLDays": 7,
    "maxMediaCacheSizeMB": 1000,
    "maxPlaylistCacheCount": 10,
    "maxEventQueueSize": 1000,
    "preloadCount": 3,
    "cleanupIntervalMinutes": 60
  }
}
```

---

## 🚀 Implementation Priority

### Phase 1: Basic Caching (MVP)
1. ✅ Playlist cache (JSON file)
2. ✅ Media cache (file system)
3. ✅ Basic offline detection
4. ✅ Fallback to cached playlist

### Phase 2: Event Queue
1. ✅ Queue events when offline
2. ✅ Sync events when online
3. ✅ Retry logic for failed syncs

### Phase 3: Advanced Features
1. ✅ Preloading
2. ✅ Cache cleanup
3. ✅ Cache statistics
4. ✅ UI indicators

---

**Status:** Ready for Implementation  
**Estimated Time:** 3-5 days for complete offline caching system

