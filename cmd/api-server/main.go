package main

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/app/http"
)

func main() {
	httpServer := http.New()
	defer httpServer.Logger.Sync()
	httpServer.Start()
}
