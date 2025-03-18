package query

import (
	"fmt"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/common"
)

const (
	NodeCondition = 1
	LinkCondition = 2
)

// Validate checks if the query is valid.
func (q *Query) Validate() error {

	// Validate the nodes
	err := q.Nodes.Validate()
	if err != nil {
		return err
	}

	// Validate the links
	err = q.Links.validate()
	if err != nil {
		return err
	}

	return nil
}

// Validate checks if the nodes are valid.
func (nodes *Nodes) Validate() error {

	// Validate the filters
	for _, filter := range nodes.Filters {
		err := filter.Validate(NodeCondition)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate checks if the links are valid.
func (links *Links) validate() error {

	// Validate the filters
	for _, filter := range links.Filters {
		err := filter.Validate(LinkCondition)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate checks if the filter is valid.
func (filter *Filter) Validate(filterType int) error {

	// Validate the condition
	err := filter.Condition.validate(filterType)
	if err != nil {
		return err
	}

	return nil
}

// Validate checks if the condition is valid.
func (condition *Condition) validate(filterType int) error {
	var err error

	allowField := false
	requireField := false

	allowValue := false
	requireValue := false

	allowCondition := false
	requireCondition := false

	// Validate the operation is 'equals' using a switch so it is easy to add more operations later
	switch condition.Operator {
	case "equals":
		requireField = true
		requireValue = true

	case "notEquals":
		requireField = true
		requireValue = true

	case "exists":
		requireField = true

	case "and":
		requireCondition = true

	case "or":
		requireCondition = true

	case "ancestorOf":
		requireValue = true

	case "descendantOf":
		requireValue = true

	case "childOf":
		requireValue = true

	case "parentOf":
		requireValue = true

	default:
		return fmt.Errorf("invalid operator: %s", condition.Operator)
	}

	// Set the allow and require flags based on the filter type
	if requireValue {
		allowValue = true
	}

	if requireField {
		allowField = true
	}

	if requireCondition {
		allowCondition = true
	}

	if requireField && condition.Field == "" {
		return fmt.Errorf("field is required")
	} else if !allowField && condition.Field != "" {
		return fmt.Errorf("field is not allowed")
	} else if allowField {
		// Validate the field
		field := condition.Field
		switch field {
		case "type":
			// Allowed for everything
		case "id":
			if filterType != NodeCondition {
				return fmt.Errorf("field 'id' is not allowed")
			}
		case "parent":
			if filterType != NodeCondition {
				return fmt.Errorf("field 'parent' is not allowed")
			}
		case "source":
			if filterType != LinkCondition {
				return fmt.Errorf("field 'source' is not allowed")
			}
		case "target":
			if filterType != LinkCondition {
				return fmt.Errorf("field 'target' is not allowed")
			}
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
	}

	if requireValue && condition.Value == "" {
		return fmt.Errorf("value is required")
	} else if !allowValue && condition.Value != "" {
		return fmt.Errorf("value is not allowed")
	} else if allowValue {
		// Validate the value
		err = common.IsValidValue(condition.Value, "value")
		if err != nil {
			return err
		}
	}

	if requireCondition && len(condition.Conditions) == 0 {
		return fmt.Errorf("condition is required")
	} else if !allowCondition && len(condition.Conditions) > 0 {
		return fmt.Errorf("condition is not allowed")
	} else if allowCondition {

		// Validate the conditions
		for _, subCondition := range condition.Conditions {
			err = subCondition.validate(filterType)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
