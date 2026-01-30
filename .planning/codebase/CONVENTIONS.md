# Coding Conventions

**Analysis Date:** 2026-01-30

## Naming Patterns

**Files:**
- Snake_case for files with underscores (e.g., `sysml_lexer.go`, `sysml_parser.go`)
- Descriptive names corresponding to package scope (e.g., `parse.go`, `model.go`, `visitor.go`, `errors.go`)
- Test files suffixed with `_test.go` (e.g., `parse_test.go`, `integration_test.go`, `visitor_test.go`)
- ANTLR-generated files prefixed with tool output (e.g., `sysmlv2_lexer.go`, `sysmlv2_parser.go`)

**Functions:**
- PascalCase for exported functions: `ParseString`, `ParseFile`, `ParseBytes`, `MustParseFile`, `Validate`
- camelCase for unexported helper functions: `parseWithSource`, `buildModel`, `extractName`, `walkElement`, `visitElement`
- Constructor functions use `New*` prefix: `NewPackage()`, `NewPart()`, `NewRequirement()`, `NewModel()`
- Receiver methods use same casing as functions
- Boolean accessor methods use `Has*` or `Is*` prefixes: `HasErrors()`, `IsDefinition`, `IsResolved()`

**Variables:**
- Short names in local scope: `i`, `err`, `pkg`, `req`, `ver`, `elem`
- Descriptive names for package-level variables and struct fields: `elementIndex`, `unresolvedSubject`, `DerivedFrom`, `SatisfiedBy`
- Private struct fields are lowercase: `model`, `currentPkg`, `elementStack`, `children`, `packages`, `parts`
- Public struct fields are PascalCase: `Model`, `Errors`, `Tree`, `Source`, `IsDefinition`, `TypeRef`
- Error variables: `err` (short form in local scope), descriptive names for error types: `ParseError`, `SyntaxError`, `Error`

**Types:**
- PascalCase for all types (structs, interfaces, type aliases): `ParseResult`, `ParseError`, `ParseOption`, `Visitor`, `Element`, `Model`, `Package`, `Part`, `Requirement`
- Constants use PascalCase with Kind/Method prefixes: `KindPackage`, `KindPart`, `KindRequirement`, `VerificationMethodTest`, `PortDirectionIn`
- Interface names follow Go convention (no `I` prefix despite ANTLR-generated code using it): `Element`, `Visitor`, `Definition`, `Usage`

## Code Style

**Formatting:**
- Uses standard Go formatting (gofmt defaults)
- 80-120 character line length observed in practice
- Tab indentation (Go standard)
- Braces on same line (Go standard)

**Linting:**
- No explicit linting configuration files found; uses Go standard conventions
- Follows Go idiomatic patterns for error handling and return types
- Receivers on type methods are typically short (1-2 characters): `(p *Package)`, `(m *Model)`, `(b *modelBuilder)`

## Import Organization

**Order:**
1. Standard library imports (`io`, `os`, `path/filepath`, `strings`, `sync`, `fmt`, `testing`)
2. Third-party imports (`github.com/antlr4-go/antlr/v4`, `github.com/dVoo/gosysml2/internal/parser`, `github.com/dVoo/gosysml2/low`)
3. Local package imports (when in different packages within same module)

**Path Aliases:**
- No aliases used; fully qualified paths throughout
- Imports from `github.com/dVoo/gosysml2` prefix packages like `low` and `internal/parser`

**Example from `parse.go`:**
```go
import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/antlr4-go/antlr/v4"
	"github.com/dVoo/gosysml2/internal/parser"
	"github.com/dVoo/gosysml2/low"
)
```

## Error Handling

**Patterns:**
- Errors returned as second return value in tuples: `(result T, err error)` or `(tree ParseTree, errors ParseErrors)`
- Error types are custom structs that implement `error` interface: `ParseError`, `SyntaxError`, `ParseErrors`
- `HasErrors()` method used to check if error collection contains errors
- `Success()` method on result types: `ParseResult.Success()` returns `true` if no errors
- Wrapped errors with context using custom error types that include `Source`, `Line`, `Column`, `Message`
- Errors collected during parsing/lexing into `ErrorCollector` which implements `antlr.ErrorListener`
- Early returns on error: File I/O errors returned immediately with wrapped message

