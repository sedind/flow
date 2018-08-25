package app

import (
	"github.com/sedind/flow/app/config"
)

// New creates instance of application Context
func New(configFile string) *Context {
	appConfig := Config{}

	err := config.LoadFromPath(configFile, &appConfig)
	if err != nil {
		panic(err)
	}

	return &Context{
		appConfig,
	}
}
