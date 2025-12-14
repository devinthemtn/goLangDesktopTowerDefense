package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// GameMode represents the different game modes available
type GameMode int

const (
	GameModeMenu GameMode = iota
	GameModeNormal
	GameModeEndless
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateGameOver
	StateVictory
	StatePaused
)

// LevelData contains information about a specific level
type LevelData struct {
	LevelNumber   int
	EnemyCount    int
	EnemyHealth   int
	EnemySpeed    float64
	SpawnDelay    float64
	StartingMoney int
	WaveBonus     int
	Description   string
	RequiredKills int
}

// GameModeManager handles game mode logic and level progression
type GameModeManager struct {
	CurrentMode       GameMode
	CurrentState      GameState
	CurrentLevel      int
	MaxLevel          int
	EndlessWave       int
	EndlessDifficulty float64
	LevelData         []LevelData
	MenuSelection     int
	MenuOptions       []string
	TransitionTimer   float64
	ShowLevelInfo     bool
	LevelInfoTimer    float64
	KeyUpPressed      bool
	KeyDownPressed    bool
	KeyEnterPressed   bool
	KeySpacePressed   bool
}

// NewGameModeManager creates a new game mode manager
func NewGameModeManager() *GameModeManager {
	gmm := &GameModeManager{
		CurrentMode:   GameModeMenu,
		CurrentState:  StateMenu,
		CurrentLevel:  1,
		MaxLevel:      10,
		MenuSelection: 0,
		MenuOptions:   []string{"Normal Mode", "Endless Mode", "Exit Game"},
		LevelData:     generateLevelData(nil),
	}
	return gmm
}

func NewGameModeManagerWithDebug(debugMode bool, config *GameConfig) *GameModeManager {
	gmm := &GameModeManager{
		CurrentMode:   GameModeMenu,
		CurrentState:  StateMenu,
		CurrentLevel:  1,
		MaxLevel:      10,
		MenuSelection: 0,
		MenuOptions:   []string{"Normal Mode", "Endless Mode", "Exit Game"},
		LevelData:     generateLevelData(config),
	}
	if debugMode {
		// Auto-start normal mode for debugging
		gmm.CurrentMode = GameModeNormal
		gmm.CurrentState = StatePlaying
		gmm.CurrentLevel = 1
		// Note: setupLevel will be called from main.go after game creation
	}
	return gmm
}

// generateLevelData creates the campaign levels for normal mode
func generateLevelData(config *GameConfig) []LevelData {
	levels := make([]LevelData, 10)

	// Use config values if available, otherwise use defaults
	baseHealth := 50
	baseEnemyCount := 5
	baseSpeed := 1.0
	baseMoney := 150

	if config != nil {
		baseHealth = config.BaseEnemyHealth
		baseEnemyCount = config.EnemiesPerWave
		baseSpeed = config.EnemySpeed
		baseMoney = config.StartingMoney
	}

	for i := 0; i < 10; i++ {
		level := i + 1
		difficultyMultiplier := 1.0 + float64(i)*0.3

		levels[i] = LevelData{
			LevelNumber:   level,
			EnemyCount:    baseEnemyCount + i*2,
			EnemyHealth:   int(float64(baseHealth) * difficultyMultiplier),
			EnemySpeed:    baseSpeed + float64(i)*0.1,
			SpawnDelay:    2.0 - float64(i)*0.1,
			StartingMoney: baseMoney + i*25,
			WaveBonus:     75 + i*25,
			RequiredKills: baseEnemyCount + i*2,
			Description:   generateLevelDescription(level),
		}

		// Ensure minimum spawn delay
		if levels[i].SpawnDelay < 0.5 {
			levels[i].SpawnDelay = 0.5
		}
	}

	return levels
}

