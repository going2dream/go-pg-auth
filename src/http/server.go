package http

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/http/router"
	SQLStore "github.com/ZeroDayDrake/go-pg-auth/src/http/store/sql"
	"github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Server struct {
	config *AppConfig
	Logger *zap.Logger
	db     *pgxpool.Pool
}

func NewHttpServer() Server {
	return Server{
		config: NewAppConfig(),
		Logger: logger.New(),
		db:     NewDBConnection(),
	}
}

func (s *Server) Start() {
	store := SQLStore.New(s.db)
	h := router.BuildRouter(store).Handler
	//if false {
	//	h = fasthttp.CompressHandler(h)
	//}

	s.Logger.Info("Binding to TCP address", zap.String("IP", s.config.BindIP), zap.String("Port", s.config.BindPort))

	if err := fasthttp.ListenAndServe(s.config.BindIP+":"+s.config.BindPort, h); err != nil {
		s.Logger.Fatal("Error in ListenAndServe", zap.Error(err))
	}
}
