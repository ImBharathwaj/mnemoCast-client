# Manual Ad Testing Guide

This guide explains how to test the ad player system with manually created ads and media files.

## Quick Start

### 1. Create Sample Ads

```bash
# Create sample ads (includes a text ad that works immediately)
go run cmd/test-ads/main.go create-sample

# Or build the utility first
go build -o bin/test-ads ./cmd/test-ads
./bin/test-ads create-sample
```

### 2. Add Your Own Media

```bash
# Add an image ad
./bin/test-ads add-image ad-001 /path/to/your/image.jpg "My Image Ad" 30

# Add a video ad
./bin/test-ads add-video ad-002 /path/to/your/video.mp4 "My Video Ad" 60

# Add a text ad
./bin/test-ads add-text ad-003 "Hello, this is a test ad!" "Text Ad" 15
```

### 3. List Current Ads

```bash
./bin/test-ads list
```

### 4. Start the Screen Client

```bash
./bin/screen
```

The player will automatically load and play the ads you've created!

## Commands Reference

### `create-sample`
Creates sample ads for testing, including a working text ad.

```bash
./bin/test-ads create-sample
```

### `add-image`
Adds an image ad to the playlist.

```bash
./bin/test-ads add-image <ad-id> <image-path> [title] [duration]
```

**Example:**
```bash
./bin/test-ads add-image ad-001 ~/Pictures/banner.jpg "Promotional Banner" 30
```

### `add-video`
Adds a video ad to the playlist.

```bash
./bin/test-ads add-video <ad-id> <video-path> [title] [duration]
```

**Example:**
```bash
./bin/test-ads add-video ad-002 ~/Videos/commercial.mp4 "Product Commercial" 60
```

### `add-text`
Adds a text ad to the playlist.

```bash
./bin/test-ads add-text <ad-id> <text-content> [title] [duration]
```

**Example:**
```bash
./bin/test-ads add-text ad-003 "Welcome to our store!" "Welcome Message" 15
```

### `list`
Lists all current ads in the playlist.

```bash
./bin/test-ads list
```

### `clear`
Clears all ads from the playlist.

```bash
./bin/test-ads clear
```

## Testing Workflow

### Step 1: Prepare Media Files

1. **For Images:**
   - Use common formats: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`
   - Recommended size: 1920x1080 or your screen resolution

2. **For Videos:**
   - Use common formats: `.mp4`, `.webm`, `.mov`, `.avi`
   - Recommended: MP4 with H.264 codec for best compatibility

3. **For Text:**
   - Just provide the text content directly
   - The utility will create a text file automatically

### Step 2: Add Ads

```bash
# Add multiple ads
./bin/test-ads add-image ad-001 ~/Pictures/ad1.jpg "Ad 1" 30
./bin/test-ads add-image ad-002 ~/Pictures/ad2.jpg "Ad 2" 30
./bin/test-ads add-text ad-003 "Special Offer Today!" "Promotion" 15
```

### Step 3: Verify Ads

```bash
# List ads to verify they were added
./bin/test-ads list
```

### Step 4: Test Player

```bash
# Start the screen client
./bin/screen
```

The player will:
1. Load ads from storage
2. Download/copy media files (if needed)
3. Start playing ads in a loop
4. Respect durations and priorities

## File Locations

- **Ads metadata:** `~/.mnemocast/ads/current_ads.json`
- **Media files:** `~/.mnemocast/ads/media/{ad-id}/`
- **Config:** `~/.mnemocast/config.json`

## Testing Different Scenarios

### Test Priority Sorting

```bash
./bin/test-ads add-image ad-low "Low Priority" ~/Pictures/ad1.jpg "Ad 1" 30
./bin/test-ads add-image ad-high "High Priority" ~/Pictures/ad2.jpg "Ad 2" 30
```

Then edit `~/.mnemocast/ads/current_ads.json` and set:
- `ad-low` priority to `1`
- `ad-high` priority to `2`

Higher priority ads will play first.

### Test Time-Based Scheduling

Edit `~/.mnemocast/ads/current_ads.json` and add `startTime` and `endTime`:

```json
{
  "id": "ad-001",
  "startTime": "2025-12-19T10:00:00Z",
  "endTime": "2025-12-19T18:00:00Z",
  ...
}
```

Ads will only play during their scheduled time window.

### Test Single Ad Loop

```bash
# Add one ad
./bin/test-ads add-image ad-single ~/Pictures/single.jpg "Single Ad" 30

# Start player - it will loop this single ad
./bin/screen
```

### Test Multiple Ad Types

```bash
# Mix different ad types
./bin/test-ads add-image ad-img ~/Pictures/image.jpg "Image" 30
./bin/test-ads add-video ad-vid ~/Videos/video.mp4 "Video" 60
./bin/test-ads add-text ad-txt "Text content" "Text" 15
```

## Troubleshooting

### Player Not Starting

1. **Check if ads exist:**
   ```bash
   ./bin/test-ads list
   ```

2. **Check logs:**
   Look for `[PLAYER]` log messages in the screen client output

3. **Verify media files:**
   ```bash
   ls -la ~/.mnemocast/ads/media/
   ```

### Media Not Displaying

1. **Check file permissions:**
   ```bash
   ls -la ~/.mnemocast/ads/media/{ad-id}/
   ```

2. **Verify renderer availability:**
   - Image: Check if `feh`, `imv`, `sxiv`, or `xdg-open` is installed
   - Video: Check if `mpv`, `vlc`, `ffplay`, or `xdg-open` is installed
   - HTML: Check if `firefox`, `chromium`, or browser is installed

3. **Check file format:**
   - Ensure media files are in supported formats
   - Try opening files manually to verify they're not corrupted

### Ads Not Updating

1. **Clear and re-add:**
   ```bash
   ./bin/test-ads clear
   ./bin/test-ads add-image ad-001 ~/Pictures/new.jpg "New Ad" 30
   ```

2. **Restart screen client:**
   - Stop the client (Ctrl+C)
   - Start again: `./bin/screen`

## Example Test Session

```bash
# 1. Create sample ads
./bin/test-ads create-sample

# 2. Add a real image ad
./bin/test-ads add-image my-ad ~/Pictures/banner.jpg "My Banner" 30

# 3. List all ads
./bin/test-ads list

# 4. Start the player
./bin/screen

# 5. Watch the ads play in a loop!

# 6. When done testing, clear ads
./bin/test-ads clear
```

## Notes

- The player supports `file://` URLs for local testing
- Media files are copied to the media directory for organization
- Ads are stored in JSON format and can be manually edited
- The player automatically loads ads on startup
- New ads are picked up when the fetcher updates (or manually via utility)

