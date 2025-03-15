package mermaid

import (
	"fmt"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/common"
)

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

	return nil
}
