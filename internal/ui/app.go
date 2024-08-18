package ui

import (
	"github.com/derailed/tview"
	"github.com/jneo8/jujuspell/internal/config"
)

type App struct {
	*tview.Application
	Main *Pages
}

func NewApp(cfg *config.Config) *App {
	a := App{
		Application: tview.NewApplication(),
		Main:        NewPages(),
	}
	return &a
}

func (a *App) Init() {
	a.SetRoot(a.Main, true).EnableMouse(true)
}
