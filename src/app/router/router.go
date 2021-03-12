package router

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/app/controllers"
	"github.com/ZeroDayDrake/go-pg-auth/src/app/http"
	"github.com/fasthttp/router"
)

func BuildRouter(s *http.Server) *router.Router {
	r := router.New()

	var AuthController = &controllers.Auth{
		Store: s.Store,
	}

	r.POST("/login", AuthController.Login)
	r.POST("/refresh.token", AuthController.RefreshToken)
	r.POST("/logout", AuthController.Logout)

	return r
}
