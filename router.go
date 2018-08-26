package main

import (
	"log"
	"net/http"

	"github.com/sedind/flow/app"
	"github.com/sedind/flow/app/middleware"
	"github.com/sedind/flow/app/router"
	"github.com/sedind/flow/features/auth"
)

func newAppRouter(ctx *app.Context) http.Handler {
	r := router.New()

	if ctx.PanicRecover {
		// Recover from panics without crashing server
		r.Use(middleware.Recoverer)
	}

	if ctx.RedirectSlashes {
		//Redirect slashes to no slash URL version
		r.Use(middleware.RedirectSlashes)
	}

	if ctx.RequestLogging {
		// log API request calls
		r.Use(middleware.Logger)
	}

	if ctx.CompressResponse {
		// Compress response body
		r.Use(middleware.DefaultCompress)
	}
	if ctx.NoCache {
		// send no-cache headers
		r.Use(middleware.NoCache)
	}

	// mount application routes
	r.Route("/v1", func(r router.Router) {
		r.Mount("/auth", auth.New(ctx).Routes())
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Walk and print all routes
		log.Printf("%s %s \n", method, route)
		return nil
	}

	if err := router.Walk(r, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	return r
}
