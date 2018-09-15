package flow

// ContextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type ContextKey struct {
	name string
}

// NewContextKey Creates New ContextKey Object
// ContextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
func NewContextKey(name string) *ContextKey {
	return &ContextKey{
		name: name,
	}
}

func (ctx *ContextKey) String() string {
	return "flow context value: " + ctx.name
}
