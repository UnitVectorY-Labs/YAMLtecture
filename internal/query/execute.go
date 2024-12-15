package query

import (
	"YAMLtecture/internal/configuration"
	"fmt"
)

// ExecuteQuery filters the configuration based on the provided query.
// It returns a new configuration containing only the nodes and links
// that match the query conditions.
func ExecuteQuery(query *Query, config *configuration.Config) (configuration.Config, error) {

	// Prepare a new Config to hold the filtered results
	filteredConfig := configuration.Config{
		Nodes: []configuration.Node{},
		Links: []configuration.Link{},
	}

	// Iterate over all nodes and apply filters
	for _, node := range config.Nodes {
		matchesAllFilters, err := nodeMatchesAllFilters(node, query.Nodes.Filters)
		if err != nil {
			return configuration.Config{}, fmt.Errorf("error applying filters to node '%s': %w", node.ID, err)
		}
		if matchesAllFilters {
			filteredConfig.Nodes = append(filteredConfig.Nodes, node)
		}
	}

	// Iterate over all links and include only those where both source and target are in filtered nodes
	for _, rel := range config.Links {
		sourceExists := false
		targetExists := false
		for _, node := range filteredConfig.Nodes {
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
		if sourceExists && targetExists {
			filteredConfig.Links = append(filteredConfig.Links, rel)
		}
	}

	return filteredConfig, nil
}

// nodeMatchesAllFilters checks if a node satisfies all the provided filters.
func nodeMatchesAllFilters(node configuration.Node, filters []Filter) (bool, error) {
	for _, filter := range filters {
		matches, err := nodeMatchesFilter(node, filter)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}
	return true, nil
}

// nodeMatchesFilter checks if a node satisfies a single filter condition.
func nodeMatchesFilter(node configuration.Node, filter Filter) (bool, error) {
	// Extract the value of the specified field from the node
	fieldValue, err := getNodeFieldValue(node, filter.Condition.Field)
	if err != nil {
		return false, err
	}

	// Perform the comparison based on the operator
	switch filter.Condition.Operator {
	case "equals":
		return fieldValue == filter.Condition.Value, nil
	default:
		return false, fmt.Errorf("unsupported operator '%s'", filter.Condition.Operator)
	}
}

// getNodeFieldValue retrieves the value of a specified field from a node.
// It first checks the top-level fields, then the Attributes map.
func getNodeFieldValue(node configuration.Node, field string) (string, error) {
	switch field {
	case "id":
		return node.ID, nil
	case "type":
		return node.Type, nil
	case "parent":
		return node.Parent, nil
	default:
		// Check in Attributes
		if val, exists := node.Attributes[field]; exists {
			// Assuming all attribute values are strings for simplicity
			if strVal, ok := val.(string); ok {
				return strVal, nil
			}
			return "", fmt.Errorf("attribute '%s' is not a string", field)
		}
		return "", fmt.Errorf("field '%s' does not exist in node '%s'", field, node.ID)
	}
}
