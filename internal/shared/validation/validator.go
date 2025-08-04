package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// validate is the global validator instance.
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct based on its 'validate' tags.
// It returns an error if validation fails, detailing the first error.
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		// This will return the first validation error encountered.
		// For more detailed error messages, you might iterate over validator.ValidationErrors.
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			if len(validationErrors) > 0 {
				firstError := validationErrors[0]
				return fmt.Errorf("validation failed for field '%s': %s", firstError.Field(), firstError.Tag())
			}
		}
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
