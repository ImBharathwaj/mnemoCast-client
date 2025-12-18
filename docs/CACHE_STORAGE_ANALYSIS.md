# 💾 Cache Storage: Filesystem vs Database Analysis

**Analysis of whether a database is needed for caching or if filesystem is sufficient.**

---

## 🎯 Short Answer

**Filesystem is sufficient** for this use case. You don't need a database.

However, a lightweight embedded database (SQLite) can provide benefits for certain operations, but it's **optional**, not required.

---

## 📊 Requirements Analysis

### What We Need to Store:

1. **Playlists** (~10-50KB each)
   - Simple JSON structure
   - Read: Load current playlist
   - Write: Save new playlist
   - No complex queries needed

2. **Media Files** (images/videos)
   - Binary files (100KB - 50MB each)
   - Read: Load file by path
   - Write: Save downloaded file
   - Filesystem is perfect for this

3. **Event Queue** (~1-10KB total)
   - Array of events
   - Read: Get all queued events
   - Write: Append event, remove after sync
   - Simple FIFO queue

4. **Cache Metadata** (~1KB)
   - File paths, timestamps, sizes
   - Read: Check if cached, get file path
   - Write: Update metadata

---

## ✅ Filesystem-Only Approach

### Advantages

✅ **Simple** - No dependencies, no setup
✅ **Lightweight** - Zero overhead
✅ **Fast** - Direct file access
✅ **Portable** - Works everywhere
✅ **Easy to debug** - Can inspect files directly
✅ **No corruption risk** - Simple file operations
✅ **Sufficient for use case** - All operations are simple

### Implementation

```
cache/
├── playlists/
│   └── current.json              # Current playlist
├── media/
│   ├── images/
│   │   └── creative-1_hash.jpg  # Media files
│   └── videos/
│       └── creative-2_hash.mp4
├── events/
│   └── queue.json                # Event queue
└── metadata.json                  # Cache metadata
```

### Code Example (Simple)

**Rust:**
```rust
// Save playlist
fs::write("cache/playlists/current.json", json)?;

// Load playlist
let json = fs::read_to_string("cache/playlists/current.json")?;
let playlist: Playlist = serde_json::from_str(&json)?;

// Check if media cached
let path = format!("cache/media/images/{}_{}.jpg", creative_id, hash);
if Path::new(&path).exists() {
    // Use cached file
}
```

**Golang:**
```go
// Save playlist
os.WriteFile("cache/playlists/current.json", jsonData, 0644)

// Load playlist
data, _ := os.ReadFile("cache/playlists/current.json")
json.Unmarshal(data, &playlist)

// Check if media cached
if _, err := os.Stat(filePath); err == nil {
    // Use cached file
}
```

### Performance

- **Read playlist:** < 1ms (small JSON file)
- **Check if media cached:** < 1ms (file exists check)
- **Write event:** < 1ms (append to JSON array)
- **Load media:** Filesystem I/O (same with database)

**Verdict:** Fast enough for this use case.

---

## 🗄️ Database Approach (Optional)

### When Database Helps

✅ **Complex queries** - "Find all creatives cached in last 24 hours"
✅ **Indexing** - Fast lookups by multiple keys
✅ **Transactions** - Atomic operations
✅ **Relationships** - Link playlists to media files
✅ **Concurrent access** - Multiple processes
✅ **Large datasets** - Thousands of items

### When Database is Overkill

❌ **Simple key-value lookups** - Filesystem is fine
❌ **Small datasets** - < 1000 items
❌ **No complex queries** - Just load/save
❌ **Single process** - No concurrency needed
❌ **Simple data structure** - JSON is sufficient

### Database Options (If Needed)

#### SQLite (Recommended if using DB)
- **Size:** ~500KB library
- **Pros:** SQL queries, ACID transactions, embedded
- **Cons:** Slight overhead, more complex

#### BoltDB (Golang)
- **Size:** ~100KB library
- **Pros:** Key-value, simple API, fast
- **Cons:** No SQL queries

#### IndexedDB (Web/Electron)
- **Size:** Built into browser
- **Pros:** Structured storage, async API
- **Cons:** Browser-only, complex API

