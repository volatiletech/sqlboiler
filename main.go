// Package main defines a command line interface for the sqlboiler package
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kat-co/vala"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const sqlBoilerVersion = "2.0.3"

var (
	cmdState  *State
	cmdConfig *Config
)

func main() {
	var err error

	// Too much happens between here and cobra's argument handling, for
	// something so simple just do it immediately.
	for _, arg := range os.Args {
		if arg == "--version" {
			fmt.Println("SQLBoiler v" + sqlBoilerVersion)
			return
		}
	}

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

	// Ignore errors here, fallback to other validation methods.
	// Users can use environment variables if a config is not found.
	_ = viper.ReadInConfig()

	// Set up the cobra root command
	var rootCmd = &cobra.Command{
		Use:   "sqlboiler [flags] <driver>",
		Short: "SQL Boiler generates an ORM tailored to your database schema.",
		Long: "SQL Boiler generates a Go ORM from template files, tailored to your database schema.\n" +
			`Complete documentation is available at http://github.com/vattle/sqlboiler`,
		Example:       `sqlboiler postgres`,
		PreRunE:       preRun,
		RunE:          run,
		PostRunE:      postRun,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Set up the cobra root command flags
	rootCmd.PersistentFlags().StringP("output", "o", "models", "The name of the folder to output to")
	rootCmd.PersistentFlags().StringP("schema", "s", "public", "The name of your database schema, for databases that support real schemas")
	rootCmd.PersistentFlags().StringP("pkgname", "p", "models", "The name you wish to assign to your generated package")
	rootCmd.PersistentFlags().StringP("basedir", "", "", "The base directory has the templates and templates_test folders")
	rootCmd.PersistentFlags().StringSliceP("blacklist", "b", nil, "Do not include these tables in your generated package")
	rootCmd.PersistentFlags().StringSliceP("whitelist", "w", nil, "Only include these tables in your generated package")
	rootCmd.PersistentFlags().StringSliceP("tag", "t", nil, "Struct tags to be included on your models in addition to json, yaml, toml")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug mode prints stack traces on error")
	rootCmd.PersistentFlags().BoolP("no-tests", "", false, "Disable generated go test files")
	rootCmd.PersistentFlags().BoolP("no-hooks", "", false, "Disable hooks feature for your models")
	rootCmd.PersistentFlags().BoolP("no-auto-timestamps", "", false, "Disable automatic timestamps for created_at/updated_at")
	rootCmd.PersistentFlags().BoolP("version", "", false, "Print the version")

	viper.SetDefault("postgres.sslmode", "require")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("mysql.sslmode", "true")
	viper.SetDefault("mysql.port", "3306")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		if e, ok := err.(commandFailure); ok {
			fmt.Printf("Error: %v\n\n", string(e))
			rootCmd.Help()
		} else if !viper.GetBool("debug") {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Error: %+v\n", err)
		}

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
		DriverName:       driverName,
		OutFolder:        viper.GetString("output"),
		Schema:           viper.GetString("schema"),
		PkgName:          viper.GetString("pkgname"),
		BaseDir:          viper.GetString("basedir"),
		Debug:            viper.GetBool("debug"),
		NoTests:          viper.GetBool("no-tests"),
		NoHooks:          viper.GetBool("no-hooks"),
		NoAutoTimestamps: viper.GetBool("no-auto-timestamps"),
	}

	// BUG: https://github.com/spf13/viper/issues/200
	// Look up the value of blacklist, whitelist & tags directly from PFlags in Cobra if we
	// detect a malformed value coming out of viper.
	// Once the bug is fixed we'll be able to move this into the init above
	cmdConfig.BlacklistTables = viper.GetStringSlice("blacklist")
	if len(cmdConfig.BlacklistTables) == 1 && strings.HasPrefix(cmdConfig.BlacklistTables[0], "[") {
		cmdConfig.BlacklistTables, err = cmd.PersistentFlags().GetStringSlice("blacklist")
		if err != nil {
			return err
		}
	}

	cmdConfig.WhitelistTables = viper.GetStringSlice("whitelist")
	if len(cmdConfig.WhitelistTables) == 1 && strings.HasPrefix(cmdConfig.WhitelistTables[0], "[") {
		cmdConfig.WhitelistTables, err = cmd.PersistentFlags().GetStringSlice("whitelist")
		if err != nil {
			return err
		}
	}

	cmdConfig.Tags = viper.GetStringSlice("tag")
	if len(cmdConfig.Tags) == 1 && strings.HasPrefix(cmdConfig.Tags[0], "[") {
		cmdConfig.Tags, err = cmd.PersistentFlags().GetStringSlice("tag")
		if err != nil {
			return err
		}
	}

	if driverName == "postgres" {
		cmdConfig.Postgres = PostgresConfig{
			User:    viper.GetString("postgres.user"),
			Pass:    viper.GetString("postgres.pass"),
			Host:    viper.GetString("postgres.host"),
			Port:    viper.GetInt("postgres.port"),
			DBName:  viper.GetString("postgres.dbname"),
			SSLMode: viper.GetString("postgres.sslmode"),
		}

		// BUG: https://github.com/spf13/viper/issues/71
		// Despite setting defaults, nested values don't get defaults
		// Set them manually
		if cmdConfig.Postgres.SSLMode == "" {
			cmdConfig.Postgres.SSLMode = "require"
			viper.Set("postgres.sslmode", cmdConfig.Postgres.SSLMode)
		}

		if cmdConfig.Postgres.Port == 0 {
			cmdConfig.Postgres.Port = 5432
			viper.Set("postgres.port", cmdConfig.Postgres.Port)
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
	}

	if driverName == "mysql" {
		cmdConfig.MySQL = MySQLConfig{
			User:    viper.GetString("mysql.user"),
			Pass:    viper.GetString("mysql.pass"),
			Host:    viper.GetString("mysql.host"),
			Port:    viper.GetInt("mysql.port"),
			DBName:  viper.GetString("mysql.dbname"),
			SSLMode: viper.GetString("mysql.sslmode"),
		}

		// MySQL doesn't have schemas, just databases
		cmdConfig.Schema = cmdConfig.MySQL.DBName

		// BUG: https://github.com/spf13/viper/issues/71
		// Despite setting defaults, nested values don't get defaults
		// Set them manually
		if cmdConfig.MySQL.SSLMode == "" {
			cmdConfig.MySQL.SSLMode = "true"
			viper.Set("mysql.sslmode", cmdConfig.MySQL.SSLMode)
		}

		if cmdConfig.MySQL.Port == 0 {
			cmdConfig.MySQL.Port = 3306
			viper.Set("mysql.port", cmdConfig.MySQL.Port)
		}

		err = vala.BeginValidation().Validate(
			vala.StringNotEmpty(cmdConfig.MySQL.User, "mysql.user"),
			vala.StringNotEmpty(cmdConfig.MySQL.Host, "mysql.host"),
			vala.Not(vala.Equals(cmdConfig.MySQL.Port, 0, "mysql.port")),
			vala.StringNotEmpty(cmdConfig.MySQL.DBName, "mysql.dbname"),
			vala.StringNotEmpty(cmdConfig.MySQL.SSLMode, "mysql.sslmode"),
		).Check()

		if err != nil {
			return commandFailure(err.Error())
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
