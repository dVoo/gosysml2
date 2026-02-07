# External Integrations

## Overview

The SysML v2 Parser is a self-contained library with minimal external dependencies. It does not integrate with external APIs, databases, or services.

## Dependencies

### ANTLR4 Runtime
- **Purpose**: Parser runtime for generated parsers
- **Package**: `github.com/antlr4-go/antlr/v4`
- **Version**: v4.13.1
- **Usage**: Tokenization, parsing, parse tree traversal

### Standard Library
The project relies heavily on Go standard library:
- `iter` — Iterator support (Go 1.23+)
- `context` — Context cancellation support
- `fmt`, `strings`, `os`, `path/filepath` — Standard utilities

## No External Services

This project does not integrate with:
- External APIs
- Databases
- Authentication providers
- Message queues
- Webhooks
- Cloud services

## File System Interactions

The parser reads from:
- Local `.sysml` and `.kerml` files
- `libraries/` directory — SysML standard libraries
- `validationdata/` directory — Test files

## Build-Time Dependencies

- **ANTLR4** — Used only at build time to generate parsers from grammar files
- **Java** — Required for ANTLR tool execution

## Summary

The SysML v2 Parser is designed to be a standalone, embeddable library with zero runtime dependencies beyond the ANTLR runtime and Go standard library.
