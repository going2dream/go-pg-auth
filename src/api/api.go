package api

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/api/router"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type HttpServer struct {
	config *Config
}

func (s *HttpServer) Start() {
	h := router.BuildRouter().Handler
	if false {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(s.config.IP+":"+s.config.Port, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func NewHttpServer() HttpServer {
	yamlConfig, err := ioutil.ReadFile("config/http.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	return HttpServer{
		config: &config,
	}
}
