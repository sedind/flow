package router

import "net/http"

// Middlewares type is a slice of standard middleware handlers
// with methods to compose middleware chains and http.Handler's
type Middlewares []func(http.Handler) http.Handler

// Handler builds and returns http.Handler from the chain
// of middlewares with `http.Handler` as the final handler
func (m Middlewares) Handler(h http.Handler) http.Handler {
	return &ChainHandler{m, h, chain(m, h)}
}

// HandlerFunc builds and returns a http.Handler from the chain
// of middlewares with `http.Handler` as the final handler
func (m Middlewares) HandlerFunc(h http.HandlerFunc) http.Handler {
	return &ChainHandler{m, h, chain(m, h)}
}
