package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Sprite represents a drawable game object with animations
type Sprite struct {
	Image        *ebiten.Image
	Width        int
	Height       int
	FrameCount   int
	CurrentFrame int
	AnimSpeed    float64
	AnimTimer    float64
}

// Particle represents a visual effect particle
type Particle struct {
	Position Point
	Velocity Point
	Life     float64
	MaxLife  float64
	Color    color.RGBA
	Size     float32
	Gravity  float64
	FadeOut  bool
	Active   bool
}

// ParticleSystem manages multiple particles
type ParticleSystem struct {
	Particles []*Particle
}

// GraphicsManager handles all visual effects and sprites
type GraphicsManager struct {
	ParticleSystem *ParticleSystem
	TowerSprites   map[int]*Sprite
	EnemySprite    *Sprite
	Textures       map[string]*ebiten.Image
}

// NewGraphicsManager creates a new graphics manager
func NewGraphicsManager() *GraphicsManager {
	gm := &GraphicsManager{
		ParticleSystem: &ParticleSystem{Particles: []*Particle{}},
		TowerSprites:   make(map[int]*Sprite),
		Textures:       make(map[string]*ebiten.Image),
	}

	gm.initializeSprites()
	gm.initializeTextures()
	return gm
}

// initializeSprites creates sprite data for game objects
func (gm *GraphicsManager) initializeSprites() {
	// Basic Tower Sprite (animated rotation)
	gm.TowerSprites[1] = &Sprite{
		Width:      30,
		Height:     30,
		FrameCount: 8,
		AnimSpeed:  0.1,
	}

	// Heavy Tower Sprite (pulsing animation)
	gm.TowerSprites[2] = &Sprite{
		Width:      35,
		Height:     35,
		FrameCount: 4,
		AnimSpeed:  0.15,
	}

	// Sniper Tower Sprite (slow tracking animation)
	gm.TowerSprites[3] = &Sprite{
		Width:      32,
		Height:     32,
		FrameCount: 12,
		AnimSpeed:  0.3,
	}

	// Laser Tower Sprite (fast spinning animation)
	gm.TowerSprites[4] = &Sprite{
		Width:      28,
		Height:     28,
		FrameCount: 16,
		AnimSpeed:  0.05,
	}

	// Splash Tower Sprite (charging animation)
	gm.TowerSprites[5] = &Sprite{
		Width:      38,
		Height:     38,
		FrameCount: 6,
		AnimSpeed:  0.2,
	}

	// Slow Tower Sprite (wave animation)
	gm.TowerSprites[6] = &Sprite{
		Width:      34,
		Height:     34,
		FrameCount: 8,
		AnimSpeed:  0.25,
	}

	// Enemy Sprite (walking animation)
	gm.EnemySprite = &Sprite{
		Width:      20,
		Height:     20,
		FrameCount: 6,
		AnimSpeed:  0.2,
	}
}

// initializeTextures creates procedural textures
func (gm *GraphicsManager) initializeTextures() {
	// Create grass texture
	gm.Textures["grass"] = gm.createGrassTexture(40, 40)

	// Create path texture
	gm.Textures["path"] = gm.createPathTexture(40, 40)
}

// createGrassTexture generates a grass-like texture
func (gm *GraphicsManager) createGrassTexture(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Base grass color
	baseColor := color.RGBA{34, 139, 34, 255}
	img.Fill(baseColor)

	// Add subtle grass texture without animation
	for x := 0; x < width; x += 8 {
		for y := 0; y < height; y += 8 {
			if (x+y)%16 == 0 { // Create a subtle pattern
				grassColor := color.RGBA{40, 150, 40, 255}
				vector.DrawFilledRect(img, float32(x), float32(y), 2, 3, grassColor, false)
			}
		}
	}

	return img
}

// createPathTexture generates a stone path texture
func (gm *GraphicsManager) createPathTexture(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Base path color (darker brown/gray)
	baseColor := color.RGBA{101, 67, 33, 255}
	img.Fill(baseColor)

	// Add static stone pattern
	for x := 0; x < width; x += 12 {
		for y := 0; y < height; y += 12 {
			if (x*y)%144 < 48 { // Create consistent stone pattern
				stoneColor := color.RGBA{85, 55, 25, 255}
				vector.DrawFilledCircle(img, float32(x+4), float32(y+4), 2, stoneColor, false)
			}
		}
	}

	return img
}

