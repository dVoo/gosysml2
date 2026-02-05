# Coding Conventions

**Analysis Date:** 2026-02-05

## Overview

This is a Go 1.22 codebase implementing a SysML v2 parser with a dual-layer API design. The project follows standard Go conventions with additional patterns for ANTLR-generated code organization.

## Naming Patterns

**Files:**
- Source files use `snake_case.go`: `parse.go`, `model.go`, `visitor.go`
- Test files use `*_test.go` suffix: `parse_test.go`, `visitor_test.go`
- Generated ANTLR files use PascalCase with `SysMLv2` prefix: `sysmlv2_parser.go`, `sysmlv2_lexer.go`

**Packages:**
- Package names are lowercase, single words: `sysml`, `low`, `parser`
- The main API package is `sysml` (high-level)
- The low-level package is `low` (direct ANTLR access)
- Generated parser code is in `internal/parser`

**Functions:**
- Public functions use PascalCase: `ParseString()`, `FindRequirements()`
- Private functions use camelCase: `buildModel()`, `extractName()`
- Constructor functions use `New` prefix: `NewParser()`, `NewRequirement()`
- Getter methods match field name without `Get` prefix: `Name()`, `Kind()`, `Children()`
- Boolean methods use `Has*` or `Is*` prefixes: `HasErrors()`, `IsResolved()`

**Types:**
- Interfaces use noun names: `Element`, `Definition`, `Usage`, `Visitor`
- Structs use noun names: `Package`, `Requirement`, `ParseResult`
- Pointer receivers for methods that modify state
- Value receivers for read-only methods

**Constants:**
- Use PascalCase for exported constants: `KindPackage`, `KindRequirement`
- Use iota for sequential constants in type enumerations

**Variables:**
- Short names in local scope: `pkg`, `req`, `err`
- Avoid single-letter names except for common patterns: `i`, `ok`, `t`
- Acronyms in all caps: `ID`, `QN` (for Qualified Name)

## Code Style

**Formatting:**
- Standard `gofmt` formatting - no custom rules detected
- No `.prettierrc` or custom formatting configuration
- No linting configuration files (`.eslintrc`, etc.) detected

**Line Length:**
- Generally under 100 characters
- Long function signatures wrap at parameter boundaries

**Import Organization:**

Standard Go grouping:
1. Standard library imports
2. Third-party imports (ANTLR)
3. Local/internal imports

```go
import (
	"io"
	"os"
	"path/filepath"

	"github.com/antlr4-go/antlr/v4"
	"github.com/dVoo/gosysml2/internal/parser"
	"github.com/dVoo/gosysml2/low"
)
```

**No Blank Imports:**
- Unused imports are not present (enforced by Go compiler)

## Error Handling

**Error Types:**
- Custom error structs in each package:
  - `sysml.Error` - User-friendly error with line/column
  - `low.SyntaxError` - Low-level syntax error from lexer/parser
  - `ParseError` - Aggregates multiple errors

**Error Patterns:**

```go
// Return error as value
func (r *ParseResult) Success() bool {
	return r.Errors == nil || !r.Errors.HasErrors()
}

// Early return on error with context
func ParseFile(filename string, opts ...ParseOption) *ParseResult {
	content, err := os.ReadFile(filename)
	if err != nil {
		return &ParseResult{
			Errors: &ParseError{
				Errors: []*Error{{
					Line:    0,
					Column:  0,
					Message: err.Error(),
				}},
				Source: filename,
			},
			Source: filename,
		}
	}
	return parseWithSource(string(content), filename, opts...)
}
```

**Panic Usage:**
- `MustParseString()` and `MustParseFile()` panic on error
- Convention: `Must` prefix indicates panic on failure
- Used only in convenience functions, not in library code

```go
func MustParseString(input string) *Model {
	result := ParseString(input)
	if !result.Success() {
		panic(result.Errors)
	}
	return result.Model
}
```

## Architecture Patterns

**Dual-Layer API Design:**

