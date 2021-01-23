package main

import "github.com/ZeroDayDrake/go-pg-auth/src/http"

func main() {
	httpServer := http.NewHttpServer()
	defer httpServer.Logger.Sync()
	httpServer.Start()
}
