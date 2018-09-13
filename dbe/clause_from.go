package dbe

import (
	"fmt"
	"strings"
)

// FromClause represets SQL SELECT clause
type FromClause struct {
	From string
	As   string
}

// FromClauses represents collection of FromClause
type FromClauses []FromClause

func (c FromClause) String() string {
	return fmt.Sprintf("%s AS %s", c.From, c.As)
}

func (c FromClauses) String() string {
	cs := []string{}
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, ", ")
}
