package model

import (
	"fmt"
	"log"
	parser "net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/proura/drlm2t/cfg"
)

type Host struct {
	Name     string    `mapstructure:"name"`
	Kvm      string    `mapstructure:"kvm"`
	Template string    `mapstructure:"template"`
	Prefix   string    `mapstructure:"prefix"`
	Nets     []Network `mapstructure:"nets"`
	Tests    []Test    `mapstructure:"tests"`
}

func (h *Host) initHost(index int) {

	h.Prefix = h.GetHostKvm().Prefix
	if h.Template == "" {
		h.Template = h.GetHostKvm().DefTmp
	}

	if strings.HasPrefix(h.Name, h.Prefix+"-"+Infrastructure.Name) {
		//n.Name = n.Name
	} else {
		h.Name = h.Prefix + "-" + Infrastructure.Name + "-" + h.Name
	}

	for i, net := range h.Nets {
		h.Nets[i].Name = h.Prefix + "-" + Infrastructure.Name + "-" + h.Nets[i].Name
		if net.Mac == "" {
			if RunningInfrastructure != nil {
				//loking for host name in running.yaml file
				for _, host := range RunningInfrastructure.Hosts {
					if host.Name == h.Name && host.Kvm == h.Kvm {
						//loking for nets name inside host in running.yaml file
						for _, hostnet := range host.Nets {
							if hostnet.Name == h.Nets[i].Name {
								h.Nets[i].Mac = hostnet.Mac
							}
						}
					}
				}
			}
			if h.Nets[i].Mac == "" {
				h.Nets[i].Mac = generateMAC()
			}
		}
		if net.Mask == "" {
			for _, inet := range Infrastructure.Nets {
				if h.Nets[i].Name == inet.Name && h.Kvm == inet.Kvm {
					h.Nets[i].Mask = inet.Mask
				}
			}
		}
		if net.DNS == "" {
			for _, inet := range Infrastructure.Nets {
				if h.Nets[i].Name == inet.Name && h.Kvm == inet.Kvm {
					h.Nets[i].DNS = inet.DNS
				}
			}
		}
		if net.Gateway == "" {
			for _, inet := range Infrastructure.Nets {
				if h.Nets[i].Name == inet.Name && h.Kvm == inet.Kvm {
					h.Nets[i].Gateway = inet.Gateway
				}
			}
		}
		if net.Prefix == "" {
			for _, inet := range Infrastructure.Nets {
				if h.Nets[i].Name == inet.Name && h.Kvm == inet.Kvm {
					h.Nets[i].Prefix = inet.Prefix
				}
			}
		}
		if net.IP == "" {
			if RunningInfrastructure != nil {
				//loking for host name in running.yaml file
				for _, host := range RunningInfrastructure.Hosts {
					if host.Name == h.Name && host.Kvm == h.Kvm {
						//loking for nets name inside host in running.yaml file
						for _, hostnet := range host.Nets {
							if hostnet.Name == h.Nets[i].Name {
								h.Nets[i].IP = hostnet.IP
							}
						}
					}
				}
			}
			if h.Nets[i].IP == "" {
				base := ""
				maxIP := 0

				for _, inet := range Infrastructure.Nets {
					if h.Nets[i].Name == inet.Name && h.Kvm == inet.Kvm {
						ip := parser.ParseIP(inet.IP).To4()
						if int(ip[3]) > maxIP {
							maxIP = int(ip[3])
						}
						base = strconv.Itoa(int(ip[0])) + "." + strconv.Itoa(int(ip[1])) + "." + strconv.Itoa(int(ip[2])) + "."
					}
				}

				for _, host := range Infrastructure.Hosts {
					if h.Kvm == host.Kvm {
						for _, inet := range host.Nets {
							if h.Nets[i].Name == inet.Name {
								if inet.IP != "" {
									ip := parser.ParseIP(inet.IP).To4()
									if int(ip[3]) > maxIP {
										maxIP = int(ip[3])
									}
								}
							}
						}
					}
				}

				h.Nets[i].IP = base + strconv.Itoa(maxIP+1)
			}
		}
	}
}

func (h *Host) createHost() {

	if h.GetHostKvm().existHost(h.Name) {
		log.Println("-", h.Name, "host already running")
	} else {
		log.Println("+ Starting", h.Name, "host at", h.GetHostKvm().HostName)
		xml := h.generateXML()
		h.createQCOW2()
		h.GetHostKvm().createHostXML(xml)
	}
}

