package config

// Config holds the values for environment loaded configuration values
type Config struct {
	EnvType                  string `validate:"required,oneof=LOCAL DEV STAGING PROD"`
	DBHost, DBUser, DBPass   string
	DBDatabase               string
	DBPort, APIHost, APIPort string
}
