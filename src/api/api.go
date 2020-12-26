package api

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/api/router"
	"github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type HttpServer struct {
	config *Config
	Logger *zap.Logger
}

func (s *HttpServer) Start() {
	h := router.BuildRouter().Handler
	if false {
		h = fasthttp.CompressHandler(h)
	}

	s.Logger.Info("Binding to TCP address", zap.String("IP", s.config.IP), zap.String("Port", s.config.Port))

	if err := fasthttp.ListenAndServe(s.config.IP+":"+s.config.Port, h); err != nil {
		s.Logger.Fatal("Error in ListenAndServe", zap.Error(err))
	}
}

func NewHttpServer() HttpServer {
	configFile, err := ioutil.ReadFile("config/http.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	return HttpServer{
		config: &config,
		Logger: logger.New(),
	}
}
