package tview

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jneo8/jujuspell/utils"
	"github.com/juju/juju/jujuclient"
	"github.com/juju/names/v5"
)

type Operations struct {
	Name string
	Ops  []Operation
}

type Operation struct {
	Name string
	Op   func()
}

func (s *Service) getRunOpsFunc(operations Operations) func() {
	return func() {
		s.RunOps(operations)
	}
}

func (s *Service) RunOps(operations Operations) {
	for _, operation := range operations.Ops {
		s.OperationChan <- operation
	}
}

func (s *Service) RunOperationHandler(ctx context.Context) {
	go func() {
		for {
			select {
			case operation := <-s.OperationChan:
				s.logger.Debugf("Run ops: %s", operation.Name)
				s.Application.QueueUpdateDraw(operation.Op)
			case <-ctx.Done():
				s.logger.Debug("Stop operation handler")
			}
		}
	}()
}

func (s *Service) getContentDataTableProcessingOp() Operation {
	return Operation{
		Name: "ContentDataTableProcessing",
		Op: func() {
			s.ContentFlex.SetBorderColor(tcell.ColorDeepPink)
			s.ContentFlex.SetTitle("Processing...")
		},
	}
}

func (s *Service) getSwitchToControllerTableOp() Operation {
	return Operation{
		Name: "SwitchToControllerTable",
		Op: func() {
			defer s.ContentFlex.SetBorderColor(tcell.ColorDarkCyan)
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
		},
	}
}

func (s *Service) getSwitchToModelTableOp(row, column int) Operation {
	return Operation{
		Name: "SwitchToModelTable",
		Op: func() {
			defer s.ContentFlex.SetBorderColor(tcell.ColorDarkCyan)
			s.logger.Debugf("row: %d Column: %d", row, column)
			// Get controller name
			cell := s.ContentDataTable.GetCell(row, column)
			controllerName := utils.RemoveWildcards(cell.Text)
			s.logger.Debugf("controlelrName: %s ", controllerName)
			s.logger.Debugf("cell.Text: %s ", cell.Text)
			if err := s.jujuClient.SetCurrentController(controllerName); err != nil {
				s.Error(err)
				return
			}
			controllerName = utils.RemoveWildcards(controllerName)
			s.logger.Debug(controllerName)

			s.logger.Debugf("Switch to %s controller", controllerName)
			s.ContentDataTable.Clear()

			// Get model list
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
			s.drawContentTable(currentModelName, fmt.Sprintf("Models(%s)", cell.Text), data)
		},
	}
}

func (s *Service) getSwitchToUnitTableOp(row, column int) Operation {
	return Operation{
		Name: "SwitchToUnitTable",
		Op: func() {
			defer s.ContentFlex.SetBorderColor(tcell.ColorDarkCyan)
			model := utils.RemoveWildcards(s.ContentDataTable.GetCell(row, column).Text)
			currentControllerName, err := s.jujuClient.CurrentController()
			if err != nil {
				s.Error(err)
				return
			}
			// modelName = {owner}/{model}
			modelName := jujuclient.JoinOwnerModelName(names.NewUserTag(s.ContentDataTable.GetCell(row, 6).Text), model)
			if err := s.jujuClient.SetCurrentModel(currentControllerName, modelName); err != nil {
				s.Error(err)
				return
			}
			s.logger.Debugf("Switch to %s model", modelName)

			unitData, err := s.jujuClient.GetUnitData(currentControllerName, modelName)
			if err != nil {
				s.Error(err)
				return
			}
			data := unitData.GetContentTableData()
			s.drawContentTable("", fmt.Sprintf("%s/%s", currentControllerName, modelName), data)
		},
	}
}
