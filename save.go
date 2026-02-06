package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SaveData represents the complete saved game state
type SaveData struct {
	Version      string          `json:"version"`
	SaveDate     time.Time       `json:"save_date"`
	GameMode     int             `json:"game_mode"` // GameMode as int
	CurrentLevel int             `json:"current_level"`
	EndlessWave  int             `json:"endless_wave,omitempty"`
	Money        int             `json:"money"`
	Lives        int             `json:"lives"`
	TowersPlaced []TowerSaveData `json:"towers"`
	HighScores   HighScores      `json:"high_scores"`
}

// TowerSaveData represents a saved tower
type TowerSaveData struct {
	Type  int     `json:"type"`
	GridX float64 `json:"grid_x"`
	GridY float64 `json:"grid_y"`
}

// HighScores tracks player achievements
type HighScores struct {
	HighestNormalLevel int `json:"highest_normal_level"`
	HighestEndlessWave int `json:"highest_endless_wave"`
	TotalGamesPlayed   int `json:"total_games_played"`
	TotalEnemiesKilled int `json:"total_enemies_killed"`
	TotalMoneyEarned   int `json:"total_money_earned"`
}

// SaveManager handles game save/load operations
type SaveManager struct {
	saveDir      string
	maxSaveSlots int
}

// NewSaveManager creates a new save manager
func NewSaveManager() (*SaveManager, error) {
	// Determine save directory based on platform
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	saveDir := filepath.Join(homeDir, ".tower-defense", "saves")

	// Create save directory if it doesn't exist
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create save directory: %w", err)
	}

	return &SaveManager{
		saveDir:      saveDir,
		maxSaveSlots: 5,
	}, nil
}

// SaveGame saves the current game state to a slot
func (sm *SaveManager) SaveGame(game *Game, slot int) error {
	if slot < 0 || slot >= sm.maxSaveSlots {
		return fmt.Errorf("invalid save slot: %d (must be 0-%d)", slot, sm.maxSaveSlots-1)
	}

	saveData := sm.createSaveData(game)

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	// Create backup of existing save if it exists
	savePath := sm.getSavePath(slot)
	if _, err := os.Stat(savePath); err == nil {
		backupPath := savePath + ".backup"
		if err := copyFile(savePath, backupPath); err != nil {
			LogWarn("Failed to create save backup: %v", err)
		}
	}

	// Write to file
	if err := os.WriteFile(savePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	LogInfo("Game saved to slot %d", slot)
	return nil
}

// LoadGame loads a game state from a slot
func (sm *SaveManager) LoadGame(game *Game, slot int) error {
	if slot < 0 || slot >= sm.maxSaveSlots {
		return fmt.Errorf("invalid save slot: %d (must be 0-%d)", slot, sm.maxSaveSlots-1)
	}

	savePath := sm.getSavePath(slot)

	// Check if save file exists
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		return fmt.Errorf("no save found in slot %d", slot)
	}

	// Read save file
	data, err := os.ReadFile(savePath)
	if err != nil {
		return fmt.Errorf("failed to read save file: %w", err)
	}

	// Unmarshal JSON
	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		// Try to restore from backup if available
		backupPath := savePath + ".backup"
		if backupData, backupErr := os.ReadFile(backupPath); backupErr == nil {
			if json.Unmarshal(backupData, &saveData) == nil {
				LogWarn("Restored save from backup after corruption")
				// Restore the backup as the main save
				os.WriteFile(savePath, backupData, 0644)
			} else {
				return fmt.Errorf("save file corrupted and backup invalid: %w", err)
			}
		} else {
			return fmt.Errorf("failed to parse save file: %w", err)
		}
	}

	// Apply save data to game
	if err := sm.applySaveData(game, &saveData); err != nil {
		return fmt.Errorf("failed to apply save data: %w", err)
	}

	LogInfo("Game loaded from slot %d", slot)
	return nil
}

// AutoSave performs an automatic save
func (sm *SaveManager) AutoSave(game *Game) error {
	return sm.SaveGame(game, 0) // Use slot 0 for autosave
}

