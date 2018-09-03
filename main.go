package main

import (
	"net/http"

	"github.com/sedind/flow/app"
)

func main() {
	a := app.New("config.yml")

	a.RegisterRouter(func(ctx *app.Context) http.Handler {
		return newAppRouter(ctx)
	})

	// serve the app
	if err := a.Serve(); err != nil {
		a.Context.Logger.Fatal(err)
	}

}
