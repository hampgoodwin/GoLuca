package validate

import (
	"github.com/go-playground/validator/v10"
)

// Validate wraps the logic for using go-playground/validator.
func Validate(i interface{}) error {
	if err := validator.New().Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil
		}
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			return err
		}
	}
	return nil
}
