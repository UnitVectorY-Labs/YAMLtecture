package mermaid

import (
	"fmt"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/common"
	query "github.com/UnitVectorY-Labs/YAMLtecture/internal/query"
	"github.com/go-playground/validator/v10"
)

// Create a singleton validator instance
var validate = validator.New()

// Validate checks if the mermaid is valid.
func (m *Mermaid) Validate() error {

	// Validate the direction is valid
	switch m.Direction {
	case "TB":
	case "TD":
	case "BT":
	case "RL":
	case "LR":
	default:
		return fmt.Errorf("invalid direction: %s", m.Direction)
	}

	// Validate the node label is valid
	if m.NodeLabel != "" {
		// Perform same validation as attribute values
		err := common.IsValidValue(m.NodeLabel, "nodeLabel")
		if err != nil {
			return err
		}
	}

	// Validate the subgraph nodes are valid
	err := m.SubgraphNodes.Validate()
	if err != nil {
		return err
	}

	// Validate all of the node styles
	for _, nodeStyle := range m.NodeStyle {
		err := nodeStyle.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NodeStyle) Validate() error {

	// Validate the filters are valid
	for _, filter := range n.Filters {
		err := filter.Validate(query.NodeCondition)
		if err != nil {
			return err
		}
	}

	// Validate the format is valid
	err := n.Format.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (n *NodeStyleFormat) Validate() error {
	// Validate the fill is valid
	if n.Fill != "" {
		err := validate.Var(n.Fill, "hexcolor")
		if err != nil {
			return fmt.Errorf("invalid fill: %s", n.Fill)
		}
	}

	return nil
}