// generateLevelDescription creates flavor text for each level
func generateLevelDescription(level int) string {
	descriptions := []string{
		"Tutorial: Basic enemy forces approach your position.",
		"Reinforcements: Enemy numbers are increasing.",
		"Advanced Scouts: Faster and tougher enemies detected.",
		"Heavy Assault: Armored units joining the attack.",
		"Coordinated Strike: Multiple enemy waves incoming.",
		"Elite Forces: Highly trained enemies with advanced gear.",
		"Siege Warfare: Massive enemy army mobilizing.",
		"Final Push: Enemy commander leads the assault.",
		"Last Stand: Overwhelming enemy forces converge.",
		"Ultimate Battle: Face the enemy's most powerful units.",
	}

	if level <= len(descriptions) {
		return descriptions[level-1]
	}
	return "Unknown threat level detected."
}

// Update handles game mode logic updates
func (gmm *GameModeManager) Update(game *Game) error {
	switch gmm.CurrentState {
	case StateMenu:
		return gmm.updateMenu(game)
	case StatePlaying:
		return gmm.updatePlaying(game)
	case StateGameOver:
		return gmm.updateGameOver(game)
	case StateVictory:
		return gmm.updateVictory(game)
	case StatePaused:
		return gmm.updatePaused(game)
	}
	return nil
}

// updateMenu handles menu navigation and mode selection
func (gmm *GameModeManager) updateMenu(game *Game) error {
	// Handle menu navigation with proper key state management
	upPressed := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	downPressed := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	enterPressed := ebiten.IsKeyPressed(ebiten.KeyEnter)
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)

	// Handle up navigation (only on key press, not hold)
	if upPressed && !gmm.KeyUpPressed {
		if gmm.MenuSelection > 0 {
			gmm.MenuSelection--
		}
	}
	gmm.KeyUpPressed = upPressed

	// Handle down navigation (only on key press, not hold)
	if downPressed && !gmm.KeyDownPressed {
		if gmm.MenuSelection < len(gmm.MenuOptions)-1 {
			gmm.MenuSelection++
		}
	}
	gmm.KeyDownPressed = downPressed

	// Handle mouse navigation
	mouseX, mouseY := ebiten.CursorPosition()
	menuY := 200
	for i := 0; i < len(gmm.MenuOptions); i++ {
		optionY := menuY + i*50
		if mouseX >= game.config.WindowWidth/2-100 && mouseX <= game.config.WindowWidth/2+100 &&
			mouseY >= optionY-10 && mouseY <= optionY+30 {
			gmm.MenuSelection = i
			break
		}
	}

	// Handle menu selection (only on key press, not hold)
	selectionMade := false
	if (enterPressed && !gmm.KeyEnterPressed) || (spacePressed && !gmm.KeySpacePressed) {
		selectionMade = true
	}
	gmm.KeyEnterPressed = enterPressed
	gmm.KeySpacePressed = spacePressed

	// Handle mouse click selection
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		selectionMade = true
	}

	if selectionMade {
		switch gmm.MenuSelection {
		case 0: // Normal Mode
			gmm.startNormalMode(game)
		case 1: // Endless Mode
			gmm.startEndlessMode(game)
		case 2: // Exit Game
			return fmt.Errorf("game exit requested")
		}
	}

	return nil
}

// updatePlaying handles gameplay state updates
func (gmm *GameModeManager) updatePlaying(game *Game) error {
	// Handle pause with proper key state management
	escPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	pPressed := ebiten.IsKeyPressed(ebiten.KeyP)

	if (escPressed || pPressed) && !(gmm.KeyUpPressed || gmm.KeyDownPressed) {
		gmm.CurrentState = StatePaused
		return nil
	}

	// Update level info timer
	if gmm.ShowLevelInfo {
		gmm.LevelInfoTimer -= 1.0 / 60.0
		if gmm.LevelInfoTimer <= 0 {
			gmm.ShowLevelInfo = false
		}
		// Allow gameplay to continue during level info display
	}

	// Check for game over
	if game.lives <= 0 {
		gmm.CurrentState = StateGameOver
		return nil
	}

	// Reduce debug output spam - only show occasionally
	if game.config.DebugMode && game.spawnTimer > 3.0 {
		fmt.Printf("Game Status: Mode=%d, State=%d, Level=%d\n",
			gmm.CurrentMode, gmm.CurrentState, gmm.CurrentLevel)
	}

	// Handle mode-specific logic
	switch gmm.CurrentMode {
	case GameModeNormal:
		return gmm.updateNormalMode(game)
	case GameModeEndless:
		return gmm.updateEndlessMode(game)
	}

	return nil
}

