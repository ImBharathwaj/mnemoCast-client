# 🚀 Running the MnemoCast Screen System

**Complete guide to build, configure, and run the screen system.**

---

## 📋 Prerequisites

### Required Software

1. **Go 1.22+**
   ```bash
   # Check if Go is installed
   go version
   
   # If not installed (Ubuntu/Debian):
   sudo apt update
   sudo apt install golang-go
   
   # Or download from: https://go.dev/dl/
   ```

2. **Backend Running**
   - Ad server must be running at `http://10.42.0.1:8080`
   - Backend should be accessible from the client machine

---

## 🔨 Building the System

### Step 1: Build the Application

```bash
# Navigate to project directory
cd /home/bharathwaj/Code/mnemoCast-client

# Build the application
go build -o bin/screen ./cmd/screen

# Verify build
ls -lh bin/screen
```

**Expected Output:**
- Binary created at `bin/screen`
- Size: ~2-3 MB

### Step 2: Make Executable (if needed)

```bash
chmod +x bin/screen
```

---

## ⚙️ First-Time Setup

### Step 1: Run the Application

```bash
./bin/screen
```

### Step 2: Initial Configuration

On first run, the system will:

1. **Create Configuration Directory**
   - Creates `~/.mnemocast/` directory
   - Sets up secure file permissions

2. **Generate Screen Identity**
   - Generates unique UUID-based screen ID
   - Creates `identity.json` file
   - Example output:
     ```
     ✅ Screen ID: screen-550e8400-e29b-41d4-a716-446655440000
     ```

3. **Create Default Configuration**
   - Creates `config.json` with default settings
   - Ad Server URL: `http://10.42.0.1:8080`
   - Heartbeat Interval: 30 seconds

4. **Prompt for API Key**
   - If credentials don't exist, prompts:
     ```
     Would you like to configure API key now? (y/n):
     ```
   - Enter `y` and provide your API key
   - Credentials are encrypted and stored securely

---

## 🔐 Configuring Credentials

### Option 1: Interactive Setup (First Run)

When prompted:
```
Would you like to configure API key now? (y/n): y
Enter API Key: your-api-key-here
```

### Option 2: Manual Configuration

If you need to update credentials later:

1. **Delete existing credentials** (if needed):
   ```bash
   rm ~/.mnemocast/credentials.json.enc
   ```

2. **Run the application again** - it will prompt for API key

### Option 3: Environment Variable (Future)

For containerized deployments, you can set:
```bash
export MNEMOCAST_API_KEY=your-api-key-here
```

---

## 🏃 Running the System

### Basic Run

```bash
./bin/screen
```

### What Happens:

**With Credentials Configured:**
1. **Initialization**
   ```
   🖥️  MnemoCast Screen System
   ============================
   Config directory: /home/user/.mnemocast
   
   📋 Loading screen identity...
   ✅ Screen ID: screen-xxx-xxx-xxx
   
   ⚙️  Loading configuration...
   ✅ Ad Server URL: http://10.42.0.1:8080
   
   🔐 Checking credentials...
   ✅ API Key: Configured
   ```

**Without Credentials:**
1. **Initialization** (same as above)
2. **Credentials Check**
   ```
   🔐 Checking credentials...
   ⚠️  Credentials: Not configured
   Would you like to configure API key now? (y/n):
   ```
3. **System Continues Running**
   ```
   ⚠️  System is running but heartbeat is not active
      Credentials are required to connect to ad server
   
   💡 To enable heartbeat system:
      1. Stop this process (Ctrl+C)
      2. Run again: ./bin/screen
      3. Enter 'y' when prompted for API key
   
   🔄 Waiting for credentials...
      Press Ctrl+C to exit
   ```
   - System stays running and waits
   - Reminds every 60 seconds to configure credentials
   - Can be stopped with Ctrl+C

2. **Ad Server Connection**
   ```
   🌐 Connecting to ad server...
      Registering screen...
      ✅ Registered successfully!
      Screen ID: screen-xxx-xxx-xxx
      Status: Online
   
      Testing heartbeat...
      ✅ Heartbeat successful!
   ```