1. **High-Level API** (`sysml` package):
   - `ParseString()`, `ParseFile()`, `ParseDirectory()`
   - Returns domain model (`*Model`, `*Package`, `*Requirement`)
   - User-friendly error messages

2. **Low-Level API** (`low` package):
   - Direct access to ANTLR lexer/parser
   - Performance-focused with minimal overhead
   - Raw parse tree access

**Visitor Pattern:**

```go
// Interface definition in visitor.go
type Visitor interface {
	VisitPackage(pkg *Package) bool
	VisitPart(part *Part) bool
	VisitRequirement(req *Requirement) bool
	// ... more methods
}

// BaseVisitor provides default implementations
type BaseVisitor struct{}
func (BaseVisitor) VisitPackage(pkg *Package) bool { return true }
```

**Functional Options Pattern:**

```go
type ParseOption func(*parseConfig)

func WithDiscardTree() ParseOption {
	return func(c *parseConfig) {
		c.discardTree = true
	}
}

func ParseString(input string, opts ...ParseOption) *ParseResult {
	cfg := &parseConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	// ... use cfg
}
```

**Type Embedding:**

```go
// Base element embedded in concrete types
type baseElement struct {
	kind     ElementKind
	name     string
	location Location
	parent   Element
	children []Element
}

type Package struct {
	baseElement
	IsLibrary bool
}
```

**Interface Segregation:**

```go
// Element is the base interface
type Element interface {
	Kind() ElementKind
	Name() string
	QualifiedName() string
	Location() Location
	Parent() Element
	Children() []Element
}

// Definition marker interface
type Definition interface {
	Element
	isDefinition()
}
```

## Comments

**Package Documentation:**

```go
// Package sysml provides a high-level, developer-friendly API for working with SysML v2 models.
// This package offers convenient functions for parsing, traversing, and manipulating SysML content.
package sysml
```

**Function Documentation:**

```go
// ParseString parses a SysML string and returns a high-level model.
func ParseString(input string, opts ...ParseOption) *ParseResult

// MustParseString parses a SysML string and panics if there are errors.
func MustParseString(input string) *Model
```

**Type Documentation:**

```go
// ParseResult contains the result of parsing SysML input.
type ParseResult struct {
	Model  *Model      // The parsed model (nil if parsing failed completely)
	Errors *ParseError // Any errors that occurred (nil if successful)
	Tree   antlr.Tree  // The raw parse tree (for advanced use)
	Source string      // Source file path or identifier
}
```

## File Organization

**Within a Package:**
- Types and interfaces first
- Constructor functions (`New*`) next
- Methods grouped by receiver type
- Helper functions at end

**Example structure in model.go:**
1. Kind constants (`ElementKind`)
2. Core interfaces (`Element`, `Definition`, `Usage`)
3. Base struct (`baseElement`)
4. Concrete types (`Package`, `Part`, `Requirement`)
5. Model struct with methods
6. Resolution methods

## Concurrency Patterns

**Worker Pool for Parallel Parsing:**

```go
func ParseDirectoryParallel(dir string, workers int, opts ...ParseOption) ([]*ParseResult, error) {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	
	results := make([]*ParseResult, len(files))
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers) // Semaphore for limiting concurrency

	for i, file := range files {
		wg.Add(1)
		sem <- struct{}{} // Acquire

		go func(idx int, path string) {
			defer wg.Done()
			defer func() { <-sem }() // Release

			results[idx] = ParseFile(path, opts...)
		}(i, file)
	}

	wg.Wait()
	return results, nil
}
```

## Testing Conventions

See `TESTING.md` for detailed testing patterns.

Key conventions:
- Tests use descriptive names: `TestParseString`, `TestParseStringWithErrors`
- Test data embedded as raw strings with backticks
- Use `t.Fatalf` for fatal errors, `t.Errorf` for non-fatal
- `t.Logf` for debugging output

---

*Convention analysis: 2026-02-05*
