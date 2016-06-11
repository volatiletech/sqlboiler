package cmds

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/nullbio/sqlboiler/dbdrivers"
)

var cmdData *CmdData
var rgxHasSpaces = regexp.MustCompile(`^\s+`)

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
				PKey: &dbdrivers.PrimaryKey{
					Name:    "pkey_thing",
					Columns: []string{"id", "fun_id"},
				},
			},
			{
				Name: "spiderman",
				Columns: []dbdrivers.Column{
					{Name: "id", Type: "int64", IsNullable: false},
				},
				PKey: &dbdrivers.PrimaryKey{
					Name:    "pkey_id",
					Columns: []string{"id"},
				},
			},
			{
				Name: "spiderman_table_two",
				Columns: []dbdrivers.Column{
					{Name: "id", Type: "int64", IsNullable: false},
					{Name: "patrick", Type: "string", IsNullable: false},
				},
				PKey: &dbdrivers.PrimaryKey{
					Name:    "pkey_id",
					Columns: []string{"id"},
				},
			},
		},
		PkgName:    "patrick",
		OutFolder:  "",
		DriverName: "postgres",
		Interface:  nil,
	}
}

func TestLoadTemplate(t *testing.T) {
	t.Parallel()

	template, err := loadTemplate("templates_test/main_test", "postgres_main.tpl")
	if err != nil {
		t.Fatalf("Unable to loadTemplate: %s", err)
	}

	if template == nil {
		t.Fatal("Unable to load template.")
	}
}

func TestTemplates(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	if err := checkPKeys(cmdData.Tables); err != nil {
		t.Fatalf("%s", err)
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

	cmdData.SingleTemplates, err = loadTemplates("templates/singleton")
	if err != nil {
		t.Fatalf("Unable to initialize singleton templates: %s", err)
	}

	if len(cmdData.SingleTemplates) == 0 {
		t.Errorf("SingleTemplates is empty.")
	}

	cmdData.TestTemplates, err = loadTemplates("templates_test")
	if err != nil {
		t.Fatalf("Unable to initialize templates: %s", err)
	}

	if len(cmdData.Templates) == 0 {
		t.Errorf("Templates is empty.")
	}

	cmdData.TestMainTemplate, err = loadTemplate("templates_test/main_test", "postgres_main.tpl")
	if err != nil {
		t.Fatalf("Unable to initialize templates: %s", err)
	}

	cmdData.OutFolder, err = ioutil.TempDir("", "templates")
	if err != nil {
		t.Fatalf("Unable to create tempdir: %s", err)
	}
	defer os.RemoveAll(cmdData.OutFolder)

	if err = cmdData.run(true); err != nil {
		t.Errorf("Unable to run SQLBoilerRun: %s", err)
	}

	buf := &bytes.Buffer{}

	cmd := exec.Command("go", "test", "-c")
	cmd.Dir = cmdData.OutFolder
	cmd.Stderr = buf

	if err = cmd.Run(); err != nil {
		t.Errorf("go test cmd execution failed: %s", err)
		outputCompileErrors(buf, cmdData.OutFolder)
		fmt.Println()
	}
}

func outputCompileErrors(buf *bytes.Buffer, outFolder string) {
	type errObj struct {
		errMsg     string
		fileName   string
		lineNumber int
	}

	var errObjects []errObj
	lineBuf := &bytes.Buffer{}

	bufLines := bytes.Split(buf.Bytes(), []byte{'\n'})
	for i := 0; i < len(bufLines); i++ {
		lineBuf.Reset()
		if !bytes.HasPrefix(bufLines[i], []byte("./")) {
			continue
		}

		fmt.Fprintf(lineBuf, "%s\n", bufLines[i])

		splits := bytes.Split(bufLines[i], []byte{':'})
		lineNum, err := strconv.Atoi(string(splits[1]))
		if err != nil {
			panic(fmt.Sprintf("Cant convert line number to int: %s", bufLines[i]))
		}

		eObj := errObj{
			fileName:   string(splits[0]),
			lineNumber: lineNum,
		}

		for y := i; y < len(bufLines); y++ {
			if !rgxHasSpaces.Match(bufLines[y]) {
				break
			}
			fmt.Fprintf(lineBuf, "%s\n", bufLines[y])
			i++
		}

		eObj.errMsg = lineBuf.String()

		errObjects = append(errObjects, eObj)
	}

	for _, eObj := range errObjects {
		fmt.Printf("-----------------\n")
		fmt.Println(eObj.errMsg)

		filePath := filepath.Join(outFolder, eObj.fileName)
		fh, err := os.Open(filePath)
		if err != nil {
			panic(fmt.Sprintf("Cant open the file: %#v", eObj))
		}

		scanner := bufio.NewScanner(fh)
		throwaway := eObj.lineNumber - 5
		for throwaway > 0 && scanner.Scan() {
			throwaway--
		}

		for i := 0; i < 6; i++ {
			if scanner.Scan() {
				b := scanner.Bytes()
				if len(b) != 0 {
					fmt.Printf("%s\n", b)
				} else {
					i--
				}
			}
		}

		fh.Close()
	}
}
