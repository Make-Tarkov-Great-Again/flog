package v4

import (
	"fmt"
	"strings"
)

const (
	Silent = "true"
)

func Info(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		logger.log(LogInfo, message)
	}
}

// SInfo is a silent info log
func SInfo(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		logger.log(LogInfo, message, true)
	}
}

func Error(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		logger.log(LogError, message)
	}
}

// Info logs an info message
func (l *Logger) Info(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		l.log(LogInfo, message)
	}
}

// SInfo is a silent info log
func (l *Logger) SInfo(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		logger.log(LogInfo, message, true)
	}
}

func (l *Logger) prepare(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	// Convert everything to strings
	parts := make([]string, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			parts[i] = v
		case fmt.Stringer:
			parts[i] = v.String()
		case error:
			parts[i] = v.Error()
		default:
			parts[i] = fmt.Sprint(v)
		}
	}

	// Handle format string if first arg is a string and we have more args
	if fmtStr, ok := args[0].(string); ok && len(args) > 1 {
		if strings.Contains(fmtStr, "%") {
			return fmt.Sprintf(fmtStr, args[1:]...)
		}
	}

	return strings.Join(parts, " ")
}

// Warn logs a warning message
func (l *Logger) Warn(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		l.log(LogWarn, message)
	}
}
func Warn(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		logger.log(LogWarn, message)
	}
}

// Error logs an error message
func (l *Logger) Error(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		l.log(LogError, message)
	}
}

// Debug logs a debug message
func Debug(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		logger.log(LogWarn, message)
	}
}

func (l *Logger) Debug(args ...any) {
	if logger != nil {
		message := logger.prepare(args...)
		l.log(LogWarn, message)
	}
}