// DrawTexturedBackground draws the game background with textures
func (gm *GraphicsManager) DrawTexturedBackground(screen *ebiten.Image, config *GameConfig, path []Point) {
	cellSize := float32(config.GridSize)
	width := config.WindowWidth
	height := config.WindowHeight

	// Create a map to track path cells
	pathCells := make(map[Point]bool)
	for _, point := range path {
		pathCells[point] = true
	}

	// Draw textured tiles
	for x := 0; x < width; x += config.GridSize {
		for y := 0; y < height; y += config.GridSize {
			gridX := float64(x / config.GridSize)
			gridY := float64(y / config.GridSize)
			point := Point{gridX, gridY}

			var texture *ebiten.Image
			if pathCells[point] {
				texture = gm.Textures["path"]
			} else {
				// Use consistent grass texture (no random water)
				texture = gm.Textures["grass"]
			}

			// Draw the texture tile
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(texture, op)
		}
	}

	// Draw path connections with decorative elements
	gm.drawPathConnections(screen, path, cellSize)
}

// drawPathConnections draws decorative elements along the path
func (gm *GraphicsManager) drawPathConnections(screen *ebiten.Image, path []Point, cellSize float32) {
	for i := 0; i < len(path)-1; i++ {
		current := path[i]
		next := path[i+1]

		startX := float32(current.X)*cellSize + cellSize/2
		startY := float32(current.Y)*cellSize + cellSize/2
		endX := float32(next.X)*cellSize + cellSize/2
		endY := float32(next.Y)*cellSize + cellSize/2

		// Draw main path line
		vector.StrokeLine(screen, startX, startY, endX, endY, 8, color.RGBA{139, 99, 61, 255}, false)

		// Add decorative border
		vector.StrokeLine(screen, startX, startY, endX, endY, 12, color.RGBA{101, 67, 33, 180}, false)

		// Add consistent stones along the path
		steps := int(math.Sqrt(float64((endX-startX)*(endX-startX)+(endY-startY)*(endY-startY))) / 15)
		for step := 0; step < steps; step++ {
			t := float32(step) / float32(steps)
			x := startX + t*(endX-startX)
			y := startY + t*(endY-startY)

			if step%3 == 0 { // Every third step, place a stone
				stoneColor := color.RGBA{80, 60, 40, 200}
				vector.DrawFilledCircle(screen, x, y, 1.5, stoneColor, false)
			}
		}
	}
}

// DrawEnhancedTower draws a tower with improved graphics
func (gm *GraphicsManager) DrawEnhancedTower(screen *ebiten.Image, tower *Tower, towerType int, config *GameConfig) {
	x := float32(tower.Position.X)
	y := float32(tower.Position.Y)

	// Update animation
	sprite := gm.TowerSprites[towerType]
	if sprite != nil {
		sprite.AnimTimer += 1.0 / 60.0
		if sprite.AnimTimer >= sprite.AnimSpeed {
			sprite.CurrentFrame = (sprite.CurrentFrame + 1) % sprite.FrameCount
			sprite.AnimTimer = 0
		}
	}

	// Draw tower base (stone foundation)
	baseColor := color.RGBA{80, 80, 80, 255}
	vector.DrawFilledCircle(screen, x, y, 18, baseColor, false)
	vector.StrokeCircle(screen, x, y, 18, 2, color.RGBA{60, 60, 60, 255}, false)

	switch towerType {
	case 1: // Basic Tower
		gm.drawBasicTower(screen, x, y, sprite)
	case 2: // Heavy Tower
		gm.drawHeavyTower(screen, x, y, sprite)
	case 3: // Sniper Tower
		gm.drawSniperTower(screen, x, y, sprite)
	case 4: // Laser Tower
		gm.drawLaserTower(screen, x, y, sprite)
	case 5: // Splash Tower
		gm.drawSplashTower(screen, x, y, sprite)
	case 6: // Slow Tower
		gm.drawSlowTower(screen, x, y, sprite)
	}

	// Draw range indicator with gradient effect
	if config.ShowRange {
		gm.drawRangeIndicator(screen, x, y, float32(tower.Range))
	}

	// Draw subtle muzzle flash effect if tower recently fired
	if tower.LastFire < 0.05 {
		gm.drawMuzzleFlash(screen, x, y, towerType)
	}
}

