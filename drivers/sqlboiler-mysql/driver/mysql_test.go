// These tests assume there is a user sqlboiler_test_user and a database
// by the name of sqlboiler_test that it has full R/W rights to.
// In order to create this you can use the following steps from a root
// mysql account:
//
//   create user sqlboiler_driver_user identified by 'sqlboiler';
//   create database sqlboiler_driver_test;
//   grant all privileges on sqlboiler_driver_test.* to sqlboiler_driver_user;

package driver

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/volatiletech/sqlboiler/v4/drivers"
)

var (
	flagOverwriteGolden = flag.Bool("overwrite-golden", false, "Overwrite the golden file with the current execution results")

	envHostname = drivers.DefaultEnv("DRIVER_HOSTNAME", "localhost")
	envPort     = drivers.DefaultEnv("DRIVER_PORT", "3306")
	envUsername = drivers.DefaultEnv("DRIVER_USER", "sqlboiler_driver_user")
	envPassword = drivers.DefaultEnv("DRIVER_PASS", "sqlboiler")
	envDatabase = drivers.DefaultEnv("DRIVER_DB", "sqlboiler_driver_test")
)

func TestDriver(t *testing.T) {
	b, err := ioutil.ReadFile("testdatabase.sql")
	if err != nil {
		t.Fatal(err)
	}

	out := &bytes.Buffer{}
	createDB := exec.Command("mysql", "-h", envHostname, "-P", envPort, "-u", envUsername, fmt.Sprintf("-p%s", envPassword), envDatabase)
	t.Logf(createDB.String())
	createDB.Stdout = out
	createDB.Stderr = out
	createDB.Stdin = bytes.NewReader(b)

	if err := createDB.Run(); err != nil {
		t.Logf("mysql output:\n%s\n", out.Bytes())
		t.Fatal(err)
	}
	t.Logf("mysql output:\n%s\n", out.Bytes())

	config := drivers.Config{
		"user":    envUsername,
		"pass":    envPassword,
		"dbname":  envDatabase,
		"host":    envHostname,
		"port":    3306,
		"sslmode": "false",
		"schema":  envDatabase,
	}

	p := &MySQLDriver{}
	info, err := p.Assemble(config)
	if err != nil {
		t.Fatal(err)
	}

	got, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	if *flagOverwriteGolden {
		if err = ioutil.WriteFile("mysql.golden.json", got, 0664); err != nil {
			t.Fatal(err)
		}
		t.Log("wrote:", string(got))
		return
	}

	want, err := ioutil.ReadFile("mysql.golden.json")
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("want:\n%s\ngot:\n%s\n", want, got)
	}
}
