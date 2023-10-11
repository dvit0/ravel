package machine

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	workerclient "github.com/valyentdev/ravel/cmd/ravel/client"
	"github.com/valyentdev/ravel/pkg/types"
	"gopkg.in/yaml.v3"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new machine",
	Long: `Create & start a new machine from a given image
The machine spec is defined in a json or yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := cmd.Flag("config").Value.String()
		file, err := os.ReadFile(configFile)
		if err != nil {
			cmd.Println("Unable to read config file ", configFile, "\nerror: ", err)
			return
		}

		var machineSpec types.RavelMachineSpec

		err = yaml.Unmarshal(file, &machineSpec)
		if err != nil {
			cmd.Println("Unable to unmarshal config file ", configFile, "\nerror: ", err)
			return
		}

		log.Debug("Creating machine ", "g", machineSpec)

		machineId, err := workerclient.Client.CreateMachine(context.Background(), machineSpec)
		if err != nil {
			cmd.Println("Unable to create machine: ", err)
			return
		}

		cmd.Println(machineId)

	},
}

func init() {
	createCmd.Flags().StringP("config", "c", "", "Config file which contains machine spec")
	createCmd.MarkFlagRequired("config")

	// If config file is not specified, we need to specify all the flags

	MachineCmd.AddCommand(createCmd)
}
