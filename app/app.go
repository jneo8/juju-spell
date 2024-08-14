package app

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jneo8/jujuspell/common"
	"github.com/jneo8/jujuspell/jujuclient"
	"github.com/jneo8/jujuspell/tview"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ExecuteAble interface {
	Execute() error
	Setup(cmd *cobra.Command) error
	Close() error
}

type App struct {
	logger       *logrus.Logger
	config       *viper.Viper
	jujuClient   jujuclient.JujuClient
	tviewService tview.Service
	logFile      *os.File
}

func (app *App) Setup(cmd *cobra.Command) error {
	app.config.BindPFlags(cmd.Flags())

	// Create log file
	if app.config.GetString("log_file") != "" {
		logFile, err := os.OpenFile(
			app.config.GetString("log_file"),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0666,
		)
		if err != nil {
			return err
		}
		app.logFile = logFile
	}
	common.SetupLogger(app.logger, app.config, app.logFile)
	app.logger.Info(app.config.AllSettings())
	return nil
}

func (app *App) Close() error {
	if app.logFile != nil {
		defer app.logFile.Close()
	}
	return nil
}

func NewApp(logger *logrus.Logger, config *viper.Viper, jujuClient jujuclient.JujuClient) ExecuteAble {
	return &App{
		logger:     logger,
		config:     config,
		jujuClient: jujuClient,
	}
}

func (app *App) Execute() error {
	app.logger.Debug("Execute")
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(2)
	errChan := make(chan error, 2)

	tviewService := tview.GetService(app.logger, app.jujuClient)
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
		app.logger.Info("Done")
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