// updateNormalMode handles normal/campaign mode progression
func (gmm *GameModeManager) updateNormalMode(game *Game) error {
	// Only show debug output occasionally to avoid spam
	if game.config.DebugMode && game.spawnTimer > 2.0 {
		fmt.Printf("updateNormalMode: Level %d, Enemies: %d, Spawned: %d/%d\n",
			gmm.CurrentLevel, len(game.enemies), game.enemiesSpawned, game.enemiesPerWave)
	}

	// Update transition timer when wave is complete
	waveComplete := len(game.enemies) == 0 && game.enemiesSpawned >= game.enemiesPerWave
	if waveComplete {
		gmm.TransitionTimer += 1.0 / 60.0
	} else {
		gmm.TransitionTimer = 0
	}

	// Check if level should advance (via spacebar or auto after delay)
	if waveComplete && (game.nextWaveRequested || gmm.shouldAutoAdvance()) {
		// Add debug output
		if game.config.DebugMode {
			fmt.Printf("*** LEVEL COMPLETION DETECTED! ***\n")
			fmt.Printf("Level %d completed! Enemies: %d, Spawned: %d/%d\n",
				gmm.CurrentLevel, len(game.enemies), game.enemiesSpawned, game.enemiesPerWave)
		}

		if gmm.CurrentLevel >= gmm.MaxLevel {
			// Game completed!
			gmm.CurrentState = StateVictory
			if game.config.DebugMode {
				fmt.Printf("*** GAME VICTORY! ***\n")
			}
		} else {
			// Calculate early completion bonus if spacebar was used
			bonus := 0
			if game.nextWaveRequested {
				bonus = gmm.calculateEarlyCompletionBonus(game)
				if game.config.DebugMode {
					fmt.Printf("*** EARLY WAVE COMPLETION! Bonus: $%d ***\n", bonus)
				}
			}

			// Advance to next level
			if game.config.DebugMode {
				fmt.Printf("*** ADVANCING TO NEXT LEVEL! ***\n")
			}
			gmm.advanceLevel(game)

			// Apply bonus money
			if bonus > 0 {
				game.money += bonus
				game.lastBonusEarned = bonus
				game.bonusDisplayTimer = 3.0 // Show bonus for 3 seconds

				// Create celebratory particle effects
				centerX := float64(game.config.WindowWidth) / 2
				centerY := float64(game.config.WindowHeight) / 2
				game.graphics.CreateExplosion(Point{X: centerX, Y: centerY}, 5, game.config)

				if game.config.DebugMode {
					fmt.Printf("*** EARLY COMPLETION BONUS: $%d ***\n", bonus)
				}
			}
		}
	}
	return nil
}

// updateEndlessMode handles endless mode scaling difficulty
func (gmm *GameModeManager) updateEndlessMode(game *Game) error {
	// Check if wave is completed
	if len(game.enemies) == 0 && game.enemiesSpawned >= game.enemiesPerWave {
		// Add debug output
		if game.config.DebugMode {
			fmt.Printf("Wave %d completed! Enemies: %d, Spawned: %d/%d\n",
				gmm.EndlessWave, len(game.enemies), game.enemiesSpawned, game.enemiesPerWave)
		}

		gmm.EndlessWave++
		gmm.EndlessDifficulty += 0.15 // Increase difficulty by 15% each wave
		gmm.setupEndlessWave(game)
	}
	return nil
}

