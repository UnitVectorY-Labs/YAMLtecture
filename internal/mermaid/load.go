package mermaid

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseYAML parses the YAML content into a Mermaid
func ParseYAML(content string) (*Mermaid, error) {
	var config Mermaid
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	// Specify the default values if they were not provided

	if config.Direction == "" {
		config.Direction = "TD"
	}

	return &config, nil
}

// LoadMermaid loads and parses a single YAML mermaid setting file from the given path.
func LoadMermaid(filePath string) (*Mermaid, error) {

	// Read the file contents to a string
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Parse the YAML
	return ParseYAML(string(data))
}
