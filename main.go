package main

import (
	"os"

	"github.com/jneo8/juju-spell/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
