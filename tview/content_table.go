package tview

import (
	"strings"

	"github.com/rivo/tview"
)

func (s *Service) ContentDataTableSelectedFunc(row int, column int) {
	title := s.ContentFlex.GetTitle()
	switch {
	case title == "Controllers":
		s.RunOps(
			Operations{
				Name: "SwitchToUnitTableOps",
				Ops: []Operation{
					s.getContentDataTableProcessingOp(),
					s.getSwitchToModelTableOp(row, column),
				},
			},
		)
	case strings.HasPrefix(title, "Models"):
		s.RunOps(
			Operations{
				Name: "SwitchToUnitTableOps",
				Ops: []Operation{
					s.getContentDataTableProcessingOp(),
					s.getSwitchToUnitTableOp(row, column),
				},
			},
		)
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
