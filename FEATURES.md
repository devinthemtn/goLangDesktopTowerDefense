# Tower Defense Game - Features & Implementation Guide

## 🎮 Game Overview

This is a fully functional tower defense game built in Go using the Ebitengine library. The game features a complete gameplay loop with enemies, towers, projectiles, waves, economy, and configuration system.

## ✨ Core Features

### 🏗️ Tower System
- **Basic Tower**: Balanced damage and range, good for early game
  - Cost: $50 (configurable)
  - Damage: 20 (configurable)
  - Range: 80 pixels (configurable)
  - Fire Rate: 1 shot/second (configurable)

- **Heavy Tower**: High damage, shorter range, slower fire rate
  - Cost: $100 (configurable)
  - Damage: 50 (configurable)
  - Range: 60 pixels (configurable)
  - Fire Rate: 0.5 shots/second (configurable)

### 👾 Enemy System
- **Dynamic Health Scaling**: Enemy health increases with each wave
- **Pathfinding**: Enemies follow predefined waypoints
- **Visual Health Bars**: Real-time health indication (toggleable)
- **Reward System**: Players earn money for each defeated enemy

### 🌊 Wave Progression
- **Scaling Difficulty**: Each wave spawns more enemies with higher health
- **Wave Bonuses**: Extra money awarded between waves
- **Dynamic Spawning**: Configurable spawn rates and enemy counts

### 💰 Economy System
- **Starting Resources**: Configurable starting money and lives
- **Tower Investment**: Strategic spending on different tower types
- **Resource Management**: Balance between saving and spending

### 🎯 Combat Mechanics
- **Auto-Targeting**: Towers automatically target nearest enemies in range
- **Projectile Physics**: Realistic projectile movement and collision
- **Range Visualization**: Visual indicators for tower placement (toggleable)

## 🛠️ Technical Features

### ⚙️ Configuration System
The game includes a comprehensive JSON-based configuration system (`config.json`) that allows customization of:

#### Display Settings
- Window dimensions and title
- Fullscreen and VSync options
- Visual toggles (range indicators, health bars, FPS counter)

#### Gameplay Balance
- Starting resources (money, lives)
- Enemy properties (health, speed, rewards)
- Tower statistics (damage, range, cost, fire rate)
- Wave progression parameters

#### Controls and Input
- Customizable key bindings
- Mouse-based tower placement
- Keyboard tower selection

#### Debug Features
- Debug mode toggles
- Path visualization
- Collision debugging
- God mode (infinite resources)

### 🏗️ Modular Architecture
- **Separation of Concerns**: Game logic, configuration, and rendering are separated
- **Extensible Design**: Easy to add new tower types or enemy varieties
- **Configuration-Driven**: Most game parameters can be modified without code changes

### 🎨 Visual System
- **Vector Graphics**: Smooth, scalable graphics using Ebitengine's vector package
- **Adaptive UI**: Interface scales with window size
- **Color-Coded Elements**: Intuitive visual design with distinct colors for different elements
- **Real-Time Feedback**: Dynamic health bars, range indicators, and projectile trails

## 🎮 Gameplay Mechanics

### Tower Placement Rules
1. **Path Blocking Prevention**: Cannot place towers on enemy paths
2. **Collision Detection**: Cannot place towers on existing towers
3. **Grid-Based System**: Towers snap to grid positions for clean placement
4. **Cost Validation**: Must have sufficient funds to place towers

### Combat Resolution
1. **Target Acquisition**: Towers scan for enemies within range
2. **Projectile Creation**: Towers fire projectiles toward targets
3. **Hit Detection**: Projectiles track and hit moving enemies
4. **Damage Application**: Health reduction and death detection

### Win/Loss Conditions
- **Loss**: All lives are depleted (enemies reach the end)
- **Progression**: Survive increasingly difficult waves
- **Score**: Based on waves survived and efficiency

## 🔧 Build System

### Multiple Build Options
1. **Make System**: `make build`, `make run`, `make clean`
2. **Shell Script**: `./build.sh` for automated building
3. **Direct Go Build**: `go build -o tower-defense *.go`

### Dependencies
- **Go 1.21+**: Modern Go version required
- **Ebitengine v2.6.3**: 2D game engine
- **System Libraries**: X11 development packages for Linux

## 🚀 Getting Started

### Quick Start
```bash
# Clone/download the project
cd golangTowerDefense

# Build and run
make run

# Or use the build script
./build.sh && ./tower-defense

# Or build manually
go build -o tower-defense *.go
./tower-defense
```

### Configuration Customization
1. Run the game once to generate `config.json`
2. Edit the configuration file to your preferences
3. Restart the game to apply changes

Example configuration changes:
```json
{
  "starting_money": 500,
  "show_fps": true,
  "basic_tower_damage": 30,
  "enemy_speed": 0.5
}
```

## 🎯 Strategic Gameplay Tips

### Early Game Strategy
- Focus on Basic Towers for cost efficiency
- Place towers at the beginning of the path for maximum damage time
- Save money for stronger towers when waves get harder

### Advanced Tactics
- Use Heavy Towers at chokepoints where enemies cluster
- Balance your economy - don't spend all money immediately
- Consider tower placement to maximize coverage overlap

### Resource Management
- Each enemy defeated gives $10 (configurable)
- Wave completion bonuses provide significant income
- Plan tower upgrades based on upcoming wave difficulty

## 🔮 Future Enhancement Ideas

### Gameplay Expansions
- **Tower Upgrades**: Multi-level tower enhancement system
- **Special Abilities**: Active powers like freeze or area damage
- **Enemy Varieties**: Different enemy types with unique abilities
- **Multiple Paths**: Complex maps with branching routes

### Technical Improvements
- **Sound System**: Audio effects and background music
- **Particle Effects**: Visual explosions and spell effects
- **Save/Load**: Game state persistence
- **Multiplayer**: Cooperative or competitive modes

### Quality of Life
- **Pause Functionality**: Game pause and resume
- **Speed Controls**: Fast-forward for experienced players
- **Statistics Tracking**: Detailed performance metrics
- **Achievement System**: Goals and unlockables

## 📊 Performance Characteristics

### Optimized Rendering
- **60 FPS Target**: Smooth gameplay experience
- **Efficient Draw Calls**: Minimal graphics overhead
- **Memory Management**: Clean object lifecycle management

### Scalability
- **Configurable Grid Size**: Adapts to different screen resolutions
- **Dynamic Pathfinding**: Efficient enemy movement algorithms
- **Collision Optimization**: Fast spatial queries for targeting

This tower defense game provides a solid foundation for both playing and learning game development in Go. The modular architecture and comprehensive configuration system make it easy to experiment with different gameplay mechanics and balance changes.