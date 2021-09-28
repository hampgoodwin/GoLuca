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
	DBPass     = "GOLUCA_DBPASS"
	DBDatabase = "GOLUCA_DBDATABASE"
	APIHost    = "GOLUCA_APIHOST"
	APIPort    = "GOLUCA_APIPORT"
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

// loadEnvironmentVariables reads environmental variables and stores then into the config.Env
func loadEnvironmentVariables() (cfg config.Config) {
	if val := os.Getenv(EnvType); val != "" {
		cfg.EnvType = val
	}
	if val := os.Getenv(DBHost); val != "" {
		cfg.DBHost = val
	}
	if val := os.Getenv(DBPort); val != "" {
		cfg.DBPort = val
	}
	if val := os.Getenv(DBUser); val != "" {
		cfg.DBUser = val
	}
	if val := os.Getenv(DBPass); val != "" {
		cfg.DBPass = val
	}
	if val := os.Getenv(DBDatabase); val != "" {
		cfg.DBDatabase = val
	}
	if val := os.Getenv(APIHost); val != "" {
		cfg.APIHost = val
	}
	if val := os.Getenv(APIPort); val != "" {
		cfg.APIPort = val
	}
	return
}

func merge(a, b config.Config) config.Config {
	if b.EnvType != "" {
		a.EnvType = b.EnvType
	}
	if b.DBHost != "" {
		a.DBHost = b.DBHost
	}
	if b.DBUser != "" {
		a.DBUser = b.DBUser
	}
	if b.DBPass != "" {
		a.DBPass = b.DBPass
	}
	if b.DBDatabase != "" {
		a.DBDatabase = b.DBDatabase
	}
	if b.DBPort != "" {
		a.DBPort = b.DBPort
	}
	if b.APIHost != "" {
		a.APIHost = b.APIHost
	}
	if b.APIPort != "" {
		a.APIPort = b.APIPort
	}
	return a
}
