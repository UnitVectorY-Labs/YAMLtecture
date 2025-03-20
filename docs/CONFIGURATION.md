---
layout: default
title: Configuration
nav_order: 3
permalink: /configuration
---

# Configuration
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

The configuration file for YAMLtecture has two main building blocks, nodes that define the various elements in the architecture and links that define the links between the nodes. Both nodes and links can have attributes that provide additional metadata.

The core concept of YAMLtecture is the ability to define a detailed a detailed architecture in a simple to maintain YAML file establishing the different components and relationships between those components.

## Nodes & Links

A node is uniquely identified by its `id`.  Within a given YAMLtecure scope, the `id` must be unique.  The `id` is used to reference the node with the optional `parent` attribute as well as links referencing it with the `source` and `target` attributes.

Each node and link have a required `type` attribute.  The `type` attribute is used to define the type of the node or link.

Both nodes and links can have optional attributes.  The attributes are used to provide additional metadata about the node or link as key value pairs.


## Node Attributes

A node defines each part of the architecture.

### id

The `id` attribute is the required unique identifier for the node within the configuration used to reference the node in the `parent` attribute and in the `source` and `target` attributes of links.

### type

The `type` attribute is a required attribute the defines the type of the node.  This is the only mandatory attribute for a node that defines the metadata for the node as the intent here is that type would commonly be used for filtering or styling the nodes.

### parent

The `parent` attribute is an optional attribute that defines the parent node of the node.  This is used to define the hierarchy of the nodes as an acyclic graph.  The parent node enables special types of filtering given the tree structure of nodes it enforces.

### attributes

The `attributes` is a key value set of additional metadata for a node where the key is not predefined.  This allows for maximum flexibility in defining the attributes of the node. These attributes can be used for styling or filtering the nodes.

## Link Attributes

A link defines the relationship between two nodes.

### source

The `source` attribute is required and references the `id` of the source node.  This is used to define the source of the link.

### target

The `target` attribute is required and references the `id` of the target node.  This is used to define the target of the link.

### type

The `type` attribute is a required attribute the defines the type of the link.  This is the only mandatory attribute for a link that defines the metadata for the link as the intent here is that type would commonly be used for filtering or styling the links.

### attributes

The `attributes` is a key value set of additional metadata for a link where the key is not predefined.  This allows for maximum flexibility in defining the attributes of the link. These attributes can be used for styling or filtering the links.

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


