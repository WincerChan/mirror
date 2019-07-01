package main

import "net/http"

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

func configStruct(w http.ResponseWriter, r *http.Request) {

}
