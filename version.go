package main

import (
	"fmt"
	"runtime"
	"time"
)

// Version information - updated during build
var (
	// Version is the semantic version of the game
	Version = "1.0.0-dev"

	// BuildDate is the date when the binary was built
	BuildDate = "unknown"

	// GitCommit is the git commit hash
	GitCommit = "unknown"

	// GoVersion is the Go version used to build
	GoVersion = runtime.Version()
)

// BuildInfo contains comprehensive build information
type BuildInfo struct {
	Version   string
	BuildDate string
	GitCommit string
	GoVersion string
	Platform  string
}

// GetBuildInfo returns the current build information
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		BuildDate: BuildDate,
		GitCommit: GitCommit,
		GoVersion: GoVersion,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetVersionString returns a formatted version string
func GetVersionString() string {
	if BuildDate != "unknown" {
		return fmt.Sprintf("v%s (built %s)", Version, BuildDate)
	}
	return fmt.Sprintf("v%s", Version)
}

// GetFullVersionString returns a detailed version string
func GetFullVersionString() string {
	info := GetBuildInfo()
	return fmt.Sprintf("Tower Defense v%s\nBuilt: %s\nCommit: %s\nGo: %s\nPlatform: %s",
		info.Version, info.BuildDate, info.GitCommit, info.GoVersion, info.Platform)
}

// IsDevBuild returns true if this is a development build
func IsDevBuild() bool {
	return BuildDate == "unknown" || GitCommit == "unknown"
}

// GetBuildTimestamp parses and returns the build date as a time.Time
func GetBuildTimestamp() (time.Time, error) {
	if BuildDate == "unknown" {
		return time.Time{}, fmt.Errorf("build date not set")
	}
	return time.Parse(time.RFC3339, BuildDate)
}
