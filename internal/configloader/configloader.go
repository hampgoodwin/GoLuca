package configloader

import (
	"log"
	"os"
	"strconv"

	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/errors"
	"github.com/pelletier/go-toml/v2"
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

	cfg = loadAndMergeEnvironmentVariables(cfg)

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
	defer func() {
		// TODO: replace with global logger
		log.Printf("closing configuration file: %v", err)
	}()

	err = toml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return config.Config{}, err
	}
	return cfg, nil
}

const (
	EnvType        = "GOLUCA_ENV_TYPE"
	DBHost         = "GOLUCA_DB_HOST"
	DBPort         = "GOLUCA_DB_PORT"
	DBUser         = "GOLUCA_DB_USER"
	DBPass         = "GOLUCA_DB_PASS"
	DBDatabase     = "GOLUCA_DB_DATABASE"
	DBSSLMode      = "GOLUCA_DB_SSLMODE"
	HTTPServerHost = "GOLUCA_HTTP_SERVER_HOST"
	HTTPServerPort = "GOLUCA_HTTP_SERVER_PORT"
	GRPCServerHost = "GOLUCA_GRPC_SERVER_HOST"
	GRPCServerPort = "GOLUCA_GRPC_SERVER_PORT"
	NATSHost       = "GOLUCA_NATS_HOST"
	NATSPort       = "GOLUCA_NATS_PORT"
	WiretapEnable  = "GOLUCA_WIRETAP_ENABLE"
	WiretapHost    = "GOLUCA_WIRETAP_HOST"
	WiretapPort    = "GOLUCA_WIRETAP_PORT"
)

var EnvironmentVariableKeys = []string{
	EnvType,
	DBHost,
	DBPort,
	DBUser,
	DBPass,
	DBDatabase,
	DBSSLMode,
	HTTPServerHost,
	HTTPServerPort,
	GRPCServerHost,
	GRPCServerPort,
	NATSHost,
	NATSPort,
	WiretapEnable,
	WiretapHost,
	WiretapPort,
}

// loadAndMergeEnvironmentVariables reads environmental variables and stores then into the
// named return cfg
func loadAndMergeEnvironmentVariables(in config.Config) (cfg config.Config) {
	if in != (config.Config{}) {
		cfg = in
	}
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
	if val := os.Getenv(DBSSLMode); val != "" {
		cfg.Database.SSLMode = val
	}
	if val := os.Getenv(HTTPServerHost); val != "" {
		cfg.HTTPServer.Host = val
	}
	if val := os.Getenv(HTTPServerPort); val != "" {
		cfg.HTTPServer.Port = val
	}
	if val := os.Getenv(GRPCServerHost); val != "" {
		cfg.GRPCServer.Host = val
	}
	if val := os.Getenv(GRPCServerPort); val != "" {
		cfg.GRPCServer.Port = val
	}
	if val := os.Getenv(NATSHost); val != "" {
		cfg.NATS.Host = val
	}
	if val := os.Getenv(NATSPort); val != "" {
		cfg.NATS.Port = val
	}
	if val := os.Getenv(WiretapEnable); val != "" {
		enabled, _ := strconv.ParseBool(val)
		cfg.NATS.Wiretap.Enable = enabled
	}
	if val := os.Getenv(WiretapHost); val != "" {
		cfg.NATS.Wiretap.Host = val
	}
	if val := os.Getenv(WiretapPort); val != "" {
		cfg.NATS.Wiretap.Port = val
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
	if b.Database.SSLMode != "" {
		a.Database.SSLMode = b.Database.SSLMode
	}
	if b.HTTPServer.Host != "" {
		a.HTTPServer.Host = b.HTTPServer.Host
	}
	if b.HTTPServer.Port != "" {
		a.HTTPServer.Port = b.HTTPServer.Port
	}
	if b.GRPCServer.Host != "" {
		a.GRPCServer.Host = b.GRPCServer.Host
	}
	if b.GRPCServer.Port != "" {
		a.GRPCServer.Port = b.GRPCServer.Port
	}
	if b.NATS.Host != "" {
		a.NATS.Host = b.NATS.Host
	}
	if b.NATS.Port != "" {
		a.NATS.Port = b.NATS.Port
	}
	a.NATS.Wiretap.Enable = b.NATS.Wiretap.Enable
	if b.NATS.Wiretap.Host != "" {
		a.NATS.Wiretap.Host = b.NATS.Wiretap.Host
	}
	if b.NATS.Wiretap.Port != "" {
		a.NATS.Wiretap.Port = b.NATS.Wiretap.Port
	}
	return a
}
