package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Point struct {
	X, Y float64
}

type Enemy struct {
	Position      Point
	Target        Point
	Health        int
	MaxHealth     int
	Speed         float64
	BaseSpeed     float64 // Original speed before slow effects
	PathIndex     int
	Alive         bool
	ReachedEnd    bool
	SlowTimer     float64 // Time remaining for slow effect
	SlowIntensity float64 // Current slow multiplier (1.0 = no slow)
}

type Tower struct {
	Position Point
	Range    float64
	Damage   int
	FireRate float64
	LastFire float64
	Cost     int
	Type     int
	Special  map[string]float64 // For special effects like splash radius, slow duration
}

type Projectile struct {
	Position    Point
	Target      *Enemy
	Speed       float64
	Damage      int
	Active      bool
	SourceTower *Tower // Track which tower fired this projectile
}

type Game struct {
	enemies           []*Enemy
	towers            []*Tower
	projectiles       []*Projectile
	path              []Point
	money             int
	lives             int
	wave              int
	spawnTimer        float64
	gameOver          bool
	selectedTowerType int
	config            *GameConfig
	enemiesSpawned    int
	enemiesPerWave    int
	graphics          *GraphicsManager
	modeManager       *GameModeManager
	waveStartTime     float64
	nextWaveRequested bool
	spacePressed      bool
	lastBonusEarned   int
	bonusDisplayTimer float64
	saveManager       *SaveManager
	audioManager      *AudioManager
	settingsMenu      *SettingsMenu
	helpScreen        *HelpScreen
	settingsKeyPressed bool // Track settings key state
	helpKeyPressed     bool // Track help key state
	mousePressed       bool // Track mouse button state for tower placement
}

func NewGame(config *GameConfig) *Game {
	cellSize := config.GridSize
	mapWidth := config.WindowWidth / cellSize
	mapHeight := config.WindowHeight / cellSize

	// Create a simple path that adapts to screen size
	path := []Point{
		{0, float64(mapHeight / 2)},
		{float64(mapWidth / 4), float64(mapHeight / 2)},
		{float64(mapWidth / 4), float64(mapHeight / 4)},
		{float64(mapWidth / 2), float64(mapHeight / 4)},
		{float64(mapWidth / 2), float64(3 * mapHeight / 4)},
		{float64(3 * mapWidth / 4), float64(3 * mapHeight / 4)},
		{float64(3 * mapWidth / 4), float64(mapHeight / 3)},
		{float64(mapWidth), float64(mapHeight / 3)},
	}

	// Initialize audio manager
	audioManager, err := NewAudioManager(48000)
	if err != nil {
		LogWarn("Audio not available: %v (continuing without audio)", err)
		audioManager = nil
	}

	game := &Game{
		enemies:           []*Enemy{},
		towers:            []*Tower{},
		projectiles:       []*Projectile{},
		path:              path,
		money:             config.StartingMoney,
		lives:             config.StartingLives,
		wave:              1,
		spawnTimer:        0,
		selectedTowerType: 1,
		config:            config,
		enemiesSpawned:    0,
		enemiesPerWave:    config.GetEnemiesInWave(1),
		graphics:          NewGraphicsManager(),
		modeManager:       NewGameModeManagerWithDebug(config.DebugMode, config),
		audioManager:      audioManager,
	}

	// Initialize settings menu
	if audioManager != nil {
		game.settingsMenu = NewSettingsMenu(config, audioManager)
		
		// Apply saved audio settings
		audioManager.SetMasterVolume(config.MasterVolume)
		audioManager.SetMusicVolume(config.MusicVolume)
		audioManager.SetSFXVolume(config.SFXVolume)
		
		if config.MuteAudio {
			audioManager.ToggleMusic()
		}
	}
	
	// Initialize help screen
	game.helpScreen = NewHelpScreen(config)

	// If debug mode auto-started playing mode, setup the first level
	if config.DebugMode && game.modeManager.CurrentState == StatePlaying {
		game.modeManager.setupLevel(game, 1)
	}

	return game
}

