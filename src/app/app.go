package app

import (
	"github.com/going2dream/go-pg-auth/src/app/logger"
	"github.com/going2dream/go-pg-auth/src/app/store"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var log = logger.New()

type app struct {
	config *Config
	Store  store.Store
}

func New() *app {
	return &app{
		config: NewAppConfig(),
	}
}

func (a *app) SetStore(store store.Store) {
	a.Store = store
}

func (a *app) Start() {
	h := BuildRouter(a.Store).Handler

	if a.config.Environment == "prod" {
		h = fasthttp.CompressHandler(h)
	}

	log.Info("Binding to TCP address", zap.String("IP", a.config.BindIP), zap.String("Port", a.config.BindPort))

	if err := fasthttp.ListenAndServe(a.config.BindIP+":"+a.config.BindPort, h); err != nil {
		log.Fatal("Error in ListenAndServe", zap.Error(err))
	}
}
