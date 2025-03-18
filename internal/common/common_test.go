package common

import (
	"fmt"
	"testing"
)

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		field    string
		expected error
	}{
		{"empty_error", "", "Name", fmt.Errorf("Name cannot be empty")},
		{"valid_value", "validName", "Name", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsValidName(test.value, test.field)
			if err != nil && err.Error() != test.expected.Error() {
				t.Errorf("IsValidName(%q, %q) = %v; want %v", test.value, test.field, err, test.expected)
			}
			if err == nil && test.expected != nil {
				t.Errorf("IsValidName(%q, %q) = nil; want %v", test.value, test.field, test.expected)
			}
		})
	}
}
func TestIsValidValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		field    string
		expected error
	}{
		{"empty_error", "", "Value", fmt.Errorf("Value cannot be empty")},
		{"valid_value", "validValue", "Value", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsValidValue(test.value, test.field)
			if err != nil && err.Error() != test.expected.Error() {
				t.Errorf("IsValidValue(%q, %q) = %v; want %v", test.value, test.field, err, test.expected)
			}
			if err == nil && test.expected != nil {
				t.Errorf("IsValidValue(%q, %q) = nil; want %v", test.value, test.field, test.expected)
			}
		})
	}
}

func TestSanitizeLabel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
		{
			name:     "no special characters",
			input:    "abcdef",
			expected: "abcdef",
		},
		{
			name:     "single special character",
			input:    "a[b",
			expected: "ab",
		},
		{
			name:     "multiple special characters",
			input:    "a[b]c",
			expected: "abc",
		},
		{
			name:     "special characters at the beginning",
			input:    "[abc",
			expected: "abc",
		},
		{
			name:     "special characters at the end",
			input:    "abc]",
			expected: "abc",
		},
		{
			name:     "special characters in the middle",
			input:    "a(b)c",
			expected: "abc",
		},
		{
			name:     "special characters in the middle",
			input:    "a{b}c",
			expected: "abc",
		},
		{
			name:     "special characters in the middle",
			input:    "a<b>c",
			expected: "abc",
		},
		{
			name:     "special characters in the middle",
			input:    "abc",
			expected: "abc",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := SanitizeLabel(test.input)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
