package main

import (
	"context"
	"errors"

	"github.com/charmbracelet/log"
	"github.com/containerd/console"
	"github.com/valyentdev/firecracker-go-sdk"
	"github.com/valyentdev/ravel/pkg/driver"
	"github.com/valyentdev/ravel/pkg/driver/proto"
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

func (ravelDriver FirecrackerRavelDriver) StartVM(vmId string, vmConfig driver.VMConfig) (*proto.StartVMResponse, error) {
	log.Debug("Starting VM", "vmId", vmId)
	return ravelDriver.startVM(vmId, vmConfigToFirecrackerConfig(vmId, vmConfig))
}

func (ravelDriver FirecrackerRavelDriver) startVM(vmId string, config firecracker.Config) (*proto.StartVMResponse, error) {
	tty, ttyPath, err := console.NewPty()
	log.Info("Creating pty", "path", ttyPath)
	if err != nil {
		log.Error("Error creating pty", "error", err)
		return nil, err
	}

	cmd := firecracker.VMCommandBuilder{}.
		WithSocketPath(config.SocketPath).
		WithStdout(tty).
		WithStderr(nil).
		Build(context.Background())

	opts := firecracker.WithProcessRunner(cmd)

	machine, err := firecracker.NewMachine(context.Background(), config, opts)
	if err != nil {
		log.Error("Error creating machine", "error", err)
		return nil, err
	}

	err = machine.Start(context.Background())

	if err != nil {
		log.Error("Error starting machine", "error", err)
		return nil, err
	}

	ravelDriver.machines[vmId] = machine

	go func() {
		err := machine.Wait(context.Background())
		tty.Close()
		if err != nil {
			log.Error("Machine exited with", "error", err, "vmId", vmId)
		}
		log.Info("Machine exited", "vmId", vmId)
	}()

	return &proto.StartVMResponse{
		Serial: ttyPath,
	}, nil
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
