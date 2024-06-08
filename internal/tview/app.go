package tview

import (
	"context"
	"fmt"
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
	PromptTextView   *tview.TextView
	HelperList       *tview.List
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
	// headerList.SetBorder(true)
	promptTextView := tview.NewTextView()
	promptTextView.
		SetChangedFunc(func() { app.Draw() }).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetRegions(true).
		SetBorder(true).
		SetTitle("Prompt")

	headerFlex.
		AddItem(headerList, 0, 30, false).
		AddItem(promptTextView, 0, 40, false).
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

	footerFlex.AddItem(logTextView, 0, 1, false)
	// End Footer

	// Root
	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rootFlex.
		AddItem(headerFlex, 0, 15, false).
		AddItem(contentFlex, 0, 80, false).
		AddItem(footerFlex, 0, 5, false)
	// End Root

	service := Service{
		logger:           logger,
		jujuClient:       jujuclient,
		Application:      app,
		RootFlex:         rootFlex,
		HeaderFlex:       headerFlex,
		HeaderList:       headerList,
		PromptTextView:   promptTextView,
		ContentFlex:      contentFlex,
		FooterFlex:       footerFlex,
		LogTextView:      logTextView,
		ContentDataTable: dataTable,
	}
	service.setUpHeaderList()
	service.setUpPromptTextView()
	service.setUpBasicInputCapture()

	return &service
}

func (s *Service) setUpBasicInputCapture() {
	s.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'd' {
			s.logger.Debug("Move to data table")
			s.Application.SetFocus(s.ContentDataTable)
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'h' {
			s.logger.Debug("Move to headerlist")
			s.Application.SetFocus(s.HeaderList)
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'l' {
			s.logger.Debug("Move to log view")
			s.Application.SetFocus(s.LogTextView)
		}
		return event
	})
}

func (s *Service) setUpHeaderList() {
	s.HeaderList.AddItem("Controllers", "", 'c', s.SwitchToControllerTable)
	s.HeaderList.AddItem("Models", "", 'm', nil)
	s.HeaderList.AddItem("Units", "", 'u', nil)
	s.HeaderList.AddItem("Integrations", "", 'i', nil)
}

func (s *Service) setUpPromptTextView() {
	color := "[blue]"
	fmt.Fprintf(s.PromptTextView, "%s<h>[-] Switch to header\n", color)
	fmt.Fprintf(s.PromptTextView, "%s<d>[-] Switch to data table\n", color)
	fmt.Fprintf(s.PromptTextView, "%s<l>[-] Switch to log\n", color)
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
	if len(data) <= 0 {
		return
	}
	center := len(data[0]) / 2
	for row, line := range controllerData.GetControllerTableData() {
		for column, cell := range line {
			currentController := false
			if cell == controllerData.CurrentController {
				currentController = true
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
			} else {
				align = tview.AlignLeft
			}

			ctrlName := cell
			if currentController {
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
			if currentController {
				s.ContentDataTable.Select(row, column)
			}
		}
	}
	s.ContentDataTable.SetSelectable(true, false)
	s.Application.SetFocus(s.ContentDataTable)
	s.ContentFlex.SetTitle("Controllers")
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

func (s *Service) Error(str string) {
	color := "[red]"
	fmt.Fprintf(s.LogTextView, "\n%s%s[-]", color, str)
	if focus := s.LogTextView.HasFocus(); !focus {
		s.LogTextView.ScrollToEnd()
	}
}
