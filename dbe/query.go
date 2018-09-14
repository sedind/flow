package dbe

import (
	"fmt"
	"regexp"
	"strings"
)

var inRegex = regexp.MustCompile(`(?i)in\s*\(\s*\?\s*\)`)

// Query is the main value that is used to build up a query
// to be executed against the `Connection`.
type Query struct {
	RawSQL        *Clause
	limitResults  int
	eager         bool
	eagerFields   []string
	whereClauses  Clauses
	orderClauses  Clauses
	fromClauses   FromClauses
	joinClauses   JoinClauses
	groupClauses  GroupClauses
	havingClauses HavingClauses
	Paginator     *Paginator
	Connection    *Connection
}

// Query Creates new Empty Query
func (c *Connection) Query() *Query {
	return NewQuery(c)
}

// NewQuery creates new "empty" query from the current connection.
func NewQuery(c *Connection) *Query {
	return &Query{
		RawSQL:     &Clause{},
		Connection: c,
	}
}

// Clone colnes current query
func (q *Query) Clone(targetQ *Query) {
	rawSQL := *q.RawSQL
	targetQ.RawSQL = &rawSQL

	targetQ.limitResults = q.limitResults
	targetQ.whereClauses = q.whereClauses
	targetQ.orderClauses = q.orderClauses
	targetQ.fromClauses = q.fromClauses
	targetQ.joinClauses = q.joinClauses
	targetQ.groupClauses = q.groupClauses
	targetQ.havingClauses = q.havingClauses

	if q.Paginator != nil {
		paginator := *q.Paginator
		targetQ.Paginator = &paginator
	}

	if q.Connection != nil {
		connection := *q.Connection
		targetQ.Connection = &connection
	}
}

// Raw will override the query building feature, and will use
// whatever query you want to execute against the `Connection`. You can continue
// to use the `?` argument syntax.
//
//	q.RawQuery("select * from foo where id = ?", 1)
func (q *Query) Raw(stmt string, args ...interface{}) *Query {
	q.RawSQL = &Clause{stmt, args}
	return q
}

// Where will append a where clause to the query. You may use `?` in place of
// arguments.
//
// 	q.Where("id = ?", 1)
// 	q.Where("id in (?)", 1, 2, 3)
func (q *Query) Where(stmt string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	if inRegex.MatchString(stmt) {
		var inq []string
		for i := 0; i < len(args); i++ {
			inq = append(inq, "?")
		}
		qs := fmt.Sprintf("(%s)", strings.Join(inq, ","))
		stmt = strings.Replace(stmt, "(?)", qs, 1)
	}
	q.whereClauses = append(q.whereClauses, Clause{stmt, args})
	return q
}

// Order will append an order clause to the query.
//
// 	q.Order("name desc")
func (q *Query) Order(stmt string) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.orderClauses = append(q.orderClauses, Clause{stmt, []interface{}{}})
	return q
}

// Limit will add a limit clause to the query.
func (q *Query) Limit(limit int) *Query {
	q.limitResults = limit
	return q
}

// ToSQL will generate SQL and the appropriate arguments for that SQL
// from the `Model` passed in.
func (q Query) ToSQL(model *Model, addColumns ...string) (string, []interface{}) {
	sb := q.toSQLBuilder(model, addColumns...)
	return sb.String(), sb.Args()
}

// ToSQLBuilder returns a new `SQLBuilder` that can be used to generate SQL,
// get arguments, and more.
func (q Query) toSQLBuilder(model *Model, addColumns ...string) *QueryBuilder {
	return NewQueryBuilder(q, model, addColumns...)
}
