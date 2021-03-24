package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type (
	Config struct {
		Environment    string   `yaml:"environment"`
		LogLevel       string   `yaml:"log_level"`
		LogOutputPaths []string `yaml:"log_outputPaths"`
	}
)

func New() *zap.Logger {
	yml, err := ioutil.ReadFile("config/app.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(yml, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	if config.Environment == "dev" {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		return logger
	}

	productionConfig := zap.NewProductionConfig()
	var l zapcore.Level
	if err := l.UnmarshalText([]byte(config.LogLevel)); err != nil {
		log.Fatalf("error: %v", err)
	}
	productionConfig.Level = zap.NewAtomicLevelAt(l)
	productionConfig.OutputPaths = config.LogOutputPaths

	logger, err := productionConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
