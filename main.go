/*
SQLBoiler is a tool to generate Go boilerplate code for database interactions.
So far this includes struct definitions and database statement helper functions.
*/

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
