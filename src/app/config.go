package app

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Environment string `yaml:"environment"`
	BindIP      string `yaml:"bind_ip"`
	BindPort    string `yaml:"bind_port"`
}

func NewAppConfig() *Config {
	configFile, err := ioutil.ReadFile("config/app.yml")
	if err != nil {
		log.Fatal("Cant read app config file", zap.Error(err))
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatal("Cant unmarshal app config file", zap.Error(err))
	}

	return &config
}
