package tview

import (
	"context"
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
	HeaderBox   *tview.Box
	ContentBox  *tview.Box
	Footer      *tview.Flex
	LogTextView *tview.TextView
}

func GetService() ViewService {
	tview.Styles = getTheme()
	app := tview.NewApplication()
	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	headerBox := tview.NewBox().SetBorder(true).SetTitle("Controller info in header")
	contentBox := tview.NewBox().SetBorder(true).SetTitle("Controller info in header")

	footer := tview.NewFlex()
	footer.SetBorder(true)

	logTextView := tview.
		NewTextView().
		SetMaxLines(1000).
		SetScrollable(true).
		SetChangedFunc(func() { app.Draw() }).
		SetDynamicColors(true)

	footer.AddItem(logTextView, 0, 1, true)
	rootFlex.
		AddItem(headerBox, 0, 15, true).
		AddItem(contentBox, 0, 80, true).
		AddItem(footer, 0, 5, true)
	service := Service{
		Application: app,
		RootFlex:    rootFlex,
		HeaderBox:   headerBox,
		ContentBox:  contentBox,
		Footer:      footer,
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
