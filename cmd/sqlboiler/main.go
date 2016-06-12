// Package main defines a command line interface for the sqlboiler package
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/nullbio/sqlboiler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var err error

	viper.SetConfigName("sqlboiler")
	viper.AddConfigPath("$HOME/.sqlboiler")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to load config file: %s\n", err)
		os.Exit(-1)
	}

	// Set up the cobra root command
	var rootCmd = &cobra.Command{
		Use:   "sqlboiler",
		Short: "SQL Boiler generates boilerplate structs and statements",
		Long: "SQL Boiler generates boilerplate structs and statements from the template files.\n" +
			`Complete documentation is available at http://github.com/nullbio/sqlboiler`,
		PreRunE:  preRun,
		RunE:     run,
		PostRunE: postRun,
	}

	// Set up the cobra root command flags
	rootCmd.PersistentFlags().StringP("driver", "d", "", "The name of the driver in your config.toml (mandatory)")
	rootCmd.PersistentFlags().StringP("table", "t", "", "A comma seperated list of table names")
	rootCmd.PersistentFlags().StringP("folder", "f", "output", "The name of the output folder")
	rootCmd.PersistentFlags().StringP("pkgname", "p", "model", "The name you wish to assign to your generated package")

	viper.BindPFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Failed to execute sqlboiler command:", err)
		os.Exit(-1)
	}
}

var state *sqlboiler.State
var config *sqlboiler.Config

func preRun(cmd *cobra.Command, args []string) error {
	config = new(sqlboiler.Config)

	config.DriverName = viper.GetString("driver")
	config.TableName = viper.GetString("table")
	config.OutFolder = viper.GetString("folder")
	config.PkgName = viper.GetString("pkgname")

	if len(config.DriverName) == 0 {
		return errors.New("Must supply a driver flag.")
	}
	if len(config.OutFolder) == 0 {
		return fmt.Errorf("No output folder specified.")
	}

	if viper.IsSet("postgres.dbname") {
		config.Postgres = sqlboiler.PostgresConfig{
			User:   viper.GetString("postgres.user"),
			Pass:   viper.GetString("postgres.pass"),
			Host:   viper.GetString("postgres.host"),
			Port:   viper.GetInt("postgres.port"),
			DBName: viper.GetString("postgres.dbname"),
		}
	}

	var err error
	state, err = sqlboiler.New(config)
	return err
}

func run(cmd *cobra.Command, args []string) error {
	return state.Run(true)
}

func postRun(cmd *cobra.Command, args []string) error {
	return state.Cleanup()
}
