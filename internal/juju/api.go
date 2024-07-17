package juju

import (
	"errors"

	"github.com/juju/juju/api"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/juju"
	"github.com/juju/juju/jujuclient"
)

func (jc *jujuClient) GetRootAPI(controllerName string, modelName string) (api.Connection, error) {
	accountDetails, err := jc.clientStore.AccountDetails(controllerName)
	if err != nil {
		return nil, err
	}
	newAPIConnectionParams, err := jc.GetNewAPIConnectionParams(
		jc.clientStore,
		controllerName,
		modelName,
		accountDetails,
		api.Open,
	)
	if err != nil {
		return nil, err
	}

	apiRoot, err := juju.NewAPIConnection(newAPIConnectionParams)
	if err != nil {
		return nil, err
	}
	return apiRoot, nil
}

var errNoNameSpecified = errors.New("no name specified")

func (jc *jujuClient) GetNewAPIConnectionParams(
	store jujuclient.ClientStore,
	controllerName string,
	modelName string,
	accountDetails *jujuclient.AccountDetails,
	apiOpen api.OpenFunc,
) (juju.NewAPIConnectionParams, error) {
	if controllerName == "" {
		return juju.NewAPIConnectionParams{}, errNoNameSpecified
	}
	var modelUUID string
	if modelName != "" {
		modelDetails, err := store.ModelByName(controllerName, modelName)
		if err != nil {
			return juju.NewAPIConnectionParams{}, err
		}
		modelUUID = modelDetails.ModelUUID
	}
	dialOpts := api.DefaultDialOpts()
	return juju.NewAPIConnectionParams{
		Store:          store,
		ControllerName: controllerName,
		AccountDetails: accountDetails,
		ModelUUID:      modelUUID,
		DialOpts:       dialOpts,
		OpenAPI:        modelcmd.OpenAPIFuncWithMacaroons(apiOpen, store, controllerName),
	}, nil
}