// updateGameOver handles game over state
func (gmm *GameModeManager) updateGameOver(game *Game) error {
	enterPressed := ebiten.IsKeyPressed(ebiten.KeyEnter)
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	rPressed := ebiten.IsKeyPressed(ebiten.KeyR)

	if (enterPressed && !gmm.KeyEnterPressed) || (spacePressed && !gmm.KeySpacePressed) {
		gmm.returnToMenu(game)
	}
	if rPressed {
		gmm.restartCurrentMode(game)
	}

	gmm.KeyEnterPressed = enterPressed
	gmm.KeySpacePressed = spacePressed
	return nil
}

// updateVictory handles victory state (normal mode completion)
func (gmm *GameModeManager) updateVictory(game *Game) error {
	enterPressed := ebiten.IsKeyPressed(ebiten.KeyEnter)
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	rPressed := ebiten.IsKeyPressed(ebiten.KeyR)

	if (enterPressed && !gmm.KeyEnterPressed) || (spacePressed && !gmm.KeySpacePressed) {
		gmm.returnToMenu(game)
	}
	if rPressed {
		gmm.startNormalMode(game) // Restart campaign
	}

	gmm.KeyEnterPressed = enterPressed
	gmm.KeySpacePressed = spacePressed
	return nil
}

// updatePaused handles pause state
func (gmm *GameModeManager) updatePaused(game *Game) error {
	// Handle resume with key state management
	escPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	pPressed := ebiten.IsKeyPressed(ebiten.KeyP)
	mPressed := ebiten.IsKeyPressed(ebiten.KeyM)

	if (escPressed || pPressed) && !(gmm.KeyUpPressed || gmm.KeyDownPressed) {
		gmm.CurrentState = StatePlaying
	}
	if mPressed && !gmm.KeyEnterPressed {
		gmm.returnToMenu(game)
	}
	return nil
}

// startNormalMode initializes normal/campaign mode
func (gmm *GameModeManager) startNormalMode(game *Game) {
	gmm.CurrentMode = GameModeNormal
	gmm.CurrentState = StatePlaying
	gmm.CurrentLevel = 1
	gmm.setupLevel(game, 1)
	gmm.ShowLevelInfo = true
	gmm.LevelInfoTimer = 1.0
}

// startEndlessMode initializes endless mode
func (gmm *GameModeManager) startEndlessMode(game *Game) {
	gmm.CurrentMode = GameModeEndless
	gmm.CurrentState = StatePlaying
	gmm.EndlessWave = 1
	gmm.EndlessDifficulty = 1.0
	gmm.setupEndlessWave(game)
	gmm.ShowLevelInfo = true
	gmm.LevelInfoTimer = 1.0
}

// setupLevel configures the game for a specific campaign level
func (gmm *GameModeManager) setupLevel(game *Game, level int) {
	if level > len(gmm.LevelData) {
		level = len(gmm.LevelData)
	}

	levelData := gmm.LevelData[level-1]

	// Reset game state
	game.enemies = []*Enemy{}
	game.projectiles = []*Projectile{}
	game.money = levelData.StartingMoney
	game.lives = game.config.StartingLives
	game.wave = level
	game.enemiesSpawned = 0
	game.enemiesPerWave = levelData.EnemyCount
	game.spawnTimer = 0
	game.gameOver = false

	// Update config for this level
	game.config.BaseEnemyHealth = levelData.EnemyHealth
	game.config.EnemySpeed = levelData.EnemySpeed
	game.config.SpawnDelay = levelData.SpawnDelay
	game.config.WaveBonus = levelData.WaveBonus
}