// drawBasicTower draws the basic tower with rotation animation
func (gm *GraphicsManager) drawBasicTower(screen *ebiten.Image, x, y float32, sprite *Sprite) {
	// Tower body (cylinder)
	bodyColor := color.RGBA{128, 128, 128, 255}
	vector.DrawFilledCircle(screen, x, y, 15, bodyColor, false)

	// Add metallic shine
	shineColor := color.RGBA{180, 180, 180, 200}
	vector.DrawFilledCircle(screen, x-3, y-3, 8, shineColor, false)

	// Rotating cannon based on animation frame
	if sprite != nil {
		angle := float64(sprite.CurrentFrame) * math.Pi / 4
		cannonLength := float32(20)
		cannonX := x + float32(math.Cos(angle))*cannonLength
		cannonY := y + float32(math.Sin(angle))*cannonLength

		// Cannon barrel
		vector.StrokeLine(screen, x, y, cannonX, cannonY, 4, color.RGBA{64, 64, 64, 255}, false)

		// Cannon tip
		vector.DrawFilledCircle(screen, cannonX, cannonY, 3, color.RGBA{40, 40, 40, 255}, false)
	}
}

// drawHeavyTower draws the heavy tower with pulsing animation
func (gm *GraphicsManager) drawHeavyTower(screen *ebiten.Image, x, y float32, sprite *Sprite) {
	// Pulsing effect based on animation
	pulseIntensity := float32(1.0)
	if sprite != nil {
		pulseIntensity = 1.0 + 0.2*float32(math.Sin(float64(sprite.CurrentFrame)*math.Pi/2))
	}

	// Tower body (larger, more imposing)
	bodyColor := color.RGBA{96, 96, 96, 255}
	vector.DrawFilledCircle(screen, x, y, 16*pulseIntensity, bodyColor, false)

	// Add armor plating effect
	armorColor := color.RGBA{120, 120, 120, 255}
	vector.StrokeCircle(screen, x, y, 14*pulseIntensity, 2, armorColor, false)
	vector.StrokeCircle(screen, x, y, 10*pulseIntensity, 1, armorColor, false)

	// Multiple cannon barrels
	for i := 0; i < 3; i++ {
		angle := float64(i) * 2 * math.Pi / 3
		cannonLength := float32(18) * pulseIntensity
		cannonX := x + float32(math.Cos(angle))*cannonLength
		cannonY := y + float32(math.Sin(angle))*cannonLength

		// Thick cannon barrel
		vector.StrokeLine(screen, x, y, cannonX, cannonY, 6, color.RGBA{48, 48, 48, 255}, false)

		// Cannon muzzle
		vector.DrawFilledCircle(screen, cannonX, cannonY, 4, color.RGBA{32, 32, 32, 255}, false)
	}

	// Central core with energy effect
	coreColor := color.RGBA{255, 100, 100, uint8(100 + 100*pulseIntensity)}
	vector.DrawFilledCircle(screen, x, y, 6*pulseIntensity, coreColor, false)
}

// drawRangeIndicator draws a gradient range circle
func (gm *GraphicsManager) drawRangeIndicator(screen *ebiten.Image, x, y, radius float32) {
	// Draw multiple circles for gradient effect
	for i := 0; i < 5; i++ {
		alpha := uint8(20 - i*3)
		ringColor := color.RGBA{255, 255, 255, alpha}
		vector.StrokeCircle(screen, x, y, radius-float32(i), 1, ringColor, false)
	}
}

// drawMuzzleFlash creates a subtle muzzle flash effect
func (gm *GraphicsManager) drawMuzzleFlash(screen *ebiten.Image, x, y float32, towerType int) {
	// Create a simple, subtle flash effect
	flashColor := color.RGBA{255, 255, 200, 100}
	vector.DrawFilledCircle(screen, x, y, 8, flashColor, false)

	// Add a smaller bright center
	centerColor := color.RGBA{255, 255, 150, 150}
	vector.DrawFilledCircle(screen, x, y, 4, centerColor, false)
}