---

## 🔍 Use Case Analysis

### Your Requirements:

| Operation | Frequency | Complexity | Filesystem OK? |
|-----------|-----------|------------|----------------|
| Load current playlist | Every 3 min | Simple read | ✅ Yes |
| Save new playlist | Every 3 min | Simple write | ✅ Yes |
| Check if media cached | Per item | File exists | ✅ Yes |
| Download & cache media | Per item | File write | ✅ Yes |
| Queue event | Per event | Append to array | ✅ Yes |
| Sync events | On reconnect | Read all, clear | ✅ Yes |
| Find cached media | Per item | File path lookup | ✅ Yes |
| Cache cleanup | Periodically | List files, delete | ✅ Yes |

**Conclusion:** All operations are simple. **Filesystem is sufficient.**

---

## 📈 Comparison

### Filesystem-Only

```
Complexity:     ████░░░░░░  Low
Performance:    ████████░░  Excellent
Dependencies:   ██░░░░░░░░  None
Setup Time:     ██░░░░░░░░  Minutes
Maintenance:    ████░░░░░░  Easy
```

### With SQLite

```
Complexity:     ████████░░  Medium
Performance:    █████████░  Excellent
Dependencies:   ██████░░░░  SQLite lib
Setup Time:     ██████░░░░  Hours
Maintenance:    ██████░░░░  Medium
```

---

## 💡 Recommendation

### **Use Filesystem-Only** ✅

**Reasons:**
1. ✅ All operations are simple (read/write JSON, check file exists)
2. ✅ No complex queries needed
3. ✅ Small dataset (< 1000 items)
4. ✅ Single process (no concurrency issues)
5. ✅ Faster development (no DB setup)
6. ✅ Easier debugging (inspect files directly)
7. ✅ Zero dependencies
8. ✅ Smaller binary size

### When to Consider Database

Consider SQLite/BoltDB **only if**:
- You need complex queries (e.g., "find all videos cached before date X")
- You have > 10,000 cached items
- You need transactions for data integrity
- You have multiple processes accessing cache
- You need advanced indexing

**For your use case:** None of these apply. **Filesystem is perfect.**

---

## 🛠️ Optimized Filesystem Implementation

### Efficient File Structure

```
cache/
├── playlists/
│   ├── current.json           # Active playlist
│   └── backup.json            # Previous playlist (fallback)
├── media/
│   ├── images/
│   │   ├── creative-1_a1b2c3.jpg
│   │   └── creative-2_d4e5f6.png
│   └── videos/
│       ├── creative-3_g7h8i9.mp4
│       └── creative-4_j0k1l2.webm
├── events/
│   └── queue.json             # Event queue
└── index.json                  # Cache index (optional optimization)
```

### Cache Index (Optional Optimization)

If you want faster lookups without a database, use a simple index file:

```json
{
  "media": {
    "creative-1": {
      "path": "cache/media/images/creative-1_a1b2c3.jpg",
      "cachedAt": "2024-12-18T15:00:00Z",
      "size": 245760,
      "type": "image"
    },
    "creative-2": {
      "path": "cache/media/videos/creative-2_d4e5f6.mp4",
      "cachedAt": "2024-12-18T15:01:00Z",
      "size": 5242880,
      "type": "video"
    }
  },
  "playlist": {
    "current": "cache/playlists/current.json",
    "cachedAt": "2024-12-18T15:00:00Z"
  }
}
```

**Benefits:**
- Fast lookup without scanning directory
- Track cache size without reading all files
- Know what's cached without file system calls

**Trade-off:** Need to keep index in sync with filesystem

---

## 🚀 Implementation Strategy

### Phase 1: Start with Filesystem-Only

1. ✅ Use JSON files for playlists
2. ✅ Use filesystem for media files
3. ✅ Use JSON array for event queue
4. ✅ Simple file existence checks

### Phase 2: Add Index (If Needed)

1. ✅ Add `index.json` for faster lookups
2. ✅ Update index when caching media
3. ✅ Use index for cache size calculation

