package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// SettingsMenu handles the settings/options interface
type SettingsMenu struct {
	visible           bool
	selectedOption    int
	options           []string
	config            *GameConfig
	audioManager      *AudioManager
	keyPressed        map[ebiten.Key]bool
	mousePressed      bool
	sliderDragging    int // -1 = none, 0 = master, 1 = music, 2 = sfx
	previousMouseX    int
}

// NewSettingsMenu creates a new settings menu
func NewSettingsMenu(config *GameConfig, audioManager *AudioManager) *SettingsMenu {
	return &SettingsMenu{
		visible:        false,
		selectedOption: 0,
		options: []string{
			"Master Volume",
			"Music Volume",
			"SFX Volume",
			"Music: On",
			"Sound Effects: On",
			"Show FPS",
			"Show Health Bars",
			"Show Range",
			"Fullscreen",
			"Back",
		},
		config:         config,
		audioManager:   audioManager,
		keyPressed:     make(map[ebiten.Key]bool),
		sliderDragging: -1,
	}
}

// Toggle toggles the settings menu visibility
func (sm *SettingsMenu) Toggle() {
	sm.visible = !sm.visible
	LogInfo("Settings menu toggled: %v", sm.visible)
}

// Show shows the settings menu
func (sm *SettingsMenu) Show() {
	sm.visible = true
}

// Hide hides the settings menu and saves settings
func (sm *SettingsMenu) Hide() {
	sm.visible = false
	sm.sliderDragging = -1
	// Auto-save settings when closing
	sm.SaveSettings()
}

// IsVisible returns whether the menu is visible
func (sm *SettingsMenu) IsVisible() bool {
	return sm.visible
}

// Update handles input for the settings menu
func (sm *SettingsMenu) Update() {
	if !sm.visible {
		return
	}
	
	// Handle keyboard navigation
	upPressed := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	downPressed := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	enterPressed := ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace)
	escPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	
	// Up navigation
	if upPressed && !sm.keyPressed[ebiten.KeyUp] {
		if sm.selectedOption > 0 {
			sm.selectedOption--
		}
	}
	sm.keyPressed[ebiten.KeyUp] = upPressed
	
	// Down navigation
	if downPressed && !sm.keyPressed[ebiten.KeyDown] {
		if sm.selectedOption < len(sm.options)-1 {
			sm.selectedOption++
		}
	}
	sm.keyPressed[ebiten.KeyDown] = downPressed
	
	// Handle selection
	if enterPressed && !sm.keyPressed[ebiten.KeyEnter] {
		sm.handleSelection()
	}
	sm.keyPressed[ebiten.KeyEnter] = enterPressed
	
	// ESC closes menu
	if escPressed && !sm.keyPressed[ebiten.KeyEscape] {
		sm.Hide()
	}
	sm.keyPressed[ebiten.KeyEscape] = escPressed
	
	// Handle volume slider dragging
	sm.handleSliders()
	
	// Update option labels
	sm.updateOptionLabels()
}

// handleSliders handles volume slider dragging
func (sm *SettingsMenu) handleSliders() {
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	mouseX, mouseY := ebiten.CursorPosition()
	
	if mousePressed && !sm.mousePressed {
		// Check if clicking on a slider
		sm.checkSliderClick(mouseX, mouseY)
	}
	
	if !mousePressed {
		sm.sliderDragging = -1
	}
	
	// Update slider values while dragging
	if sm.sliderDragging >= 0 && mousePressed {
		sm.updateSliderValue(mouseX)
	}
	
	sm.mousePressed = mousePressed
	sm.previousMouseX = mouseX
}

// checkSliderClick checks if a slider was clicked
func (sm *SettingsMenu) checkSliderClick(mouseX, mouseY int) {
	centerX := sm.config.WindowWidth / 2
	startY := 150
	
	for i := 0; i < 3; i++ {
		sliderY := startY + i*60
		sliderStartX := centerX + 100
		sliderEndX := sliderStartX + 200
		
		if mouseX >= sliderStartX && mouseX <= sliderEndX &&
			mouseY >= sliderY-10 && mouseY <= sliderY+10 {
			sm.sliderDragging = i
			sm.updateSliderValue(mouseX)
			break
		}
	}
}

// updateSliderValue updates a slider value based on mouse position
func (sm *SettingsMenu) updateSliderValue(mouseX int) {
	centerX := sm.config.WindowWidth / 2
	sliderStartX := centerX + 100
	sliderWidth := 200
	
	// Calculate value (0.0 to 1.0)
	relativeX := float64(mouseX - sliderStartX)
	value := relativeX / float64(sliderWidth)
	value = clamp(value, 0.0, 1.0)
	
	// Apply to appropriate volume
	switch sm.sliderDragging {
	case 0: // Master
		sm.audioManager.SetMasterVolume(value)
		sm.config.MasterVolume = value
	case 1: // Music
		sm.audioManager.SetMusicVolume(value)
		sm.config.MusicVolume = value
	case 2: // SFX
		sm.audioManager.SetSFXVolume(value)
		sm.config.SFXVolume = value
	}
}

