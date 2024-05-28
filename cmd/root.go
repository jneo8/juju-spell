package cmd

import (
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
	return nil
}

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: run,
}

func init() {
	rootCmd.AddCommand(jujudryCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
