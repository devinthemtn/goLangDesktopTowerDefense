#!/bin/bash

# Tower Defense Game - Enhanced Graphics Demo Script
# This script demonstrates the enhanced graphics features

echo "🎮 Tower Defense - Enhanced Graphics & Game Modes Demo"
echo "====================================================="
echo

# Check if game is built
if [ ! -f "./tower-defense" ]; then
    echo "Building game with enhanced graphics..."
    go build -o tower-defense *.go
    if [ $? -ne 0 ]; then
        echo "❌ Build failed! Please check the error messages above."
        exit 1
    fi
    echo "✅ Build successful!"
    echo
fi

# Create demo configuration with enhanced graphics
echo "📝 Creating demo configuration with enhanced graphics..."
cat > demo-config.json << 'EOF'
{
  "window_width": 1024,
  "window_height": 768,
  "window_title": "Tower Defense - Graphics & Modes Demo",
  "fullscreen": false,
  "vsync": true,
  "starting_money": 500,
  "starting_lives": 20,
  "enemy_speed": 0.8,
  "spawn_delay": 1.5,
  "basic_tower_cost": 40,
  "basic_tower_damage": 25,
  "basic_tower_range": 90,
  "basic_tower_fire_rate": 1.2,
  "heavy_tower_cost": 80,
  "heavy_tower_damage": 60,
  "heavy_tower_range": 70,
  "heavy_tower_fire_rate": 0.6,
  "base_enemy_health": 40,
  "health_per_wave": 8,
  "enemy_reward": 12,
  "wave_bonus": 60,
  "enemies_per_wave": 4,
  "show_range": true,
  "show_health_bars": true,
  "show_fps": true,
  "grid_size": 48,
  "master_volume": 1,
  "sfx_volume": 0.8,
  "music_volume": 0.6,
  "mute_audio": false,
  "pause_key": "Space",
  "restart_key": "R",
  "tower_select_1_key": "1",
  "tower_select_2_key": "2",
  "debug_mode": false,
  "show_path_points": false,
  "show_collision": false,
  "god_mode": false
}
EOF

echo "✨ Demo configuration created with:"
echo "   • Larger resolution (1024x768)"
echo "   • Enhanced starting resources for quick tower testing"
echo "   • FPS counter enabled"
echo "   • Balanced gameplay for both Normal and Endless modes"
echo "   • All visual effects enabled"
echo

echo "🎮 Game Features to Experience:"
echo "==============================="
echo
echo "🎯 Game Modes:"
echo "   • Normal Mode: 10-level campaign with progressive difficulty"
echo "   • Endless Mode: Infinite waves with exponential scaling"
echo "   • Mode Selection: Navigate with ↑/↓, select with ENTER"
echo "   • In-game controls: ESC to pause, M for menu, R to restart"
echo
echo "🎨 Enhanced Graphics:"
echo
echo "🏞️  Background & Textures:"
echo "   • Procedural grass texture with individual blades"
echo "   • Stone path with embedded rocks and dirt"
echo "   • Random water tiles for visual variety"
echo "   • Seamless texture tiling across the map"
echo
echo "🏗️  Tower Graphics (6 Unique Types):"
echo "   • Basic Tower (Key 1): Rotating cannon with metallic shine"
echo "   • Heavy Tower (Key 2): Pulsing energy core with triple barrels"
echo "   • Sniper Tower (Key 3): Elevated platform with long-range scope"
echo "   • Laser Tower (Key 4): Crystalline structure with spinning emitters"
echo "   • Splash Tower (Key 5): Heavy mortar with explosive shell loading"
echo "   • Slow Tower (Key 6): Ice crystals with freezing wave effects"
echo
echo "👾  Enemy Animations:"
echo "   • 6-frame walking animation cycle"
echo "   • Breathing/pulsing size effects"
echo "   • Health-based color changes (green → yellow → red)"
echo "   • Drop shadows beneath enemies"
echo "   • Glowing yellow eyes"
echo
echo "🎆  Particle Effects:"
echo "   • Explosion bursts when enemies die"
echo "   • Impact effects when projectiles hit"
echo "   • Movement dust trails behind walking enemies"
echo "   • Glowing spark trails behind projectiles"
echo "   • Physics simulation with gravity and fade-out"
echo
echo "💚  Enhanced UI:"
echo "   • Gradient health bars with smooth color transitions"
echo "   • Professional borders and shine effects"
echo "   • FPS counter in top-left corner"
echo "   • Real-time enemy spawn tracking"
echo
echo "🎯  Projectile Effects:"
echo "   • Glowing yellow projectile cores"
echo "   • Particle spark trails"
echo "   • Impact explosion effects on contact"
echo "   • Variable projectile sizes by tower type"
echo

echo "🎮 How to Experience the Demo:"
echo "============================="
echo "📋 Game Mode Testing:"
echo "   1. At menu: Try both Normal Mode and Endless Mode"
echo "   2. Normal Mode: Experience structured 10-level campaign"
echo "   3. Endless Mode: Test infinite scaling difficulty"
echo "   4. Use ESC to pause and switch between modes"
echo
echo "🏗️ Tower System Testing:"
echo "   5. Start with Basic Towers (Key 1) - rotating cannon animations"
echo "   6. Try Heavy Towers (Key 2) - pulsing energy cores with armor"
echo "   7. Test Sniper Towers (Key 3) - long-range scoped precision"
echo "   8. Experiment with Laser Towers (Key 4) - rapid-fire crystals"
echo "   9. Use Splash Towers (Key 5) - explosive area damage mortars"
echo "   10. Deploy Slow Towers (Key 6) - ice crystals that freeze enemies"
echo
echo "🎨 Visual Effects Testing:"
echo "   11. Watch unique projectile types and explosion effects"
echo "   12. Notice enemy speed changes from slow tower effects"
echo "   13. Observe mode-specific UI and progression systems"
echo "   14. Check the FPS counter showing smooth 60fps performance"
echo

echo "🚀 Starting Enhanced Graphics Demo..."
echo "Press Ctrl+C to exit the game when done exploring."
echo

# Backup original config if it exists
if [ -f "config.json" ]; then
    echo "💾 Backing up original config.json..."
    cp config.json config-backup.json
fi

# Use demo config
cp demo-config.json config.json

echo "🎬 Launching game with enhanced graphics and dual modes..."
echo "   Resolution: 1024x768"
echo "   Game Modes: Normal Campaign + Endless Survival"
echo "   Enhanced effects: ALL ENABLED"
echo "   Starting money: $500 (for quick tower placement)"
echo "   FPS display: ON"
echo
echo "Enjoy the visual experience! 🌟"
echo

# Run the game
./tower-defense

# Restore original config if backup exists
if [ -f "config-backup.json" ]; then
    echo
    echo "🔄 Restoring original configuration..."
    mv config-backup.json config.json
else
    echo
    echo "ℹ️  Demo configuration left in place."
    echo "   Delete config.json to regenerate defaults, or"
    echo "   Edit it to customize your graphics preferences."
fi

# Cleanup demo config
rm -f demo-config.json

echo
echo "🎉 Thanks for trying the Enhanced Graphics & Game Modes Demo!"
echo "   Check out MODES.md for detailed game mode strategies"
echo "   See TOWERS.md for complete tower guide and tactics"
echo "   Check GRAPHICS.md for technical visual details"
echo "   Modify config.json to customize gameplay and effects"
echo "   See README.md for complete gameplay instructions"
