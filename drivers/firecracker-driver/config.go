package main

import (
	"github.com/valyentdev/firecracker-go-sdk"
	"github.com/valyentdev/firecracker-go-sdk/client/models"
	"github.com/valyentdev/ravel/pkg/driver"
)

func vmConfigToFirecrackerConfig(vmId string, vmConfig driver.VMConfig) firecracker.Config {

	firecrackerConfig := firecracker.Config{
		SocketPath:      "/var/lib/ravel/machines/" + vmId + ".sock",
		KernelImagePath: vmConfig.Kernel,
		KernelArgs:      vmConfig.KernelArgs,
		VMID:            vmId,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(vmConfig.VcpuCount),
			MemSizeMib: firecracker.Int64(vmConfig.Memory),
		},
		Drives:            []models.Drive{},
		NetworkInterfaces: []firecracker.NetworkInterface{},
	}

	for _, drive := range vmConfig.Drives {

		firecrackerConfig.Drives = append(firecrackerConfig.Drives, models.Drive{
			DriveID:      firecracker.String(drive.DriveId),
			PathOnHost:   firecracker.String(drive.DrivePathOnHost),
			IsRootDevice: firecracker.Bool(drive.IsRoot),
			IsReadOnly:   firecracker.Bool(drive.IsReadOnly),
		})
	}

	for _, networkInterface := range vmConfig.NetworkInterfaces {

		firecrackerConfig.NetworkInterfaces = append(firecrackerConfig.NetworkInterfaces, firecracker.NetworkInterface{
			StaticConfiguration: &firecracker.StaticNetworkConfiguration{
				MacAddress:      networkInterface.MacAddress,
				HostDevName:     networkInterface.HostDevName,
				IPConfiguration: (*firecracker.IPConfiguration)(networkInterface.IPConfiguration),
			},
		})
	}

	return firecrackerConfig
}
