package validator

import (
	"fmt"
	"strings"

	playvalidator "github.com/go-playground/validator/v10"
)

var validate = playvalidator.New()

// ValidationError represents a field-level validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate validates a struct and returns a list of field errors
func Validate(s interface{}) []ValidationError {
	var errs []ValidationError
	if err := validate.Struct(s); err != nil {
		for _, e := range err.(playvalidator.ValidationErrors) {
			errs = append(errs, ValidationError{
				Field:   toSnakeCase(e.Field()),
				Message: msgFor(e),
			})
		}
	}
	return errs
}

func msgFor(e playvalidator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", e.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", e.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", e.Param())
	default:
		return fmt.Sprintf("Validation failed on %s", e.Tag())
	}
}

// toSnakeCase converts FieldName to field_name
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
