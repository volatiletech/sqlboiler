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

	"github.com/volatiletech/sqlboiler/drivers"
)

var (
	flagOverwriteGolden = flag.Bool("overwrite-golden", false, "Overwrite the golden file with the current execution results")
	flagHostname        = flag.String("hostname", "", "Connect to the server on the given host")
	flagUsername        = flag.String("username", "", "Username to use when connecting to server")
	flagPassword        = flag.String("password", "", "Password to use when connecting to server")
	flagDatabase        = flag.String("database", "", "The database to use")
)

func TestDriver(t *testing.T) {
	b, err := ioutil.ReadFile("testdatabase.sql")
	if err != nil {
		t.Fatal(err)
	}

	out := &bytes.Buffer{}
	createDB := exec.Command("mysql", "-h", *flagHostname, "-u", *flagUsername, fmt.Sprintf("-p%s", *flagPassword), *flagDatabase)
	createDB.Stdout = out
	createDB.Stderr = out
	createDB.Stdin = bytes.NewReader(b)

	if err := createDB.Run(); err != nil {
		t.Logf("mysql output:\n%s\n", out.Bytes())
		t.Fatal(err)
	}
	t.Logf("mysql output:\n%s\n", out.Bytes())

	config := drivers.Config{
		"user":    *flagUsername,
		"pass":    *flagPassword,
		"dbname":  *flagDatabase,
		"host":    *flagHostname,
		"port":    3306,
		"sslmode": "false",
		"schema":  *flagDatabase,
	}

	p := &MySQLDriver{}
	info, err := p.Assemble(config)
	if err != nil {
		t.Fatal(err)
	}

	got, err := json.Marshal(info)
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