3. **Heartbeat System Starts**
   ```
   💓 Starting heartbeat scheduler...
      ✅ Heartbeat scheduler started (interval: 30 seconds)
   
   ✅ Screen system initialized successfully!
   
   💓 Heartbeat system is running...
      Press Ctrl+C to stop
   ```

4. **Status Updates** (every 30 seconds)
   ```
   [✅] Heartbeat Status: Connected (last sent: 5s ago)
   [✅] Heartbeat Status: Connected (last sent: 30s ago)
   ```

### Stop the System

Press `Ctrl+C` for graceful shutdown:

**With Heartbeat Running:**
```
🛑 Shutting down...
Heartbeat scheduler stopped
✅ Shutdown complete
```

**Without Credentials:**
```
🛑 Shutting down...
✅ Shutdown complete
```

**Note:** The system runs continuously until you press Ctrl+C, even without credentials configured.

---

## 🔧 Configuration

### Configuration File Location

`~/.mnemocast/config.json`

### Default Configuration

```json
{
  "identity": {
    "id": "screen-xxx-xxx-xxx",
    "name": "Default Screen",
    "location": {
      "city": "Chennai",
      "area": "Airport",
      "venueType": "airport"
    },
    "classification": 1
  },
  "adServerUrl": "http://10.42.0.1:8080",
  "heartbeatInterval": 30,
  "retryAttempts": 3,
  "retryDelay": 5
}
```

### Modifying Configuration

1. **Edit config file:**
   ```bash
   nano ~/.mnemocast/config.json
   ```

2. **Or use environment variables** (if supported):
   ```bash
   export MNEMOCAST_AD_SERVER_URL=http://10.42.0.1:8080
   export MNEMOCAST_HEARTBEAT_INTERVAL=30
   ```

3. **Restart the application** to apply changes

---

## 🧪 Testing the System

### Test 1: Verify Identity

```bash
./bin/screen
# Check that screen ID is displayed
# Verify identity.json exists in ~/.mnemocast/
```

### Test 2: Verify Credentials

```bash
./bin/screen
# Check that API key is configured
# Verify credentials.json.enc exists
```

### Test 3: Test Backend Connection

```bash
# Make sure backend is running
curl http://10.42.0.1:8080/api/v1/health

# Run screen system
./bin/screen
# Should see "Registered successfully!"
```

### Test 4: Verify Heartbeat

```bash
./bin/screen
# Wait 30 seconds
# Should see status updates showing "Connected"
```

---

## 🐛 Troubleshooting

### Issue: "Failed to get home directory"

**Solution:**
```bash
# Ensure HOME environment variable is set
echo $HOME
export HOME=/home/your-username
```

### Issue: "Registration failed"

**Possible Causes:**
1. Backend not running
   ```bash
   # Check backend status
   curl http://10.42.0.1:8080/api/v1/health
   ```

2. Wrong API key
   - Delete credentials and reconfigure:
     ```bash
     rm ~/.mnemocast/credentials.json.enc
     ./bin/screen
     ```

3. Network connectivity
   ```bash
   # Test connectivity
   ping 10.42.0.1
   curl http://10.42.0.1:8080/api/v1/health
   ```

### Issue: "Heartbeat failed"

**Possible Causes:**
1. Backend went offline
   - System will retry automatically
   - Check backend logs

2. Network issues
   - Check network connectivity
   - Verify firewall rules

3. Invalid screen ID
   - Check `~/.mnemocast/identity.json`
   - Verify screen is registered in backend

### Issue: "Credentials not found"

**Solution:**
```bash
# Run application - it will prompt for API key
./bin/screen
# Enter 'y' when prompted
# Provide your API key
```

### Issue: Build fails

**Solution:**
```bash
# Update dependencies
go mod tidy

# Clean build
rm -rf bin/
go build -o bin/screen ./cmd/screen
```

---

## 🔄 Running as a Service

### Systemd Service (Linux)

Create `/etc/systemd/system/mnemocast-screen.service`:

```ini
[Unit]
Description=MnemoCast Screen System
After=network.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/home/bharathwaj/Code/mnemoCast-client
ExecStart=/home/bharathwaj/Code/mnemoCast-client/bin/screen
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Enable and start:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable mnemocast-screen
sudo systemctl start mnemocast-screen
sudo systemctl status mnemocast-screen
```

### Docker (Container)

