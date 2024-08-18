package view

import (
	"context"

	"github.com/jneo8/jujuspell/internal/client"
)

type Browser struct {
	*Table
	resource client.Resource
}

func NewBrowser(resource client.Resource) ResourceReviewer {
	return &Browser{
		Table:    NewTable(resource),
		resource: resource,
	}
}

func (b *Browser) Init(ctx context.Context) error {
	if err := b.Table.Init(ctx); err != nil {
		return err
	}
	return nil
}
func (b *Browser) Name() string {
	return b.resource.ResourceName()
}
func (b *Browser) Start() {
}
func (b *Browser) Stop() {
}
