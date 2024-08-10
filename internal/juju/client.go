package juju

import (
	"github.com/juju/errors"
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/base"
	apiclient "github.com/juju/juju/api/client/client"
	"github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/juju/jujuclient"
	"github.com/juju/names/v5"
	"github.com/sirupsen/logrus"
)

type JujuClientStore interface {
	CurrentController() (string, error)
	SetCurrentController(controllerName string) error
	SetCurrentModel(controllerName, modelName string) error
}

type JujuClient interface {
	GetControllerData() (ControllerData, error)
	GetModelData(controllerName string) (ModelData, error)
	GetUnitData(controllerName string, modelName string) (UnitData, error)
	JujuClientStore
}

type jujuClient struct {
	apiConnection api.Connection
	logger        *logrus.Logger
	clientStore   jujuclient.ClientStore
}

func NewJujuClient(clientStore jujuclient.ClientStore, logger *logrus.Logger) (JujuClient, error) {
	return &jujuClient{
		logger:      logger,
		clientStore: clientStore,
	}, nil
}

func (jc *jujuClient) CurrentController() (string, error) {
	return jc.clientStore.CurrentController()
}

func (jc *jujuClient) SetCurrentController(controllerName string) error {
	return jc.clientStore.SetCurrentController(controllerName)
}

func (jc *jujuClient) SetCurrentModel(controllerName, modelName string) error {
	return jc.clientStore.SetCurrentModel(controllerName, modelName)
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

func (jc *jujuClient) GetModelData(controllerName string) (ModelData, error) {
	data := ModelData{}

	accountDetails, err := jc.clientStore.AccountDetails(controllerName)
	if err != nil {
		jc.logger.Error(err)
		return data, err
	}

	apiRoot, err := jc.GetAPIConnection(controllerName, "")
	if err != nil {
		jc.logger.Error(err)
		return data, err
	}
	defer apiRoot.Close()

	modelManagerClient := modelmanager.NewClient(apiRoot)
	summaries, err := modelManagerClient.ListModelSummaries(accountDetails.User, true)
	if err != nil {
		jc.logger.Error(err)
		return data, err
	}
	data.ModelSummaries = summaries

	// Update models with ClientStore
	if err := jc.updateModels(controllerName, summaries); err != nil {
		jc.logger.Error(err)
		return data, err
	}

	currentModel, err := jc.clientStore.CurrentModel(controllerName)
	if err != nil {
		jc.logger.Warning(err)
	}
	data.CurrentModel = currentModel
	return data, nil
}

func (jc *jujuClient) updateModels(controllerName string, summaries []base.UserModelSummary) error {
	modelsToStore := map[string]jujuclient.ModelDetails{}
	for _, summary := range summaries {
		name := jujuclient.JoinOwnerModelName(names.NewUserTag(summary.Owner), summary.Name)
		modelsToStore[name] = jujuclient.ModelDetails{ModelUUID: summary.UUID, ModelType: summary.Type}
	}
	if err := jc.clientStore.SetModels(controllerName, modelsToStore); err != nil {
		return err
	}
	return nil
}

func (jc *jujuClient) GetUnitData(controllerName, modelName string) (UnitData, error) {
	data := UnitData{}
	apiRoot, err := jc.GetAPIConnection(controllerName, modelName)
	if err != nil {
		jc.logger.Error(err)
		return data, err
	}
	defer apiRoot.Close()

	args := apiclient.StatusArgs{}
	client := apiclient.NewClient(apiRoot, jc.logger)
	status, err := client.Status(&args)
	if err != nil {
		return data, err
	}
	data.FullStatus = status
	return data, nil
}
