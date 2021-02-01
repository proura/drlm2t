package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

var port = 6060
var myClient = &http.Client{Timeout: 10 * time.Second}

type Config struct {
	Name       string    `mapstructure:"name"`
	HostConfig Host      `mapstructure:"host"`
	NetConfig  []Network `mapstructure:"nets"`
}

type Host struct {
	Name string `mapstructure:"name"`
}

type Network struct {
	Name      string `mapstructure:"name"`
	Mac       string `mapstructure:"mac"`
	IP        string `mapstructure:"ip"`
	Mask      string `mapstructure:"mask"`
	Gateway   string `mapstructure:"gateway"`
	DNS       string `mapstructure:"dns"`
	Interface string `mapstructure:"interface"`
}

var cfg *Config
var Gateway string

func main() {

	// Mirem si s'ha especificat el Gataway
	if len(os.Args) == 1 {
		log.Fatalln("You need to specificy the Gateway")
	} else {
		Gateway = os.Args[1]
	}

	// Si no hi ha el fitxer de config el descarregem
	if !fileExists("./drlm2t.cfg") {
		for generateConfig() == "" {
			time.Sleep(1 * time.Second)
		}
		err := ioutil.WriteFile("drlm2t.cfg", []byte(generateConfig()), 0755)
		if err != nil {
			log.Printf("Unable to write file: %v", err)
		}
	} else {
		for generateConfig() == "" {
			time.Sleep(1 * time.Second)
		}
	}

	urlGet := "http://" + Gateway + ":" + strconv.Itoa(port) + "/tests/" + cfg.Name + "/tests/" + cfg.HostConfig.Name
	urlPost := "http://" + Gateway + ":" + strconv.Itoa(port) + "/upload/" + cfg.Name + "/tests/" + cfg.HostConfig.Name

	for {

		//////////////////////////////////////////////////////////////////
		/////////Mirem si hi ha tests per decarregar//////////////////////
		log.Println("Check for tests to run...")

		doc, err := htmlquery.LoadURL(urlGet)
		if err == nil {
			nodes, err := htmlquery.QueryAll(doc, "//a")
			if err == nil {
				for _, n := range nodes {
					test := htmlquery.InnerText(n)
					file := strings.Split(test, ".")

					if len(file) > 1 && file[1] == "test" {
						if !fileExists(file[0] + ".output") {
							err := DownloadFile(test, urlGet+"/"+test)
							if err != nil {
								log.Panicln(err)
							}
							log.Println("Test " + test + " downloaded")

							log.Println("Running test " + test)
							cmd := exec.Command("/bin/bash", "-c", "source "+test+" > "+file[0]+".output")
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr
							err = cmd.Run()
							if err != nil {
								log.Println("cmd.Run() failed with %s\n", err)
							}
						}

						//////////////////////////////////////////////////////////////////
						////////Mirem si hi ha resultat per pujar/////////////////////////
						log.Println("Uploading results of " + test)
						err = UploadFile(file[0]+".output", urlPost+"/"+file[0]+".output")
						if err != nil {
							log.Panicln(err)
						}

						// Once test is done ant updloaded delete.
						deleteFile(test)
						deleteFile(file[0] + ".output")
					}

				}
			}
		}

		//////////////////////////////////////////////////////////////////
		////////Stop abans seguent test          /////////////////////////
		time.Sleep(1 * time.Second)
	}

}

// func keepLines(s string, n int) string {
// 	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
// 	return strings.Replace(result, "\r", "", -1)
// }

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func generateConfig() string {

	err := getJSON("http://"+Gateway+":"+strconv.Itoa(port)+"/config/", &cfg)
	if err != nil {
		log.Println("- URL " + "http://" + Gateway + ":" + strconv.Itoa(port) + "/config/" + " no reached")
	} else {
		configuration := ""
		configuration = configuration + "TESTNAME " + cfg.Name + "\n"
		configuration = configuration + "HOSTNAME " + cfg.HostConfig.Name + "\n"

		localInterfaces, _ := net.Interfaces()
		for _, n := range cfg.NetConfig {
			for _, nl := range localInterfaces {
				if n.Mac == nl.HardwareAddr.String() {
					n.Interface = nl.Name
				}
			}
			mask, _ := net.IPMask(net.ParseIP(n.Mask).To4()).Size()
			configuration = configuration + "NETWORK " + n.Name + " " + n.Interface + " " + n.IP + " " + n.Mask + " " + strconv.Itoa(mask) + " " + n.Gateway + " " + n.Mac + " " + n.DNS + "\n"
		}

		return configuration
	}

	return ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func deleteFile(filename string) {
	if fileExists(filename) {
		err := os.Remove(filename)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(filename + " file deleted")
		}
	} else {
		log.Println(filename + " does not exist")
	}
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func UploadFile(filepath string, url string) error {
	fileUP, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fileUP.Close()

	res, err := http.Post(url, "binary/octet-stream", fileUP)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	ioutil.ReadAll(res.Body)
	return err
}

//eval $(./drlm2t-client $(ip r | grep default | awk '{ print $3 }'))
//go build && scp -i ../cfg/sshkey rlm2t-client root@192.168.75.148:~

//=====New Version==========================
//==========================================
// [Unit]
// Description=drlm2t-client
// Wants=network-online.target
// After=network-online.target

// [Service]
// Type=idle
// Restart=on-failure
// RestartSec=5s
// WorkingDirectory=/root
// ExecStartPre=/usr/bin/sleep 5
// ExecStart=/root/drlm2t-client 192.168.75.1

// [Install]
// WantedBy=default.target
