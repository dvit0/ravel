package machines

import (
	"context"

	"github.com/charmbracelet/log"
)

func (machineManager *MachineManager) DeleteMachine(ctx context.Context, machineId string) error {
	err := machineManager.machines.DeleteMachine(ctx, machineId)

	if err != nil {
		log.Error("Error deleting firecracker machine", "error", err)
		return err
	}

	err = machineManager.store.DeleteRavelMachine(machineId)

	if err != nil {
		log.Error("Error deleting machine from store", "error", err)
		return err
	}

	return nil
}
