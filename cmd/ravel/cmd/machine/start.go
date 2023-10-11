package machine

import (
	"context"

	"github.com/spf13/cobra"
	workerclient "github.com/valyentdev/ravel/cmd/ravel/client"
)

var startCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Start a previously stopped machine",
	Long:                  `Start a previously stopped machine. This is a no-op if the machine is already running.`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Please specify a machineId")
			cmd.Help()
			return
		}

		machineId := args[0]

		err := workerclient.Client.StartMachine(context.Background(), machineId)

		if err != nil {
			cmd.Println("Error while starting machine: ", err)
			return
		}

		cmd.Println(machineId)

	},
}

func init() {
	MachineCmd.AddCommand(startCmd)

}
