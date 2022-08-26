package validate

import (
	"github.com/go-playground/validator/v10"
)

// Validate wraps the logic for using go-playground/validator.
func Validate(i interface{}) error {
	v := validator.New()
	registerCustomerFunctions(v)
	if err := v.Struct(i); err != nil {
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

func registerCustomerFunctions(v *validator.Validate) {
	_ = v.RegisterValidation("int64", int64)
	_ = v.RegisterValidation("KSUID", KSUID)
}
