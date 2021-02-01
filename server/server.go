package testserver

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/proura/drlm2t/model"
)

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

var runMode string

func RunServer(rm string) {

	runMode = rm

	log.Println("+ Starting server in " + runMode + " mode (http://" + model.GetMgmtIP() + ":6060/tests/)")

	if runMode == "testing" {
		for _, h := range model.Infrastructure.GetLocalHosts() {
			t := h.GetNextTest()
			if t.Name != "" {
				t.PublishTest(h)
			}
		}
	}

	fs := http.FileServer(http.Dir("./tests"))
	http.Handle("/tests/", http.StripPrefix("/tests/", fs))

	http.HandleFunc("/upload/", UploadOutput)
	http.HandleFunc("/config/", ConfigHandler)
	if err := http.ListenAndServe(":6060", nil); err != nil {
		panic(err)
	}
}

func UploadOutput(w http.ResponseWriter, r *http.Request) {

	URL := strings.Split(r.URL.String(), "/")
	testURL := URL[2]
	hostURL := URL[4]
	fileURL := URL[5]

	dst := "tests/" + testURL + "/tests/" + hostURL + "/" + fileURL

	file, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))

	hostname, err := os.Hostname()

	for ih, h := range model.Infrastructure.Hosts {
		if h.Name == hostURL && (h.Kvm == "localhost" || h.Kvm == hostname) {
			for it, t := range h.Tests {
				var splitfile []string
				if t.FileToRun != "" {
					splitfile = strings.Split(t.FileToRun, ".")
				} else {
					splitfile = append(splitfile, t.Name, "test")
				}

				if strconv.Itoa(t.Index)+"-"+t.Name+".output" == fileURL || strconv.Itoa(t.Index)+"-"+splitfile[0]+".output" == fileURL {

					contentOoutput, err := ioutil.ReadFile("tests/" + testURL + "/tests/" + hostURL + "/" + fileURL)
					if err != nil {
						log.Fatal(err)
					}

					// Save output to Intrastructure
					model.Infrastructure.Hosts[ih].Tests[it].Output = strings.TrimSpace(string(contentOoutput))

					if t.TestType == model.TestOutputOnly {
						log.Println("+ " + model.Infrastructure.Hosts[ih].Name + ": Test <" + strconv.Itoa(t.Index) + "-" + splitfile[0] + "> done succesfully")
						model.Infrastructure.Hosts[ih].Tests[it].Status = 1

					} else {
						// Check if output is equal to Expeted
						if model.Infrastructure.Hosts[ih].Tests[it].Output == model.Infrastructure.Hosts[ih].Tests[it].Expect {
							log.Println("+ " + model.Infrastructure.Hosts[ih].Name + ": Test <" + strconv.Itoa(t.Index) + "-" + splitfile[0] + "> done succesfully")
							model.Infrastructure.Hosts[ih].Tests[it].Status = 1
						} else {
							var startLog string
							if model.Infrastructure.Hosts[ih].Tests[it].TestType == model.TestDontStop {
								startLog = "! "
							} else {
								startLog = "- "
							}
							log.Println(startLog + model.Infrastructure.Hosts[ih].Name + ": Test <" + strconv.Itoa(t.Index) + "-" + splitfile[0] + "> NOT give the expected output")
							model.Infrastructure.Hosts[ih].Tests[it].Status = -1
						}
					}

					// Get the content of the test
					contentTest, err := ioutil.ReadFile("tests/" + testURL + "/tests/" + hostURL + "/" + strconv.Itoa(t.Index) + "-" + splitfile[0] + ".test")
					if err != nil {
						log.Fatal(err)
					}
					// Get the content of the test output
					contentExpected := t.Expect

					///RESULT OF TEST/////////////////
					resPrefix := ""
					if model.Infrastructure.Hosts[ih].Tests[it].Status == 1 {
						resPrefix = "OK"
					} else {
						resPrefix = "ERR"
					}
					dst := "tests/" + testURL + "/tests/" + hostURL + "/done/" + strconv.Itoa(t.Index) + "-" + resPrefix + "-" + splitfile[0] + ".test"

					f, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						log.Println(err)
					}
					defer f.Close()

					currentTime := time.Now()

					if _, err := f.WriteString("=======================================================\n" + currentTime.String() + "\n=======================================================\n\n"); err != nil {
						log.Println(err)
					}
					if _, err := f.WriteString("========================TEST===========================\n\n"); err != nil {
						log.Println(err)
					}
					if _, err := f.WriteString(string(contentTest) + "\n\n"); err != nil {
						log.Println(err)
					}
					if _, err := f.WriteString("========================EXPECTED=======================\n\n"); err != nil {
						log.Println(err)
					}
					if _, err := f.WriteString(string(contentExpected) + "\n\n"); err != nil {
						log.Println(err)
					}
					if _, err := f.WriteString("========================OUTPUT=========================\n\n"); err != nil {
						log.Println(err)
					}
					if _, err := f.WriteString(string(contentOoutput) + "\n\n"); err != nil {
						log.Println(err)
					}
					/////////////////////////////////////

					err = os.Remove("tests/" + testURL + "/tests/" + hostURL + "/" + strconv.Itoa(t.Index) + "-" + splitfile[0] + ".test")
					if err != nil {
						log.Println("-", err)
					}
					err = os.Remove("tests/" + testURL + "/tests/" + hostURL + "/" + fileURL)
					if err != nil {
						log.Println("-", err)
					}

					//////////////////////////////////////////////////////////////////
					// Check if there are more tests to run //////////////////////////
					if model.Infrastructure.Hosts[ih].Tests[it].Status == 1 || t.TestType == model.TestDontStop || t.TestType == model.TestOutputOnly {
						if splitfile[0] != "config" && len(model.Infrastructure.Hosts[ih].Tests) > it {
							if len(h.Tests) > it+1 {
								h.Tests[it+1].PublishTest(h)
							}

							//Look if there are some test with dependencies of this test wainting
							for irh, rh := range model.RunningInfrastructure.Hosts {
								oldStatus := 0
								for ith, th := range rh.Tests {
									if len(th.Dependencies) > 0 {
										for _, dh := range th.Dependencies {
											if model.Infrastructure.Prefix+"-"+model.Infrastructure.Name+"-"+dh.Host == h.Name && oldStatus > 0 {
												model.RunningInfrastructure.Hosts[irh].Tests[ith].PublishTest(model.RunningInfrastructure.Hosts[irh])
											}
										}
									}
									oldStatus = th.Status
								}
							}
						}
					}
					//////////////////////////////////////////////////////////////////

					model.SaveRunningIfrastructure()

				}
			}
		}
	}

}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {

	resp := Config{}
	resp.Name = model.Infrastructure.Name

	// Obtenim la IP
	ip := strings.Split(GetIP(r), ":")[0]

	// Mirem a quin host estem executant el testing
	hostkvm, _ := os.Hostname()
	hostname := ""

	for _, k := range model.Infrastructure.Kvms {
		if k.HostName == "localhost" || k.HostName == hostkvm {
			hostname = k.GetHostByIP(ip)
			resp.HostConfig.Name = hostname
		}
	}

	w.Header().Add("Content-Type", "application/json")

	for _, h := range model.Infrastructure.Hosts {
		if h.Name == hostname && (h.Kvm == "localhost" || h.Kvm == hostkvm) {
			for _, n := range h.Nets {
				newNet := Network{
					Name:    n.Name,
					IP:      n.IP,
					Mac:     n.Mac,
					Mask:    n.Mask,
					DNS:     n.DNS,
					Gateway: n.Gateway,
				}
				resp.NetConfig = append(resp.NetConfig, newNet)
			}
		}
	}
	respf, _ := json.Marshal(resp)
	w.Write(respf)
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}
