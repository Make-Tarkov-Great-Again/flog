package flog

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
)

func (c Config) SetLogPath(path string) {
	c.LogFolder = path
}

func Default() Config {
	c := Config{}
	folder, _ := os.UserCacheDir()
	logFolder = filepath.Join(folder, "FLog") //Default ot FLog
	c.LogConsole = true
	c.LogFolder = path.Join(logFolder, "logs")
	c.Colors = Colors{}.Default()
	c.FormatPrefix = "!"
	c.LogFilePrefix = ""
	return c
}

func isColorsOmitted(config *Config) bool {
	return reflect.ValueOf(config.Colors).IsZero()
}

func isFolderOmitted(config *Config) bool {
	return reflect.ValueOf(config.LogFolder).IsZero()
}