func (h *Host) deleteHost() {
	if h.GetHostKvm().existHost(h.Name) {
		log.Println("+ Stopping", h.Name, "host from", h.GetHostKvm().HostName)
		h.GetHostKvm().deleteHost(h.Name)
	} else {
		log.Println("-", h.Name, "host is not running!")
	}
}

func (h *Host) createQCOW2() {
	templates := h.GetHostKvm().Templates

	if !fileExists(templates + "/" + Infrastructure.Name + "/" + h.Name + ".qcow2") {
		// If templates dir does not exist, create it
		if !dirExists(templates + "/templates") {
			createDir(templates + "/templates")
		}

		// Check if template exists
		if !fileExists(templates + "/templates/" + h.Template + "/" + h.Template + ".qcow2") {
			fileURL := Infrastructure.URL + "/" + h.Template + ".tar.gz"
			fileDst := templates + "/templates/" + h.Template + ".tar.gz"

			// Check if exists tar.gz downloaded file, if not download
			if !fileExists(templates + "/templates/" + h.Template + ".tar.gz") {
				log.Println("+ " + h.Template + " not found, will be downloaded from the server")
				err := downloadFile(fileDst, fileURL)
				if err != nil {
					log.Fatal("- Template " + h.Template + " not found to download in \"" + fileURL + "/templates/\".")
				}
			}
			// Extract tar.gz file
			log.Println("+ Uncompressing file " + fileDst)
			ExtractTarGz(templates+"/templates", fileDst)
		}

		// Create destination Dir for test if no exists
		if !dirExists(templates + "/" + Infrastructure.Name) {
			createDir(templates + "/" + Infrastructure.Name)
		}

		// Crear la imatge base de disc a lloc
		cmd := exec.Command("/bin/bash", "-c", "qemu-img create -f qcow2 -b "+templates+"/templates/"+h.Template+"/"+h.Template+".qcow2 "+templates+"/"+Infrastructure.Name+"/"+h.Name+".qcow2 > /dev/null 2>&1")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println("cmd.Run() failed with %s\n", err)
		}

		if err := os.Chmod(templates+"/"+Infrastructure.Name+"/"+h.Name+".qcow2", 0666); err != nil {
			log.Fatal(err)
		}

		// Copiar el testing de config a lloc
		execCopy(templates+"/templates/"+h.Template+"/tests/config.test", cfg.Config.Drlm2tPath+"/tests/"+Infrastructure.Name+"/tests/"+h.Name+"/0-config.test")

	}
}

func (h *Host) deleteQCOW2() {
	templates := h.GetHostKvm().Templates

	err := os.Remove(templates + "/" + Infrastructure.Name + "/" + h.Name + ".qcow2")
	if err != nil {
		log.Println("-", err)
	} else {
		log.Println("+ File " + templates + "/" + Infrastructure.Name + "/" + h.Name + ".qcow2 file deleted")
	}
}

