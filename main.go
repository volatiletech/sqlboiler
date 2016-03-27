/*
SQLBoiler is a tool to generate Go boilerplate code for database interactions.
So far this includes struct definitions and database statement helper functions.
*/

package main

import (
	"fmt"
	"os"

	"github.com/pobri19/sqlboiler/cmds"
	"github.com/spf13/cobra"
)

func main() {
	var err error
	cmdData := &cmds.CmdData{}

	// Load the "config.toml" global config
	err = cmdData.LoadConfigFile("config.toml")
	if err != nil {
		fmt.Printf("Failed to load config file: %s\n", err)
		os.Exit(-1)
	}

	// Load all templates
	err = cmdData.LoadTemplates()
	if err != nil {
		fmt.Printf("Failed to load templates: %s\n", err)
		os.Exit(-1)
	}

	// Set up the cobra root command
	var rootCmd = &cobra.Command{
		Use:   "sqlboiler",
		Short: "SQL Boiler generates boilerplate structs and statements",
		Long: "SQL Boiler generates boilerplate structs and statements from the template files.\n" +
			`Complete documentation is available at http://github.com/pobri19/sqlboiler`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cmdData.SQLBoilerPreRun(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdData.SQLBoilerRun(cmd, args)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return cmdData.SQLBoilerPostRun(cmd, args)
		},
	}

	// Set up the cobra root command flags
	rootCmd.PersistentFlags().StringP("driver", "d", "", "The name of the driver in your config.toml (mandatory)")
	rootCmd.PersistentFlags().StringP("table", "t", "", "A comma seperated list of table names")
	rootCmd.PersistentFlags().StringP("folder", "f", "output", "The name of the output folder")
	rootCmd.PersistentFlags().StringP("pkgname", "p", "model", "The name you wish to assign to your generated package")

	// Execute SQLBoiler
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
