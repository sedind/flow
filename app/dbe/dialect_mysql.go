package dbe

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type mysql struct {
	commonDialect
}

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (mysql) Name() string {
	return "mysql"
}

func (s mysql) URL() string {
	d := s.details
	if d.URL != "" {
		return strings.TrimPrefix(d.URL, "mysql://")
	}
	url := "%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(url, d.User, d.Password, d.Host, d.Port, d.Database)
}

func (s mysql) urlWithoutDb() string {
	d := s.details
	if d.URL != "" {
		// respect user's own URL definition (with options).
		url := strings.TrimPrefix(d.URL, "mysql://")
		return strings.Replace(url, "/"+d.Database+"?", "/?", 1)
	}
	url := "%s:%s@(%s:%s)/?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(url, d.User, d.Password, d.Host, d.Port)
}

func (mysql) Quote(key string) string {
	return fmt.Sprintf("`%s`", key)
}

func (s mysql) RemoveIndex(tableName string, indexName string) error {
	_, err := s.store.Exec(fmt.Sprintf("DROP INDEX %v ON %v", indexName, s.Quote(tableName)))
	return err
}

func (s mysql) ModifyColumn(tableName string, columnName string, typ string) error {
	_, err := s.store.Exec(fmt.Sprintf("ALTER TABLE %v MODIFY COLUMN %v %v", tableName, columnName, typ))
	return err
}

func (s mysql) LimitAndOffsetSQL(limit, offset interface{}) (sql string) {
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

func (s mysql) HasForeignKey(tableName string, foreignKeyName string) bool {
	var count int
	currentDatabase, tableName := s.currentDatabaseAndTable(tableName)
	s.store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA=? AND TABLE_NAME=? AND CONSTRAINT_NAME=? AND CONSTRAINT_TYPE='FOREIGN KEY'", currentDatabase, tableName, foreignKeyName).Scan(&count)
	return count > 0
}

func (s mysql) CurrentDatabase() (name string) {
	s.store.QueryRow("SELECT DATABASE()").Scan(&name)
	return
}

func (mysql) SelectFromDummyTable() string {
	return "FROM DUAL"
}

func (s mysql) BuildKeyName(kind, tableName string, fields ...string) string {
	keyName := s.commonDialect.BuildKeyName(kind, tableName, fields...)
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

func (mysql) DefaultValueStr() string {
	return "VALUES()"
}
