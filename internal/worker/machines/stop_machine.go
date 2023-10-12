package machines

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) StopMachine(machineId string) error {
	machine, found, err := machineManager.store.GetRavelMachine(machineId)
	if err != nil {
		log.Error("Error getting machine", "error", err)
		return err
	}
	if !found {
		log.Error("Machine not found", "machineId", machineId)
		return errors.New("machine not found")
	}

	if machine.Status != "running" {
		log.Error("Machine is not running", "machine", machine)
		return errors.New("machine is not running")
	}

	driver, err := machineManager.driversManager.GetDriver("firecracker-driver")
	if err != nil {
		log.Error("Error getting driver", "error", err)
		return errors.New("error getting driver")
	}
	machineManager.LogsManager.RemoveLogBroadcaster(machineId)

	err = driver.StopVM(machineId)
	if err != nil {
		log.Error("Error stopping VM", "error", err)
		return errors.New("error stopping VM")
	}

	err = machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
		m.InitDriveId = ""
		m.RootDriveId = ""
		m.Status = types.RavelMachineStatusStopped
	})
	if err != nil {
		log.Error("Error updating machine in the store", "error", err)
	}

	err = machineManager.drives.DeleteDrive(machine.InitDriveId)
	if err != nil {
		log.Error("Error deleting init drive", "error", err)
	}

	err = machineManager.drives.DeleteDrive(machine.RootDriveId)
	if err != nil {
		log.Error("Error deleting root drive", "error", err)
	}

	return nil
}
