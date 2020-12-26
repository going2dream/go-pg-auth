package logger

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func New() *zap.Logger {
	configFile, err := ioutil.ReadFile("config/logger.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config zap.Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
