package machine

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	workerclient "github.com/valyentdev/ravel/cmd/ravel/client"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running machine",
	Long:  `Stop a running machine from a given config file`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Please specify a machine ID")
			return
		}

		machineId := args[0]

		err := workerclient.Client.StopMachine(context.Background(), machineId)

		if err != nil {
			cmd.Println("Unable to stop machine: ", err)
			os.Exit(1)
		}

		cmd.Println("Machine stopped")

	},
}

func init() {
	MachineCmd.AddCommand(stopCmd)
}
