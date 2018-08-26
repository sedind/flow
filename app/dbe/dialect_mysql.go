package dbe

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// DialectMySQL dialect implementation speciffic to mysql engine
type DialectMySQL struct {
	DialectCommon
}

//Name - returns dialect name
func (d *DialectMySQL) Name() string {
	return "mysql"
}

//URL gets connection string url
func (d *DialectMySQL) URL() string {
	c := d.details
	if c.URL != "" {
		return strings.TrimPrefix(c.URL, "mysql://")
	}
	s := "%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

// Quote -
func (d *DialectMySQL) Quote(key string) string {
	return fmt.Sprintf("`%s`", key)
}

// RemoveIndex -
func (d *DialectMySQL) RemoveIndex(tableName string, indexName string) error {
	_, err := d.store.Exec(fmt.Sprintf("DROP INDEX %v ON %v", indexName, d.Quote(tableName)))
	return err
}

// ModifyColumn -
func (d *DialectMySQL) ModifyColumn(tableName string, columnName string, typ string) error {
	_, err := d.store.Exec(fmt.Sprintf("ALTER TABLE %v MODIFY COLUMN %v %v", tableName, columnName, typ))
	return err
}

// LimitAndOffsetSQL -
func (d *DialectMySQL) LimitAndOffsetSQL(limit, offset interface{}) (sql string) {
	if limit != nil {
		if parsedLimit, err := strconv.ParseInt(fmt.Sprint(limit), 0, 0); err == nil && parsedLimit >= 0 {
			sql += fmt.Sprintf(" LIMIT %d", parsedLimit)

			if offset != nil {
				if parsedOffset, err := strconv.ParseInt(fmt.Sprint(offset), 0, 0); err == nil && parsedOffset >= 0 {
					sql += fmt.Sprintf(" OFFSET %d", parsedOffset)
				}
			}
		}
	}
	return
}

// HasForeignKey -
func (d *DialectMySQL) HasForeignKey(tableName string, foreignKeyName string) bool {
	var count int
	currentDatabase, tableName := d.currentDatabaseAndTable(tableName)
	d.store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA=? AND TABLE_NAME=? AND CONSTRAINT_NAME=? AND CONSTRAINT_TYPE='FOREIGN KEY'", currentDatabase, tableName, foreignKeyName).Scan(&count)
	return count > 0
}

// CurrentDatabase -
func (d *DialectMySQL) CurrentDatabase() (name string) {
	d.store.QueryRow("SELECT DATABASE()").Scan(&name)
	return
}

//SelectFromDummyTable -
func (d *DialectMySQL) SelectFromDummyTable() string {
	return "FROM DUAL"
}

//BuildKeyName -
func (d *DialectMySQL) BuildKeyName(kind, tableName string, fields ...string) string {
	keyName := d.DialectCommon.BuildKeyName(kind, tableName, fields...)
	if utf8.RuneCountInString(keyName) <= 64 {
		return keyName
	}
	h := sha1.New()
	h.Write([]byte(keyName))
	bs := h.Sum(nil)

	// sha1 is 40 characters, keep first 24 characters of destination
	destRunes := []rune(regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(fields[0], "_"))
	if len(destRunes) > 24 {
		destRunes = destRunes[:24]
	}

	return fmt.Sprintf("%s%x", string(destRunes), bs)
}

//DefaultValueStr -
func (d *DialectMySQL) DefaultValueStr() string {
	return "VALUES()"
}
