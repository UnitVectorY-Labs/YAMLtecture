package mermaid

import (
	"fmt"
	"sort"
	"strings"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/configuration"
)

type Mermaid struct {
	Direction string `yaml:"direction"`
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
		// Represent as a standalone node
		line := fmt.Sprintf("    %s\n", id)
		mermaid.WriteString(line)
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
