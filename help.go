package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HelpScreen displays game tutorial and controls
type HelpScreen struct {
	visible    bool
	config     *GameConfig
	keyPressed map[ebiten.Key]bool
	page       int
	maxPages   int
}

// NewHelpScreen creates a new help screen
func NewHelpScreen(config *GameConfig) *HelpScreen {
	return &HelpScreen{
		visible:    false,
		config:     config,
		keyPressed: make(map[ebiten.Key]bool),
		page:       0,
		maxPages:   3, // 3 pages of help
	}
}

// Toggle toggles the help screen visibility
func (hs *HelpScreen) Toggle() {
	hs.visible = !hs.visible
	hs.page = 0 // Reset to first page
	LogInfo("Help screen toggled: %v", hs.visible)
}

// Show shows the help screen
func (hs *HelpScreen) Show() {
	hs.visible = true
	hs.page = 0
}

// Hide hides the help screen
func (hs *HelpScreen) Hide() {
	hs.visible = false
}

// IsVisible returns whether the screen is visible
func (hs *HelpScreen) IsVisible() bool {
	return hs.visible
}

// Update handles input for the help screen
func (hs *HelpScreen) Update() {
	if !hs.visible {
		return
	}

	escPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	leftPressed := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	enterPressed := ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace)

	// Close on ESC or ENTER
	if (escPressed && !hs.keyPressed[ebiten.KeyEscape]) || (enterPressed && !hs.keyPressed[ebiten.KeyEnter]) {
		hs.Hide()
	}
	hs.keyPressed[ebiten.KeyEscape] = escPressed
	hs.keyPressed[ebiten.KeyEnter] = enterPressed

	// Navigate pages
	if leftPressed && !hs.keyPressed[ebiten.KeyLeft] && hs.page > 0 {
		hs.page--
	}
	hs.keyPressed[ebiten.KeyLeft] = leftPressed

	if rightPressed && !hs.keyPressed[ebiten.KeyRight] && hs.page < hs.maxPages-1 {
		hs.page++
	}
	hs.keyPressed[ebiten.KeyRight] = rightPressed
}

// Draw renders the help screen
func (hs *HelpScreen) Draw(screen *ebiten.Image) {
	if !hs.visible {
		return
	}

	// Semi-transparent overlay
	vector.DrawFilledRect(screen, 0, 0, float32(hs.config.WindowWidth), float32(hs.config.WindowHeight),
		color.RGBA{0, 0, 0, 220}, false)

	centerX := hs.config.WindowWidth / 2

	// Title
	titleText := fmt.Sprintf("HELP - Page %d/%d", hs.page+1, hs.maxPages)
	ebitenutil.DebugPrintAt(screen, titleText, centerX-80, 40)

	// Content based on page
	switch hs.page {
	case 0:
		hs.drawPage1Controls(screen, centerX)
	case 1:
		hs.drawPage2Towers(screen, centerX)
	case 2:
		hs.drawPage3Strategy(screen, centerX)
	}

	// Navigation hints
	navText := "←/→: Change Page | ENTER/ESC: Close"
	ebitenutil.DebugPrintAt(screen, navText, centerX-140, hs.config.WindowHeight-50)
}

// drawPage1Controls draws the controls page
func (hs *HelpScreen) drawPage1Controls(screen *ebiten.Image, centerX int) {
	y := 100

	title := "=== GAME CONTROLS ==="
	ebitenutil.DebugPrintAt(screen, title, centerX-100, y)
	y += 40

	controls := []string{
		"MOUSE CLICK - Place tower at cursor position",
		"",
		"KEYS 1-6 - Select tower type:",
		"  1: Basic Tower ($50)    2: Heavy Tower ($100)",
		"  3: Sniper Tower ($150)  4: Laser Tower ($200)",
		"  5: Splash Tower ($180)  6: Slow Tower ($120)",
		"",
		"SPACEBAR - Start next wave early (BONUS MONEY!)",
		"",
		"S or F1 - Open Settings menu",
		"H or F2 - Toggle this Help screen",
		"ESC/P - Pause game",
		"M - Return to main menu (when paused)",
		"R - Restart game (on game over)",
	}

	for _, line := range controls {
		ebitenutil.DebugPrintAt(screen, line, 50, y)
		y += 25
	}
}

// drawPage2Towers draws the tower information page
func (hs *HelpScreen) drawPage2Towers(screen *ebiten.Image, centerX int) {
	y := 100

	title := "=== TOWER TYPES ==="
	ebitenutil.DebugPrintAt(screen, title, centerX-100, y)
	y += 40

	towers := []string{
		"BASIC TOWER - Balanced damage and range",
		"  Cost: $50 | Damage: 20 | Range: 80px | Rate: 1.0/s",
		"  Good for: Early game, general defense",
		"",
		"HEAVY TOWER - High damage, short range",
		"  Cost: $100 | Damage: 50 | Range: 60px | Rate: 0.5/s",
		"  Good for: Chokepoints, armored enemies",
		"",
		"SNIPER TOWER - Long range, high damage",
		"  Cost: $150 | Damage: 100 | Range: 150px | Rate: 0.3/s",
		"  Good for: Path corners, boss enemies",
		"",
		"LASER TOWER - Rapid fire, low damage",
		"  Cost: $200 | Damage: 15 | Range: 70px | Rate: 3.0/s",
		"  Good for: Multiple weak enemies",
		"",
		"SPLASH TOWER - Area damage",
		"  Cost: $180 | Damage: 40 (AoE) | Range: 65px",
		"  Good for: Clustered enemies",
		"",
		"SLOW TOWER - Reduces enemy speed",
		"  Cost: $120 | Damage: 10 | Range: 90px",
		"  Good for: Crowd control at path start",
	}

	for _, line := range towers {
		ebitenutil.DebugPrintAt(screen, line, 40, y)
		y += 20
	}
}

// drawPage3Strategy draws the strategy tips page
func (hs *HelpScreen) drawPage3Strategy(screen *ebiten.Image, centerX int) {
	y := 100

	title := "=== STRATEGY TIPS ==="
	ebitenutil.DebugPrintAt(screen, title, centerX-100, y)
	y += 40

	tips := []string{
		"EARLY GAME:",
		"  • Start with 1-2 Basic Towers near the start",
		"  • Save money for stronger towers later",
		"  • Place towers to maximize path coverage",
		"",
		"MID GAME:",
		"  • Add Slow Towers at path beginning",
		"  • Place Sniper Towers at corners",
		"  • Use Splash Towers where enemies cluster",
		"",
		"LATE GAME:",
		"  • Upgrade to Heavy/Laser Towers",
		"  • Balance your tower composition",
		"  • Use SPACEBAR for wave bonuses",
		"",
		"GENERAL TIPS:",
		"  • Don't spend all your money at once",
		"  • Towers can't be placed on the path",
		"  • Each enemy killed gives $10",
		"  • Wave completion gives $50 bonus",
		"  • Spacebar bonus can give up to 75% extra!",
		"  • Normal Mode: 10 levels to complete",
		"  • Endless Mode: Survive as long as you can",
		"",
		"GOOD LUCK, DEFENDER!",
	}

	for _, line := range tips {
		ebitenutil.DebugPrintAt(screen, line, 50, y)
		y += 20
	}
}
