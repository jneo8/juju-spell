//go:build wireinject
// +build wireinject

package juju

import (
	"github.com/google/wire"
	"github.com/jneo8/juju-spell/internal/common"
)

func InitializeJujuClient() (JujuClient, error) {
	wire.Build(NewJujuClient, common.NewLogger, GetJujuContent)
	return &jujuClient{}, nil
}