func (g *Game) Update() error {
	// Add error recovery for update loop
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in Update loop: %v", r)
			// Try to auto-save before crash
			if g.saveManager != nil {
				if err := g.saveManager.AutoSave(g); err != nil {
					LogError("Failed to auto-save after panic: %v", err)
				}
			}
			panic(r) // Re-panic to trigger main recovery
		}
	}()

	// Handle Q key to quit the game
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	// Update game mode system
	if err := g.modeManager.Update(g); err != nil {
		LogError("Mode manager update error: %v", err)
		return err
	}

	// Only update game logic if we're in playing state
	if g.modeManager.CurrentState != StatePlaying {
		return nil
	}

	if g.gameOver {
		return nil
	}

	// Handle spacebar for next wave (only when all enemies are dead and spawned)
	spaceCurrentlyPressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	if spaceCurrentlyPressed && !g.spacePressed && len(g.enemies) == 0 && g.enemiesSpawned >= g.enemiesPerWave {
		g.nextWaveRequested = true
		if g.config.DebugMode {
			fmt.Printf("Next wave requested via spacebar!\n")
		}
	}
	g.spacePressed = spaceCurrentlyPressed

	// Handle settings menu toggle (S key or F1) - with key state tracking
	if g.settingsMenu != nil {
		settingsKeyCurrentlyPressed := ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyF1)
		if settingsKeyCurrentlyPressed && !g.settingsKeyPressed {
			g.settingsMenu.Toggle()
		}
		g.settingsKeyPressed = settingsKeyCurrentlyPressed
		
		// Update settings menu if visible
		if g.settingsMenu.IsVisible() {
			g.settingsMenu.Update()
			return nil // Don't update game when settings menu is open
		}
	}
	
	// Handle help screen toggle (H key or F2) - with key state tracking
	if g.helpScreen != nil {
		helpKeyCurrentlyPressed := ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyF2)
		if helpKeyCurrentlyPressed && !g.helpKeyPressed {
			g.helpScreen.Toggle()
		}
		g.helpKeyPressed = helpKeyCurrentlyPressed
		
		// Update help screen if visible
		if g.helpScreen.IsVisible() {
			g.helpScreen.Update()
			return nil // Don't update game when help screen is open
		}
	}

	// Update particle system
	g.graphics.ParticleSystem.Update()

	// Spawn enemies
	g.spawnTimer += 1.0 / 60.0

	if g.spawnTimer > g.config.SpawnDelay && g.enemiesSpawned < g.enemiesPerWave {
		if g.config.DebugMode {
			fmt.Printf("Spawning enemy %d/%d for wave %d\n", g.enemiesSpawned+1, g.enemiesPerWave, g.wave)
		}
		g.spawnEnemy()
		g.spawnTimer = 0
		g.enemiesSpawned++
	}

	// Check wave completion here in main game loop as backup
	if len(g.enemies) == 0 && g.enemiesSpawned >= g.enemiesPerWave {
		if g.config.DebugMode {
			fmt.Printf("Main loop detected wave completion: enemies=%d, spawned=%d/%d\n",
				len(g.enemies), g.enemiesSpawned, g.enemiesPerWave)
			fmt.Printf("Mode: %d, State: %d, Wave: %d\n",
				g.modeManager.CurrentMode, g.modeManager.CurrentState, g.wave)
		}
	}

	// Update enemies
	for i := len(g.enemies) - 1; i >= 0; i-- {
		enemy := g.enemies[i]
		if !enemy.Alive {
			continue
		}

		g.moveEnemy(enemy)

		if enemy.ReachedEnd {
			g.lives--
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			if g.lives <= 0 {
				g.gameOver = true
			}
		} else if enemy.Health <= 0 {
			// Create explosion effect when enemy dies
			g.graphics.CreateExplosion(enemy.Position, 3, g.config)
			g.money += g.config.EnemyReward
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			if g.config.DebugMode {
				fmt.Printf("Enemy killed! Remaining: %d, Spawned: %d/%d\n", len(g.enemies)-1, g.enemiesSpawned, g.enemiesPerWave)
				if len(g.enemies)-1 == 0 && g.enemiesSpawned >= g.enemiesPerWave {
					fmt.Printf("*** WAVE SHOULD COMPLETE NOW! ***\n")
				}
			}
		}
	}

	// Update towers
	for _, tower := range g.towers {
		tower.LastFire += 1.0 / 60.0
		if tower.LastFire >= tower.FireRate {
			target := g.findNearestEnemy(tower)
			if target != nil {
				g.fireTower(tower, target)
				tower.LastFire = 0
			}
		}
	}

	// Update projectiles
	for i := len(g.projectiles) - 1; i >= 0; i-- {
		proj := g.projectiles[i]
		if !proj.Active {
			continue
		}

		g.moveProjectile(proj)

		if !proj.Active {
			g.projectiles = append(g.projectiles[:i], g.projectiles[i+1:]...)
		}
	}

	// Handle mouse input for tower placement (single click, not continuous)
	mouseCurrentlyPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if mouseCurrentlyPressed && !g.mousePressed {
		x, y := ebiten.CursorPosition()
		cellSize := g.config.GridSize
		gridX := x / cellSize
		gridY := y / cellSize
		g.placeTower(float64(gridX), float64(gridY))
	}
	g.mousePressed = mouseCurrentlyPressed

	// Handle key input for tower selection (only in playing state)
	if g.modeManager.CurrentState == StatePlaying {
		if ebiten.IsKeyPressed(ebiten.Key1) {
			g.selectedTowerType = 1
		} else if ebiten.IsKeyPressed(ebiten.Key2) {
			g.selectedTowerType = 2
		} else if ebiten.IsKeyPressed(ebiten.Key3) {
			g.selectedTowerType = 3
		} else if ebiten.IsKeyPressed(ebiten.Key4) {
			g.selectedTowerType = 4
		} else if ebiten.IsKeyPressed(ebiten.Key5) {
			g.selectedTowerType = 5
		} else if ebiten.IsKeyPressed(ebiten.Key6) {
			g.selectedTowerType = 6
		}
	}

	// Update wave timer for bonus calculation
	g.waveStartTime += 1.0 / 60.0

	return nil
}

