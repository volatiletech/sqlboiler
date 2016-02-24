package cmds

import "github.com/spf13/cobra"

// init the "all" command
func init() {
	SQLBoiler.AddCommand(allCmd)
	allCmd.Run = allRun
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Generate all templates from table definitions",
}

// allRun executes every sqlboiler command, starting with structs
func allRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateStructs())
	if err != nil {
		errorQuit(err)
	}

	err = outHandler(generateDeletes())
	if err != nil {
		errorQuit(err)
	}

	err = outHandler(generateInserts())
	if err != nil {
		errorQuit(err)
	}

	err = outHandler(generateSelects())
	if err != nil {
		errorQuit(err)
	}
}
