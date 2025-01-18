package flog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

import (
	"log"
)

var (
	logger *Logger
)

func init() {
	folder, _ := os.UserCacheDir()
	logFolder = filepath.Join(folder, "Make Tarkov Great Again", "Server")
	fmt.Println(path.Join(logFolder, "logs"))
	logger = &Logger{
		logFolder:  path.Join(logFolder, "logs"),
		logConsole: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	_ = logger.initFolders()

	_ = logger.initLogFiles()

}

func (l *Logger) initFolders() error {
	if _, err := os.Stat(l.logFolder); os.IsNotExist(err) {
		return os.MkdirAll(l.logFolder, 0755)
	}
	return nil
}

func (l *Logger) initLogFiles() error {
	l.logFileMap = make(map[string]*os.File)

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
		l.logFileMap[logType] = file
	}

	return nil
}

func (l *Logger) logDating() string {
	now := time.Now()
	milliseconds := now.Nanosecond() / 1e6
	return fmt.Sprintf("%02d:%02d:%02d:%03d", now.Hour(), now.Minute(), now.Second(), milliseconds)
}

func getCallerInfo(skip int) CallerInfo {
	pc, _, line, _ := runtime.Caller(skip)
	callerFunc := runtime.FuncForPC(pc)

	if callerFunc == nil {
		return CallerInfo{"Anonymous", 0}
	}

	fullName := callerFunc.Name()
	parts := strings.Split(fullName, ".")

	if len(parts) <= 0 {
		return CallerInfo{"Anonymous", 0}
	}
	return CallerInfo{fullName, line}
}

func (l *Logger) log(logType LogLevel, message string, args ...interface{}) {
	// Extract silent logging flag
	silentLogging := false
	if len(args) > 0 {
		if val, ok := args[len(args)-1].(bool); ok {
			silentLogging = val
			args = args[:len(args)-1]
		}
	}

	// Get caller information (maybe)
	caller := getCallerInfo(3) // Skip two frames for the wrapper functions (again maybe, i just learned about this package)

	// Format the message
	formattedMsg := fmt.Sprintf(message, args...)
	timestamp := l.logDating()
	prefix := fmt.Sprintf("[%s] [%s → %d]:", timestamp, caller.funcName, caller.line)

	// Write to log file
	if file, exists := l.logFileMap[string(logType)]; exists {
		fmt.Fprintf(file, "%s %s\n", prefix, formattedMsg)
	}

	// Write to console if not silent
	if !silentLogging {
		logTypePrefix := fmt.Sprintf("[%s]", strings.ToUpper(string(logType)))
		logTypeColor := colorMap[logType]
		fmt.Printf("%s%s %s%s %s\n",
			logTypeColor,
			logTypePrefix,
			resetColor,
			prefix,
			formattedMsg)
	}
}
