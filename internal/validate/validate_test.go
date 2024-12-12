package validate

import (
	"testing"

	"YAMLtecture/internal/configuration"
)

func TestValidateConfigFiles(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		expectError bool
	}{
		{
			name:        "Valid configuration",
			configPath:  "../../example/simple/",
			expectError: false,
		},
		// ...add more test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := configuration.LoadYAMLFiles(tt.configPath)
			if err != nil {
				t.Fatalf("Failed to load YAML files: %v", err)
			}
			err = ValidateConfig(config)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateConfig() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