// setupEndlessWave configures the game for the next endless wave
func (gmm *GameModeManager) setupEndlessWave(game *Game) {
	// Reset enemies and projectiles but keep towers and money
	game.enemies = []*Enemy{}
	game.projectiles = []*Projectile{}
	game.wave = gmm.EndlessWave
	game.enemiesSpawned = 0
	game.spawnTimer = 0

	// Scale difficulty
	baseHealth := 50
	baseEnemyCount := 5
	baseSpeed := 1.0

	// Exponential scaling for endless mode
	scaledHealth := int(float64(baseHealth) * gmm.EndlessDifficulty)
	scaledEnemyCount := int(float64(baseEnemyCount) * (1.0 + float64(gmm.EndlessWave)*0.2))
	scaledSpeed := baseSpeed + float64(gmm.EndlessWave)*0.05
	scaledSpawnDelay := math.Max(0.3, 2.0-float64(gmm.EndlessWave)*0.05)

	game.enemiesPerWave = scaledEnemyCount
	game.config.BaseEnemyHealth = scaledHealth
	game.config.EnemySpeed = scaledSpeed
	game.config.SpawnDelay = scaledSpawnDelay

	// Bonus money for surviving longer
	waveBonus := 50 + gmm.EndlessWave*10
	game.money += waveBonus
	game.config.WaveBonus = waveBonus
}

// advanceLevel moves to the next campaign level
func (gmm *GameModeManager) advanceLevel(game *Game) {
	// Add debug output
	if game.config.DebugMode {
		fmt.Printf("*** ADVANCE LEVEL CALLED! ***\n")
		fmt.Printf("Advancing from level %d to %d\n", gmm.CurrentLevel, gmm.CurrentLevel+1)
	}

	oldLevel := gmm.CurrentLevel
	gmm.CurrentLevel++
	game.money += gmm.LevelData[gmm.CurrentLevel-2].WaveBonus // Previous level bonus
	gmm.setupLevel(game, gmm.CurrentLevel)
	gmm.ShowLevelInfo = true
	gmm.LevelInfoTimer = 1.0

	if game.config.DebugMode {
		fmt.Printf("*** LEVEL ADVANCED: %d -> %d, Wave: %d ***\n", oldLevel, gmm.CurrentLevel, game.wave)
	}

	// Reset wave timing and request flags
	game.waveStartTime = 0
	game.nextWaveRequested = false
	gmm.TransitionTimer = 0
}

// shouldAutoAdvance determines if wave should advance automatically
func (gmm *GameModeManager) shouldAutoAdvance() bool {
	// Auto-advance after a brief delay to allow players to see completion
	return gmm.TransitionTimer > 2.0
}

// calculateEarlyCompletionBonus calculates bonus money for early wave completion
func (gmm *GameModeManager) calculateEarlyCompletionBonus(game *Game) int {
	if !game.nextWaveRequested {
		return 0
	}

	// Calculate expected wave duration based on spawn timing and difficulty
	baseTime := float64(game.enemiesPerWave) * game.config.SpawnDelay
	killTime := float64(game.enemiesPerWave) * 2.0 // Assume 2 seconds per enemy to kill
	expectedDuration := baseTime + killTime + 5.0  // Add buffer time

	actualDuration := game.waveStartTime

	if actualDuration >= expectedDuration {
		return 25 // Minimum bonus for using spacebar even if not faster
	}

	// Calculate time saved and bonus percentage
	timeSaved := expectedDuration - actualDuration
	bonusPercentage := timeSaved / expectedDuration

	// More generous bonus scaling: 10% to 75% of wave bonus
	if bonusPercentage > 0.8 {
		bonusPercentage = 0.75
	} else if bonusPercentage < 0.1 {
		bonusPercentage = 0.1
	} else {
		bonusPercentage = bonusPercentage * 0.75 // Scale up the percentage
	}

	// Calculate bonus based on wave bonus + base reward
	baseBonus := gmm.LevelData[gmm.CurrentLevel-1].WaveBonus
	speedBonus := int(float64(baseBonus) * bonusPercentage)
	minimumBonus := 25

	// Ensure minimum bonus
	if speedBonus < minimumBonus {
		speedBonus = minimumBonus
	}

	if game.config.DebugMode {
		fmt.Printf("Early completion: %.1fs saved (%.1fs expected), %.1f%% bonus, $%d earned\n",
			timeSaved, expectedDuration, bonusPercentage*100, speedBonus)
	}

	return speedBonus
}

