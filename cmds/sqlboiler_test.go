package cmds

import "github.com/pobri19/sqlboiler/dbdrivers"

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
