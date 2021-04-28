package config

import (
	"sync"

	"gopkg.in/yaml.v2"
	"io/ioutil"
	T "mirror/tool"
)

var (
	once   sync.Once
	config *Config
)

func GetConfig() *Config {
	if config != nil {
		return config
	}
	once.Do(loadConfig)
	return config
}

func loadConfig() {
	config = new(Config)
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	T.CheckErr(err)
	err = yaml.Unmarshal(yamlFile, config)
	T.CheckErr(err)
	if config.EnableSSL {
		config.Protocol = "https://"
	} else {
		config.Protocol = "http://"
	}
}
