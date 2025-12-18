#!/bin/bash

# MnemoCast Client Run Script

echo "🚀 Starting MnemoCast Display Client..."
echo "Backend URL: http://10.42.0.1:8080"
echo ""

# Check if Wails is in PATH
if ! command -v wails &> /dev/null; then
    echo "❌ Wails CLI not found in PATH"
    echo "Adding Go bin to PATH..."
    export PATH=$PATH:$(go env GOPATH)/bin
    
    if ! command -v wails &> /dev/null; then
        echo "❌ Wails still not found. Please install it:"
        echo "   go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        exit 1
    fi
fi

echo "✅ Wails found: $(wails version)"
echo ""
echo "Starting development server..."
echo ""

wails dev

