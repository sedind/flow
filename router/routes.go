package router

// Routes interface adds two methods for router traversal
type Routes interface {
	// Routes returns routing tree in an easily traversable structure
	Routes() []Route

	// Middlewares returns list of middlewares used my Router
	Middlewares() Middlewares

	// Match searches routing tree for handler that matches
	// the method/path
	Match(rctx *Context, method, path string) bool
}
