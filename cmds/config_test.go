package cmds

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cmdData := &CmdData{}

	file, _ := ioutil.TempFile(os.TempDir(), "sqlboilercfgtest")
	defer os.Remove(file.Name())

	fContents := `[postgres]
    host="localhost"
    port=5432
    user="user"
    pass="pass"
    dbname="mydb"`

	file.WriteString(fContents)
	err := cmdData.LoadConfigFile(file.Name())
	if err != nil {
		t.Errorf("Unable to load config file: %s", err)
	}

	if cmdData.Config.Postgres.Host != "localhost" ||
		cmdData.Config.Postgres.User != "user" || cmdData.Config.Postgres.Pass != "pass" ||
		cmdData.Config.Postgres.DBName != "mydb" || cmdData.Config.Postgres.Port != 5432 {
		t.Errorf("Config failed to load properly, got: %#v", cmdData.Config.Postgres)
	}
}
