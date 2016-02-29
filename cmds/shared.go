package cmds

import (
	"fmt"
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
	TablesInfo [][]dbdrivers.DBTable
	TableNames []string
	DBDriver   dbdrivers.DBDriver
	OutFile    *os.File
}

// tplData is used to pass data to the template
type tplData struct {
	TableName string
	TableData []dbdrivers.DBTable
}

// errorQuit displays an error message and then exits the application.
func errorQuit(err error) {
	fmt.Println(fmt.Sprintf("Error: %s\n---\n\nRun 'sqlboiler --help' for usage.", err))
	os.Exit(-1)
}

// defaultRun is the default function passed to the commands cobra.Command.Run.
// It will generate the specific commands template and send it to outHandler for output.
func defaultRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateTemplate(cmd.Name()))
	if err != nil {
		errorQuit(fmt.Errorf("Unable to generate the template for command %s: %s", cmd.Name(), err))
	}
}

// outHandler loops over the slice of byte slices, outputting them to either
// the OutFile if it is specified with a flag, or to Stdout if no flag is specified.
func outHandler(data [][]byte) error {
	nl := []byte{'\n'}

	// Use stdout if no outfile is specified
	var out *os.File
	if cmdData.OutFile == nil {
		out = os.Stdout
	} else {
		out = cmdData.OutFile
	}

	for _, v := range data {
		if _, err := out.Write(v); err != nil {
			return err
		}
		if _, err := out.Write(nl); err != nil {
			return err
		}
	}

	return nil
}
