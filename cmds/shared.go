package cmds

import (
	"fmt"
	"io"
	"os"

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

		err := outHandler(cmdData.OutFolder, out, &data)
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
func outHandler(outFolder string, output [][]byte, data *tplData) error {
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

	for _, templateOutput := range output {
		if _, err := fmt.Fprintf(out, "%s\n", templateOutput); err != nil {
			errorQuit(fmt.Errorf("Unable to write template output to file handle: %v", err))
		}
	}

	return nil
}
