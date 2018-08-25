package router

import "net/http"

// Route describes the details of a routing handler
type Route struct {
	Pattern   string
	Handlers  map[string]http.Handler
	SubRoutes Routes
}
