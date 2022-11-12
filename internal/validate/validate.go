package validate

import (
	"github.com/go-playground/validator/v10"

	modelv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/model/v1"
	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
)

// Validate wraps the logic for using go-playground/validator.
func Validate(i interface{}) error {
	v := validator.New()
	registerCustomValidations(v)
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

func Var(i interface{}, tag string) error {
	v := validator.New()
	registerCustomValidations(v)
	if err := v.Var(i, tag); err != nil {
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

func registerCustomValidations(v *validator.Validate) {
	// custom validations
	_ = v.RegisterValidation("stringAsInt64", stringAsInt64)
	_ = v.RegisterValidation("KSUID", KSUID)

	// custom structs
	v.RegisterStructValidationMapRules(account, &modelv1.Account{})
	v.RegisterStructValidationMapRules(getAccountRequest, &servicev1.GetAccountRequest{})
	v.RegisterStructValidationMapRules(createAccountRequest, &servicev1.CreateAccountRequest{})

	v.RegisterStructValidationMapRules(transaction, &modelv1.Transaction{})
	v.RegisterStructValidationMapRules(getTransactionRequest, &servicev1.GetTransactionRequest{})

}
