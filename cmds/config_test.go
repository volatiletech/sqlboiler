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

	if cfg != nil {
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

	if cfg.TestPostgres != nil || cfg.Postgres.Host != "localhost" ||
		cfg.Postgres.User != "user" || cfg.Postgres.Pass != "pass" ||
		cfg.Postgres.DBName != "mydb" || cfg.Postgres.Port != 5432 {
		t.Errorf("Config failed to load properly, got: %#v", cfg.Postgres)
	}

	fContents = `
	[postgres_test]
	host="localhost"
	port=5432
	user="testuser"
	pass="testpass"`

	file.WriteString(fContents)
	LoadConfigFile(file.Name())

	if cfg.TestPostgres != nil {
		t.Errorf("Test config failed to load properly, got: %#v", cfg.Postgres)
	}

	fContents = `
	dbname="testmydb"`

	file.WriteString(fContents)
	LoadConfigFile(file.Name())

	if cfg.TestPostgres.DBName != "testmydb" || cfg.TestPostgres.Host != "localhost" ||
		cfg.TestPostgres.User != "testuser" || cfg.TestPostgres.Pass != "testpass" ||
		cfg.TestPostgres.Port != 5432 {
		t.Errorf("Test config failed to load properly, got: %#v", cfg.Postgres)
	}
}
