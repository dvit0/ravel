package machines

import (
	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) Recover() error {
	machines, err := machineManager.store.ListRavelMachines()
	if err != nil {
		return err
	}

	for _, machine := range machines {
		go func(machine types.RavelMachine) {
			if machine.Status == types.RavelMachineStatusRunning || machine.Status == types.RavelMachineStatusStarting || machine.Status == types.RavelMachineStatusStopped {
				err = machineManager.cleanupMachineDrives(machine.Id)
				if err != nil {
					log.Error("Error while recovering machine", machine.Id, "error", err)
					return
				}

				err = machineManager.StartMachine(machine.Id)
				if err != nil {
					log.Error("Error while recovering machine", machine.Id, "error", err)
				}
			}
		}(machine)
	}

	return nil
}
