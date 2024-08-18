package view

import (
	"github.com/jneo8/jujuspell/internal/client"
	"github.com/jneo8/jujuspell/internal/model"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	app *App
}

func NewCommand(app *App) *Command {
	return &Command{
		app: app,
	}
}

func (c *Command) runDefaultCommand() error {
	log.Debug("run default command")
	c.run(true)
	return nil
}

func (c *Command) run(clearStack bool) error {
	resource := client.Resource{
		Controller: "test-controller-1",
		Model:      "test-model-1",
		Name:       "application",
	}
	comp := c.componentFor(resource)
	return c.exec(comp, clearStack)
}

func (c *Command) exec(comp model.Component, clearStack bool) (err error) {
	if err := c.app.inject(comp, clearStack); err != nil {
		return err
	}
	return
}

func (c *Command) componentFor(resource client.Resource) ResourceReviewer {
	browser := NewBrowser(resource)
	return browser
}
