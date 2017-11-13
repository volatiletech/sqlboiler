// These tests assume there is a user sqlboiler_driver_user and a database
// by the name of sqlboiler_driver_test that it has full R/W rights to.
// In order to create this you can use the following steps from a root
// psql account:
//
//   create role sqlboiler_driver_user login nocreatedb nocreaterole nocreateuser password 'sqlboiler';
//   create database sqlboiler_driver_test owner = sqlboiler_driver_user;

package driver

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/volatiletech/sqlboiler/drivers"
)

var (
	flagGolden = flag.Bool("test.golden", false, "Overwrite the golden file with the current execution results")
)

func TestDriver(t *testing.T) {
	b, err := ioutil.ReadFile("testdatabase.sql")
	if err != nil {
		t.Fatal(err)
	}

	out := &bytes.Buffer{}
	createDB := exec.Command("psql", "--user", "sqlboiler_driver_user", "sqlboiler_driver_test")
	createDB.Env = append([]string{"PGPASSWORD=sqlboiler"}, os.Environ()...)
	createDB.Stdout = out
	createDB.Stderr = out
	createDB.Stdin = bytes.NewReader(b)

	if err := createDB.Run(); err != nil {
		t.Logf("psql output:\n%s\n", out.Bytes())
		t.Fatal(err)
	}
	t.Logf("psql output:\n%s\n", out.Bytes())

	config := drivers.Config{
		"user":    "sqlboiler_driver_user",
		"pass":    "sqlboiler",
		"dbname":  "sqlboiler_driver_test",
		"host":    "localhost",
		"port":    5432,
		"sslmode": "require",
		"schema":  "public",
	}

	p := &PostgresDriver{}
	info, err := p.Assemble(config)
	if err != nil {
		t.Fatal(err)
	}

	got, err := json.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}

	if *flagGolden {
		if err = ioutil.WriteFile("psql.golden.json", got, 0664); err != nil {
			t.Fatal(err)
		}
		t.Log("wrote:", string(got))
		return
	}

	want, err := ioutil.ReadFile("psql.golden.json")
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("want:\n%s\ngot:\n%s\n", want, got)
	}
}
