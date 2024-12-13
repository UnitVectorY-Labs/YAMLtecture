package configuration

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Node represents a system component
type Node struct {
	ID         string                 `yaml:"id"`
	Type       string                 `yaml:"type"`
	Parent     string                 `yaml:"parent,omitempty"`
	Attributes map[string]interface{} `yaml:"attributes,omitempty"`
}

// Relationship represents an interaction between nodes
type Relationship struct {
	Source     string                 `yaml:"source"`
	Target     string                 `yaml:"target"`
	Type       string                 `yaml:"type"`
	Attributes map[string]interface{} `yaml:"attributes,omitempty"`
}

// Architecture holds all nodes and relationships
type Architecture struct {
	Nodes         []Node         `yaml:"nodes"`
	Relationships []Relationship `yaml:"relationships"`
	Includes      []string       `yaml:"includes,omitempty"` // For future modularity
}

// Config holds the aggregated architecture
type Config struct {
	Nodes         map[string]Node
	Relationships []Relationship
}

// Print out the configuration method for the Config struct as the YAML
// representation of the architecture.
func (c *Config) YamlString() string {
	// Marshall the config to a string
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("error marshalling config: %v", err)
	}
	return string(data)
}

// LoadYAMLFiles loads and parses YAML configuration files from the specified directory.
func LoadYAMLFiles(dir string) (*Config, error) {
	config := &Config{
		Nodes:         make(map[string]Node),
		Relationships: []Relationship{},
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Process only .yaml or .yml files
		if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			var yf struct {
				Architecture Architecture `yaml:"architecture"`
			}
			err = yaml.Unmarshal(data, &yf)
			if err != nil {
				return fmt.Errorf("failed to parse YAML file %s: %w", path, err)
			}

			arch := yf.Architecture

			// Aggregate nodes
			for _, node := range arch.Nodes {
				if node.ID == "" {
					return fmt.Errorf("node in file %s is missing 'id'", path)
				}
				if _, exists := config.Nodes[node.ID]; exists {
					return fmt.Errorf("duplicate node ID '%s' found in file %s", node.ID, path)
				}
				config.Nodes[node.ID] = node
			}

			// Aggregate relationships
			config.Relationships = append(config.Relationships, arch.Relationships...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return config, nil
}
