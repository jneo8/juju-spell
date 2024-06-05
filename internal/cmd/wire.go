//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/jneo8/juju-spell/internal/common"
)

func InitializeExecute() func() {
	wire.Build(GetExecute, common.NewLogger)
	return func() {}
}
