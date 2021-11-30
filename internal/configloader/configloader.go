package configloader

import (
	"os"

	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/errors"
	"github.com/pelletier/go-toml"
)

// Load at rest configuration into memory
// First load configuration files into the local configuration store.
// Second, load environmental variables to the local configuration store, overwriting pre-existing values, if any.
// Lastly, set configuration values with cli flags, overwriting pre-existing values, if any.
func Load(c config.Config, fp string) (config.Config, error) {
	cfg := config.Config{}
	if c != (config.Config{}) {
		cfg = c
	}

	fileCfg, err := loadConfigurationFile(fp)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to load configuration from file")
	}
	cfg = merge(cfg, fileCfg)

	envCfg := loadEnvironmentVariables()
	cfg = merge(cfg, envCfg)

	// TODO: Load command line flags

	if err := validate.Validate(cfg); err != nil {
		return cfg, errors.WithErrorMessage(err, errors.NotValid, "validating configuration")
	}
	return cfg, nil
}

func loadConfigurationFile(fp string) (config.Config, error) {
	cfg := config.Config{}
	if fp == "" {
		return cfg, nil
	}
	f, err := os.Open(fp)
	if err != nil {
		return config.Config{}, err
	}
	err = toml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return config.Config{}, err
	}
	return cfg, nil
}

const (
	EnvType    = "GOLUCA_ENVTYPE"
	DBHost     = "GOLUCA_DBHOST"
	DBPort     = "GOLUCA_DBPORT"
	DBUser     = "GOLUCA_DBUSER"
	DBPass     = "GOLUCA_DBPass"
	DBDatabase = "GOLUCA_DBDatabase"
	APIHost    = "GOLUCA_APIHost"
	APIPort    = "GOLUCA_APIPort"
)

var EnvironmentVariableKeys = []string{
	EnvType,
	DBHost,
	DBPort,
	DBUser,
	DBPass,
	DBDatabase,
	APIHost,
	APIPort,
}

// loadEnvironmentVariables reads environmental variables and stores then into the
// named return cfg
func loadEnvironmentVariables() (cfg config.Config) {
	if val := os.Getenv(EnvType); val != "" {
		cfg.Environment.Type = val
	}
	if val := os.Getenv(DBHost); val != "" {
		cfg.Database.Host = val
	}
	if val := os.Getenv(DBPort); val != "" {
		cfg.Database.Port = val
	}
	if val := os.Getenv(DBUser); val != "" {
		cfg.Database.User = val
	}
	if val := os.Getenv(DBPass); val != "" {
		cfg.Database.Pass = val
	}
	if val := os.Getenv(DBDatabase); val != "" {
		cfg.Database.Database = val
	}
	if val := os.Getenv(APIHost); val != "" {
		cfg.HTTPAPI.Host = val
	}
	if val := os.Getenv(APIPort); val != "" {
		cfg.HTTPAPI.Port = val
	}
	return
}

func merge(a, b config.Config) config.Config {
	if b.Environment.Type != "" {
		a.Environment.Type = b.Environment.Type
	}
	if b.Database.Host != "" {
		a.Database.Host = b.Database.Host
	}
	if b.Database.Port != "" {
		a.Database.Port = b.Database.Port
	}
	if b.Database.User != "" {
		a.Database.User = b.Database.User
	}
	if b.Database.Pass != "" {
		a.Database.Pass = b.Database.Pass
	}
	if b.Database.Database != "" {
		a.Database.Database = b.Database.Database
	}
	if b.HTTPAPI.Host != "" {
		a.HTTPAPI.Host = b.HTTPAPI.Host
	}
	if b.HTTPAPI.Port != "" {
		a.HTTPAPI.Port = b.HTTPAPI.Port
	}
	return a
}
