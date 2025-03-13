package mermaid

import "testing"

// Test the sanitizeLabel function

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
			actual := sanitizeLabel(test.input)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
