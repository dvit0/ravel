package machines

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) StartMachine(ctx context.Context, machineId string) error {
	err := machineManager.machines.StartMachine(ctx, machineId)
	if err != nil {
		log.Error("Error starting firecracker machine", "error", err)
		return err
	}

	err = machineManager.store.UpdateRavelMachine(machineId, func(rm *types.RavelMachine) {
		rm.Status = types.RavelMachineStatusRunning
	})
	if err != nil {
		log.Error("Error updating machine status in store", "error", err)
		return err
	}

	return nil
}
