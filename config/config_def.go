package config

type Replaced struct {
	Old string `yaml:"old"`
	New string `yaml:"new"`
}

type Config struct {
	Host struct {
		Self  string `yaml:"self"`
		Proxy string `yaml:"proxy"`
	}
	ReplacedURLs   []Replaced `yaml:"replaced_urls"`
	EnableSSL      bool       `yaml:"enable_ssl"`
	HandleCookie   bool       `yaml:"handle_cookie"`
	Protocol       string     `yaml:"-"`
	Token          string     `yaml:"token"`
	HeaderTokenKey string     `yaml:"header_token_key"`
}
