package tview

import (
	"fmt"
	"strings"

	"github.com/jneo8/juju-spell/utils"
	"github.com/juju/juju/jujuclient"
	"github.com/juju/names/v5"
	"github.com/rivo/tview"
)

func (s *Service) ContentDataTableSelectedFunc(row int, column int) {
	s.logger.Debug(row, column)
	title := s.ContentFlex.GetTitle()
	switch {
	case title == "Controllers":
		cell := s.ContentDataTable.GetCell(row, column)
		controllerName := utils.RemoveWildcards(cell.Text)
		if err := s.jujuClient.SetCurrentController(controllerName); err != nil {
			s.Error(err)
			return
		}
		s.SwitchToModelTable(cell.Text)
	case strings.HasPrefix(title, "Models"):
		model := utils.RemoveWildcards(s.ContentDataTable.GetCell(row, column).Text)
		currentController, err := s.jujuClient.CurrentController()
		if err != nil {
			s.Error(err)
			return
		}
		// modelName = {owner}/{model}
		modelName := jujuclient.JoinOwnerModelName(names.NewUserTag(s.ContentDataTable.GetCell(row, 6).Text), model)
		if err := s.jujuClient.SetCurrentModel(currentController, modelName); err != nil {
			s.Error(err)
			return
		}
		s.SwitchToUnitTable(currentController, modelName)
	default:
	}
}

func (s *Service) drawContentTable(current, title string, data [][]string) {
	s.ContentDataTable.Clear()

	s.ContentFlex.SetTitle(title)
	s.ContentDataTable.SetSelectable(true, false)
	s.Application.SetFocus(s.ContentDataTable)
	if len(data) <= 0 {
		s.logger.Debug("Empty data")
		return
	}
	center := len(data[0]) / 2
	for row, line := range data {
		for column, cell := range line {
			currentModel := false
			if cell == current {
				currentModel = true
			}
			color := DataColor
			align := tview.AlignLeft
			selectable := true

			if row == 0 {
				color = ColumnColor
				selectable = false
			}
			if column > center {
				align = tview.AlignRight
			}

			ctrlName := cell
			if currentModel {
				ctrlName = "*" + ctrlName
				color = CurrentControllerColor
			}
			tableCell := tview.
				NewTableCell(ctrlName).
				SetAlign(align).
				SetTextColor(color).
				SetSelectable(selectable).
				SetExpansion(1)
			s.ContentDataTable.SetCell(row, column, tableCell)
			if currentModel {
				s.ContentDataTable.Select(row, column)
			}
		}
	}
}

func (s *Service) SwitchToControllerTable() {
	s.ContentDataTable.Clear()
	controllerData, err := s.jujuClient.GetControllerData()

	if err != nil {
		s.Error(err)
		return
	}
	for _, err := range controllerData.Errors {
		s.Error(err)
	}
	data := controllerData.GetControllerTableData()
	s.drawContentTable(controllerData.CurrentController, "Controllers", data)
}

func (s *Service) SwitchToModelTable(controllerName string) {
	controllerName = utils.RemoveWildcards(controllerName)
	s.logger.Debugf("Switch to %s controller", controllerName)

	s.ContentDataTable.Clear()
	modelData, err := s.jujuClient.GetModelData(controllerName)
	if err != nil {
		s.Error(err)
		return
	}

	data := modelData.GetModelTableData()
	currentModelName := ""
	if modelData.CurrentModel != "" {
		currentModelName = strings.Split(modelData.CurrentModel, "/")[1]
	}
	s.drawContentTable(currentModelName, fmt.Sprintf("Models(%s)", controllerName), data)
}

func (s *Service) SwitchToUnitTable(controllerName, modelName string) {

	s.logger.Debugf("Switch to %s model", modelName)

	unitData, err := s.jujuClient.GetUnitData(controllerName, modelName)
	if err != nil {
		s.Error(err)
		return
	}
	data := unitData.GetContentTableData(s.logger)
	s.drawContentTable("", fmt.Sprintf("%s/%s", controllerName, modelName), data)
}
