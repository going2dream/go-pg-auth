package http

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/http/store"
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
	store  store.Store
}

func NewHttpServer() Server {
	db := NewDBConnection()
	return Server{
		config: NewAppConfig(),
		Logger: logger.New(),
		db:     db,
		store:  SQLStore.New(db),
	}
}

func (s *Server) Start() {
	h := BuildRouter(s).Handler
	//if false {
	//	h = fasthttp.CompressHandler(h)
	//}

	s.Logger.Info("Binding to TCP address", zap.String("IP", s.config.BindIP), zap.String("Port", s.config.BindPort))

	if err := fasthttp.ListenAndServe(s.config.BindIP+":"+s.config.BindPort, h); err != nil {
		s.Logger.Fatal("Error in ListenAndServe", zap.Error(err))
	}
}
