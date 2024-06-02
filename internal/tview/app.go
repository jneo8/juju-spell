package tview

import (
	"context"
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	Application *tview.Application
	RootFlex    *tview.Flex
	HeaderFlex  *tview.Flex
	ContentFlex *tview.Flex
	FooterFlex  *tview.Flex
	LogTextView *tview.TextView
}

func GetService() ViewService {
	tview.Styles = getTheme()
	app := tview.NewApplication()

	// HeaderFlex
	headerFlex := tview.NewFlex()
	logoTextView := tview.NewTextView()
	logoTextView.SetText(logo)

	headerList := tview.NewList()
	headerList.AddItem("Controllers", "", 'c', nil)
	headerList.AddItem("Models", "", 'm', nil)
	headerList.AddItem("Units", "", 'u', nil)
	headerList.AddItem("Integrations", "", 'i', nil)

	headerFlex.
		AddItem(headerList, 0, 70, false).
		AddItem(logoTextView, 40, 30, false)
	// End HeaderFlex

	// Content
	contentFlex := tview.NewFlex()
	contentFlex.SetBorder(true).SetTitle("Controller info in header")
	DataTable := tview.NewTable()
	contentFlex.AddItem(DataTable, 0, 100, false)
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
		Application: app,
		RootFlex:    rootFlex,
		HeaderFlex:  headerFlex,
		ContentFlex: contentFlex,
		FooterFlex:  footerFlex,
		LogTextView: logTextView,
	}
	return &service
}

func (s *Service) Run(ctx context.Context, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	go func() {
		if err := s.Application.SetRoot(s.RootFlex, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
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
