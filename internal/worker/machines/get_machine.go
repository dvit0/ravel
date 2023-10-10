package machines

import (
	"github.com/valyentdev/ravel/pkg/types"
)

func (machineManager *MachineManager) GetMachine(machineId string) (*types.RavelMachine, bool, error) {
	return machineManager.store.GetRavelMachine(machineId)
}
