package cmds

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	SQLBoiler.AddCommand(updateCmd)
	updateCmd.Run = updateRun
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Generate update statement helpers from table definitions",
}

func updateRun(cmd *cobra.Command, args []string) {
	out := generateUpdates()

	for _, v := range out {
		os.Stdout.Write(v)
	}
}

func generateUpdates() [][]byte {
	return [][]byte{}
}
