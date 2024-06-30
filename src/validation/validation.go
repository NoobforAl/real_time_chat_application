package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateStruct(structData any) []*fiber.Error {
	var errors []*fiber.Error
	err := validate.Struct(structData)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: err.Field() + " is " + err.Tag(),
			}
			errors = append(errors, &element)
		}
	}
	return errors
}
