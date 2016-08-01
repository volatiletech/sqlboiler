// Package main defines a command line interface for the sqlboiler package
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kat-co/vala"
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

	// Find and read config, ignore errors because we'll fall back to defaults
	// and other validation mechanisms
	_ = viper.ReadInConfig()

	// Set up the cobra root command
	var rootCmd = &cobra.Command{
		Use:   "sqlboiler [flags] <driver>",
		Short: "SQL Boiler generates boilerplate structs and statements",
		Long: "SQL Boiler generates boilerplate structs and statements from template files.\n" +
			`Complete documentation is available at http://github.com/nullbio/sqlboiler`,
		Example:       `sqlboiler -o models -p models postgres`,
		PreRunE:       preRun,
		RunE:          run,
		PostRunE:      postRun,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Set up the cobra root command flags
	rootCmd.PersistentFlags().StringSliceP("tables", "t", nil, "Tables to generate models for, all tables if empty")
	rootCmd.PersistentFlags().StringP("output", "o", "models", "The name of the folder to output to")
	rootCmd.PersistentFlags().StringP("pkgname", "p", "models", "The name you wish to assign to your generated package")

	viper.SetDefault("postgres.sslmode", "required")
	viper.BindPFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		if e, ok := err.(commandFailure); ok {
			rootCmd.HelpFunc()
			fmt.Printf("\n%s\n", string(e))
			return
		}

		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

type commandFailure string

func (c commandFailure) Error() string {
	return string(c)
}

func preRun(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		return commandFailure("must provide a driver name")
	}

	driverName := args[0]

	cmdConfig = &Config{
		DriverName: driverName,
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

	if viper.IsSet("postgres.dbname") {
		cmdConfig.Postgres = PostgresConfig{
			User:    viper.GetString("postgres.user"),
			Pass:    viper.GetString("postgres.pass"),
			Host:    viper.GetString("postgres.host"),
			Port:    viper.GetInt("postgres.port"),
			DBName:  viper.GetString("postgres.dbname"),
			SSLMode: viper.GetString("postgres.sslmode"),
		}

		// Set the default SSLMode value
		if cmdConfig.Postgres.SSLMode == "" {
			viper.Set("postgres.sslmode", "require")
			cmdConfig.Postgres.SSLMode = viper.GetString("postgres.sslmode")
		}

		err = vala.BeginValidation().Validate(
			vala.StringNotEmpty(cmdConfig.Postgres.User, "postgres.user"),
			vala.StringNotEmpty(cmdConfig.Postgres.Host, "postgres.host"),
			vala.Not(vala.Equals(cmdConfig.Postgres.Port, 0, "postgres.port")),
			vala.StringNotEmpty(cmdConfig.Postgres.DBName, "postgres.dbname"),
			vala.StringNotEmpty(cmdConfig.Postgres.SSLMode, "postgres.sslmode"),
		).Check()

		if err != nil {
			return commandFailure(err.Error())
		}
	} else if driverName == "postgres" {
		return errors.New("postgres driver requires a postgres section in the config")
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
