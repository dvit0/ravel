package machines

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) cleanupMachineDrives(machineId string) error {
	machine, found, err := machineManager.GetMachine(machineId)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("machine not found")
	}

	err = machineManager.drives.DeleteDrive(machine.InitDriveId)
	if err != nil {
		log.Error("Error deleting init drive", "error", err)
	}

	err = machineManager.drives.DeleteDrive(machine.RootDriveId)
	if err != nil {
		log.Error("Error deleting root drive", "error", err)
		return err
	}

	err = machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
		m.InitDriveId = ""
		m.RootDriveId = ""
	})

	return err
}
