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

	for i := 0; i < len(cmdData.Columns); i++ {
		data := tplData{
			Table:   cmdData.Tables[i],
			Columns: cmdData.Columns[i],
			PkgName: cmdData.PkgName,
		}

		var out [][]byte
		var imps imports

		imps.standard = sqlBoilerDefaultImports.standard
		imps.thirdparty = sqlBoilerDefaultImports.thirdparty
		// Loop through and generate every command template (excluding skipTemplates)
		for _, command := range commandNames {
			imps = combineImports(imps, sqlBoilerCustomImports[command])
			out = append(out, generateTemplate(command, &data))
		}

		err := outHandler(cmdData.OutFolder, out, &data, &imps)
		if err != nil {
			errorQuit(err)
		}
	}
}

func buildCommandList() []string {
	// Exclude these commands from the output
	skipTemplates := []string{
		"boil",
	}

	var commandNames []string

	// Build a list of template names
	for _, c := range sqlBoilerCommands {
		skip := false
		for _, s := range skipTemplates {
			// Skip name if it's in the exclude list.
			// Also skip "struct" so that it can be added manually at the beginning
			// of the slice. Structs should always go to the top of the file.
			if s == c.Name() || c.Name() == "struct" {
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