// drawSniperTower draws the sniper tower with long barrel and scope
func (gm *GraphicsManager) drawSniperTower(screen *ebiten.Image, x, y float32, sprite *Sprite) {
	// Tower base (elevated platform)
	baseColor := color.RGBA{70, 70, 70, 255}
	vector.DrawFilledCircle(screen, x, y, 17, baseColor, false)

	// Elevated platform
	platformColor := color.RGBA{90, 90, 90, 255}
	vector.DrawFilledCircle(screen, x, y-2, 14, platformColor, false)

	// Long sniper barrel based on animation frame
	if sprite != nil {
		angle := float64(sprite.CurrentFrame) * math.Pi / 6 // Slower, more precise tracking
		barrelLength := float32(35)                         // Extra long barrel
		barrelX := x + float32(math.Cos(angle))*barrelLength
		barrelY := y + float32(math.Sin(angle))*barrelLength

		// Main barrel (thick and long)
		vector.StrokeLine(screen, x, y, barrelX, barrelY, 6, color.RGBA{50, 50, 50, 255}, false)

		// Barrel tip with scope
		vector.DrawFilledCircle(screen, barrelX, barrelY, 4, color.RGBA{30, 30, 30, 255}, false)

		// Scope on top of barrel
		scopeX := x + float32(math.Cos(angle))*25
		scopeY := y + float32(math.Sin(angle))*25
		vector.DrawFilledRect(screen, scopeX-2, scopeY-4, 4, 8, color.RGBA{40, 40, 40, 255}, false)

		// Scope lens glint
		vector.DrawFilledCircle(screen, scopeX, scopeY, 1, color.RGBA{200, 200, 255, 200}, false)
	}

	// Central targeting system
	vector.DrawFilledCircle(screen, x, y, 6, color.RGBA{100, 50, 50, 255}, false)
}

// drawLaserTower draws the laser tower with energy crystals and beam emitters
func (gm *GraphicsManager) drawLaserTower(screen *ebiten.Image, x, y float32, sprite *Sprite) {
	// Tower base (crystalline structure)
	baseColor := color.RGBA{60, 80, 120, 255}
	vector.DrawFilledCircle(screen, x, y, 16, baseColor, false)

	// Energy crystal core
	if sprite != nil {
		intensity := 0.5 + 0.5*float32(math.Sin(float64(sprite.CurrentFrame)*math.Pi/8))
		coreColor := color.RGBA{100, 150, 255, uint8(150 + 100*intensity)}
		vector.DrawFilledCircle(screen, x, y, 8*intensity, coreColor, false)
	}

	// Multiple laser emitters arranged in a circle
	for i := 0; i < 6; i++ {
		angle := float64(i) * math.Pi / 3
		if sprite != nil {
			angle += float64(sprite.CurrentFrame) * math.Pi / 8 // Fast rotation
		}

		emitterX := x + float32(math.Cos(angle))*12
		emitterY := y + float32(math.Sin(angle))*12

		// Laser emitter
		emitterColor := color.RGBA{150, 200, 255, 255}
		vector.DrawFilledCircle(screen, emitterX, emitterY, 3, emitterColor, false)

		// Energy beam effect
		beamColor := color.RGBA{100, 200, 255, 100}
		vector.StrokeLine(screen, x, y, emitterX, emitterY, 2, beamColor, false)
	}

	// Central control unit
	vector.DrawFilledCircle(screen, x, y, 4, color.RGBA{200, 220, 255, 255}, false)
}

