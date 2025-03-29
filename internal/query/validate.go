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

	conditionType := ""
	if filterType == NodeCondition {
		conditionType = "nodes"
	} else if filterType == LinkCondition {
		conditionType = "links"
	} else {
		return fmt.Errorf("invalid filter type: %d", filterType)
	}

	allowCommand := false

	allowField := false
	requireField := false

	allowValue := false
	requireValue := false

	allowCondition := false
	requireCondition := false

	// Validate the operation is 'equals' using a switch so it is easy to add more operations later
	switch condition.Operator {
	case "equals":
		allowCommand = true
		requireField = true
		requireValue = true

	case "notEquals":
		allowCommand = true
		requireField = true
		requireValue = true

	case "exists":
		allowCommand = true
		requireField = true

	case "and":
		allowCommand = true
		requireCondition = true

	case "or":
		allowCommand = true
		requireCondition = true

	case "ancestorOf":
		allowCommand = filterType == NodeCondition
		requireValue = true

	case "descendantOf":
		allowCommand = filterType == NodeCondition
		requireValue = true

	case "childOf":
		allowCommand = filterType == NodeCondition
		requireValue = true

	case "parentOf":
		allowCommand = filterType == NodeCondition
		requireValue = true

	default:
		return fmt.Errorf("invalid operator: '%s'", condition.Operator)
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

	if !allowCommand {
		return fmt.Errorf("operator '%s' is not allowed for %s", condition.Operator, conditionType)
	} else if requireField && condition.Field == "" {
		return fmt.Errorf("'field' property is required for operator '%s'", condition.Operator)
	} else if !allowField && condition.Field != "" {
		return fmt.Errorf("'field' property is not allowed for operator '%s'", condition.Operator)
	} else if allowField {
		// Validate the field
		field := condition.Field
		switch field {
		case "type":
			// Allowed for everything
		case "id":
			if filterType != NodeCondition {
				return fmt.Errorf("field '%s' is not allowed for %s", field, conditionType)
			}
		case "parent":
			if filterType != NodeCondition {
				return fmt.Errorf("field '%s' is not allowed for %s", field, conditionType)
			}
		case "source":
			if filterType != LinkCondition {
				return fmt.Errorf("field '%s' is not allowed for %s", field, conditionType)
			}
		case "target":
			if filterType != LinkCondition {
				return fmt.Errorf("field '%s' is not allowed for %s", field, conditionType)
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
				return fmt.Errorf("invalid field: '%s'", field)
			}
		}
	}

	if requireValue && condition.Value == "" {
		return fmt.Errorf("'value' property is required for operator '%s'", condition.Operator)
	} else if !allowValue && condition.Value != "" {
		return fmt.Errorf("'value' property is not allowed for operator '%s'", condition.Operator)
	} else if allowValue {
		// Validate the value
		err = common.IsValidValue(condition.Value, "value")
		if err != nil {
			return err
		}
	}

	if requireCondition && len(condition.Conditions) == 0 {
		return fmt.Errorf("'conditions' property is required for operator '%s'", condition.Operator)
	} else if !allowCondition && len(condition.Conditions) > 0 {
		return fmt.Errorf("'conditions' property is not allowed for operator '%s'", condition.Operator)
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
