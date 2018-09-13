package dbe

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// QueryBuilder builds SQL query from Query object
type QueryBuilder struct {
	Query      Query
	Model      *Model
	AddColumns []string
	sql        string
	args       []interface{}
}

// NewQueryBuilder creates QueryBuilder instance
func NewQueryBuilder(q Query, m *Model, addColumns ...string) *QueryBuilder {
	return &QueryBuilder{
		Query:      q,
		Model:      m,
		AddColumns: addColumns,
		args:       []interface{}{},
	}
}

func (qb *QueryBuilder) String() string {
	if qb.sql == "" {
		qb.compile()
	}
	return qb.sql
}

// Args returns query arguments
func (qb *QueryBuilder) Args() []interface{} {
	if len(qb.args) == 0 {
		if len(qb.Query.RawSQL.Arguments) > 0 {
			qb.args = qb.Query.RawSQL.Arguments
		} else {
			qb.compile()
		}
	}
	return qb.args
}

func (qb *QueryBuilder) compile() {
	if qb.sql == "" {
		if qb.Query.RawSQL.Fragment != "" {
			qb.sql = qb.Query.RawSQL.Fragment
		} else {
			qb.sql = qb.buildSelectSQL()
		}

		if inRegex.MatchString(qb.sql) {
			s, _, err := sqlx.In(qb.sql, qb.Args())
			if err == nil {
				qb.sql = s
			}
		}
		qb.sql = qb.Query.Connection.Dialect.TranslateSQL(qb.sql)
	}
}

func (qb *QueryBuilder) buildSelectSQL() string {
	cols := qb.buildColumns()

	fc := qb.buildfromClauses()

	sql := fmt.Sprintf("SELECT %s FROM %s", cols.String(), fc)

	sql = qb.buildJoinClauses(sql)
	sql = qb.buildWhereClauses(sql)
	sql = qb.buildGroupClauses(sql)
	sql = qb.buildOrderClauses(sql)
	sql = qb.buildPaginationClauses(sql)

	return sql
}

func (qb *QueryBuilder) buildfromClauses() FromClauses {
	models := []*Model{
		qb.Model,
	}

	fc := qb.Query.fromClauses
	for _, m := range models {
		tableName := m.TableName()
		asName := m.As
		if asName == "" {
			asName = strings.Replace(tableName, ".", "_", -1)
		}
		fc = append(fc, FromClause{
			From: tableName,
			As:   asName,
		})
	}

	return fc
}

func (qb *QueryBuilder) buildWhereClauses(sql string) string {
	wc := qb.Query.whereClauses
	if len(wc) > 0 {
		sql = fmt.Sprintf("%s WHERE %s", sql, wc.Join(" AND "))
		for _, arg := range wc.Args() {
			qb.args = append(qb.args, arg)
		}
	}
	return sql
}

func (qb *QueryBuilder) buildJoinClauses(sql string) string {
	oc := qb.Query.joinClauses
	if len(oc) > 0 {
		sql += " " + oc.String()
		for i := range oc {
			for _, arg := range oc[i].Arguments {
				qb.args = append(qb.args, arg)
			}
		}
	}

	return sql
}

func (qb *QueryBuilder) buildGroupClauses(sql string) string {
	gc := qb.Query.groupClauses
	if len(gc) > 0 {
		sql = fmt.Sprintf("%s GROUP BY %s", sql, gc.String())

		hc := qb.Query.havingClauses
		if len(hc) > 0 {
			sql = fmt.Sprintf("%s HAVING %s", sql, hc.String())
		}

		for i := range hc {
			for _, arg := range hc[i].Arguments {
				qb.args = append(qb.args, arg)
			}
		}
	}

	return sql
}

func (qb *QueryBuilder) buildOrderClauses(sql string) string {
	oc := qb.Query.orderClauses
	if len(oc) > 0 {
		sql = fmt.Sprintf("%s ORDER BY %s", sql, oc.Join(", "))
		for _, arg := range oc.Args() {
			qb.args = append(qb.args, arg)
		}
	}
	return sql
}

func (qb *QueryBuilder) buildPaginationClauses(sql string) string {
	if qb.Query.limitResults > 0 && qb.Query.Paginator == nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, qb.Query.limitResults)
	}
	if qb.Query.Paginator != nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, qb.Query.Paginator.PerPage)
		sql = fmt.Sprintf("%s OFFSET %d", sql, qb.Query.Paginator.Offset)
	}
	return sql
}

func (qb *QueryBuilder) buildColumns() Columns {
	cols := qb.Model.Columns()
	if len(qb.AddColumns) > 0 {
		cols.Add(qb.AddColumns...)
	}
	return cols
}
