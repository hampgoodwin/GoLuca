package validate

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
)

// int64 validates that a field is an int64 number
func stringAsInt64(fl validator.FieldLevel) bool {
	_, err := strconv.ParseInt(fl.Field().String(), 10, 64)
	return err == nil
}

// KSUID does a simple regex string comparison for char type
// and length
func KSUID(fl validator.FieldLevel) bool {
	_, err := ksuid.Parse(fl.Field().String())
	return err == nil
}
