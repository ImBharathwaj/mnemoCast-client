# Player Troubleshooting Guide

## Why No GUI Appears

The player uses external applications to display ads:
- **Images**: `xdg-open` (opens in default image viewer - may not be fullscreen)
- **Videos**: `vlc` (should open in fullscreen)
- **HTML**: `firefox` (opens in browser)
- **Text**: Prints to terminal (no GUI)

## Common Issues

### 1. No Ads in Storage

**Symptom:** Player doesn't start, or shows "No active ads available"

**Check:**
```bash
./scripts/check-player-status.sh
```

**Fix:**
```bash
# Add sample ads
./bin/test-ads create-sample

# Or add your own
./bin/test-ads add-image ad-001 /path/to/image.jpg "My Ad" 30
```

### 2. Ads Exist But Player Not Starting

**Check logs when running `./bin/screen`** - Look for:
- `[PLAYER] Player started`
- `[PLAYER] Playing ad: ...`
- `[RENDER] Rendering ...`

**Common causes:**
- Ads file has 0 ads (check with `./bin/test-ads list`)
- Media files not found
- Renderer commands not available

### 3. GUI Opens But Not Fullscreen

**For Images:**
- `xdg-open` uses your default image viewer
- It may not support fullscreen automatically
- **Solution:** Install a fullscreen-capable viewer:
  ```bash
  sudo apt install feh  # Lightweight fullscreen image viewer
  ```

**For Videos:**
- `vlc` should open fullscreen automatically
- If not, check VLC settings

### 4. No Renderer Available

**Check available renderers:**
```bash
./bin/check-renderers
```

**Install missing renderers:**
```bash
# Image viewer (fullscreen capable)
sudo apt install feh

# Video player (if vlc not available)
sudo apt install mpv
```

## Step-by-Step Testing

### 1. Verify Ads Exist
```bash
./bin/test-ads list
```

If empty:
```bash
./bin/test-ads create-sample
./bin/test-ads list
```

### 2. Check Media Files
```bash
ls -la ~/.mnemocast/ads/media/
```

### 3. Start Player with Verbose Logging
```bash
./bin/screen 2>&1 | grep -E "\[PLAYER\]|\[RENDER\]|\[DOWNLOAD\]"
```

You should see:
```
[PLAYER] Player started
[PLAYER] Playing ad: sample-text-001 (type: text, duration: 15s)
[RENDER] Rendering text ad: sample-text-001
```

### 4. For Image/Video Ads

**Add an image ad:**
```bash
./bin/test-ads add-image test-img /path/to/image.jpg "Test" 30
```

**Start player:**
```bash
./bin/screen
```

**Expected behavior:**
- `xdg-open` will launch your default image viewer
- Window should appear (may not be fullscreen)
- After 30 seconds, next ad should play

**For fullscreen images, install feh:**
```bash
sudo apt install feh
```

Then restart the player - it will use `feh` which supports fullscreen.

## Debugging Commands

### Check Player Status
```bash
./scripts/check-player-status.sh
```

### Check Available Renderers
```bash
./bin/check-renderers
```

### List Current Ads
```bash
./bin/test-ads list
```

### View Ads JSON
```bash
cat ~/.mnemocast/ads/current_ads.json | jq .
```

### Check Media Files
```bash
find ~/.mnemocast/ads/media -type f
```

## Expected Behavior

### Text Ads
- Display in terminal with formatting
- No GUI window
- Should see formatted text output

### Image Ads
- Opens image viewer window (via xdg-open)
- May not be fullscreen (depends on default viewer)
- Window should appear and stay open for ad duration

### Video Ads
- Opens VLC in fullscreen
- Video plays for specified duration
- Automatically closes when done

### HTML Ads
- Opens Firefox browser
- Displays HTML content
- Browser window appears

## Quick Fix Checklist

- [ ] Ads exist: `./bin/test-ads list` shows ads
- [ ] Media files exist: `ls ~/.mnemocast/ads/media/`
- [ ] Renderers available: `./bin/check-renderers`
- [ ] Player started: Check logs for `[PLAYER] Player started`
- [ ] Rendering attempted: Check logs for `[RENDER]` messages
- [ ] No errors: Check logs for `[ERROR]` messages

## Still Not Working?

1. **Check logs** - Run `./bin/screen` and look for error messages
2. **Verify ads** - Run `./bin/test-ads list`
3. **Test renderer** - Try manually opening a file:
   ```bash
   xdg-open ~/.mnemocast/ads/media/sample-text-001/sample-text-001.txt
   ```
4. **Check file permissions** - Ensure media files are readable

