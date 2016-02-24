package main

import (
	"os"

	"github.com/pobri19/sqlboiler/cmds"
)

func main() {
	if err := cmds.SQLBoiler.Execute(); err != nil {
		os.Exit(-1)
	}
}
