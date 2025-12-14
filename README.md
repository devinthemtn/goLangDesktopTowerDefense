# Tower Defense Game in Go

A visually stunning tower defense game built with Go using the Ebitengine library, featuring **2 game modes**, 6 unique tower types, enhanced graphics, animations, particle effects, and **spacebar wave acceleration** with bonus rewards.

## 🎮 Game Modes

### 🎯 **Normal Mode (Campaign)**
- **10 Progressive Levels** with unique objectives and challenges
- **Story-driven progression** with increasing difficulty
- **Level-specific bonuses** and starting resources
- **Victory condition**: Complete all 10 levels
- **Strategic planning**: Each level requires different tower strategies

### ♾️ **Endless Mode** 
- **Infinite waves** of increasingly difficult enemies
- **Exponential scaling**: 15% difficulty increase per wave
- **Survival challenge**: How long can you last?
- **Dynamic scaling**: Enemy health, speed, and count all increase
- **High score competition**: Track your best wave performance

## 🚀 **NEW: Spacebar Wave Acceleration**
- **Press SPACE** to immediately start the next wave after clearing enemies
- **Earn bonus money** based on how quickly you complete waves
- **Minimum $25 bonus** guaranteed, up to 75% of wave bonus for fast completion
- **Strategic depth**: Balance tower investment vs speed completion
- **Risk-free feature**: Always beneficial to use!

## ✨ Enhanced Graphics Features

- **🎨 Procedural Textures**: Dynamic grass, stone paths, and water tiles
- **🏗️ Animated Towers**: Rotating cannons, pulsing effects, and muzzle flashes
- **👾 Enemy Animations**: Walking cycles, breathing effects, and damage visualization
- **🎆 Particle Systems**: Explosions, trails, and impact effects
- **💚 Enhanced UI**: Gradient health bars, visual feedback, and smooth animations
- **🌟 Visual Effects**: Glowing projectiles, shadows, and range indicators

## 🎮 Core Gameplay Features

- **Dual Game Modes**: Choose between Campaign progression or Endless survival
- **Path-based enemy movement**: Enemies follow a predefined path with visual trails
- **6 Unique Tower Types**: Each with distinct visuals, abilities, and strategic purposes
- **Dynamic difficulty**: Mode-specific scaling and progression systems
- **Economy system**: Earn money by defeating enemies, spend it on towers
- **Real-time combat**: Towers automatically target with muzzle flash effects
- **Health system**: Visual health bars with color-coded damage states
- **Particle effects**: Explosions on death and impact effects on hits

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone or download this project
3. Navigate to the project directory
4. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Game

```bash
go run main.go
```

## How to Play

### Game Start
When you launch the game, you'll see a **Mode Selection Menu**:
- **↑/↓ Keys**: Navigate between options
- **Enter/Space**: Select game mode
- **ESC**: Exit game

### Game Modes

#### 🎯 Normal Mode Objectives
- **Complete 10 campaign levels** with unique challenges
- **Progress through difficulties**: Each level has specific enemy counts and stats  
- **Level completion**: Defeat all enemies in each level to advance
- **Victory condition**: Complete all 10 levels successfully

#### ♾️ Endless Mode Objectives  
- **Survive infinite waves** of increasingly difficult enemies
- **Track your progress**: See how many waves you can complete
- **Exponential scaling**: Each wave becomes 15% more difficult
- **Challenge yourself**: Compete for highest wave count

### Controls

- **Mouse Click**: Place a tower at the clicked grid position
- **Keys 1-6**: Select different tower types (see Tower Types below)
- **SPACE**: Send next wave immediately (when wave complete) - **EARNS BONUS MONEY!**
- **ESC/P**: Pause game (during gameplay)
- **M**: Return to main menu (when paused)
- **R**: Restart current mode (on game over)

### Tower Types

**Key 1 - Basic Tower** ($50)
   - Damage: 20 | Range: 80px | Rate: 1.0/sec
   - **Rotating cannon** with metallic shine and stone base
   - Good for early waves and general defense

**Key 2 - Heavy Tower** ($100)
   - Damage: 50 | Range: 60px | Rate: 0.5/sec  
   - **Pulsing energy core** with triple-barrel design
   - High damage, shorter range, armored appearance

**Key 3 - Sniper Tower** ($150)
   - Damage: 100 | Range: 150px | Rate: 0.3/sec
   - **Long-range scope** with precision targeting
   - Elevated platform, slow but devastating shots

**Key 4 - Laser Tower** ($200)
   - Damage: 15 | Range: 70px | Rate: 3.0/sec
   - **Rapid-fire energy beams** with crystal emitters
   - Fast continuous damage with blue energy effects

**Key 5 - Splash Tower** ($180)
   - Damage: 40 | Range: 65px | Rate: 0.8/sec
   - **Area damage mortar** with explosive shells
   - Damages multiple enemies in splash radius

