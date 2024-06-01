package tview

import (
	"github.com/rivo/tview"
)

func NewApplication() *tview.Application {
	return tview.NewApplication()
}

type Layout struct {
	RootFlex   *tview.Flex
	HeaderBox  *tview.Box
	ContentBox *tview.Box
	FooterBox  *tview.Box
}

func GetLayout() *Layout {
	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	headerBox := tview.NewBox().SetBorder(true).SetTitle("Controller info in header")
	contentBox := tview.NewBox().SetBorder(true).SetTitle("Controller info in header")
	footerBox := tview.NewBox().SetBorder(true).SetTitle("Footer")
	rootFlex.
		AddItem(headerBox, 0, 2, true).
		AddItem(contentBox, 0, 7, true).
		AddItem(footerBox, 0, 1, true)
	layout := Layout{
		RootFlex:   rootFlex,
		HeaderBox:  headerBox,
		ContentBox: contentBox,
		FooterBox:  footerBox,
	}
	return &layout
}