func (g *Game) spawnEnemy() {
	if len(g.path) == 0 {
		return
	}

	cellSize := float64(g.config.GridSize)
	health := g.config.GetEnemyHealth(g.wave)

	enemy := &Enemy{
		Position:      Point{g.path[0].X*cellSize + cellSize/2, g.path[0].Y*cellSize + cellSize/2},
		Health:        health,
		MaxHealth:     health,
		Speed:         g.config.EnemySpeed,
		BaseSpeed:     g.config.EnemySpeed, // Store original speed
		PathIndex:     0,
		Alive:         true,
		SlowTimer:     0,
		SlowIntensity: 1.0, // No slow effect initially
	}

	if len(g.path) > 1 {
		enemy.Target = Point{g.path[1].X*cellSize + cellSize/2, g.path[1].Y*cellSize + cellSize/2}
	}

	g.enemies = append(g.enemies, enemy)
}

func (g *Game) moveEnemy(enemy *Enemy) {
	if enemy.PathIndex >= len(g.path)-1 {
		enemy.ReachedEnd = true
		return
	}

	// Update slow effect timer and restore speed if expired
	if enemy.SlowTimer > 0 {
		enemy.SlowTimer -= 1.0 / 60.0
		if enemy.SlowTimer <= 0 {
			enemy.SlowTimer = 0
			enemy.SlowIntensity = 1.0
			enemy.Speed = enemy.BaseSpeed // Restore original speed
		}
	}

	// Move towards target
	dx := enemy.Target.X - enemy.Position.X
	dy := enemy.Target.Y - enemy.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	cellSize := float64(g.config.GridSize)

	if distance < 5 {
		// Reached current target, move to next waypoint
		enemy.PathIndex++
		if enemy.PathIndex < len(g.path) {
			enemy.Target = Point{g.path[enemy.PathIndex].X*cellSize + cellSize/2, g.path[enemy.PathIndex].Y*cellSize + cellSize/2}
		}
	} else {
		// Move towards target
		enemy.Position.X += (dx / distance) * enemy.Speed
		enemy.Position.Y += (dy / distance) * enemy.Speed
	}
}

