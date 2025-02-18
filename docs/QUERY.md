---
layout: default
title: Query
nav_order: 5
permalink: /query
---

# Query

One core capability of YAMLtecture is the ability to apply a query to a configuration file to filter down to a subset of the nodes and links.  This is useful for taking a single larger definition configuration and applying different queries to filter down to a subset for different use cases.

## Query Syntax

Operator: `equals`

```yaml
nodes:
  filters:
    - condition:
        field: type
        operator: equals
        value: "Microservice"
```

Operator: `notEquals`

```yaml
nodes:
  filters:
    - condition:
        field: type
        operator: notEquals
        value: "Microservice"
```

Operator: `and`

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

Operator: `or`

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
