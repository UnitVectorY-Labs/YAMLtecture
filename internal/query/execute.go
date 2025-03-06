package query

import (
	"fmt"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/configuration"
)

// ConfigContext holds the configuration with pre-calculated maps for efficient querying
type ConfigContext struct {
	Config      *configuration.Config
	NodesById   map[string]*configuration.Node
	ChildrenMap map[string][]string // Maps parent node ID to list of child node IDs
}

// NewConfigContext creates a new ConfigContext with pre-calculated maps for efficient querying
func NewConfigContext(config *configuration.Config) *ConfigContext {
	// Create a map for faster node lookups
	nodesById := make(map[string]*configuration.Node)
	for i := range config.Nodes {
		nodesById[config.Nodes[i].ID] = &config.Nodes[i]
	}

	// Create a map of parent node IDs to their children's IDs
	childrenMap := make(map[string][]string)
	for _, node := range config.Nodes {
		if node.Parent != "" {
			childrenMap[node.Parent] = append(childrenMap[node.Parent], node.ID)
		}
	}

	return &ConfigContext{
		Config:      config,
		NodesById:   nodesById,
		ChildrenMap: childrenMap,
	}
}

// ExecuteQuery filters the configuration based on the provided query.
// It returns a new configuration containing only the nodes and links
// that match the query conditions.
func ExecuteQuery(query *Query, config *configuration.Config) (configuration.Config, error) {
	// Create context with pre-calculated maps for efficient querying
	ctx := NewConfigContext(config)

	// Prepare a new Config to hold the filtered results
	filteredConfig := configuration.Config{
		Nodes: []configuration.Node{},
		Links: []configuration.Link{},
	}

	// Iterate over all nodes and apply filters
	for _, node := range config.Nodes {
		matchesAllFilters, err := nodeMatchesAllFilters(node, query.Nodes.Filters, ctx)
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
func nodeMatchesAllFilters(node configuration.Node, filters []Filter, ctx *ConfigContext) (bool, error) {
	for _, filter := range filters {
		matches, err := nodeMatchesFilter(node, filter, ctx)
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
func nodeMatchesFilter(node configuration.Node, filter Filter, ctx *ConfigContext) (bool, error) {
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
			matches, err := nodeMatchesFilter(node, Filter{Condition: condition}, ctx)
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
			matches, err := nodeMatchesFilter(node, Filter{Condition: condition}, ctx)
			if err != nil {
				return false, err
			}
			if matches {
				return true, nil
			}
		}
		return false, nil
	case "ancestorOf":
		return isAncestorOf(node.ID, filter.Condition.Value, ctx)
	case "descendantOf":
		return isDescendantOf(node.ID, filter.Condition.Value, ctx)
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

// isDescendantOf checks if nodeID is a descendant of targetNodeID in the given configuration context.
// A descendant is a node in the direct child chain (child, child's child, etc.)
func isDescendantOf(nodeID string, targetNodeID string, ctx *ConfigContext) (bool, error) {
	// Find the target node
	_, exists := ctx.NodesById[targetNodeID]
	if !exists {
		return false, fmt.Errorf("target node with ID %s not found", targetNodeID)
	}

	// Use a stack-based approach to traverse the descendant tree
	// Start with the children of the target node
	children, exists := ctx.ChildrenMap[targetNodeID]
	if !exists {
		// Target node has no children
		return false, nil
	}

	// Check direct children first
	for _, childID := range children {
		if childID == nodeID {
			return true, nil
		}
	}

	// If not found in direct children, check deeper descendants
	for _, childID := range children {
		// For each child, check if nodeID is in its descendant tree
		childChildren, hasChildren := ctx.ChildrenMap[childID]
		if !hasChildren {
			continue
		}

		// Create a stack for depth-first traversal
		stack := append([]string{}, childChildren...)

		// Process the stack
		for len(stack) > 0 {
			// Pop the last item from the stack
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// Check if this is the node we're looking for
			if current == nodeID {
				return true, nil
			}

			// Add its children to the stack
			if descendants, hasDescendants := ctx.ChildrenMap[current]; hasDescendants {
				stack = append(stack, descendants...)
			}
		}
	}

	// If we get here, nodeID is not a descendant of targetNodeID
	return false, nil
}

// isAncestorOf checks if nodeID is an ancestor of targetNodeID in the given configuration context.
// An ancestor is a node in the direct parent chain (parent, parent's parent, etc.)
func isAncestorOf(nodeID string, targetNodeID string, ctx *ConfigContext) (bool, error) {
	// Find the target node
	targetNode, exists := ctx.NodesById[targetNodeID]
	if !exists {
		return false, fmt.Errorf("target node with ID %s not found", targetNodeID)
	}

	// Start checking with the direct parent
	currentParentID := targetNode.Parent

	// Traverse up the ancestor chain
	for currentParentID != "" {
		// If this ancestor matches the nodeID we're checking, return true
		if currentParentID == nodeID {
			return true, nil
		}

		// Move up to the next parent
		parent, exists := ctx.NodesById[currentParentID]
		if !exists {
			// If parent reference doesn't exist in the nodes, stop traversal
			break
		}
		currentParentID = parent.Parent
	}

	// If we get here, nodeID is not an ancestor of targetNodeID
	return false, nil
}
