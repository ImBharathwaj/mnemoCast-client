# 🚀 Getting Started with MnemoCast Client

## ✅ Project Setup Complete!

All project files have been created. Here's what's been set up:

### 📁 Project Structure
```
✅ app/internal/models/     - Data models (Screen, Playlist, Events)
✅ app/internal/api/        - API client (Screen, Playlist, Events APIs)
✅ app/internal/app.go       - Wails app with exposed methods
✅ app/main.go               - Application entry point
✅ frontend/                 - Frontend UI (HTML, CSS, JS)
✅ go.mod                    - Go module definition
✅ wails.json                - Wails configuration
```

## 🔧 Installation Steps

### Step 1: Install Go

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install golang-go
```

**Or download from:** https://go.dev/dl/

**Verify installation:**
```bash
go version  # Should show go1.21 or later
```

### Step 2: Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Add Go bin to PATH (if not already)
export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
wails version
```

**Or use the setup script:**
```bash
./setup.sh
```

### Step 3: Install Dependencies

```bash
cd /home/bharathwaj/Code/mnemoCast-client
go mod tidy
```

### Step 4: Platform-Specific Setup

**Linux:**
```bash
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev
```

**Windows:**
- Install TDM-GCC or MinGW-w64
- Wails will guide you through setup

**macOS:**
```bash
xcode-select --install
```

## 🚀 Running the Application

### Development Mode

```bash
# Make sure backend is running first
cd ../backend
sbt run

# In another terminal, run client
cd /home/bharathwaj/Code/mnemoCast-client
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

## 🧪 Testing

1. **Start Backend:**
   ```bash
   cd ../backend
   sbt run
   ```

2. **Run Client:**
   ```bash
   wails dev
   ```

3. **Expected Behavior:**
   - ✅ Screen registers automatically on startup
   - ✅ Playlist fetches from backend
   - ✅ Media displays (images/videos)
   - ✅ Events are sent (impression, play)
   - ✅ Heartbeat sends every 30 seconds

## 📝 Configuration

Default configuration is in `app/internal/app.go`:

```go
ScreenID:                "screen-1"
Name:                    "Default Screen"
Location:                {City: "Chennai", Area: "Airport", VenueType: "airport"}
BackendURL:              "http://localhost:8080"
PlaylistRefreshInterval: 3  // minutes
HeartbeatInterval:       30 // seconds
```

## 🐛 Troubleshooting

### Issue: `go: command not found`
- **Solution:** Install Go (see Step 1)

### Issue: `wails: command not found`
- **Solution:** Install Wails CLI (see Step 2)
- **Check:** Make sure `$(go env GOPATH)/bin` is in your PATH

### Issue: `wails dev` fails
- **Check:** Node.js is installed (`node --version`)
- **Check:** Frontend directory exists
- **Check:** Platform-specific dependencies installed

### Issue: Cannot connect to backend
- **Check:** Backend is running on port 8080
- **Check:** CORS is enabled in backend
- **Check:** BackendURL in config is correct

### Issue: Build fails
- **Check:** Go version is 1.21+
- **Check:** All dependencies installed (`go mod tidy`)
- **Check:** Platform-specific build tools installed

## 📚 Next Steps

1. ✅ **Test basic functionality** - Run `wails dev` and verify everything works
2. ⏭️ **Add offline caching** - Implement filesystem cache (see `OFFLINE_CACHING_STRATEGY.md`)
3. ⏭️ **Add error handling** - Retry logic, offline detection
4. ⏭️ **Add preloading** - Preload next items for smooth playback
5. ⏭️ **Add event queue** - Queue events when offline

## 📖 Documentation

- **Quick Start:** `GOLANG_QUICKSTART.md`
- **Full Plan:** `GOLANG_WAILS_IMPLEMENTATION_PLAN.md`
- **Offline Caching:** `OFFLINE_CACHING_STRATEGY.md`
- **Storage Analysis:** `CACHE_STORAGE_ANALYSIS.md`

---

**Status:** ✅ Project structure created, ready for Go installation and testing!

