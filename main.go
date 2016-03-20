/*
SQLBoiler is a tool to generate Go boilerplate code for database interactions.
So far this includes struct definitions and database statement helper functions.
*/

package main

import (
	"fmt"
	"os"

	"github.com/pobri19/sqlboiler/cmds"
)

func main() {
	// Load the config.toml file
	cmds.LoadConfigFile("config.toml")

	// Execute SQLBoiler
	if err := cmds.SQLBoiler.Execute(); err != nil {
		fmt.Printf("Failed to execute SQLBoiler: %s", err)
		os.Exit(-1)
	}
}
