package model

import (
	"io/ioutil"
	"log"
	parser "net"
	"os"
	"path/filepath"
	"strings"

	"github.com/proura/drlm2t/cfg"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var Infrastructure *DRLMTestingConfig
var RunningInfrastructure *DRLMTestingConfig

type DRLMTestingConfig struct {
	Name        string `mpastructure:"name"`
	Description string `mpastructure:"description"`
	Prefix      string `mapstructure:"prefix"`
	Templates   string `mpastructure:"templates"`
	URL         string `mpastructure:"url"`
	DefIP       string `mpastructure:"defip"`
	DefMask     string `mpastructure:"defmask"`
	DefDNS      string `mpastructure:"defdns"`
	DefTem      string `mpastructure:"deftmp"`

	Kvms  []Kvm     `mpastructure:"kvms"`
	Nets  []Network `mapstructure:"nets"`
	Hosts []Host    `mapstructure:"hosts"`
}

func LoadInfrastructure(cfgName string) {
	viper.SetConfigName("infrastructure")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfg.Config.Drlm2tPath + "/tests/" + cfgName)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&Infrastructure)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func LoadRunningInfrastructure(cfgName string) {
	if fileExists(cfg.Config.Drlm2tPath + "/tests/" + cfgName + "/running.yaml") {
		viper.SetConfigName("running")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(cfg.Config.Drlm2tPath + "/tests/" + cfgName)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}
		err := viper.Unmarshal(&RunningInfrastructure)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
	}
}