func (g *Game) placeTower(gridX, gridY float64) {
	// Check if position is valid (not on path and not occupied)
	if g.isOnPath(gridX, gridY) || g.isTowerAt(gridX, gridY) {
		return
	}

	cellSize := float64(g.config.GridSize)
	cost, damage, rangeVal, fireRate := g.config.GetTowerStats(g.selectedTowerType)

	if g.money >= cost {
		tower := &Tower{
			Position: Point{gridX*cellSize + cellSize/2, gridY*cellSize + cellSize/2},
			Range:    rangeVal,
			Damage:   damage,
			FireRate: fireRate,
			Cost:     cost,
			Type:     g.selectedTowerType,
			Special:  make(map[string]float64),
		}

		// Set special properties based on tower type
		switch g.selectedTowerType {
		case 5: // Splash Tower
			tower.Special["splash_radius"] = g.config.SplashRadius
		case 6: // Slow Tower
			tower.Special["slow_effect"] = g.config.SlowEffect
			tower.Special["slow_duration"] = g.config.SlowDuration
		}

		g.money -= cost
		g.towers = append(g.towers, tower)
	}
}

func (g *Game) isOnPath(gridX, gridY float64) bool {
	// Check waypoints
	for _, point := range g.path {
		if point.X == gridX && point.Y == gridY {
			return true
		}
	}

	// Check path segments between waypoints
	for i := 0; i < len(g.path)-1; i++ {
		if g.isOnPathSegment(gridX, gridY, g.path[i], g.path[i+1]) {
			return true
		}
	}
	return false
}

// isOnPathSegment checks if a grid cell lies on the path between two waypoints
func (g *Game) isOnPathSegment(gridX, gridY float64, start, end Point) bool {
	// For horizontal segments
	if start.Y == end.Y && gridY == start.Y {
		minX := math.Min(start.X, end.X)
		maxX := math.Max(start.X, end.X)
		if gridX >= minX && gridX <= maxX {
			return true
		}
	}

	// For vertical segments
	if start.X == end.X && gridX == start.X {
		minY := math.Min(start.Y, end.Y)
		maxY := math.Max(start.Y, end.Y)
		if gridY >= minY && gridY <= maxY {
			return true
		}
	}

	return false
}

func (g *Game) isTowerAt(gridX, gridY float64) bool {
	cellSize := float64(g.config.GridSize)
	for _, tower := range g.towers {
		towerGridX := (tower.Position.X - cellSize/2) / cellSize
		towerGridY := (tower.Position.Y - cellSize/2) / cellSize
		if math.Abs(towerGridX-gridX) < 0.1 && math.Abs(towerGridY-gridY) < 0.1 {
			return true
		}
	}
	return false
}

func (g *Game) findNearestEnemy(tower *Tower) *Enemy {
	var nearest *Enemy
	minDistance := tower.Range

	for _, enemy := range g.enemies {
		if !enemy.Alive {
			continue
		}

		dx := enemy.Position.X - tower.Position.X
		dy := enemy.Position.Y - tower.Position.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance <= minDistance {
			nearest = enemy
			minDistance = distance
		}
	}

	return nearest
}

// applyProjectileDamage applies damage and special effects from projectiles
func (g *Game) applyProjectileDamage(proj *Projectile) {
	if proj.Target == nil || !proj.Target.Alive {
		return
	}

	// Apply base damage
	proj.Target.Health -= proj.Damage
	if proj.Target.Health <= 0 {
		proj.Target.Alive = false
	}

	// Apply special effects from the source tower
	if proj.SourceTower != nil {
		switch proj.SourceTower.Type {
		case 5: // Splash Tower - damage nearby enemies
			g.applySplashDamage(proj.Target.Position, proj.SourceTower)
		case 6: // Slow Tower - apply slow effect
			g.applySlowEffect(proj.Target, proj.SourceTower)
		}
	}
}

