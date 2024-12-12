# YAMLtecture

A lightweight CLI tool for generating outputs, including Mermaid diagrams, from YAML-defined system architectures.

## Overview

YAMLtecture is an open-source CLI tool designed for architects who need a simple yet powerful way to define, validate, and visualize system architectures. Using modular YAML files, YAMLtecture helps you map out components, their hierarchies, and interactions, all while keeping everything flexible and easy to manage. It’s built to handle the complexity of modern systems without locking you into rigid formats. With YAMLtecture, you can effortlessly transform your YAML definitions into diagrams (like Mermaid), keeping your architecture clear, version-controlled, and up-to-date.

## Configuration Files

At the core of YAMLtecture are YAML configuration files that define your system architecture. These files are split into two main sections: nodes and relationships. The nodes section describes the individual components of your system, while the relationships section defines how those components interact.

```yaml
architecture:
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

  relationships:
    - source: service_foo
      target: service_bar
      type: "API"
      attributes:
        payload: "example"
```

The key concept here is that each node has a unique `id`. The `type` field is flexible and can be set to whatever suits your architecture—common examples include Microservice, Database, Queue, etc. The `parent` field is optional and defines hierarchical relationships between nodes. These relationships are validated to form an acyclic tree, but not every node needs to be part of the same hierarchy.

The `attributes` field is a dictionary of key-value pairs that can store additional metadata about a node, useful for filtering and querying.

Relationships tie nodes together with the `source` and `target` fields. Each relationship is one-way, so for a bidirectional connection, you’ll need to define two separate relationships. Like nodes, relationships have a `type` field to classify the interaction (e.g., API, MessageQueue), and `attributes` can store metadata about the relationship, such as protocols, payloads, or endpoints.
