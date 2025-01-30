package flog

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Buffered writers for each log file
var logger *Logger = nil

/*
Initializes the logging system with the provided configuration. Must be called before any logging operations.

Does the following:
  - Creates necessary directories
  - Sets up log files
  - Initializes color mapping
  - Starts periodic flush routine
*/
func Init(config Config) {
	//defer measurment.Un(measurment.Trace("init"))

	logger = &Logger{
		logFolder: path.Join(config.LogFolder, "logs"),
		bufPool: sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		},
	}

	if isColorsOmitted(&config) {
		config.Colors.Default()
	}
	if isFolderOmitted(&config) {
		folder, _ := os.UserCacheDir()
		config.LogFolder = filepath.Join(folder, "FLog") //Default ot FLog
	}

	colorMap = map[LogLevel]Color{
		LogError:   config.Colors.LogError,
		LogWarn:    config.Colors.LogWarn,
		LogInfo:    config.Colors.LogInfo,
		LogSuccess: config.Colors.LogSuccess,
		LogDebug:   config.Colors.LogDebug,
	}

	logger.Config = config

	_ = logger.initFolders()
	_ = logger.initLogFiles()

}

func (l *Logger) initFolders() error {
	if _, err := os.Stat(l.Config.LogFolder); os.IsNotExist(err) {
		return os.MkdirAll(l.Config.LogFolder, 0755)
	}
	return nil
}

func GetConfig() Config {
	return logger.Config
}

func SetConfig(config Config) {
	Init(config)
}

func (l *Logger) initLogFiles() error {
	//defer measurment.Un(measurment.Trace("Init log files"))

	l.logFileMap = make(map[string]*bufio.Writer)
	l.files = make(map[string]*os.File)

	logTypes := []string{"warn", "error", "info", "debug"}

	for _, logType := range logTypes {
		logTypeFolder := path.Join(l.logFolder, logType)
		if _, err := os.Stat(logTypeFolder); os.IsNotExist(err) {
			if err := os.Mkdir(logTypeFolder, 0755); err != nil {
				return err
			}
		}

		logFileName := fmt.Sprintf("log_%s_%d.log", logType, time.Now().UnixNano())
		logFilePath := path.Join(logTypeFolder, logFileName)
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		l.files[logType] = file
		l.logFileMap[logType] = bufio.NewWriterSize(file, 32*1024) // 32KB buffer
	}

	// Start periodic flush
	go l.periodicFlush()

	return nil
}

func (l *Logger) periodicFlush() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		l.mu.Lock()
		for _, writer := range l.logFileMap {
			writer.Flush()
		}
		l.mu.Unlock()
	}
}

// Pre-allocated measurment format
var timeFormat = "15:04:05.000"

func (l *Logger) logDating() string {
	return time.Now().Format(timeFormat)
}

// Caller info cache
var (
	callerCache = make(map[uintptr]CallerInfo)
	callerMu    sync.RWMutex
)

/*
Internal function that retrieves information about the calling function:

Explanation:
  - Uses runtime.Caller to get call stack information
  - Implements caching for better performance
  - Returns CallerInfo with function name and line number
*/
func getCallerInfo(skip int) CallerInfo {
	//defer measurment.Un(measurment.Trace("Get caller info"))

	pc, _, line, _ := runtime.Caller(skip)

	callerMu.RLock()
	if info, ok := callerCache[pc]; ok {
		callerMu.RUnlock()
		info.line = line // Update line number
		return info
	}
	callerMu.RUnlock()

	callerFunc := runtime.FuncForPC(pc)
	if callerFunc == nil {
		return CallerInfo{"Anonymous", line}
	}

	fullName := callerFunc.Name()

	callerMu.Lock()
	callerCache[pc] = CallerInfo{fullName, line}
	info := callerCache[pc]
	callerMu.Unlock()

	return info
}

func (l *Logger) log(logType LogLevel, message string, args ...interface{}) {

	//defer measurment.Un(measurment.Trace("Log time"))
	// Get string builder from pool
	builder := l.bufPool.Get().(*strings.Builder)
	builder.Reset()
	defer l.bufPool.Put(builder)

	// Extract silent logging flag
	silentLogging := false
	if len(args) > 0 {
		if val, ok := args[len(args)-1].(bool); ok {
			silentLogging = val
			args = args[:len(args)-1]
		}
	}

	caller := getCallerInfo(3)

	// Format message
	formattedMsg := fmt.Sprintf(message, args...)
	timestamp := l.logDating()

	// Build log entry
	builder.WriteString("[")
	builder.WriteString(timestamp)
	builder.WriteString("] [")
	builder.WriteString(caller.funcName)
	builder.WriteString(" → ")
	builder.WriteString(fmt.Sprint(caller.line))
	builder.WriteString("]: ")
	builder.WriteString(formattedMsg)
	builder.WriteString("\n")

	logEntry := builder.String()

	// Write to log file
	l.mu.Lock()
	if writer, exists := l.logFileMap[string(logType)]; exists {
		writer.WriteString(logEntry)
	}
	l.mu.Unlock()

	// Write to console if not silent
	if !silentLogging && l.Config.LogConsole {
		var logTypeColor Color
		logTypePrefix := fmt.Sprintf("[%s]", strings.ToUpper(string(logType)))

		if logType == LogPanic {
			logTypeColor = colorMap[LogError]
		} else {
			logTypeColor = colorMap[logType]

		}

		//fmt.Printf("%s%s %s%s",
		//	logTypeColor,
		//	logTypePrefix,
		//	resetColor,
		//	logEntry)

		st := fmt.Sprintf("%s%s %s%s", logTypeColor, logTypePrefix, resetColor, logEntry)

		_, _ = os.Stdout.WriteString(st)
	}
}

func (l *Logger) format(logType LogLevel, message string, args ...interface{}) string {
	// Get string builder from pool
	builder := l.bufPool.Get().(*strings.Builder)
	builder.Reset()
	defer l.bufPool.Put(builder)

	// Get caller information
	caller := getCallerInfo(3)

	// Format the message
	formattedMsg := fmt.Sprintf(message, args...)
	timestamp := l.logDating()

	// Build the log entry
	builder.WriteString("[")
	builder.WriteString(timestamp)
	builder.WriteString("] [")
	builder.WriteString(caller.funcName)
	builder.WriteString(" → ")
	builder.WriteString(fmt.Sprint(caller.line))
	builder.WriteString("]: ")
	builder.WriteString(formattedMsg)
	builder.WriteString("\n")

	logEntry := builder.String()

	// Write to the log file
	l.mu.Lock()

	var logTypeThing LogLevel

	if logType == LogPanic {
		logTypeThing = LogError
	} else {
		logTypeThing = logType
	}

	if writer, exists := l.logFileMap[string(logTypeThing)]; exists {
		writer.WriteString(logEntry)
	}
	l.mu.Unlock()

	// Return the formatted log entry
	return logEntry
}

// Cleanup function to be called when shutting down
func (l *Logger) Cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, writer := range l.logFileMap {
		writer.Flush()
	}

	for _, file := range l.files {
		file.Close()
	}
}

func Cleanup() {
	l := logger
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, writer := range l.logFileMap {
		writer.Flush()
	}

	for _, file := range l.files {
		file.Close()
	}
}
