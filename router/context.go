package router

import (
	"context"
	"net"
	"net/http"
	"strings"
)

var (
	// RouteCtxKey is the context.Context key to store the request context.
	RouteCtxKey = &contextKey{"RouteContext"}
)

// Context is the default routing context set on the root node of a
// request context to track route patterns, URL parameters and
// an optional routing path.
type Context struct {
	Routes Routes

	// Routing path/method override used during the route search.
	// See Mux#routeHTTP method.
	RoutePath   string
	RouteMethod string

	// Routing pattern stack throughout the lifecycle of the request,
	// across all connected routers. It is a record of all matching
	// patterns across a stack of sub-routers.
	RoutePatterns []string

	// URLParams are the stack of routeParams captured during the
	// routing lifecycle across a stack of sub-routers.
	URLParams RouteParams

	// The endpoint routing pattern that matched the request URI path
	// or `RoutePath` of the current sub-router. This value will update
	// during the lifecycle of a request passing through a stack of
	// sub-routers.
	routePattern string

	// Route parameters matched for the current sub-router. It is
	// intentionally unexported so it cant be tampered.
	routeParams RouteParams

	// methodNotAllowed hint
	methodNotAllowed bool
}

// NewRouteContext returns new Routing Context object
func NewRouteContext() *Context {
	return &Context{}
}

// Reset routing context to its initial state
func (c *Context) Reset() {
	c.Routes = nil
	c.RoutePath = ""
	c.RouteMethod = ""
	c.RoutePatterns = c.RoutePatterns[:0]
	c.URLParams.Keys = c.URLParams.Keys[:0]
	c.URLParams.Values = c.URLParams.Values[:0]

	c.routePattern = ""
	c.routeParams.Keys = c.routeParams.Keys[:0]
	c.routeParams.Values = c.routeParams.Values[:0]
	c.methodNotAllowed = false
}

// URLParam returns the corresponding URL parameter value from the request
// routing context.
func (c *Context) URLParam(key string) string {
	for k := len(c.URLParams.Keys) - 1; k >= 0; k-- {
		if c.URLParams.Keys[k] == key {
			return c.URLParams.Values[k]
		}
	}
	return ""
}

// RoutePattern builds the routing pattern string for the particular
// request, at the particular point during routing. This means, the value
// will change throughout the execution of a request in a router. That is
// why its advised to only use this value after calling the next handler.
func (c *Context) RoutePattern() string {
	routePattern := strings.Join(c.RoutePatterns, "")
	return strings.Replace(routePattern, "/*/", "/", -1)
}

// RouteContext returns routing Context object from a
// http.Request Context.
func RouteContext(ctx context.Context) *Context {
	return ctx.Value(RouteCtxKey).(*Context)
}

// URLParam returns the url parameter from a http.Request object.
func URLParam(r *http.Request, key string) string {
	if rctx := RouteContext(r.Context()); rctx != nil {
		return rctx.URLParam(key)
	}
	return ""
}

// URLParamFromCtx returns the url parameter from a http.Request Context.
func URLParamFromCtx(ctx context.Context, key string) string {
	if rctx := RouteContext(ctx); rctx != nil {
		return rctx.URLParam(key)
	}
	return ""
}

// ServerBaseContext wraps an http.Handler to set the request context to the
// `baseCtx`.
func ServerBaseContext(baseCtx context.Context, h http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		baseCtx := baseCtx

		// Copy over default net/http server context keys
		if v, ok := ctx.Value(http.ServerContextKey).(*http.Server); ok {
			baseCtx = context.WithValue(baseCtx, http.ServerContextKey, v)
		}
		if v, ok := ctx.Value(http.LocalAddrContextKey).(net.Addr); ok {
			baseCtx = context.WithValue(baseCtx, http.LocalAddrContextKey, v)
		}

		h.ServeHTTP(w, r.WithContext(baseCtx))
	})
	return fn
}

// contextKey is a value for use with context.WithValue.
// it is used as a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "context value " + k.name
}
