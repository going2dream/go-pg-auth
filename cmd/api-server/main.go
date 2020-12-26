package main

import "github.com/ZeroDayDrake/go-pg-auth/src/api"

func main() {
	httpServer := api.NewHttpServer()
	defer httpServer.Logger.Sync()
	httpServer.Start()
}
