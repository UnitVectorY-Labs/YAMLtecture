package query

import (
	"YAMLtecture/internal/common"
	"fmt"
)

// Validate checks if the query is valid.
func (q *Query) Validate() error {

	// Validate the nodes
	err := q.Nodes.Validate()
	if err != nil {
		return err
	}

	return nil
}

// Validate checks if the nodes are valid.
func (nodes *Nodes) Validate() error {

	// Validate the filters
	for _, filter := range nodes.Filters {
		err := filter.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate checks if the filter is valid.
func (filter *Filter) Validate() error {

	// Validate the condition
	err := filter.Condition.Validate()
	if err != nil {
		return err
	}

	return nil
}

// Validate checks if the condition is valid.
func (condition *Condition) Validate() error {
	var err error

	// Validate the field
	field := condition.Field
	switch field {
	case "id":
	case "type":
	case "parent":
	default:

		// Check if key starts with 'attribute.'
		if len(field) > 10 && field[:10] == "attribute." {
			// Validate the attribute key
			err = common.IsValidName(field[10:], "attribute.key")
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid field: %s", field)
		}
	}

	// Validate the operation is 'equals' using a switch so it is easy to add more operations later
	switch condition.Operator {
	case "equals":
	default:
		return fmt.Errorf("invalid operator: %s", condition.Operator)
	}

	// Validate the value
	err = common.IsValidValue(condition.Value, "value")
	if err != nil {
		return err
	}

	return nil
}
