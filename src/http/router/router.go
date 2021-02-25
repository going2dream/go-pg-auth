package router

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/http/controllers"
	"github.com/ZeroDayDrake/go-pg-auth/src/http/store"
	"github.com/fasthttp/router"
)

func BuildRouter(store *store.Store) *router.Router {
	r := router.New()

	var AuthController = &controllers.Auth{}

	r.POST("/login", AuthController.Login)
	r.POST("/refresh.token", AuthController.RefreshToken)
	r.POST("/logout", AuthController.Logout)

	return r
}
