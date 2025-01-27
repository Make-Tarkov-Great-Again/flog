package flog

import (
	"bufio"
	"os"
	"strconv"
	"sync"
)

type CallerInfo struct {
	funcName string
	line     int
}

const (
	LogError   LogLevel = "error"   //A error warning with red text to stand out
	LogWarn    LogLevel = "warn"    // A non-fatal error that Go can recover from
	LogInfo    LogLevel = "info"    // Normal log
	LogDebug   LogLevel = "debug"   // Debug log
	LogSuccess LogLevel = "success" //Success
	resetColor          = "\033[0m" // Reset color to default

)

type LogLevel string

var logFolder string

var colorMap = make(map[LogLevel]string)

type Logger struct {
	Config     Config
	logFolder  string
	logFileMap map[string]*bufio.Writer
	files      map[string]*os.File
	mu         sync.Mutex
	bufPool    sync.Pool
}

type Config struct {
	LogFolder     string
	Colors        Colors `json:"colors,omitempty"`
	LogConsole    bool
	LogFilePrefix string
	FormatPrefix  string //Defaults to !
}

type Colors struct {
	LogError   string `json:"log_error,omitempty"`
	LogWarn    string `json:"log_warn,omitempty"`
	LogInfo    string `json:"log_info,omitempty"`
	LogSuccess string `json:"log_success,omitempty"`
	LogDebug   string `json:"log_debug,omitempty"`
}

func (co Colors) Default() Colors {
	co.LogError = AnsiRGB(RGB{R: 234, G: 1, B: 1})
	co.LogWarn = AnsiRGB(RGB{R: 234, G: 173, B: 1})
	co.LogInfo = AnsiRGB(RGB{R: 0, G: 86, B: 234})
	co.LogSuccess = AnsiRGB(RGB{R: 1, G: 235, B: 110})
	return co
}

//func (co Colors) Default() {
//	co.LogError = hexToAnsi("#EA0101")
//	co.LogWarn = hexToAnsi("#EAAD01")
//	co.LogInfo = hexToAnsi("#0056EA")
//	co.LogSuccess = hexToAnsi("#01EB6E")
//
//}

type RGB struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
}

// AnsiRGB creates a new ansi escape code based on an RGB input
func AnsiRGB(rgb RGB) string {
	s := "\x1b[38;2;" + strconv.FormatInt(int64(rgb.R), 10) + ";" + strconv.FormatInt(int64(rgb.G), 10) + ";" + strconv.FormatInt(int64(rgb.B), 10) + "m"
	return s
}
