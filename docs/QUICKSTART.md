---
layout: default
title: Quick Start
nav_order: 2
permalink: /quickstart
---

# Quick Start
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Validate Config

```bash
YAMLtecture --configIn=./tests/simple/config.yaml --validateConfig
```

## Merge Configs

```bash
YAMLtecture --in=./tests/complex/configs/ --mergeConfig
```

## Generate Graph

```bash
YAMLtecture --configIn=./tests/simple/config.yaml --generateMermaid
```

Rendering a graph with mermaid can be done on the CLI with [mermaid-cli](https://github.com/mermaid-js/mermaid-cli).

Piping these togeher you can generate and open a graph in a single command:

```bash
rm -f out.png && YAMLtecture -file=./tests/complex/config.yaml -graph | mmdc -i - -o ./out.png && open ./out.png
```
