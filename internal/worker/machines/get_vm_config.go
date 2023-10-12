package machines

import (
	"github.com/valyentdev/ravel/internal/worker/config"
	"github.com/valyentdev/ravel/pkg/driver"
	"github.com/valyentdev/ravel/pkg/types"
)

func getVMConfig(ravelMachine *types.RavelMachine) driver.VMConfig {
	return driver.VMConfig{
		Kernel:     "./kernel/vmlinux-5.10",
		KernelArgs: "ro console=ttyS0,115200n8 noapic reboot=k panicOD=1  pci=off nomodules init=/ravel/init quiet",
		VcpuCount:  ravelMachine.Spec.Vcpus,
		Memory:     ravelMachine.Spec.Memory,

		Drives: []driver.Drive{
			{
				DriveId:         "rootfs",
				DrivePathOnHost: config.GetDriveImagePath(ravelMachine.InitDriveId),
				IsRoot:          true,
				IsReadOnly:      false,
			},
			{
				DriveId:         "mainfs",
				DrivePathOnHost: config.GetDriveImagePath(ravelMachine.RootDriveId),
				IsRoot:          false,
				IsReadOnly:      false,
			},
		},
	}
}
