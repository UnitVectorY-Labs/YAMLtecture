---
layout: default
title: Development
nav_order: 7
permalink: /development
---

# Development
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Compile

Build the YAMLtecture command line application.

```bash
go build
```

## Run Tests

Run the test cases for YAMLtecture.

```bash
go test ./...
```

## Generate Tests

Generate the test cases output based on the configuration. The intended use case is that the generated files should pass the test cases and will be manually verified for correctness before being committed to the repository.

```bash
./generate.sh
```