// drawSplashTower draws the splash tower with mortar design and explosive elements
func (gm *GraphicsManager) drawSplashTower(screen *ebiten.Image, x, y float32, sprite *Sprite) {
	// Heavy base platform
	baseColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x, y, 20, baseColor, false)

	// Reinforcement rings
	vector.StrokeCircle(screen, x, y, 18, 2, color.RGBA{60, 40, 20, 255}, false)
	vector.StrokeCircle(screen, x, y, 14, 1, color.RGBA{100, 80, 60, 255}, false)

	// Mortar tube (short and wide)
	tubeColor := color.RGBA{40, 40, 40, 255}
	vector.DrawFilledCircle(screen, x, y, 12, tubeColor, false)

	// Charging animation effect
	if sprite != nil {
		chargeLevel := float32(sprite.CurrentFrame) / float32(sprite.FrameCount)

		// Explosive energy building up
		if chargeLevel > 0.5 {
			energyColor := color.RGBA{255, 150, 50, uint8(100 + 100*chargeLevel)}
			vector.DrawFilledCircle(screen, x, y-5, 6*chargeLevel, energyColor, false)
		}

		// Loading shells around the base
		for i := 0; i < 4; i++ {
			angle := float64(i) * math.Pi / 2
			shellX := x + float32(math.Cos(angle))*16
			shellY := y + float32(math.Sin(angle))*16

			shellColor := color.RGBA{150, 100, 50, 255}
			vector.DrawFilledCircle(screen, shellX, shellY, 2, shellColor, false)
		}
	}

	// Barrel opening
	vector.DrawFilledCircle(screen, x, y, 8, color.RGBA{20, 20, 20, 255}, false)
}

// drawSlowTower draws the slow tower with ice crystals and freezing effects
func (gm *GraphicsManager) drawSlowTower(screen *ebiten.Image, x, y float32, sprite *Sprite) {
	// Ice crystal base
	baseColor := color.RGBA{150, 200, 255, 200}
	vector.DrawFilledCircle(screen, x, y, 18, baseColor, false)

	// Crystalline structure
	crystalColor := color.RGBA{200, 230, 255, 180}
	for i := 0; i < 6; i++ {
		angle := float64(i) * math.Pi / 3
		crystalX := x + float32(math.Cos(angle))*14
		crystalY := y + float32(math.Sin(angle))*14

		// Ice crystal spikes
		vector.StrokeLine(screen, x, y, crystalX, crystalY, 3, crystalColor, false)
		vector.DrawFilledCircle(screen, crystalX, crystalY, 2, crystalColor, false)
	}

	// Freezing wave animation
	if sprite != nil {
		waveRadius := 10 + float32(sprite.CurrentFrame)*2
		waveAlpha := uint8(200 - sprite.CurrentFrame*25)

		if waveAlpha > 0 {
			waveColor := color.RGBA{100, 150, 255, waveAlpha}
			vector.StrokeCircle(screen, x, y, waveRadius, 2, waveColor, false)
		}
	}

	// Central ice core
	coreColor := color.RGBA{180, 220, 255, 255}
	vector.DrawFilledCircle(screen, x, y, 6, coreColor, false)

	// Frost particles around tower
	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4
		if sprite != nil {
			angle += float64(sprite.CurrentFrame) * math.Pi / 16
		}

		particleX := x + float32(math.Cos(angle))*20
		particleY := y + float32(math.Sin(angle))*20

		frostColor := color.RGBA{200, 230, 255, 150}
		vector.DrawFilledCircle(screen, particleX, particleY, 1, frostColor, false)
	}
}

// DrawEnhancedEnemy draws an enemy with animation and effects
func (gm *GraphicsManager) DrawEnhancedEnemy(screen *ebiten.Image, enemy *Enemy, config *GameConfig) {
	if !enemy.Alive {
		return
	}

	x := float32(enemy.Position.X)
	y := float32(enemy.Position.Y)

	// Update walking animation
	if gm.EnemySprite != nil {
		gm.EnemySprite.AnimTimer += enemy.Speed / 60.0
		if gm.EnemySprite.AnimTimer >= gm.EnemySprite.AnimSpeed {
			gm.EnemySprite.CurrentFrame = (gm.EnemySprite.CurrentFrame + 1) % gm.EnemySprite.FrameCount
			gm.EnemySprite.AnimTimer = 0
		}
	}

	// Draw shadow
	shadowColor := color.RGBA{0, 0, 0, 100}
	vector.DrawFilledCircle(screen, x+2, y+2, 12, shadowColor, false)

	// Enemy body with breathing animation
	breathEffect := 1.0 + 0.1*math.Sin(float64(gm.EnemySprite.CurrentFrame)*math.Pi/3)
	bodySize := float32(10) * float32(breathEffect)

	// Health-based color (red when damaged)
	healthRatio := float64(enemy.Health) / float64(enemy.MaxHealth)
	enemyColor := color.RGBA{
		255,
		uint8(255 * healthRatio),
		uint8(255 * healthRatio),
		255,
	}

	// Draw enemy body
	vector.DrawFilledCircle(screen, x, y, bodySize, enemyColor, false)

	// Add armor/detail effects
	armorColor := color.RGBA{200, 200, 200, 200}
	vector.StrokeCircle(screen, x, y, bodySize-2, 1, armorColor, false)

	// Draw eyes
	eyeColor := color.RGBA{255, 255, 0, 255}
	vector.DrawFilledCircle(screen, x-3, y-2, 2, eyeColor, false)
	vector.DrawFilledCircle(screen, x+3, y-2, 2, eyeColor, false)

	// Draw movement trail particles
	if enemy.Speed > 0 && config.ParticleDensity > 0 {
		gm.createMovementTrail(enemy, config)
	}

	// Enhanced health bar
	if config.ShowHealthBars {
		gm.drawEnhancedHealthBar(screen, x, y, enemy)
	}
}