**Key 6 - Slow Tower** ($120)
   - Damage: 10 | Range: 90px | Rate: 1.5/sec
   - **Ice crystals** that slow enemy movement
   - Reduces enemy speed by 50% for crowd control

### Game Mechanics

- **Starting Resources**: $100, 10 lives
- **Enemy Rewards**: $10 per enemy defeated
- **Wave Bonus**: $50 per completed wave
- **Enemy Scaling**: Each wave has stronger enemies with more health
- **Tower Placement**: Cannot place towers on the path or on existing towers

### Strategy Tips

#### 🎯 Normal Mode Strategy
1. **Level Planning**: Each level has preset enemy counts - plan accordingly
2. **Resource Management**: Starting money varies per level - spend wisely
3. **Progressive Difficulty**: Later levels need stronger tower combinations
4. **Completion Focus**: Eliminate all enemies to advance to next level

#### ♾️ Endless Mode Strategy  
1. **Early Investment**: Build strong foundation with basic towers
2. **Scaling Preparation**: Save money for expensive towers as waves get harder
3. **Efficiency Focus**: Maximize damage per dollar spent
4. **Survival Mindset**: Plan for exponentially increasing difficulty

#### 🏗️ Universal Tower Tips
1. **Early Game**: Start with Basic Towers for cost-effective general defense
2. **Long Range**: Use Sniper Towers at path corners for maximum damage time  
3. **Crowd Control**: Place Slow Towers at the start to reduce enemy speed
4. **Area Denial**: Use Splash Towers where enemies cluster together
5. **Continuous DPS**: Laser Towers excel against multiple weak enemies
6. **Heavy Assault**: Heavy Towers for high single-target damage
7. **Economy**: Balance tower costs - expensive towers need good positioning

### Enhanced Visual Elements

- **Textured Background**: Static procedural grass and stone path textures
- **6 Unique Tower Designs**: 
  - **Basic**: Rotating cannon with metallic shine
  - **Heavy**: Pulsing energy core with armor plating  
  - **Sniper**: Elevated platform with long barrel and scope
  - **Laser**: Crystalline structure with spinning energy emitters
  - **Splash**: Heavy mortar with explosive shell loading
  - **Slow**: Ice crystal spikes with freezing wave effects
- **Animated Enemies**: Walking cycles, breathing, and damage-based color changes
- **Smart Projectiles**: Different visual styles based on tower type
- **Particle Effects**: Controlled explosions, trails, and visual feedback
- **Professional UI**: Gradient health bars and tower selection display

## Code Structure

- `main.go`: Core game logic and gameplay systems
- `gamemode.go`: Game mode system with:
  - Mode selection menu and navigation
  - Normal mode level progression (10 levels)
  - Endless mode infinite scaling
  - Game state management (menu, playing, paused, game over)
  - Level data and difficulty curves
- `config.go`: Comprehensive configuration system with JSON support  
- `graphics.go`: Enhanced graphics system with:
  - Procedural texture generation
  - Sprite animation system
  - Particle effects and physics
  - Visual enhancement rendering
  - Advanced drawing utilities

## Dependencies

- `github.com/hajimehoshi/ebiten/v2`: 2D game engine with hardware acceleration
- Standard Go libraries for math, graphics, and JSON configuration
- Vector graphics support for scalable visual effects

## Future Enhancements

Potential improvements you could add:

- **More Game Modes**: Time attack, survival challenges, puzzle levels
- **Advanced Graphics**: Shader effects, dynamic lighting, weather systems  
- **Tower Upgrades**: Multi-level enhancement system for existing towers
- **Enemy Varieties**: Flying enemies, armored units, boss enemies, enemy abilities
- **Campaign Features**: Story elements, cutscenes, character progression
- **Audio system**: Sound effects, music, and spatial audio
- **Map editor**: Custom level creation tools
- **Multiplayer**: Cooperative and competitive modes
- **Advanced effects**: Screen shake, bloom, post-processing
- **Achievements**: Unlock system for completing challenges

## Technical Details

- **Resolution**: Configurable (default 800x600 pixels)
- **Grid Size**: Configurable texture resolution (default 40x40 pixels)
- **Frame Rate**: 60 FPS with VSync support
- **Graphics**: Hardware-accelerated vector rendering with particle systems
- **Textures**: Procedurally generated at runtime for variety
- **Animation**: Frame-based sprite animation with configurable timing

## 📖 Additional Documentation

- **TOWERS.md**: Complete tower guide with strategies and stats
- **GRAPHICS.md**: Detailed graphics system documentation  
- **FEATURES.md**: Comprehensive feature overview
- **SHOWCASE.md**: Complete visual transformation showcase
- **config.json**: Runtime configuration options
- **demo.sh**: Interactive graphics demonstration script

Enjoy defending your tower with enhanced visuals! 🎮✨