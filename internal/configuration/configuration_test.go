package configuration

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadYAMLFiles(t *testing.T) {
	err := filepath.Walk("../../example", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "config.yaml" {

			relDir, err := filepath.Rel("../../example", filepath.Dir(path))
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
