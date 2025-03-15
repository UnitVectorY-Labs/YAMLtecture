package common

import (
	"fmt"
	"strings"
)

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

// SanitizeLabel escapes reserved characters in a Mermaid flowchart node label.
// It prefixes each reserved character with a backslash.
func SanitizeLabel(label string) string {
	// Define the set of reserved characters.
	reserved := map[rune]bool{
		'[': true,
		']': true,
		'(': true,
		')': true,
		'{': true,
		'}': true,
		'<': true,
		'>': true,
		'"': true,
	}

	var builder strings.Builder

	for _, char := range label {
		if !reserved[char] {
			// Only include characters that aren't reserved
			builder.WriteRune(char)
		}
	}
	return builder.String()
}
