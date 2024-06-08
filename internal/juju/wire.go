//go:build wireinject
// +build wireinject

package juju

import (
	"github.com/google/wire"
	"github.com/jneo8/juju-spell/internal/common"
	"github.com/juju/juju/jujuclient"
)

func InitializeJujuClient() (JujuClient, error) {
	wire.Build(NewJujuClient, jujuclient.NewFileClientStore, common.NewLogger)
	return &jujuClient{}, nil
}
