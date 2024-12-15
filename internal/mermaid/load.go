package mermaid

import (
	"fmt"

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
