package configloader

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-playground/validator"
	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/errors"
	"github.com/matryer/is"
)

func TestLoad(t *testing.T) {
	testCases := []struct {
		description string
		cfg         config.Config
		filepath    string
		envVars     map[string]string
		expected    config.Config
		err         error // TODO add specific err case catches
	}{
		{
			description: "json-file-empty-vars-empty-config-error-decoding",
			filepath:    "../../test/data/configloader/json.env.toml",
			err:         errors.New("(1, 1): parsing error: keys cannot contain { character"),
		},
		{
			description: "empty-file-empty-vars-empty-config-err-validation",
			filepath:    "../../test/data/configloader/empty.env.toml",
			err:         validator.ValidationErrors{},
		},
		{
			description: "full-file-empty-vars-empty-config",
			filepath:    "../../.env.toml.example",
			expected: config.Config{
				Environment: config.Environment{Type: "LOCAL"},
				Database: config.Database{
					Host:     "127.0.0.1",
					Port:     "5432",
					User:     "user",
					Pass:     "password",
					Database: "goluca",
					SSLMode:  "disable",
				},
				HTTPServer: config.HTTPServer{
					Host: "localhost",
					Port: "3333",
				},
				GRPCServer: config.GRPCServer{
					Host: "localhost",
					Port: "5000",
				},
				NATS: config.NATS{
					Host: "localhost",
					Port: "4222",
					Wiretap: config.NATSWiretap{
						Enable: config.WiretapEnabled,
						Host:   "localhost",
						Port:   "4222",
					},
				},
			},
		},
		{
			description: "full-file-full-vars-empty-config-overwrite-file",
			filepath:    "../../.env.toml.example",
			envVars: map[string]string{
				EnvType:        "DEV",
				DBHost:         "GOLUCA_DB_HOST",
				DBPort:         "GOLUCA_DB_PORT",
				DBUser:         "GOLUCA_DB_USER",
				DBPass:         "GOLUCA_DB_PASS",
				DBDatabase:     "GOLUCA_DB_DATABASE",
				DBSSLMode:      "GOLUCA_DB_SSLMODE",
				HTTPServerHost: "GOLUCA_HTTP_SERVER_HOST",
				HTTPServerPort: "GOLUCA_HTTP_SERVER_PORT",
				GRPCServerHost: "GOLUCA_GRPC_SERVER_HOST",
				GRPCServerPort: "GOLUCA_GRPC_SERVER_PORT",
				NATSHost:       "GOLUCA_NATS_HOST",
				NATSPort:       "GOLUCA_NATS_PORT",
				WiretapEnable:  "GOLUCA_WIRETAP_ENABLE",
				WiretapHost:    "GOLUCA_WIRETAP_HOST",
				WiretapPort:    "GOLUCA_WIRETAP_PORT",
			},
			expected: config.Config{
				Environment: config.Environment{Type: "DEV"},
				Database: config.Database{
					Host:     "GOLUCA_DB_HOST",
					Port:     "GOLUCA_DB_PORT",
					User:     "GOLUCA_DB_USER",
					Pass:     "GOLUCA_DB_PASS",
					Database: "GOLUCA_DB_DATABASE",
					SSLMode:  "GOLUCA_DB_SSLMODE",
				},
				HTTPServer: config.HTTPServer{
					Host: "GOLUCA_HTTP_SERVER_HOST",
					Port: "GOLUCA_HTTP_SERVER_PORT",
				},
				GRPCServer: config.GRPCServer{
					Host: "GOLUCA_GRPC_SERVER_HOST",
					Port: "GOLUCA_GRPC_SERVER_PORT",
				},
				NATS: config.NATS{
					Host: "GOLUCA_NATS_HOST",
					Port: "GOLUCA_NATS_PORT",
					Wiretap: config.NATSWiretap{
						Enable: config.WiretapEnabled,
						Host:   "GOLUCA_WIRETAP_HOST",
						Port:   "GOLUCA_WIRETAP_PORT",
					},
				},
			},
		},
		{
			description: "full-file-partial-vars-empty-config-overwrite-file",
			filepath:    "../../.env.toml.example",
			envVars: map[string]string{
				EnvType:        "DEV",
				HTTPServerHost: "GOLUCA_HTTP_SERVER_HOST",
				HTTPServerPort: "GOLUCA_HTTP_SERVER_PORT",
			},
			expected: config.Config{
				Environment: config.Environment{Type: "DEV"},
				Database: config.Database{
					Host:     "127.0.0.1",
					Port:     "5432",
					User:     "user",
					Pass:     "password",
					Database: "goluca",
					SSLMode:  "disable",
				},
				HTTPServer: config.HTTPServer{
					Host: "GOLUCA_HTTP_SERVER_HOST",
					Port: "GOLUCA_HTTP_SERVER_PORT",
				},
				GRPCServer: config.GRPCServer{
					Host: "localhost",
					Port: "5000",
				},
				NATS: config.NATS{
					Host: "localhost",
					Port: "4222",
					Wiretap: config.NATSWiretap{
						Enable: config.WiretapEnabled,
						Host:   "localhost",
						Port:   "4222",
					},
				},
			},
		},
		{
			description: "full-file-partial-vars-partial-config-overwrite-file",
			cfg:         config.Config{Environment: config.Environment{Type: "DEV"}},
			filepath:    "../../.env.toml.example",
			envVars: map[string]string{
				EnvType:        "LOCAL",
				HTTPServerHost: "GOLUCA_HTTP_SERVER_HOST",
				HTTPServerPort: "GOLUCA_HTTP_SERVER_PORT",
			},
			expected: config.Config{
				Environment: config.Environment{Type: "LOCAL"},
				Database: config.Database{
					Host:     "127.0.0.1",
					Port:     "5432",
					User:     "user",
					Pass:     "password",
					Database: "goluca",
					SSLMode:  "disable",
				},
				HTTPServer: config.HTTPServer{
					Host: "GOLUCA_HTTP_SERVER_HOST",
					Port: "GOLUCA_HTTP_SERVER_PORT",
				},
				GRPCServer: config.GRPCServer{
					Host: "localhost",
					Port: "5000",
				},
				NATS: config.NATS{
					Host: "localhost",
					Port: "4222",
					Wiretap: config.NATSWiretap{
						Enable: config.WiretapEnabled,
						Host:   "localhost",
						Port:   "4222",
					},
				},
			},
		},
		{
			description: "full-file-partial-vars-partial-config-persist-merge",
			cfg:         config.Config{Environment: config.Environment{Type: "DEV"}},
			filepath:    "../../.env.toml.example",
			envVars: map[string]string{
				HTTPServerHost: "GOLUCA_HTTP_SERVER_HOST",
				HTTPServerPort: "GOLUCA_HTTP_SERVER_PORT",
			},
			expected: config.Config{
				Environment: config.Environment{Type: "LOCAL"},
				Database: config.Database{
					Host:     "127.0.0.1",
					Port:     "5432",
					User:     "user",
					Pass:     "password",
					Database: "goluca",
					SSLMode:  "disable",
				},
				HTTPServer: config.HTTPServer{
					Host: "GOLUCA_HTTP_SERVER_HOST",
					Port: "GOLUCA_HTTP_SERVER_PORT",
				},
				GRPCServer: config.GRPCServer{
					Host: "localhost",
					Port: "5000",
				},
				NATS: config.NATS{
					Host: "localhost",
					Port: "4222",
					Wiretap: config.NATSWiretap{
						Enable: config.WiretapEnabled,
						Host:   "localhost",
						Port:   "4222",
					},
				},
			},
		},
	}

	a := is.New(t)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			resetApplicationEnvironmentVariables()
			defer resetApplicationEnvironmentVariables()

			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			actual, err := Load(tc.cfg, tc.filepath)
			if tc.err != nil {
				a.True(err != nil)
				// Using errors.As because it detects validator.ValidationErrors
				a.True(errors.As(err, &tc.err))
				return
			}
			a.NoErr(err)
			a.Equal(tc.expected, actual)
		})
	}
}

