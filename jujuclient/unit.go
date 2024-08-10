package jujuclient

import (
	"strconv"
	"strings"

	"github.com/juju/juju/rpc/params"
)

type UnitData struct {
	FullStatus *params.FullStatus
}

func (d *UnitData) GetContentTableData() [][]string {
	data := [][]string{}
	columns := []string{"Unit", "leader", "Workload", "Agent", "Machine", "Public address", "Ports", "Workload Message"}
	data = append(data, columns)
	for _, appStatus := range d.FullStatus.Applications {
		for unitName, unitStatus := range appStatus.Units {
			data = append(
				data,
				[]string{
					unitName,
					strconv.FormatBool(unitStatus.Leader),
					unitStatus.WorkloadStatus.Status,
					unitStatus.AgentStatus.Status,
					unitStatus.Machine,
					unitStatus.PublicAddress,
					strings.Join(unitStatus.OpenedPorts[:], ","),
					unitStatus.WorkloadStatus.Info,
				},
			)
			for subordinateName, subordinateStatus := range unitStatus.Subordinates {
				data = append(
					data,
					[]string{
						subordinateName,
						strconv.FormatBool(subordinateStatus.Leader),
						subordinateStatus.WorkloadStatus.Status,
						subordinateStatus.AgentStatus.Status,
						unitStatus.Machine,
						subordinateStatus.PublicAddress,
						strings.Join(subordinateStatus.OpenedPorts[:], ","),
						subordinateStatus.WorkloadStatus.Info,
					},
				)
			}
		}
	}
	return data
}
