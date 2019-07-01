package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
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