func TestLoadConfigurationFile(t *testing.T) {
	testCases := []struct {
		description string
		filepath    string
		expected    config.Config
		err         error
	}{
		{
			description: "not-toml-file",
			filepath:    "../../test/data/configloader/json.env.toml",
			err:         errors.New(""),
		},
		{
			description: "empty-file-empty-config",
			filepath:    "../../test/data/configloader/empty.env.toml",
		},
		{
			description: "full-file-full-config",
			filepath:    "../../.env.toml.example",
			expected: config.Config{
				Environment: config.Environment{Type: "LOCAL"},
				Database: config.Database{
					Host:     "127.0.0.1",
					Port:     "5432",
					User:     "user",
					Pass:     "password",
					Database: "goluca",
					SSLMode:  "disable",
				},
				HTTPServer: config.HTTPServer{
					Host: "localhost",
					Port: "3333",
				},
				GRPCServer: config.GRPCServer{
					Host: "localhost",
					Port: "5000",
				},
				NATS: config.NATS{
					Host: "localhost",
					Port: "4222",
					Wiretap: config.NATSWiretap{
						Enable: config.WiretapEnabled,
						Host:   "localhost",
						Port:   "4222",
					},
				},
			},
		},
		{
			description: "partial-file-partial-config",
			filepath:    "../../test/data/configloader/partial.env.toml",
			expected: config.Config{
				Environment: config.Environment{Type: "DEV"},
				HTTPServer: config.HTTPServer{
					Host: "localhost",
					Port: "3333",
				},
			},
		},
	}

	a := is.New(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			actual, err := loadConfigurationFile(tc.filepath)
			if tc.err != nil {
				a.True(err != nil)
				a.True(errors.As(err, &tc.err))
				return
			}
			a.NoErr(err)
			a.Equal(tc.expected, actual)
		})
	}
}

