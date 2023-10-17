package machine

import (
	"github.com/spf13/cobra"
)

// machineCmd represents the machine command
var MachineCmd = &cobra.Command{
	Use:   "machine",
	Short: "Manage ravel machines",
	Long:  `Commands to manage ravel machines.`,
}
