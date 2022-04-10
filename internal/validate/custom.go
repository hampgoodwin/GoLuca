package validate

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

// int64 validates that a field is an int64 number
func int64(fl validator.FieldLevel) bool {
	_, err := strconv.ParseInt(fl.Field().String(), 10, 64)
	return err == nil
}
