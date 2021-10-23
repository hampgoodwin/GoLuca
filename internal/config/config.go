package config

import "fmt"

// Config holds the values for environment loaded configuration values
type Config struct {
	Environment
	Database
	HTTPAPI
}

type Environment struct {
	Type string `validate:"required,oneof=LOCAL DEV STAGING PROD"`
}

type Database struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

func (d *Database) ConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		d.User,
		d.Pass,
		d.Host,
		d.Port,
		d.Database,
	)
}

type HTTPAPI struct {
	Host string
	Port string
}

func (a *HTTPAPI) AddressString() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}
