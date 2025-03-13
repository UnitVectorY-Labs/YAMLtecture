package mermaid

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/configuration"
)

func TestGenerateMermaid(t *testing.T) {
	err := filepath.Walk("../../tests", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			configPath := filepath.Join(path, "config.yaml")
			mermaidConfigPath := filepath.Join(path, "mermaid.yaml")
			mermaidPath := filepath.Join(path, "mermaid.mmd")

			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return nil
			}
			if _, err := os.Stat(mermaidConfigPath); os.IsNotExist(err) {
				return nil
			}
			if _, err := os.Stat(mermaidPath); os.IsNotExist(err) {
				return nil
			}

			relDir, err := filepath.Rel("../../tests", filepath.Dir(mermaidPath))
			if err != nil {
				return err
			}

			sanitizedRelDir := strings.ReplaceAll(relDir, string(filepath.Separator), "#")

			t.Run(sanitizedRelDir, func(t *testing.T) {
				config, err := configuration.LoadConfig(configPath)
				if err != nil {
					t.Fatalf("Failed to load config: %v", err)
				}

				mermaidConfig, err := LoadMermaid(mermaidConfigPath)
				if err != nil {
					t.Fatalf("Failed to load mermaid config: %v", err)
				}

				expectedBytes, err := os.ReadFile(mermaidPath)
				if err != nil {
					t.Fatalf("Failed to read Mermaid file: %v", err)
				}
				expectedOutput := string(expectedBytes)

				output, err := GenerateMermaid(config, mermaidConfig)
				if err != nil {
					t.Fatalf("GenerateMermaid returned error: %v", err)
				}
				if output != expectedOutput {
					t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
				}
			})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking through example folder: %v", err)
	}
}
