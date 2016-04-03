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

var cmdData *CmdData

func init() {
	cmdData = &CmdData{
		Tables: []dbdrivers.Table{
			{
				Name: "patrick_table",
				Columns: []dbdrivers.Column{
					{Name: "patrick_column", Type: "string", IsNullable: false},
					{Name: "aaron_column", Type: "null.String", IsNullable: true},
					{Name: "id", Type: "null.Int", IsNullable: true},
					{Name: "fun_id", Type: "int64", IsNullable: false},
					{Name: "time", Type: "null.Time", IsNullable: true},
					{Name: "fun_time", Type: "time.Time", IsNullable: false},
					{Name: "cool_stuff_forever", Type: "[]byte", IsNullable: false},
				},
			},
			{
				Name: "spiderman",
				Columns: []dbdrivers.Column{
					{Name: "patrick", Type: "string", IsNullable: false},
				},
			},
		},
		PkgName:   "patrick",
		OutFolder: "",
		Interface: nil,
	}
}

func TestTemplates(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// Initialize the templates
	var err error
	cmdData.Templates, err = loadTemplates("templates")
	if err != nil {
		t.Fatalf("Unable to initialize templates: %s", err)
	}

	if len(cmdData.Templates) == 0 {
		t.Errorf("Templates is empty.")
	}

	cmdData.TestTemplates, err = loadTemplates("templates_test")
	if err != nil {
		t.Fatalf("Unable to initialize templates: %s", err)
	}

	if len(cmdData.Templates) == 0 {
		t.Errorf("Templates is empty.")
	}

	cmdData.OutFolder, err = ioutil.TempDir("", "templates")
	if err != nil {
		t.Fatalf("Unable to create tempdir: %s", err)
	}
	defer os.RemoveAll(cmdData.OutFolder)

	if err = cmdData.run(true); err != nil {
		t.Errorf("Unable to run SQLBoilerRun: %s", err)
	}

	tplFile := cmdData.OutFolder + "/templates_test.go"
	tplTestHandle, err := os.Create(tplFile)
	if err != nil {
		t.Errorf("Unable to create %s: %s", tplFile, err)
	}
	defer tplTestHandle.Close()

	fmt.Fprintf(tplTestHandle, "package %s\n", cmdData.PkgName)

	buf := bytes.Buffer{}
	buf2 := bytes.Buffer{}

	cmd := exec.Command("go", "test")
	cmd.Dir = cmdData.OutFolder
	cmd.Stderr = &buf
	cmd.Stdout = &buf2

	if err = cmd.Run(); err != nil {
		t.Errorf("go test cmd execution failed: %s\n\n%s", err, buf.String())
	}
}
