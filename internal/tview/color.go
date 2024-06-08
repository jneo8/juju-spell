package tview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	ColumnColor            = tcell.ColorYellow
	DataColor              = tcell.ColorWhite
	CurrentControllerColor = tcell.ColorGreen
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
