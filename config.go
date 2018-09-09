package flow

import (
	"github.com/sedind/flow/dbe"
)

// Config -
type Config struct {
	Name              string                  `yaml:"name"`
	Addr              string                  `yaml:"addr"`
	LogLevel          string                  `yaml:"log_level"`
	RequestLogging    bool                    `yaml:"request_logging"`
	CompressResponse  bool                    `yaml:"compress_response"`
	RedirectSlashes   bool                    `yaml:"redirect_slashes"`
	PanicRecover      bool                    `yaml:"panic_recover"`
	NoCache           bool                    `yaml:"no_cache"`
	CORS              CORSConfig              `yaml:"cors"`
	MigrationsPath    string                  `yaml:"migrations_path"`
	DefaultConnection string                  `yaml:"default_connection"`
	ConnectionStrings map[string]*dbe.Details `yaml:"connection_strings"`
	AppSettings       map[string]string       `yaml:"app_settings"`
}

// CORSConfig -
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}
