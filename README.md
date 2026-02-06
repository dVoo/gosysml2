# SysML v2 Language Reference & Tools

This repository contains the SysML v2 language reference, including BNF grammar specifications, ANTLR parser definitions, and a Go library for parsing and working with SysML v2 models.

## Project Status

- **Grammar Coverage**: ~73% (58 of ~80 grammar elements implemented)
- **Validation Success Rate**: 96.4% (54 of 56 validation files pass)
- **Library Support**: Full SysML v2 standard library resolution (52 packages, 2605 elements)
- **Go Version**: 1.25+ with modern generics and iterator support

## Repository Structure

```
.
├── cmd/                    # Command-line tools
│   ├── test-attrs/         # Test requirement attributes parsing
│   ├── test-low-level/     # Low-level parser testing
│   ├── test-requirement-attributes/  # Test requirement attribute handling
│   ├── verify-completeness/          # Comprehensive model analysis tool
│   └── verify-parser/                # Parser verification tool
├── examples/               # Example SysML files for testing
├── analysis/               # Analysis reports and documentation
├── gosysml2/              # Go library for SysML v2 parsing (separate module)
│   ├── cmd/               # Library-specific tools
│   ├── examples/          # Library usage examples
│   ├── internal/          # Internal parser implementation
│   ├── low/               # Low-level parsing API
│   └── sysml/             # High-level SysML API
├── code/                  # ANTLR4 grammar files
│   ├── SysMLv2Lexer.g4   # Shared lexer
│   ├── KerMLParser.g4     # KerML parser
│   └── SysMLv2Parser.g4   # Full SysML parser
├── docs/                  # Language specification and documentation
│   ├── bnf/               # BNF grammar specifications
│   └── validationdata/    # SysML validation test files (18 categories)
├── libraries/             # SysML v2 standard library files
└── testdata/              # Test data for tools

```

## Getting Started

### Development Environment

This project uses Nix flakes for environment management:

```bash
nix develop          # Enter dev shell with antlr, openjdk, claude-code
```

### Building Tools

```bash
# Build all command-line tools
go build ./cmd/...

# Build a specific tool
go build ./cmd/verify-completeness
```

### Using the gosysml2 Library

The `gosysml2/` directory contains a production-ready Go library for parsing SysML v2 models:

```bash
cd gosysml2
go test ./...        # Run tests
go doc ./sysml       # View high-level API documentation
go doc ./low         # View low-level API documentation
```

#### Features

- **Two-tier API**: High-level model API (`sysml` package) and low-level parser API (`low` package)
- **Type-safe references**: Generic `Ref[T]` types with automatic resolution
- **Standard library support**: Automatic resolution of SysML v2 standard library imports
- **Multiple parsing modes**: Sequential, parallel, and streaming for different performance needs
- **Comprehensive element types**: Parts, requirements, verifications, use cases, state machines, and more
- **Visitor pattern**: Built-in visitor support for model traversal

See `gosysml2/README.md` for detailed API documentation and `gosysml2/examples/` for working code samples.

## Command-Line Tools

### verify-completeness

Analyzes SysML files and provides detailed statistics about the parsed model:

```bash
go run ./cmd/verify-completeness examples/test_example_req_verification.sysml
go run ./cmd/verify-completeness testdata/simple/
```

### test-attrs

Tests requirement attribute parsing:

```bash
go run ./cmd/test-attrs examples/test_example_req_verification.sysml
```

### verify-parser

Verifies parser completeness and coverage:

```bash
go run ./cmd/verify-parser
```

## Validation

The parser is validated against the official SysML v2 validation suite:

- **Location**: `validationdata/` directory (18 categories of test cases)
- **Coverage**: 96.4% success rate (54 of 56 files pass)
- **Categories**: Parts Tree, Requirements, Verification, State-based Behavior, Use Cases, and more

Run validation:
```bash
cd gosysml2
go run ./examples/validation
```

## ANTLR Grammar Files

The `code/` directory contains ANTLR4 grammar files that can be used to generate parsers for various target languages:

```bash
cd code
antlr -Dlanguage=Go -o parser SysMLv2Lexer.g4 SysMLv2Parser.g4
```

## Language Specification

The `docs/bnf/` directory contains the official BNF grammar specifications for:
- KerML (Kernel Modeling Language)
- SysML v2 textual syntax
- SysML v2 graphical notation

## License

See LICENSE file for details.
