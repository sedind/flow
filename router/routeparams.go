package router

// RouteParams is a structure to track URL routing parameters
type RouteParams struct {
	Keys, Values []string
}

// Add appends URL parameter to the end of the route params
func (s *RouteParams) Add(key, value string) {
	(*s).Keys = append((*s).Keys, key)
	(*s).Values = append((*s).Values, value)
}
