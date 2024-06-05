package juju

import "github.com/juju/juju/jujuclient"

type JujuContent struct {
	Controllers *jujuclient.Controllers
}

func GetControllers() (*jujuclient.Controllers, error) {
	jujuControllersPath := jujuclient.JujuControllersPath()
	return jujuclient.ReadControllersFile(jujuControllersPath)
}

func (content *JujuContent) LoadFromFiles() error {
	controllers, err := GetControllers()
	if err != nil {
		return nil
	}
	content.Controllers = controllers
	return nil
}

func GetJujuContent() (*JujuContent, error) {
	content := &JujuContent{}
	if err := content.LoadFromFiles(); err != nil {
		return nil, err
	}
	return content, nil
}
