package machines

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/containerd/console"
	logger "github.com/valyentdev/ravel/internal/worker/logs"
	"github.com/valyentdev/ravel/pkg/driver/proto"
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) StartMachine(machineId string) error {
	machine, found, err := machineManager.store.GetRavelMachine(machineId)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("machine not found")
	}
	if machine.Status == types.RavelMachineStatusRunning {
		return errors.New("machine is not stopped")
	}
	err = machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
		m.Status = types.RavelMachineStatusStarting
	})
	if err != nil {
		log.Error("Error updating machine in the store", "error", err)
		return err
	}

	initDriveId, rootDriveId, err := machineManager.buildMachineDrives(*machine.RavelMachineSpec)
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

	_, err = driver.StartVM(machineId, getVMConfig(machine))
	if err != nil {
		return errors.New("error starting VM")
	}

	machineManager.store.UpdateRavelMachine(machineId, func(m *types.RavelMachine) {
		m.Status = types.RavelMachineStatusRunning
	})

	return nil
}

func startCollectingLogs(machineId string, infos *proto.StartVMResponse) {
	path := infos.Serial
	file, err := os.Open(path)
	if err != nil {
		log.Error("Error opening serial", "error", err)
		return
	}
	defer file.Close()

	cons, err := console.ConsoleFromFile(file)
	if err != nil {
		log.Error("Error getting console", "error", err)
		return
	}

	logWriter := logger.NewRotateWriter(logger.RotateWriterOptions{
		Filename:      "machine.log",
		MaxSizeByFile: 1024,
		MaxFiles:      3,
		Directory:     "/tmp/ravel/" + machineId,
	})

	for {
		log.Info("Reading from console")
		data := make([]byte, 128)
		n, err := cons.Read(data)
		if err != nil {
			log.Info("Stopping reading the logs for machine", "machineId", machineId)
			return
		}

		_, err = logWriter.Write(data[:n])
		if err != nil {
			log.Error("Error writing to log", "error", err)
			return
		}

	}
}
