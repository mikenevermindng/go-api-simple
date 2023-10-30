package customvalidator

import (
	"github.com/go-playground/validator/v10"
)

var CustomValidator *validator.Validate

func SetupValidator() {
	CustomValidator = validator.New()
	CustomValidator.RegisterValidation("password", PasswordValidator)
}
