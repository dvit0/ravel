package machine

import (
	"github.com/spf13/cobra"
	workerclient "github.com/valyentdev/ravel/cmd/ravel/client"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all  machines",
	Long:  `List all machines`,
	Run: func(cmd *cobra.Command, args []string) {
		machines, err := workerclient.Client.ListMachines()
		if err != nil {
			cmd.PrintErrln("Unable to list machines: ", err)
		}
		const format = "%24s %8s\t%s\n"

		cmd.Printf(format, "MACHINE ID", "STATUS", "IMAGE")
		for _, machine := range machines {
			cmd.Printf(format, machine.Id, machine.Status, machine.Spec.Image)
		}

	},
}

func init() {
	MachineCmd.AddCommand(lsCmd)
}
