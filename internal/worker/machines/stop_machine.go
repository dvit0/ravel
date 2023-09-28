package machines

import "context"

func (machineManager *MachineManager) StopMachine(ctx context.Context, machineId string) error {
	return machineManager.machines.StopMachine(ctx, machineId)
}
