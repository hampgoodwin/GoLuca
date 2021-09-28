package configloader

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-playground/validator"
	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/stretchr/testify/require"
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
			description: "no-file-empty-vars-empty-config-err-reading-file",
			filepath:    "",
			err:         &os.PathError{},
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
				EnvType:    "DEV",
				DBHost:     "127.0.0.1",
				DBUser:     "user",
				DBPass:     "password",
				DBDatabase: "goluca",
				DBPort:     "5432",
				APIHost:    "localhost",
				APIPort:    "3333",
			},
		},
		{
			description: "full-file-full-vars-empty-config-overwrite-file",
			filepath:    "../../.env.toml.example",
			envVars: map[string]string{
				EnvType:    "DEV",
				DBHost:     "GOLUCA_DBHOST",
				DBPort:     "GOLUCA_DBPORT",
				DBUser:     "GOLUCA_DBUSER",
				DBPass:     "GOLUCA_DBPASS",
				DBDatabase: "GOLUCA_DBDATABASE",
				APIHost:    "GOLUCA_APIHOST",
				APIPort:    "GOLUCA_APIPORT",
			},
			expected: config.Config{
				EnvType:    "DEV",
				DBHost:     "GOLUCA_DBHOST",
				DBPort:     "GOLUCA_DBPORT",
				DBUser:     "GOLUCA_DBUSER",
				DBPass:     "GOLUCA_DBPASS",
				DBDatabase: "GOLUCA_DBDATABASE",
				APIHost:    "GOLUCA_APIHOST",
				APIPort:    "GOLUCA_APIPORT",
			},
		},
		{
			description: "full-file-partial-vars-empty-config-overwrite-file",
			filepath:    "../../.env.toml.example",
			envVars: map[string]string{
				EnvType: "DEV",
				APIHost: "GOLUCA_APIHOST",
				APIPort: "GOLUCA_APIPORT",
			},
			expected: config.Config{
				EnvType:    "DEV",
				DBHost:     "127.0.0.1",
				DBUser:     "user",
				DBPass:     "password",
				DBDatabase: "goluca",
				DBPort:     "5432",
				APIHost:    "GOLUCA_APIHOST",
				APIPort:    "GOLUCA_APIPORT",
			},
		},
		{
			description: "full-file-partial-vars-partial-config-overwrite-file",
			cfg:         config.Config{EnvType: "DEV"},
			filepath:    "../../.env.toml.example",
			envVars: map[string]string{
				EnvType: "LOCAL",
				APIHost: "GOLUCA_APIHOST",
				APIPort: "GOLUCA_APIPORT",
			},
			expected: config.Config{
				EnvType:    "LOCAL",
				DBHost:     "127.0.0.1",
				DBUser:     "user",
				DBPass:     "password",
				DBDatabase: "goluca",
				DBPort:     "5432",
				APIHost:    "GOLUCA_APIHOST",
				APIPort:    "GOLUCA_APIPORT",
			},
		},
		{
			description: "full-file-partial-vars-partial-config-persist-merge",
			cfg:         config.Config{EnvType: "DEV"},
			filepath:    "../../.env.toml.example",
			envVars: map[string]string{
				APIHost: "GOLUCA_APIHOST",
				APIPort: "GOLUCA_APIPORT",
			},
			expected: config.Config{
				EnvType:    "DEV",
				DBHost:     "127.0.0.1",
				DBUser:     "user",
				DBPass:     "password",
				DBDatabase: "goluca",
				DBPort:     "5432",
				APIHost:    "GOLUCA_APIHOST",
				APIPort:    "GOLUCA_APIPORT",
			},
		},
	}

	a := require.New(t)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			resetApplicationEnvironmentVariables()
			defer resetApplicationEnvironmentVariables()

			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			actual, err := Load(tc.cfg, tc.filepath)
			if tc.err != nil {
				a.NotNil(err)
				// Using errors.As because it detects validator.ValidationErrors
				a.True(errors.As(err, &tc.err))
				return
			}
			a.NoError(err)
			a.Equal(tc.expected, actual)
		})
	}
}

func TestLoadConfigurationFile(t *testing.T) {
	testCases := []struct {
		description string
		filepath    string
		expected    config.Config
		err         bool
	}{
		{
			description: "not-toml-file",
			filepath:    "../../test/data/configloader/json.env.toml",
			err:         true,
		},
		// {
		// 	description: "empty-file-empty-config",
		// 	filepath:    "../../test/data/configloader/empty.env.toml",
		// },
		// {
		// 	description: "full-file-full-config",
		// 	filepath:    "../../.env.toml.example",
		// 	expected: config.Config{
		// 		EnvType:    "DEV",
		// 		DBHost:     "127.0.0.1",
		// 		DBUser:     "user",
		// 		DBPass:     "password",
		// 		DBDatabase: "goluca",
		// 		DBPort:     "5432",
		// 		APIHost:    "localhost",
		// 		APIPort:    "3333",
		// 	},
		// },
		// {
		// 	description: "partial-file-partial-config",
		// 	filepath:    "../../test/data/configloader/partial.env.toml",
		// 	expected: config.Config{
		// 		EnvType: "DEV",
		// 		APIHost: "localhost",
		// 		APIPort: "3333",
		// 	},
		// },
	}

	a := require.New(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			actual, err := loadConfigurationFile(tc.filepath)
			if tc.err {
				a.Error(err)
				return
			}
			a.NoError(err)
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
				EnvType:    "",
				DBHost:     "",
				DBUser:     "",
				DBPass:     "",
				DBDatabase: "",
				DBPort:     "",
				APIHost:    "",
				APIPort:    "",
			},
		},
		{
			description: "filled-vars",
			envVars: map[string]string{
				EnvType:    "GOLUCA_ENVTYPE",
				DBHost:     "GOLUCA_DBHOST",
				DBPort:     "GOLUCA_DBPORT",
				DBUser:     "GOLUCA_DBUSER",
				DBPass:     "GOLUCA_DBPASS",
				DBDatabase: "GOLUCA_DBDATABASE",
				APIHost:    "GOLUCA_APIHOST",
				APIPort:    "GOLUCA_APIPORT",
			},
			expected: config.Config{
				EnvType:    "GOLUCA_ENVTYPE",
				DBHost:     "GOLUCA_DBHOST",
				DBPort:     "GOLUCA_DBPORT",
				DBUser:     "GOLUCA_DBUSER",
				DBPass:     "GOLUCA_DBPASS",
				DBDatabase: "GOLUCA_DBDATABASE",
				APIHost:    "GOLUCA_APIHOST",
				APIPort:    "GOLUCA_APIPORT",
			},
		},
		{
			description: "partial-vars",
			envVars: map[string]string{
				EnvType: "GOLUCA_ENVTYPE",
				APIHost: "GOLUCA_APIHOST",
				APIPort: "GOLUCA_APIPORT",
			},
			expected: config.Config{
				EnvType:    "GOLUCA_ENVTYPE",
				DBHost:     "",
				DBPort:     "",
				DBUser:     "",
				DBPass:     "",
				DBDatabase: "",
				APIHost:    "GOLUCA_APIHOST",
				APIPort:    "GOLUCA_APIPORT",
			},
		},
	}

	a := require.New(t)
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
