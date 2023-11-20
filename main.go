package main

import (
	"os"

	"github.com/jneo8/juju-man/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
