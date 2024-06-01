package app

import (
	_tview "github.com/jneo8/juju-spell/internal/tview"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ExecuteAble interface {
	Execute() error
	BindFlags(cmd *cobra.Command) error
}

type App struct {
	logger   *logrus.Logger
	config   *viper.Viper
	tviewApp *tview.Application
}

func (app *App) BindFlags(cmd *cobra.Command) error {
	app.config.BindPFlags(cmd.Flags())
	app.logger.Info(config.AllSettings())
	return nil
}

func NewApp(logger *logrus.Logger, config *viper.Viper, app *tview.Application) ExecuteAble {
	logger.Info("NewApp")
	return &App{
		logger:   logger,
		config:   config,
		tviewApp: app,
	}
}

func (app *App) Execute() error {
	app.logger.Info("Execute")
	layout := _tview.GetLayout()
	if err := app.tviewApp.SetRoot(layout.RootFlex, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		return err
	}
	return nil
}
