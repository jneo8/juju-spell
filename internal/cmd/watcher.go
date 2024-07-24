package cmd

import (
	"github.com/jneo8/juju-spell/internal/juju"
	"github.com/spf13/cobra"
)

var watcherCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watcher",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := juju.InitializeJujuClient()
		if err != nil {
			return err
		}
		client.Watcher()
		return nil
	},
}
