package common

import (
	"fmt"
	"testing"
)

func TestIsValidColor(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		color    string
		expected error
	}{
		{"empty_color", "fill", "", nil},
		{"invalid_hex", "fill", "#12345Q", fmt.Errorf("invalid color for 'fill': '#12345Q'")},
		{"valid_hex", "fill", "#000000", nil},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			err := IsValidColor(test.field, test.color)
			if err != nil && err.Error() != test.expected.Error() {
				t.Errorf("IsValidColor(%q, %q) = %v; want %v", test.field, test.color, err, test.expected)
			}
			if err == nil && test.expected != nil {
				t.Errorf("IsValidColor(%q, %q) = nil; want %v", test.field, test.color, test.expected)
			}
		})
	}
}

func TestIsValidPixel(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		value    string
		expected error
	}{
		{"empty_pixel", "stroke-width", "", nil},
		{"invalid_pixel", "stroke-width", "1", fmt.Errorf("invalid pixel value for 'stroke-width': '1'")},
		{"invalid_pixel_bad_number", "stroke-width", "apx", fmt.Errorf("invalid pixel value for 'stroke-width': 'apx'")},
		{"invalid_pixel_no_number", "stroke-width", "px", fmt.Errorf("invalid pixel value for 'stroke-width': 'px'")},
		{"valid_pixel", "stroke-width", "1px", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsValidPixel(test.field, test.value)
			if err != nil && err.Error() != test.expected.Error() {
				t.Errorf("IsValidPixel(%q, %q) = %v; want %v", test.field, test.value, err, test.expected)
			}
			if err == nil && test.expected != nil {
				t.Errorf("IsValidPixel(%q, %q) = nil; want %v", test.field, test.value, test.expected)
			}
		})
	}
}
