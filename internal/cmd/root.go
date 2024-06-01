package cmd

import (
	"os"

	"github.com/jneo8/juju-spell/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.cobra.yaml)")
}

var rootCmd = &cobra.Command{
	Use:   "juju-spell",
	Short: "Juju Spell",
	Long:  "This is a sample Cobra CLI application",
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := app.InitializeRootApp()
		if err != nil {
			return err
		}
		if err := app.BindFlags(cmd); err != nil {
			return err
		}
		return app.Execute()
	},
}

func GetExecute(logger *logrus.Logger) func() {
	return func() {
		if err := rootCmd.Execute(); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	}
}
