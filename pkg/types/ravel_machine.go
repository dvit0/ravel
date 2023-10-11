package types

import api "github.com/valyentdev/ravel/pkg/api/worker"

type RavelMachineSpec = api.RavelMachineSpec
type RavelMachine = api.RavelMachine

type RavelMachineStatus = api.RavelMachineStatus

const (
	RavelMachineStatusStarting RavelMachineStatus = RavelMachineStatus(api.Starting)
	RavelMachineStatusRunning  RavelMachineStatus = RavelMachineStatus(api.Running)
	RavelMachineStatusStopped  RavelMachineStatus = RavelMachineStatus(api.Stopped)
	RavelMachineStatusError    RavelMachineStatus = RavelMachineStatus(api.Error)
)
