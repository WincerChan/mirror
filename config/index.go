package config

import (
	"io/ioutil"
	T "mirror/tool"

	"gopkg.in/yaml.v2"
)

type Replaced struct {
	Old string `yaml:"old"`
	New string `yaml:"new"`
}

type Yaml struct {
	Host struct {
		Self  string `yaml:"self"`
		Proxy string `yaml:"proxy"`
	}
	ReplacedURLs []Replaced `yaml:"replaced_urls"`
	EnableSSL    bool       `yaml:"enable_ssl"`
	HandleCookie bool       `yaml:"handle_cookie"`
}

var Config *Yaml
var Protocal string

func LoadConfig() {
	Config = new(Yaml)
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	T.CheckErr(err)
	err = yaml.Unmarshal(yamlFile, Config)
	T.CheckErr(err)
	if Config.EnableSSL {
		Protocal = "https://"
	} else {
		Protocal = "http://"
	}
}
