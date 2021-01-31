package model

import (
	"fmt"
	"log"
	parser "net"
	"os"
	"strconv"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

type Network struct {
	Name        string `mapstructure:"name"`
	Kvm         string `mapstructure:"kvm"`
	Mac         string `mapstructure:"mac"`
	IP          string `mapstructure:"ip"`
	Mask        string `mapstructure:"mask"`
	Gateway     string `mapstructure:"gateway"`
	DNS         string `mapstructure:"dns"`
	DhcpStartIP string `mapstructure:"dhcpstartip"`
	DhcpEndIP   string `mapstructure:"dhcpendip"`
	Prefix      string `mapstructure:"prefix"`
}

func (n *Network) initNetwork(index int) {

	n.Prefix = n.GetNetKvm().Prefix
	n.Name = n.Prefix + "-" + Infrastructure.Name + "-" + n.Name

	if n.Mac == "" {
		if RunningInfrastructure != nil {
			for _, net := range RunningInfrastructure.Nets {
				if net.Name == n.Name && net.Kvm == n.Kvm {
					n.Mac = net.Mac
				}
			}
		} else {
			n.Mac = generateMAC()
		}
	}

	if n.IP == "" {
		if RunningInfrastructure != nil {
			for _, net := range RunningInfrastructure.Nets {
				if net.Name == n.Name && net.Kvm == n.Kvm {
					n.IP = net.IP
				}
			}
		} else {
			base := ""
			maxIP := 0
			for _, net := range Infrastructure.Nets {
				if net.Kvm == n.Kvm {
					if net.IP != "" {
						ip := parser.ParseIP(net.IP).To4()
						if int(ip[2]) > maxIP {
							maxIP = int(ip[2])
						}
						base = strconv.Itoa(int(ip[0])) + "." + strconv.Itoa(int(ip[1])) + "."
					}
				}
			}

			n.IP = base + strconv.Itoa(maxIP+1) + ".1"
			n.Gateway = n.IP
		}
	}

	if n.Gateway == "" {
		if RunningInfrastructure != nil {
			for _, net := range RunningInfrastructure.Nets {
				if net.Name == n.Name && net.Kvm == n.Kvm {
					n.Gateway = net.Gateway
				}
			}
		} else {
			n.Gateway = n.IP
		}
	}

	if n.Mask == "" {
		if RunningInfrastructure != nil {
			for _, net := range RunningInfrastructure.Nets {
				if net.Name == n.Name && net.Kvm == n.Kvm {
					n.Mask = net.Mask
				}
			}
		} else {
			n.Mask = Infrastructure.DefMask
			for _, kvm := range Infrastructure.Kvms {
				if n.Kvm == kvm.HostName {
					n.Mask = kvm.DefMask
				}
			}
		}
	}

	if n.DNS == "" {
		if RunningInfrastructure != nil {
			for _, net := range RunningInfrastructure.Nets {
				if net.Name == n.Name && net.Kvm == n.Kvm {
					n.DNS = net.DNS
				}
			}
		} else {
			n.DNS = Infrastructure.DefDNS
			for _, kvm := range Infrastructure.Kvms {
				if n.Kvm == kvm.HostName {
					n.DNS = kvm.DefDNS
				}
			}
		}
	}

}

func (n *Network) createNetwork() {
	if n.GetNetKvm().existNetwork(n.Name) {
		log.Println("-", n.Name, "network already running in", n.GetNetKvm().HostName)
	} else {
		log.Println("+ Starting", n.Name, "network at", n.GetNetKvm().HostName)
		n.GetNetKvm().createNetworkXML(n.generateXML())
	}
}

func (n *Network) deleteNetwork() {
	if n.GetNetKvm().existNetwork(n.Name) {
		log.Println("+ Stopping", n.Name, "network from", n.GetNetKvm().HostName)
		n.GetNetKvm().deleteNetwork(n.Name)
	} else {
		log.Println("-", n.Name, "network is not running in", n.GetNetKvm().HostName)
	}
}

func (n *Network) GetNetKvm() *Kvm {
	for _, element := range Infrastructure.Kvms {
		if n.Kvm == element.HostName {
			return &element
		}
	}
	return nil
}

func (n *Network) generateXML() string {

	netcfg := &libvirtxml.Network{
		Name: n.Name,
		Forward: &libvirtxml.NetworkForward{
			Mode: "nat"},
		MAC: &libvirtxml.NetworkMAC{
			Address: n.Mac}}

	if n.DhcpStartIP != "" {
		if n.DhcpEndIP == "" {
			log.Fatal("Dhcp end UP for network " + n.Name + "not especified")
		}

		netcfg.IPs = append(netcfg.IPs, libvirtxml.NetworkIP{
			Address: n.IP,
			Netmask: n.Mask,
			DHCP: &libvirtxml.NetworkDHCP{
				Ranges: []libvirtxml.NetworkDHCPRange{{
					Start: n.DhcpStartIP,
					End:   n.DhcpEndIP}}}})
	} else {
		netcfg.IPs = append(netcfg.IPs, libvirtxml.NetworkIP{
			Address: n.IP,
			Netmask: n.Mask})
	}

	xml, err := netcfg.Marshal()
	if err != nil {
		panic(err)
	}

	return xml
}

func (n *Network) ShowNetwork() {
	fmt.Println("kvm:\t\t", n.Kvm)
	fmt.Println("name:\t\t", n.Name)
	fmt.Println("mac:\t\t", n.Mac)
	fmt.Println("ip:\t\t", n.IP)
	fmt.Println("mask:\t\t", n.Mask)
	fmt.Println("dhcpstartip:\t", n.DhcpStartIP)
	fmt.Println("dhcpendip:\t", n.DhcpEndIP)
}

func GetMgmtIP() string {
	hostname, _ := os.Hostname()
	for _, n := range Infrastructure.Nets {
		if n.Kvm == "localhost" || n.Kvm == hostname {
			if n.Name == n.Prefix+"-"+Infrastructure.Name+"-mgmt" {
				return n.IP
			}
		}
	}
	return ""
}
