package juju

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/juju/api/controller/controller"
	"github.com/juju/juju/jujuclient"
	"github.com/juju/juju/rpc/params"
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
	GetModelData(ControllerName string) (ModelData, error)
	GetUnitData(ControllerName string) (UnitData, error)
	Watcher() error
	JujuClientStore
}

type jujuClient struct {
	apiConnection api.Connection
	logger        *logrus.Logger
	clientStore   jujuclient.ClientStore
	states        map[string]state
}

func NewJujuClient(clientStore jujuclient.ClientStore, logger *logrus.Logger) (JujuClient, error) {
	return &jujuClient{
		logger:      logger,
		clientStore: clientStore,
		states:      make(map[string]state),
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

func (jc *jujuClient) GetUnitData(modelName string) (UnitData, error) {
	data := UnitData{}
	return data, nil
}

func (jc *jujuClient) Watch() error {
	return nil
}

func (jc *jujuClient) Watcher() error {
	apiRoot, err := jc.GetAPIConnection("local-3.4", "")
	if err != nil {
		jc.logger.Error(err)
		return nil
	}
	client := controller.NewClient(apiRoot)
	watcher, err := client.WatchAllModels()
	if err != nil {
		jc.logger.Error(err)
		return nil
	}
	jc.logger.Info(watcher)
	for {
		time.Sleep(1)
		deltas, err := watcher.Next()
		if err != nil {
			jc.logger.Error(err)
		}
		for _, delta := range deltas {
			// jc.logger.Infof("Delta: %#v", delta)
			// jc.logger.Infof("Entity: %#v", delta.Entity)
			// jc.logger.Infof("EntityId: %v", delta.Entity.EntityId())
			entityId := delta.Entity.EntityId()
			jc.logger.Infof("EntityId: %#v", entityId)
			switch entityId.Kind {
			case "action":
				jc.logger.Info("action")
				actionInfo := delta.Entity.(*params.ActionInfo)
				jc.logger.Infof("%#v", actionInfo.ModelUUID)
			case "annotation":
				jc.logger.Info("annotation")
			case "application":
				jc.logger.Info("application")
				applicationInfo := delta.Entity.(*params.ApplicationInfo)
				jc.logger.Infof("%#v", applicationInfo)
			case "applicationOffer":
				jc.logger.Info("applicationOffer")
			case "block":
				jc.logger.Info("block")
			case "branch":
				jc.logger.Info("branch")
			case "charm":
				jc.logger.Info("charm")
			case "machine":
				jc.logger.Info("machine")
			case "model":
				jc.logger.Info("model")
			case "relation":
				jc.logger.Info("relation")
			case "remoteApplication":
				jc.logger.Info("remoteApplication")
			case "unit":
				jc.logger.Info("unit")
			default:
				jc.logger.Errorf("Unexcepted entity name %q", entityId.Kind)
			}
		}
	}
}
