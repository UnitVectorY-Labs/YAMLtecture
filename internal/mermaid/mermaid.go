package mermaid

import (
	"fmt"
	"strings"

	"YAMLtecture/internal/configuration"
)

// GenerateMermaid creates a Mermaid diagram based on the relationships and includes all nodes.
func GenerateMermaid(config *configuration.Config) (string, error) {
	var sb strings.Builder
	sb.WriteString("flowchart TD\n")

	// Add all of the nodes
	for _, node := range config.Nodes {
		id := node.ID
		// Represent as a standalone node
		line := fmt.Sprintf("    %s\n", id)
		sb.WriteString(line)
	}

	// Add all of the relationships
	for _, rel := range config.Relationships {
		source := rel.Source
		target := rel.Target
		label := rel.Type

		// Mermaid syntax: source -->|label| target
		line := fmt.Sprintf("    %s -->|%s| %s\n", source, label, target)
		sb.WriteString(line)
	}

	return sb.String(), nil
}
