package configuration

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseYAML parses a YAML configuration string into a Config struct.
func ParseYAML(config string) (*Config, error) {
	// Initialize an empty configuration
	c := &Config{
		Nodes: []Node{},
		Links: []Link{},
	}

	// Unmarshal the YAML into the Config struct
	err := yaml.Unmarshal([]byte(config), c)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	return c, nil
}

// LoadConfig loads and parses a single YAML configuration file from the given path.
func LoadConfig(filePath string) (*Config, error) {

	// Read the file contents to a string
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Parse the YAML
	return ParseYAML(string(data))
}

func MergeConfigs(configs ...*Config) (*Config, error) {
	merged := &Config{
		Nodes: []Node{},
		Links: []Link{},
	}

	nodeMap := make(map[string]Node)
	for _, config := range configs {
		// Merge nodes
		for _, node := range config.Nodes {
			if _, exists := nodeMap[node.ID]; exists {
				return nil, fmt.Errorf("duplicate node ID '%s' found", node.ID)
			}
			nodeMap[node.ID] = node
			merged.Nodes = append(merged.Nodes, node)
		}

		// Merge links
		merged.Links = append(merged.Links, config.Links...)
	}

	return merged, nil
}