// applySplashDamage damages enemies in splash radius
func (g *Game) applySplashDamage(center Point, tower *Tower) {
	radius := tower.Special["splash_radius"]
	splashDamage := float64(tower.Damage) / 2 // Half damage for splash

	for _, enemy := range g.enemies {
		if !enemy.Alive {
			continue
		}

		dx := enemy.Position.X - center.X
		dy := enemy.Position.Y - center.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance <= radius {
			enemy.Health -= int(splashDamage)
			if enemy.Health <= 0 {
				enemy.Alive = false
			}
			// Create small explosion for splash effect
			g.graphics.CreateExplosion(enemy.Position, 1, g.config)
		}
	}
}

// applySlowEffect applies slowing effect to enemy with duration
func (g *Game) applySlowEffect(enemy *Enemy, tower *Tower) {
	slowEffect := tower.Special["slow_effect"]
	slowDuration := tower.Special["slow_duration"]

	// Apply or refresh slow effect
	enemy.SlowTimer = slowDuration
	enemy.SlowIntensity = slowEffect
	enemy.Speed = enemy.BaseSpeed * slowEffect
}

func (g *Game) fireTower(tower *Tower, target *Enemy) {
	projectile := &Projectile{
		Position:    Point{tower.Position.X, tower.Position.Y},
		Target:      target,
		Speed:       5.0,
		Damage:      tower.Damage,
		Active:      true,
		SourceTower: tower, // Track which tower fired this projectile
	}
	g.projectiles = append(g.projectiles, projectile)
}

