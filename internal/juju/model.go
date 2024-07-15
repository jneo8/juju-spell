package juju

import (
	"fmt"
	"strconv"

	"github.com/juju/juju/api/base"
)

type ModelData struct {
	ModelSummaries []base.UserModelSummary
	CurrentModel   string
}

func (m *ModelData) GetModelTableData() [][]string {
	data := [][]string{}
	columns := []string{"Model", "Cloud/Region", "Type", "Status", "Machines", "Units", "Access", "Last connection"}

	data = append(data, columns)

	for idx, model := range m.ModelSummaries {
		idx++
		data = append(
			data,
			[]string{
				model.Name,
				fmt.Sprintf("%s/%s", model.Cloud, model.CloudRegion),
				model.ProviderType,
				model.Status.Status.String(),
				strconv.FormatInt(model.Counts[0].Count, 10),
				strconv.FormatInt(model.Counts[1].Count, 10),
				model.ModelUserAccess,
				model.UserLastConnection.String(),
			},
		)
	}
	return data
}
