package query

import (
	"fmt"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/configuration"
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

	// Create a map lf all of the node IDs so we can easily check if a node exists
	nodeIDs := make(map[string]bool)
	for _, node := range filteredConfig.Nodes {
		nodeIDs[node.ID] = true
	}

	// Iterate over all links and include only those where both source and target are in filtered nodes
	for _, link := range config.Links {
		matchesAllFilters, err := linkMatchesAllFilters(link, query.Links.Filters)
		if err != nil {
			return configuration.Config{}, fmt.Errorf("error applying filters to link '%s' -> '%s': %w", link.Source, link.Target, err)
		}
		if !matchesAllFilters {
			continue
		}

		if nodeIDs[link.Source] && nodeIDs[link.Target] {
			filteredConfig.Links = append(filteredConfig.Links, link)
		}
	}

	// Loop through all of the nodes, if the parent is not in the nodeIds map, remove the reference to the parent
	for i, node := range filteredConfig.Nodes {
		if _, exists := nodeIDs[node.Parent]; !exists {
			filteredConfig.Nodes[i].Parent = ""
		}
	}

	return filteredConfig, nil
}

// linkMatchesAllFilters checks if a link satisfies all the provided filters.
func linkMatchesAllFilters(link configuration.Link, filters []Filter) (bool, error) {
	for _, filter := range filters {
		matches, err := linkMatchesFilter(link, filter)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}
	return true, nil
}

// linkMatchesFilter checks if a link satisfies a single filter condition.
func linkMatchesFilter(link configuration.Link, filter Filter) (bool, error) {
	// Extract the value of the specified field from the link
	fieldValue, err := getLinkFieldValue(link, filter.Condition.Field)
	if err != nil {
		return false, err
	}

	// Perform the comparison based on the operator
	switch filter.Condition.Operator {
	case "equals":
		return fieldValue == filter.Condition.Value, nil
	case "notEquals":
		return fieldValue != filter.Condition.Value, nil
	case "and":
		// Check if all conditions are met
		for _, condition := range filter.Condition.Conditions {
			matches, err := linkMatchesFilter(link, Filter{Condition: condition})
			if err != nil {
				return false, err
			}
			if !matches {
				return false, nil
			}
		}
		return true, nil
	case "or":
		// Check if any condition is met
		for _, condition := range filter.Condition.Conditions {
			matches, err := linkMatchesFilter(link, Filter{Condition: condition})
			if err != nil {
				return false, err
			}
			if matches {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported operator '%s'", filter.Condition.Operator)
	}
}

// getLinkFieldValue retrieves the value of a specified field from a link.
// It first checks the top-level fields, then the Attributes map.
func getLinkFieldValue(link configuration.Link, field string) (string, error) {
	switch field {
	case "source":
		return link.Source, nil
	case "target":
		return link.Target, nil
	case "type":
		return link.Type, nil
	default:
		// Check in Attributes, if 'attribute.' prefix remove it and look up attribute
		if len(field) > 10 && field[:10] == "attribute." {
			field = field[10:]
		} else {
			return "", nil
		}

		if val, exists := link.Attributes[field]; exists {
			// Assuming all attribute values are strings for simplicity
			if strVal, ok := val.(string); ok {
				return strVal, nil
			}
			return "", fmt.Errorf("attribute '%s' is not a string", field)
		}

		return "", nil
	}
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
	case "notEquals":
		return fieldValue != filter.Condition.Value, nil
	case "and":
		// Check if all conditions are met
		for _, condition := range filter.Condition.Conditions {
			matches, err := nodeMatchesFilter(node, Filter{Condition: condition})
			if err != nil {
				return false, err
			}
			if !matches {
				return false, nil
			}
		}
		return true, nil
	case "or":
		// Check if any condition is met
		for _, condition := range filter.Condition.Conditions {
			matches, err := nodeMatchesFilter(node, Filter{Condition: condition})
			if err != nil {
				return false, err
			}
			if matches {
				return true, nil
			}
		}
		return false, nil
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
		// Check in Attributes, if 'attribute.' prefix remove it and look up attribute
		if len(field) > 10 && field[:10] == "attribute." {
			field = field[10:]
		} else {
			return "", nil
		}

		if val, exists := node.Attributes[field]; exists {
			// Assuming all attribute values are strings for simplicity
			if strVal, ok := val.(string); ok {
				return strVal, nil
			}
			return "", fmt.Errorf("attribute '%s' is not a string", field)
		}

		return "", nil
	}
}
