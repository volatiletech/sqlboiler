// Package main defines a command line interface for the sqlboiler package
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
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
		Use:   "sqlboiler [flags] <driver>",
		Short: "SQL Boiler generates boilerplate structs and statements",
		Long: "SQL Boiler generates boilerplate structs and statements from template files.\n" +
			`Complete documentation is available at http://github.com/nullbio/sqlboiler`,
		Example:  `sqlboiler -o models -p models postgres`,
		PreRunE:  preRun,
		RunE:     run,
		PostRunE: postRun,
	}

	// Set up the cobra root command flags
	rootCmd.PersistentFlags().StringSliceP("tables", "t", nil, "Tables to generate models for, all tables if empty")
	rootCmd.PersistentFlags().StringP("output", "o", "models", "The name of the folder to output to")
	rootCmd.PersistentFlags().StringP("pkgname", "p", "models", "The name you wish to assign to your generated package")

	viper.BindPFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("\n%+v\n", err)
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		return errors.New("must provide a driver name")
	}

	cmdConfig = &Config{
		DriverName: args[0],
		OutFolder:  viper.GetString("output"),
		PkgName:    viper.GetString("pkgname"),
	}

	// BUG: https://github.com/spf13/viper/issues/200
	// Look up the value of TableNames directly from PFlags in Cobra if we
	// detect a malformed value coming out of viper.
	// Once the bug is fixed we'll be able to move this into the init above
	cmdConfig.TableNames = viper.GetStringSlice("tables")
	if len(cmdConfig.TableNames) == 1 && strings.HasPrefix(cmdConfig.TableNames[0], "[") {
		cmdConfig.TableNames, err = cmd.PersistentFlags().GetStringSlice("tables")
		if err != nil {
			return err
		}
	}

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

	cmdState, err = New(cmdConfig)
	return err
}

func run(cmd *cobra.Command, args []string) error {
	return cmdState.Run(true)
}

func postRun(cmd *cobra.Command, args []string) error {
	return cmdState.Cleanup()
}
