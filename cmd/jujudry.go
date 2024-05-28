/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/juju/juju/api/client/client"
	"github.com/juju/juju/api/connector"
	"github.com/spf13/cobra"
)

type Logger struct{}

func (l Logger) Errorf(s string, args ...interface{}) {
	log.Println(args)
	log.Println(s)
}

// jujudryCmd represents the jujudry command
var jujudryCmd = &cobra.Command{
	Use:   "jujudry",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// connr, err := connector.NewClientStore(connector.ClientStoreConfig{
		// 	ControllerName: "local-lxd-3",
		// })
		// modelcmdbase := modelcmd.ModelCommandBase{}
		// modelcmdbase.SetClientStore(connr)
		// api, err := modelcmdbase.NewAPIClient()
		// if err != nil {
		// 	return err
		// }
		// status, err := api.Status([]string{})
		// if err != nil {
		// 	return err
		// }
		// log.Println(status)
		// return nil

		connr, err := connector.NewClientStore(connector.ClientStoreConfig{
			ControllerName: "local-lxd-3",
		})
		log.Printf("%#v\n", connr)

		conn, err := connr.Connect()
		if err != nil {
			log.Fatalf("Error opening connection: %s", err)
		}
		defer conn.Close()

		logger := Logger{}

		client := client.NewClient(conn, logger)
		defer client.Close()

		status, err := client.Status(nil)
		if err != nil {
			log.Fatalf("Error requesting status: %s", err)
		}

		b, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling response: %s", err)
		}
		fmt.Printf("%s\n", b)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(jujudryCmd)
}
