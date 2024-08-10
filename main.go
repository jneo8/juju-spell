package main

import (
	"github.com/jneo8/jujuspell/cmd"
)

func main() {
	execFunc := cmd.InitializeExecute()
	execFunc()
}
