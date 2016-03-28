package cmds

import (
	"bytes"
	"io"
	"sort"
	"testing"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

func TestOutHandler(t *testing.T) {
	buf := &bytes.Buffer{}

	saveTestHarnessStdout := testHarnessStdout
	testHarnessStdout = buf
	defer func() {
		testHarnessStdout = saveTestHarnessStdout
	}()

	data := tplData{
		Table: dbdrivers.Table{
			Name: "patrick",
		},
	}

	templateOutputs := [][]byte{[]byte("hello world"), []byte("patrick's dreams")}

	if err := outHandler(&CmdData{PkgName: "patrick"}, templateOutputs, &data, imports{}, false); err != nil {
		t.Error(err)
	}

	if out := buf.String(); out != "package patrick\n\nhello world\npatrick's dreams\n" {
		t.Errorf("Wrong output: %q", out)
	}
}

type NopWriteCloser struct {
	io.Writer
}

func (NopWriteCloser) Close() error {
	return nil
}

func nopCloser(w io.Writer) io.WriteCloser {
	return NopWriteCloser{w}
}

func TestOutHandlerFiles(t *testing.T) {
	saveTestHarnessFileOpen := testHarnessFileOpen
	defer func() {
		testHarnessFileOpen = saveTestHarnessFileOpen
	}()

	file := &bytes.Buffer{}
	testHarnessFileOpen = func(path string) (io.WriteCloser, error) {
		return nopCloser(file), nil
	}

	data := tplData{
		Table: dbdrivers.Table{Name: "patrick"},
	}

	templateOutputs := [][]byte{[]byte("hello world"), []byte("patrick's dreams")}

	if err := outHandler(&CmdData{OutFolder: "folder", PkgName: "patrick"}, templateOutputs, &data, imports{}, false); err != nil {
		t.Error(err)
	}
	if out := file.String(); out != "package patrick\n\nhello world\npatrick's dreams\n" {
		t.Errorf("Wrong output: %q", out)
	}

	a1 := imports{
		standard: importList{
			`"fmt"`,
		},
	}
	file = &bytes.Buffer{}

	if err := outHandler(&CmdData{OutFolder: "folder", PkgName: "patrick"}, templateOutputs, &data, a1, false); err != nil {
		t.Error(err)
	}
	if out := file.String(); out != "package patrick\n\nimport \"fmt\"\nhello world\npatrick's dreams\n" {
		t.Errorf("Wrong output: %q", out)
	}

	a2 := imports{
		thirdparty: []string{
			`"github.com/spf13/cobra"`,
		},
	}
	file = &bytes.Buffer{}

	if err := outHandler(&CmdData{OutFolder: "folder", PkgName: "patrick"}, templateOutputs, &data, a2, false); err != nil {
		t.Error(err)
	}
	if out := file.String(); out != "package patrick\n\nimport \"github.com/spf13/cobra\"\nhello world\npatrick's dreams\n" {
		t.Errorf("Wrong output: %q", out)
	}

	a3 := imports{
		standard: importList{
			`"fmt"`,
			`"errors"`,
		},
		thirdparty: importList{
			`_ "github.com/lib/pq"`,
			`_ "github.com/gorilla/n"`,
			`"github.com/gorilla/mux"`,
			`"github.com/gorilla/websocket"`,
		},
	}
	file = &bytes.Buffer{}

	sort.Sort(a3.standard)
	sort.Sort(a3.thirdparty)

	if err := outHandler(&CmdData{OutFolder: "folder", PkgName: "patrick"}, templateOutputs, &data, a3, false); err != nil {
		t.Error(err)
	}

	expectedOut := `package patrick

import (
	"errors"
	"fmt"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/n"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

hello world
patrick's dreams
`

	if out := file.String(); out != expectedOut {
		t.Errorf("Wrong output (len %d, len %d): \n\n%q\n\n%q", len(out), len(expectedOut), out, expectedOut)
	}
}
