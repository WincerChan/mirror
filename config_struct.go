package main

type Replaced struct {
	Old string `yaml:"old"`
	New string `yaml:"new"`
}

type Yaml struct {
	Host struct {
		Self  string `yaml:"self"`
		Proxy string `yaml:"proxy"`
	}
	ReplacedStrings []Replaced `yaml:"replaced_strings"`
	EnableSSL bool `yaml:"enable_ssl"`
	HandleCookie bool `yaml:"handle_cookie"`
}
