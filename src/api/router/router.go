package router

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/api/controllers"
	"github.com/fasthttp/router"
)

func BuildRouter() *router.Router {
	r := router.New()

	r.GET("/", SiteController.Index)

	return r
}
