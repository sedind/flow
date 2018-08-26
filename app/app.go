package app

import (
	"github.com/sedind/flow/app/config"
	"github.com/sedind/flow/app/dbe"
)

// New creates instance of application Context
func New(configFile string) *Context {
	appConfig := Config{}

	err := config.LoadFromPath(configFile, &appConfig)
	if err != nil {
		panic(err)
	}

	connections := map[string]*dbe.Connection{}

	for k, d := range appConfig.ConnectionStrings {
		c, err := dbe.NewConnection(d)
		if err != nil {
			panic(err)
		}
		connections[k] = c

	}

	return &Context{
		appConfig,
		connections,
	}
}
