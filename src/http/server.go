package http

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/http/router"
	"github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Server struct {
	config *Config
	Logger *zap.Logger
	db     *pgxpool.Pool
}

func (s *Server) Start() {
	h := router.BuildRouter().Handler
	if false {
		h = fasthttp.CompressHandler(h)
	}

	s.Logger.Info("Binding to TCP address", zap.String("IP", s.config.BindIP), zap.String("Port", s.config.BindPort))

	if err := fasthttp.ListenAndServe(s.config.BindIP+":"+s.config.BindPort, h); err != nil {
		s.Logger.Fatal("Error in ListenAndServe", zap.Error(err))
	}
}

func NewHttpServer() Server {
	configFile, err := ioutil.ReadFile("config/app.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	return Server{
		config: &config,
		Logger: logger.New(),
		db:     NewDBConnection(),
	}
}
