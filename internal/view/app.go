package view

import (
	"github.com/jneo8/jujuspell/internal/config"
	log "github.com/sirupsen/logrus"
)

type App struct {
}

func NewApp(cfg *config.Config) *App {
	return &App{}
}

func (a *App) Run() error {
	log.Debug("Start app")
	return nil
}
