package configuration

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Node represents a system component
type Node struct {
	ID         string                 `yaml:"id"`
	Type       string                 `yaml:"type"`
	Parent     string                 `yaml:"parent,omitempty"`
	Attributes map[string]interface{} `yaml:"attributes,omitempty"`
}

// Link represents an interaction between nodes
type Link struct {
	Source     string                 `yaml:"source"`
	Target     string                 `yaml:"target"`
	Type       string                 `yaml:"type"`
	Attributes map[string]interface{} `yaml:"attributes,omitempty"`
}

// Config holds the aggregated architecture
type Config struct {
	Nodes map[string]Node `yaml:"nodes"`
	Links []Link          `yaml:"links"`
}

// YamlString returns the YAML representation of the configuration
func (c *Config) YamlString() string {
	// Marshall the config to a string
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("error marshalling config: %v", err)
	}
	return string(data)
}
