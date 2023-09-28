package machines

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) GetMachine(ctx context.Context, machineId string) (*types.RavelMachine, error) {
	machine, err := machineManager.store.GetRavelMachine(machineId)

	if err != nil {
		log.Error("Error getting firecracker machine", "error", err)
		return nil, err
	}

	return &machine, nil
}
