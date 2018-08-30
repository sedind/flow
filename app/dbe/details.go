package dbe

import (
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
	Dialect  string            `yaml:"dialect"`
	Database string            `yaml:"database"`
	Host     string            `yaml:"host"`
	Port     string            `yaml:"port"`
	User     string            `yaml:"user"`
	Password string            `yaml:"password"`
	URL      string            `yaml:"url"`
	Pool     int               `yaml:"pool"`
	IdlePool int               `yaml:"idle_pool"`
	Options  map[string]string `yaml:"options"`
}

var dialectRegex = regexp.MustCompile(`\s+:\/\/`)

// Finalize cleans up the connection details by normalizing names
func (d *Details) Finalize() error {

	switch strings.ToLower(d.Dialect) {
	case "mysql":
		if d.URL != "" {
			cfg, err := _mysql.ParseDSN(d.URL)
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
