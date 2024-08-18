package config

import (
	"github.com/spf13/viper"
)

const (
	AppName         = "jujuspell"
	DefaultLogLevel = "info"
	DefaultLogFile  = "jujuspell.log"
)

type Config struct {
	LogLevel string
	LogFile  string
}

func NewConfig() *Config {
	config := Config{
		LogLevel: viper.GetString("logLevel"),
		LogFile:  viper.GetString("logFile"),
	}
	return &config
}
