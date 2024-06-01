//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
)

func InitializeRootApp() (ExecuteAble, error) {
	wire.Build(NewApp, NewLogger, NewViper)
	return &App{}, nil
}
