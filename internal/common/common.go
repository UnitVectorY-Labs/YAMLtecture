package common

import "fmt"

// IsValidName checks if the value is a valid name for YAMLtecture conventions.
func IsValidName(value string, field string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", field)

		// TODO: Validate character set to apply restrictions on the name
	} else {
		return nil
	}
}

// IsValidValue checks if the value is a valid attribute value for YAMLtecture conventions.
func IsValidValue(value string, field string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", field)

		// TODO: Less strict validation will be applied to attribute values
	} else {
		return nil
	}
}
