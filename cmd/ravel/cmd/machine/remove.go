package machine

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	workerclient "github.com/valyentdev/ravel/cmd/ravel/client"
)

var removeCmd = &cobra.Command{

	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove a machine",
	Long:    `Remove a machine. The machine must be stopped.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Please specify a machine ID")
			return
		}

		machineId := args[0]

		err := workerclient.Client.DeleteMachine(context.Background(), machineId)

		if err != nil {
			cmd.Println("Unable to remove machine: ", err)
			os.Exit(1)
		}

		cmd.Println("Machine removed")

	},
}

func init() {
	MachineCmd.AddCommand(removeCmd)
}
