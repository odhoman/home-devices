package common

import (
	"errors"
	"fmt"
	response "lambdas/common/response"

	"github.com/go-playground/validator/v10"
)

func ValidateDeviceRequestStruct(s interface{}) []string {
	validate := validator.New()
	var validationErrors []string
	if err := validate.Struct(s); err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, getMessageForFieldError(err.Tag(), err.Field()))
		}

	}
	return validationErrors
}

func getMessageForFieldError(tag, field string) string {
	switch field {
	case "MAC":
		if tag == "min" || tag == "max" {
			return "MAC address must be between 12 and 17 characters"
		}
	case "Name":
		if tag == "min" || tag == "max" {
			return "Name must be between 3 and 50 characters"
		}
	case "Type":
		if tag == "min" || tag == "max" {
			return "Type must be between 3 and 20 characters"
		}
	case "HomeID":
		if tag == "min" || tag == "max" {
			return "Home ID must be between 5 and 30 characters"
		}
	}

	return getDefaultValidationErrorMessage(tag, field)
}

func getDefaultValidationErrorMessage(tag, field string) string {
	return fmt.Sprintf("Validation failed for field '%s': %s", field, tag)
}

func ValidateAndResponseBadRequestErrors(s interface{}) map[string]interface{} {
	if errors := ValidateDeviceRequestStruct(s); len(errors) > 0 {
		return response.ReturnErrorResponse(errors, 400)
	}

	return map[string]interface{}{}
}

func CheckEmptyString(fieldName, value string) error {
	if value == "" {
		return errors.New(fmt.Sprintf("Field '%s' is empty. Please enter a value", fieldName))
	}
	return nil
}

func CheckEmptyStringAndResponseBadRequestErrors(fieldName, value string) map[string]interface{} {

	if error := CheckEmptyString(fieldName, value); error != nil {
		return response.ReturnErrorResponse([]string{error.Error()}, 400)
	}

	return nil
}
