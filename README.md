# Tower Defense Game in Go

A simple tower defense game built with Go using the Ebitengine library.

## Features

- **Path-based enemy movement**: Enemies follow a predefined path
- **Multiple tower types**: Basic and Heavy towers with different stats
- **Wave system**: Progressively harder waves of enemies
- **Economy system**: Earn money by defeating enemies, spend it on towers
- **Real-time combat**: Towers automatically target and shoot nearby enemies
- **Health system**: Both enemies and player lives
- **Visual feedback**: Health bars, range indicators, and projectiles

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

### Objective
Defend your base by preventing enemies from reaching the end of the path. You lose lives when enemies reach the end, and the game ends when you run out of lives.

### Controls

- **Mouse Click**: Place a tower at the clicked grid position
- **Key 1**: Select Basic Tower ($50)
- **Key 2**: Select Heavy Tower ($100)

### Tower Types

1. **Basic Tower** ($50)
   - Damage: 20
   - Range: 80 pixels
   - Fire Rate: 1 shot per second
   - Good for early waves and crowd control

2. **Heavy Tower** ($100)
   - Damage: 50
   - Range: 60 pixels
   - Fire Rate: 2 shots per second
   - High damage but shorter range

### Game Mechanics

- **Starting Resources**: $100, 10 lives
- **Enemy Rewards**: $10 per enemy defeated
- **Wave Bonus**: $50 per completed wave
- **Enemy Scaling**: Each wave has stronger enemies with more health
- **Tower Placement**: Cannot place towers on the path or on existing towers

### Strategy Tips

1. **Early Game**: Place Basic Towers near the beginning of the path for maximum damage time
2. **Chokepoints**: Use Heavy Towers at curves where enemies bunch up
3. **Economy**: Balance spending on towers vs. saving for stronger towers
4. **Range Management**: Use the white range circles to optimize tower placement

### Game Elements

- **Green Background**: Grass/terrain
- **Brown Path**: Enemy walking path
- **Gray Circles**: Towers with white range indicators
- **Red Circles**: Enemies with green/red health bars
- **Yellow Dots**: Projectiles fired by towers

## Code Structure

- `main.go`: Complete game implementation including:
  - Game state management
  - Enemy spawning and movement
  - Tower placement and targeting
  - Projectile physics
  - UI rendering
  - Input handling

## Dependencies

- `github.com/hajimehoshi/ebiten/v2`: 2D game engine for Go
- Standard Go libraries for math and graphics

## Future Enhancements

Potential improvements you could add:

- More tower types (splash damage, slowing, etc.)
- Different enemy types with special abilities
- Upgradeable towers
- Sound effects and music
- More complex maps with multiple paths
- Save/load game state
- High score system
- Particle effects for explosions

## Technical Details

- **Resolution**: 800x600 pixels
- **Grid Size**: 40x40 pixel cells (20x15 grid)
- **Frame Rate**: 60 FPS
- **Coordinate System**: Top-left origin

Enjoy defending your tower!