package configuration

import (
	"fmt"

	"YAMLtecture/internal/common"
)

// ValidateConfig performs all required validations on the configuration.
func (config *Config) Validate() error {

	// Validate nodes
	for _, node := range config.Nodes {
		err := node.validate()
		if err != nil {
			return fmt.Errorf("node '%s' is invalid: %w", node.ID, err)
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
	for _, node := range config.Nodes {
		if node.Parent != "" {
			parentMap[node.ID] = node.Parent
			// Check if parent exists
			parentExists := false
			for _, n := range config.Nodes {
				if n.ID == node.Parent {
					parentExists = true
					break
				}
			}
			if !parentExists {
				return fmt.Errorf("node '%s' has non-existent parent '%s'", node.ID, node.Parent)
			}
		}
	}

	// Detect cycles in parent links
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, node := range config.Nodes {
		if !visited[node.ID] {
			if isCyclic(node.ID, parentMap, visited, recStack) {
				return fmt.Errorf("cycle detected in parent links involving node '%s'", node.ID)
			}
		}
	}

	// Validate links
	for _, rel := range config.Links {
		sourceExists := false
		targetExists := false
		for _, node := range config.Nodes {
			if node.ID == rel.Source {
				sourceExists = true
			}
			if node.ID == rel.Target {
				targetExists = true
			}
			if sourceExists && targetExists {
				break
			}
		}
		if !sourceExists {
			return fmt.Errorf("link has non-existent source node '%s'", rel.Source)
		}
		if !targetExists {
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