// returnToMenu resets to main menu
func (gmm *GameModeManager) returnToMenu(game *Game) {
	gmm.CurrentMode = GameModeMenu
	gmm.CurrentState = StateMenu
	gmm.MenuSelection = 0

	// Reset key states
	gmm.KeyUpPressed = false
	gmm.KeyDownPressed = false
	gmm.KeyEnterPressed = false
	gmm.KeySpacePressed = false

	// Reset game state
	game.enemies = []*Enemy{}
	game.towers = []*Tower{}
	game.projectiles = []*Projectile{}
	game.money = game.config.StartingMoney
	game.lives = game.config.StartingLives
	game.wave = 1
	game.enemiesSpawned = 0
	game.gameOver = false
}

// restartCurrentMode restarts the current game mode
func (gmm *GameModeManager) restartCurrentMode(game *Game) {
	switch gmm.CurrentMode {
	case GameModeNormal:
		gmm.startNormalMode(game)
	case GameModeEndless:
		gmm.startEndlessMode(game)
	}
}

// DrawMenu renders the main menu
func (gmm *GameModeManager) DrawMenu(screen *ebiten.Image, config *GameConfig) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{20, 30, 40, 255})

	// Title
	titleText := "TOWER DEFENSE"
	ebitenutil.DebugPrintAt(screen, titleText, config.WindowWidth/2-100, 100)

	// Subtitle
	subtitleText := "Choose Your Battle Mode"
	ebitenutil.DebugPrintAt(screen, subtitleText, config.WindowWidth/2-80, 140)

	// Menu options
	menuY := 200
	for i, option := range gmm.MenuOptions {
		x := config.WindowWidth/2 - 80
		y := menuY + i*50

		// Highlight selected option
		if i == gmm.MenuSelection {
			// Draw selection background
			vector.DrawFilledRect(screen, float32(x-20), float32(y-5), 220, 30, color.RGBA{50, 100, 150, 150}, false)
			vector.StrokeRect(screen, float32(x-20), float32(y-5), 220, 30, 2, color.RGBA{100, 150, 200, 255}, false)
			optionText := "► " + option + " ◄"
			ebitenutil.DebugPrintAt(screen, optionText, x-5, y)
		} else {
			ebitenutil.DebugPrintAt(screen, option, x, y)
		}
	}

	// Mode descriptions
	descY := menuY + len(gmm.MenuOptions)*50 + 40
	switch gmm.MenuSelection {
	case 0: // Normal Mode
		desc := "Campaign Mode: Complete 10 progressively challenging levels\nEach level has unique objectives and difficulty scaling\nComplete all levels to achieve victory!"
		ebitenutil.DebugPrintAt(screen, desc, 50, descY)
	case 1: // Endless Mode
		desc := "Endless Mode: Survive infinite waves of enemies\nDifficulty increases with each wave\nHow long can you survive?"
		ebitenutil.DebugPrintAt(screen, desc, 50, descY)
	case 2: // Exit
		desc := "Exit the game"
		ebitenutil.DebugPrintAt(screen, desc, 50, descY)
	}

	// Show current selection clearly
	selectionIndicator := fmt.Sprintf("Selected: %s (Press ENTER/SPACE or Click to confirm)", gmm.MenuOptions[gmm.MenuSelection])
	ebitenutil.DebugPrintAt(screen, selectionIndicator, 50, config.WindowHeight-80)

	// Controls
	controlsText := "Controls: ↑/↓ Navigate | ENTER/SPACE Select | ESC Exit"
	ebitenutil.DebugPrintAt(screen, controlsText, 50, config.WindowHeight-50)
}

// DrawGameState renders game state overlays
func (gmm *GameModeManager) DrawGameState(screen *ebiten.Image, game *Game) {
	config := game.config

	switch gmm.CurrentState {
	case StatePlaying:
		gmm.drawPlayingOverlay(screen, game)
	case StateGameOver:
		gmm.drawGameOverScreen(screen, config)
	case StateVictory:
		gmm.drawVictoryScreen(screen, config)
	case StatePaused:
		gmm.drawPausedOverlay(screen, config)
	}
}

