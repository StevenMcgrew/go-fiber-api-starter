package validation

import (
	"go-fiber-api-starter/internal/model"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateUser(user *model.User) error {
	err := validate.Struct(user)
	if err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

func ValidateSomething(something *model.Something) error {
	err := validate.Struct(something)
	if err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
