---
layout: default
title: Query
nav_order: 5
permalink: /query
---

# Query
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

One core capability of YAMLtecture is the ability to apply a query to a configuration file to filter down to a subset of the nodes and links.  This is useful for taking a single larger definition configuration and applying different queries to filter down to a subset for different use cases.

## Why Query?

The intended use of a YAMLtecture configuration file is that it represents a superset of the overall system with all of the detail represented.  The challenge this provides is that consuming this at full detail as in creating a mermaid flowchart results in too much detail.  Therefore, this information can be selected down with a query.

## Query Syntax

Queries are represented as YAML files and apply filters in the form of various operators that can take the source config file and reduce it down into the desired subset.

### Operator: `equals`

```yaml
nodes:
  filters:
    - condition:
        field: type
        operator: equals
        value: "Microservice"
```

### Operator: `notEquals`

```yaml
nodes:
  filters:
    - condition:
        field: type
        operator: notEquals
        value: "Microservice"
```

### Operator: `and`

```yaml
nodes:
  filters:
    - condition:
        operator: and
        conditions:
          - field: type
            operator: equals
            value: "Microservice"
          - field: attribute.name
            operator: equals
            value: "Service A"
```

### Operator: `or`

```yaml
nodes:
  filters:
    - condition:
        operator: or
        conditions:
          - field: attribute.name
            operator: equals
            value: "Service A"
          - field: attribute.name
            operator: equals
            value: "Service B"
```
