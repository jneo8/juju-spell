package ui

import "github.com/derailed/tview"

type SelectTable struct {
	*tview.Table
	model Tabular
}
