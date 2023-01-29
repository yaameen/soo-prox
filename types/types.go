package types

type Config struct {
	Secure     bool    `yaml:"secure"`
	ConfigFile string  `yaml:"configFile"`
	Name       string  `yaml:"name"`
	Port       int     `yaml:"port"`
	Host       string  `yaml:"host"`
	Proxies    []Proxy `yaml:"proxies"`
}

type Proxy struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Prefix string `yaml:"prefix"`
}
