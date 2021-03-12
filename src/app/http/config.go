package http

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Environment string `yaml:"environment"`
	BindIP      string `yaml:"bind_ip"`
	BindPort    string `yaml:"bind_port"`
}

func NewAppConfig() *Config {
	configFile, err := ioutil.ReadFile("config/app.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	return &config
}
