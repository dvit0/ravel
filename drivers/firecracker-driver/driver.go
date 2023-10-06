package main

import (
	"context"
	"errors"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/firecracker-go-sdk"
	"github.com/valyentdev/ravel/pkg/driver"
)

type FirecrackerRavelDriver struct {
	machines map[string]*firecracker.Machine
}

func NewFirecrackerDriver() *FirecrackerRavelDriver {
	return &FirecrackerRavelDriver{
		machines: map[string]*firecracker.Machine{},
	}
}

func (ravelDriver FirecrackerRavelDriver) StopVM(vmId string) error {
	return ravelDriver.stopVM(vmId)
}

func (ravelDriver FirecrackerRavelDriver) StartVM(vmId string, vmConfig driver.VMConfig) error {
	return ravelDriver.startVM(vmId, vmConfigToFirecrackerConfig(vmId, vmConfig))
}

func (ravelDriver FirecrackerRavelDriver) startVM(vmId string, config firecracker.Config) error {
	machine, err := firecracker.NewMachine(context.Background(), config)
	if err != nil {
		log.Error("Error creating machine", "error", err)
		return err
	}

	err = machine.Start(context.Background())

	if err != nil {
		return err
	}

	ravelDriver.machines[vmId] = machine
	return nil
}

func (machineManager *FirecrackerRavelDriver) cleanup() {
	log.Info("Cleaning up machines")
	for _, machine := range machineManager.machines {
		machine.StopVMM()
	}
}

func (machineManager *FirecrackerRavelDriver) stopVM(vmId string) error {
	machine, ok := machineManager.machines[vmId]
	if !ok {
		return errors.New("machine not found")
	}
	err := machine.Shutdown(context.Background())
	if err != nil {
		return err
	}

	err = machine.StopVMM()
	if err != nil {
		return err
	}

	return nil
}
