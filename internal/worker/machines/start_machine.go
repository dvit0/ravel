package machines

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/worker/logsmanager"
	"github.com/valyentdev/ravel/pkg/driver/proto"
	"github.com/valyentdev/ravel/pkg/types"
	"github.com/valyentdev/ravel/pkg/units"
)

func (machineManager *MachineManager) StartMachine(machineId string) error {
	machine, found, err := machineManager.store.GetRavelMachine(machineId)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("machine not found")
	}

	err = machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
		m.Status = types.RavelMachineStatusStarting
	})
	if err != nil {
		log.Error("Error updating machine in the store", "error", err)
		return err
	}

	initDriveId, rootDriveId, err := machineManager.buildMachineDrives(machine.Spec)
	if err != nil {
		machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
			m.Status = types.RavelMachineStatusError
		})
		return err
	}

	err = machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
		m.InitDriveId = initDriveId
		m.RootDriveId = rootDriveId
	})
	if err != nil {
		log.Error("Error updating machine in the store", "error", err)
		return err
	}

	machine.RootDriveId = rootDriveId
	machine.InitDriveId = initDriveId

	driver, err := machineManager.driversManager.GetDriver("firecracker-driver")
	if err != nil {
		return errors.New("error getting driver")
	}

	vminfos, err := driver.StartVM(machineId, getVMConfig(machine))
	if err != nil {
		machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
			m.Status = types.RavelMachineStatusError
		})
		return errors.New("error starting VM")
	}

	machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
		m.Status = types.RavelMachineStatusRunning
	})

	machineManager.startCollectingLogs(machineId, vminfos)

	return nil
}

func (machineManager *MachineManager) startCollectingLogs(machineId string, infos *proto.StartVMResponse) {

	broadcaster := machineManager.LogsManager.NewLogBroadcaster(machineId, infos.Serial, logsmanager.RotateWriterOptions{
		Filename:      "machine.log",
		Directory:     "/var/log/ravel/machines/" + machineId,
		MaxFiles:      5,
		MaxSizeByFile: 1 * units.MB,
	})

	go broadcaster.Start()

}
