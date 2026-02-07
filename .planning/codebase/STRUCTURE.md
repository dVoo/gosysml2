# Directory Structure

## Overview

The project follows standard Go project layout with clear separation of concerns.

## Directory Layout

```
.
├── cmd/                          # Command-line tools
│   └── check_validation.go       # Validation checker
├── docs/                         # Documentation
│   ├── bnf/                      # BNF grammar specifications
│   ├── PERFORMANCE.md            # Performance guide
│   ├── README.md                 # API reference
│   ├── USAGE.md                  # Usage guide
│   └── usage-guide.md            # Project-level guide
├── examples/                     # Example programs
│   ├── basic/                    # Basic parsing
│   ├── parallel/                 # Parallel parsing
│   ├── requirements/             # Requirements traceability
│   ├── validation/               # Validation example
│   └── visitor/                  # Visitor pattern
├── internal/                     # Internal implementation
│   └── parser/                   # ANTLR-generated parser
├── libraries/                    # SysML standard libraries
├── low/                          # Low-level parsing API
│   ├── doc.go
│   ├── errors.go
│   ├── lexer.go
│   └── parser.go
├── sysml/                        # High-level SysML API
│   ├── doc.go
│   ├── errors.go
│   ├── model.go
│   ├── parse.go
│   └── visitor.go
├── validationdata/               # Validation test files
├── go.mod                        # Go module
├── go.sum                        # Dependency checksums
└── README.md                     # Project entry point
```

## Key Locations

### Source Code
- `low/` — Low-level parsing API (lexer, parser wrappers)
- `sysml/` — High-level model API (types, parsing, visitor)
- `internal/parser/` — ANTLR-generated code (do not modify)

### Tests
- `*_test.go` files alongside source files
- `sysml/integration_test.go` — Integration tests
- `sysml/bench_test.go` — Benchmarks

### Documentation
- `README.md` — Project overview (root)
- `docs/README.md` — Full API reference
- `docs/USAGE.md` — Usage examples
- `docs/PERFORMANCE.md` — Performance guide

### Examples
- `examples/basic/` — Basic parsing example
- `examples/parallel/` — Concurrent parsing
- `examples/requirements/` — Requirements handling
- `examples/validation/` — Validation example
- `examples/visitor/` — Custom visitor

### Configuration
- `go.mod` — Module definition
- `flake.nix` — Nix development environment

## Naming Conventions

### Files
- `*_test.go` — Test files
- `*_bench_test.go` — Benchmark files
- `doc.go` — Package documentation

### Packages
- `low` — Low-level API
- `sysml` — High-level API
- `internal/parser` — Internal generated code
