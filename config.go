package main

import (
	"encoding/json"
	"os"
)

// GameConfig holds all configurable game settings
type GameConfig struct {
	// Display settings
	WindowWidth  int    `json:"window_width"`
	WindowHeight int    `json:"window_height"`
	WindowTitle  string `json:"window_title"`
	Fullscreen   bool   `json:"fullscreen"`
	VSync        bool   `json:"vsync"`

	// Gameplay settings
	StartingMoney int     `json:"starting_money"`
	StartingLives int     `json:"starting_lives"`
	EnemySpeed    float64 `json:"enemy_speed"`
	SpawnDelay    float64 `json:"spawn_delay"`

	// Tower settings
	BasicTowerCost   int     `json:"basic_tower_cost"`
	BasicTowerDamage int     `json:"basic_tower_damage"`
	BasicTowerRange  float64 `json:"basic_tower_range"`
	BasicTowerRate   float64 `json:"basic_tower_fire_rate"`

	HeavyTowerCost   int     `json:"heavy_tower_cost"`
	HeavyTowerDamage int     `json:"heavy_tower_damage"`
	HeavyTowerRange  float64 `json:"heavy_tower_range"`
	HeavyTowerRate   float64 `json:"heavy_tower_fire_rate"`

	// Enemy settings
	BaseEnemyHealth int `json:"base_enemy_health"`
	HealthPerWave   int `json:"health_per_wave"`
	EnemyReward     int `json:"enemy_reward"`
	WaveBonus       int `json:"wave_bonus"`
	EnemiesPerWave  int `json:"enemies_per_wave"`

	// Visual settings
	ShowRange      bool `json:"show_range"`
	ShowHealthBars bool `json:"show_health_bars"`
	ShowFPS        bool `json:"show_fps"`
	GridSize       int  `json:"grid_size"`

	// Audio settings (for future use)
	MasterVolume float64 `json:"master_volume"`
	SFXVolume    float64 `json:"sfx_volume"`
	MusicVolume  float64 `json:"music_volume"`
	MuteAudio    bool    `json:"mute_audio"`

	// Controls
	PauseKey        string `json:"pause_key"`
	RestartKey      string `json:"restart_key"`
	TowerSelect1Key string `json:"tower_select_1_key"`
	TowerSelect2Key string `json:"tower_select_2_key"`

	// Debug settings
	DebugMode      bool `json:"debug_mode"`
	ShowPathPoints bool `json:"show_path_points"`
	ShowCollision  bool `json:"show_collision"`
	GodMode        bool `json:"god_mode"`
}

