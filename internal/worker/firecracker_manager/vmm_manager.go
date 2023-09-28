package vmm_manager

import (
	"context"
	"errors"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/firecracker-go-sdk"
	"github.com/valyentdev/ravel/internal/worker/store"
)

type VMMManager struct {
	machines map[string]*firecracker.Machine
}

func NewMachinesManager(store *store.Store) *VMMManager {
	return &VMMManager{
		machines: make(map[string]*firecracker.Machine),
	}
}

func (machineManager *VMMManager) Cleanup(ctx context.Context) {
	log.Info("Cleaning up machines")
	for _, machine := range machineManager.machines {
		machine.Shutdown(ctx)
		machine.StopVMM()
	}
}

func (machinesManager *VMMManager) CreateMachine(ctx context.Context, machineId string, config firecracker.Config) (*firecracker.Machine, error) {
	machine, err := firecracker.NewMachine(ctx, config)
	if err != nil {
		log.Error("Error creating machine", "error", err)
		return nil, err
	}

	machinesManager.machines[machineId] = machine

	return machine, nil
}

func (machinesManager *VMMManager) StartMachine(ctx context.Context, machineId string) error {
	machine, ok := machinesManager.machines[machineId]
	if !ok {
		return errors.New("machine not found")
	}

	err := machine.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (machineManager *VMMManager) StopMachine(ctx context.Context, machineId string) error {
	machine, ok := machineManager.machines[machineId]
	if !ok {
		return errors.New("machine not found")
	}

	return machine.Shutdown(ctx)
}

func (machineManager *VMMManager) GetMachine(machineId string) (*firecracker.Machine, error) {
	machine, ok := machineManager.machines[machineId]
	if !ok {
		return nil, errors.New("machine not found")
	}

	return machine, nil
}

func (machineManager *VMMManager) DeleteMachine(ctx context.Context, machineId string) error {
	machine, ok := machineManager.machines[machineId]
	if !ok {
		return errors.New("machine not found")
	}

	err := machine.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = machine.StopVMM()
	if err != nil {
		return err
	}

	delete(machineManager.machines, machineId)

	return nil
}
