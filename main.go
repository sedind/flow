package main

import (
	"log"
	"net/http"

	"github.com/sedind/flow/app"
)

func main() {
	a := app.New("config.yaml")

	a.RegisterRouter(func(ctx *app.Context) http.Handler {
		return newAppRouter(ctx)
	})
	// serve the app
	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}
}
