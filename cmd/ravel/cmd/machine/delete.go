package machine

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	workerclient "github.com/valyentdev/ravel/cmd/ravel/client"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a machine",
	Long:  `Delete a machine. The machine must be stopped.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Please specify a machine ID")
			return
		}

		machineId := args[0]

		err := workerclient.Client.DeleteMachine(context.Background(), machineId)

		if err != nil {
			cmd.Println("Unable to delete machine: ", err)
			os.Exit(1)
		}

		cmd.Println("Machine deleted")

	},
}

func init() {
	MachineCmd.AddCommand(deleteCmd)
}
