package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/valyentdev/ravel/cmd/ravel/cmd/machine"
)

var rootCmd = &cobra.Command{
	Use:   "ravel",
	Short: "A cli tool for managing a ravel worker.",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(machine.MachineCmd)
}