// createSaveData creates a SaveData structure from current game state
func (sm *SaveManager) createSaveData(game *Game) *SaveData {
	cellSize := float64(game.config.GridSize)

	// Save tower positions
	towers := make([]TowerSaveData, len(game.towers))
	for i, tower := range game.towers {
		towers[i] = TowerSaveData{
			Type:  tower.Type,
			GridX: (tower.Position.X - cellSize/2) / cellSize,
			GridY: (tower.Position.Y - cellSize/2) / cellSize,
		}
	}

	return &SaveData{
		Version:      Version,
		SaveDate:     time.Now(),
		GameMode:     int(game.modeManager.CurrentMode),
		CurrentLevel: game.modeManager.CurrentLevel,
		EndlessWave:  game.modeManager.EndlessWave,
		Money:        game.money,
		Lives:        game.lives,
		TowersPlaced: towers,
		HighScores: HighScores{
			HighestNormalLevel: game.modeManager.CurrentLevel,
			HighestEndlessWave: game.modeManager.EndlessWave,
		},
	}
}

// applySaveData applies loaded save data to the game
func (sm *SaveManager) applySaveData(game *Game, saveData *SaveData) error {
	// Validate version compatibility
	if saveData.Version != Version && !IsDevBuild() {
		LogWarn("Loading save from different version: %s (current: %s)", saveData.Version, Version)
	}

	cellSize := float64(game.config.GridSize)

	// Restore game mode
	game.modeManager.CurrentMode = GameMode(saveData.GameMode)
	game.modeManager.CurrentLevel = saveData.CurrentLevel
	game.modeManager.EndlessWave = saveData.EndlessWave

	// Restore resources
	game.money = saveData.Money
	game.lives = saveData.Lives

	// Clear current state
	game.enemies = []*Enemy{}
	game.projectiles = []*Projectile{}
	game.towers = []*Tower{}

	// Restore towers
	for _, towerSave := range saveData.TowersPlaced {
		cost, damage, rangeVal, fireRate := game.config.GetTowerStats(towerSave.Type)
		tower := &Tower{
			Position: Point{
				X: towerSave.GridX*cellSize + cellSize/2,
				Y: towerSave.GridY*cellSize + cellSize/2,
			},
			Range:    rangeVal,
			Damage:   damage,
			FireRate: fireRate,
			Cost:     cost,
			Type:     towerSave.Type,
			Special:  make(map[string]float64),
		}

		// Set special properties
		switch towerSave.Type {
		case 5: // Splash Tower
			tower.Special["splash_radius"] = game.config.SplashRadius
		case 6: // Slow Tower
			tower.Special["slow_effect"] = game.config.SlowEffect
			tower.Special["slow_duration"] = game.config.SlowDuration
		}

		game.towers = append(game.towers, tower)
	}

	// Set up the level based on game mode
	switch game.modeManager.CurrentMode {
	case GameModeNormal:
		game.modeManager.setupLevel(game, saveData.CurrentLevel)
	case GameModeEndless:
		game.modeManager.setupEndlessWave(game)
	}

	return nil
}

// GetSaveInfo returns information about a save slot
func (sm *SaveManager) GetSaveInfo(slot int) (*SaveData, error) {
	if slot < 0 || slot >= sm.maxSaveSlots {
		return nil, fmt.Errorf("invalid save slot: %d", slot)
	}

	savePath := sm.getSavePath(slot)

	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		return nil, nil // No save in this slot
	}

	data, err := os.ReadFile(savePath)
	if err != nil {
		return nil, err
	}

	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		return nil, err
	}

	return &saveData, nil
}

// DeleteSave deletes a save slot
func (sm *SaveManager) DeleteSave(slot int) error {
	if slot < 0 || slot >= sm.maxSaveSlots {
		return fmt.Errorf("invalid save slot: %d", slot)
	}

	savePath := sm.getSavePath(slot)
	backupPath := savePath + ".backup"

	// Remove both save and backup
	os.Remove(savePath)
	os.Remove(backupPath)

	LogInfo("Deleted save slot %d", slot)
	return nil
}

// getSavePath returns the file path for a save slot
func (sm *SaveManager) getSavePath(slot int) string {
	return filepath.Join(sm.saveDir, fmt.Sprintf("save_%d.json", slot))
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
