//configuration.go
package main

import (
	"fmt"
)

type Configuration struct {
	Path          string
	WwwPath       string
	CertPath      string
	TestPath      string
	TemplatesPath string
	APIUser       string
	APIPasswd     string
	Certificate   string
	Key           string
}

var configDRLM2T Configuration

func printDRLMConfiguration() {
	fmt.Println("================================")
	fmt.Println("=== DRLM2T API CONFIGURATION ===")
	fmt.Println("================================")
	fmt.Println("DRLM2T_PATH = " + configDRLM2T.Path)
	fmt.Println("DRLM2T_WWW_PATH = " + configDRLM2T.Path)
	fmt.Println("DRLM2T_CERT_PATH = " + configDRLM2T.CertPath)
	fmt.Println("DRLM2T_TEST_PATH = " + configDRLM2T.CertPath)
	fmt.Println("DRLM2T_CERT = " + configDRLM2T.Certificate)
	fmt.Println("DRLM2T_KEY = " + configDRLM2T.Key)
	fmt.Println("DRLM2T_USER = " + configDRLM2T.APIUser)
	fmt.Println("DRLM2T_PASS = " + configDRLM2T.APIPasswd)
	fmt.Println("")
}

func initDRLM2TConfiguration() {
	configDRLM2T.Path = "/home/pau/src/proura/drlm2t"
	configDRLM2T.WwwPath = configDRLM2T.Path + "/drlm2t-front/www"
	configDRLM2T.CertPath = configDRLM2T.Path + "/cfg"
	configDRLM2T.TestPath = configDRLM2T.Path + "/tests"
	configDRLM2T.TemplatesPath = configDRLM2T.Path + "/drlm2t-templates/templates"
	configDRLM2T.APIUser = "admindrlm2t"
	configDRLM2T.APIPasswd = "admindrlm2t"
	configDRLM2T.Certificate = configDRLM2T.CertPath + "/drlm2t.crt"
	configDRLM2T.Key = configDRLM2T.CertPath + "/drlm2t.key"
}
