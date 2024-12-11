package validation

import (
	"fmt"
	"go-fiber-api-starter/internal/model"
	"log"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateUser(user *model.User) []string {
	err := validate.Struct(user)
	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Fatal("InvalidValidationError: Your code may have produced an invalid value for validation such as an interface with nil value.")
			return []string{err.(validator.ValidationErrors).Error()}
		}

		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, (err.Field() + " is invalid"))
		}

		return errors
	}
	return nil
}

func ValidateSomething(something *model.Something) error {
	err := validate.Struct(something)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value. most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println("ERROR:", err)
			return err.(validator.ValidationErrors)
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println("NAMESPACE:", err.Namespace())
			fmt.Println("FIELD:", err.Field())
			fmt.Println("STRUCTNAMESPACE:", err.StructNamespace())
			fmt.Println("STRUCTFIELD:", err.StructField())
			fmt.Println("TAG:", err.Tag())
			fmt.Println("ACTUALTAG:", err.ActualTag())
			fmt.Println("KIND:", err.Kind())
			fmt.Println("TYPE:", err.Type())
			fmt.Println("VALUE:", err.Value())
			fmt.Println("PARAM:", err.Param())
			fmt.Println()
		}

		return err.(validator.ValidationErrors)
	}
	return nil
}
