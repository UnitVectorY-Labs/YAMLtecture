package main

import (
	"flag"
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
	FilePaths     []string // To keep track of files read
}

var (
	validateFlag = flag.Bool("validate", false, "Validate the YAML architecture files")
	graphFlag    = flag.Bool("graph", false, "Generate a Mermaid diagram from the YAML architecture files")
	debugFlag    = flag.Bool("debug", false, "Enable debug output")
	dirFlag      = flag.String("dir", ".", "Directory containing YAML configuration files")
)

func main() {
	flag.Parse()

	// Ensure at least one command is provided
	if !*validateFlag && !*graphFlag {
		fmt.Println("Error: You must specify at least one command: --validate or --graph")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize configuration
	config := Config{
		Nodes:         make(map[string]Node),
		Relationships: []Relationship{},
		FilePaths:     []string{},
	}

	// Load and parse YAML files
	err := loadYAMLFiles(*dirFlag, &config)
	if err != nil {
		fmt.Printf("Error loading YAML files: %v\n", err)
		os.Exit(1)
	}

	// Perform validation
	err = validateConfig(&config)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

	if *debugFlag {
		fmt.Println("Validation successful.")
	}

	// Execute commands
	if *validateFlag && !*graphFlag {
		// Only validate
		fmt.Println("Validation passed.")
	}

	if *graphFlag {
		// Generate Mermaid diagram
		mermaid, err := generateMermaid(&config)
		if err != nil {
			fmt.Printf("Error generating Mermaid diagram: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(mermaid)
	}
}

// Add a wrapper struct to match the YAML structure
type YAMLFile struct {
	Architecture Architecture `yaml:"architecture"`
}

func loadYAMLFiles(dir string, config *Config) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Process only .yaml or .yml files
		if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
			if *debugFlag {
				fmt.Printf("Reading file: %s\n", path)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			var yf YAMLFile
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
				if *debugFlag {
					fmt.Printf("Added node: %s (Type: %s)\n", node.ID, node.Type)
				}
			}

			// Aggregate relationships
			for _, rel := range arch.Relationships {
				config.Relationships = append(config.Relationships, rel)
				if *debugFlag {
					fmt.Printf("Added relationship: %s -> %s (Type: %s)\n", rel.Source, rel.Target, rel.Type)
				}
			}

			// Track files read
			config.FilePaths = append(config.FilePaths, path)
		}
		return nil
	})
}

// validateConfig performs all required validations on the configuration
func validateConfig(config *Config) error {
	// Check that all nodes have unique IDs and IDs are present (already done during loading)

	// Validate parent relationships (acyclic)
	parentMap := make(map[string]string)
	for id, node := range config.Nodes {
		if node.Parent != "" {
			parentMap[id] = node.Parent
			// Check if parent exists
			if _, exists := config.Nodes[node.Parent]; !exists {
				return fmt.Errorf("node '%s' has non-existent parent '%s'", id, node.Parent)
			}
		}
	}

	// Detect cycles in parent relationships
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for id := range config.Nodes {
		if !visited[id] {
			if isCyclic(id, parentMap, visited, recStack) {
				return fmt.Errorf("cycle detected in parent relationships involving node '%s'", id)
			}
		}
	}

	// Validate relationships
	for _, rel := range config.Relationships {
		if _, exists := config.Nodes[rel.Source]; !exists {
			return fmt.Errorf("relationship has non-existent source node '%s'", rel.Source)
		}
		if _, exists := config.Nodes[rel.Target]; !exists {
			return fmt.Errorf("relationship has non-existent target node '%s'", rel.Target)
		}
	}

	return nil
}

// isCyclic is a helper function to detect cycles in parent relationships
func isCyclic(node string, parentMap map[string]string, visited, recStack map[string]bool) bool {
	visited[node] = true
	recStack[node] = true

	parent, hasParent := parentMap[node]
	if hasParent {
		if !visited[parent] {
			if isCyclic(parent, parentMap, visited, recStack) {
				return true
			}
		} else if recStack[parent] {
			return true
		}
	}

	recStack[node] = false
	return false
}

// generateMermaid creates a Mermaid diagram based on the relationships and includes all nodes
func generateMermaid(config *Config) (string, error) {
	var sb strings.Builder
	sb.WriteString("graph LR\n")

	// To keep track of nodes that have been connected via relationships
	connectedNodes := make(map[string]bool)

	for _, rel := range config.Relationships {
		source := rel.Source
		target := rel.Target
		label := rel.Type

		// Mermaid syntax: source -->|label| target
		line := fmt.Sprintf("  %s -->|%s| %s\n", source, label, target)
		sb.WriteString(line)

		connectedNodes[source] = true
		connectedNodes[target] = true
	}

	// Add nodes without relationships
	for id, node := range config.Nodes {
		if !connectedNodes[id] {
			// Represent as a standalone node
			line := fmt.Sprintf("  %s[%s]\n", id, node.Type)
			sb.WriteString(line)
		}
	}

	return sb.String(), nil
}
