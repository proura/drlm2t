package model

import (
	"log"
	"time"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
	libvirt "libvirt.org/libvirt-go"
)

//Kvm struct to store KVM server settings
type Kvm struct {
	HostName  string `mapstructure:"hostname"`
	User      string `mapstructure:"user"`
	URI       string `mapstructure:"uri"`
	Prefix    string `mapstructure:"prefix"`
	Templates string `mapstructure:"templates"`
	DefIP     string `mpastructure:"defip"`
	DefMask   string `mpastructure:"defmask"`
	DefDNS    string `mpastructure:"defdns"`
	DefTmp    string `mpastructure:"deftmp"`
}

var conn *libvirt.Connect

func (k *Kvm) initKvm() {
	if k.HostName == "" {
		k.HostName = "localhost"
	}

	if k.URI == "" {
		k.URI = "qemu:///system"
	}

	if k.Prefix == "" {
		k.Prefix = Infrastructure.Prefix
	}

	if k.Templates == "" {
		k.Templates = Infrastructure.Templates
	}

	if k.DefIP == "" {
		k.DefIP = Infrastructure.DefIP
	}

	if k.DefMask == "" {
		k.DefMask = Infrastructure.DefMask
	}

	if k.DefDNS == "" {
		k.DefDNS = Infrastructure.DefDNS
	}

	if k.DefTmp == "" {
		k.DefTmp = Infrastructure.DefTem
	}
}

func (k Kvm) connect() {
	var err error
	conn, err = libvirt.NewConnect(k.URI)
	if err != nil {
		log.Fatal("Error connectant a ", k.URI)
	}
}

func (k Kvm) close() {
	conn.Close()
}

func (k Kvm) existNetwork(name string) bool {
	k.connect()
	defer k.close()

	net, err := conn.LookupNetworkByName(name)

	if err != nil {
		return false
	}

	net.Free()
	return true
}

func (k Kvm) createNetworkXML(xml string) {
	k.connect()
	defer k.close()

	_, err := conn.NetworkCreateXML(xml)
	if err != nil {
		log.Println(err)
		log.Fatal("Error creant la xarxa!")
	}
}

func (k Kvm) deleteNetwork(name string) {
	k.connect()
	defer k.close()

	net, err := conn.LookupNetworkByName(name)
	if err != nil {
		log.Println(err)
		log.Fatal("Error eliminant la xarxa!")
	}

	net.Destroy()
	net.Free()
}

func (k Kvm) existHost(name string) bool {
	k.connect()
	defer k.close()

	dom, err := conn.LookupDomainByName(name)

	if err != nil {
		return false
	}

	dom.Free()
	return true
}

func (k Kvm) createHostXML(xml string) {
	k.connect()
	defer k.close()

	_, err := conn.DomainCreateXML(xml, libvirt.DOMAIN_NONE)
	if err != nil {
		log.Println("-", err)
		log.Fatal("- Error creating host!")
	}
}

func (k Kvm) deleteHost(name string) {
	k.connect()
	defer k.close()

	host, err := conn.LookupDomainByName(name)
	if err != nil {
		log.Println("-", err)
		log.Fatal("- Error stopping Host!")
	}

	host.Destroy()
	host.Free()
}

func (k Kvm) createSnap(host, test string) {
	k.connect()
	defer k.close()

	domain, _ := conn.LookupDomainByName(host)
	domcfg := &libvirtxml.DomainSnapshot{}
	domcfg.Name = test

	xml, err := domcfg.Marshal()
	if err != nil {
		panic(err)
	}

	_, err = domain.CreateSnapshotXML(xml, libvirt.DOMAIN_SNAPSHOT_CREATE_DISK_ONLY)
	if err != nil {
		log.Println(err)
	}
}

func (k Kvm) checkMode(host string, testMode Mode) {
	k.connect()
	defer k.close()

	domain, _ := conn.LookupDomainByName(host)

	domXML, _ := domain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	domcfg := &libvirtxml.Domain{}
	domcfg.Unmarshal(domXML)

	nnets := len(domcfg.Devices.Interfaces)

	if domcfg.Devices.Disks[0].Boot.Order == 1 {
		if testMode != 0 {
			log.Println("+ " + host + ": Switching from Normal mode to Netboot mode")
			domain.Destroy()

			state, _, _ := domain.GetState()
			for state != 0 {
				time.Sleep(1 * time.Second)
				state, _, _ = domain.GetState()
			}

			domcfg.Devices.Disks[0].Boot.Order = uint(nnets + 1)
			domcfg.Devices.Interfaces[nnets-1].Boot.Order = 1
		}

	} else {
		if testMode != 1 {
			log.Println("+ " + host + ": Switching from Netboot mode to Normal mode")
			domain.Destroy()

			state, _, _ := domain.GetState()
			for state != 0 {
				time.Sleep(1 * time.Second)
				state, _, _ = domain.GetState()
			}

			domcfg.Devices.Interfaces[nnets-1].Boot.Order = uint(nnets + 1)
			domcfg.Devices.Disks[0].Boot.Order = 1
		}

	}

	state, _, _ := domain.GetState()
	if state == 0 {
		xml, err := domcfg.Marshal()
		if err != nil {
			panic(err)
		}
		_, err = conn.DomainCreateXML(xml, libvirt.DOMAIN_NONE)
		if err != nil {
			log.Println("-", err)
			log.Fatal("- Error creating host!")
		}
	}
}

func (k Kvm) deleteDirs() {
	if k.Templates != "" && Infrastructure.Name != "" && dirExists(k.Templates+"/"+Infrastructure.Name) {
		removeDir(k.Templates + "/" + Infrastructure.Name)
		log.Println("+ Test dir " + k.Templates + "/" + Infrastructure.Name + " deleted")
	} else {
		log.Println("- Test dir " + k.Templates + "/" + Infrastructure.Name + " no exists")
	}
}

func (k Kvm) getHostIP(host string) string {
	k.connect()
	defer k.close()

	var ifaces libvirt.DomainInterfaceAddressesSource

	domain, err := conn.LookupDomainByName(host)
	if err != nil {
		log.Println("-", err)
		log.Println("- Error looking for Host!")
	}

	inter, err := domain.ListAllInterfaceAddresses(ifaces)
	if len(inter) == 0 {
		for _, h := range Infrastructure.GetLocalHosts() {
			if h.Name == host {
				for _, n := range h.Nets {
					if n.Name == n.Prefix+"-"+Infrastructure.Name+"-mgmt" {
						return n.IP
					}
				}
			}
		}
		return ""
	}

	return inter[0].Addrs[0].Addr
}

// GetHostByIP returnt de the host name giving and IP
func (k Kvm) GetHostByIP(IP string) string {
	k.connect()
	defer k.close()

	runningHosts, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		log.Fatalln("Error connecting KVM to get running Hosts")
	}

	name := ""

	for _, h := range runningHosts {
		var ifaces libvirt.DomainInterfaceAddressesSource
		inter, _ := h.ListAllInterfaceAddresses(ifaces)
		if len(inter) == 0 {
			name = ""
		} else if inter[0].Addrs[0].Addr == IP {
			name, _ := h.GetName()
			return name
		}
	}

	if name == "" {
		for _, h := range Infrastructure.Hosts {
			for _, n := range h.Nets {
				if n.IP == IP {
					name = h.Name
				}
			}
		}
	}

	return name
}
