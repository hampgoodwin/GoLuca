package config

// Environment holds the values for environment loaded configuration
type Environment struct {
	Type         string `validate:"required,oneof=DEV QA PROD"`
	DBDriverName string `validate:"required,eq=postgres"`
	DBConnString string `validate:"required"`
}
