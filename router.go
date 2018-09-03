package main

import (
	"log"
	"net/http"

	"github.com/sedind/flow/app/middleware/cors"

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

	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   ctx.CORS.AllowedOrigins,
		AllowedMethods:   ctx.CORS.AllowedMethods,
		AllowedHeaders:   ctx.CORS.AllowedHeaders,
		ExposedHeaders:   ctx.CORS.ExposedHeaders,
		AllowCredentials: ctx.CORS.AllowCredentials,
		MaxAge:           ctx.CORS.MaxAge,
	})
	r.Use(cors.Handler)

	auth := auth.New(ctx)

	//public application routes
	r.Group(func(r router.Router) {
		r.Mount("/auth", auth.Routes())
	})

	//protected application routes
	r.Group(func(r router.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(auth.JWTVerifierHandler())

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(auth.JWTAuthenticatorMiddleware)

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
