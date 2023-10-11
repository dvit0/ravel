package machines

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) DeleteMachine(machineId string) error {
	machine, found, err := machineManager.store.GetRavelMachine(machineId)
	if err != nil {
		log.Error("Error getting firecracker machine", "error", err)
		return err
	}
	if !found {
		log.Error("Machine not found", "machineId", machineId)
		return errors.New("machine not found")
	}

	if machine.Status == types.RavelMachineStatusRunning {
		log.Error("Machine is not stopped", "machine", machine)
		return errors.New("machine is not stopped")
	}

	err = machineManager.store.DeleteRavelMachine(machineId)
	if err != nil {
		log.Error("Error deleting machine from store", "error", err)
		return err
	}

	return nil
}
