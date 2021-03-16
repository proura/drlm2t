package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/proura/drlm2t/model"
	"gopkg.in/yaml.v2"
)

func apiGetInfrastructures(w http.ResponseWriter, r *http.Request) {

	var infrastructures []model.DRLMTestingConfig

	files, err := ioutil.ReadDir(configDRLM2T.TestPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filesinf, err := ioutil.ReadDir(configDRLM2T.TestPath + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		found := false
		for _, finf := range filesinf {
			if finf.Name() == "infrastructure.yaml" {
				found = true
				break
			}
		}
		if found {
			infra := getInfrastructure(configDRLM2T.TestPath + "/" + f.Name())
			infra.Name = f.Name()
			infrastructures = append(infrastructures, *infra)
		}
	}

	response := generateJSONResponse(infrastructures)
	fmt.Fprintln(w, response)
}

func apiGetInfrastructure(w http.ResponseWriter, r *http.Request) {

	receivedInfrastructure := getField(r, 0)

	content, err := ioutil.ReadFile(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/infrastructure.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, string(content))
}

func apiSetInfrastructure(w http.ResponseWriter, r *http.Request) {
	receivedInfrastructure := getField(r, 0)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
	}
	bodyString := string(bodyBytes)

	f, err := os.Create(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/infrastructure.yaml")
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()

	_, err = f.WriteString(bodyString)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("updated " + configDRLM2T.TestPath + "/" + receivedInfrastructure + "/infrastructure.yaml")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func apiPutInfrastructure(w http.ResponseWriter, r *http.Request) {
	receivedInfrastructure := getField(r, 0)

	path := configDRLM2T.TestPath + "/" + receivedInfrastructure

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0775)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad")
		return
	}

	bodyString := "description: \"Template infrastructure\"" + "\n" +
		"nets:" + "\n" +
		"- name: net1" + "\n" +
		"  ip: 192.168.181.1" + "\n" +
		"" + "\n" +
		"hosts:" + "\n" +
		"#################" + "\n" +
		"## DRLM SERVER ##" + "\n" +
		"#################" + "\n" +
		"- name: srv1" + "\n" +
		"  template: deb10" + "\n" +
		"  nets:" + "\n" +
		"  - name: net1" + "\n" +
		"    ip: 192.168.181.2" + "\n" +
		"  tests:" + "\n" +
		"" + "\n" +
		"  # host: srv1 - test: 1" + "\n" +
		"  # Update DRLM server making apt update & apt upgrade" + "\n" +
		"  - name: update" + "\n" +
		"    testtype: 1" + "\n" +
		"    filetorun: update.test" + "\n" +
		"    expect: \"OK\"" + "\n" +
		"" + "\n" +
		"#################" + "\n" +
		"## DRLM CLIENT ##" + "\n" +
		"#################" + "\n" +
		"- name: cli1" + "\n" +
		"  template: deb10" + "\n" +
		"  nets:" + "\n" +
		"  - name: net1" + "\n" +
		"    ip: 192.168.181.53" + "\n" +
		"  tests:" + "\n" +
		"" + "\n" +
		"  # host: deb10 - test: 1" + "\n" +
		"  # Update client making apt update & apt upgrade" + "\n" +
		"  - name: update" + "\n" +
		"    filetorun: update.test" + "\n" +
		"    expect: \"OK\"" + "\n"

	f, err := os.Create(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/infrastructure.yaml")
	if err != nil {
		log.Println("===>" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad")
		return
	}
	defer f.Close()

	_, err2 := f.WriteString(bodyString)
	if err2 != nil {
		log.Println(err2.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad")
		return
	}

	log.Println("created " + configDRLM2T.TestPath + "/" + receivedInfrastructure)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func apiDeleteInfrastructure(w http.ResponseWriter, r *http.Request) {

	receivedInfrastructure := getField(r, 0)

	err := os.RemoveAll(configDRLM2T.TestPath + "/" + receivedInfrastructure)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad")
		return
	}

	log.Println("deleted " + configDRLM2T.TestPath + "/" + receivedInfrastructure + "/infrastructure.yaml")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func apiGetRunning(w http.ResponseWriter, r *http.Request) {

	var runningInfrastructures []model.DRLMTestingConfig

	files, err := ioutil.ReadDir(configDRLM2T.TestPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filesinf, err := ioutil.ReadDir(configDRLM2T.TestPath + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		found := false
		for _, finf := range filesinf {
			if finf.Name() == "running.yaml" {
				found = true
				break
			}
		}
		if found {
			infra := getRunningInfrastructure(configDRLM2T.TestPath + "/" + f.Name())
			infra.Name = f.Name()
			runningInfrastructures = append(runningInfrastructures, *infra)
		}
	}

	response := generateJSONResponse(runningInfrastructures)
	fmt.Fprintln(w, response)
}

func getInfrastructure(file string) *model.DRLMTestingConfig {

	c := new(model.DRLMTestingConfig)

	yamlFile, err := ioutil.ReadFile(file + "/infrastructure.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Println("Unmarshal: %v", err)
	}

	return c
}

func getRunningInfrastructure(file string) *model.DRLMTestingConfig {

	c := new(model.DRLMTestingConfig)

	yamlFile, err := ioutil.ReadFile(file + "/running.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Println("Unmarshal: %v", err)
	}

	return c
}

func apiGetTemplates(w http.ResponseWriter, r *http.Request) {
	var templates []model.Template
	files, err := ioutil.ReadDir(configDRLM2T.TemplatesPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			infra := new(model.Template)
			infra.Name = f.Name()

			filesInDir, err := ioutil.ReadDir(configDRLM2T.TemplatesPath + "/" + f.Name() + "/tests")
			if err != nil {
				log.Fatal(err)
			}
			for _, fd := range filesInDir {
				templateTest := new(model.TemplateTest)
				templateTest.Name = fd.Name()
				content, err := ioutil.ReadFile(configDRLM2T.TemplatesPath + "/" + f.Name() + "/tests/" + fd.Name())
				if err != nil {
					log.Println(err)
					templateTest.Content = ""
				} else {
					templateTest.Content = string(content)
				}
				infra.TemplateTests = append(infra.TemplateTests, *templateTest)
			}
			templates = append(templates, *infra)
		}
	}
	response := generateJSONResponse(templates)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, response)
}

func apiGetTestResult(w http.ResponseWriter, r *http.Request) {

	receivedInfrastructure := getField(r, 0)
	receivedHostIndex, _ := strconv.Atoi(getField(r, 1))
	receivedTestIndex, _ := strconv.Atoi(getField(r, 2))

	infrastructure := getRunningInfrastructure(configDRLM2T.TestPath + "/" + receivedInfrastructure)

	if len(infrastructure.Hosts) > 1 {
		testFilePath := ""

		if _, err := os.Stat(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-OK-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].Name + ".test"); err == nil {
			testFilePath = configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-OK-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].Name + ".test"
		} else if _, err := os.Stat(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-OK-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].FileToRun); err == nil {
			testFilePath = configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-OK-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].FileToRun
		} else if _, err := os.Stat(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-ERR-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].Name + ".test"); err == nil {
			testFilePath = configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-ERR-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].Name + ".test"
		} else if _, err := os.Stat(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-ERR-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].FileToRun); err == nil {
			testFilePath = configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-ERR-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].FileToRun
		} else if _, err := os.Stat(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/" + getField(r, 2) + "-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].Name + ".test"); err == nil {
			testFilePath = configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/" + getField(r, 2) + "-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].Name + ".test"
		} else if _, err := os.Stat(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/" + getField(r, 2) + "-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].FileToRun); err == nil {
			testFilePath = configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/" + getField(r, 2) + "-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].FileToRun
		} else {
			log.Println(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/tests/" + infrastructure.Hosts[receivedHostIndex].Name + "/done/" + getField(r, 2) + "-OK-" + infrastructure.Hosts[receivedHostIndex].Tests[receivedTestIndex].Name + ".test")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "ok")
			return
		}

		testFileContent, err := ioutil.ReadFile(testFilePath)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(testFileContent))
	}
}

func apiUpTest(w http.ResponseWriter, r *http.Request) {
	receivedTestName := getField(r, 0)
	cmd := exec.Command("bash", "-c", "drlm2t up "+receivedTestName)
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func apiDownTest(w http.ResponseWriter, r *http.Request) {
	receivedTestName := getField(r, 0)
	cmd := exec.Command("bash", "-c", "drlm2t down "+receivedTestName)
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func apiRunTest(w http.ResponseWriter, r *http.Request) {
	receivedTestName := getField(r, 0)
	cmd := exec.Command("bash", "-c", "drlm2t run "+receivedTestName)
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func apiCleanTest(w http.ResponseWriter, r *http.Request) {
	receivedTestName := getField(r, 0)
	cmd := exec.Command("bash", "-c", "drlm2t clean "+receivedTestName)
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}
