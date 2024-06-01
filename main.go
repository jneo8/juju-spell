package main

import (
	"github.com/jneo8/juju-spell/internal/cmd"
)

func main() {
	execFunc := cmd.InitializeExecute()
	execFunc()
}