// drawEnhancedHealthBar draws a detailed health bar
func (gm *GraphicsManager) drawEnhancedHealthBar(screen *ebiten.Image, x, y float32, enemy *Enemy) {
	barWidth := float32(24)
	barHeight := float32(6)
	barX := x - barWidth/2
	barY := y - 18

	// Background (black border)
	vector.DrawFilledRect(screen, barX-1, barY-1, barWidth+2, barHeight+2, color.RGBA{0, 0, 0, 200}, false)

	// Health bar background
	vector.DrawFilledRect(screen, barX, barY, barWidth, barHeight, color.RGBA{100, 0, 0, 255}, false)

	// Health bar fill with gradient
	healthRatio := float32(enemy.Health) / float32(enemy.MaxHealth)
	healthWidth := barWidth * healthRatio

	// Color transitions: Green -> Yellow -> Red
	var barColor color.RGBA
	if healthRatio > 0.6 {
		barColor = color.RGBA{0, 255, 0, 255} // Green
	} else if healthRatio > 0.3 {
		barColor = color.RGBA{255, 255, 0, 255} // Yellow
	} else {
		barColor = color.RGBA{255, 0, 0, 255} // Red
	}

	vector.DrawFilledRect(screen, barX, barY, healthWidth, barHeight, barColor, false)

	// Add shine effect
	shineColor := color.RGBA{255, 255, 255, 100}
	vector.DrawFilledRect(screen, barX, barY, healthWidth, 2, shineColor, false)
}

// DrawEnhancedProjectile draws a projectile with trail effects
func (gm *GraphicsManager) DrawEnhancedProjectile(screen *ebiten.Image, proj *Projectile, config *GameConfig) {
	if !proj.Active {
		return
	}

	x := float32(proj.Position.X)
	y := float32(proj.Position.Y)

	// Create trail effect
	gm.createProjectileTrail(proj, config)

	// Different projectile appearance based on damage (tower type indicator)
	if proj.Damage >= 100 { // Sniper projectile
		// Long, thin projectile
		coreColor := color.RGBA{255, 100, 100, 255}
		vector.DrawFilledRect(screen, x-2, y-8, 4, 16, coreColor, false)

		glowColor := color.RGBA{255, 150, 150, 100}
		vector.DrawFilledRect(screen, x-3, y-9, 6, 18, glowColor, false)
	} else if proj.Damage >= 40 { // Splash projectile
		// Round, explosive projectile
		coreColor := color.RGBA{255, 150, 50, 255}
		vector.DrawFilledCircle(screen, x, y, 5, coreColor, false)

		// Sparkling effect
		sparkColor := color.RGBA{255, 200, 100, 200}
		vector.DrawFilledCircle(screen, x, y, 7, sparkColor, false)
	} else if proj.Damage <= 15 { // Laser projectile
		// Bright energy beam
		coreColor := color.RGBA{100, 200, 255, 255}
		vector.DrawFilledCircle(screen, x, y, 3, coreColor, false)

		// Energy glow
		glowColor := color.RGBA{150, 220, 255, 150}
		vector.DrawFilledCircle(screen, x, y, 8, glowColor, false)
	} else {
		// Standard projectile
		coreColor := color.RGBA{255, 255, 100, 255}
		vector.DrawFilledCircle(screen, x, y, 4, coreColor, false)

		// Add glow effect
		glowColor := color.RGBA{255, 200, 100, 150}
		vector.DrawFilledCircle(screen, x, y, 6, glowColor, false)
	}
}

