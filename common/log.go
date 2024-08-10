package common

import (
	"os"
	"sync"

	filename "github.com/keepeye/logrus-filename"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	loggerOnce sync.Once
	logger     *logrus.Logger
)

func NewLogger() *logrus.Logger {
	loggerOnce.Do(
		func() {
			logger = logrus.New()
			logger.Info("New logger")
		},
	)
	return logger
}

func SetupLogger(logger *logrus.Logger, cfg *viper.Viper, logFile *os.File) {
	if logFile != nil {
		logger.SetOutput(logFile)
	}
	SetLoggerLevel(logger, cfg)
	SetFilenameHook(logger)
	SetFormatter(logger)
}

func SetFilenameHook(logger *logrus.Logger) {
	logger.AddHook(filename.NewHook())
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
		FullTimestamp:   true,
		ForceColors:     true,
	})
}

func SetLoggerLevel(logger *logrus.Logger, cfg *viper.Viper) {
	level := cfg.GetString("log_level")
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.DebugLevel)
	case "error":
		logger.SetLevel(logrus.DebugLevel)
	case "fatal":
		logger.SetLevel(logrus.DebugLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	logger.Infof("Log level: %s", logger.GetLevel().String())
}
