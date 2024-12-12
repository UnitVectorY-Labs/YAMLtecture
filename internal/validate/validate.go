package validate

import (
	"fmt"

	"YAMLtecture/internal/configuration"
)

// ValidateConfig performs all required validations on the configuration.
func ValidateConfig(config *configuration.Config) error {
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

// isCyclic is a helper function to detect cycles in parent relationships.
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
