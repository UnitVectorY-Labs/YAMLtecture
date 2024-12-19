package configuration

import (
	"fmt"
	"os"
	"path/filepath"

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

func LoadFolder(folderPath string) (*Config, error) {
	// Verify the folderPath is a folder
	fileInfo, err := os.Stat(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading folder: %v", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("folderPath is not a directory")
	}

	// Loop through the folder contents loading in all .yaml files with LoadConfig
	configs := []*Config{}
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading folder: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".yaml" {
			config, err := LoadConfig(filepath.Join(folderPath, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("error loading config file: %v", err)
			}
			configs = append(configs, config)
		}
	}

	// Merge the loaded configs
	return MergeConfigs(configs...)
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
