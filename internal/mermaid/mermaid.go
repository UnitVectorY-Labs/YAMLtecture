package mermaid

import (
	"fmt"
	"sort"
	"strings"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/common"
	"github.com/UnitVectorY-Labs/YAMLtecture/internal/configuration"
	query "github.com/UnitVectorY-Labs/YAMLtecture/internal/query"
)

// Mermaid contains the settings for generating the diagram.
type Mermaid struct {
	// The direction of the flowchart (TB, TD, BT, RL, LR)
	Direction string `yaml:"direction"`
	// The attribute to use as the node label (if set)
	NodeLabel string `yaml:"nodeLabel"`
	// The query to identify nodes to treat as subgraphs (explicit containers)
	SubgraphNodes query.Nodes `yaml:"subgraphNodes,omitempty"`
	// The style to apply to nodes
	NodeStyle []NodeStyle `yaml:"nodeStyles,omitempty"`
}

type NodeStyle struct {
	// The query to identify nodes to format with the style
	Filters []query.Filter `yaml:"filters"`
	// The style to apply to the nodes
	Format NodeStyleFormat `yaml:"format"`
}

type NodeStyleFormat struct {
	Fill string `yaml:"fill,omitempty"`
}

// SubgraphContainer holds a subgraph’s details, its nested explicit subgraphs,
// and any non‐explicit (leaf) node IDs that should be rendered inside it.
type subgraphContainer struct {
	ID        string
	Label     string // optional label if available
	Subgraphs []*subgraphContainer
	Nodes     []string // non‐explicit node IDs that fall under this container
}

