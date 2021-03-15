package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

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
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	f, err := os.Create(configDRLM2T.TestPath + "/" + receivedInfrastructure + "/infrastructure.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err2 := f.WriteString(bodyString)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("updated " + configDRLM2T.TestPath + "/" + receivedInfrastructure + "/infrastructure.yaml")

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
		log.Fatalf("Unmarshal: %v", err)
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
		log.Fatalf("Unmarshal: %v", err)
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
	fmt.Fprintln(w, response)
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
