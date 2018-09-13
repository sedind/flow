package dbe

import (
	"fmt"
	"strings"
)

// JoinClause represets SQL JOIN clause
type JoinClause struct {
	JoinType  string
	Table     string
	On        string
	Arguments []interface{}
}

// JoinClauses represents collection of JoinClause
type JoinClauses []JoinClause

func (c JoinClause) String() string {
	sql := fmt.Sprintf("%s %s", c.JoinType, c.Table)

	if len(c.On) > 0 {
		sql += " ON " + c.On
	}

	return sql
}

func (c JoinClauses) String() string {
	cs := []string{}
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, " ")
}
