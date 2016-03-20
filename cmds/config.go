package cmds

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// cfg holds the configuration file data from config.toml
var cfg = struct {
	Postgres struct {
		User   string
		Pass   string
		Host   string
		Port   int
		DBName string
	}
}{}

// testCfg holds the test configuration file data from config.toml
var testCfg = struct {
	Found    bool
	Postgres struct {
		// Test details for template test file generation
		TestUser   string
		TestPass   string
		TestHost   string
		TestPort   int
		TestDBName string
	}
}{}

func LoadConfigFile(filename string) {
	var tmpCfg = struct {
		Postgres struct {
			User   string `toml:"user"`
			Pass   string `toml:"pass"`
			Host   string `toml:"host"`
			Port   int    `toml:"port"`
			DBName string `toml:"dbname"`
			// Test details for template test file generation
			TestUser   string `toml:"test_user"`
			TestPass   string `toml:"test_pass"`
			TestHost   string `toml:"test_host"`
			TestPort   int    `toml:"test_port"`
			TestDBName string `toml:"test_dbname"`
		} `toml:"postgres"`
	}{}

	_, err := toml.DecodeFile(filename, &tmpCfg)

	if os.IsNotExist(err) {
		fmt.Printf("Failed to find the toml configuration file %s: %s", filename, err)
		return
	}

	if err != nil {
		fmt.Println("Failed to decode toml configuration file:", err)
	}

	cfg.Postgres.User = tmpCfg.Postgres.User
	cfg.Postgres.Pass = tmpCfg.Postgres.Pass
	cfg.Postgres.Host = tmpCfg.Postgres.Host
	cfg.Postgres.Port = tmpCfg.Postgres.Port
	cfg.Postgres.DBName = tmpCfg.Postgres.DBName

	testCfg.Postgres.TestUser = tmpCfg.Postgres.TestUser
	testCfg.Postgres.TestPass = tmpCfg.Postgres.TestPass
	testCfg.Postgres.TestHost = tmpCfg.Postgres.TestHost
	testCfg.Postgres.TestPort = tmpCfg.Postgres.TestPort
	testCfg.Postgres.TestDBName = tmpCfg.Postgres.TestDBName

	// If all test cfg variables are present set found flag to true
	if testCfg.Postgres.TestUser != "" && testCfg.Postgres.TestPass != "" &&
		testCfg.Postgres.TestHost != "" && testCfg.Postgres.TestPort != 0 &&
		testCfg.Postgres.TestDBName != "" {
		testCfg.Found = true
	}

	// As a safety precaution, set found to false if
	// the dbname is the same as the cfg dbname. This will prevent the test
	// from erasing the production database tables if someone accidently
	// configures the config.toml incorrectly.
	if testCfg.Postgres.TestDBName == cfg.Postgres.DBName {
		testCfg.Found = false
		testCfg.Postgres.TestDBName = ""
	}
}