Create `Dockerfile`:
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o screen ./cmd/screen

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/screen .
CMD ["./screen"]
```

**Build and run:**
```bash
docker build -t mnemocast-screen .
docker run -d --name mnemocast-screen mnemocast-screen
```

---

## 📊 Monitoring & Status Checking

### How to Check if System is Running

#### Method 1: Check Running Process

```bash
# Check if screen process is running
ps aux | grep screen | grep -v grep

# Or using pgrep
pgrep -f "bin/screen"

# Check with process details
ps -ef | grep "[b]in/screen"
```

**Expected Output:**
```
bharathwaj  12345  0.1  0.2  7456784  12345 ?  Sl   22:38   0:05 ./bin/screen
```

**If running:** You'll see the process  
**If not running:** No output

#### Method 2: Check Process by Name

```bash
# Check if process exists
pgrep -f "mnemoCast-client" || echo "Not running"

# Count processes
pgrep -c -f "bin/screen"
```

#### Method 3: Check Systemd Service (if running as service)

```bash
# Check service status
sudo systemctl status mnemocast-screen

# Check if service is active
sudo systemctl is-active mnemocast-screen

# Check if service is enabled
sudo systemctl is-enabled mnemocast-screen
```

**Expected Output (if running):**
```
● mnemocast-screen.service - MnemoCast Screen System
     Loaded: loaded (/etc/systemd/system/mnemocast-screen.service)
     Active: active (running) since Wed 2024-12-18 22:38:44 UTC
```

#### Method 4: Check Configuration Files

```bash
# Check if configuration directory exists
ls -la ~/.mnemocast/

# Check identity file
cat ~/.mnemocast/identity.json

# Check config file
cat ~/.mnemocast/config.json

# Check credentials (encrypted)
ls -la ~/.mnemocast/credentials.json.enc
```

**If files exist:** System has been configured  
**If files missing:** System hasn't run yet or was cleaned

#### Method 5: Check Backend for Registration

```bash
# Get screen ID from identity file
SCREEN_ID=$(cat ~/.mnemocast/identity.json | grep -o '"id":"[^"]*' | cut -d'"' -f4)

# Check if screen is registered in backend
curl http://10.42.0.1:8080/api/v1/screens/$SCREEN_ID

# Or check all screens
curl http://10.42.0.1:8080/api/v1/screens
```

**If registered:** Backend returns screen information  
**If not registered:** 404 or error response

#### Method 6: Check Logs

```bash
# If running as systemd service
sudo journalctl -u mnemocast-screen -n 50 --no-pager

# Check recent logs
sudo journalctl -u mnemocast-screen --since "10 minutes ago"

# Follow logs in real-time
sudo journalctl -u mnemocast-screen -f
```

**Look for:**
- "Heartbeat scheduler started"
- "Heartbeat Status: Connected"
- Recent heartbeat timestamps

#### Method 7: Check Network Connections

```bash
# Check if process has open connection to backend
netstat -tulpn | grep 10.42.0.1:8080

# Or using ss
ss -tulpn | grep 10.42.0.1:8080

# Check with lsof
sudo lsof -i :8080 | grep screen
```

**If connected:** You'll see established connections  
**If not connected:** No output

#### Method 8: Quick Status Script

Create a status check script:

```bash
#!/bin/bash
# save as check-status.sh

echo "🔍 Checking MnemoCast Screen System Status..."
echo ""

# Check process
if pgrep -f "bin/screen" > /dev/null; then
    echo "✅ Process: RUNNING"
    ps aux | grep "[b]in/screen" | awk '{print "   PID:", $2, "CPU:", $3"%", "MEM:", $4"%"}'
else
    echo "❌ Process: NOT RUNNING"
fi

echo ""

