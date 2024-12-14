package mermaid

import (
	"os"
	"testing"

	"YAMLtecture/internal/configuration"
)

func TestGenerateMermaid(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		mermaidPath string
	}{
		{
			name:        "Single link",
			configPath:  "../../example/simple/architecture.yaml",
			mermaidPath: "../../example/simple/mermaid.mmd",
		},
		// ...add more test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := configuration.LoadConfig(tt.configPath)
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}

			expectedBytes, err := os.ReadFile(tt.mermaidPath)
			if err != nil {
				t.Fatalf("Failed to read Mermaid file: %v", err)
			}
			expectedOutput := string(expectedBytes)

			output, err := GenerateMermaid(config)
			if err != nil {
				t.Fatalf("GenerateMermaid returned error: %v", err)
			}
			if output != expectedOutput {
				t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
			}
		})
	}
}
