package juju

import (
	"github.com/juju/errors"
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/juju/api/connector"
	"github.com/juju/juju/jujuclient"
	"github.com/sirupsen/logrus"
)

type JujuClient interface {
	GetControllerData() (ControllerData, error)
	GetModelData(string) (ModelData, error)
}

type jujuClient struct {
	clientStoreConnector *connector.ClientStoreConnector
	apiConnection        api.Connection
	logger               *logrus.Logger
	clientStore          jujuclient.ClientStore
}

func NewJujuClient(clientStore jujuclient.ClientStore, logger *logrus.Logger) (JujuClient, error) {
	return &jujuClient{
		logger:      logger,
		clientStore: clientStore,
	}, nil
}

func (jc *jujuClient) GetModelData(controllerName string) (ModelData, error) {
	data := ModelData{}

	accountDetails, err := jc.clientStore.AccountDetails(controllerName)
	if err != nil {
		return data, nil
	}

	apiRoot, err := jc.GetRootAPI(controllerName)
	if err != nil {
		return data, err
	}
	defer apiRoot.Close()
	currentModel, err := jc.clientStore.CurrentModel(controllerName)
	if err != nil {
		jc.logger.Warning(err)
	}
	data.CurrentModel = currentModel

	modelManagerClient := modelmanager.NewClient(apiRoot)
	models, err := modelManagerClient.ListModelSummaries(accountDetails.User, true)
	if err != nil {
		return data, err
	}
	data.ModelSummaries = models
	return data, nil
}

func (jc *jujuClient) GetControllerData() (ControllerData, error) {
	data := ControllerData{Errors: []error{}}
	if allControllers, err := jc.clientStore.AllControllers(); err != nil {
		return data, err
	} else {
		jc.logger.Debug("Run AllControllers success")
		controllerItems, errs := convertControllerDetails(jc.clientStore, allControllers, jc.logger)
		data.ControllerItems = controllerItems
		data.Errors = append(data.Errors, errs...)
	}
	if currentController, err := jc.clientStore.CurrentController(); err != nil {
		if errors.IsNotFound(err) {
			jc.logger.Debug("Current controller not found")
		} else {
			jc.logger.Error(err)
			data.Errors = append(data.Errors, err)
		}
	} else {
		data.CurrentController = currentController
	}
	return data, nil
}