func (g *Game) moveProjectile(proj *Projectile) {
	if proj.Target == nil || !proj.Target.Alive {
		proj.Active = false
		return
	}

	dx := proj.Target.Position.X - proj.Position.X
	dy := proj.Target.Position.Y - proj.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 5 {
		// Hit target - create impact effect
		g.graphics.CreateExplosion(proj.Target.Position, 2, g.config)

		// Apply damage and special effects based on projectile type
		g.applyProjectileDamage(proj)

		proj.Active = false
	} else {
		// Move towards target
		proj.Position.X += (dx / distance) * proj.Speed
		proj.Position.Y += (dy / distance) * proj.Speed
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Handle different drawing based on game state
	switch g.modeManager.CurrentState {
	case StateMenu:
		g.modeManager.DrawMenu(screen, g.config)
		return
	case StatePlaying, StatePaused, StateGameOver, StateVictory:
		// Draw game content
		g.drawGameContent(screen)

		// Draw game state overlays
		g.modeManager.DrawGameState(screen, g)
		
		// Draw settings menu on top of everything
		if g.settingsMenu != nil && g.settingsMenu.IsVisible() {
			g.settingsMenu.Draw(screen)
		}
		
		// Draw help screen on top of everything
		if g.helpScreen != nil && g.helpScreen.IsVisible() {
			g.helpScreen.Draw(screen)
		}
	}
}

func (g *Game) drawGameContent(screen *ebiten.Image) {
	// Draw enhanced textured background
	g.graphics.DrawTexturedBackground(screen, g.config, g.path)

	// Draw enhanced towers with their types
	for _, tower := range g.towers {
		g.graphics.DrawEnhancedTower(screen, tower, tower.Type, g.config)
	}

	// Draw enhanced enemies
	for _, enemy := range g.enemies {
		g.graphics.DrawEnhancedEnemy(screen, enemy, g.config)
	}

	// Draw enhanced projectiles
	for _, proj := range g.projectiles {
		g.graphics.DrawEnhancedProjectile(screen, proj, g.config)
	}

	// Draw particle effects
	g.graphics.ParticleSystem.Draw(screen)

	// Draw UI with all tower types (only if playing)
	if g.modeManager.CurrentState == StatePlaying {
		cost1, _, _, _ := g.config.GetTowerStats(1)
		cost2, _, _, _ := g.config.GetTowerStats(2)
		cost3, _, _, _ := g.config.GetTowerStats(3)
		cost4, _, _, _ := g.config.GetTowerStats(4)
		cost5, _, _, _ := g.config.GetTowerStats(5)
		cost6, _, _, _ := g.config.GetTowerStats(6)

		// Add wave progress feedback
		waveStatus := ""
		if g.enemiesSpawned < g.enemiesPerWave {
			waveStatus = fmt.Sprintf(" - Spawning: %d/%d", g.enemiesSpawned, g.enemiesPerWave)
		} else if len(g.enemies) > 0 {
			waveStatus = fmt.Sprintf(" - Kill remaining: %d", len(g.enemies))
		} else if len(g.enemies) == 0 && g.enemiesSpawned >= g.enemiesPerWave {
			waveStatus = " - Press SPACE for next wave (BONUS!)"
		}

		uiText := fmt.Sprintf("Money: $%d | Lives: %d | Wave: %d%s\n\n"+
			"1: Basic ($%d)  2: Heavy ($%d)  3: Sniper ($%d)\n"+
			"4: Laser ($%d)  5: Splash ($%d)  6: Slow ($%d)\n\n"+
			"Selected: %s Tower\n"+
			"Click to place towers and defend against enemies!\n"+
			"Press SPACE when wave complete for bonus money!\n"+
			"Press S/F1 for Settings | H/F2 for Help",
			g.money, g.lives, g.wave, waveStatus,
			cost1, cost2, cost3, cost4, cost5, cost6,
			g.config.GetTowerName(g.selectedTowerType))

		// Add bonus display if recently earned
		if g.bonusDisplayTimer > 0 {
			uiText += fmt.Sprintf("\n\n🎉 EARLY WAVE BONUS: +$%d!", g.lastBonusEarned)
			g.bonusDisplayTimer -= 1.0 / 60.0
		}

		if g.config.ShowFPS {
			uiText += fmt.Sprintf("\nFPS: %.1f", ebiten.ActualFPS())
		}

		ebitenutil.DebugPrint(screen, uiText)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.WindowWidth, g.config.WindowHeight
}

func main() {
	// Set up panic recovery
	defer func() {
		if r := recover(); r != nil {
			LogError("PANIC: %v", r)
			fmt.Fprintf(os.Stderr, "Game crashed with error: %v\n", r)
			fmt.Fprintf(os.Stderr, "Please check the log file in the 'logs' directory for details.\n")
			os.Exit(1)
		}
	}()

	// Initialize logger
	isProduction := !IsDevBuild()
	if err := InitLogger(isProduction); err != nil {
		log.Printf("Warning: Failed to initialize logger: %v", err)
	}
	defer CloseLogger()

	// Log startup information
	LogInfo("=== Tower Defense Game Starting ===")
	LogInfo("Version: %s", GetVersionString())
	LogInfo("Build Info: %s", GetFullVersionString())

	// Determine config file to use
	configFile := "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
		LogInfo("Using custom config file: %s", configFile)
	}

	// Load configuration with error handling
	config, err := LoadConfig(configFile)
	if err != nil {
		LogWarn("Error loading config: %v, using defaults", err)
		config = DefaultConfig()
		// Try to save default config
		if saveErr := config.SaveConfig(configFile); saveErr != nil {
			LogWarn("Failed to save default config: %v", saveErr)
		}
	} else {
		LogInfo("Configuration loaded successfully from %s", configFile)
	}

	// Validate configuration
	config.ValidateConfig()

	// Initialize save manager
	saveManager, err := NewSaveManager()
	if err != nil {
		LogError("Failed to initialize save manager: %v", err)
		fmt.Fprintf(os.Stderr, "Error: Could not initialize save system: %v\n", err)
		os.Exit(1)
	}
	LogInfo("Save manager initialized")

	// Set window properties
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle(fmt.Sprintf("%s - %s", config.WindowTitle, GetVersionString()))
	if config.Fullscreen {
		ebiten.SetFullscreen(true)
	}
	if config.VSync {
		ebiten.SetVsyncEnabled(true)
	}

	LogInfo("Window initialized: %dx%d", config.WindowWidth, config.WindowHeight)

	game := NewGame(config)
	game.saveManager = saveManager

	LogInfo("Starting game loop...")
	if err := ebiten.RunGame(game); err != nil {
		LogError("Game error: %v", err)
		fmt.Fprintf(os.Stderr, "Game error: %v\n", err)
		os.Exit(1)
	}

	LogInfo("=== Game Shutdown Complete ===")
}
