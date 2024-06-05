// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/jneo8/juju-spell/internal/common"
	"github.com/jneo8/juju-spell/internal/juju"
)

// Injectors from wire.go:

func InitializeRootApp() (ExecuteAble, error) {
	logger := common.NewLogger()
	viper := common.NewViper()
	jujuClient, err := juju.InitializeJujuClient()
	if err != nil {
		return nil, err
	}
	executeAble := NewApp(logger, viper, jujuClient)
	return executeAble, nil
}
