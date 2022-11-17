// These tests assume there is a user sqlboiler_driver_user and a database
// by the name of sqlboiler_driver_test that it has full R/W rights to.
// In order to create this you can use the following steps from a root
// mssql account:
//
//   create database sqlboiler_driver_test;
//   go
//   use sqlboiler_driver_test;
//   go
//   create user sqlboiler_driver_user with password = 'sqlboiler';
//   go
//   exec sp_configure 'contained database authentication', 1;
//   go
//   reconfigure
//   go
//   alter database sqlboiler_driver_test set containment = partial;
//   go
//   create user sqlboiler_driver_user with password = 'Sqlboiler@1234';
//   go
//   grant alter, control to sqlboiler_driver_user;
//   go

package driver

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/volatiletech/sqlboiler/v4/drivers"
)

var (
	flagOverwriteGolden = flag.Bool("overwrite-golden", false, "Overwrite the golden file with the current execution results")

	envHostname = drivers.DefaultEnv("DRIVER_HOSTNAME", "localhost")
	envPort     = drivers.DefaultEnv("DRIVER_PORT", "1433")
	envUsername = drivers.DefaultEnv("DRIVER_USER", "sqlboiler_driver_user")
	envPassword = drivers.DefaultEnv("DRIVER_PASS", "Sqlboiler@1234")
	envDatabase = drivers.DefaultEnv("DRIVER_DB", "sqlboiler_driver_test")

	rgxKeyIDs = regexp.MustCompile(`__[A-F0-9]+$`)
)

func TestDriver(t *testing.T) {
	out := &bytes.Buffer{}
	createDB := exec.Command("sqlcmd", "-S", envHostname, "-U", envUsername, "-P", envPassword, "-d", envDatabase, "-b", "-i", "testdatabase.sql")
	createDB.Stdout = out
	createDB.Stderr = out

	if err := createDB.Run(); err != nil {
		t.Logf("mssql output:\n%s\n", out.Bytes())
		t.Fatal(err)
	}
	t.Logf("mssql output:\n%s\n", out.Bytes())

	config := drivers.Config{
		"user":    envUsername,
		"pass":    envPassword,
		"dbname":  envDatabase,
		"host":    envHostname,
		"port":    envPort,
		"sslmode": "disable",
		"schema":  "dbo",
	}

	p := &MSSQLDriver{}
	info, err := p.Assemble(config)
	if err != nil {
		t.Fatal(err)
	}

	for _, t := range info.Tables {
		if t.IsView {
			continue
		}

		t.PKey.Name = rgxKeyIDs.ReplaceAllString(t.PKey.Name, "")
		for i := range t.FKeys {
			t.FKeys[i].Name = rgxKeyIDs.ReplaceAllString(t.FKeys[i].Name, "")
		}
	}

	got, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	if *flagOverwriteGolden {
		if err = os.WriteFile("mssql.golden.json", got, 0664); err != nil {
			t.Fatal(err)
		}
		t.Log("wrote:", string(got))
		return
	}

	want, err := os.ReadFile("mssql.golden.json")
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("want:\n%s\ngot:\n%s\n", want, got)
	}
}
