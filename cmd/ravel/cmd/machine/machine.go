/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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
