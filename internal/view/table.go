package view

import (
	"context"

	"github.com/jneo8/jujuspell/internal/client"
	"github.com/jneo8/jujuspell/internal/ui"
)

type Table struct {
	*ui.Table
	app *App
}

func NewTable(resource client.Resource) *Table {
	t := Table{
		Table: ui.NewTable(resource),
	}
	return &t
}

func (t *Table) Init(ctx context.Context) error {
	t.Table.Init(ctx)
	return nil
}
