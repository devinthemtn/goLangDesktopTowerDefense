# Spacebar Next Wave Feature

## Overview
The spacebar feature allows players to immediately advance to the next wave after completing the current wave, earning bonus money based on how quickly they completed it.

## How It Works

### Activation
- **Key**: Press `SPACE` (spacebar)
- **When**: Only available after all enemies in the current wave are killed
- **Visual Cue**: UI shows "Press SPACE for next wave (BONUS!)"

### Requirements
To use the spacebar feature, both conditions must be met:
1. ✅ All enemies for the wave must be spawned (`enemiesSpawned >= enemiesPerWave`)
2. ✅ All enemies must be killed (`len(enemies) == 0`)

### Bonus Calculation

The bonus system rewards players for efficient wave completion:

#### Formula
```
Expected Duration = (Enemies × Spawn Delay) + (Enemies × 2s) + 5s buffer
Time Saved = Expected Duration - Actual Duration
Bonus Percentage = (Time Saved / Expected Duration) × 0.75
Bonus Money = Wave Bonus × Bonus Percentage
```

#### Bonus Range
- **Minimum**: $25 (even if not faster than expected)
- **Maximum**: 75% of the wave bonus
- **Typical Range**: 10-75% of wave bonus

#### Example
- Wave 1 with 3 enemies, spawn delay 0.5s, wave bonus $100
- Expected time: `(3 × 0.5) + (3 × 2) + 5 = 12.5 seconds`
- If completed in 8 seconds: `4.5s saved = 36% faster = $27 bonus`

### Visual Feedback

#### UI Changes
- **Wave Status**: Shows "Press SPACE for next wave (BONUS!)" when ready
- **Bonus Display**: "🎉 EARLY WAVE BONUS: +$X!" appears for 3 seconds
- **Particle Effects**: Celebratory explosion at screen center
- **Instructions**: Permanent reminder "Press SPACE when wave complete for bonus money!"

#### Debug Output (when enabled)
```
Next wave requested via spacebar!
*** EARLY WAVE COMPLETION! Bonus: $27 ***
Early completion: 4.5s saved (12.5s expected), 36.0% bonus, $27 earned
*** EARLY COMPLETION BONUS: $27 ***
```

## Gameplay Strategy

### When to Use
- **Always use it** - minimum $25 bonus even if not faster
- **Maximize efficiency** by placing towers strategically
- **Clear waves quickly** with powerful tower combinations

### Optimization Tips
1. **Place towers early** in the enemy path
2. **Use multiple towers** for faster kills
3. **Upgrade to stronger towers** as you earn more money
4. **Don't wait** - press spacebar immediately when last enemy dies

### Risk vs Reward
- **No Risk**: Spacebar is always beneficial (minimum bonus)
- **High Reward**: Up to 75% bonus for very fast completion
- **Strategic Depth**: Balances tower investment vs speed completion

## Technical Implementation

### State Management
```go
type Game struct {
    waveStartTime     float64  // Tracks wave duration
    nextWaveRequested bool     // Spacebar pressed flag
    spacePressed      bool     // Key state management
    lastBonusEarned   int      // For display purposes
    bonusDisplayTimer float64  // Bonus message timing
}
```

### Auto-Advance Fallback
If spacebar isn't pressed, waves auto-advance after 2 seconds to prevent getting stuck.

### Bonus Persistence
The bonus amount is displayed for 3 seconds and added immediately to the player's money.

## Configuration Compatibility

Works with all configuration modes:
- **Normal Mode**: Standard wave bonuses and timing
- **Debug Mode**: Faster testing with reduced spawn delays
- **Endless Mode**: Scaling bonuses with increasing difficulty

## Testing

### Debug Configuration
Use `debug-wave-config.json` for easier testing:
- Only 3 enemies per wave
- 0.5s spawn delay
- More starting money for tower placement
- Enhanced debug output

### Test Commands
```bash
# Normal game
go run .

# Debug mode testing  
go run . debug-wave-config.json

# Build and run
go build -o tower-defense .
./tower-defense debug-wave-config.json
```

## Player Benefits

1. **Increased Income**: Extra money for efficient play
2. **Faster Progression**: Skip waiting time between waves
3. **Strategic Depth**: Rewards optimal tower placement
4. **Player Control**: Choose your own pacing
5. **Risk-Free**: Always beneficial to use

## Future Enhancements

Potential improvements:
- Sound effects for bonus earning
- Streak bonuses for consecutive fast completions
- Different bonus tiers with visual indicators
- Leaderboards for fastest wave completions
- Achievement system for bonus milestones

## Code Files Modified

- `main.go`: Spacebar input handling, UI updates, bonus display
- `gamemode.go`: Bonus calculation, wave advancement logic
- Game state variables for tracking timing and bonuses

## Usage Summary

**Simple**: Kill all enemies → Press SPACE → Get bonus money → Next wave starts!

**Advanced**: Optimize tower placement → Kill enemies quickly → Press SPACE immediately → Maximize bonus → Repeat for higher income!