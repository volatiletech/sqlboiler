package main

import (
	"os"

	"github.com/pobri19/sqlboiler/cmds"
)

func main() {
	// Execute SQLBoiler
	if err := cmds.SQLBoiler.Execute(); err != nil {
		os.Exit(-1)
	}
}
