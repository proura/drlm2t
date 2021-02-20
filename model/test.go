package model

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type Test struct {
	Index        int      `mapstructure:"index"`
	Status       int      `mapstructure:"status"`
	Name         string   `mapstructure:"name"`
	TestType     Tipus    `mapstructure:"testtype"`
	Mode         Mode     `mapstructure:"mode"`
	CommandToRun string   `mapstructure:"commandtorun"`
	FileToRun    string   `mapstructure:"filetorun"`
	Args         []string `mapstructure:"args"`
	Expect       string   `mapstructure:"expect"`
	Output       string   `mapstructure:"output"`
	LandMark     bool     `mapstructure:"landmark"`
	Dependencies []Deps   `mapstructure:"dependencies"`
}

type Deps struct {
	Host string `mapstructure:"host"`
	Test int    `mapstructure:"test"`
}

/*
(0) TestNormal 		=> (default) executa el test i mira si concorda amb l'expected, sinó para.
(1) TestOutputOnly 	=> executa el test i nomes guarda l'output.
(2) TestDontStop	=> executa el test i continua tan si dona error com si no.
*/
type Tipus int

const (
	TestNormal = iota
	TestOutputOnly
	TestDontStop
)

/*
(0) RunNormal	=> (default) Execució del test en mode normal (hdd)
(1) RunRecover	=> Execució del test en mode recover (netboot)
*/
type Mode int

const (
	RunNormal = iota
	RunRecover
)

func (t Tipus) String() string {
	return [...]string{"Normal", "OutputOnly", "DontStop"}[t]
}