func TestLoadEnvironmentVariables(t *testing.T) {
	testCases := []struct {
		description string
		envVars     map[string]string
		expected    config.Config
	}{
		{
			description: "empty-vars",
			envVars: map[string]string{
				EnvType:        "",
				DBHost:         "",
				DBUser:         "",
				DBPass:         "",
				DBDatabase:     "",
				DBPort:         "",
				HTTPServerHost: "",
				HTTPServerPort: "",
				GRPCServerHost: "",
				GRPCServerPort: "",
				NATSHost:       "",
				NATSPort:       "",
				WiretapEnable:  "",
				WiretapHost:    "",
				WiretapPort:    "",
			},
		},
		{
			description: "filled-vars",
			envVars: map[string]string{
				EnvType:        "GOLUCA_ENVTYPE",
				DBHost:         "GOLUCA_DBHOST",
				DBPort:         "GOLUCA_DBPORT",
				DBUser:         "GOLUCA_DBUSER",
				DBPass:         "GOLUCA_DBPASS",
				DBDatabase:     "GOLUCA_DBDATABASE",
				HTTPServerHost: "GOLUCA_HTTP_SERVER_HOST",
				HTTPServerPort: "GOLUCA_HTTP_SERVER_PORT",
				GRPCServerHost: "GOLUCA_GRPC_SERVER_HOST",
				GRPCServerPort: "GOLUCA_GRPC_SERVER_PORT",
				NATSHost:       "GOLUCA_NATS_HOST",
				NATSPort:       "GOLUCA_NATS_PORT",
				WiretapEnable:  "enabled",
				WiretapHost:    "GOLUCA_WIRETAP_HOST",
				WiretapPort:    "GOLUCA_WIRETAP_PORT",
			},
			expected: config.Config{
				Environment: config.Environment{Type: "GOLUCA_ENVTYPE"},
				Database: config.Database{
					Host:     "GOLUCA_DBHOST",
					Port:     "GOLUCA_DBPORT",
					User:     "GOLUCA_DBUSER",
					Pass:     "GOLUCA_DBPASS",
					Database: "GOLUCA_DBDATABASE",
				},
				HTTPServer: config.HTTPServer{
					Host: "GOLUCA_HTTP_SERVER_HOST",
					Port: "GOLUCA_HTTP_SERVER_PORT",
				},
				GRPCServer: config.GRPCServer{
					Host: "GOLUCA_GRPC_SERVER_HOST",
					Port: "GOLUCA_GRPC_SERVER_PORT",
				},
				NATS: config.NATS{
					Host: "GOLUCA_NATS_HOST",
					Port: "GOLUCA_NATS_PORT",
					Wiretap: config.NATSWiretap{
						Enable: config.WiretapEnabled,
						Host:   "GOLUCA_WIRETAP_HOST",
						Port:   "GOLUCA_WIRETAP_PORT",
					},
				},
			},
		},
		{
			description: "partial-vars",
			envVars: map[string]string{
				EnvType:        "GOLUCA_ENVTYPE",
				HTTPServerHost: "GOLUCA_HTTP_SERVER_HOST",
				HTTPServerPort: "GOLUCA_HTTP_SERVER_PORT",
			},
			expected: config.Config{
				Environment: config.Environment{Type: "GOLUCA_ENVTYPE"},
				HTTPServer: config.HTTPServer{
					Host: "GOLUCA_HTTP_SERVER_HOST",
					Port: "GOLUCA_HTTP_SERVER_PORT",
				},
			},
		},
	}

	a := is.New(t)
	for i, tc := range testCases {
		// clean the
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			resetApplicationEnvironmentVariables()
			defer resetApplicationEnvironmentVariables()

			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			actual := loadEnvironmentVariables()
			a.Equal(tc.expected, actual)
		})

	}
}

func resetApplicationEnvironmentVariables() {
	for _, k := range EnvironmentVariableKeys {
		os.Unsetenv(k)
	}
}
