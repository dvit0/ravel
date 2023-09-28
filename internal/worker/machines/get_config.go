package machines

import (
	"github.com/valyentdev/firecracker-go-sdk"
	"github.com/valyentdev/firecracker-go-sdk/client/models"
	"github.com/valyentdev/ravel/internal/worker/config"
	"github.com/valyentdev/ravel/pkg/types"
)

func GetFirecrackerConfig(ravelMachine *types.RavelMachine) firecracker.Config {
	return firecracker.Config{
		SocketPath:      config.GetMachineSocketPath(ravelMachine.Id),
		KernelImagePath: "./kernel/vmlinux-5.10",
		KernelArgs:      "ro console=ttyS0,115200n8 noapic reboot=k panicOD=1  pci=off nomodules init=/ravel/init",
		VMID:            ravelMachine.Id,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(ravelMachine.VCpus),
			MemSizeMib: firecracker.Int64(ravelMachine.Memory),
		},
		Drives: []models.Drive{
			{
				DriveID:      firecracker.String("rootfs"),
				PathOnHost:   firecracker.String(config.GetDriveImagePath(ravelMachine.InitDriveId)),
				IsRootDevice: firecracker.Bool(true),
				IsReadOnly:   firecracker.Bool(false),
			},
			{
				DriveID:      firecracker.String("mainfs"),
				PathOnHost:   firecracker.String(config.GetDriveImagePath(ravelMachine.RootDriveId)),
				IsRootDevice: firecracker.Bool(false),
				IsReadOnly:   firecracker.Bool(false),
			},
		},

		LogLevel:      "Error",
		FifoLogWriter: nil,
	}
}
