//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/jneo8/juju-spell/internal/app"
)

func InitializeExecute() func() {
	wire.Build(GetExecute, app.NewLogger)
	return func() {}
}
