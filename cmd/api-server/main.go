package main

import (
	"github.com/going2dream/go-pg-auth/src/app"
	"github.com/going2dream/go-pg-auth/src/app/store/pgsql"
)

func main() {
	a := app.New()
	a.SetStore(pgsql.NewStore())

	a.Start()
}
