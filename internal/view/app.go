package view

import (
	"context"

	"github.com/derailed/tview"
	"github.com/jneo8/jujuspell/internal"
	"github.com/jneo8/jujuspell/internal/config"
	"github.com/jneo8/jujuspell/internal/model"
	"github.com/jneo8/jujuspell/internal/ui"
	log "github.com/sirupsen/logrus"
)

type App struct {
	*ui.App
	command *Command
	Content *PageStack
}

func NewApp(cfg *config.Config) *App {
	a := App{
		App:     ui.NewApp(cfg),
		Content: NewPageStack(),
	}
	return &a
}

func (a *App) Init() error {
	a.App.Init()
	a.command = NewCommand(a)
	a.layout()
	return nil
}

func (a *App) Run() error {
	a.Main.SwitchToPage("main")
	if err := a.command.runDefaultCommand(); err != nil {
		return err
	}
	if err := a.Application.Run(); err != nil {
		return nil
	}
	return nil
}

func (a *App) layout() {
	main := tview.NewFlex().SetDirection(tview.FlexRow)
	main.AddItem(a.Content, 0, 10, true)
	a.Main.AddPage("main", main, true, false)
}

func (a *App) inject(comp model.Component, clearStack bool) error {
	ctx := context.WithValue(context.Background(), internal.KeyApp, a)
	if err := comp.Init(ctx); err != nil {
		log.Error(err)
		return err
	}
	if clearStack {
		a.Content.Stack.Clear()
	}
	a.Content.Push(comp)
	return nil
}
