#!/bin/bash

# Cleanup script - Remove existing Golang implementation

echo "🧹 Cleaning up existing implementation..."
echo ""

# Remove Go source files
echo "Removing Go source files..."
rm -rf app/
rm -f main.go

# Remove frontend files
echo "Removing frontend files..."
rm -rf frontend/

# Remove Wails config (we'll create new one)
echo "Removing old Wails config..."
rm -f wails.json

# Keep documentation and config files
echo "Keeping documentation and config files..."

# Clean go.mod but keep it
echo "Resetting go.mod..."
cat > go.mod << 'EOF'
module mnemoCast-client

go 1.22.0

toolchain go1.22.2
EOF

echo ""
echo "✅ Cleanup complete!"
echo ""
echo "Remaining files:"
ls -la | grep -E "\.(md|sh|mod|gitignore)$|^[A-Z]"
echo ""
echo "Ready to start fresh implementation!"