**Example from `parse.go`:**
```go
if !result.Success() {
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
```

## Logging

**Framework:** Uses standard `testing.T` for test logging; no runtime logging framework observed

**Patterns:**
- `t.Fatalf()` for fatal test failures that stop execution
- `t.Errorf()` for non-fatal test failures
- `t.Logf()` for informational logging in tests
- `t.Error()` for simple error reporting
- `t.Skip()` to skip tests when preconditions not met
- Comments used for documentation, not for runtime behavior explanation

## Comments

**When to Comment:**
- Function comments for all exported functions (Go standard doc comment pattern)
- Type comments for all exported types explaining their purpose
- Inline comments for non-obvious logic (e.g., "Removed tokens.Fill() - let ANTLR consume tokens lazily")
- Inline comments for performance considerations and workarounds
- Comments on constants and struct fields that aren't self-documenting

**GoDoc/Doc Comments:**
- All exported functions preceded by comment: `// FunctionName does something`
- All exported types preceded by comment: `// TypeName represents...`
- All exported struct fields lack doc comments (relying on field names being descriptive)
- Package-level doc comments at top: `// Package sysml provides...`
- Examples use inline comments and test cases to document behavior

**Example from `parse.go`:**
```go
// ParseResult contains the result of parsing SysML input.
type ParseResult struct {
	Model  *Model      // The parsed model (nil if parsing failed completely)
	Errors *ParseError // Any errors that occurred (nil if successful)
	Tree   antlr.Tree  // The raw parse tree (for advanced use)
	Source string      // Source file path or identifier
}

// Success returns true if parsing was successful (no errors).
func (r *ParseResult) Success() bool {
	return r.Errors == nil || !r.Errors.HasErrors()
}
```

## Function Design

**Size:** Functions generally 5-50 lines; tree-walking and parsing logic can reach 50-100 lines with clear purpose

**Parameters:**
- Constructors take essential fields: `NewPackage(name string, loc Location)`
- Parse functions take input with variadic options: `ParseString(input string, opts ...ParseOption)`
- Options pattern used: `ParseOption func(*parseConfig)` for flexible configuration
- Visitor methods take typed elements: `VisitPackage(pkg *Package) bool`

**Return Values:**
- Single value for successful operations: `ParseRootNamespace() parser.IEntryRuleRootNamespaceContext`
- Tuple returns for fallible operations: `(result T, err error)`
- Pointer receivers for methods that modify state: `(r *ParseResult)`, `(m *Model)`
- Value receivers for immutable types: `(k ElementKind) String() string`

## Module Design

**Exports:**
- Clear public API in package-level exported types and functions
- Lowercase private helpers and internal builders: `parseConfig`, `modelBuilder`, `extractName`
- Two-package architecture: `low` (direct ANTLR) and `sysml` (high-level API)

**Barrel Files:**
- `visitor.go` exports multiple visitor types and functions: `Visitor` interface, `BaseVisitor`, `Visit()`, `Walk()`, `Filter()`, `FindByKind()`, `FindByName()`, `FindPackages()`, etc.
- `model.go` exports model navigation and building: `Model`, `Package`, `Element` interfaces, all element types (`Part`, `Requirement`, `Verification`, etc.)
- `parse.go` exports parsing entry points and configuration
- `errors.go` exports error types and formatting

**Private Implementation Details:**
- ANTLR parser integration hidden in `low` package
- `modelBuilder` is unexported but documented as the tree walker
- Reference resolution methods private: `resolveRequirementRefs()`, `resolvePartRefs()`, `findElement()`
- Tree construction details private in `parseWithSource()`, `buildModel()`

---

*Convention analysis: 2026-01-30*
