# Coding Conventions

**Analysis Date:** 2026-02-05

## Naming Patterns

**Files:**
- All lowercase, no underscores: `parse.go`, `model.go`, `visitor.go`
- Test files follow Go convention: `*_test.go` suffix
- Package structure: `sysml/`, `low/`, `internal/parser/`

**Packages:**
- Lowercase, single word where possible
- Main packages: `sysml`, `low`, `parser`
- Package comment on first line explaining purpose

**Functions:**
- **Exported:** PascalCase, descriptive: `ParseString`, `NewCounter`, `FindByKind`
- **Unexported:** camelCase: `buildModel`, `extractName`, `walkDepth`
- **Constructors:** `New{Type}` pattern: `NewParser`, `NewRequirement`, `NewModel`
- **Interface methods:** PascalCase, verb-noun pattern: `VisitPackage`, `AddChild`

**Variables:**
- Short names for local scope: `pkg`, `err`, `ctx`
- Descriptive names for exported: `parseConfig`, `elementIndex`
- Constants: PascalCase for exported, camelCase for unexported: `KindPackage`, `buildParseTree`

**Types:**
- **Interfaces:** PascalCase describing capability: `Visitor`, `Element`, `Definition`, `Usage`
- **Structs:** PascalCase nouns: `Parser`, `ParseResult`, `Requirement`, `baseElement`
- **Type parameters:** Single uppercase: `T` in `Ref[T Element]`
- **Function types:** PascalCase ending in `Func` or descriptive: `WalkFunc`, `ParseOption`

## Code Style

**Formatting:**
- Standard Go formatting (`gofmt` assumed)
- Tab indentation
- Line length ~120 characters typical maximum

**Linting:**
- No explicit linter configuration detected
- Code follows standard Go conventions

## Import Organization

**Order:**
1. Standard library imports (grouped)
2. Third-party imports (blank line separator)
3. Internal package imports

**Pattern:**
```go
import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/dVoo/gosysml2/internal/parser"
	"github.com/dVoo/gosysml2/low"
)
```

## Error Handling

**Patterns:**
- Custom error types implementing `error` interface
- Error wrapping with context: `fmt.Errorf("parse failed: %s", errors)`
- Error collector pattern for accumulating multiple errors
- Guard clauses for early returns on errors

**Error Types:**
```go
// Simple error with fields
type Error struct {
    Line    int
    Column  int
    Message string
    Context string
}

// Error collection
type ParseError struct {
    Errors []*Error
    Source string
}

// HasErrors pattern
func (e *ParseError) HasErrors() bool {
    return len(e.Errors) > 0
}
```

**Error Return Pattern:**
```go
result := ParseString(input)
if !result.Success() {
    return nil, result.Errors
}
```

## Logging

**Framework:** `testing.T` logging in tests only

**Test Logging Pattern:**
```go
t.Logf("Found %d SysML test files", len(files))
t.Logf("  PASS: %s", name)
t.Logf("  FAIL: %s - %s", name, result.Errors.First().Message)
```

## Comments

**Package Comments:**
```go
// Package sysml provides a high-level, developer-friendly API for working with SysML v2 models.
// This package offers convenient functions for parsing, traversing, and manipulating SysML content.
package sysml
```

**Interface Documentation:**
```go
// Visitor defines the interface for visiting SysML elements.
// Implement this interface to traverse and process a model.
type Visitor interface {
    // VisitPackage is called for each package element.
    // Return false to skip visiting children.
    VisitPackage(pkg *Package) bool
}
```

**Method Comments:**
```go
// ParseDirectoryParallel parses all .sysml files in a directory using multiple workers.
// This can significantly speed up parsing of large repositories on multi-core machines.
// Set workers to 0 to use runtime.NumCPU().
```

## Function Design

**Size:**
- Small, focused functions (typical 5-30 lines)
- Builder pattern for complex construction: `modelBuilder`
- Option pattern for configuration: `ParseOption`, `WithDiscardTree()`

**Parameters:**
- Context first if needed (not used here)
- Input sources second: `input string`, `r io.Reader`
- Variadic options last: `opts ...ParseOption`

**Return Values:**
- Result struct pattern: `*ParseResult` containing success/error/data
- Named returns for documentation (rarely used)
- Multiple returns: `(result, error)` or `(tree, errors)`

**Method Receivers:**
- Pointer receivers for mutable state: `func (p *Parser) ParseRootNamespace()`
- Value receivers for read-only: `func (r Ref[T]) Name() string`

## Module Design

**Exports:**
- Clear public API surface in each package
- Internal packages for generated code: `internal/parser/`
- Type-safe wrappers around generated code

**Barrel Files:**
- No barrel files; each file declares its package
- Related types grouped in single files (e.g., `model.go` contains ~2000 lines)

## Interface Patterns

**Marker Interfaces:**
```go
type Definition interface {
    Element
    isDefinition()
}

type Usage interface {
    Element
    Type() Element
    isUsage()
}
```

**Implementation Pattern:**
```go
func (p *Part) isDefinition() {}
func (p *Part) isUsage()      {}
```

**Generic References:**
```go
type Ref[T Element] struct {
    name     string
    resolved T
}
```

## State Management

**Struct Composition:**
```go
type modelBuilder struct {
    *parser.BaseSysMLv2ParserListener  // Embedded for ANTLR
    model        *Model
    currentPkg   *Package
    packageStack []*Package
    elementStack []Element
}
```

**Initialization:**
```go
func NewModel() *Model {
    return &Model{
        Packages:     make([]*Package, 0),
        Imports:      make([]*Import, 0),
        elementIndex: make(map[string]Element),
    }
}
```

---

*Convention analysis: 2026-02-05*
