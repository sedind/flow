package dbe

import "strings"

// Clause represents Query clause
type Clause struct {
	Fragment  string
	Arguments []interface{}
}

// Clauses represents a Clause collection
type Clauses []Clause

// Join forms string of clauses separatet by given separator
func (c Clauses) Join(sep string) string {
	out := make([]string, 0, len(c))
	for _, clause := range c {
		out = append(out, clause.Fragment)
	}
	return strings.Join(out, sep)
}

// Args gets array of Clause arguments
func (c Clauses) Args() (args []interface{}) {
	for _, clause := range c {
		for _, arg := range clause.Arguments {
			args = append(args, arg)
		}
	}
	return
}
