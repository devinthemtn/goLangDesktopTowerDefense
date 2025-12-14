#!/bin/bash

# Tower Defense Game - Wave Progress Debug Script
# This script helps debug wave progression issues

echo "🐛 Tower Defense - Wave Progression Debug"
echo "========================================"
echo

# Check if game is built
if [ ! -f "./tower-defense" ]; then
    echo "Building game for debug..."
    go build -o tower-defense *.go
    if [ $? -ne 0 ]; then
        echo "❌ Build failed! Please check the error messages above."
        exit 1
    fi
    echo "✅ Build successful!"
    echo
fi

# Backup original config if it exists
if [ -f "config.json" ]; then
    echo "💾 Backing up original config.json..."
    cp config.json config-backup.json
fi

# Use debug config for wave testing
echo "📝 Setting up debug configuration..."
cp debug-wave-config.json config.json

echo "🎮 Debug Configuration Features:"
echo "   • Debug mode enabled (console output)"
echo "   • Fast enemy spawning (0.5s delay)"
echo "   • Fast enemy movement (2x speed)"
echo "   • Powerful towers for quick kills"
echo "   • Reduced enemy health for faster testing"
echo "   • Only 3 enemies per wave for quick completion"
echo

echo "🔍 What to Look For:"
echo "=============================="
echo "1. Console output showing enemy spawning"
echo "2. Wave completion detection messages"
echo "3. Level/wave advancement notifications"
echo "4. Enemy kill count tracking"
echo "5. UI showing wave number progression"
echo

echo "🎯 Testing Instructions:"
echo "========================"
echo "1. Try both Normal Mode and Endless Mode"
echo "2. Place a few Basic Towers (Key 1) near the start"
echo "3. Watch console output for debug messages"
echo "4. Verify waves advance after all enemies are killed"
echo "5. Check that wave counter increases in UI"
echo

echo "📊 Expected Console Output:"
echo "Spawning enemy 1/3 for wave 1"
echo "Spawning enemy 2/3 for wave 1"
echo "Spawning enemy 3/3 for wave 1"
echo "Enemy killed! Remaining: 2"
echo "Enemy killed! Remaining: 1"
echo "Enemy killed! Remaining: 0"
echo "Wave 1 completed! Enemies: 0, Spawned: 3/3"
echo "Advancing to next wave..."
echo

echo "🚀 Starting Wave Debug Session..."
echo "Press Ctrl+C to exit when done testing"
echo

# Run the game with debug output
echo "🎬 Launching game with debug configuration..."
./tower-defense

# Restore original config if backup exists
echo
if [ -f "config-backup.json" ]; then
    echo "🔄 Restoring original configuration..."
    mv config-backup.json config.json
else
    echo "ℹ️  Debug configuration left in place."
    echo "   Delete config.json to regenerate defaults, or"
    echo "   Edit it to customize your debug settings."
fi

# Cleanup debug config
rm -f debug-wave-config.json

echo
echo "🎉 Debug session completed!"
echo "If waves still don't progress, check for these issues:"
echo "   1. Enemies spawning but not being killed"
echo "   2. Wave completion logic not triggering"
echo "   3. Level info display blocking progression"
echo "   4. Enemy counter mismatch (spawned vs. expected)"
echo
echo "Report any console error messages for further debugging."
