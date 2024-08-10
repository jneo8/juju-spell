//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/jneo8/jujuspell/common"
	"github.com/jneo8/jujuspell/jujuclient"
)

func InitializeRootApp() (ExecuteAble, error) {
	wire.Build(NewApp, common.NewLogger, common.NewViper, jujuclient.InitializeJujuClient)
	return &App{}, nil
}