func InitTesting(mode int) {

	for index, h := range Infrastructure.Hosts {
		var t Test
		t.Index = 0
		t.Status = 0
		t.Name = "config"
		t.TestType = TestNormal
		t.Mode = RunNormal
		t.CommandToRun = ""
		t.FileToRun = "config.test"
		t.Args = nil
		t.Expect = "OK"
		t.Output = ""
		t.LandMark = false

		Infrastructure.Hosts[index].Tests = append([]Test{t}, Infrastructure.Hosts[index].Tests...)

		for it := range Infrastructure.Hosts[index].Tests {
			Infrastructure.Hosts[index].Tests[it].Index = it
		}

		if !dirExists("tests/" + Infrastructure.Name + "/tests/" + h.Name) {
			createDir("tests/" + Infrastructure.Name + "/tests/" + h.Name + "/done")
		}

	}

	if mode == 0 {
		// Normal mode, look for last test results and load

		if RunningInfrastructure != nil {
			for index, h := range Infrastructure.Hosts {
				for _, hr := range RunningInfrastructure.Hosts {
					if h.Name == hr.Name && h.Kvm == hr.Kvm {
						for it, t := range h.Tests {
							for _, tr := range hr.Tests {
								if t.Index == tr.Index {
									Infrastructure.Hosts[index].Tests[it].Status = tr.Status
									t.Status = tr.Status
									Infrastructure.Hosts[index].Tests[it].TestType = tr.TestType
									t.TestType = tr.TestType
									Infrastructure.Hosts[index].Tests[it].Output = tr.Output
									t.Output = tr.Output
									Infrastructure.Hosts[index].Tests[it].Mode = tr.Mode
									t.Mode = tr.Mode
								}
							}

							//Remove all snap files not needed
							if t.Status < 1 && t.TestType != 2 {
								if _, err := os.Stat(h.GetHostKvm().Templates + "/" + Infrastructure.Name + "/" + h.Name + "." + strconv.Itoa(t.Index) + "-" + t.Name); err == nil {
									err := os.Remove(h.GetHostKvm().Templates + "/" + Infrastructure.Name + "/" + h.Name + "." + strconv.Itoa(t.Index) + "-" + t.Name)
									if err != nil {
										log.Println(err)
									}
								}
							}

						}
					}
				}
			}
		}

	} else if mode == 1 {
		// Reset mode, restart all tests
		// Only load Config test from RunningInfrastructure Tests.
		for index, h := range Infrastructure.Hosts {
			for _, hr := range RunningInfrastructure.Hosts {
				if h.Name == hr.Name && h.Kvm == hr.Kvm {
					Infrastructure.Hosts[index].Tests[0].Status = RunningInfrastructure.Hosts[index].Tests[0].Status
					Infrastructure.Hosts[index].Tests[0].TestType = RunningInfrastructure.Hosts[index].Tests[0].TestType
					Infrastructure.Hosts[index].Tests[0].Output = RunningInfrastructure.Hosts[index].Tests[0].Output
					Infrastructure.Hosts[index].Tests[0].Mode = RunningInfrastructure.Hosts[index].Tests[0].Mode
				}
			}
			//Remove all snap files not needed
			for _, t := range h.Tests {
				if _, err := os.Stat(h.GetHostKvm().Templates + "/" + Infrastructure.Name + "/" + h.Name + "." + strconv.Itoa(t.Index) + "-" + t.Name); err == nil {
					err := os.Remove(h.GetHostKvm().Templates + "/" + Infrastructure.Name + "/" + h.Name + "." + strconv.Itoa(t.Index) + "-" + t.Name)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	} else if mode == 2 {
		// Look for first landmark and start from there.
		if RunningInfrastructure != nil {
			for index, h := range Infrastructure.Hosts {
				landmarkFound := false
				for _, hr := range RunningInfrastructure.Hosts {
					if h.Name == hr.Name && h.Kvm == hr.Kvm {
						for it, t := range h.Tests {
							if landmarkFound == false {
								for _, tr := range hr.Tests {
									if t.Index == tr.Index {
										Infrastructure.Hosts[index].Tests[it].Status = tr.Status
										Infrastructure.Hosts[index].Tests[it].TestType = tr.TestType
										Infrastructure.Hosts[index].Tests[it].Output = tr.Output
										Infrastructure.Hosts[index].Tests[it].Mode = tr.Mode
									}
								}
								//set to true landmarkFound
								landmarkFound = t.LandMark
							} else {
								//Remove all snap files not needed.
								if _, err := os.Stat(h.GetHostKvm().Templates + "/" + Infrastructure.Name + "/" + h.Name + "." + strconv.Itoa(t.Index) + "-" + t.Name); err == nil {
									err := os.Remove(h.GetHostKvm().Templates + "/" + Infrastructure.Name + "/" + h.Name + "." + strconv.Itoa(t.Index) + "-" + t.Name)
									if err != nil {
										log.Println(err)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func cleanTests() {
	if Infrastructure.Name != "" && dirExists("tests/"+Infrastructure.Name+"/tests") {
		RemoveDirEsp("tests/" + Infrastructure.Name + "/tests")
		log.Println("+ Tests dir tests/" + Infrastructure.Name + "/tests deleted")
	} else {
		log.Println("- Tests dir tests/" + Infrastructure.Name + "/tests no exists")
	}
}

func (infra *DRLMTestingConfig) CountConfigRun() int {
	hostsDone := 0
	for _, h := range Infrastructure.Hosts {
		hostname, _ := os.Hostname()
		if h.Kvm == "localhost" || h.Kvm == hostname {
			for _, t := range h.Tests {
				if t.Name == "config" && t.Status != 0 {
					hostsDone++
				}
			}
		}
	}
	return hostsDone
}

func (infra *DRLMTestingConfig) CountConfigOk() int {
	hostsDone := 0
	for _, h := range Infrastructure.Hosts {
		hostname, _ := os.Hostname()
		if h.Kvm == "localhost" || h.Kvm == hostname {
			for _, t := range h.Tests {
				if t.Name == "config" && t.Status == 1 {
					hostsDone++
				}
			}
		}
	}
	return hostsDone
}

func (infra *DRLMTestingConfig) AllDone() bool {
	done := true

	for _, h := range Infrastructure.Hosts {
		hostname, _ := os.Hostname()
		if h.Kvm == "localhost" || h.Kvm == hostname {
			for _, t := range h.Tests {
				if t.Status == 0 {
					done = false
					return done
				}
			}
		}
	}
	return done
}

func (t *Test) initTest() {
	t.Index = 0
	t.Status = 0
	t.Name = ""
	t.TestType = TestNormal
	t.Mode = 0
	t.CommandToRun = ""
	t.FileToRun = ""
	t.Args = nil
	t.Expect = ""
	t.Output = ""
	t.LandMark = false
}

func (h *Host) GetNextTest() Test {

	test := &Test{}
	test.initTest()

	for _, t := range h.Tests {
		if t.Status == 0 || (t.Status == -1 && t.TestType != 2) {
			return t
		}
	}

	return *test
}

func (t *Test) PublishTest(host Host) {

	// Check for dependencies
	done := true

	if len(t.Dependencies) > 0 {
		for _, td := range t.Dependencies {
			for _, rh := range RunningInfrastructure.Hosts {
				if Infrastructure.Prefix+"-"+Infrastructure.Name+"-"+td.Host == rh.Name {
					if rh.Tests[td.Test].Status < 1 {
						done = false
					}
				}
			}
		}
		if !done {
			log.Println("* " + host.Name + ": Test <" + strconv.Itoa(t.Index) + "-" + t.Name + "> wainting for dependencies")
			return
		}
	}

	// mirar is existeix l'snap i si existeix vol dir que el test ja està en run i marxem.
	// sino existeix crear-lo
	if _, err := os.Stat(host.GetHostKvm().Templates + "/" + Infrastructure.Name + "/" + host.Name + "." + strconv.Itoa(t.Index) + "-" + t.Name); err == nil {
		return
	}

	//Crea un Snap abans d'executar el test
	host.GetHostKvm().createSnap(host.Name, strconv.Itoa(t.Index)+"-"+t.Name)
	time.Sleep(2 * time.Second)

	//Mirem en quin mode hem d'executar el test (rescue o normal)
	host.GetHostKvm().checkMode(host.Name, t.Mode)

	if t.FileToRun == "" && t.CommandToRun == "" {
		log.Println("- Test " + t.Name + " has nothing to run")
	} else if t.CommandToRun != "" {
		newLocation := "tests/" + Infrastructure.Name + "/tests/" + host.Name + "/" + strconv.Itoa(t.Index) + "-" + t.Name + ".test"
		f, err := os.Create(newLocation)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = f.WriteString(t.CommandToRun)
		if err != nil {
			log.Println(err)
			f.Close()
			return
		}
		log.Println("+ " + host.Name + ": Test <" + strconv.Itoa(t.Index) + "-" + t.Name + "> created successfully")
		err = f.Close()
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		srcLocation := "tests/" + Infrastructure.Name + "/" + t.FileToRun
		if !fileExists(srcLocation) {
			srcLocation = Infrastructure.Templates + "/templates/" + host.Template + "/tests/" + t.FileToRun
			if !fileExists(srcLocation) {
				log.Println("- Test " + t.FileToRun + "not found")
				return
			}
		}

		log.Println("+ " + host.Name + ": Test <" + strconv.Itoa(t.Index) + "-" + t.Name + "> created successfully")

		dstLocation := "tests/" + Infrastructure.Name + "/tests/" + host.Name + "/" + strconv.Itoa(t.Index) + "-" + t.FileToRun

		from, err := os.Open(srcLocation)
		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()

		to, err := os.OpenFile(dstLocation, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			log.Fatal(err)
		}

		if len(t.Args) > 0 {
			for i, arg := range t.Args {
				input, err := ioutil.ReadFile(dstLocation)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

				output := bytes.Replace(input, []byte("{{var"+strconv.Itoa(i)+"}}"), []byte(arg), -1)

				if err = ioutil.WriteFile(dstLocation, output, 0666); err != nil {
					log.Println(err)
					os.Exit(1)
				}
			}
		}

	}
}
