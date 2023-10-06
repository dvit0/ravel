package driver

import (
	"net"

	"github.com/valyentdev/ravel/pkg/driver/proto"
)

type VMConfig struct {
	VcpuCount         int64
	Memory            int64
	Kernel            string
	KernelArgs        string
	Drives            []Drive
	NetworkInterfaces []NetworkInterface
}

type Drive struct {
	DriveId         string
	DrivePathOnHost string // Path to the drive on the host ( ext4 file system )
	IsRoot          bool
	IsReadOnly      bool
}

type NetworkInterface struct {
	MacAddress      string
	HostDevName     string
	IPConfiguration *IPConfiguration
}

type IPConfiguration struct {
	IPAddr      net.IPNet
	Gateway     net.IP
	Nameservers []string
	IfName      string
}

func (vmConfig *VMConfig) ToProto() *proto.StartVMRequest {

	createVMRequest := proto.StartVMRequest{
		VcpuCount:         vmConfig.VcpuCount,
		Memory:            vmConfig.Memory,
		Kernel:            vmConfig.Kernel,
		KernelArgs:        vmConfig.KernelArgs,
		Drives:            []*proto.Drive{},
		NetworkInterfaces: []*proto.NetworkInterface{},
	}

	for _, drive := range vmConfig.Drives {
		createVMRequest.Drives = append(createVMRequest.Drives, &proto.Drive{
			DriveId:         drive.DriveId,
			DrivePathOnHost: drive.DrivePathOnHost,
			IsRoot:          drive.IsRoot,
			IsReadOnly:      drive.IsReadOnly,
		})
	}

	for _, networkInterface := range vmConfig.NetworkInterfaces {
		createVMRequest.NetworkInterfaces = append(createVMRequest.NetworkInterfaces, &proto.NetworkInterface{
			MacAddress:  networkInterface.MacAddress,
			HostDevName: networkInterface.HostDevName,
			IpConfiguration: &proto.IPConfiguration{
				IpAddr: &proto.IPNet{
					Ip:   networkInterface.IPConfiguration.IPAddr.IP,
					Mask: networkInterface.IPConfiguration.IPAddr.Mask,
				},
				Gateway:     networkInterface.IPConfiguration.Gateway,
				Nameservers: networkInterface.IPConfiguration.Nameservers,
				IfName:      networkInterface.IPConfiguration.IfName,
			},
		})
	}

	return &createVMRequest

}

func VMConfigFromProto(vmProto *proto.StartVMRequest) VMConfig {

	vmConfig := VMConfig{
		VcpuCount:         vmProto.VcpuCount,
		Memory:            vmProto.Memory,
		Kernel:            vmProto.Kernel,
		KernelArgs:        vmProto.KernelArgs,
		Drives:            []Drive{},
		NetworkInterfaces: []NetworkInterface{},
	}

	for _, drive := range vmProto.Drives {
		vmConfig.Drives = append(vmConfig.Drives, Drive{
			DriveId:         drive.DriveId,
			DrivePathOnHost: drive.DrivePathOnHost,
			IsRoot:          drive.IsRoot,
			IsReadOnly:      drive.IsReadOnly,
		})
	}

	for _, networkInterface := range vmProto.NetworkInterfaces {
		vmConfig.NetworkInterfaces = append(vmConfig.NetworkInterfaces, NetworkInterface{
			MacAddress:  networkInterface.MacAddress,
			HostDevName: networkInterface.HostDevName,
			IPConfiguration: &IPConfiguration{
				IPAddr: net.IPNet{
					IP:   networkInterface.IpConfiguration.IpAddr.Ip,
					Mask: networkInterface.IpConfiguration.IpAddr.Mask,
				},

				Gateway:     networkInterface.IpConfiguration.Gateway,
				Nameservers: networkInterface.IpConfiguration.Nameservers,
				IfName:      networkInterface.IpConfiguration.IfName,
			},
		})
	}

	return vmConfig

}
