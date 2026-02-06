# SysML v2 Parser

A Go library for parsing SysML v2 models. Provides both a low-level performance-oriented API and a high-level developer-friendly API.

## Project Status

- **Grammar Coverage**: ~73% (58 of ~80 grammar elements implemented)
- **Validation Success Rate**: 96.4% (54 of 56 validation files pass)
- **Library Support**: Full SysML v2 standard library resolution (52 packages, 2605 elements)
- **Go Version**: 1.25+ with modern generics and iterator support

## Features

- Full SysML v2 grammar support via ANTLR4-generated parser
- Two-tier API design:
  - **Low-level API** (`low` package): Direct access to lexer, parser, and parse trees
  - **High-level API** (`sysml` package): Idiomatic Go model with visitor pattern
- Comprehensive error handling with source locations
- Memory-efficient parsing options for large repositories
- Parallel and streaming parsing modes
- Element finder functions for requirements, verifications, parts, and more

## Installation

```bash
go get github.com/dVoo/gosysml2_oc_oc
```

## Quick Start

### Parse a SysML String

```go
package main

import (
    "fmt"
    "github.com/dVoo/gosysml2_oc_oc/sysml"
)

func main() {
    input := `
        package Vehicle {
            part def Engine {
                attribute power : Real;
            }
            part def Car {
                part engine : Engine;
            }
        }
    `

    result := sysml.ParseString(input)
    if !result.Success() {
        fmt.Printf("Parse error: %s\n", result.Errors)
        return
    }

    // Access the model
    for _, pkg := range result.Model.Packages {
        fmt.Printf("Package: %s\n", pkg.Name())
    }
}
```

## Documentation

- **[API Reference](docs/README.md)** — Complete API documentation
- **[Usage Guide](docs/USAGE.md)** — Comprehensive usage examples and patterns
- **[Performance Guide](docs/PERFORMANCE.md)** — Performance optimization tips
- **[Project Usage Guide](docs/usage-guide.md)** — Project-level documentation

## Repository Structure

```
.
├── cmd/                    # Command-line tools
├── examples/               # Example programs
│   ├── basic/              # Basic parsing example
│   ├── requirements/       # Requirements traceability
│   ├── validation/         # Validation example
│   ├── parallel/           # Parallel parsing
│   └── visitor/            # Visitor pattern example
├── internal/               # Internal parser implementation
├── low/                    # Low-level parsing API
├── sysml/                  # High-level SysML API
├── docs/                   # Documentation
│   ├── bnf/                # BNF grammar specifications
│   ├── README.md           # API reference
│   ├── USAGE.md            # Usage guide
│   └── PERFORMANCE.md      # Performance guide
├── code/                   # ANTLR4 grammar files
├── libraries/              # SysML v2 standard library files
└── validationdata/         # Validation test files (18 categories)
```

## Development Environment

This project uses Nix flakes for environment management:

```bash
nix develop          # Enter dev shell with antlr, openjdk, claude-code
```

## Building

```bash
# Build all packages
go build ./...

# Run tests
go test ./...

# Build command-line tools
go build ./cmd/...
```

## Command-Line Tools

### verify-completeness

Analyzes SysML files and provides detailed statistics:

```bash
go run ./cmd/verify-completeness validationdata/parts-tree/
```

### check_validation

Validates the parser against the test suite:

```bash
go run ./cmd/check_validation
```

## Validation

The parser is validated against the official SysML v2 validation suite:

- **Location**: `validationdata/` directory (18 categories of test cases)
- **Coverage**: 96.4% success rate (54 of 56 files pass)
- **Categories**: Parts Tree, Requirements, Verification, State-based Behavior, Use Cases, and more

## License

See LICENSE file for details.
