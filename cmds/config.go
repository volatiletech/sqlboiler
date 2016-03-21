package cmds

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type PostgresCfg struct {
	User   string `toml:"user"`
	Pass   string `toml:"pass"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	DBName string `toml:"dbname"`
}

type Config struct {
	Postgres PostgresCfg `toml:"postgres"`
	// The TestPostgres object holds the configuration pointing to a test database.
	// This test database is used to test all of the commands that have an accompanying
	// command.testtpl file. These template files generate the go test functions.
	//
	// Note: These test templates will only be generated for the boil command,
	// if an OutFolder is provided, and if all test config variables are present.
	TestPostgres *PostgresCfg `toml:"postgres_test"`
}

var cfg *Config

// LoadConfigFile loads the toml config file into the cfg object
func LoadConfigFile(filename string) {
	_, err := toml.DecodeFile(filename, &cfg)

	if os.IsNotExist(err) {
		fmt.Printf("Failed to find the toml configuration file %s: %s", filename, err)
		return
	}

	if err != nil {
		fmt.Println("Failed to decode toml configuration file:", err)
	}

	// If any of the test cfg variables are not present then test TestPostgres to nil
	//
	// As a safety precaution, set TestPostgres to nil if
	// the dbname is the same as the cfg dbname. This will prevent the test
	// from erasing the production database tables if someone accidently
	// configures the config.toml incorrectly.
	if cfg.TestPostgres != nil {
		if cfg.TestPostgres.User == "" || cfg.TestPostgres.Pass == "" ||
			cfg.TestPostgres.Host == "" || cfg.TestPostgres.Port == 0 ||
			cfg.TestPostgres.DBName == "" || cfg.Postgres.DBName == cfg.TestPostgres.DBName {
			cfg.TestPostgres = nil
		}
	}
}
