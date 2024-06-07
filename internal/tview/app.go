package tview

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/jneo8/juju-spell/internal/juju"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

const logo = `
    _       _                       _ _ 
   (_)_  _ (_)_  _ ___ ____ __  ___| | |
   | | || || | || |___(_-< '_ \/ -_) | |
  _/ |\_,_|/ |\_,_|   /__/ .__/\___|_|_|
 |__/    |__/            |_|            
`

func getTheme() tview.Theme {
	return tview.Theme{
		PrimitiveBackgroundColor:    tcell.Color16,
		ContrastBackgroundColor:     tcell.ColorBlack,
		MoreContrastBackgroundColor: tcell.ColorBlack,
		BorderColor:                 tcell.ColorDarkCyan,
		TitleColor:                  tcell.ColorDeepSkyBlue,
		GraphicsColor:               tcell.ColorDeepSkyBlue,
		PrimaryTextColor:            tcell.ColorWhite,
		SecondaryTextColor:          tcell.ColorMediumPurple,
		TertiaryTextColor:           tcell.ColorDeepSkyBlue,
		InverseTextColor:            tcell.ColorBlack,
		ContrastSecondaryTextColor:  tcell.ColorDeepPink,
	}
}

type ViewService interface {
	Run(ctx context.Context, wg *sync.WaitGroup, errChan chan<- error)
	Info(string)
	Debug(string)
}

func NewApplication() *tview.Application {
	return tview.NewApplication()
}

type Service struct {
	Application      *tview.Application
	RootFlex         *tview.Flex
	HeaderFlex       *tview.Flex
	HeaderList       *tview.List
	ContentFlex      *tview.Flex
	FooterFlex       *tview.Flex
	LogTextView      *tview.TextView
	ContentDataTable *tview.Table
	jujuClient       juju.JujuClient
	logger           *logrus.Logger
}

func GetService(logger *logrus.Logger, jujuclient juju.JujuClient) ViewService {
	tview.Styles = getTheme()
	app := tview.NewApplication()

	// HeaderFlex
	headerFlex := tview.NewFlex()
	logoTextView := tview.NewTextView()
	logoTextView.SetText(logo)

	headerList := tview.NewList()

	headerFlex.
		AddItem(headerList, 0, 70, false).
		AddItem(logoTextView, 40, 30, false)
	// End HeaderFlex

	// Content
	contentFlex := tview.NewFlex()
	contentFlex.SetBorder(true).SetTitle("Controller info in header")
	dataTable := tview.NewTable()
	contentFlex.AddItem(dataTable, 0, 100, false)
	// End Content

	// Footer
	footerFlex := tview.NewFlex()
	footerFlex.SetBorder(true)
	footerFlex.SetTitle("Log message")

	logTextView := tview.
		NewTextView().
		SetMaxLines(1000).
		SetScrollable(true).
		SetChangedFunc(func() { app.Draw() }).
		SetDynamicColors(true)

	footerFlex.AddItem(logTextView, 0, 1, true)
	// End Footer

	// Root
	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rootFlex.
		AddItem(headerFlex, 0, 15, true).
		AddItem(contentFlex, 0, 80, true).
		AddItem(footerFlex, 0, 5, true)
	// End Root

	service := Service{
		logger:           logger,
		jujuClient:       jujuclient,
		Application:      app,
		RootFlex:         rootFlex,
		HeaderFlex:       headerFlex,
		HeaderList:       headerList,
		ContentFlex:      contentFlex,
		FooterFlex:       footerFlex,
		LogTextView:      logTextView,
		ContentDataTable: dataTable,
	}
	service.setUpHeaderItem()
	return &service
}

func (s *Service) setUpHeaderItem() {
	s.HeaderList.AddItem("Controllers", "", 'c', s.SwitchToControllerTable)
	s.HeaderList.AddItem("Models", "", 'm', nil)
	s.HeaderList.AddItem("Units", "", 'u', nil)
	s.HeaderList.AddItem("Integrations", "", 'i', nil)
}

func (s *Service) SwitchToControllerTable() {
	s.logger.Info("Get controller")
	controllers := s.jujuClient.GetControllers()
	for idx, colume := range []string{"Controller", "Model", "User", "Access", "Cloud/Region", "Models", "Nodes", "HA", "Version"} {
		color := tcell.ColorYellow
		align := tview.AlignCenter
		tableCell := tview.NewTableCell(colume).SetAlign(align).SetTextColor(color).SetSelectable(false)
		s.ContentDataTable.SetCell(0, idx, tableCell)
	}
	for ctrlName, ctrl := range controllers.Controllers {
		color := tcell.ColorWhite
		align := tview.AlignRight
		controllerCell := tview.NewTableCell(ctrlName).SetAlign(align).SetTextColor(color).SetSelectable(false)
		cloudRegionCell := tview.NewTableCell(fmt.Sprintf("%s/%s", ctrl.Cloud, ctrl.CloudRegion)).SetAlign(align).SetTextColor(color).SetSelectable(false)
		nodesCell := tview.NewTableCell(strconv.Itoa(*ctrl.MachineCount)).SetAlign(align).SetTextColor(color).SetSelectable(false)
		versionCell := tview.NewTableCell(ctrl.AgentVersion).SetAlign(align).SetTextColor(color).SetSelectable(false)
		s.ContentDataTable.SetCell(1, 0, controllerCell)
		s.ContentDataTable.SetCell(1, 1, controllerCell)
		s.ContentDataTable.SetCell(1, 4, cloudRegionCell)
		s.ContentDataTable.SetCell(1, 6, nodesCell)
		s.ContentDataTable.SetCell(1, 8, versionCell)
		s.logger.Debugf("%s %#v", ctrlName, ctrl)
	}
}

func (s *Service) Run(ctx context.Context, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	go func() {
		if err := s.Application.SetRoot(s.RootFlex, true).EnableMouse(true).EnablePaste(true).SetFocus(s.HeaderList).Run(); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()
	select {
	case <-ctx.Done():
		s.Application.Stop()
		return
	}
}

func (s *Service) Debug(str string) {
	color := ""
	fmt.Fprintf(s.LogTextView, "%s%s\n", color, str)
	if focus := s.LogTextView.HasFocus(); !focus {
		s.LogTextView.ScrollToEnd()
	}
}

func (s *Service) Info(str string) {
	color := "[blue]"
	fmt.Fprintf(s.LogTextView, "\n%s%s[-]", color, str)
	if focus := s.LogTextView.HasFocus(); !focus {
		s.LogTextView.ScrollToEnd()
	}
}
