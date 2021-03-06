package dbe

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

var rLimitOffset = regexp.MustCompile("(?i)(limit [0-9]+ offset [0-9]+)$")
var rLimit = regexp.MustCompile("(?i)(limit [0-9]+)$")

// Exec runs the given query
func (q *Query) Exec() error {

	sql, args := q.ToSQL(nil)
	Logger.Info(fmt.Sprintf("%s | %s", sql, args))
	_, err := q.Connection.Store.Exec(sql, args...)
	return err
}

// ExecWithCount Execute and count
func (q *Query) ExecWithCount() (int64, error) {
	sql, args := q.ToSQL(nil)
	Logger.Info(fmt.Sprintf("%s | %s", sql, args))
	result, err := q.Connection.Store.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// First record of the model in the database that matches the query.
//
//	q.Where("name = ?", "Irfan").First(&User{})
func (q *Query) First(model interface{}) error {
	m := &Model{Value: model}
	q.Limit(1)
	sql, args := q.ToSQL(m)
	Logger.Info(fmt.Sprintf("%s | %s", sql, args))

	return q.Connection.Store.Get(m.Value, sql, args...)
}

// Last record of the model in the database that matches the query.
//
//	q.Where("name = ?", "Irfan").Last(&User{})
func (q *Query) Last(model interface{}) error {
	m := &Model{Value: model}
	q.Order("id DESC")
	q.Limit(1)
	sql, args := q.ToSQL(m)
	Logger.Info(fmt.Sprintf("%s | %s", sql, args))
	return q.Connection.Store.Get(m.Value, sql, args...)
}

// Find the first record of the model in the database with a particular id.
//
//	q.Find(&User{}, 1)
func (q *Query) Find(model interface{}, id interface{}) error {
	m := &Model{Value: model}
	idq := fmt.Sprintf("%s.id = ?", m.TableName())
	switch t := id.(type) {
	case string:
		var err error
		id, err = strconv.Atoi(t)
		if err != nil {
			return q.Where(idq, t).First(model)
		}
	}

	return q.Where(idq, id).First(model)
}

// All retrieves all of the records in the database that match the query.
//
//	q.Where("name = ?", "Irfan").All(&[]User{})
func (q *Query) All(models interface{}) error {
	m := &Model{Value: models}
	sql, args := q.ToSQL(m)
	Logger.Info(fmt.Sprintf("%s | %s", sql, args))
	err := q.Connection.Store.Select(m.Value, sql, args...)
	if err == nil && q.Paginator != nil {
		ct, err := q.Count(models)
		if err == nil {
			q.Paginator.TotalEntriesSize = ct
			st := reflect.ValueOf(models).Elem()
			q.Paginator.CurrentEntriesSize = st.Len()
			q.Paginator.TotalPages = (q.Paginator.TotalEntriesSize / q.Paginator.PerPage)
			if q.Paginator.TotalEntriesSize%q.Paginator.PerPage > 0 {
				q.Paginator.TotalPages = q.Paginator.TotalPages + 1
			}
		}
	}
	return err
}

// Exists returns true/false if a record exists in the database that matches
// the query.
//
// 	q.Where("name = ?", "Irfan").Exists(&User{})
func (q *Query) Exists(model interface{}) (bool, error) {
	i, err := q.Count(model)
	return i != 0, err
}

// Count the number of records in the database.
//
//	q.Where("name = ?", "Irfan").Count(&User{})
func (q Query) Count(model interface{}) (int, error) {
	return q.CountByField(model, "*")
}

// CountByField counts the number of records in the database, for a given field.
//
//	q.Where("gender = ?", "male").Count(&User{}, "name")
func (q Query) CountByField(model interface{}, field string) (int, error) {
	tmpQuery := NewQuery(q.Connection)
	q.Clone(tmpQuery) //avoid mendling with original query

	res := &rowCount{}

	tmpQuery.Paginator = nil
	tmpQuery.orderClauses = Clauses{}
	tmpQuery.limitResults = 0
	query, args := tmpQuery.ToSQL(&Model{Value: model})
	//when query contains custom selected fields / executed using RawQuery,
	//	sql may already contains limit and offset

	if rLimitOffset.MatchString(query) {
		foundLimit := rLimitOffset.FindString(query)
		query = query[0 : len(query)-len(foundLimit)]
	} else if rLimit.MatchString(query) {
		foundLimit := rLimit.FindString(query)
		query = query[0 : len(query)-len(foundLimit)]
	}

	countQuery, err := q.Connection.Dialect.CountStmt(field, query)
	if err != nil {
		return 0, err
	}
	Logger.Info(fmt.Sprintf("%s | %s", countQuery, args))
	err = q.Connection.Store.Get(res, countQuery, args...)
	if err != nil {
		return 0, err
	}

	return res.Count, err
}

// rowCount is helper struct for CountByField query
type rowCount struct {
	Count int `db:"row_count"`
}