// GenerateMermaid creates a Mermaid diagram from the config and mermaid settings.
// When a subgraph query is provided, nodes that have a parent will be placed
// inside the nearest explicit (query–matched) ancestor.
func GenerateMermaid(config *configuration.Config, setting *Mermaid) (string, error) {
	var mermaid strings.Builder

	// Write the header.
	mermaid.WriteString(fmt.Sprintf("flowchart %s\n", setting.Direction))

	// Write the node styles.
	styleMap := make(map[string][]string)
	if len(setting.NodeStyle) > 0 {
		mermaid.WriteString("    %% Node Styles\n")
		for i, style := range setting.NodeStyle {
			styleClassName := fmt.Sprintf("style%d", i)

			syntheticQuery := query.Query{
				Nodes: query.Nodes{
					Filters: style.Filters,
				},
			}

			// TODO: Get the nodes that need this style applied, we need to append those later but need to save them here
			nodes, err := query.ExecuteQuery(&syntheticQuery, config)
			if err != nil {
				return "", fmt.Errorf("error executing subgraph query: %v", err)
			}

			for _, node := range nodes.Nodes {
				styleMap[styleClassName] = append(styleMap[styleClassName], node.ID)
			}

			mermaid.WriteString(style.Format.print(styleClassName))
			mermaid.WriteString("\n")
		}

		mermaid.WriteString("\n")
	}

	mermaid.WriteString("    %% Nodes\n")

	// Build a lookup for nodes and a parent map.
	nodeLookup := make(map[string]configuration.Node)
	parentMap := make(map[string]string)
	for _, node := range config.Nodes {
		nodeLookup[node.ID] = node
		parentMap[node.ID] = node.Parent
	}

	// Determine which nodes are "explicit" subgraphs based on the query.
	explicit := make(map[string]bool)
	if len(setting.SubgraphNodes.Filters) > 0 {
		syntheticQuery := query.Query{
			Nodes: setting.SubgraphNodes,
		}
		subgraphConfig, err := query.ExecuteQuery(&syntheticQuery, config)
		if err != nil {
			return "", fmt.Errorf("error executing subgraph query: %v", err)
		}
		for _, node := range subgraphConfig.Nodes {
			explicit[node.ID] = true
		}
	}

	// Build explicit subgraph containers.
	containerMap := make(map[string]*subgraphContainer)
	for id := range explicit {
		label := ""
		if setting.NodeLabel != "" {
			if node, ok := nodeLookup[id]; ok {
				if val, ok := node.Attributes[setting.NodeLabel].(string); ok && val != "" {
					label = common.SanitizeLabel(val)
				}
			}
		}
		containerMap[id] = &subgraphContainer{
			ID:        id,
			Label:     label,
			Subgraphs: []*subgraphContainer{},
			Nodes:     []string{},
		}
	}

	// Helper: find the nearest explicit ancestor given a starting parent id.
	findExplicitAncestor := func(start string) string {
		cur := start
		for cur != "" {
			if explicit[cur] {
				return cur
			}
			cur = parentMap[cur]
		}
		return ""
	}

	// Lists for top-level explicit containers and standalone nodes.
	var topLevelExplicit []*subgraphContainer
	var topLevelNodes []string

	// Process non-explicit nodes: if they have a parent with an explicit container,
	// add them there; otherwise, treat them as top-level.
	for _, node := range config.Nodes {
		if explicit[node.ID] {
			continue // explicit nodes will be processed later
		}
		if node.Parent != "" {
			ancestor := findExplicitAncestor(node.Parent)
			if ancestor != "" {
				containerMap[ancestor].Nodes = append(containerMap[ancestor].Nodes, node.ID)
				continue
			}
		}
		topLevelNodes = append(topLevelNodes, node.ID)
	}

	// Process explicit nodes: nest them if their parent chain includes an explicit container.
	for id := range explicit {
		node := nodeLookup[id]
		if node.Parent != "" {
			ancestor := findExplicitAncestor(node.Parent)
			if ancestor != "" {
				containerMap[ancestor].Subgraphs = append(containerMap[ancestor].Subgraphs, containerMap[id])
				continue
			}
		}
		topLevelExplicit = append(topLevelExplicit, containerMap[id])
	}

	// Sort top-level explicit containers and standalone nodes for deterministic output.
	sort.Slice(topLevelExplicit, func(i, j int) bool {
		return topLevelExplicit[i].ID < topLevelExplicit[j].ID
	})
	sort.Strings(topLevelNodes)

	// Recursive helper to output an explicit container.
	var outputContainer func(cont *subgraphContainer, indent string)
	outputContainer = func(cont *subgraphContainer, indent string) {
		mermaid.WriteString(fmt.Sprintf("%ssubgraph %s\n", indent, cont.ID))
		// If a label is available, output it.
		if cont.Label != "" {
			mermaid.WriteString(fmt.Sprintf("%s    %s\n", indent, cont.Label))
		}
		// Output contained non-explicit nodes.
		sort.Strings(cont.Nodes)
		for _, nid := range cont.Nodes {
			node := nodeLookup[nid]
			if setting.NodeLabel != "" {
				if val, ok := node.Attributes[setting.NodeLabel].(string); ok && val != "" {
					mermaid.WriteString(fmt.Sprintf("%s    %s[%s]\n", indent, nid, common.SanitizeLabel(val)))
					continue
				}
			}
			mermaid.WriteString(fmt.Sprintf("%s    %s\n", indent, nid))
		}
		// Output nested explicit containers.
		sort.Slice(cont.Subgraphs, func(i, j int) bool {
			return cont.Subgraphs[i].ID < cont.Subgraphs[j].ID
		})
		for _, sub := range cont.Subgraphs {
			outputContainer(sub, indent+"    ")
		}
		mermaid.WriteString(fmt.Sprintf("%send\n", indent))
	}

	// Output all top-level explicit containers.
	for _, cont := range topLevelExplicit {
		outputContainer(cont, "    ")
	}
	// Output remaining top-level nodes.
	for _, nid := range topLevelNodes {
		node := nodeLookup[nid]
		if setting.NodeLabel != "" {
			if val, ok := node.Attributes[setting.NodeLabel].(string); ok && val != "" {
				mermaid.WriteString(fmt.Sprintf("    %s[%s]\n", nid, common.SanitizeLabel(val)))
				continue
			}
		}
		mermaid.WriteString(fmt.Sprintf("    %s\n", nid))
	}

	// Output the classes to format the nodes
	if len(styleMap) > 0 {
		mermaid.WriteString("\n")
		mermaid.WriteString("    %% Node Styles\n")
		for nodeID, styles := range styleMap {
			mermaid.WriteString(fmt.Sprintf("    class %s %s\n", strings.Join(styles, ","), nodeID))
		}
	}

	// Output the links.
	mermaid.WriteString("\n")
	mermaid.WriteString("    %% Links\n")
	sort.Slice(config.Links, func(i, j int) bool {
		if config.Links[i].Source == config.Links[j].Source {
			return config.Links[i].Target < config.Links[j].Target
		}
		return config.Links[i].Source < config.Links[j].Source
	})
	for _, rel := range config.Links {
		line := fmt.Sprintf("    %s -->|%s| %s\n", rel.Source, rel.Type, rel.Target)
		mermaid.WriteString(line)
	}

	return mermaid.String(), nil
}

func (f NodeStyleFormat) print(name string) string {
	var style strings.Builder

	style.WriteString("    classDef ")
	style.WriteString(name)
	style.WriteString(" ")
	if f.Fill != "" {
		style.WriteString(fmt.Sprintf("fill:%s;", f.Fill))
	}
	return style.String()
}
