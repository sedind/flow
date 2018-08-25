package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

var errWrongConfigurationType = errors.New("Configuration type must be a pointer to a struct")

// LoadFromPath reads configuration from path and stores it to obj interface
// The format is deduced from the file extension
//	* .yml     - is decoded as yaml
func LoadFromPath(path string, obj interface{}) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	data, err := os.Open(path)
	if err != nil {
		return err
	}

	return LoadFromReader(data, obj)
}

// LoadFromReader reads configuration from reader and stores it to obj interface
// The format is deduced from the file extension
//	* .yml     - is decoded as yaml
func LoadFromReader(reader io.Reader, obj interface{}) error {
	err := checkConfigObj(obj)
	if err != nil {
		return errors.WithStack(err)
	}

	tmpl := template.New("app_config")
	tmpl.Funcs(map[string]interface{}{
		"envOr": func(envKey, defaultVal string) string {
			return getEnv(envKey, defaultVal)
		},
		"env": func(envKey string) string {
			return getEnv(envKey, "")
		},
	})

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.WithStack(err)
	}

	t, err := tmpl.Parse(string(b))
	if err != nil {
		return errors.Wrap(err, "couldn't parse config")
	}

	var bb bytes.Buffer
	err = t.Execute(&bb, nil)
	if err != nil {
		return errors.Wrap(err, "couldn't execute config")
	}
	err = yaml.Unmarshal(bb.Bytes(), obj)
	if err != nil {
		return errors.Wrap(err, "couldn't unmarshal config to yaml")
	}
	return nil
}

func getEnv(envKey, defaultValue string) string {
	val := os.Getenv(envKey)
	if len(val) > 0 {
		return val
	}
	return defaultValue
}

func checkConfigObj(obj interface{}) error {
	// check if type is a pointer
	objVal := reflect.ValueOf(obj)
	if objVal.Kind() != reflect.Ptr || objVal.IsNil() {
		return errWrongConfigurationType
	}

	// get and confirm struct value
	objVal = objVal.Elem()
	if objVal.Kind() != reflect.Struct {
		return errWrongConfigurationType
	}
	return nil
}
