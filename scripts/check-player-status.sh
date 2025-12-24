#!/bin/bash

echo "=== Player Status Check ==="
echo

# Check if ads exist
ADS_FILE="$HOME/.mnemocast/ads/current_ads.json"
if [ -f "$ADS_FILE" ]; then
    echo "✓ Ads file exists: $ADS_FILE"
    AD_COUNT=$(jq '.ads | length' "$ADS_FILE" 2>/dev/null || echo "0")
    echo "  Ads count: $AD_COUNT"
    if [ "$AD_COUNT" -gt 0 ]; then
        echo "  First ad:"
        jq '.ads[0] | {id, type, contentUrl, duration}' "$ADS_FILE" 2>/dev/null || echo "    (Could not parse)"
    fi
else
    echo "✗ Ads file not found: $ADS_FILE"
fi
echo

# Check media directory
MEDIA_DIR="$HOME/.mnemocast/ads/media"
if [ -d "$MEDIA_DIR" ]; then
    echo "✓ Media directory exists: $MEDIA_DIR"
    MEDIA_COUNT=$(find "$MEDIA_DIR" -type f | wc -l)
    echo "  Media files: $MEDIA_COUNT"
    echo "  Media directories:"
    ls -1 "$MEDIA_DIR" 2>/dev/null | head -5 | sed 's/^/    /'
else
    echo "✗ Media directory not found: $MEDIA_DIR"
fi
echo

# Check available renderers
echo "Available Renderers:"
if command -v xdg-open &> /dev/null; then
    echo "  ✓ xdg-open (images)"
else
    echo "  ✗ xdg-open"
fi

if command -v vlc &> /dev/null; then
    echo "  ✓ vlc (videos)"
else
    echo "  ✗ vlc"
fi

if command -v firefox &> /dev/null; then
    echo "  ✓ firefox (HTML)"
else
    echo "  ✗ firefox"
fi
echo

echo "To test the player:"
echo "  1. Add ads: ./bin/test-ads create-sample"
echo "  2. Start player: ./bin/screen"
echo "  3. Check logs for [PLAYER] and [RENDER] messages"

