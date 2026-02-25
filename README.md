> [!CAUTION]
> This library and all of the contents are purely vibe-coded without further checks. Use at your own risk!

# gosysml2

Go library for parsing SysML v2 and KerML textual models with two layers:
- `low`: ANTLR-oriented low-level parser access
- `sysml`: high-level typed model with reference resolution

Current version: `v0.4.1` (from `VERSION`).

## Quick Start

Install:

```bash
go get github.com/dVoo/gosysml2
```

Parse a SysML string:

```go
package main

import (
    "fmt"

    "github.com/dVoo/gosysml2/sysml"
)

func main() {
    input := `
package Vehicle {
    part def Engine;
    part car : Engine;
}
`

    result := sysml.ParseString(input, sysml.WithSource("memory://Vehicle.sysml"))
    if err := result.Err(); err != nil {
        fmt.Println(err)
        return
    }

    parts := sysml.FindAll[*sysml.Part](result.Model)
    fmt.Printf("parts: %d\n", len(parts))
}
```

Parse a directory with one unified API:

```go
ctx := context.Background()
opts := sysml.DirOptions{
    Workers:      0, // 0 = NumCPU
    ParseOptions: []sysml.ParseOption{sysml.WithDiscardTree()},
}
for r := range sysml.ParseDir(ctx, "./models", opts) {
    if err := r.Err(); err != nil {
        fmt.Printf("failed %s: %v\n", r.Source, err)
        continue
    }
    fmt.Printf("ok: %s (%d top-level packages)\n", r.Source, len(r.Model.Packages()))
}
```

Semantic role and type-name APIs (canonical for tooling):

```go
for elem := range sysml.All(result.Model) {
    if elem.Role() == sysml.RoleUsage {
        if usage, ok := elem.(sysml.Usage); ok {
            fmt.Println(usage.TypeName())
        }
    }
}
```

## Documentation

Read in this order:

1. [`docs/QUICKSTART.md`](docs/QUICKSTART.md)
2. [`docs/LOW_LEVEL_GUIDE.md`](docs/LOW_LEVEL_GUIDE.md)
3. [`docs/HIGH_LEVEL_GUIDE.md`](docs/HIGH_LEVEL_GUIDE.md)
4. [`docs/EXAMPLES.md`](docs/EXAMPLES.md)

Reference docs:

- [`docs/README.md`](docs/README.md)
- [`docs/PARSER_LAYERS.md`](docs/PARSER_LAYERS.md)
- [`docs/SYSML_DATA_MODEL.md`](docs/SYSML_DATA_MODEL.md)
- [`docs/PERFORMANCE.md`](docs/PERFORMANCE.md)
- [`llms.txt`](llms.txt) for AI/LLM agent context

## Repository Structure

```
.
├── antlr/                  # ANTLR grammar files (.g4)
├── cmd/                    # Command-line tools
├── docs/                   # User and reference documentation
├── examples/               # Runnable examples
├── libraries/              # SysML/KerML library models
├── low/                    # Low-level parser layer
├── sysml/                  # High-level model layer
└── validationdata/         # Validation corpus
```

## Version Control (jj)

This repository uses `jj` with the Git backend.

```bash
jj st
jj describe -m "change description"
jj bookmark list
jj git push --bookmark main
```

## Build and Test

```bash
go build ./...
go test ./...
```