// createMovementTrail creates particles behind moving enemies
func (gm *GraphicsManager) createMovementTrail(enemy *Enemy, config *GameConfig) {
	// Use particle density setting to control frequency
	baseFrequency := 0.1 * config.ParticleDensity
	if rand.Float64() < baseFrequency {
		particle := &Particle{
			Position: Point{enemy.Position.X, enemy.Position.Y + 5}, // No random offset
			Velocity: Point{-enemy.Speed * 0.3, 0},                  // Reduced velocity
			Life:     0.3,                                           // Shorter life
			MaxLife:  0.3,
			Color:    color.RGBA{139, 69, 19, 80}, // More transparent
			Size:     1,                           // Smaller size
			FadeOut:  true,
			Active:   true,
		}
		gm.ParticleSystem.AddParticle(particle)
	}
}

// createProjectileTrail creates particles behind projectiles
func (gm *GraphicsManager) createProjectileTrail(proj *Projectile, config *GameConfig) {
	baseFrequency := 0.3 * config.ParticleDensity
	if rand.Float64() < baseFrequency {
		particle := &Particle{
			Position: Point{proj.Position.X, proj.Position.Y}, // No random offset
			Velocity: Point{0, 0},                             // No random movement
			Life:     0.2,                                     // Shorter life
			MaxLife:  0.2,
			Color:    color.RGBA{255, 200, 100, 150}, // More transparent
			Size:     1,
			FadeOut:  true,
			Active:   true,
		}
		gm.ParticleSystem.AddParticle(particle)
	}
}

// CreateExplosion creates explosion particles when enemies die
func (gm *GraphicsManager) CreateExplosion(position Point, intensity int, config *GameConfig) {
	// Create particles based on density setting
	particleCount := int(float64(intensity*2) * config.ParticleDensity)
	for i := 0; i < particleCount; i++ {
		angle := float64(i) * (2 * math.Pi / float64(intensity*2)) // Even distribution

		particle := &Particle{
			Position: Point{
				position.X,
				position.Y,
			},
			Velocity: Point{
				math.Cos(angle) * 2.0,
				math.Sin(angle) * 2.0,
			},
			Life:    0.6,
			MaxLife: 0.6,
			Color: color.RGBA{
				255,
				150,
				50,
				200,
			},
			Size:    2,
			Gravity: 0.05,
			FadeOut: true,
			Active:  true,
		}
		gm.ParticleSystem.AddParticle(particle)
	}
}

// AddParticle adds a particle to the system
func (ps *ParticleSystem) AddParticle(particle *Particle) {
	ps.Particles = append(ps.Particles, particle)
}

// Update updates all particles in the system
func (ps *ParticleSystem) Update() {
	for i := len(ps.Particles) - 1; i >= 0; i-- {
		particle := ps.Particles[i]
		if !particle.Active {
			continue
		}

		// Update particle physics
		particle.Position.X += particle.Velocity.X
		particle.Position.Y += particle.Velocity.Y
		particle.Velocity.Y += particle.Gravity

		// Update life
		particle.Life -= 1.0 / 60.0
		if particle.Life <= 0 {
			particle.Active = false
			ps.Particles = append(ps.Particles[:i], ps.Particles[i+1:]...)
		}
	}
}

// Draw renders all particles
func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for _, particle := range ps.Particles {
		if !particle.Active {
			continue
		}

		x := float32(particle.Position.X)
		y := float32(particle.Position.Y)

		// Apply fade out effect
		alpha := particle.Color.A
		if particle.FadeOut {
			alpha = uint8(float64(particle.Color.A) * (particle.Life / particle.MaxLife))
		}

		particleColor := color.RGBA{
			particle.Color.R,
			particle.Color.G,
			particle.Color.B,
			alpha,
		}

		vector.DrawFilledCircle(screen, x, y, particle.Size, particleColor, false)
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
