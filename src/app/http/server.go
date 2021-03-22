package http

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/app/router"
	"github.com/ZeroDayDrake/go-pg-auth/src/app/store"
	SQLStore "github.com/ZeroDayDrake/go-pg-auth/src/app/store/sql"
	"github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Server struct {
	config *Config
	Logger *zap.Logger
	Store  store.Store
}

func NewHTTPServer() Server {
	pool := store.NewPoolInstance()

	return Server{
		config: NewAppConfig(),
		Logger: logger.New(),
		Store:  SQLStore.New(pool),
	}
}

func (s *Server) Start() {
	h := router.BuildRouter(s).Handler

	if s.config.Environment == "prod" {
		h = fasthttp.CompressHandler(h)
	}

	s.Logger.Info("Binding to TCP address", zap.String("IP", s.config.BindIP), zap.String("Port", s.config.BindPort))

	if err := fasthttp.ListenAndServe(s.config.BindIP+":"+s.config.BindPort, h); err != nil {
		s.Logger.Fatal("Error in ListenAndServe", zap.Error(err))
	}
}
