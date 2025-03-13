package mermaid

import "strings"

// SanitizeLabel escapes reserved characters in a Mermaid flowchart node label.
// It prefixes each reserved character with a backslash.
func sanitizeLabel(label string) string {
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
