// Package main defines a command line interface for the sqlboiler package
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cmdState  *State
	cmdConfig *Config
)

func main() {
	var err error

	viper.SetConfigName("sqlboiler")

	configHome := os.Getenv("XDG_CONFIG_HOME")
	homePath := os.Getenv("HOME")
	wd, err := os.Getwd()
	if err != nil {
		wd = "./"
	}

	configPaths := []string{wd}
	if len(configHome) > 0 {
		configPaths = append(configPaths, filepath.Join(configHome, "sqlboiler"))
	} else {
		configPaths = append(configPaths, filepath.Join(homePath, ".config/sqlboiler"))
	}

	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}

	// Find and read config
	err = viper.ReadInConfig()

	// Set up the cobra root command
	var rootCmd = &cobra.Command{
		Use:   "sqlboiler [options] <driver>",
		Short: "SQL Boiler generates boilerplate structs and statements",
		Long: "SQL Boiler generates boilerplate structs and statements from template files.\n" +
			`Complete documentation is available at http://github.com/nullbio/sqlboiler`,
		Example:  `sqlboiler -o mymodels -p mymodelpackage postgres`,
		PreRunE:  preRun,
		RunE:     run,
		PostRunE: postRun,
	}

	// Set up the cobra root command flags
	rootCmd.PersistentFlags().StringSliceP("table", "t", nil, "Tables to generate models for, all tables if empty")
	rootCmd.PersistentFlags().StringP("output", "o", "output", "The name of the folder to output to")
	rootCmd.PersistentFlags().StringP("pkgname", "p", "model", "The name you wish to assign to your generated package")

	viper.BindPFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Failed to execute sqlboiler command:", err)
		os.Exit(-1)
	}
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		_ = cmd.Help()
		fmt.Println("\nmust provide a driver")
		os.Exit(1)
	}

	cmdConfig = new(Config)

	cmdConfig.DriverName = args[0]
	cmdConfig.TableName = viper.GetString("table")
	cmdConfig.OutFolder = viper.GetString("folder")
	cmdConfig.PkgName = viper.GetString("pkgname")

	if len(cmdConfig.DriverName) == 0 {
		return errors.New("Must supply a driver flag.")
	}
	if len(cmdConfig.OutFolder) == 0 {
		return fmt.Errorf("No output folder specified.")
	}

	if viper.IsSet("postgres.dbname") {
		cmdConfig.Postgres = PostgresConfig{
			User:   viper.GetString("postgres.user"),
			Pass:   viper.GetString("postgres.pass"),
			Host:   viper.GetString("postgres.host"),
			Port:   viper.GetInt("postgres.port"),
			DBName: viper.GetString("postgres.dbname"),
		}
	}

	var err error
	cmdState, err = New(cmdConfig)
	return err
}

func run(cmd *cobra.Command, args []string) error {
	return cmdState.Run(true)
}

func postRun(cmd *cobra.Command, args []string) error {
	return cmdState.Cleanup()
}