### Phase 3: Consider Database (Only If Needed)

1. ✅ Only if you hit performance issues
2. ✅ Only if you need complex queries
3. ✅ Migrate to SQLite/BoltDB if necessary

**Start simple, optimize later.**

---

## 📝 Code Examples: Filesystem-Only

### Rust (Tauri) - Simple Approach

```rust
use std::fs;
use std::path::PathBuf;

pub struct Cache {
    cache_dir: PathBuf,
}

impl Cache {
    pub fn new(cache_dir: PathBuf) -> Self {
        fs::create_dir_all(&cache_dir.join("playlists")).ok();
        fs::create_dir_all(&cache_dir.join("media/images")).ok();
        fs::create_dir_all(&cache_dir.join("media/videos")).ok();
        fs::create_dir_all(&cache_dir.join("events")).ok();
        Self { cache_dir }
    }

    // Save playlist
    pub fn save_playlist(&self, playlist: &Playlist) -> Result<(), Box<dyn std::error::Error>> {
        let path = self.cache_dir.join("playlists/current.json");
        let json = serde_json::to_string_pretty(playlist)?;
        fs::write(path, json)?;
        Ok(())
    }

    // Load playlist
    pub fn load_playlist(&self) -> Option<Playlist> {
        let path = self.cache_dir.join("playlists/current.json");
        fs::read_to_string(path)
            .ok()
            .and_then(|json| serde_json::from_str(&json).ok())
    }

    // Check if media cached
    pub fn is_media_cached(&self, creative_id: &str, hash: &str, ext: &str) -> Option<PathBuf> {
        let path = self.cache_dir.join(format!("media/images/{}_{}.{}", creative_id, hash, ext));
        if path.exists() {
            Some(path)
        } else {
            None
        }
    }

    // Save media file
    pub fn save_media(&self, creative_id: &str, hash: &str, ext: &str, data: &[u8]) -> Result<PathBuf, Box<dyn std::error::Error>> {
        let path = self.cache_dir.join(format!("media/images/{}_{}.{}", creative_id, hash, ext));
        fs::write(&path, data)?;
        Ok(path)
    }
}
```

### Golang (Wails) - Simple Approach

```go
package storage

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type Cache struct {
    cacheDir string
}

func NewCache(cacheDir string) *Cache {
    os.MkdirAll(filepath.Join(cacheDir, "playlists"), 0755)
    os.MkdirAll(filepath.Join(cacheDir, "media/images"), 0755)
    os.MkdirAll(filepath.Join(cacheDir, "media/videos"), 0755)
    os.MkdirAll(filepath.Join(cacheDir, "events"), 0755)
    return &Cache{cacheDir: cacheDir}
}

// Save playlist
func (c *Cache) SavePlaylist(playlist *Playlist) error {
    path := filepath.Join(c.cacheDir, "playlists/current.json")
    data, err := json.MarshalIndent(playlist, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

// Load playlist
func (c *Cache) LoadPlaylist() (*Playlist, error) {
    path := filepath.Join(c.cacheDir, "playlists/current.json")
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var playlist Playlist
    return &playlist, json.Unmarshal(data, &playlist)
}

// Check if media cached
func (c *Cache) IsMediaCached(creativeID, hash, ext string) (string, bool) {
    path := filepath.Join(c.cacheDir, "media/images", fmt.Sprintf("%s_%s.%s", creativeID, hash, ext))
    _, err := os.Stat(path)
    return path, err == nil
}
```

---

## ✅ Final Recommendation

### **Use Filesystem-Only** ✅

**Why:**
- ✅ Simple to implement
- ✅ Fast enough for your use case
- ✅ No dependencies
- ✅ Easy to debug
- ✅ Sufficient for all operations

**When to reconsider:**
- ❌ If you need complex queries (you don't)
- ❌ If you have > 10,000 items (you won't)
- ❌ If you need transactions (you don't)
- ❌ If multiple processes access cache (single client)

**Start with filesystem. Add database only if you actually need it.**

---

**Status:** Filesystem is sufficient  
**Recommendation:** Start simple, optimize later if needed

