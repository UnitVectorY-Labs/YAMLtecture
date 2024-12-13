package query

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"YAMLtecture/internal/configuration"
)

type Query struct {
	Nodes Nodes `yaml:"nodes"`
}

type Nodes struct {
	Filters []Filter `yaml:"filters"`
}

type Filter struct {
	Condition Condition `yaml:"condition"`
}

type Condition struct {
	Field    string `yaml:"field"`
	Operator string `yaml:"operator"`
	Value    string `yaml:"value"`
}

// LoadQuery loads the query from a YAML file.
func LoadQuery(filepath string) (*Query, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var query Query
	err = yaml.Unmarshal(data, &query)
	if err != nil {
		return nil, err
	}
	return &query, nil
}

func (q *Query) Validate() error {
	if len(q.Nodes.Filters) == 0 {
		return errors.New("no filters found")
	}

	// Loop through filters
	for _, filter := range q.Nodes.Filters {

		// Check if field is empty
		if filter.Condition.Field == "" {
			return errors.New("field is empty")
		}
		// Operator must be "equals"
		if filter.Condition.Operator != "equals" {
			return errors.New("operator must be 'equals'")
		}
		// Check if value is empty
		if filter.Condition.Value == "" {
			return errors.New("value is empty")
		}
	}

	return nil
}

// ApplyQuery filters the configuration based on the provided query.
// It returns a new configuration containing only the nodes and relationships
// that match the query conditions.
func ApplyQuery(query *Query, config *configuration.Config) (configuration.Config, error) {

	// Prepare a new Config to hold the filtered results
	filteredConfig := configuration.Config{
		Nodes:         make(map[string]configuration.Node),
		Relationships: []configuration.Relationship{},
	}

	// Iterate over all nodes and apply filters
	for id, node := range config.Nodes {
		matchesAllFilters, err := nodeMatchesAllFilters(node, query.Nodes.Filters)
		if err != nil {
			return configuration.Config{}, fmt.Errorf("error applying filters to node '%s': %w", id, err)
		}
		if matchesAllFilters {
			filteredConfig.Nodes[id] = node
		}
	}

	// Iterate over all relationships and include only those where both source and target are in filtered nodes
	for _, rel := range config.Relationships {
		if _, sourceExists := filteredConfig.Nodes[rel.Source]; sourceExists {
			if _, targetExists := filteredConfig.Nodes[rel.Target]; targetExists {
				filteredConfig.Relationships = append(filteredConfig.Relationships, rel)
			}
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
