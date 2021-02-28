package http

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/http/controllers"
	"github.com/fasthttp/router"
)

func BuildRouter(s *Server) *router.Router {
	r := router.New()

	var AuthController = &controllers.Auth{
		Store: s.store,
	}

	r.POST("/login", AuthController.Login)
	r.POST("/refresh.token", AuthController.RefreshToken)
	r.POST("/logout", AuthController.Logout)

	return r
}
