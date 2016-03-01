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
	// Exclude these commands from the output
	skipTemplates := []string{
		"boil",
	}

	var templateNames []string

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
			templateNames = append(templateNames, c.Name())
		}
	}

	// Sort all names alphabetically
	sort.Strings(templateNames)

	// Prepend "struct" command to templateNames slice so it sits at top of sort
	templateNames = append([]string{"struct"}, templateNames...)

	for i := 0; i < len(cmdData.TablesInfo); i++ {
		data := tplData{
			TableName: cmdData.TableNames[i],
			TableData: cmdData.TablesInfo[i],
		}

		var out [][]byte
		// Loop through and generate every command template (excluding skipTemplates)
		for _, n := range templateNames {
			out = append(out, generateTemplate(n, &data))
		}

		err := outHandler(out, &data)
		if err != nil {
			errorQuit(err)
		}
	}
}
