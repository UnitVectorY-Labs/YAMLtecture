package mermaid

import (
	"fmt"
	"sort"
	"strings"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/configuration"
)

type Mermaid struct {
	// The direction of the flowchart (TB, TD, BT, RL, LR)
	Direction string `yaml:"direction"`
	// The attribute value to use as the node label if set
	NodeLabel string `yaml:"nodeLabel"`
}

// GenerateMermaid creates a Mermaid diagram based on the links and includes all nodes.
func GenerateMermaid(config *configuration.Config, setting *Mermaid) (string, error) {
	var mermaid strings.Builder
	mermaid.WriteString(fmt.Sprintf("flowchart %s\n", setting.Direction))

	mermaid.WriteString("    %% Nodes\n")

	// Sort nodes by ID
	sort.Slice(config.Nodes, func(i, j int) bool {
		return config.Nodes[i].ID < config.Nodes[j].ID
	})

	// Add all of the nodes
	for _, node := range config.Nodes {
		id := node.ID

		if setting.NodeLabel == "" {
			// Represent as a standalone node with no label
			mermaid.WriteString(fmt.Sprintf("    %s\n", id))
		} else {
			// Try to get the specified attribute value
			if val, ok := node.Attributes[setting.NodeLabel].(string); ok && val != "" {
				// Attribute exists and has a value - show node with label
				mermaid.WriteString(fmt.Sprintf("    %s[%s]\n", id, sanitizeLabel(val)))
			} else {
				// Attribute doesn't exist or is empty - show just the node
				mermaid.WriteString(fmt.Sprintf("    %s\n", id))
			}
		}
	}

	mermaid.WriteString("\n")
	mermaid.WriteString("    %% Links\n")

	// Sort links by source + target
	sort.Slice(config.Links, func(i, j int) bool {
		if config.Links[i].Source == config.Links[j].Source {
			return config.Links[i].Target < config.Links[j].Target
		}
		return config.Links[i].Source < config.Links[j].Source
	})

	// Add all of the links
	for _, rel := range config.Links {
		source := rel.Source
		target := rel.Target
		label := rel.Type

		// Mermaid syntax: source -->|label| target
		line := fmt.Sprintf("    %s -->|%s| %s\n", source, label, target)
		mermaid.WriteString(line)
	}

	return mermaid.String(), nil
}
