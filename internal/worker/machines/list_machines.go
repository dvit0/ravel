package machines

import "github.com/valyentdev/ravel/pkg/types"

func (machineManager *MachineManager) ListMachines() ([]types.RavelMachine, error) {
	machines, err := machineManager.store.ListRavelMachines()
	if err != nil {
		return nil, err
	}

	return machines, nil
}
