package juju

import (
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/connector"
	"github.com/juju/juju/jujuclient"
	"github.com/sirupsen/logrus"
)

type JujuClient interface {
	Connect() error
	GetContent() JujuContent
	GetControllers() jujuclient.Controllers
}

type jujuClient struct {
	clientStoreConnector *connector.ClientStoreConnector
	clientStoreConfig    connector.ClientStoreConfig
	apiConnection        api.Connection
	logger               *logrus.Logger
	content              *JujuContent
}

func NewJujuClient(jujuContent *JujuContent, logger *logrus.Logger) (JujuClient, error) {
	clientStoreConfig := connector.ClientStoreConfig{}
	return &jujuClient{
		clientStoreConfig: clientStoreConfig,
		logger:            logger,
		content:           jujuContent,
	}, nil
}

func (jc *jujuClient) Connect() error {
	conn, err := jc.clientStoreConnector.Connect()
	if err != nil {
		return err
	}
	jc.apiConnection = conn
	return nil
}

func (jc *jujuClient) GetContent() JujuContent {
	return *jc.content
}

func (jc *jujuClient) GetControllers() jujuclient.Controllers {
	return *jc.content.Controllers
}
