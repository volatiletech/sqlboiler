package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var testBinaryDriver = fmt.Sprintf("#!/bin/sh\ncat <<EOF%s\nEOF\n", testBinaryJSON)
var testWarningBinaryDriver = `#!/bin/sh
echo "warning binary" 1>&2
echo "{}"
`
var testBadBinaryDriver = `#!/bin/sh
echo "bad binary" 1>&2
exit 1
`

func TestBinaryDriver(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cannot run binary test on windows (needs bin/sh)")
	}

	var want, got *DBInfo
	if err := json.Unmarshal([]byte(testBinaryJSON), &want); err != nil {
		t.Fatal(err)
	}

	bin, err := os.CreateTemp("", "test_binary_driver")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprint(bin, testBinaryDriver)
	if err := bin.Chmod(0774); err != nil {
		t.Fatal(err)
	}
	if err := bin.Close(); err != nil {
		t.Fatal(err)
	}

	name := bin.Name()

	exe := binaryDriver(name)
	got, err = exe.Assemble(nil)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want:\n%#v\ngot:\n%#v\n", want, got)
	}
}

func TestBinaryWarningDriver(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cannot run binary test on windows (needs bin/sh)")
	}

	bin, err := os.CreateTemp("", "test_binary_driver")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprint(bin, testWarningBinaryDriver)
	if err := bin.Chmod(0774); err != nil {
		t.Fatal(err)
	}
	if err := bin.Close(); err != nil {
		t.Fatal(err)
	}

	stderr := &bytes.Buffer{}
	err = execute(bin.Name(), "method", nil, nil, stderr)
	if err != nil {
		t.Error(err)
	} else if !strings.Contains(stderr.String(), "warning binary") {
		t.Error("it should have written to stderr")
	}
}

func TestBinaryBadDriver(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cannot run binary test on windows (needs bin/sh)")
	}

	bin, err := os.CreateTemp("", "test_binary_driver")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprint(bin, testBadBinaryDriver)
	if err := bin.Chmod(0774); err != nil {
		t.Fatal(err)
	}
	if err := bin.Close(); err != nil {
		t.Fatal(err)
	}

	stderr := &bytes.Buffer{}
	err = execute(bin.Name(), "method", nil, nil, stderr)
	if err == nil {
		t.Error("it should have failed when the program exited 1")
	} else if !strings.Contains(stderr.String(), "bad binary") {
		t.Error("it should have written to stderr")
	} else if !strings.Contains(err.Error(), "non-zero") {
		t.Error("it should have reported non-zero exit")
	}
}

var testBinaryJSON = `
{
	"tables": [
		{
			"name": "users",
			"schema_name": "dbo",
			"columns": [
				{
					"name": "id",
					"type": "int",
					"db_type": "integer",
					"default": "",
					"nullable": false,
					"unique": true,
					"validated": false,
					"arr_type": null,
					"udt_name": "",
					"full_db_type": "",
					"auto_generated": false
				},
				{
					"name": "profile_id",
					"type": "int",
					"db_type": "integer",
					"default": "",
					"nullable": false,
					"unique": true,
					"validated": false,
					"arr_type": null,
					"udt_name": "",
					"full_db_type": "",
					"auto_generated": false
				}
			],
			"p_key": {
				"name": "pk_users",
				"columns": ["id"] 
			},
			"f_keys": [
				{
					"table": "users",
					"name": "fk_users_profile",
					"column": "profile_id",
					"nullable": false,
					"unique": true,
					"foreign_table": "profiles",
					"foreign_column": "id",
					"foreign_column_nullable": false,
					"foreign_column_unique": true
				}
			],
			"is_join_table": false,
			"to_one_relationships": [
				{
					"table": "users",
					"name": "fk_users_profile",
					"column": "profile_id",
					"nullable": false,
					"unique": true,
					"foreign_table": "profiles",
					"foreign_column": "id",
					"foreign_column_nullable": false,
					"foreign_column_unique": true
				}
			]
		}
	],
	"dialect": {
		"lq": 91,
		"rq": 93,

		"use_index_placeholders": false,
		"use_last_insert_id": true,
		"use_top_clause": false
	}
}
`
