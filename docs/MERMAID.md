---
layout: default
title: Mermaid
nav_order: 6
permalink: /mermaid
---

# Mermaid
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

One core capability of YAMLtecture is the ability to transform your YAML definitions into Mermaid.

## Generate Mermaid Flowchart

```bash
./YAMLtecture -configIn=./tests/simple/architecture.yaml -mermaidIn=./tests/simple/mermaid.yaml -generateMermaid
```

## Setting Configuration

An optional setting YAML file can be provided with the `--mermaidIn` flag. This file can contain the following settings:

- `direction` - The direction of the flowchart
- `nodeLabel` - The attribute to use as the node label
- `subgraphNodes` - The attribute to filter to identify nodes that will be used as subgraphs

All settings are optional but to generate a configuration file must be specified, even if it is empty.

### Direction

The `direction` setting can be set to one of the following values:

- `TB` - Top to bottom
- `TD` - Top-down (same as top to bottom) - default
- `BT` - Bottom to top
- `RL` - Right to left
- `LR` - Left to right

```yaml
direction: "LR"
```

### Node Label

The `nodeLabel` attribute can be set to the name of the attribute that is set for a node that will be applied as the name of the node in the Mermaid flowchart. For example, using the "name" attribute allows the name of the node to be set to override `id` which is used by default.

```yaml
nodeLabel: "name"
```

### Subgraph Nodes

The `subgraphNodes` attribute uses the same syntax as a query but instead of selecting the nodes to be include, it selects the nodes that will be used as subgraphs. For example, the following setting will create subgraphs for all nodes that have a `type` attribute set to `Application`.

```yaml
subgraphNodes:
  filters:
    - condition:
        field: type
        operator: equals
        value: "Application"
```

## Node Styles

The `nodeStyles` attribute is used to define the Mermaid styles that will be applied to the rendered nodes. The selection of which nodes to apply uses the same syntax as a query. Multiple styles can be applied to the same node but this behavior is non-deterministic and therefore should be avoided.  There are multiple attributes that can be set for a node style which each match the attributes that can be set in Mermaid for the class definition.

- `fill` - The fill color of the node background in RGB hex format.
- `color` - The text color of the node in RGB hex format.
- `stroke-width` - The thickness of the border of the node in pixels.

```yaml
nodeStyles:
  - style:
      filters:
        - condition:
            field: type
            operator: equals
            value: "Application"
      style:
        fill: "#f9f9f9"
        color: "#000000"
        stroke-width: 2px
```
