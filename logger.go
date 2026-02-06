package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger handles structured logging for the game
type Logger struct {
	level      LogLevel
	fileLogger *log.Logger
	file       *os.File
	production bool
}

var gameLogger *Logger

// InitLogger initializes the global game logger
func InitLogger(production bool) error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := filepath.Join(logsDir, fmt.Sprintf("game_%s.log", timestamp))

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer to write to both file and stdout in dev mode
	var writer io.Writer
	if production {
		writer = file
	} else {
		writer = io.MultiWriter(os.Stdout, file)
	}

	logger := &Logger{
		level:      INFO,
		fileLogger: log.New(writer, "", 0), // We'll add our own prefix
		file:       file,
		production: production,
	}

	// Set debug level in development builds
	if !production {
		logger.level = DEBUG
	}

	gameLogger = logger

	// Clean up old log files (keep last 10)
	go cleanOldLogs(logsDir, 10)

	return nil
}

// CloseLogger closes the log file
func CloseLogger() {
	if gameLogger != nil && gameLogger.file != nil {
		gameLogger.file.Close()
	}
}

// SetLogLevel sets the minimum log level
func SetLogLevel(level LogLevel) {
	if gameLogger != nil {
		gameLogger.level = level
	}
}

// log writes a log message with the specified level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if l == nil || level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)

	l.fileLogger.Println(logLine)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	if gameLogger != nil {
		gameLogger.log(DEBUG, format, args...)
	}
}

// LogInfo logs an info message
func LogInfo(format string, args ...interface{}) {
	if gameLogger != nil {
		gameLogger.log(INFO, format, args...)
	}
}

// LogWarn logs a warning message
func LogWarn(format string, args ...interface{}) {
	if gameLogger != nil {
		gameLogger.log(WARN, format, args...)
	}
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	if gameLogger != nil {
		gameLogger.log(ERROR, format, args...)
	}
}

// LogPerformance logs performance metrics
func LogPerformance(operation string, duration time.Duration) {
	LogDebug("PERF: %s took %v", operation, duration)
}

// cleanOldLogs removes old log files, keeping only the most recent count
func cleanOldLogs(logsDir string, keepCount int) {
	files, err := filepath.Glob(filepath.Join(logsDir, "game_*.log"))
	if err != nil {
		return
	}

	if len(files) <= keepCount {
		return
	}

	// Sort by modification time (oldest first)
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	var fileInfos []fileInfo
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, fileInfo{path: file, modTime: info.ModTime()})
	}

	// Sort by mod time
	for i := 0; i < len(fileInfos)-1; i++ {
		for j := i + 1; j < len(fileInfos); j++ {
			if fileInfos[i].modTime.After(fileInfos[j].modTime) {
				fileInfos[i], fileInfos[j] = fileInfos[j], fileInfos[i]
			}
		}
	}

	// Delete oldest files
	deleteCount := len(fileInfos) - keepCount
	for i := 0; i < deleteCount; i++ {
		os.Remove(fileInfos[i].path)
	}
}
