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

	SniperTowerCost   int     `json:"sniper_tower_cost"`
	SniperTowerDamage int     `json:"sniper_tower_damage"`
	SniperTowerRange  float64 `json:"sniper_tower_range"`
	SniperTowerRate   float64 `json:"sniper_tower_fire_rate"`

	LaserTowerCost   int     `json:"laser_tower_cost"`
	LaserTowerDamage int     `json:"laser_tower_damage"`
	LaserTowerRange  float64 `json:"laser_tower_range"`
	LaserTowerRate   float64 `json:"laser_tower_fire_rate"`

	SplashTowerCost   int     `json:"splash_tower_cost"`
	SplashTowerDamage int     `json:"splash_tower_damage"`
	SplashTowerRange  float64 `json:"splash_tower_range"`
	SplashTowerRate   float64 `json:"splash_tower_fire_rate"`
	SplashRadius      float64 `json:"splash_radius"`

	SlowTowerCost   int     `json:"slow_tower_cost"`
	SlowTowerDamage int     `json:"slow_tower_damage"`
	SlowTowerRange  float64 `json:"slow_tower_range"`
	SlowTowerRate   float64 `json:"slow_tower_fire_rate"`
	SlowEffect      float64 `json:"slow_effect"`
	SlowDuration    float64 `json:"slow_duration"`

	// Enemy settings
	BaseEnemyHealth int `json:"base_enemy_health"`
	HealthPerWave   int `json:"health_per_wave"`
	EnemyReward     int `json:"enemy_reward"`
	WaveBonus       int `json:"wave_bonus"`
	EnemiesPerWave  int `json:"enemies_per_wave"`

	// Visual settings
	ShowRange       bool    `json:"show_range"`
	ShowHealthBars  bool    `json:"show_health_bars"`
	ShowFPS         bool    `json:"show_fps"`
	GridSize        int     `json:"grid_size"`
	ParticleDensity float64 `json:"particle_density"`

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
	TowerSelect3Key string `json:"tower_select_3_key"`
	TowerSelect4Key string `json:"tower_select_4_key"`
	TowerSelect5Key string `json:"tower_select_5_key"`
	TowerSelect6Key string `json:"tower_select_6_key"`

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

		SniperTowerCost:   150,
		SniperTowerDamage: 100,
		SniperTowerRange:  150,
		SniperTowerRate:   0.3,

		LaserTowerCost:   200,
		LaserTowerDamage: 15,
		LaserTowerRange:  70,
		LaserTowerRate:   3.0,

		SplashTowerCost:   180,
		SplashTowerDamage: 40,
		SplashTowerRange:  65,
		SplashTowerRate:   0.8,
		SplashRadius:      30,

		SlowTowerCost:   120,
		SlowTowerDamage: 10,
		SlowTowerRange:  90,
		SlowTowerRate:   1.5,
		SlowEffect:      0.5,
		SlowDuration:    2.0,

		// Enemy settings
		BaseEnemyHealth: 50,
		HealthPerWave:   10,
		EnemyReward:     10,
		WaveBonus:       50,
		EnemiesPerWave:  3,

		// Visual settings
		ShowRange:       true,
		ShowHealthBars:  true,
		ShowFPS:         false,
		GridSize:        40,
		ParticleDensity: 1.0,

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
		TowerSelect3Key: "3",
		TowerSelect4Key: "4",
		TowerSelect5Key: "5",
		TowerSelect6Key: "6",

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
	if c.ParticleDensity < 0 {
		c.ParticleDensity = 0
	}
	if c.ParticleDensity > 2 {
		c.ParticleDensity = 2
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
	case 3: // Sniper Tower
		return c.SniperTowerCost, c.SniperTowerDamage, c.SniperTowerRange, c.SniperTowerRate
	case 4: // Laser Tower
		return c.LaserTowerCost, c.LaserTowerDamage, c.LaserTowerRange, c.LaserTowerRate
	case 5: // Splash Tower
		return c.SplashTowerCost, c.SplashTowerDamage, c.SplashTowerRange, c.SplashTowerRate
	case 6: // Slow Tower
		return c.SlowTowerCost, c.SlowTowerDamage, c.SlowTowerRange, c.SlowTowerRate
	default:
		return c.BasicTowerCost, c.BasicTowerDamage, c.BasicTowerRange, c.BasicTowerRate
	}
}

// GetTowerName returns the name of a tower type
func (c *GameConfig) GetTowerName(towerType int) string {
	switch towerType {
	case 1:
		return "Basic"
	case 2:
		return "Heavy"
	case 3:
		return "Sniper"
	case 4:
		return "Laser"
	case 5:
		return "Splash"
	case 6:
		return "Slow"
	default:
		return "Unknown"
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
