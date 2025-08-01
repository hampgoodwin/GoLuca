package config

import "fmt"

// Config holds the values for environment loaded configuration values
type Config struct {
	Environment Environment
	Database    Database
	HTTPServer  HTTPServer
	GRPCServer  GRPCServer
	NATS        NATS
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
	SSLMode  string
}

func (d *Database) ConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		d.User,
		d.Pass,
		d.Host,
		d.Port,
		d.Database,
		d.SSLMode,
	)
}

type HTTPServer struct {
	Host string
	Port string
}

func (a *HTTPServer) AddressString() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

type GRPCServer struct {
	Host string
	Port string
}

func (s *GRPCServer) URL() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

type NATS struct {
	Host    string
	Port    string
	Wiretap NATSWiretap
}

func (n *NATS) URL() string {
	return fmt.Sprintf("nats://%s:%s", n.Host, n.Port)
}

const (
	WiretapEnabled  = "enabled"
	WiretapDisabled = "disabled"
)

// NATS Wiretap is by default enabled and host+port are expected
type NATSWiretap struct {
	Enable bool
	Host   string
	Port   string
}

func (nw *NATSWiretap) URL() string {
	return fmt.Sprintf("nats://%s:%s", nw.Host, nw.Port)
}
