# Architecture

## Overview

The SysML v2 Parser uses a **two-tier API design** with clear separation between low-level parsing and high-level model manipulation.

## Architecture Pattern

```
┌─────────────────────────────────────────────────────────────┐
│                      High-Level API                         │
│                     (sysml package)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │    Model    │  │   Visitor   │  │   Element Types     │  │
│  │  (model.go) │  │ (visitor.go)│  │ (dependency.go, ...)│  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Low-Level API                          │
│                      (low package)                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │    Lexer    │  │   Parser    │  │      Errors         │  │
│  │ (lexer.go)  │  │ (parser.go) │  │    (errors.go)      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  ANTLR-Generated Parser                     │
│                   (internal/parser)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │    Lexer    │  │   Parser    │  │     Listener        │  │
│  │(generated)  │  │(generated)  │  │   (generated)       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

1. **Input** — SysML source (string or file)
2. **Lexical Analysis** — Tokenization via ANTLR lexer
3. **Parsing** — Build parse tree via ANTLR parser
4. **Model Building** — Parse tree listener builds Go model
5. **Reference Resolution** — Resolve cross-references between elements
6. **Output** — Typed Go model with full element hierarchy

## Key Abstractions

### Model
- `sysml.Model` — Root container with packages and imports
- `sysml.Package` — Namespace container
- `sysml.Element` — Base interface for all model elements

### References
- `sysml.Ref[T]` — Type-safe generic reference with resolution
- Automatic resolution after parsing completes

### Visitor Pattern
- `sysml.Visitor` interface for model traversal
- `sysml.BaseVisitor` provides default implementations
- Supports early termination (return `false` to stop)

## Entry Points

### High-Level API
- `sysml.ParseString(input)` — Parse from string
- `sysml.ParseFile(path)` — Parse from file
- `sysml.ParseDirectory(dir)` — Parse directory

### Low-Level API
- `low.Parse(input)` — Direct parse tree access
- `low.NewLexer(input)` — Token stream access

## Design Decisions

1. **Type Safety**: Generic `Ref[T]` for compile-time type checking
2. **Zero Discarding**: All parsed elements appear in model
3. **Two-Tier API**: Low-level for performance, high-level for ergonomics
4. **Visitor Pattern**: Extensible traversal without modifying elements
