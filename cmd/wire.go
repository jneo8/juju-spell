//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/jneo8/jujuspell/common"
)

func InitializeExecute() func() {
	wire.Build(GetExecute, common.NewLogger)
	return func() {}
}
