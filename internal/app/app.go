package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ExecuteAble interface {
	Execute() error
	BindFlags(cmd *cobra.Command) error
}

type App struct {
	logger *logrus.Logger
	config *viper.Viper
}

func (app *App) BindFlags(cmd *cobra.Command) error {
	app.config.BindPFlags(cmd.Flags())
	app.logger.Info(config.AllSettings())
	return nil
}

func NewApp(logger *logrus.Logger, config *viper.Viper) ExecuteAble {
	logger.Info("NewApp")
	return &App{
		logger: logger,
		config: config,
	}
}

func (app *App) Execute() error {
	app.logger.Info("Execute")
	return nil
}
