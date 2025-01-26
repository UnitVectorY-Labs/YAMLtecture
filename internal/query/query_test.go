package query

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/configuration"
)

func TestGenerateQuery(t *testing.T) {
	err := filepath.Walk("../../example", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			configPath := filepath.Join(path, "config.yaml")
			queryPath := filepath.Join(path, "query.yaml")
			configInputPath := filepath.Join(path, "../../config.yaml")

			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return nil
			}
			if _, err := os.Stat(queryPath); os.IsNotExist(err) {
				return nil
			}
			if _, err := os.Stat(configInputPath); os.IsNotExist(err) {
				return nil
			}

			relDir, err := filepath.Rel("../../example", filepath.Dir(queryPath))
			if err != nil {
				return err
			}

			sanitizedRelDir := strings.ReplaceAll(relDir, string(filepath.Separator), "#")

			t.Run(sanitizedRelDir, func(t *testing.T) {
				config, err := configuration.LoadConfig(configInputPath)
				if err != nil {
					t.Fatalf("Failed to load config: %v", err)
				}

				err = config.Validate()
				if err != nil {
					t.Fatalf("Failed to validate config: %v", err)
				}

				query, err := LoadQuery(queryPath)
				if err != nil {
					t.Fatalf("Failed to load query: %v", err)
				}

				output, err := ExecuteQuery(query, config)
				if err != nil {
					t.Fatalf("ExecuteQuery returned error: %v", err)
				}

				err = output.Validate()
				if err != nil {
					t.Fatalf("Failed to validate output: %v", err)
				}

				outputYaml := output.YamlString()

				expectedBytes, err := os.ReadFile(configPath)
				if err != nil {
					t.Fatalf("Failed to read expected output: %v", err)
				}
				expectedOutput := string(expectedBytes)

				if outputYaml != expectedOutput {
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
