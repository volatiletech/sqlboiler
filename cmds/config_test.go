package cmds

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	file, _ := ioutil.TempFile(os.TempDir(), "sqlboilercfgtest")
	defer os.Remove(file.Name())

	if cfg.Postgres.DBName != "" || testCfg.Postgres.TestDBName != "" {
		t.Errorf("Expected cfgs to be empty for the time being.")
	}

	fContents := `[postgres]
    host="localhost"
    port=5432
    user="user"
    pass="pass"
    dbname="mydb"`

	file.WriteString(fContents)
	LoadConfigFile(file.Name())

	if cfg.Postgres.Host != "localhost" || cfg.Postgres.Port != 5432 ||
		cfg.Postgres.User != "user" || cfg.Postgres.Pass != "pass" ||
		cfg.Postgres.DBName != "mydb" || testCfg.Found == true {
		t.Errorf("Config failed to load properly, got: %#v", cfg.Postgres)
	}

	fContents = `
	test_host="localhost"
	test_port=5432
	test_user="testuser"
	test_pass="testpass"`

	file.WriteString(fContents)
	LoadConfigFile(file.Name())

	if testCfg.Postgres.TestHost != "localhost" || testCfg.Postgres.TestPort != 5432 ||
		testCfg.Postgres.TestUser != "testuser" || testCfg.Postgres.TestPass != "testpass" ||
		testCfg.Postgres.TestDBName != "" || testCfg.Found == true {
		t.Errorf("Test config failed to load properly, got: %#v", testCfg.Postgres)
	}

	fContents = `
	test_dbname="testmydb"`

	file.WriteString(fContents)
	LoadConfigFile(file.Name())

	if testCfg.Postgres.TestDBName != "testmydb" || testCfg.Found != true {
		t.Errorf("Test config failed to load properly, got: %#v", testCfg.Postgres)
	}
}
