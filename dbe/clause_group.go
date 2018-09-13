package dbe

import "strings"

// GroupClause represets SQL GROUP BY clause
type GroupClause struct {
	Field string
}

// GroupClauses represents collection of GroupClause
type GroupClauses []GroupClause

func (c GroupClause) String() string {
	return c.Field
}

func (c GroupClauses) String() string {
	cs := []string{}
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, ", ")
}
