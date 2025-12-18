# 🚀 Quick Start - MnemoCast Client

## ✅ Setup Complete!

Your project is ready. The backend URL is configured to: **http://10.42.0.1:8080**

## 🏃 Run the Application

### Option 1: Using the run script
```bash
./run.sh
```

### Option 2: Direct command
```bash
# Make sure Wails is in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Run in development mode
wails dev
```

## 📋 What's Configured

- ✅ Backend URL: `http://10.42.0.1:8080`
- ✅ Screen ID: `screen-1`
- ✅ Playlist refresh: Every 3 minutes
- ✅ Heartbeat: Every 30 seconds

## 🔧 If Wails Command Not Found

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Or add to ~/.bashrc permanently
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

## 🧪 Testing

1. **Make sure backend is running** on `10.42.0.1:8080`
2. **Run the client:**
   ```bash
   ./run.sh
   ```
3. **Expected behavior:**
   - Screen registers automatically
   - Playlist fetches from backend
   - Media displays
   - Events are sent
   - Heartbeat runs every 30 seconds

## 📝 Change Backend URL

Edit `app/internal/app.go`:
```go
BackendURL: "http://10.42.0.1:8080",  // Change this line
```

Or use the `SetConfig` method from frontend after app starts.

## 🐛 Troubleshooting

### Issue: `wails: command not found`
- **Solution:** Add Go bin to PATH (see above)

### Issue: Cannot connect to backend
- **Check:** Backend is running on `10.42.0.1:8080`
- **Check:** CORS is enabled in backend
- **Check:** Network connectivity to `10.42.0.1`

### Issue: Build errors
- **Run:** `go mod tidy`
- **Check:** Go version is 1.22+

---

**Ready to run!** Execute `./run.sh` to start the client.

