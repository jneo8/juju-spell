package juju

import (
	"fmt"
	"strconv"

	jujucontroller "github.com/juju/juju/cmd/juju/controller"
)

type ControllerData struct {
	ControllerItems map[string]jujucontroller.ControllerItem
	errors          []error
}

func (c *ControllerData) GetControllerTableData() [][]string {
	data := [][]string{}
	columns := []string{"Controller", "Model", "User", "Access", "Cloud/Region", "Models", "Nodes", "HA", "Version"}
	data = append(data, columns)

	idx := 0
	for ctrlName, ctrl := range c.ControllerItems {
		idx++
		ha := "none"
		if ctrl.ControllerMachines != nil && ctrl.ControllerMachines.Total > 1 {
			ha = "yes"
		}
		data = append(data, []string{
			ctrlName,
			ctrl.ModelName,
			ctrl.User,
			ctrl.Access,
			fmt.Sprintf("%s/%s", ctrl.Cloud, ctrl.CloudRegion),
			strconv.Itoa(*ctrl.ModelCount),
			strconv.Itoa(*ctrl.MachineCount),
			ha,
			ctrl.AgentVersion,
		})
	}
	return data
}
