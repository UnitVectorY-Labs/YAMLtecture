package configuration

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/UnitVectorY-Labs/yamlequal"
)

func TestLoadYAMLFiles(t *testing.T) {
	err := filepath.Walk("../../tests", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "configs" && info.IsDir() {

			relDir, err := filepath.Rel("../../tests", filepath.Dir(path))
			if err != nil {
				return err
			}

			sanitizedRelDir := strings.ReplaceAll(relDir, string(filepath.Separator), "#") + "#configs"

			// Log the paths for debugging
			log.Printf("Relative Directory: %s", relDir)
			log.Printf("Sanitized Relative Directory: %s", sanitizedRelDir)

			// Get the parent directory of 'path'
			parentDir := filepath.Dir(path)
			configPath := filepath.Join(parentDir, "config.yaml")

			// If config path does not exist, skip this iteration
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				log.Printf("Skipping %s as config.yaml does not exist", sanitizedRelDir)
				return nil
			}

			t.Run(sanitizedRelDir, func(t *testing.T) {

				folderConfig, err := LoadFolder(path)
				if err != nil {
					t.Errorf("Failed to load folder %s: %v", path, err)
					return
				}

				folderConfigBytes := []byte(folderConfig.YamlString())

				config, err := LoadConfig(configPath)
				if err != nil {
					t.Errorf("Failed to load %s: %v", path, err)
					return
				}

				configBytes := []byte(config.YamlString())

				// Compare two YAML content strings directly
				equal, diff, err := yamlequal.CompareYAML(folderConfigBytes, configBytes)
				if err != nil {
					t.Error("Error comparing files.", err)
					return
				}

				if !equal {
					t.Errorf("Files are not the same:\n%s", diff)
					return
				}
			})
		} else if info.Name() == "config.yaml" {

			relDir, err := filepath.Rel("../../tests", filepath.Dir(path))
			if err != nil {
				return err
			}

			sanitizedRelDir := strings.ReplaceAll(relDir, string(filepath.Separator), "#")

			// Log the paths for debugging
			log.Printf("Relative Directory: %s", relDir)
			log.Printf("Sanitized Relative Directory: %s", sanitizedRelDir)

			t.Run(sanitizedRelDir, func(t *testing.T) {
				config, err := LoadConfig(path)
				if err != nil {
					t.Errorf("Failed to load %s: %v", path, err)
					return
				}

				// Validate the configuration
				err = config.Validate()
				if err != nil {
					t.Errorf("Failed to validate %s: %v", path, err)
					return
				}
			})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking the example directory: %v", err)
	}
}
