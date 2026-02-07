# Coding Conventions

## Overview

The project follows standard Go conventions with project-specific patterns for error handling and model design.

## Go Style

### Standard Conventions
- **gofmt** — All code formatted with standard Go formatter
- **golint** — Follows standard Go naming conventions
- **go vet** — Passes static analysis (except ANTLR-generated code)

### Naming

#### Packages
- Lowercase, single word: `sysml`, `low`
- No underscores or mixedCase

#### Types
- PascalCase for exported: `Model`, `ParseResult`, `Element`
- camelCase for unexported: `baseElement`, `elementKind`

#### Functions
- PascalCase for exported: `ParseString`, `FindParts`
- camelCase for unexported: `resolveRefs`, `buildIndex`

#### Constants
- PascalCase for exported: `KindPackage`, `KindPart`
- Use iota for sequential enums

## Error Handling

### Pattern
```go
// Wrap errors with context
return fmt.Errorf("parsing %s: %w", filename, err)

// Support errors.Is/errors.As
func (e *ParseError) Unwrap() error {
    return e.Cause
}
```

### Error Types
- `sysml.Errors` — Collection of parse errors
- `low.Error` — Low-level parse error with location
- Use `errors.Is()` and `errors.As()` for error checking

## Model Patterns

### Element Interface
```go
type Element interface {
    Name() string
    QualifiedName() string
    Location() Location
    Parent() Element
    Children() []Element
    Kind() ElementKind
}
```

### Reference Pattern
```go
// Type-safe reference
type Ref[T Element] struct {
    name string
    resolved T
}

func (r Ref[T]) Resolved() T
func (r Ref[T]) IsResolved() bool
```

### Base Element Embedding
```go
type Part struct {
    baseElement
    IsDefinition bool
    TypeRef Ref[*Part]
}
```

## Code Organization

### File Organization
- One type per file (e.g., `part.go`, `requirement.go`)
- Tests in `*_test.go`
- Package doc in `doc.go`

### Import Organization
```go
import (
    // Standard library
    "fmt"
    "strings"
    
    // Third-party
    "github.com/antlr4-go/antlr/v4"
    
    // Internal
    "github.com/dVoo/gosysml2_oc/internal/parser"
)
```

## Documentation

### Package Documentation
```go
// Package sysml provides a high-level API for parsing SysML v2 models.
//
// Quick start:
//
//     result := sysml.ParseString(input)
//     if result.Success() {
//         // Use result.Model
//     }
//
package sysml
```

### Function Documentation
```go
// ParseString parses SysML source from a string.
// Returns a ParseResult containing the model or errors.
func ParseString(input string, opts ...ParseOption) *ParseResult
```

## ANTLR-Generated Code

**Policy**: Do not modify files in `internal/parser/`

- Accept `go vet` warnings as generation artifacts
- Use wrapper types in `low/` package instead
