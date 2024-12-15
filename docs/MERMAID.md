# Mermaid

One core capability of YAMLtecture is the ability to transform your YAML definitions into Mermaid.

## Generate Mermaid Flowchart

```bash
./YAMLtecture -configIn=./example/simple/architecture.yaml -mermaidIn=./example/simple/mermaid.yaml -generateMermaid
```

## Setting Configuration

An optional setting YAML file can be provided with the `--mermaidIn` flag. This file can contain the following settings:

```yaml
direction: "LR"
```

## Setting Attributes

The `direction` setting can be set to one of the following values:

- `TB` - Top to bottom
- `TD` - Top-down (same as top to bottom) - default
- `BT` - Bottom to top
- `RL` - Right to left
- `LR` - Left to right
