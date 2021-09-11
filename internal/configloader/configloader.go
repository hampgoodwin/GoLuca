package configloader

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/pelletier/go-toml"
)

// Load at rest configuration into memory
// First load configuration files into the local configuration store.
// Second, load environmental variables to the local configuration store, overwriting pre-existing values, if any.
// Lastly, set configuration values with cli flags, overwriting pre-existing values, if any.
func Load() error {
	if err := loadConfigurationFile(); err != nil {
		return errors.Wrap(err, "failed to load configuration from file")
	}
	if err := loadEnvironmentVariables(); err != nil {
		return errors.Wrap(err, "failed to load configuration from environment variables")
	}
	// load flags here at some point
	validate := validator.New()
	if err := validate.Struct(config.Env); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			return err
		}
	}
	return nil
}

func loadConfigurationFile() error {
	f, err := os.Open("/etc/goluca/.env.toml")
	if err != nil {
		return err
	}
	err = toml.NewDecoder(f).Decode(config.Env)
	if err != nil {
		return err
	}
	return nil
}

// loadEnvironmentVariables reads environmental variables and stores then into the config.Env
func loadEnvironmentVariables() error {
	if val := os.Getenv("GOLUCA_DBHOST"); val != "" {
		config.Env.DBHost = val
	}
	if val := os.Getenv("GOLUCA_DBPORT"); val != "" {
		config.Env.DBPort = val
	}
	if val := os.Getenv("GOLUCA_DBUSER"); val != "" {
		config.Env.DBUser = val
	}
	if val := os.Getenv("GOLUCA_DBPASS"); val != "" {
		config.Env.DBPass = val
	}
	if val := os.Getenv("GOLUCA_DBDATABASE"); val != "" {
		config.Env.DBDatabase = val
	}
	if val := os.Getenv("GOLUCA_APIHOST"); val != "" {
		config.Env.APIHost = val
	}
	if val := os.Getenv("GOLUCA_APIPORT"); val != "" {
		config.Env.APIPort = val
	}
	if val := os.Getenv("GOLUCA_ENVTYPE"); val != "" {
		config.Env.EnvType = val
	}
	return nil
}
