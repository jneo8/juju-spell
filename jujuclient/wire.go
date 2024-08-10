//go:build wireinject
// +build wireinject

package jujuclient

import (
	"github.com/google/wire"
	"github.com/jneo8/jujuspell/common"
	"github.com/juju/juju/jujuclient"
)

func InitializeJujuClient() (JujuClient, error) {
	wire.Build(NewJujuClient, jujuclient.NewFileClientStore, common.NewLogger)
	return &jujuClient{}, nil
}
