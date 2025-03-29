package configuration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInvalidConfig(t *testing.T) {
	configDir := "../../tests/invalid/config"

	entries, err := os.ReadDir(configDir)
	if err != nil {
		t.Fatalf("Error reading the invalid config directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(configDir, entry.Name())

			t.Run(path, func(t *testing.T) {
				// Verify the "input.yaml" and "expected_error.txt" files both exist
				inputFile := filepath.Join(path, "input.yaml")
				if _, err := os.Stat(inputFile); os.IsNotExist(err) {
					t.Fatalf("input.yaml file does not exist in %s", path)
				}

				expectedErrorFile := filepath.Join(path, "expected_error.txt")
				if _, err := os.Stat(expectedErrorFile); os.IsNotExist(err) {
					t.Fatalf("expected_error.txt file does not exist in %s", path)
				}

				// Load the configuration
				config, err := LoadConfig(inputFile)
				if err != nil {
					t.Fatalf("Failed to load %s: %v", inputFile, err)
				}

				// Validate the configuration
				err = config.Validate()
				if err == nil {
					t.Fatalf("Expected validation error for %s, but got none", inputFile)
				}

				actualErrorStr := "YAMLtecture\nError: Error validating configuration\n" + strings.TrimSpace(err.Error())

				// Read the expected error message
				expectedError, err := os.ReadFile(expectedErrorFile)
				if err != nil {
					t.Fatalf("Failed to read %s: %v", expectedErrorFile, err)
				}

				// Guard against nil error and trim whitespace from expected error
				expectedErrorStr := strings.TrimSpace(string(expectedError))

				// Check if the error message equals the expected error
				if actualErrorStr != expectedErrorStr {
					t.Errorf("Expected error message for %s: %q, but got: %q",
						inputFile, expectedErrorStr, actualErrorStr)
				}
			})
		}
	}
}
