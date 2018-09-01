package migrate

import (
	"github.com/pkg/errors"
	"github.com/sedind/flow/app"
	"github.com/sedind/flow/app/config"
	"github.com/sedind/flow/app/dbe"
	"github.com/sedind/flow/app/dotenv"
	"github.com/spf13/cobra"
)

// upCmd generates sql migration files
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all of the 'up' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// load environment variables - this is needed as
		// config package utilizes environment variables loading
		dotenv.Load()

		if configFile == "" {
			return errors.New("config file not provided")
		}
		// get app config
		appConfig := app.Config{}
		err := config.LoadFromPath(configFile, &appConfig)
		if err != nil {
			return errors.Wrapf(err, "Unable to load configuration %s", configFile)
		}

		// get connection details for default connection string
		cd, ok := appConfig.ConnectionStrings[appConfig.DefaultConnection]
		if !ok {
			return errors.Errorf("Default Connection String configuration not provided in %s", configFile)
		}

		// ceate new DB connection
		dbConn, err := dbe.NewConnection(*cd)
		if err != nil {
			return errors.Wrap(err, "Unable to create database connection")
		}

		// open DB connection
		err = dbConn.Open()
		if err != nil {
			return errors.Wrapf(err, "Unable to connect to `%s` connection", appConfig.DefaultConnection)
		}

		fm, err := dbe.NewFileMigrator(appConfig.MigrationsPath, dbConn)
		if err != nil {
			return errors.Wrap(err, "Unable to create File Migration")
		}

		return fm.Up()
	},
}
