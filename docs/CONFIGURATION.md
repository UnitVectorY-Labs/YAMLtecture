# Configuration

The configuration file for YAMLtecture has two main building blocks, nodes that define the various elements in the architecture and links that define the links between the nodes. Both nodes and links can have attributes that provide additional metadata.

The core concept of YAMLtecture is the ability to define a detailed a detailed architecture in a simple to maintain YAML file establishing the different components and relationships between those components.

## Nodes & Links

A node is uniquely identified by its `id`.  Within a given YAMLtecure scope, the `id` must be unique.  The `id` is used to reference the node with the optional `parent` attribute as well as links referencing it with the `source` and `target` attributes.

Each node and link have a required `type` attribute.  The `type` attribute is used to define the type of the node or link.

Both nodes and links can have optional attributes.  The attributes are used to provide additional metadata about the node or link as key value pairs.

## Example Configuration

```yaml
nodes:
  - id: cluster
    type: Infrastructure
    attributes:
      name: "Container Hosting"
  - id: service_foo
    type: Microservice
    parent: cluster
    attributes:
      name: "Foo Service"
      language: "Java"
  - id: service_bar
    type: Microservice
    parent: cluster
    attributes:
      name: "Bar Service"
      language: "Go"

links:
  - source: service_foo
    target: service_bar
    type: "API"
    attributes:
      payload: "example"
```
