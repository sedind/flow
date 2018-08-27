package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Configuration object
type Configuration struct {
	AppRoot string                   `yaml:"app_root"`
	Watcher map[string]WatcherConfig `yaml:"watcher"`
}

// WatcherConfig object
type WatcherConfig struct {
	Name              string   `yaml:"name"`
	Watch             string   `yaml:"watch"`
	ChangeCommand     string   `yaml:"change_command"`
	PostChangeCommand string   `yaml:"post_change_command"`
	Extensions        []string `yaml:"extensions"`
	Ignore            []string `yaml:"ignore"`
}

// Load loads configuration from file
func (c *Configuration) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

// Dump configuration object to file
func (c *Configuration) Dump(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}