// drawPlayingOverlay renders gameplay UI overlays
func (gmm *GameModeManager) drawPlayingOverlay(screen *ebiten.Image, game *Game) {
	// Mode-specific UI
	switch gmm.CurrentMode {
	case GameModeNormal:
		modeText := fmt.Sprintf("CAMPAIGN MODE - Level %d/%d", gmm.CurrentLevel, gmm.MaxLevel)
		ebitenutil.DebugPrintAt(screen, modeText, 10, game.config.WindowHeight-100)

		if gmm.CurrentLevel <= len(gmm.LevelData) {
			levelData := gmm.LevelData[gmm.CurrentLevel-1]
			progressText := fmt.Sprintf("Progress: %d/%d enemies",
				maxInt(0, levelData.EnemyCount-len(game.enemies)), levelData.EnemyCount)
			ebitenutil.DebugPrintAt(screen, progressText, 10, game.config.WindowHeight-80)
		}

	case GameModeEndless:
		modeText := fmt.Sprintf("ENDLESS MODE - Wave %d", gmm.EndlessWave)
		ebitenutil.DebugPrintAt(screen, modeText, 10, game.config.WindowHeight-100)

		difficultyText := fmt.Sprintf("Difficulty: %.1fx", gmm.EndlessDifficulty)
		ebitenutil.DebugPrintAt(screen, difficultyText, 10, game.config.WindowHeight-80)
	}

	// Show level info at start of level
	if gmm.ShowLevelInfo {
		gmm.drawLevelInfo(screen, game)
	}

	// Game controls
	controlsText := "ESC/P: Pause | M: Menu"
	ebitenutil.DebugPrintAt(screen, controlsText, 10, game.config.WindowHeight-40)
}

// drawLevelInfo shows level information at the start
func (gmm *GameModeManager) drawLevelInfo(screen *ebiten.Image, game *Game) {
	// Semi-transparent overlay
	vector.DrawFilledRect(screen, 0, 0, float32(game.config.WindowWidth), float32(game.config.WindowHeight),
		color.RGBA{0, 0, 0, 150}, false)

	centerX := game.config.WindowWidth / 2
	centerY := game.config.WindowHeight / 2

	switch gmm.CurrentMode {
	case GameModeNormal:
		if gmm.CurrentLevel <= len(gmm.LevelData) {
			levelData := gmm.LevelData[gmm.CurrentLevel-1]

			titleText := fmt.Sprintf("LEVEL %d", gmm.CurrentLevel)
			ebitenutil.DebugPrintAt(screen, titleText, centerX-50, centerY-80)

			ebitenutil.DebugPrintAt(screen, levelData.Description, centerX-150, centerY-40)

			statsText := fmt.Sprintf("Enemies: %d | Health: %d | Speed: %.1fx",
				levelData.EnemyCount, levelData.EnemyHealth, levelData.EnemySpeed)
			ebitenutil.DebugPrintAt(screen, statsText, centerX-120, centerY)

			bonusText := fmt.Sprintf("Starting Money: $%d | Wave Bonus: $%d",
				levelData.StartingMoney, levelData.WaveBonus)
			ebitenutil.DebugPrintAt(screen, bonusText, centerX-120, centerY+20)
		}

	case GameModeEndless:
		titleText := fmt.Sprintf("WAVE %d", gmm.EndlessWave)
		ebitenutil.DebugPrintAt(screen, titleText, centerX-50, centerY-60)

		if gmm.EndlessWave == 1 {
			descText := "Endless Mode: Survive as long as possible!"
			ebitenutil.DebugPrintAt(screen, descText, centerX-120, centerY-20)
		} else {
			diffText := fmt.Sprintf("Difficulty increased to %.1fx", gmm.EndlessDifficulty)
			ebitenutil.DebugPrintAt(screen, diffText, centerX-100, centerY-20)
		}
	}

	if gmm.LevelInfoTimer > 0.5 {
		countdownText := fmt.Sprintf("Starting in %.1f seconds...", gmm.LevelInfoTimer)
		ebitenutil.DebugPrintAt(screen, countdownText, centerX-80, centerY+60)
	} else {
		readyText := "Ready! Game starting..."
		ebitenutil.DebugPrintAt(screen, readyText, centerX-80, centerY+60)
	}
}

