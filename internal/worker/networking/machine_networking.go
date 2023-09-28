package machine_networking

import (
	"fmt"
	"net"
	"os/exec"
)

var baseMacAddress net.HardwareAddr = []byte{0x02, 0xFC, 0x00, 0x00, 0x00, 0x00}

func SetupNetworking(machineId string) error {

	tapDevice := "tap-" + machineId
	bridge := "br-" + machineId

	exec.Command("ip", "link", "del", tapDevice).Run()

	// Create the bridge
	if err := CreateBridge(bridge); err != nil {
		return fmt.Errorf("failed creating ip link: %s", err)
	}

	// Create the tap device
	if err := exec.Command("ip", "tuntap", "add", "dev", tapDevice, "mode", "tap").Run(); err != nil {
		return fmt.Errorf("failed creating ip link: %s", err)
	}
	// Add the tap device to the bridge
	if err := exec.Command("ip", "link", "set", tapDevice, "master", "firecracker0").Run(); err != nil {
		return fmt.Errorf("failed adding tap device to bridge: %s", err)
	}
	if err := exec.Command("ip", "link", "set", tapDevice, "up").Run(); err != nil {
		return fmt.Errorf("failed creating ip link: %s", err)
	}
	if err := exec.Command("sysctl", "-w", fmt.Sprintf("net.ipv4.conf.%s.proxy_arp=1", tapDevice)).Run(); err != nil {
		return fmt.Errorf("failed doing first sysctl: %s", err)
	}
	if err := exec.Command("sysctl", "-w", fmt.Sprintf("net.ipv6.conf.%s.disable_ipv6=1", tapDevice)).Run(); err != nil {
		return fmt.Errorf("failed doing second sysctl: %s", err)
	}
	return nil
}

func CleanupNetworking(machineId string) error {

	tapDevice := "tap-" + machineId

	exec.Command("ip", "link", "del", tapDevice).Run()

	return nil
}

func CreateBridge(bridgeName string) error {
	// Create the bridge
	if err := exec.Command("ip", "link", "add", bridgeName, "type", "bridge").Run(); err != nil {
		return fmt.Errorf("failed creating ip link: %s", err)
	}
	return nil
}

func DeleteBridge(bridgeName string) error {
	// Create the bridge
	if err := exec.Command("ip", "link", "del", bridgeName).Run(); err != nil {
		return fmt.Errorf("failed creating ip link: %s", err)
	}
	return nil
}

func CreateTapDevice(tapDevice string) error {
	// Create the tap device
	if err := exec.Command("ip", "tuntap", "add", "dev", tapDevice, "mode", "tap").Run(); err != nil {
		return fmt.Errorf("failed creating ip link: %s", err)
	}
	return nil
}

func DeleteTapDevice(tapDevice string) error {
	// Create the tap device
	if err := exec.Command("ip", "tuntap", "del", "dev", tapDevice, "mode", "tap").Run(); err != nil {
		return fmt.Errorf("failed creating ip link: %s", err)
	}
	return nil
}
