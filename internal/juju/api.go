package juju

import (
	"github.com/juju/juju/api"
	"github.com/juju/names/v5"
)

func (jc *jujuClient) GetRootAPI(controllerName string) (api.Connection, error) {
	ctrl, err := jc.clientStore.ControllerByName(controllerName)
	if err != nil {
		return nil, err
	}

	accountDetails, err := jc.clientStore.AccountDetails(controllerName)
	if err != nil {
		return nil, err
	}

	info := api.Info{
		Addrs:          ctrl.APIEndpoints,
		CACert:         ctrl.CACert,
		ControllerUUID: ctrl.ControllerUUID,
		Password:       accountDetails.Password,
		Tag:            names.NewUserTag(accountDetails.User),
	}
	apiRoot, err := api.Open(&info, api.DefaultDialOpts())
	if err != nil {
		return nil, err
	}
	return apiRoot, nil
}
