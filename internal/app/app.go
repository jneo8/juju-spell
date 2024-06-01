package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jneo8/juju-spell/internal/tview"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ExecuteAble interface {
	Execute() error
	BindFlags(cmd *cobra.Command) error
}

type App struct {
	logger *logrus.Logger
	config *viper.Viper
}

func (app *App) BindFlags(cmd *cobra.Command) error {
	app.config.BindPFlags(cmd.Flags())
	app.logger.Info(config.AllSettings())
	return nil
}

func NewApp(logger *logrus.Logger, config *viper.Viper) ExecuteAble {
	logger.Info("NewApp")
	return &App{
		logger: logger,
		config: config,
	}
}

func (app *App) Execute() error {
	app.logger.Info("Execute")
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(2)
	errChan := make(chan error, 2)

	tviewService := tview.GetService()
	go tviewService.Run(ctx, &wg, errChan)
	go RunDummyService(ctx, &wg, errChan, tviewService)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case err := <-errChan:
		app.logger.Error(err)
		cancel()
	case <-ctx.Done():
		logger.Info("Done")
	}
	return nil
}

func RunDummyService(ctx context.Context, wg *sync.WaitGroup, errChan chan<- error, service tview.ViewService) {
	defer wg.Done()
	go func() {
		for i := range 10000 {
			time.Sleep(1 * time.Second)
			service.Info(fmt.Sprintf("Info: %d", i))
		}
	}()
	select {
	case <-ctx.Done():
		return
	}
}
