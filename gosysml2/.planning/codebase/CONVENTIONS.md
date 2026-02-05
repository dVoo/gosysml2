# Coding Conventions

**Analysis Date:** 2026-02-05

## Naming Patterns

**Files:**
- Go source files use lowercase with underscores: `parse.go`, `lexer.go`, `parser_test.go`
- Test files follow Go convention: `*_test.go` (e.g., `parse_test.go`, `parser_test.go`, `integration_test.go`, `visitor_test.go`)
- Package-level files organized by functional area: `errors.go`, `model.go`, `parse.go`, `visitor.go`

**Functions:**
- Public functions use PascalCase: `ParseString`, `ParseFile`, `FindRequirements`, `NewPackage`
- Private functions use camelCase: `buildModel`, `extractName`, `extractRedefinitionName`, `convertFromLowLevel`
- Constructor functions use `New` prefix: `NewParser`, `NewLexer`, `NewPackage`, `NewRequirement`, `NewRef`, `NewCounter`
- Receiver method names in constructors: `NewParserFromLexer`, `NewLexerFromStream`, `NewLexerFromBytes`
- Getter methods without `Get` prefix: `Name()`, `Kind()`, `Location()`, `Parent()` (not `GetName()`)
- Setter methods use `Set` prefix: `SetParent()`, `SetDocumentation()` but non-public setters are direct assignment

**Types:**
- Public types use PascalCase: `Element`, `Package`, `Part`, `Requirement`, `Parser`, `Lexer`
- Interface types: `Element`, `Definition`, `Usage`, `Visitor` (no `I` prefix, despite ANTLR context types using `IContext` suffix)
- Private types use camelCase: `baseElement`, `modelBuilder`, `parseConfig`, `parserConfig`
- Constants and enums use PascalCase with descriptive prefix: `KindPackage`, `KindPart`, `KindRequirement` (ElementKind iota constants)
- Options type: `ParseOption` (function type for options pattern)

**Variables:**
- Loop variables: short form acceptable (`i`, `idx`, `err`, `d`, `r`)
- Receiver names: typically single letter `p` for Package, `r` for Requirement, `l` for Lexer, `a` for Attribute
- Package-scoped errors: use descriptive names even if private (`lexerErrors`, `parserErrors`)
- Short-lived values: `sb` for `strings.Builder`, `ec` for `ErrorCollector`
- Slices for collections: `parts`, `packages`, `requirements`, `files`, `results` (plural names)

## Code Style

**Formatting:**
- Standard Go formatting (no external formatter configuration detected; uses `go fmt` default)
- Indentation: 1 tab (Go standard)
- Line length: No hard limit enforced; code generally keeps under 100-120 columns
- Spacing: Single blank lines between function definitions, double blank lines between major sections within functions (discouraged in Go, but used in builder loops)

**Linting:**
- No `.golangci.yml` or ESLint configuration detected
- Code follows Go idioms and best practices (errors as return values, defer for cleanup)
- No type assertions without comma-ok checking pattern
- Proper nil checking throughout

## Import Organization

**Order:**
1. Standard library imports (e.g., `fmt`, `io`, `os`, `strings`, `sync`, `testing`)
2. Third-party imports (e.g., `github.com/antlr4-go/antlr/v4`)
3. Local imports (e.g., `github.com/dVoo/gosysml2/internal/parser`, `github.com/dVoo/gosysml2/low`)

**Path Aliases:**
- `sysml` package is the high-level public API (`github.com/dVoo/gosysml2/sysml`)
- `low` package is the low-level parser API (`github.com/dVoo/gosysml2/low`)
- `internal/parser` contains ANTLR-generated code (not exported)
- No aliases used in imports; full paths are typical

## Error Handling

**Patterns:**
- Errors returned as explicit return values (Go idiomatic): `func Parse(input string) (Tree, *ParseErrors)`
- Success status checked via method: `result.Success()` returns `true` if `result.Errors == nil || !result.Errors.HasErrors()`
- Type-specific error types: `Error` (single error), `ParseError` (collection with source info), `SyntaxError` (low-level), `ParseErrors` (aggregated)
- Errors implement `error` interface with `Error() string` method
- First error accessible via `.First()` method for quick error reporting
- Error aggregation: `ParseErrors.All()` combines lexer and parser errors
- Custom error types include context information: `Error` has `Context` field, `ParseError` has `Source` field
- Panic used only in intentional error cases: `MustParseString()` and `MustParseFile()` panic on error

**Low-level error collection:**
- `ErrorCollector` implements ANTLR `ErrorListener` interface
- Captured at construction time and passed to ANTLR: `lexer.AddErrorListener(errors)`
- Errors collected during parse (not thrown/panicked)

## Logging

**Framework:** No logging framework detected; uses standard `testing.T.Log()` and `testing.T.Logf()` for test diagnostics

