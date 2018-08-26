package main

import (
	"log"
	"net/http"

	"github.com/sedind/flow/app"
	"github.com/sedind/flow/app/middleware"
	"github.com/sedind/flow/app/router"
	"github.com/sedind/flow/features/auth"
)

func newAppRouter(ctx *app.Context) *router.Mux {
	appRouter := router.New()

	if ctx.PanicRecover {
		// Recover from panics without crashing server
		appRouter.Use(middleware.Recoverer)
	}

	if ctx.RedirectSlashes {
		//Redirect slashes to no slash URL version
		appRouter.Use(middleware.RedirectSlashes)
	}

	if ctx.RequestLogging {
		// log API request calls
		appRouter.Use(middleware.Logger)
	}

	if ctx.CompressResponse {
		// Compress response body
		appRouter.Use(middleware.DefaultCompress)
	}

	// mount application routes
	appRouter.Route("/v1", func(r router.Router) {
		r.Mount("/auth", auth.New(ctx).Routes())
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Walk and print all routes
		log.Printf("%s %s \n", method, route)
		return nil
	}

	if err := router.Walk(appRouter, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	return appRouter
}
