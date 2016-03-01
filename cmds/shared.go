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
	PkgName    string
	OutFolder  string
	DBDriver   dbdrivers.DBDriver
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
	// Generate the template for every table
	for i := 0; i < len(cmdData.TablesInfo); i++ {
		data := tplData{
			TableName: cmdData.TableNames[i],
			TableData: cmdData.TablesInfo[i],
		}

		// outHandler takes a slice of byte slices, so append the Template
		// execution output to a [][]byte before sending it to outHandler.
		out := [][]byte{
			0: generateTemplate(cmd.Name(), &data),
		}

		err := outHandler(out, &data)
		if err != nil {
			errorQuit(fmt.Errorf("Unable to generate the template for command %s: %s", cmd.Name(), err))
		}
	}
}

// outHandler loops over the slice of byte slices, outputting them to either
// the OutFile if it is specified with a flag, or to Stdout if no flag is specified.
func outHandler(output [][]byte, data *tplData) error {
	nl := []byte{'\n'}

	if cmdData.OutFolder == "" {
		for _, v := range output {
			if _, err := os.Stdout.Write(v); err != nil {
				return err
			}

			if _, err := os.Stdout.Write(nl); err != nil {
				return err
			}
		}
	} else { // If not using stdout, attempt to create the model file.
		path := cmdData.OutFolder + "/" + data.TableName + ".go"
		out, err := os.Create(path)
		if err != nil {
			errorQuit(fmt.Errorf("Unable to create output file %s: %s", path, err))
		}

		// Combine the slice of slice into a single byte slice.
		var newOutput []byte
		for _, v := range output {
			newOutput = append(newOutput, v...)
			newOutput = append(newOutput, nl...)
		}

		if _, err := out.Write(newOutput); err != nil {
			return err
		}
	}

	return nil
}
