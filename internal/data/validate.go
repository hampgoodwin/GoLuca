package data

import (
	"github.com/go-playground/validator"
)

// Validate is a DRY function which wraps the logic for using go-playground/validator.
func Validate(i interface{}) error {
	if err := validator.New().Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			return err
		}
	}
	return nil
}
