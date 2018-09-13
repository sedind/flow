package dbe

import (
	"fmt"
	"strings"
)

// HavingClause defines a condition and its arguments for a HAVING clause
type HavingClause struct {
	Condition string
	Arguments []interface{}
}

// HavingClauses represents collection of HavingClause
type HavingClauses []HavingClause

func (c HavingClause) String() string {
	sql := fmt.Sprintf("%s", c.Condition)

	return sql
}

func (c HavingClauses) String() string {
	if len(c) == 0 {
		return ""
	}

	cs := []string{}
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, " AND ")
}
