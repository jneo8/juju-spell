package app

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	loggerOnce sync.Once
	logger     *logrus.Logger

	config     *viper.Viper
	configOnce sync.Once
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

func NewViper() *viper.Viper {
	configOnce.Do(
		func() {
			config = viper.New()
			logger.Info("New config")
		},
	)
	return config
}
