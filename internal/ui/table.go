package ui

import (
	"context"

	"github.com/derailed/tview"
	"github.com/jneo8/jujuspell/internal/client"
	"github.com/jneo8/jujuspell/internal/model"
)

type Table struct {
	resource client.Resource
	*SelectTable
}

func NewTable(resource client.Resource) *Table {
	return &Table{
		SelectTable: &SelectTable{
			Table: tview.NewTable(),
			model: model.NewTable(resource),
		},
		resource: resource,
	}
}

func (t *Table) Init(ctx context.Context) {
	t.SetBorder(true)
	t.SetTitle("Table")
}
