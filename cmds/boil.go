package cmds

import (
	"sort"

	"github.com/spf13/cobra"
)

var boilCmd = &cobra.Command{
	Use:   "boil",
	Short: "Generates ALL templates by running every command alphabetically",
}

// boilRun executes every sqlboiler command, starting with structs.
func boilRun(cmd *cobra.Command, args []string) {
	commandNames := buildCommandList()

	// Prepend "struct" command to templateNames slice so it sits at top of sort
	commandNames = append([]string{"struct"}, commandNames...)

	// Create a testCommandNames with "driverName_main" prepended to the front for the test templates
	// the main template initializes all of the testing assets
	testCommandNames := append([]string{cmdData.DriverName + "_main"}, commandNames...)

	for _, table := range cmdData.Tables {
		data := &tplData{
			Table:   table,
			PkgName: cmdData.PkgName,
		}

		var out [][]byte
		var imps imports

		imps.standard = sqlBoilerDefaultImports.standard
		imps.thirdparty = sqlBoilerDefaultImports.thirdparty

		// Loop through and generate every command template (excluding skipTemplates)
		for _, command := range commandNames {
			imps = combineImports(imps, sqlBoilerCustomImports[command])
			imps = combineConditionalTypeImports(imps, sqlBoilerConditionalTypeImports, data.Table.Columns)
			out = append(out, generateTemplate(command, data))
		}

		err := outHandler(cmdData.OutFolder, out, data, imps, false)
		if err != nil {
			errorQuit(err)
		}

		// Generate the test templates for all commands
		if len(testTemplates) != 0 {
			var testOut [][]byte
			var testImps imports

			testImps.standard = sqlBoilerDefaultTestImports.standard
			testImps.thirdparty = sqlBoilerDefaultTestImports.thirdparty

			testImps = combineImports(testImps, sqlBoilerConditionalDriverTestImports[cmdData.DriverName])

			// Loop through and generate every command test template (excluding skipTemplates)
			for _, command := range testCommandNames {
				testImps = combineImports(testImps, sqlBoilerCustomTestImports[command])
				testOut = append(testOut, generateTestTemplate(command, data))
			}

			err = outHandler(cmdData.OutFolder, testOut, data, testImps, true)
			if err != nil {
				errorQuit(err)
			}
		}
	}
}

func buildCommandList() []string {
	// Exclude these commands from the output
	skipCommands := []string{
		"boil",
		"struct",
	}

	var commandNames []string

	// Build a list of template names
	for _, c := range sqlBoilerCommands {
		skip := false
		for _, s := range skipCommands {
			// Skip name if it's in the exclude list.
			if s == c.Name() {
				skip = true
				break
			}
		}

		if !skip {
			commandNames = append(commandNames, c.Name())
		}
	}

	// Sort all names alphabetically
	sort.Strings(commandNames)
	return commandNames
}