# Check configuration
if [ -f ~/.mnemocast/identity.json ]; then
    echo "✅ Configuration: EXISTS"
    SCREEN_ID=$(cat ~/.mnemocast/identity.json | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "   Screen ID: $SCREEN_ID"
else
    echo "⚠️  Configuration: NOT FOUND"
fi

echo ""

# Check credentials
if [ -f ~/.mnemocast/credentials.json.enc ]; then
    echo "✅ Credentials: CONFIGURED"
else
    echo "⚠️  Credentials: NOT CONFIGURED"
fi

echo ""

# Check backend connectivity
if curl -s -o /dev/null -w "%{http_code}" http://10.42.0.1:8080/api/v1/health | grep -q "200"; then
    echo "✅ Backend: ACCESSIBLE"
else
    echo "❌ Backend: NOT ACCESSIBLE"
fi
```

**Usage:**
```bash
chmod +x check-status.sh
./check-status.sh
```

### Verify Heartbeat

The application displays status every 30 seconds when running:
```
[✅] Heartbeat Status: Connected (last sent: 30s ago)
```

**To verify heartbeat is working:**
1. Check if process is running (Method 1)
2. Check backend logs for heartbeat requests
3. Check last heartbeat time in backend

### Check Backend

```bash
# Verify screen is registered
SCREEN_ID=$(cat ~/.mnemocast/identity.json | grep -o '"id":"[^"]*' | cut -d'"' -f4)
curl http://10.42.0.1:8080/api/v1/screens/$SCREEN_ID

# Check heartbeat status in backend
# (depends on your backend API)
curl http://10.42.0.1:8080/api/v1/screens/$SCREEN_ID/status
```

### Quick Status Check Commands

```bash
# One-liner: Is it running?
pgrep -f "bin/screen" && echo "✅ Running" || echo "❌ Not running"

# One-liner: Check with details
ps aux | grep "[b]in/screen" && echo "✅ Running" || echo "❌ Not running"

# One-liner: Check service (if systemd)
systemctl is-active mnemocast-screen 2>/dev/null && echo "✅ Running" || echo "❌ Not running"
```

### Automated Status Check Script

Use the provided status check script:

```bash
# Run status check
./scripts/check-status.sh
```

**Output Example:**
```
🔍 Checking MnemoCast Screen System Status...
==============================================

📋 Process Status:
   ✅ RUNNING
   PID: 12345
   CPU: 0.1%
   Memory: 0.2%
   Runtime: 00:05:23

⚙️  Configuration Status:
   ✅ Identity file exists
   Screen ID: screen-00029ab6-755b-4099-ba86-866050f2eec5
   Screen Name: Unnamed Screen
   ✅ Config file exists
   Ad Server: http://10.42.0.1:8080
   Heartbeat Interval: 30s

🔐 Credentials Status:
   ✅ Credentials file exists (encrypted)
   ✅ Encryption key exists

🌐 Backend Connectivity:
   ✅ Backend is ACCESSIBLE
   URL: http://10.42.0.1:8080

📡 Registration Status:
   ✅ Screen is REGISTERED in backend
```

---

## 🔒 Security Checklist

- ✅ Credentials encrypted (AES-256-GCM)
- ✅ Configuration files have restricted permissions (0600)
- ✅ No plain-text credential storage
- ✅ API key masked in logs
- ✅ Secure key generation

---

## 📝 Quick Reference

### Build
```bash
go build -o bin/screen ./cmd/screen
```

### Run
```bash
./bin/screen
```

### Stop
```
Press Ctrl+C
```

### Configuration Directory
```
~/.mnemocast/
```

### Configuration Files
- `identity.json` - Screen identity
- `config.json` - Application config
- `credentials.json.enc` - Encrypted credentials
- `.encryption_key` - Encryption key

### Backend URL
```
http://10.42.0.1:8080
```

---

## ✅ Verification Checklist

Before considering the system ready:

- [ ] Application builds successfully
- [ ] Screen identity generated
- [ ] Configuration created
- [ ] API key configured
- [ ] Backend connection successful
- [ ] Registration successful
- [ ] Heartbeat scheduler running
- [ ] Status updates appearing
- [ ] Graceful shutdown works

---

## 🆘 Getting Help

### Check Logs

The application logs to stdout/stderr. If running as service:
```bash
sudo journalctl -u mnemocast-screen -n 100
```

### Verify Backend

```bash
# Health check
curl http://10.42.0.1:8080/api/v1/health

# Test registration endpoint
curl -X POST http://10.42.0.1:8080/api/v1/screens/register \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"id":"test","name":"Test","location":{"city":"Test","area":"Test","venueType":"test"},"classification":1}'
```

### Common Issues

See **Troubleshooting** section above for detailed solutions.

---

**Status:** Ready to Run  
**Last Updated:** December 2024

