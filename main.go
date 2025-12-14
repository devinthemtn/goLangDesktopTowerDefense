package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Point struct {
	X, Y float64
}

type Enemy struct {
	Position   Point
	Target     Point
	Health     int
	MaxHealth  int
	Speed      float64
	PathIndex  int
	Alive      bool
	ReachedEnd bool
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
	Position Point
	Target   *Enemy
	Speed    float64
	Damage   int
	Active   bool
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
		modeManager:       NewGameModeManager(),
	}

	return game
}

func (g *Game) Update() error {
	// Update game mode system
	if err := g.modeManager.Update(g); err != nil {
		return err
	}

	// Only update game logic if we're in playing state
	if g.modeManager.CurrentState != StatePlaying {
		return nil
	}

	if g.gameOver {
		return nil
	}

	// Update particle system
	g.graphics.ParticleSystem.Update()

	// Spawn enemies
	g.spawnTimer += 1.0 / 60.0
	if g.spawnTimer > g.config.SpawnDelay && g.enemiesSpawned < g.enemiesPerWave {
		g.spawnEnemy()
		g.spawnTimer = 0
		g.enemiesSpawned++
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

	// Handle mouse input for tower placement
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cellSize := g.config.GridSize
		gridX := x / cellSize
		gridY := y / cellSize
		g.placeTower(float64(gridX), float64(gridY))
	}

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

	// Wave completion logic is now handled by mode manager
	// This section is kept for compatibility but may be overridden

	return nil
}

func (g *Game) spawnEnemy() {
	if len(g.path) == 0 {
		return
	}

	cellSize := float64(g.config.GridSize)
	health := g.config.GetEnemyHealth(g.wave)

	enemy := &Enemy{
		Position:  Point{g.path[0].X*cellSize + cellSize/2, g.path[0].Y*cellSize + cellSize/2},
		Health:    health,
		MaxHealth: health,
		Speed:     g.config.EnemySpeed,
		PathIndex: 0,
		Alive:     true,
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
	for _, point := range g.path {
		if point.X == gridX && point.Y == gridY {
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

	// Find the tower that fired this projectile to get special effects
	var sourceTower *Tower
	for _, tower := range g.towers {
		// This is simplified - in a real game you'd track projectile source
		if tower.Type >= 3 { // For now, assume special towers
			sourceTower = tower
			break
		}
	}

	// Apply base damage
	proj.Target.Health -= proj.Damage
	if proj.Target.Health <= 0 {
		proj.Target.Alive = false
	}

	// Apply special effects if source tower has them
	if sourceTower != nil {
		switch sourceTower.Type {
		case 5: // Splash Tower - damage nearby enemies
			g.applySplashDamage(proj.Target.Position, sourceTower)
		case 6: // Slow Tower - apply slow effect
			g.applySlowEffect(proj.Target, sourceTower)
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

// applySlowEffect applies slowing effect to enemy
func (g *Game) applySlowEffect(enemy *Enemy, tower *Tower) {
	// Store original speed and apply slow
	if enemy.Speed >= 1.0 { // Only slow if not already slowed
		slowEffect := tower.Special["slow_effect"]
		enemy.Speed *= slowEffect
		// In a more complex system, you'd track slow duration and restore speed
	}
}

func (g *Game) fireTower(tower *Tower, target *Enemy) {
	projectile := &Projectile{
		Position: Point{tower.Position.X, tower.Position.Y},
		Target:   target,
		Speed:    5.0,
		Damage:   tower.Damage,
		Active:   true,
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

		uiText := fmt.Sprintf("Money: $%d | Lives: %d | Wave: %d | Enemies: %d/%d\n\n"+
			"1: Basic ($%d)  2: Heavy ($%d)  3: Sniper ($%d)\n"+
			"4: Laser ($%d)  5: Splash ($%d)  6: Slow ($%d)\n\n"+
			"Selected: %s Tower",
			g.money, g.lives, g.wave, g.enemiesSpawned, g.enemiesPerWave,
			cost1, cost2, cost3, cost4, cost5, cost6,
			g.config.GetTowerName(g.selectedTowerType))

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
	// Load configuration
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Printf("Error loading config: %v, using defaults", err)
		config = DefaultConfig()
	}

	// Validate configuration
	config.ValidateConfig()

	// Set window properties
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle(config.WindowTitle)
	if config.Fullscreen {
		ebiten.SetFullscreen(true)
	}
	if config.VSync {
		ebiten.SetVsyncEnabled(true)
	}

	game := NewGame(config)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
