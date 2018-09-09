package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sedind/flow"
	"github.com/sedind/flow/flow/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var appConfig, flowConfig bool

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&appConfig, "app", false, "generate app only project configuration")
	initCmd.Flags().BoolVar(&flowConfig, "flow", false, "Projectgenerate flow only project configuration")
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize Flow project",
	RunE: func(cmd *cobra.Command, args []string) error {

		projectName := ""
		if len(args) <= 0 {
			pwd, _ := os.Getwd()
			projectName = filepath.Base(pwd)
		} else {
			projectName = args[0]
		}

		if appConfig {
			err := createAppConfig(projectName)
			if err != nil {
				return err
			}
		}

		if flowConfig {
			err := createFlowConfig(projectName)
			if err != nil {
				return err
			}
		}

		if !appConfig && !flowConfig {
			//generate flow and app config files
			err := createAppConfig(projectName)
			if err != nil {
				return err
			}
			err = createFlowConfig(projectName)
		}

		return nil
	},
}

func createAppConfig(name string) error {
	appConfig := flow.Config{}
	appConfig.Name = name
	appConfig.Addr = "0.0.0.0:3000"
	appConfig.LogLevel = "debug"
	appConfig.RequestLogging = true
	appConfig.CompressResponse = true
	appConfig.RedirectSlashes = true
	appConfig.PanicRecover = true
	appConfig.CORS = flow.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	appConfig.MigrationsPath = "migrations"
	appConfig.AppSettings = map[string]string{}

	return saveObjToFile("config.yml", &appConfig)
}

func createFlowConfig(name string) error {
	bName := strings.Replace(name, " ", "", -1)
	tmpPath := fmt.Sprintf("tmp/%s", bName)

	flowConfig := config.Configuration{}
	flowConfig.AppRoot = "."
	flowConfig.Watcher = map[string]config.WatcherConfig{
		"app": config.WatcherConfig{
			Name:              "Application watcher",
			Watch:             ".",
			ChangeCommand:     "go",
			ChangeArgs:        []string{"build", "-v", "-i", "-o", tmpPath},
			PostChangeCommand: tmpPath,
			PostChangeArgs:    []string{},
			Extensions:        []string{".go", ".yml"},
			Ignore:            []string{},
		},
	}
	return saveObjToFile("flow.yml", &flowConfig)
}

func saveObjToFile(path string, obj interface{}) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}
