package cmds

import (
	"testing"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

func init() {
	cmdData = &CmdData{
		Tables: []string{"patrick_table"},
		Columns: [][]dbdrivers.DBColumn{
			[]dbdrivers.DBColumn{
				{Name: "patrick_column", IsNullable: false},
			},
		},
		PkgName:   "patrick",
		OutFolder: "",
		DBDriver:  nil,
	}
}

// ioutil.TempDir
// os.TempDir
// set the temp dir to outfolder
// generate all the stuffs

// create a file in the tempdir folder named templates_test.go
// use exec package to run go test in that folder (exec go test in that temp folder)

// when i use the exec theres a special thing. if i look here https://golang.org/pkg/os/exec/#Cmd
// stderr (create bytes.buf, shove it into that) (use Command for initialization of obj)
// use Run (not start) on the command. run the thing which will give an error
// check that error, if its nil it completed successfully and test should pass
// if not nil, compile failed. check stderr and pump it out and fail test.
//
// use Dir to set working dir of test.
// ALWAYs REMBerR To DEFerR DleELtEE The FOoFldER
// miGtihtr WaNnaAu leAVae around  iwhen testing

func TestTemplates(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
}

/*
var testHeader = `package main

import (
)
`

func TestInitTemplates(t *testing.T) {
	templates, err := initTemplates("./templates")
	if err != nil {
		t.Errorf("Unable to init templates: %s", err)
	}

	testData := tplData{
		Table: "hello_world",
		Columns: []dbdrivers.DBColumn{
			{Name: "hello_there", Type: "int64", IsNullable: true},
			{Name: "enemy_friend_list", Type: "string", IsNullable: false},
		},
	}

	for _, tpl := range templates {
		file, err := ioutil.TempFile(os.TempDir(), "boilertemplatetest")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprintln(testHeader)

		if err = tpl.Execute(tpl, testData); err != nil {
			t.Error(err)
		}

		if err = file.Close(); err != nil {
			t.Error(err)
		}
	}
}

*/
