package api

import (
	"simple_bank/util"

	"github.com/go-playground/validator/v10"
)

// validCurrency is a custom validator function to check if a currency is valid.
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return util.IsValidCurrency(currency)
}