// DefaultConfig returns the default game configuration
func DefaultConfig() *GameConfig {
	return &GameConfig{
		// Display settings
		WindowWidth:  800,
		WindowHeight: 600,
		WindowTitle:  "Tower Defense",
		Fullscreen:   false,
		VSync:        true,

		// Gameplay settings
		StartingMoney: 100,
		StartingLives: 10,
		EnemySpeed:    1.0,
		SpawnDelay:    2.0,

		// Tower settings
		BasicTowerCost:   50,
		BasicTowerDamage: 20,
		BasicTowerRange:  80,
		BasicTowerRate:   1.0,

		HeavyTowerCost:   100,
		HeavyTowerDamage: 50,
		HeavyTowerRange:  60,
		HeavyTowerRate:   0.5,

		// Enemy settings
		BaseEnemyHealth: 50,
		HealthPerWave:   10,
		EnemyReward:     10,
		WaveBonus:       50,
		EnemiesPerWave:  3,

		// Visual settings
		ShowRange:      true,
		ShowHealthBars: true,
		ShowFPS:        false,
		GridSize:       40,

		// Audio settings
		MasterVolume: 1.0,
		SFXVolume:    0.8,
		MusicVolume:  0.6,
		MuteAudio:    false,

		// Controls
		PauseKey:        "Space",
		RestartKey:      "R",
		TowerSelect1Key: "1",
		TowerSelect2Key: "2",

		// Debug settings
		DebugMode:      false,
		ShowPathPoints: false,
		ShowCollision:  false,
		GodMode:        false,
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*GameConfig, error) {
	config := DefaultConfig()

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist, create it with default config
		if err := config.SaveConfig(filename); err != nil {
			return config, err
		}
		return config, nil
	}

	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	// Parse JSON
	if err := json.Unmarshal(data, config); err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig saves the configuration to a JSON file
func (c *GameConfig) SaveConfig(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// ValidateConfig ensures all config values are within valid ranges
func (c *GameConfig) ValidateConfig() {
	// Clamp window size
	if c.WindowWidth < 640 {
		c.WindowWidth = 640
	}
	if c.WindowHeight < 480 {
		c.WindowHeight = 480
	}
	if c.WindowWidth > 1920 {
		c.WindowWidth = 1920
	}
	if c.WindowHeight > 1080 {
		c.WindowHeight = 1080
	}

	// Clamp gameplay values
	if c.StartingMoney < 0 {
		c.StartingMoney = 0
	}
	if c.StartingLives < 1 {
		c.StartingLives = 1
	}
	if c.EnemySpeed <= 0 {
		c.EnemySpeed = 0.1
	}
	if c.SpawnDelay < 0.1 {
		c.SpawnDelay = 0.1
	}

	// Clamp tower values
	if c.BasicTowerCost < 1 {
		c.BasicTowerCost = 1
	}
	if c.BasicTowerDamage < 1 {
		c.BasicTowerDamage = 1
	}
	if c.BasicTowerRange < 10 {
		c.BasicTowerRange = 10
	}
	if c.BasicTowerRate <= 0 {
		c.BasicTowerRate = 0.1
	}

	if c.HeavyTowerCost < 1 {
		c.HeavyTowerCost = 1
	}
	if c.HeavyTowerDamage < 1 {
		c.HeavyTowerDamage = 1
	}
	if c.HeavyTowerRange < 10 {
		c.HeavyTowerRange = 10
	}
	if c.HeavyTowerRate <= 0 {
		c.HeavyTowerRate = 0.1
	}

	// Clamp enemy values
	if c.BaseEnemyHealth < 1 {
		c.BaseEnemyHealth = 1
	}
	if c.HealthPerWave < 0 {
		c.HealthPerWave = 0
	}
	if c.EnemyReward < 0 {
		c.EnemyReward = 0
	}
	if c.WaveBonus < 0 {
		c.WaveBonus = 0
	}
	if c.EnemiesPerWave < 1 {
		c.EnemiesPerWave = 1
	}

	// Clamp visual values
	if c.GridSize < 20 {
		c.GridSize = 20
	}
	if c.GridSize > 80 {
		c.GridSize = 80
	}

	// Clamp audio values
	if c.MasterVolume < 0 {
		c.MasterVolume = 0
	}
	if c.MasterVolume > 1 {
		c.MasterVolume = 1
	}
	if c.SFXVolume < 0 {
		c.SFXVolume = 0
	}
	if c.SFXVolume > 1 {
		c.SFXVolume = 1
	}
	if c.MusicVolume < 0 {
		c.MusicVolume = 0
	}
	if c.MusicVolume > 1 {
		c.MusicVolume = 1
	}
}

// GetTowerStats returns tower statistics based on tower type
func (c *GameConfig) GetTowerStats(towerType int) (cost int, damage int, rangeVal float64, fireRate float64) {
	switch towerType {
	case 1: // Basic Tower
		return c.BasicTowerCost, c.BasicTowerDamage, c.BasicTowerRange, c.BasicTowerRate
	case 2: // Heavy Tower
		return c.HeavyTowerCost, c.HeavyTowerDamage, c.HeavyTowerRange, c.HeavyTowerRate
	default:
		return c.BasicTowerCost, c.BasicTowerDamage, c.BasicTowerRange, c.BasicTowerRate
	}
}

// GetEnemyHealth returns the health for an enemy in the given wave
func (c *GameConfig) GetEnemyHealth(wave int) int {
	return c.BaseEnemyHealth + (wave-1)*c.HealthPerWave
}

// GetEnemiesInWave returns the number of enemies that should spawn in the given wave
func (c *GameConfig) GetEnemiesInWave(wave int) int {
	return c.EnemiesPerWave * wave
}
