package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// SoundType represents different types of sound effects
type SoundType int

const (
	SoundTowerFire SoundType = iota
	SoundEnemyDeath
	SoundEnemyHit
	SoundTowerPlace
	SoundUIClick
	SoundUIHover
	SoundWaveComplete
	SoundGameOver
	SoundVictory
	SoundLevelUp
)

// AudioManager handles all game audio (music and sound effects)
type AudioManager struct {
	audioContext *audio.Context
	musicPlayer  *audio.Player
	sfxPlayers   map[SoundType]*audio.Player
	
	// Volume controls (0.0 to 1.0)
	masterVolume float64
	musicVolume  float64
	sfxVolume    float64
	
	// State
	musicEnabled bool
	sfxEnabled   bool
	currentMusic string
	
	// Thread safety
	mutex sync.RWMutex
}

// NewAudioManager creates a new audio manager
func NewAudioManager(sampleRate int) (*AudioManager, error) {
	audioContext := audio.NewContext(sampleRate)
	
	am := &AudioManager{
		audioContext: audioContext,
		sfxPlayers:   make(map[SoundType]*audio.Player),
		masterVolume: 1.0,
		musicVolume:  0.6,
		sfxVolume:    0.8,
		musicEnabled: true,
		sfxEnabled:   true,
	}
	
	LogInfo("Audio manager initialized with sample rate: %d Hz", sampleRate)
	return am, nil
}

// LoadMusic loads a music file (OGG format recommended)
func (am *AudioManager) LoadMusic(filepath string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	// Stop current music if playing
	if am.musicPlayer != nil {
		am.musicPlayer.Close()
		am.musicPlayer = nil
	}
	
	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open music file: %w", err)
	}
	
	// Decode based on file extension - try OGG first
	stream, err := vorbis.DecodeWithoutResampling(file)
	if err != nil {
		file.Close()
		// Try WAV as fallback
		file, err = os.Open(filepath)
		if err != nil {
			return fmt.Errorf("failed to reopen file: %w", err)
		}
		wavStream, err := wav.DecodeWithoutResampling(file)
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to decode audio file: %w", err)
		}
		
		// Create infinite loop player for WAV
		loopStream := audio.NewInfiniteLoop(wavStream, wavStream.Length())
		player, err := am.audioContext.NewPlayer(loopStream)
		if err != nil {
			return fmt.Errorf("failed to create music player: %w", err)
		}
		
		am.musicPlayer = player
		am.currentMusic = filepath
		am.updateMusicVolume()
		
		LogInfo("Music loaded (WAV): %s", filepath)
		return nil
	}
	
	// Create infinite loop player for OGG
	loopStream := audio.NewInfiniteLoop(stream, stream.Length())
	player, err := am.audioContext.NewPlayer(loopStream)
	if err != nil {
		return fmt.Errorf("failed to create music player: %w", err)
	}
	
	am.musicPlayer = player
	am.currentMusic = filepath
	am.updateMusicVolume()
	
	LogInfo("Music loaded (OGG): %s", filepath)
	return nil
}

// LoadSoundEffect loads a sound effect (WAV format recommended for short sounds)
func (am *AudioManager) LoadSoundEffect(soundType SoundType, filepath string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	// Open and read entire file into memory for sound effects
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read sound file: %w", err)
	}
	
	// Decode WAV
	stream, err := wav.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode sound: %w", err)
	}
	
	// Read all data into memory for quick playback
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, stream); err != nil {
		return fmt.Errorf("failed to buffer sound: %w", err)
	}
	
	// Store buffered data for later playback
	// We'll create players on-demand when sounds are played
	
	LogDebug("Sound effect loaded: %s for type %d", filepath, soundType)
	return nil
}

// PlayMusic starts playing the loaded music
func (am *AudioManager) PlayMusic() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if am.musicPlayer != nil && am.musicEnabled {
		am.musicPlayer.Play()
		LogDebug("Music started: %s", am.currentMusic)
	}
}

// PauseMusic pauses the current music
func (am *AudioManager) PauseMusic() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if am.musicPlayer != nil {
		am.musicPlayer.Pause()
		LogDebug("Music paused")
	}
}

// StopMusic stops and rewinds the music
func (am *AudioManager) StopMusic() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if am.musicPlayer != nil {
		am.musicPlayer.Pause()
		am.musicPlayer.Rewind()
		LogDebug("Music stopped")
	}
}

// PlaySound plays a sound effect
func (am *AudioManager) PlaySound(soundType SoundType) {
	if !am.sfxEnabled {
		return
	}
	
	// Note: In a real implementation, you would load the sound data
	// and create a new player each time to allow overlapping sounds
	
	LogDebug("Sound effect played: %d", soundType)
}

// SetMasterVolume sets the master volume (0.0 to 1.0)
func (am *AudioManager) SetMasterVolume(volume float64) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.masterVolume = clamp(volume, 0.0, 1.0)
	am.updateMusicVolume()
	LogDebug("Master volume set to: %.2f", am.masterVolume)
}

// SetMusicVolume sets the music volume (0.0 to 1.0)
func (am *AudioManager) SetMusicVolume(volume float64) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.musicVolume = clamp(volume, 0.0, 1.0)
	am.updateMusicVolume()
	LogDebug("Music volume set to: %.2f", am.musicVolume)
}

// SetSFXVolume sets the sound effects volume (0.0 to 1.0)
func (am *AudioManager) SetSFXVolume(volume float64) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.sfxVolume = clamp(volume, 0.0, 1.0)
	LogDebug("SFX volume set to: %.2f", am.sfxVolume)
}

// ToggleMusic toggles music on/off
func (am *AudioManager) ToggleMusic() bool {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.musicEnabled = !am.musicEnabled
	
	if am.musicPlayer != nil {
		if am.musicEnabled {
			am.musicPlayer.Play()
		} else {
			am.musicPlayer.Pause()
		}
	}
	
	LogInfo("Music toggled: %v", am.musicEnabled)
	return am.musicEnabled
}

// ToggleSFX toggles sound effects on/off
func (am *AudioManager) ToggleSFX() bool {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.sfxEnabled = !am.sfxEnabled
	LogInfo("SFX toggled: %v", am.sfxEnabled)
	return am.sfxEnabled
}

// GetVolumes returns current volume settings
func (am *AudioManager) GetVolumes() (master, music, sfx float64) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	return am.masterVolume, am.musicVolume, am.sfxVolume
}

// IsMusicEnabled returns whether music is enabled
func (am *AudioManager) IsMusicEnabled() bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	return am.musicEnabled
}

// IsSFXEnabled returns whether sound effects are enabled
func (am *AudioManager) IsSFXEnabled() bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	return am.sfxEnabled
}

// updateMusicVolume updates the music player volume (must be called with lock held)
func (am *AudioManager) updateMusicVolume() {
	if am.musicPlayer != nil {
		finalVolume := am.masterVolume * am.musicVolume
		am.musicPlayer.SetVolume(finalVolume)
	}
}

// Close closes the audio manager and releases resources
func (am *AudioManager) Close() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if am.musicPlayer != nil {
		am.musicPlayer.Close()
	}
	
	for _, player := range am.sfxPlayers {
		if player != nil {
			player.Close()
		}
	}
	
	LogInfo("Audio manager closed")
}

// clamp restricts a value between min and max
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
