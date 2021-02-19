package configloader

import (
	"os"

	"github.com/abelgoodwin1988/GoLuca/internal/config"
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml"
)

// Load at rest configuration into memory
// First load configuration files into the local configuration store.
// Second, load environmental variables to the local configuration store, overwriting pre-existing values, if any.
// Lastly, set configuration values with cli flags, overwriting pre-existing values, if any.
func Load() error {
	if err := loadConfigurationFile(); err != nil {
		return err
	}
	if err := loadEnvironmentVariables(); err != nil {
		return err
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
	if val := os.Getenv("GOLUCA_DBDRIVERNAME"); val != "" {
		config.Env.DBConnString = val
	}
	if val := os.Getenv("GOLUCA_DBCONNSTRING"); val != "" {
		config.Env.DBDriverName = val
	}
	return nil
}
