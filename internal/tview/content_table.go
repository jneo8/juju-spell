package tview

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

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

func (s *Service) SwitchToModelTable(controllerName string) {
	s.ContentDataTable.Clear()
	controllerName = strings.ReplaceAll(controllerName, "*", "")
	modelData, err := s.jujuClient.GetModelData(controllerName)
	if err != nil {
		s.logger.Error(err)
	}

	data := modelData.GetModelTableData()
	currentModelName := ""
	if modelData.CurrentModel != "" {
		currentModelName = strings.Split(modelData.CurrentModel, "/")[1]
	}
	s.drawContentTable(currentModelName, fmt.Sprintf("Models(%s)", controllerName), data)
}

func (s *Service) SwitchToControllerTable() {
	s.logger.Info("Get controller")
	controllerData, err := s.jujuClient.GetControllerData()

	if err != nil {
		s.logger.Error(err)
		s.Error(fmt.Sprint(err))
		return
	}
	for _, err := range controllerData.Errors {
		s.Error(err.Error())
	}
	data := controllerData.GetControllerTableData()
	s.drawContentTable(controllerData.CurrentController, "Controllers", data)
}
