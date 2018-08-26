package dbe

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	_mysql "github.com/go-sql-driver/mysql"
	"github.com/sedind/flow/app/defaults"

	"github.com/pkg/errors"
)

// Details represents data needed to connect to a datasource
type Details struct {
	// Example: "postgres" or "sqlite3" or "mysql"
	Dialect string
	// The name of your database. Example: "foo_development"
	Database string
	// The host of your database. Example: "127.0.0.1"
	Host string
	// The port of your database. Example: 1234
	// Will default to the "default" port for each dialect.
	Port string
	// The username of the database user. Example: "root"
	User string
	// The password of the database user. Example: "password"
	Password string
	// Instead of specifying each individual piece of the
	// connection you can instead just specify the URL of the
	// database. Example: "postgres://postgres:postgres@localhost:5432/pop_test?sslmode=disable"
	URL string
	// Defaults to 0 "unlimited". See https://golang.org/pkg/database/sql/#DB.SetMaxOpenConns
	Pool int
	// Defaults to 0 "unlimited". See https://golang.org/pkg/database/sql/#DB.SetMaxIdleConns
	IdlePool int
	Options  map[string]string
}

var dialectRegex = regexp.MustCompile(`\s+:\/\/`)

// Finalize cleans up the connection details by normalizing names
func (d *Details) Finalize() error {
	// check if db connection is passed in form of URL and construct details from URL
	if d.URL != "" {
		ul := d.URL
		if d.Dialect != "" {
			if !dialectRegex.MatchString(ul) {
				ul = d.Dialect + "://" + ul
			}
		}
		d.Database = d.URL
		if !strings.HasPrefix(d.Dialect, "sqlite") {
			u, err := url.Parse(ul)
			if err != nil {
				return errors.Wrapf(err, "could not parse %s", ul)
			}
			d.Dialect = u.Scheme
			d.Database = u.Path

			hp := strings.Split(u.Host, ":")
			d.Host = hp[0]
			if len(hp) > 1 {
				d.Port = hp[1]
			}

			if u.User != nil {
				d.User = u.User.Username()
				d.Password, _ = u.User.Password()
			}
		}
	}
	switch strings.ToLower(d.Dialect) {
	case "postgres", "postgresql", "pg":
		d.Dialect = "postgres"
		d.Port = defaults.String(d.Port, "5432")
		d.Database = strings.TrimPrefix(d.Database, "/")
	case "mysql":
		if d.URL != "" {
			cfg, err := _mysql.ParseDSN(strings.TrimPrefix(d.URL, "mysql://"))
			if err != nil {
				return errors.Wrap(err, "The URL is not supported by MySQL driver")
			}
			d.User = cfg.User
			d.Password = cfg.Passwd
			d.Database = cfg.DBName
			addr := strings.TrimSuffix(strings.TrimPrefix(cfg.Addr, "("), ")")
			if cfg.Net == "unix" {
				d.Port = "socket"
				d.Host = addr
			} else {
				tmp := strings.Split(addr, ":")
				switch len(tmp) {
				case 1:
					d.Host = tmp[0]
					d.Port = "3306"
				case 2:
					d.Host = tmp[0]
					d.Port = tmp[1]
				}
			}
		} else {
			d.Port = defaults.String(d.Port, "3306")
			d.Database = strings.TrimPrefix(d.Database, "/")
		}
	case "sqlite", "sqlite3":
		d.Dialect = "sqlite3"
	default:
		return errors.Errorf("Unsupported dialect `%s`!", d.Dialect)
	}
	return nil
}

// RetrySleep returns the amount of time to wait between two connection retries
func (d *Details) RetrySleep() time.Duration {
	dur, err := time.ParseDuration(defaults.String(d.Options["retry_sleep"], "1ms"))
	if err != nil {
		return 1 * time.Millisecond
	}
	return dur
}

// RetryLimit returns the maximum number of accepted connection retries
func (d *Details) RetryLimit() int {
	i, err := strconv.Atoi(defaults.String(d.Options["retry_limit"], "1000"))
	if err != nil {
		return 100
	}
	return i
}

// MigrationTableName returns the name of the table to track migrations
func (d *Details) MigrationTableName() string {
	return defaults.String(d.Options["migration_table_name"], "schema_migration")
}
