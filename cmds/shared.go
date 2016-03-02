package cmds

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

// CobraRunFunc declares the cobra.Command.Run function definition
type CobraRunFunc func(cmd *cobra.Command, args []string)

// CmdData holds the table schema a slice of (column name, column type) slices.
// It also holds a slice of all of the table names sqlboiler is generating against,
// the database driver chosen by the driver flag at runtime, and a pointer to the
// output file, if one is specified with a flag.
type CmdData struct {
	Tables    []string
	Columns   [][]dbdrivers.DBColumn
	PkgName   string
	OutFolder string
	DBDriver  dbdrivers.DBDriver
}

// tplData is used to pass data to the template
type tplData struct {
	Table   string
	Columns []dbdrivers.DBColumn
}

// errorQuit displays an error message and then exits the application.
func errorQuit(err error) {
	fmt.Println(fmt.Sprintf("Error: %s\n---\n\nRun 'sqlboiler --help' for usage.", err))
	os.Exit(-1)
}

// defaultRun is the default function passed to the commands cobra.Command.Run.
// It will generate the specific commands template and send it to outHandler for output.
func defaultRun(cmd *cobra.Command, args []string) {
	// Generate the template for every table
	for i := 0; i < len(cmdData.Columns); i++ {
		data := tplData{
			Table:   cmdData.Tables[i],
			Columns: cmdData.Columns[i],
		}

		// outHandler takes a slice of byte slices, so append the Template
		// execution output to a [][]byte before sending it to outHandler.
		out := [][]byte{generateTemplate(cmd.Name(), &data)}

		imps := combineImports(sqlBoilerDefaultImports, sqlBoilerCustomImports[cmd.Name()])
		err := outHandler(cmdData.OutFolder, out, &data, &imps)
		if err != nil {
			errorQuit(fmt.Errorf("Unable to generate the template for command %s: %s", cmd.Name(), err))
		}
	}
}

var testHarnessStdout io.Writer = os.Stdout
var testHarnessFileOpen = func(filename string) (io.WriteCloser, error) {
	file, err := os.Create(filename)
	return file, err
}

// outHandler loops over the slice of byte slices, outputting them to either
// the OutFile if it is specified with a flag, or to Stdout if no flag is specified.
func outHandler(outFolder string, output [][]byte, data *tplData, imps *imports) error {
	out := testHarnessStdout

	if len(outFolder) != 0 {
		path := outFolder + "/" + data.Table + ".go"
		outFile, err := testHarnessFileOpen(path)
		if err != nil {
			errorQuit(fmt.Errorf("Unable to create output file %s: %s", path, err))
		}
		defer outFile.Close()
		out = outFile
	}

	impStr := buildImportString(imps)
	if len(impStr) > 0 {
		if _, err := fmt.Fprintf(out, "%s\n", impStr); err != nil {
			errorQuit(fmt.Errorf("Unable to write imports to file handle: %v", err))
		}
	}

	for _, templateOutput := range output {
		if _, err := fmt.Fprintf(out, "%s\n", templateOutput); err != nil {
			errorQuit(fmt.Errorf("Unable to write template output to file handle: %v", err))
		}
	}

	return nil
}

func combineImports(a, b imports) imports {
	var c imports

	c.standard = removeDuplicates(combineStringSlices(a.standard, b.standard))
	c.thirdparty = removeDuplicates(combineStringSlices(a.thirdparty, b.thirdparty))

	c.standard = sortImports(c.standard)
	c.thirdparty = sortImports(c.thirdparty)

	return c
}

// sortImports sorts the import strings alphabetically.
// If the import begins with an underscore, it temporarily
// strips it so that it does not impact the sort.
func sortImports(data []string) []string {
	sorted := make([]string, len(data))
	copy(sorted, data)

	var underscoreImports []string
	for i, v := range sorted {
		if string(v[0]) == "_" && len(v) > 1 {
			s := strings.Split(v, "_")
			underscoreImports = append(underscoreImports, s[1])
			sorted[i] = s[1]
		}
	}

	sort.Strings(sorted)

AddUnderscores:
	for i, v := range sorted {
		for _, underImp := range underscoreImports {
			if v == underImp {
				sorted[i] = "_" + sorted[i]
				continue AddUnderscores
			}
		}
	}

	return sorted
}

func combineStringSlices(a, b []string) []string {
	c := make([]string, len(a)+len(b))
	if len(a) > 0 {
		copy(c, a)
	}
	if len(b) > 0 {
		copy(c[len(a):], b)
	}

	return c
}

func removeDuplicates(dedup []string) []string {
	if len(dedup) <= 1 {
		return dedup
	}

	for i := 0; i < len(dedup)-1; i++ {
		for j := i + 1; j < len(dedup); j++ {
			if dedup[i] != dedup[j] {
				continue
			}

			if j != len(dedup)-1 {
				dedup[j] = dedup[len(dedup)-1]
				j--
			}
			dedup = dedup[:len(dedup)-1]
		}
	}

	return dedup
}

func buildImportString(imps *imports) []byte {
	stdlen, thirdlen := len(imps.standard), len(imps.thirdparty)
	if stdlen+thirdlen < 1 {
		return []byte{}
	}

	if stdlen+thirdlen == 1 {
		var imp string
		if stdlen == 1 {
			imp = imps.standard[0]
		} else {
			imp = imps.thirdparty[0]
		}
		return []byte(fmt.Sprintf(`import "%s"`, imp))
	}

	buf := &bytes.Buffer{}
	buf.WriteString("import (")
	for _, std := range imps.standard {
		fmt.Fprintf(buf, "\n\t\"%s\"", std)
	}
	if stdlen != 0 && thirdlen != 0 {
		buf.WriteString("\n")
	}
	for _, third := range imps.thirdparty {
		fmt.Fprintf(buf, "\n\t\"%s\"", third)
	}
	buf.WriteString("\n)\n")

	return buf.Bytes()
}
