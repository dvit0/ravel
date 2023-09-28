package machines

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) CreateMachine(ctx context.Context, ravelMachineSpec types.RavelMachineSpec) (string, error) {
	ravelMachine, err := machineManager.buildMachine(ravelMachineSpec)
	if err != nil {
		log.Error("Error building machine", "error", err)
		return "", err
	}

	firecrackerConfig := GetFirecrackerConfig(&ravelMachine)

	_, err = machineManager.machines.CreateMachine(ctx, ravelMachine.Id, firecrackerConfig)
	if err != nil {
		log.Error("Error creating firecracker machine", "error", err)
		return "", err
	}

	err = machineManager.store.StoreRavelMachine(&ravelMachine)

	if err != nil {
		log.Error("Error storing machine", "error", err)
		return "", err
	}

	return ravelMachine.Id, nil
}
