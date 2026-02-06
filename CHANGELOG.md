# Changelog

All notable changes to the Tower Defense game will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0-dev] - 2026-02-06

### Added - Phase 2: Audio System & UI/UX Enhancements
- **Audio System**: Complete audio infrastructure
  - AudioManager with music and SFX support
  - OGG/Vorbis music playback with infinite looping
  - WAV sound effect infrastructure
  - Volume controls for Master, Music, and SFX (0.0-1.0 range)
  - Thread-safe audio operations
  - Graceful degradation if audio unavailable
  
- **Settings Menu**: Interactive configuration interface
  - Accessible via S key or F1 during gameplay
  - Interactive volume sliders with click & drag support
  - Visual feedback for all settings
  - Persistent settings (saves to config.json)
  - Toggle options: Music, SFX, FPS, Health Bars, Range, Fullscreen
  - Semi-transparent overlay design
  
- **Help/Tutorial Screen**: 3-page comprehensive guide
  - Accessible via H key or F2 during gameplay
  - Page 1: Complete controls reference
  - Page 2: All tower types with detailed stats
  - Page 3: Strategy tips for players
  - Keyboard navigation between pages (← →)
  
- **Dependencies**: Audio support libraries
  - github.com/ebitengine/oto/v3 for audio output
  - github.com/jfreymuth/oggvorbis for music codec

### Added - Phase 1: Foundation & Stability
- **Logging System**: Structured logging with DEBUG, INFO, WARN, ERROR levels
  - Logs written to `logs/` directory with timestamp
  - Automatic log rotation (keeps last 10 log files)
  - Multi-output in development mode (console + file)
  - Performance metrics logging capability
  
- **Save/Load System**: Complete game state persistence
  - JSON-based save format for compatibility
  - Multiple save slots (5 slots supported)
  - Auto-save on level completion and victory
  - Backup save files for corruption recovery
  - Save directory: `~/.tower-defense/saves/`
  - High scores tracking
  
- **Version & Build Information**: Comprehensive versioning system
  - Semantic versioning (v1.0.0-dev)
  - Build metadata (commit hash, build date, Go version, platform)
  - Version displayed in main menu
  - Build info embedded during compilation
  
- **Error Handling & Recovery**:
  - Panic recovery in main game loop
  - Graceful error handling for file I/O operations
  - Error boundaries for rendering failures
  - User-friendly error messages
  - Emergency auto-save before crash

### Changed
- Main menu now displays version information
- Window title includes version number
- Makefile updated with version-aware build system
- Build process now embeds git commit and build date
- UI hints updated with settings and help hotkeys

### Fixed
- Config loading errors are now handled gracefully with defaults
- Missing config file automatically creates default configuration
- Corrupted save files attempt recovery from backup

## [Future Releases]

### Planned for v1.1.0 - Phase 2 Completion
- Actual music and sound effect assets
- Sound effects hooked to game events (tower fire, enemy death, etc.)
- Confirmation dialogs for exits
- Enhanced visual feedback (hover states, click animations)
- Loading screen for level transitions

### Planned for v1.2.0 - Phase 2 Polish
- Difficulty selection (Easy/Normal/Hard)
- Tower sell/refund mechanism
- Game speed controls (1x, 2x speed)
- Achievement/stats tracking
- Enhanced visual effects (screen shake, better particles)

### Planned for v1.3.0 - Phase 3: Cross-Platform
- Windows, macOS, Linux builds
- WebAssembly (WASM) browser version
- Platform-specific installers
- Optimized distribution packages

### Planned for v1.4.0 - Phase 4: Monetization
- In-app purchase framework
- Analytics integration
- Player progression system
- Achievement tracking

---

[1.0.0-dev]: https://github.com/yourusername/golangTowerDefense/compare/v0.0.0...v1.0.0-dev
