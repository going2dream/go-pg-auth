package app

import (
	"github.com/fasthttp/router"
	"github.com/going2dream/go-pg-auth/src/app/controllers"
	"github.com/going2dream/go-pg-auth/src/app/store"
)

func BuildRouter(store store.Store) *router.Router {
	r := router.New()

	var AuthController = &controllers.Auth{
		Store: store,
	}

	r.POST("/login", AuthController.Login)
	r.POST("/refresh.token", AuthController.RefreshToken)
	r.POST("/logout", AuthController.Logout)

	return r
}
