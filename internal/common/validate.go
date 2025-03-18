package common

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Create a singleton validator instance
var validate = validator.New()

func init() {
	// Register the custom validator function
	validate.RegisterValidation("pixel", isPixelValue)
}

// Custom validator function
func isPixelValue(fl validator.FieldLevel) bool {
	// Regular expression to match a number followed by "px"
	re := regexp.MustCompile(`^\d+px$`)
	return re.MatchString(fl.Field().String())
}

func IsValidColor(field string, color string) error {
	if color == "" {
		return nil
	}

	// Validate the color is valid
	err := validate.Var(color, "hexcolor")
	if err != nil {
		return fmt.Errorf("invalid color for %s: %s", field, color)
	}

	return nil
}

func IsValidPixel(field string, value string) error {
	if value == "" {
		return nil
	}

	// Validate the stroke width is valid integer suffixed with 'px'
	err := validate.Var(value, "pixel")
	if err != nil {
		return fmt.Errorf("invalid pixel for %s: %s", field, value)
	}

	return nil
}
