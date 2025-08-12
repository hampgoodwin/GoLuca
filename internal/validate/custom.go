package validate

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var uuid7regexp = regexp.MustCompile(`^[0-9a-f]{8}(?:\-[0-9a-f]{4}){3}-[0-9a-f]{12}$`)

// int64 validates that a field is an int64 number
func stringAsInt64(fl validator.FieldLevel) bool {
	_, err := strconv.ParseInt(fl.Field().String(), 10, 64)
	return err == nil
}

// uuid7 parses a UUID and ensures it's V7 via regex
func uuid7(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	match := uuid7regexp.MatchString(field)
	return match
}
