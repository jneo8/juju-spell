package juju

import (
	"fmt"
	"strconv"
	"time"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/cmd/juju/common"
)

type ModelData struct {
	ModelSummaries []base.UserModelSummary
	CurrentModel   string
}

func (m *ModelData) GetModelTableData() [][]string {
	data := [][]string{}
	columns := []string{"Model", "Cloud/Region", "Type", "Status", "Machines", "Units", "Owner", "Access", "Last connection"}

	data = append(data, columns)
	for idx, summary := range m.ModelSummaries {
		idx++
		var machineCount int64 = 0
		var unitCount int64 = 0
		for _, count := range summary.Counts {
			switch count.Entity {
			case "units":
				unitCount = count.Count
			case "machines":
				machineCount = count.Count
			}
		}
		userLastConnection := "never connected"
		if summary.UserLastConnection != nil {
			userLastConnection = common.UserFriendlyDuration(*summary.UserLastConnection, time.Now())
		}
		data = append(
			data,
			[]string{
				summary.Name,
				fmt.Sprintf("%s/%s", summary.Cloud, summary.CloudRegion),
				summary.ProviderType,
				summary.Status.Status.String(),
				strconv.FormatInt(machineCount, 10),
				strconv.FormatInt(unitCount, 10),
				summary.Owner,
				summary.ModelUserAccess,
				userLastConnection,
			},
		)
	}
	return data
}