func InitInfrastructure(cfgName string) {
	// Set Default Global cfgName
	if Infrastructure.Name == "" {
		Infrastructure.Name = cfgName
	}
	// Set Default Global Prefix
	if Infrastructure.Prefix == "" {
		Infrastructure.Prefix = "drlm2t"
	}
	// Set Default Global path to qcow2 files (with libvirt permissions)
	if Infrastructure.Templates == "" {
		Infrastructure.Templates, _ = filepath.Abs(cfg.Config.Drlm2tPath + "/drlm2t-templates")
	}

	// Set Default Global path to qcow2 files (with libvirt permissions)
	if Infrastructure.URL == "" {
		Infrastructure.URL = "http://www.etnalan.es/"
	}

	// Set Default Global starting IP
	if Infrastructure.DefIP == "" {
		Infrastructure.DefIP = "192.168.75.0"
	}
	// Set Defautl Global mask
	if Infrastructure.DefMask == "" {
		Infrastructure.DefMask = "255.255.255.0"
	}
	// Set Defautl Global DNS
	if Infrastructure.DefDNS == "" {
		Infrastructure.DefDNS = "8.8.8.8"
	}
	// Set Defautl Template
	if Infrastructure.DefTem == "" {
		Infrastructure.DefTem = "deb9"
	}
	// Add default kvm in networks without kvm
	for index, net := range Infrastructure.Nets {
		if net.Kvm == "" {
			Infrastructure.Nets[index].Kvm = "localhost"
		}
	}
	// Add default kvm in hosts without kvm
	for index, host := range Infrastructure.Hosts {
		if host.Kvm == "" {
			Infrastructure.Hosts[index].Kvm = "localhost"
		}
	}
	// Look for default KVM configuration in nets and host and insert if needed
	insertDefaultKvm := false
	for _, net := range Infrastructure.Nets {
		if net.Kvm == "localhost" {
			insertDefaultKvm = true
		}
	}
	for _, host := range Infrastructure.Hosts {
		if host.Kvm == "localhost" {
			insertDefaultKvm = true
		}
	}
	for _, kvm := range Infrastructure.Kvms {
		if kvm.HostName == "localhost" {
			insertDefaultKvm = false
		}
	}
	// Add Kvm default struct because is needed and don't exist
	if insertDefaultKvm {
		var kvm Kvm
		kvm.HostName = "localhost"
		Infrastructure.Kvms = append(Infrastructure.Kvms, kvm)
	}
	// Look for network kvm names without configuration in kvm section
	for _, net := range Infrastructure.Nets {
		found := false
		for _, kvm := range Infrastructure.Kvms {
			if net.Kvm == kvm.HostName {
				found = true
			}
		}
		if !found {
			log.Fatalf("Kvm name \"" + net.Kvm + "\" specified in network \"" + net.Name + "\" not found!")
		}
	}
	// Look for host kvm names without configuration in kvm section
	for _, host := range Infrastructure.Hosts {
		found := false

		for _, kvm := range Infrastructure.Kvms {
			if host.Kvm == kvm.HostName {
				found = true
			}
		}
		if !found {
			log.Fatalf("Kvm name \"" + host.Kvm + "\" specified in vm \"" + host.Name + "\" not found!")
		}
	}
	// Initlilize KVMs
	for index := range Infrastructure.Kvms {
		Infrastructure.Kvms[index].initKvm()
	}

	// Add management network for each KVM
	for _, kvm := range Infrastructure.Kvms {
		foundMgmtNet := false
		for _, netikvm := range Infrastructure.Nets {
			if netikvm.Kvm == kvm.HostName {
				if strings.HasSuffix(netikvm.Name, "-mgmt") {
					foundMgmtNet = true
				}
			}
		}
		if !foundMgmtNet {
			var net Network
			net.Name = "mgmt"
			net.Kvm = kvm.HostName
			ip := parser.ParseIP(kvm.DefIP).To4()
			ip[3] = 1
			net.IP = ip.String()
			net.Gateway = ip.String()
			net.DNS = kvm.DefDNS
			ip[3] = 50
			net.DhcpStartIP = ip.String()
			ip[3] = 200
			net.DhcpEndIP = ip.String()
			log.Println("Appen1 net =========> " + net.Name)
			Infrastructure.Nets = append(Infrastructure.Nets, net)
		}
	}
	// Add management network in hosts
	for index := range Infrastructure.Hosts {
		foundMgmtNet := false
		for _, netInHost := range Infrastructure.Hosts[index].Nets {
			if strings.HasSuffix(netInHost.Name, "-mgmt") {
				foundMgmtNet = true
			}
		}
		if !foundMgmtNet {
			var net Network
			net.Name = "mgmt"
			log.Println("Appen2 net =========> " + net.Name)

			Infrastructure.Hosts[index].Nets = append(Infrastructure.Hosts[index].Nets, net)
		}
	}
	// Add default network in hosts without networks
	for index, host := range Infrastructure.Hosts {
		if len(host.Nets) == 1 {
			var net Network
			net.Name = "default"
			log.Println("Appen3 net =========> " + net.Name)

			Infrastructure.Hosts[index].Nets = append(Infrastructure.Hosts[index].Nets, net)
		}
	}
	// Add networks found in host if not exist
	for _, host := range Infrastructure.Hosts {
		for _, nethost := range host.Nets {
			found := false
			//mirar si existeix i si no existeix crear
			for _, net := range Infrastructure.Nets {
				if net.Name == nethost.Name && net.Kvm == host.Kvm {
					found = true
				}
				if strings.HasSuffix(net.Name, nethost.Name) && net.Kvm == host.Kvm {
					found = true
				}
				if strings.HasSuffix(nethost.Name, net.Name) && net.Kvm == host.Kvm {
					found = true
				}

			}
			if !found {
				var net Network
				net.Name = nethost.Name
				net.Kvm = host.Kvm
				log.Println("Appen4 net =========> " + net.Name)

				Infrastructure.Nets = append(Infrastructure.Nets, net)
			}
		}
	}
	// Initlilize Networks
	for index := range Infrastructure.Nets {
		Infrastructure.Nets[index].initNetwork(index)
	}
	// Initialize Hosts
	for index := range Infrastructure.Hosts {
		Infrastructure.Hosts[index].initHost(index)
	}

}

