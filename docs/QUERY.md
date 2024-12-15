# Query

One core capability of YAMLtecture is the ability to apply a query to a configuration file to filter down to a subset of the nodes and links.  This is useful for taking a single larger definition configuration and applying different queries to filter down to a subset for different use cases.

## Query Syntax

```yaml
nodes:
  filters:
    - condition:
        field: type
        operator: equals
        value: "Microservice"
```
