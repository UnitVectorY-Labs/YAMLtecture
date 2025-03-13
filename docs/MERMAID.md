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

## Setting Attributes

The `direction` setting can be set to one of the following values:

- `TB` - Top to bottom
- `TD` - Top-down (same as top to bottom) - default
- `BT` - Bottom to top
- `RL` - Right to left
- `LR` - Left to right

```yaml
direction: "LR"
```

## Node Label

The `nodeLabel` attribute can be set to the name of the attribute that is set for a node that will be applied as the name of the node in the Mermaid flowchart. For example, using the "name" attribute allows the name of the node to be set to override `id` which is used by default.

```yaml
nodeLabel: "name"
```