**Patterns:**
- Diagnostics only in test code: `t.Logf("Token count: %d", len(tokens))`
- No production logging; library is silent on success
- Error messages are descriptive with line/column information: `"line %d, column %d: %s"`
- Source tracking includes filename or identifier: `source: filename`

## Comments

**When to Comment:**
- Package documentation comments: Present on package declarations
  - `// Package low provides a low-level, performance-oriented API for parsing SysML v2.`
  - `// Package sysml provides a high-level, developer-friendly API for working with SysML v2 models.`
- Type documentation: Comments above public types and interfaces explain purpose
- Method documentation: Comments above public receiver methods
- Complex logic: Comments explain non-obvious algorithm choices
  - Example: `// Help GC by clearing references` before `result.Model = nil`
- Implementation notes: Comments explain why something is done a certain way
  - Example: `// NOTE: Removed tokens.Fill() - let ANTLR consume tokens lazily`
- Stack operations: Comments document element/package stack management in builder
- TODO/FIXME: Not commonly used (none found in main code)

**Documentation Style:**
- GoDoc format: `// FunctionName does X` (first word is function name)
- Multi-line docs start with function/type name: `// ParseString parses a SysML string and returns a high-level model.`
- Receiver method comment format: `// EnterPackage_ captures package declarations and builds the model.`
- No special markers; plain English description

## Function Design

**Size:** Functions are reasonably scoped (typical 10-50 lines); builder visitor methods (Enter/Exit) are 20-30 lines

**Parameters:**
- Options pattern used: `func ParseString(input string, opts ...ParseOption) *ParseResult`
- Options are function closures: `type ParseOption func(*parseConfig)`
- Variadic parameters: limited use, primarily for options
- Context receivers implicit in methods (no explicit ctx parameter)

**Return Values:**
- Single return value is common for simple operations: `func (l *Lexer) NextToken() antlr.Token`
- Multiple return values for error cases: `func Parse(input string) (Tree, *ParseErrors)`
- Pointer returns for complex types: `*ParseResult`, `*Model`, `*Package`
- Result struct used for compound results: `ParseResult` contains `Model`, `Errors`, `Tree`, `Source`
- No explicit error variables; wrapped in result structs or returned alongside data

## Module Design

**Exports:**
- Clear separation: `sysml` package re-exports high-level API, `low` package exports low-level API
- Internal ANTLR code in `internal/parser` (not exported)
- Public API types exported from package root: `sysml.ParseString`, `sysml.Element`, `sysml.Model`
- Helper functions exported: `sysml.FindRequirements`, `sysml.FindParts`, `sysml.Walk`, `sysml.Visit`

**Barrel Files:**
- No barrel files (`index.ts` style) used
- Each Go file is independent; related types are grouped in same file
- `parse.go` contains parsing functions and ParseResult type
- `model.go` contains all element types (Package, Part, Requirement, etc.)
- `visitor.go` contains Visitor interface and traversal functions

## API Patterns

**Two-tier design:**
- **Low-level** (`low` package): Direct ANTLR access with error collection
  - `NewLexer(input)`, `NewParser(input)`, `Parse(input)` return raw ANTLR types
  - Suitable for grammar exploration and performance-critical code
- **High-level** (`sysml` package): Idiomatic Go model with visitor pattern
  - `ParseString(input)`, `ParseFile(filename)` return `*ParseResult` with semantic model
  - Suitable for typical use cases

**Options pattern:**
- Parsing functions accept variadic options: `ParseString(input, opts...)`
- Options are functions that modify config: `WithDiscardTree()` returns closure
- Allows backward-compatible addition of new parsing modes
- Config struct pattern: `type parseConfig struct { discardTree bool }`

**Streaming and parallel APIs:**
- Sequential: `ParseDirectory(dir)` returns all results
- Parallel: `ParseDirectoryParallel(dir, workers)` with semaphore control
- Streaming: `ParseDirectoryStream(dir, handler)` with callback for each file
- All variants share the same parsing logic

**Reference resolution:**
- Generic `Ref[T Element]` type for type-safe references
- References are resolved automatically after parsing: `model.ResolveReferences()`
- Unresolved references stored during parse, resolved in post-processing phase
- Getter: `r.Resolved()`, checker: `r.IsResolved()`, setter: `r.Resolve(elem)`

## Generics Use

**Pattern:**
- Generic `Ref[T Element]` for flexible reference handling
- Single use case: element references with type safety
- Go 1.22 constraint style (minimal; no complex bounds)

## Interface Implementation

**Pattern:**
- Interfaces define contracts; implementations use embedded structs
- `baseElement` embedded in all element types provides common interface implementation
- Interface method set determined by receiver methods only (no explicit `implements`)
- Type assertions for specialized behavior: `switch p := parent.(type) { case *Package: ... }`
