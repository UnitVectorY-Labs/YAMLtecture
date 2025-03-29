package mermaid

import (
	"fmt"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/common"
	query "github.com/UnitVectorY-Labs/YAMLtecture/internal/query"
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

	// Validate all of the node styles
	for _, nodeStyle := range m.NodeStyle {
		err := nodeStyle.Validate()
		if err != nil {
			return err
		}
	}

	// Validate all of the link styles
	for _, linkStyle := range m.LinkStyle {
		err := linkStyle.Validate()
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

	hasAttribute := false

	// Validate the fill is valid
	err := common.IsValidColor("fill", n.Fill)
	if err != nil {
		return err
	} else if n.Fill != "" {
		hasAttribute = true
	}

	// Validate the color is valid
	err = common.IsValidColor("color", n.Color)
	if err != nil {
		return err
	} else if n.Color != "" {
		hasAttribute = true
	}

	// Validate the stroke width is valid integer suffixed with 'px'
	err = common.IsValidPixel("stroke-width", n.StrokeWidth)
	if err != nil {
		return err
	} else if n.StrokeWidth != "" {
		hasAttribute = true
	}

	// Validate the font size is valid integer suffixed with 'px'
	err = common.IsValidPixel("font-size", n.FontSize)
	if err != nil {
		return err
	} else if n.FontSize != "" {
		hasAttribute = true
	}

	// Validate the padding is valid integer suffixed with 'px'
	err = common.IsValidPixel("padding", n.Padding)
	if err != nil {
		return err
	} else if n.Padding != "" {
		hasAttribute = true
	}

	// Validate the rx is valid integer suffixed with 'px'
	err = common.IsValidPixel("rx", n.Rx)
	if err != nil {
		return err
	} else if n.Rx != "" {
		hasAttribute = true
	}

	// Validate the ry is valid integer suffixed with 'px'
	err = common.IsValidPixel("ry", n.Ry)
	if err != nil {
		return err
	} else if n.Ry != "" {
		hasAttribute = true
	}

	// Ensure at least one attribute is set
	if !hasAttribute {
		return fmt.Errorf("at least one 'format' attribute must be set")
	}

	return nil
}

func (l *LinkStyle) Validate() error {

	// Validate the filters are valid
	for _, filter := range l.Filters {
		err := filter.Validate(query.LinkCondition)
		if err != nil {
			return err
		}
	}

	// Validate the format is valid
	err := l.Format.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkStyleFormat) Validate() error {

	hasAttribute := false

	// Validate the stroke is valid
	err := common.IsValidColor("stroke", l.Stroke)
	if err != nil {
		return err
	} else if l.Stroke != "" {
		hasAttribute = true
	}

	// Validate the stroke width is valid integer suffixed with 'px'
	err = common.IsValidPixel("stroke-width", l.StrokeWidth)
	if err != nil {
		return err
	} else if l.StrokeWidth != "" {
		hasAttribute = true
	}

	// Ensure at least one attribute is set
	if !hasAttribute {
		return fmt.Errorf("at least one 'format' attribute must be set")
	}

	return nil
}
