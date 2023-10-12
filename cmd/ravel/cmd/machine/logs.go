package machine

import (
	"context"

	"github.com/spf13/cobra"
	workerclient "github.com/valyentdev/ravel/cmd/ravel/client"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get logs from a machine",
	Long:  `Get logs from a machine. The machine must be running.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Please specify a machine ID")
			return
		}
		workerclient.Client.GetStreamedLogs(context.Background(), args[0])

	},
}

func init() {
	MachineCmd.AddCommand(logsCmd)
}
