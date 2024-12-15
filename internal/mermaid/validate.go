package mermaid

import "fmt"

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

	return nil
}
