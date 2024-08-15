package cmd

import (
	"fmt"
	"os"

	"github.com/jneo8/jujuspell/internal/config"
	"github.com/jneo8/jujuspell/internal/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	appName = config.AppName
)

var (
	flags   *config.Flags
	rootCmd = &cobra.Command{
		Use:  appName,
		RunE: run,
	}
)

func init() {
	if err := config.InitLogLoc(); err != nil {
		fmt.Printf("Fail to init logs location %s\n", err)
	}
	initFlags()
}

func run(cmd *cobra.Command, args []string) error {
	file, err := os.OpenFile(
		config.AppLogFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		config.DefaultFileMod,
	)
	if err != nil {
		return fmt.Errorf("Log file %q init failed: %w", config.AppLogFile, err)
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	config.InitLogger()
	log.SetOutput(file)

	cfg, err := loadConfiguration()
	app := view.NewApp(cfg)

	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func loadConfiguration() (*config.Config, error) {
	config := config.NewConfig(flags)
	return config, nil
}

func initFlags() {
	flags = config.NewFlags()
	rootCmd.Flags().StringVarP(
		flags.LogLevel,
		"logLevel", "l",
		config.DefaultLogLevel,
		"Specify a log level (info, warn, debug, trace, error)",
	)
}
