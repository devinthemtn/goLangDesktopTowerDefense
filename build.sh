#!/bin/bash

# Tower Defense Game Build Script

echo "Building Tower Defense Game..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
REQUIRED_VERSION="1.21"

if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
    echo "Warning: Go version $GO_VERSION detected. Go 1.21+ is recommended."
fi

# Install dependencies
echo "Installing dependencies..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "Error: Failed to install dependencies"
    exit 1
fi

# Build the game
echo "Compiling game..."
go build -o tower-defense main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "Run the game with: ./tower-defense"
    echo ""
    echo "Controls:"
    echo "- Mouse click: Place tower"
    echo "- Key 1: Select Basic Tower ($50)"
    echo "- Key 2: Select Heavy Tower ($100)"
    echo ""
    echo "Note: Make sure you have X11 display available to run the game."
else
    echo "❌ Build failed!"
    echo ""
    echo "If you get X11/OpenGL related errors, try installing:"
    echo "sudo apt-get install libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev"
    exit 1
fi
