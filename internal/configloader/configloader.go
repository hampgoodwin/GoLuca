package configloader

import (
	"os"

	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/pelletier/go-toml"
)

// Load at rest configuration into memory
// First load configuration files into the local configuration store.
// Second, load environmental variables to the local configuration store, overwriting pre-existing values, if any.
// Lastly, set configuration values with cli flags, overwriting pre-existing values, if any.
func Load(c *config.Config) (*config.Config, error) {
	cfg := &config.Config{}
	if c != nil {
		cfg = c
	}

	if err := loadConfigurationFile(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to load configuration from file")
	}
	loadEnvironmentVariables(cfg)
	// TODO: Load command line flags
	if err := validate.Validate(cfg); err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotValid, "validating configuration")
	}
	return cfg, nil
}

func loadConfigurationFile(cfg *config.Config) error {
	f, err := os.Open("/etc/goluca/.env.toml")
	if err != nil {
		return err
	}
	err = toml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

// loadEnvironmentVariables reads environmental variables and stores then into the config.Env
func loadEnvironmentVariables(cfg *config.Config) {
	if val := os.Getenv("GOLUCA_DBHOST"); val != "" {
		cfg.DBHost = val
	}
	if val := os.Getenv("GOLUCA_DBPORT"); val != "" {
		cfg.DBPort = val
	}
	if val := os.Getenv("GOLUCA_DBUSER"); val != "" {
		cfg.DBUser = val
	}
	if val := os.Getenv("GOLUCA_DBPASS"); val != "" {
		cfg.DBPass = val
	}
	if val := os.Getenv("GOLUCA_DBDATABASE"); val != "" {
		cfg.DBDatabase = val
	}
	if val := os.Getenv("GOLUCA_APIHOST"); val != "" {
		cfg.APIHost = val
	}
	if val := os.Getenv("GOLUCA_APIPORT"); val != "" {
		cfg.APIPort = val
	}
	if val := os.Getenv("GOLUCA_ENVTYPE"); val != "" {
		cfg.EnvType = val
	}
}
