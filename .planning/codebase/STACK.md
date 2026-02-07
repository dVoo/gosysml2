# Technology Stack

## Overview

SysML v2 Parser is a Go-based library for parsing SysML v2 models, built with modern Go features and the ANTLR4 parser generator.

## Core Technologies

### Language & Runtime
- **Go 1.25+** — Primary language with generics and iterator support
- **Module**: `github.com/dVoo/gosysml2_oc`

### Parser Generation
- **ANTLR4** — Parser generator for SysML v2 grammar
- **antlr4-go/antlr/v4 v4.13.1** — Go runtime for ANTLR

### Dependencies

```go
require (
    github.com/antlr4-go/antlr/v4 v4.13.1
    golang.org/x/exp v0.0.0-20260112195511-716be5621a96 // indirect
)
```

## Build Tools

- **Standard Go toolchain** — `go build`, `go test`, `go mod`
- **Nix flakes** — Development environment management
- **ANTLR** — Grammar file compilation (via Nix)

## Development Environment

### Nix Shell
```bash
nix develop  # Enter dev shell with antlr, openjdk, claude-code
```

### Available Tools
- `antlr` — ANTLR4 parser generator
- `openjdk` — Java runtime for ANTLR
- Standard Go toolchain

## Configuration Files

- `go.mod` — Go module definition
- `go.sum` — Dependency checksums
- `flake.nix` — Nix development environment

## Generated Code

The parser implementation includes ANTLR-generated files:
- `internal/parser/sysmlv2_lexer.go` — Generated lexer
- `internal/parser/sysmlv2_parser.go` — Generated parser
- `internal/parser/sysmlv2parser_listener.go` — Parse tree listener interface
- `internal/parser/sysmlv2parser_base_listener.go` — Base listener implementation

**Note**: Generated files in `internal/parser/` should not be manually modified.
