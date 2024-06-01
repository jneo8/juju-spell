//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/rivo/tview"
)

func InitializeRootApp() (ExecuteAble, error) {
	wire.Build(NewApp, NewLogger, NewViper, tview.NewApplication)
	return &App{}, nil
}
