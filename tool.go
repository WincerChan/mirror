package main

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func hasGziped(coding string) bool {
	return strings.HasPrefix(coding, "gz")
}

func isTextType(typeName string) bool {
	return strings.HasPrefix(typeName, "text") ||
		strings.HasPrefix(typeName, "appli")
}

func loadConfig() {
	Config = new(Yaml)
	yamlFile, err := ioutil.ReadFile("config.yaml")
	checkErr(err)
	err = yaml.Unmarshal(yamlFile, Config)
	checkErr(err)
	if Config.EnableSSL {
		protocal = "https://"
	} else {
		protocal = "http://"
	}
	log.Println("conf", Config.ReplacedURLs)
}