func (h *Host) generateXML() string {

	templates := h.GetHostKvm().Templates

	qcowFile := templates + "/" + Infrastructure.Name + "/" + h.Name + ".qcow2"
	backQcowFile := templates + "/templates/" + h.Template + "/" + h.Template + ".qcow2"

	domcfg := &libvirtxml.Domain{
		Type: "kvm",
		Name: h.Name,
		Memory: &libvirtxml.DomainMemory{
			Value: 2048,
			Unit:  "MB"},
		VCPU: &libvirtxml.DomainVCPU{
			Value: 1},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Type: "hvm"}},
		Features: &libvirtxml.DomainFeatureList{
			ACPI: &libvirtxml.DomainFeature{},
		},
		Devices: &libvirtxml.DomainDeviceList{
			Disks:      []libvirtxml.DomainDisk{},
			Interfaces: []libvirtxml.DomainInterface{},
			Graphics: []libvirtxml.DomainGraphic{
				{
					Spice: &libvirtxml.DomainGraphicSpice{},
				}},
			Videos: []libvirtxml.DomainVideo{
				{
					Model: libvirtxml.DomainVideoModel{
						Type: "qxl"},
				}}}}

	var slot uint
	slot = 3

	for i, net := range h.Nets {
		xmlNet := libvirtxml.DomainInterface{
			MAC: &libvirtxml.DomainInterfaceMAC{
				Address: net.Mac},
			Source: &libvirtxml.DomainInterfaceSource{
				Network: &libvirtxml.DomainInterfaceSourceNetwork{
					Network: net.Name}},
			Boot: &libvirtxml.DomainDeviceBoot{
				Order: uint(i + 2)}}

		if net.Name == net.Prefix+"-"+Infrastructure.Name+"-"+"mgmt" {
			xmlNet.Address = &libvirtxml.DomainAddress{
				PCI: &libvirtxml.DomainAddressPCI{
					Slot: &slot}}
		}

		domcfg.Devices.Interfaces = append(domcfg.Devices.Interfaces, xmlNet)
	}

	//Get how many tests are done.
	testsDone := 0
	for _, test := range h.Tests {
		if test.Status == 1 || (test.TestType == 2 && test.Status == -1) {
			testsDone++
		} else {
			break
		}
	}

	if testsDone <= 1 {
		xmlDisk := libvirtxml.DomainDisk{
			Driver: &libvirtxml.DomainDiskDriver{
				Name: "qemu",
				Type: "qcow2"},
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: qcowFile}},
			Target: &libvirtxml.DomainDiskTarget{
				Dev: "hda",
				Bus: "ide"},
			Boot: &libvirtxml.DomainDeviceBoot{
				Order: 1},
			BackingStore: &libvirtxml.DomainDiskBackingStore{
				Format: &libvirtxml.DomainDiskFormat{
					Type: "qcow2"},
				Source: &libvirtxml.DomainDiskSource{
					File: &libvirtxml.DomainDiskSourceFile{
						File: backQcowFile}}}}
		domcfg.Devices.Disks = append(domcfg.Devices.Disks, xmlDisk)
	} else {
		xmlBackingStoreOld := &libvirtxml.DomainDiskBackingStore{
			Format: &libvirtxml.DomainDiskFormat{
				Type: "qcow2"},
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: qcowFile}},
			BackingStore: &libvirtxml.DomainDiskBackingStore{
				Format: &libvirtxml.DomainDiskFormat{
					Type: "qcow2"},
				Source: &libvirtxml.DomainDiskSource{
					File: &libvirtxml.DomainDiskSourceFile{
						File: backQcowFile}}}}

		for i := 1; i < testsDone-1; i++ {
			actSnap := templates + "/" + Infrastructure.Name + "/" + h.Name + "." + strconv.Itoa(h.Tests[i].Index) + "-" + h.Tests[i].Name
			xmlBackingStore := libvirtxml.DomainDiskBackingStore{
				Format: &libvirtxml.DomainDiskFormat{
					Type: "qcow2"},
				Source: &libvirtxml.DomainDiskSource{
					File: &libvirtxml.DomainDiskSourceFile{
						File: actSnap}},
				BackingStore: xmlBackingStoreOld}

			xmlBackingStoreOld = &xmlBackingStore
		}

		actSnap := templates + "/" + Infrastructure.Name + "/" + h.Name + "." + strconv.Itoa(h.Tests[testsDone-1].Index) + "-" + h.Tests[testsDone-1].Name

		xmlDisk := libvirtxml.DomainDisk{
			Driver: &libvirtxml.DomainDiskDriver{
				Name: "qemu",
				Type: "qcow2"},
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: actSnap}},
			Target: &libvirtxml.DomainDiskTarget{
				Dev: "hda",
				Bus: "ide"},
			Boot: &libvirtxml.DomainDeviceBoot{
				Order: 1},
			BackingStore: xmlBackingStoreOld}

		domcfg.Devices.Disks = append(domcfg.Devices.Disks, xmlDisk)

	}

	xml, err := domcfg.Marshal()
	if err != nil {
		panic(err)
	}

	return xml
}

func (h *Host) ShowHost() {
	fmt.Println("kvm:\t\t", h.Kvm)
	fmt.Println("hostname:\t", h.Name)
	fmt.Println("template:\t", h.Template)
	fmt.Println("networks:")
	for _, element := range h.Nets {
		fmt.Println("    name:\t", element.Name)
		fmt.Println("    ip:\t\t", element.IP)
		fmt.Println("    mac:\t", element.Mac)
	}
}

func (h *Host) GetHostMgmtIP() string {
	return h.GetHostKvm().getHostIP(h.Name)
}

func (h *Host) GetHostKvm() *Kvm {
	for _, element := range Infrastructure.Kvms {
		if h.Kvm == element.HostName {
			return &element
		}
	}
	return nil
}
