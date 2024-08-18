package view

import (
	"github.com/jneo8/jujuspell/internal/model"
	"github.com/jneo8/jujuspell/internal/ui"
)

type PageStack struct {
	*ui.Pages
	app *App
}

func NewPageStack() *PageStack {
	return &PageStack{
		Pages: ui.NewPages(),
	}
}

// StackPushed notifies a new page was added.
func (p *PageStack) StackPushed(c model.Component) {
	c.Start()
	p.app.SetFocus(c)
}

// StackPopped notifies a page was removed.
func (p *PageStack) StackPopped(o, top model.Component) {
	o.Stop()
	p.StackTop(top)
}

// StackTop notifies for the top component.
func (p *PageStack) StackTop(top model.Component) {
	if top == nil {
		return
	}
	top.Start()
	p.app.SetFocus(top)
}
