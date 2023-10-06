package machine_networking

import (
	"github.com/charmbracelet/log"
	"github.com/vishvananda/netlink"
)

// var baseMacAddress net.HardwareAddr = []byte{0xAA, 0xFC, 0x00, 0x00, 0x00, 0x01}

/*
sudo ip tuntap add tap0 mode tap
sudo ip addr add 172.16.0.1/24 dev tap0
sudo ip link set tap0 up
sudo iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
sudo iptables -A FORWARD -i tap0 -o eth0 -j ACCEPT
*/

func SetupBaseMachineNetworking(machineId string) error {
	// Create tap interface
	tap := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: "tap" + machineId,
		},
		Mode: netlink.TUNTAP_MODE_TAP,
	}

	if err := netlink.LinkAdd(tap); err != nil {
		log.Error("Failed to create tap interface", "error", err)
		return err
	}

	addr, _ := netlink.ParseAddr("172.16.0.1/24")

	if err := netlink.AddrAdd(tap, addr); err != nil {
		log.Error("Failed to assign IP address to tap interface", "error", err)
		return err
	}

	if err := netlink.LinkSetUp(tap); err != nil {
		log.Error("Failed to set tap interface up", "error", err)
		return err
	}

	// Setup NAT & iptables

	return nil
}
