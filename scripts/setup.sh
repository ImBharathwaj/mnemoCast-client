#!/bin/bash

# MnemoCast Client Setup Script

echo "🚀 Setting up MnemoCast Digital Display Client..."
echo ""

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo "📦 Installing Go..."
    
    # Detect OS and install Go
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "Detected Linux. Installing Go via apt..."
        sudo apt update
        sudo apt install -y golang-go
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "Detected macOS. Please install Go from https://go.dev/dl/"
        exit 1
    else
        echo "Please install Go from https://go.dev/dl/"
        exit 1
    fi
else
    echo "✅ Go is installed: $(go version)"
fi

# Check Wails CLI
if ! command -v wails &> /dev/null; then
    echo "📦 Installing Wails CLI..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    
    # Add Go bin to PATH if not already there
    export PATH=$PATH:$(go env GOPATH)/bin
    echo "export PATH=\$PATH:\$(go env GOPATH)/bin" >> ~/.bashrc
else
    echo "✅ Wails CLI is installed: $(wails version)"
fi

# Install dependencies
echo ""
echo "📦 Installing Go dependencies..."
go mod tidy

# Check Node.js
if command -v node &> /dev/null; then
    echo "✅ Node.js is installed: $(node --version)"
else
    echo "❌ Node.js is not installed. Please install Node.js 18+"
    exit 1
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Make sure your Scala/Pekko backend is running on http://localhost:8080"
echo "2. Run 'wails dev' to start the client in development mode"
echo "3. Run 'wails build' to build for production"
echo ""