// drawGameOverScreen renders game over screen
func (gmm *GameModeManager) drawGameOverScreen(screen *ebiten.Image, config *GameConfig) {
	// Semi-transparent overlay
	vector.DrawFilledRect(screen, 0, 0, float32(config.WindowWidth), float32(config.WindowHeight),
		color.RGBA{0, 0, 0, 200}, false)

	centerX := config.WindowWidth / 2
	centerY := config.WindowHeight / 2

	titleText := "GAME OVER"
	ebitenutil.DebugPrintAt(screen, titleText, centerX-50, centerY-60)

	var statsText string
	switch gmm.CurrentMode {
	case GameModeNormal:
		statsText = fmt.Sprintf("Reached Level: %d/%d", gmm.CurrentLevel, gmm.MaxLevel)
	case GameModeEndless:
		statsText = fmt.Sprintf("Survived Waves: %d", gmm.EndlessWave-1)
	}
	ebitenutil.DebugPrintAt(screen, statsText, centerX-80, centerY-20)

	controlsText := "ENTER: Return to Menu | R: Restart"
	ebitenutil.DebugPrintAt(screen, controlsText, centerX-120, centerY+20)
}

// drawVictoryScreen renders victory screen (normal mode completion)
func (gmm *GameModeManager) drawVictoryScreen(screen *ebiten.Image, config *GameConfig) {
	// Semi-transparent overlay
	vector.DrawFilledRect(screen, 0, 0, float32(config.WindowWidth), float32(config.WindowHeight),
		color.RGBA{0, 100, 0, 200}, false)

	centerX := config.WindowWidth / 2
	centerY := config.WindowHeight / 2

	titleText := "VICTORY!"
	ebitenutil.DebugPrintAt(screen, titleText, centerX-40, centerY-60)

	congratsText := "Campaign Completed Successfully!"
	ebitenutil.DebugPrintAt(screen, congratsText, centerX-120, centerY-20)

	controlsText := "ENTER: Return to Menu | R: Play Again"
	ebitenutil.DebugPrintAt(screen, controlsText, centerX-120, centerY+20)
}

// drawPausedOverlay renders pause screen
func (gmm *GameModeManager) drawPausedOverlay(screen *ebiten.Image, config *GameConfig) {
	// Semi-transparent overlay
	vector.DrawFilledRect(screen, 0, 0, float32(config.WindowWidth), float32(config.WindowHeight),
		color.RGBA{0, 0, 0, 150}, false)

	centerX := config.WindowWidth / 2
	centerY := config.WindowHeight / 2

	titleText := "PAUSED"
	ebitenutil.DebugPrintAt(screen, titleText, centerX-30, centerY-20)

	controlsText := "ESC/P: Resume | M: Return to Menu"
	ebitenutil.DebugPrintAt(screen, controlsText, centerX-100, centerY+20)
}

// GetCurrentModeInfo returns information about the current game mode
func (gmm *GameModeManager) GetCurrentModeInfo() (GameMode, GameState, int) {
	switch gmm.CurrentMode {
	case GameModeNormal:
		return gmm.CurrentMode, gmm.CurrentState, gmm.CurrentLevel
	case GameModeEndless:
		return gmm.CurrentMode, gmm.CurrentState, gmm.EndlessWave
	default:
		return gmm.CurrentMode, gmm.CurrentState, 0
	}
}

// Helper function for max
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Helper function for max integers
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
