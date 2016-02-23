package cmds

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	SQLBoiler.AddCommand(deleteCmd)
	deleteCmd.Run = deleteRun
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Generate delete statement helpers from table definitions",
}

func deleteRun(cmd *cobra.Command, args []string) {
	out := generateDeletes()

	for _, v := range out {
		os.Stdout.Write(v)
	}
}

func generateDeletes() [][]byte {
	return [][]byte{}
}
