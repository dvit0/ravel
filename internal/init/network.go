package init

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/vishvananda/netlink"
)

func NetworkSetup() error {
	log.Debug("setting network interfaces up")

	lo, err := netlink.LinkByName("lo")
	if err != nil {
		return fmt.Errorf("error getting loopback interface: %v", err)
	}

	if err := netlink.LinkSetUp(lo); err != nil {
		return fmt.Errorf("error configuring loopback interface: %v", err)
	}

	eth0, err := netlink.LinkByName("eth0")
	if err != nil {
		return fmt.Errorf("error getting eth0 interface: %v", err)
	}

	if err := netlink.LinkSetUp(eth0); err != nil {
		return fmt.Errorf("error configuring eth0 interface: %v", err)
	}

	return nil
}

func WriteEtcResolv(entries EtcResolv) error {
	log.Debug("populating /etc/resolv.conf")

	f, err := os.OpenFile("/etc/resolv.conf", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm0755)
	if err != nil {
		return fmt.Errorf("error opening resolv.conf: %v", err)
	}
	defer f.Close()

	for _, entry := range entries.Nameservers {
		if _, err := fmt.Fprintf(f, "\nnameserver\t%s", entry); err != nil {
			return fmt.Errorf("error writing to resolv.conf file: %v", err)
		}
	}

	if _, err := f.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

var (
	defaultHosts = []EtcHost{
		{IP: "127.0.0.1", Host: "localhost localhost.localdomain localhost4 localhost4.localdomain4"},
		{IP: "::1", Host: "localhost localhost.localdomain localhost6 localhost6.localdomain6"},
		{IP: "fe00::0", Host: "ip6-localnet"},
		{IP: "ff00::0", Host: "ip6-mcastprefix"},
		{IP: "ff02::1", Host: "ip6-allnodes"},
		{IP: "ff02::2", Host: "ip6-allrouters"},
	}

	etchostPath = "/etc/hosts"
)

func WriteEtcHost(hosts []EtcHost) error {
	log.Debug("populating /etc/hosts")

	records := append(defaultHosts, hosts...)

	f, err := os.OpenFile(etchostPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm0755)
	if err != nil {
		return fmt.Errorf("error opening /etc/hosts file: %v", err)
	}
	defer f.Close()

	for _, entry := range records {
		if entry.Desc != "" {
			_, err := fmt.Fprintf(f, "# %s\n%s\t%s\n", entry.Desc, entry.IP, entry.Host)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(f, "%s\t%s\n", entry.IP, entry.Host)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
