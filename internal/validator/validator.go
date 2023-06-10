package validator

import (
	"game-server-selector/internal/models"

	"github.com/go-playground/validator"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateServer(user models.CreateServerRequest) []*ErrorResponse {
	var validate = validator.New()
	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateServerUpdate(user models.UpdateServerRequest) []*ErrorResponse {
	var validate = validator.New()
	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
