package common

import (
	"fmt"
	"testing"
)

func TestIsValidName(t *testing.T) {
	tests := []struct {
		value    string
		field    string
		expected error
	}{
		{"", "Name", fmt.Errorf("Name cannot be empty")},
		{"validName", "Name", nil},
	}

	for _, test := range tests {
		err := IsValidName(test.value, test.field)
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("IsValidName(%q, %q) = %v; want %v", test.value, test.field, err, test.expected)
		}
		if err == nil && test.expected != nil {
			t.Errorf("IsValidName(%q, %q) = nil; want %v", test.value, test.field, test.expected)
		}
	}
}
func TestIsValidValue(t *testing.T) {
	tests := []struct {
		value    string
		field    string
		expected error
	}{
		{"", "Value", fmt.Errorf("Value cannot be empty")},
		{"validValue", "Value", nil},
	}

	for _, test := range tests {
		err := IsValidValue(test.value, test.field)
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("IsValidValue(%q, %q) = %v; want %v", test.value, test.field, err, test.expected)
		}
		if err == nil && test.expected != nil {
			t.Errorf("IsValidValue(%q, %q) = nil; want %v", test.value, test.field, test.expected)
		}
	}
}
