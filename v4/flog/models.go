package flog

import (
	"log"
	"os"
)

type CallerInfo struct {
	funcName string
	line     int
}

const (
	LogError   LogLevel = "error"   //A error warning with red text to stand out
	LogWarn    LogLevel = "warn"    // A non-fatal error that Go can recover from
	LogInfo    LogLevel = "info"    // Normal log
	LogSuccess LogLevel = "success" //Success
	resetColor          = "\033[0m" // Reset color to default

)

type LogLevel string

var logFolder string

var colorMap = map[LogLevel]string{
	LogError:   "\033[38;2;250;0;0m",
	LogWarn:    "\033[38;2;255;204;0m",
	LogInfo:    "\033[38;2;84;175;190m",
	LogSuccess: "\033[38;2;0;175;0m",
}

type Logger struct {
	logFolder     string
	logFileMap    map[string]*os.File
	logConsole    *log.Logger
	logFilePrefix string
}