func SaveRunningIfrastructure() {
	running, _ := yaml.Marshal(Infrastructure)
	//Save Infrastructure to running.yaml file
	err := ioutil.WriteFile(cfg.Config.Drlm2tPath+"/tests/"+Infrastructure.Name+"/running.yaml", running, 0644)
	if err != nil {
		log.Fatal("- Error saving " + cfg.Config.Drlm2tPath + "/tests/" + Infrastructure.Name + "/running.yaml file")
	}
}

func (infra *DRLMTestingConfig) CreateNetworks() {
	log.Println("STARTING NETWORKS")

	hostname, _ := os.Hostname()

	for index := range infra.Nets {
		if infra.Nets[index].Kvm == "localhost" || infra.Nets[index].Kvm == hostname {
			infra.Nets[index].createNetwork()
		}
	}
}

func (infra *DRLMTestingConfig) DeleteNetworks() {
	log.Println("STOPPING NETWORKS")

	hostname, _ := os.Hostname()

	for index := range infra.Nets {
		if infra.Nets[index].Kvm == "localhost" || infra.Nets[index].Kvm == hostname {
			infra.Nets[index].deleteNetwork()
		}
	}
}

func (infra *DRLMTestingConfig) CreateHosts() {
	log.Println("STARTING HOSTS")

	hostname, _ := os.Hostname()

	for index := range infra.Hosts {
		if infra.Hosts[index].Kvm == "localhost" || infra.Hosts[index].Kvm == hostname {
			infra.Hosts[index].createHost()
		}
	}
}

func (infra *DRLMTestingConfig) DeleteHosts() {
	log.Println("STOPPING HOSTS")

	hostname, _ := os.Hostname()

	if infra.Hosts == nil {
		return
	}

	for index := range infra.Hosts {
		if infra.Hosts[index].Kvm == "localhost" || infra.Hosts[index].Kvm == hostname {
			infra.Hosts[index].deleteHost()
		}
	}
}

func (infra *DRLMTestingConfig) Clean() {
	log.Println("DELETING RUNNING FILES")

	if Infrastructure == nil {
		log.Println("- No Infrastructure to clean")
		return
	}

	hostname, _ := os.Hostname()

	//Delete qcow2 file for each host in Infrastructure
	for _, host := range infra.Hosts {
		if host.Kvm == "localhost" || host.Kvm == hostname {
			host.deleteQCOW2()
		}
	}

	for _, kvm := range infra.Kvms {
		if kvm.HostName == "localhost" || kvm.HostName == hostname {
			kvm.deleteDirs()
		}
	}

	cleanTests()

	//Delete running.yaml file
	if fileExists(cfg.Config.Drlm2tPath + "/tests/" + infra.Name + "/running.yaml") {
		err := os.Remove(cfg.Config.Drlm2tPath + "/tests/" + infra.Name + "/running.yaml")
		if err != nil {
			log.Println("-", err)
		} else {
			log.Println("+ " + cfg.Config.Drlm2tPath + "/tests/" + infra.Name + "/running.yaml file deleted")
		}
	} else {
		log.Println("- " + cfg.Config.Drlm2tPath + "/tests/" + infra.Name + "/running.yaml does not exist")
	}

}

func (infra *DRLMTestingConfig) CountLocalHosts() int {
	hosts := 0
	for _, h := range Infrastructure.Hosts {
		hostname, _ := os.Hostname()
		if h.Kvm == "localhost" || h.Kvm == hostname {
			hosts++
		}
	}
	return hosts
}

func (infra *DRLMTestingConfig) GetLocalHosts() []Host {
	var hosts []Host
	for _, h := range Infrastructure.Hosts {
		hostname, _ := os.Hostname()
		if h.Kvm == "localhost" || h.Kvm == hostname {
			hosts = append(hosts, h)
		}
	}
	return hosts
}
