package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	}

	return game
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

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

	// Handle key input for tower selection
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.selectedTowerType = 1
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.selectedTowerType = 2
	}

	// Check wave completion
	if len(g.enemies) == 0 && g.enemiesSpawned >= g.enemiesPerWave {
		g.wave++
		g.money += g.config.WaveBonus
		g.spawnTimer = 0
		g.enemiesSpawned = 0
		g.enemiesPerWave = g.config.GetEnemiesInWave(g.wave)
	}

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
		// Hit target
		proj.Target.Health -= proj.Damage
		if proj.Target.Health <= 0 {
			proj.Target.Alive = false
		}
		proj.Active = false
	} else {
		// Move towards target
		proj.Position.X += (dx / distance) * proj.Speed
		proj.Position.Y += (dy / distance) * proj.Speed
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	cellSize := float32(g.config.GridSize)

	// Draw background
	screen.Fill(color.RGBA{34, 139, 34, 255}) // Forest green

	// Draw path
	for i, point := range g.path {
		x := float32(point.X) * cellSize
		y := float32(point.Y) * cellSize
		vector.DrawFilledRect(screen, x, y, cellSize, cellSize, color.RGBA{139, 69, 19, 255}, false) // Brown

		if i < len(g.path)-1 {
			next := g.path[i+1]
			nextX := float32(next.X)*cellSize + cellSize/2
			nextY := float32(next.Y)*cellSize + cellSize/2
			currX := x + cellSize/2
			currY := y + cellSize/2
			vector.StrokeLine(screen, currX, currY, nextX, nextY, 3, color.RGBA{101, 67, 33, 255}, false)
		}
	}

	// Draw towers
	for _, tower := range g.towers {
		vector.DrawFilledCircle(screen, float32(tower.Position.X), float32(tower.Position.Y), 15, color.RGBA{128, 128, 128, 255}, false) // Gray tower
		if g.config.ShowRange {
			vector.StrokeCircle(screen, float32(tower.Position.X), float32(tower.Position.Y), float32(tower.Range), 1, color.RGBA{255, 255, 255, 50}, false) // Range indicator
		}
	}

	// Draw enemies
	for _, enemy := range g.enemies {
		if enemy.Alive {
			x := float32(enemy.Position.X)
			y := float32(enemy.Position.Y)
			vector.DrawFilledCircle(screen, x, y, 10, color.RGBA{255, 0, 0, 255}, false) // Red enemy

			// Health bar
			if g.config.ShowHealthBars {
				barWidth := float32(20)
				barHeight := float32(4)
				healthRatio := float32(enemy.Health) / float32(enemy.MaxHealth)
				vector.DrawFilledRect(screen, x-barWidth/2, y-15, barWidth, barHeight, color.RGBA{255, 0, 0, 255}, false)
				vector.DrawFilledRect(screen, x-barWidth/2, y-15, barWidth*healthRatio, barHeight, color.RGBA{0, 255, 0, 255}, false)
			}
		}
	}

	// Draw projectiles
	for _, proj := range g.projectiles {
		if proj.Active {
			x := float32(proj.Position.X)
			y := float32(proj.Position.Y)
			vector.DrawFilledCircle(screen, x, y, 3, color.RGBA{255, 255, 0, 255}, false) // Yellow projectile
		}
	}

	// Draw UI
	cost1, _, _, _ := g.config.GetTowerStats(1)
	cost2, _, _, _ := g.config.GetTowerStats(2)

	uiText := fmt.Sprintf("Money: $%d\nLives: %d\nWave: %d\nEnemies: %d/%d\n\n1: Basic Tower ($%d)\n2: Heavy Tower ($%d)\n\nSelected: %d",
		g.money, g.lives, g.wave, g.enemiesSpawned, g.enemiesPerWave, cost1, cost2, g.selectedTowerType)

	if g.config.ShowFPS {
		uiText += fmt.Sprintf("\nFPS: %.1f", ebiten.ActualFPS())
	}

	ebitenutil.DebugPrint(screen, uiText)

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER!", g.config.WindowWidth/2-50, g.config.WindowHeight/2)
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
