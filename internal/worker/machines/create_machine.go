package machines

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/utils"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) CreateMachine(ctx context.Context, ravelMachineSpec types.RavelMachineSpec) (string, error) {
	now := time.Now()
	ravelMachine := types.RavelMachine{
		Id:               utils.NewId(),
		RavelMachineSpec: &ravelMachineSpec,
	}

	err := machineManager.store.StoreRavelMachine(&ravelMachine)
	if err != nil {
		log.Error("Error storing machine", "error", err)
		return "", err
	}

	machineManager.StartMachine(ravelMachine.Id)
	log.Info("Created machine", "machineId", ravelMachine.Id, "duration", time.Since(now))

	return ravelMachine.Id, nil
}
