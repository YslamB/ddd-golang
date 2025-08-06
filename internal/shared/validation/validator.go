package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {

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
