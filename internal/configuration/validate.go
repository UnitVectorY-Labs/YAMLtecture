package configuration

import (
	"fmt"

	"YAMLtecture/internal/common"
)

// ValidateConfig performs all required validations on the configuration.
func (config *Config) Validate() error {

	// Validate nodes
	for id, node := range config.Nodes {
		err := node.validate()
		if err != nil {
			return fmt.Errorf("node '%s' is invalid: %w", id, err)
		}
	}

	// Validate links
	for i, rel := range config.Links {
		err := rel.validate()
		if err != nil {
			return fmt.Errorf("link at index %d is invalid: %w", i, err)
		}
	}

	// Validate parent links (acyclic)
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

	// Detect cycles in parent links
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for id := range config.Nodes {
		if !visited[id] {
			if isCyclic(id, parentMap, visited, recStack) {
				return fmt.Errorf("cycle detected in parent links involving node '%s'", id)
			}
		}
	}

	// Validate links
	for _, rel := range config.Links {
		if _, exists := config.Nodes[rel.Source]; !exists {
			return fmt.Errorf("link has non-existent source node '%s'", rel.Source)
		}
		if _, exists := config.Nodes[rel.Target]; !exists {
			return fmt.Errorf("link has non-existent target node '%s'", rel.Target)
		}
	}

	return nil
}

func (node *Node) validate() error {
	var err error

	err = common.IsValidName(node.ID, "node.id")
	if err != nil {
		return err
	}

	err = common.IsValidName(node.Type, "node.type")
	if err != nil {
		return err
	}

	if node.Parent != "" {
		err = common.IsValidName(node.Parent, "node.parent")
		if err != nil {
			return err
		}
	}

	// validate the attribute keys
	for key, value := range node.Attributes {
		err = common.IsValidName(key, "attribute.key")
		if err != nil {
			return err
		}

		err = common.IsValidValue(value.(string), "attribute.value")
		if err != nil {
			return err
		}
	}

	return nil
}

func (rel *Link) validate() error {
	var err error

	err = common.IsValidName(rel.Source, "link.source")
	if err != nil {
		return err
	}

	err = common.IsValidName(rel.Target, "link.target")
	if err != nil {
		return err
	}

	err = common.IsValidName(rel.Type, "link.type")
	if err != nil {
		return err
	}

	// validate the attribute keys
	for key, value := range rel.Attributes {
		err = common.IsValidName(key, "attribute.key")
		if err != nil {
			return err
		}

		err = common.IsValidValue(value.(string), "attribute.value")
		if err != nil {
			return err
		}
	}

	return nil
}

// isCyclic is a helper function to detect cycles in parent links.
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
