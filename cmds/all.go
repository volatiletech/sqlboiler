package cmds

import "github.com/spf13/cobra"

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Generate all templates from table definitions",
}

// allRun executes every sqlboiler command, starting with structs.
func allRun(cmd *cobra.Command, args []string) {
	// Exclude these commands from the output
	skipTemplates := []string{
		"all",
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

	// Prepend "struct" command to templateNames slice
	templateNames = append([]string{"struct"}, templateNames...)

	// Loop through and generate every command template (excluding skipTemplates)
	for _, n := range templateNames {
		err := outHandler(generateTemplate(n))
		if err != nil {
			errorQuit(err)
		}
	}
}
