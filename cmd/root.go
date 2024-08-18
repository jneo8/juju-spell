package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jneo8/jujuspell/internal/config"
	"github.com/jneo8/jujuspell/internal/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = config.AppName
)

var (
	// flags   *config.Flags
	cfgFile string
	cfgName = "config"
	cfgPath = "."
	rootCmd = &cobra.Command{
		Use:  appName,
		RunE: run,
	}
)

func init() {
	initFlags()
	initCfg()
}

func run(cmd *cobra.Command, args []string) error {
	cfg := config.NewConfig()

	// Init logger
	config.InitLogger()
	config.SetLogLevel(cfg.LogLevel)
	if err := config.InitLogLoc(cfg.LogFile); err != nil {
		fmt.Printf("Fail to init logs location %s\n", err)
	}
	file, err := os.OpenFile(
		cfg.LogFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		config.DefaultFileMod,
	)
	if err != nil {
		return fmt.Errorf("Log file %q init failed: %w", cfg.LogFile, err)
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	log.SetOutput(file)
	app := view.NewApp(cfg)

	if err := app.Init(); err != nil {
		return err
	}

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

func initFlags() {
	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")

	// logLevel
	rootCmd.Flags().String(
		"logLevel", config.DefaultLogLevel, "Specify a log level (info, warn, debug, trace, error)")
	viper.BindPFlag("logLevel", rootCmd.Flags().Lookup("logLevel"))
	// logFile
	rootCmd.Flags().String(
		"logFile", filepath.Join(config.GetXDGStateFile(), config.DefaultLogFile), "Specify the log file")
	viper.BindPFlag("logFile", rootCmd.Flags().Lookup("logFile"))
}

func initCfg() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(cfgName)
		viper.AddConfigPath(cfgPath)
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Error reading config file: %s", err)
	}
}
