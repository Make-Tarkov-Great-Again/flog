package flog

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Silent = "true"
)

func Info(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		logger.log(LogInfo, message)
	}
}

func InfoF(format string, data ...any) {
	if logger != nil {
		message := formatlog(format, data...)
		logger.log(LogInfo, message)
	}
}

// SInfo is a silent info log
func SInfo(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		logger.log(LogInfo, message, true)
	}
}

func Error(data ...any) {
	//defer measurment.Un(measurment.Trace("Error"))

	if logger != nil {
		message := logger.prepare(data...)
		logger.log(LogError, message)
	}
}

func ErrorF(format string, data ...any) {
	if logger != nil {
		message := formatlog(format, data...)
		logger.log(LogError, message)
	}
}

// Info logs an info message
func (l *Logger) Info(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		l.log(LogInfo, message)
	}
}

// SInfo is a silent info log
func (l *Logger) SInfo(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		logger.log(LogInfo, message, true)
	}
}
func formatlog(format string, data ...any) string {
	// If no formatting needed or no data, return original
	if !strings.Contains(format, logger.Config.FormatPrefix) || len(data) == 0 {
		return format
	}

	result := format
	// Replace em
	for i, arg := range data {
		placeholder := logger.Config.FormatPrefix + strconv.Itoa(i+1)
		// Convert arg to string based on its type
		var strVal string
		switch v := arg.(type) {
		case string:
			strVal = v
		case error:
			strVal = v.Error()
		case fmt.Stringer:
			strVal = v.String()
		case []byte:
			strVal = string(v)
		default:
			strVal = fmt.Sprint(v) // Fallback to fmt.Sprint for unknown types
		}
		result = strings.Replace(result, placeholder, strVal, -1)
	}
	return result
}

/*
Internal function that prepares formatted messages:

Does:
  - Converts arguments to strings
  - Handles format strings if present
  - Joins multiple arguments with spaces
*/
func (l *Logger) prepare(data ...any) string {
	if len(data) == 0 {
		return ""
	}

	// Convert everything to strings
	parts := make([]string, len(data))
	for i, arg := range data {
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

	// Handle format string if first arg is a string and we have more data
	if fmtStr, ok := data[0].(string); ok && len(data) > 1 {
		if strings.Contains(fmtStr, "%") {
			return fmt.Sprintf(fmtStr, data[1:]...)
		}
	}

	return strings.Join(parts, " ")
}

// Warn logs a warning message
func (l *Logger) Warn(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		l.log(LogWarn, message)
	}
}

func WarnF(format string, data ...any) {
	if logger != nil {
		message := formatlog(format, data...)
		logger.log(LogWarn, message)
	}
}

func Warn(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		logger.log(LogWarn, message)
	}
}

// Error logs an error message
func (l *Logger) Error(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		l.log(LogError, message)
	}
}

// Debug logs a debug message
func Debug(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		logger.log(LogWarn, message)
	}
}

func DebugF(format string, data ...any) {
	if logger != nil {
		message := formatlog(format, data...)
		logger.log(LogDebug, message)
	}
}

func (l *Logger) Debug(data ...any) {
	if logger != nil {
		message := logger.prepare(data...)
		l.log(LogWarn, message)
	}
}

func AnsiRGB(rgb RGB) Color {
	s := "\x1b[38;2;" + strconv.FormatInt(int64(rgb.R), 10) + ";" + strconv.FormatInt(int64(rgb.G), 10) + ";" + strconv.FormatInt(int64(rgb.B), 10) + "m"
	return Color(s)
}

func (c Color) Set(in Color) {
	c = in
}
