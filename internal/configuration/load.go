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
		Nodes: make(map[string]Node),
		Links: []Link{},
	}

	// Define a temporary struct to unmarshal the YAML
	temp := struct {
		Nodes []Node `yaml:"nodes"`
		Links []Link `yaml:"links"`
	}{}

	// Unmarshal the YAML into the temporary struct
	err := yaml.Unmarshal([]byte(config), &temp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	// Convert the temporary struct to the Config struct
	for _, node := range temp.Nodes {
		c.Nodes[node.ID] = node
	}
	c.Links = temp.Links

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
		Nodes: make(map[string]Node),
		Links: []Link{},
	}

	for _, config := range configs {
		// Merge nodes
		for id, node := range config.Nodes {
			if _, exists := merged.Nodes[id]; exists {
				return nil, fmt.Errorf("duplicate node ID '%s' found", id)
			}
			merged.Nodes[id] = node
		}

		// Merge links
		merged.Links = append(merged.Links, config.Links...)
	}

	return merged, nil
}
