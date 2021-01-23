package router

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/http/controllers"
	"github.com/fasthttp/router"
)

func BuildRouter() *router.Router {
	r := router.New()

	r.POST("/login", AuthController.Login)
	r.POST("/refresh.token", AuthController.RefreshToken)
	r.POST("/logout", AuthController.Logout)

	return r
}
