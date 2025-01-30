package flog

import (
	"bufio"
	"os"
	"sync"
)

// CallerInfo is used internally to store function caller information for log messages.
type CallerInfo struct {
	funcName string // Name of the calling function
	line     int    // Line number in the source file
}

const (
	LogPanic   LogLevel = "panic"   //A panic with pitch red text that go cannot recover from
	LogError   LogLevel = "error"   //A error warning with red text to stand out
	LogWarn    LogLevel = "warn"    // A non-fatal error that Go can recover from
	LogInfo    LogLevel = "info"    // Normal log
	LogDebug   LogLevel = "debug"   // Debug log
	LogSuccess LogLevel = "success" //Success
	resetColor          = "\033[0m" // Reset color to default

)

type LogLevel string

var logFolder string

type Color string

var colorMap = make(map[LogLevel]Color)

type Logger struct {
	Config     Config                   //Current logger configuration
	logFolder  string                   //Base directory for log files. Appends a log folder automaticly to the path
	logFileMap map[string]*bufio.Writer //Map of buffered writers for each log level
	files      map[string]*os.File      //Map of open file handles
	mu         sync.Mutex               //Mutex for thread-safe operations
	bufPool    sync.Pool                //Pool of string builders for efficient string operations
}

type Config struct {
	LogFolder     string // Base folder for log files
	Colors        Colors // Color configuration for console output
	LogConsole    bool   // Enable/disable console logging
	LogFilePrefix string // Prefix for log files
	FormatPrefix  string // Prefix for format specifiers (default: "!")
}

type Colors struct {
	LogPanic   Color `json:"log_panic,omitempty"`
	LogError   Color `json:"log_error,omitempty"`
	LogWarn    Color `json:"log_warn,omitempty"`
	LogInfo    Color `json:"log_info,omitempty"`
	LogSuccess Color `json:"log_success,omitempty"`
	LogDebug   Color `json:"log_debug,omitempty"`
}

func (co Colors) Default() Colors {
	co.LogError = AnsiRGB(RGB{R: 234, G: 1, B: 1})
	co.LogPanic = AnsiRGB(RGB{R: 255, G: 0, B: 0})
	co.LogError = AnsiRGB(RGB{R: 234, G: 1, B: 1})
	co.LogWarn = AnsiRGB(RGB{R: 234, G: 173, B: 1})
	co.LogInfo = AnsiRGB(RGB{R: 0, G: 86, B: 234})
	co.LogSuccess = AnsiRGB(RGB{R: 1, G: 235, B: 110})
	return co
}

type RGB struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
}
