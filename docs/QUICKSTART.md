# Quickstart

## Compile

```bash
go build
```

## Run Tests

```
go test ./...
```

## Validate Example

```bash
./YAMLtecture -file=./example/simple/architecture.yaml -validate
```

## Generate Graph

```bash
./YAMLtecture -file=./example/simple/architecture.yaml -graph
```

Rendering a graph with mermaid can be done on the CLI with [mermaid-cli](https://github.com/mermaid-js/mermaid-cli).

Piping these togeher you can generate and open a graph in a single command:

```bash
rm -f out.png && ./YAMLtecture -file=./example/complex/architecture.yaml -graph | mmdc -i - -o ./out.png && open ./out.png
```