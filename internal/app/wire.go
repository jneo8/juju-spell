//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/jneo8/juju-spell/internal/common"
	"github.com/jneo8/juju-spell/internal/juju"
)

func InitializeRootApp() (ExecuteAble, error) {
	wire.Build(NewApp, common.NewLogger, common.NewViper, juju.InitializeJujuClient)
	return &App{}, nil
}
