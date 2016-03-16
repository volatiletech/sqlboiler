package cmds

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

func init() {
	cmdData = &CmdData{
		Tables: []string{"patrick_table", "spiderman"},
		Columns: [][]dbdrivers.DBColumn{
			[]dbdrivers.DBColumn{
				{Name: "patrick_column", Type: "string", IsNullable: false},
				{Name: "aaron_column", Type: "null.String", IsNullable: true},
				{Name: "id", Type: "null.Int", IsNullable: true},
				{Name: "fun_id", Type: "int64", IsNullable: false},
				{Name: "time", Type: "null.Time", IsNullable: true},
				{Name: "fun_time", Type: "time.Time", IsNullable: false},
				{Name: "cool_stuff_forever", Type: "[]byte", IsNullable: false},
			},
			[]dbdrivers.DBColumn{
				{Name: "patrick", Type: "string", IsNullable: false},
			},
		},
		PkgName:   "patrick",
		OutFolder: "",
		DBDriver:  nil,
	}
}

func TestTemplates(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// Initialize the templates
	var err error
	templates, err = initTemplates("templates")
	if err != nil {
		t.Fatalf("Unable to initialize templates: %s", err)
	}

	cmdData.OutFolder, err = ioutil.TempDir("", "templates")
	if err != nil {
		t.Fatalf("Unable to create tempdir: %s", err)
	}
	defer os.RemoveAll(cmdData.OutFolder)

	boilRun(sqlBoilerCommands["boil"], []string{})

	tplFile := cmdData.OutFolder + "/templates_test.go"
	tplTestHandle, err := os.Create(tplFile)
	if err != nil {
		t.Errorf("Unable to create %s: %s", tplFile, err)
	}
	defer tplTestHandle.Close()

	fmt.Fprintf(tplTestHandle, "package %s", cmdData.PkgName)

	buf := bytes.Buffer{}
	cmd := exec.Command("go", "test", tplFile)
	cmd.Dir = cmdData.OutFolder
	cmd.Stderr = &buf

	if err = cmd.Run(); err != nil {
		t.Errorf("go test cmd execution failed: %s\n\n%s", err, buf.String())
	}
}