// handleSelection handles menu item selection
func (sm *SettingsMenu) handleSelection() {
	switch sm.selectedOption {
	case 3: // Toggle Music
		enabled := sm.audioManager.ToggleMusic()
		sm.config.MuteAudio = !enabled
		
	case 4: // Toggle SFX
		sm.audioManager.ToggleSFX()
		
	case 5: // Toggle FPS
		sm.config.ShowFPS = !sm.config.ShowFPS
		
	case 6: // Toggle Health Bars
		sm.config.ShowHealthBars = !sm.config.ShowHealthBars
		
	case 7: // Toggle Range
		sm.config.ShowRange = !sm.config.ShowRange
		
	case 8: // Toggle Fullscreen
		sm.config.Fullscreen = !sm.config.Fullscreen
		ebiten.SetFullscreen(sm.config.Fullscreen)
		
	case 9: // Back
		sm.Hide()
	}
}

// updateOptionLabels updates dynamic option labels
func (sm *SettingsMenu) updateOptionLabels() {
	sm.options[3] = fmt.Sprintf("Music: %s", boolToOnOff(sm.audioManager.IsMusicEnabled()))
	sm.options[4] = fmt.Sprintf("Sound Effects: %s", boolToOnOff(sm.audioManager.IsSFXEnabled()))
	sm.options[5] = fmt.Sprintf("Show FPS: %s", boolToOnOff(sm.config.ShowFPS))
	sm.options[6] = fmt.Sprintf("Show Health Bars: %s", boolToOnOff(sm.config.ShowHealthBars))
	sm.options[7] = fmt.Sprintf("Show Range: %s", boolToOnOff(sm.config.ShowRange))
	sm.options[8] = fmt.Sprintf("Fullscreen: %s", boolToOnOff(sm.config.Fullscreen))
}

// Draw renders the settings menu
func (sm *SettingsMenu) Draw(screen *ebiten.Image) {
	if !sm.visible {
		return
	}
	
	// Semi-transparent overlay
	vector.DrawFilledRect(screen, 0, 0, float32(sm.config.WindowWidth), float32(sm.config.WindowHeight),
		color.RGBA{0, 0, 0, 200}, false)
	
	centerX := sm.config.WindowWidth / 2
	
	// Title
	titleText := "SETTINGS"
	ebitenutil.DebugPrintAt(screen, titleText, centerX-40, 80)
	
	// Draw options
	startY := 150
	
	for i, option := range sm.options {
		y := startY + i*60
		x := centerX - 150
		
		// Volume sliders (first 3 options)
		if i < 3 {
			// Draw label
			ebitenutil.DebugPrintAt(screen, option, x, y)
			
			// Draw slider
			sliderX := float32(centerX + 100)
			sliderY := float32(y)
			sliderWidth := float32(200)
			
			// Slider background
			vector.DrawFilledRect(screen, sliderX, sliderY, sliderWidth, 4, color.RGBA{60, 60, 60, 255}, false)
			
			// Slider fill
			var volume float64
			switch i {
			case 0:
				volume, _, _ = sm.audioManager.GetVolumes()
			case 1:
				_, volume, _ = sm.audioManager.GetVolumes()
			case 2:
				_, _, volume = sm.audioManager.GetVolumes()
			}
			
			fillWidth := float32(volume) * sliderWidth
			vector.DrawFilledRect(screen, sliderX, sliderY, fillWidth, 4, color.RGBA{100, 150, 200, 255}, false)
			
			// Slider handle
			handleX := sliderX + fillWidth
			vector.DrawFilledCircle(screen, handleX, sliderY+2, 8, color.RGBA{150, 200, 250, 255}, false)
			
			// Volume percentage
			percentText := fmt.Sprintf("%d%%", int(volume*100))
			ebitenutil.DebugPrintAt(screen, percentText, int(sliderX+sliderWidth+10), y-5)
			
		} else {
			// Regular options
			if i == sm.selectedOption {
				// Highlight
				vector.DrawFilledRect(screen, float32(x-10), float32(y-5), 320, 30, color.RGBA{50, 100, 150, 150}, false)
				optionText := "► " + option + " ◄"
				ebitenutil.DebugPrintAt(screen, optionText, x, y)
			} else {
				ebitenutil.DebugPrintAt(screen, option, x, y)
			}
		}
	}
	
	// Instructions
	instructionsText := "↑/↓: Navigate | ENTER: Select | ESC: Close | Click & Drag: Adjust Volume"
	ebitenutil.DebugPrintAt(screen, instructionsText, 50, sm.config.WindowHeight-50)
}

// SaveSettings saves current settings to config
func (sm *SettingsMenu) SaveSettings() error {
	// Update config with current audio settings
	master, music, sfx := sm.audioManager.GetVolumes()
	sm.config.MasterVolume = master
	sm.config.MusicVolume = music
	sm.config.SFXVolume = sfx
	sm.config.MuteAudio = !sm.audioManager.IsMusicEnabled()
	
	// Save to file
	if err := sm.config.SaveConfig("config.json"); err != nil {
		LogError("Failed to save settings: %v", err)
		return err
	}
	
	LogInfo("Settings saved successfully")
	return nil
}

// boolToOnOff converts boolean to On/Off string
func boolToOnOff(b bool) string {
	if b {
		return "On"
	}
	return "Off"
}
