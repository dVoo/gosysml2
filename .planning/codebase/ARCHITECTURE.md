# Architecture

**Analysis Date:** 2026-02-05

## Pattern Overview

**Overall:** Layered Parser Architecture with Visitor Pattern

This is a SysML v2 parser library implemented in Go. It uses ANTLR4 for grammar-based parsing and implements a layered architecture that separates low-level parsing from high-level model construction.

**Key Characteristics:**
- **ANTLR4-generated parser** for SysML v2 grammar recognition
- **Two-tier API design**: Low-level direct access and high-level developer-friendly interface
- **Visitor pattern** for model traversal and processing
- **Reference resolution** system for cross-element relationships
- **Lazy token consumption** for memory-efficient parsing of large files

## Layers

### Layer 1: Generated Parser (internal/parser)

- **Purpose:** ANTLR-generated lexer and parser for SysML v2 grammar
- **Location:** `gosysml2/internal/parser/`
- **Contains:** Generated Go code from ANTLR grammar (lexer, parser, listener interfaces)
- **Depends on:** `github.com/antlr4-go/antlr/v4`
- **Used by:** `gosysml2/low` package

**Key files:**
- `sysmlv2_lexer.go` - Generated lexer (~86KB, ANTLR)
- `sysmlv2_parser.go` - Generated parser (~2MB, main grammar implementation)
- `sysmlv2parser_listener.go` - Listener interface definitions
- `sysmlv2parser_base_listener.go` - Base listener implementation

### Layer 2: Low-Level API (low)

- **Purpose:** Thin wrapper around ANTLR parser providing direct access with error collection
- **Location:** `gosysml2/low/`
- **Contains:** Lexer wrapper, Parser wrapper, Error types
- **Depends on:** `internal/parser`, `antlr4-go/antlr/v4`
- **Used by:** `gosysml2/sysml`, CLI tools

**Key abstractions:**
```go
// Parser wraps ANTLR parser with configuration options
type Parser struct {
    parser       *parser.SysMLv2Parser
    lexerErrors  *ErrorCollector
    parserErrors *ErrorCollector
}

// Lexer wraps ANTLR lexer with error collection
type Lexer struct {
    lexer  *parser.SysMLv2Lexer
    errors *ErrorCollector
}

// ErrorCollector implements antlr.ErrorListener
type ErrorCollector struct {
    errors []*SyntaxError
}
```

**Main entry points:**
- `low.Parse(input)` - Parse string to AST
- `low.ParseBytes(input)` - Parse byte slice to AST
- `low.Validate(input)` - Validate without building parse tree
- `low.NewParser(input)` - Create configurable parser

### Layer 3: High-Level API (sysml)

- **Purpose:** Developer-friendly model construction and manipulation
- **Location:** `gosysml2/sysml/`
- **Contains:** Domain model types, parsing facade, visitor pattern, reference resolution
- **Depends on:** `gosysml2/low`, `gosysml2/internal/parser`
- **Used by:** Applications, CLI tools, tests

**Key abstractions:**

**Element hierarchy:**
```go
// Element is the base interface for all SysML elements
type Element interface {
    Kind() ElementKind
    Name() string
    QualifiedName() string
    Location() Location
    Parent() Element
    Children() []Element
    Documentation() string
}

// Definition marks definition elements
type Definition interface {
    Element
    isDefinition()
}

// Usage marks usage elements
type Usage interface {
    Element
    Type() Element
    isUsage()
}
```

**Reference system:**
```go
// Ref represents a reference that may be resolved or unresolved
type Ref[T Element] struct {
    name     string
    resolved T
}
```

**Model root:**
```go
type Model struct {
    Packages     []*Package
    Imports      []*Import
    Comments     []*Comment
    Elements     []Element
    elementIndex map[string]Element // For O(1) qualified name lookup
}
```

### Layer 4: CLI Tools (cmd)

- **Purpose:** Command-line utilities for testing and verification
- **Location:** `cmd/*/`
- **Contains:** Standalone executables for various testing purposes

**Tools:**
- `cmd/verify-parser/` - Parser completeness verification (parse tree analysis)
- `cmd/verify-completeness/` - Model extraction verification (element stats, unresolved refs)
- `cmd/test-low-level/` - Low-level API testing
- `cmd/test-attrs/` - Attribute parsing tests
- `cmd/test-requirement-attributes/` - Requirement attribute tests

## Data Flow

### Parse Flow (High-Level API)

```
1. Input (string/file/bytes)
   |
   v
2. low.Parse() / low.ParseBytes()
   |
   v
3. ANTLR Lexer -> Token Stream
   |
   v
4. ANTLR Parser -> Parse Tree
   |
   v
5. modelBuilder (listener) -> Model
   |
   v
6. Model.BuildIndex() -> elementIndex
   |
   v
7. Model.ResolveReferences() -> resolved Ref[T]
   |
   v
8. ParseResult{Model, Errors, Tree}
```

### Visitor Flow

```
1. Visit(model, visitor)
   |
   v
2. For each top-level element:
   visitElement(elem, visitor)
   |
   v
3. Type switch on element:
   - Call specific VisitXxx() method
   - Returns bool (continue?)
   |
   v
4. If continue, recurse to children
```

## Key Abstractions

### Element Types

**Core structural elements:**
- `Package` - Namespace container with typed child accessors
- `Part` - Part definition/usage with attributes and nested parts
- `Item` - Item definition/usage
- `Attribute` - Typed attributes with optional default values
- `Port` - Port definition/usage with direction

**Behavioral elements:**
- `Action` - Action definition/usage
- `State` - State definition/usage with transitions
- `Transition` - State machine transitions
- `Connection` - Connection between elements

**Requirements:**
- `Requirement` - Requirement definition/usage with constraints
- `RequirementConstraint` - assume/require constraints
- `Verification` - Verification case definition/usage
- `Concern` - Stakeholder concern

**Analysis:**
- `UseCase` - Use case definition/usage
- `AnalysisCase` - Analysis case definition/usage
- `Calculation` - Calculation definition/usage
- `Constraint` - Constraint definition/usage

**Views:**
- `View` - View definition/usage exposing elements
- `Viewpoint` - Viewpoint definition with concerns

### Reference Resolution

Two-phase resolution system:

**Phase 1 (Parsing):** Store unresolved references as strings
- `unresolvedDerivedFrom []string`
- `unresolvedSatisfiedBy []string`
- `unresolvedSubject string`
- etc.

**Phase 2 (Resolution):** Resolve using elementIndex
```go
func (m *Model) ResolveReferences() {
    m.Walk(func(elem Element) bool {
        switch e := elem.(type) {
        case *Requirement:
            m.resolveRequirementRefs(e)
        // ... other types
        }
        return true
    })
}
```

Resolution strategy:
1. Direct qualified name lookup in `elementIndex`
2. Relative lookup walking up parent chain
3. Simple name lookup in any package

## Entry Points

### Library Entry Points

**High-level parsing:**
- `gosysml2/sysml.ParseString(input)` - Parse from string
- `gosysml2/sysml.ParseFile(filename)` - Parse from file
- `gosysml2/sysml.ParseBytes(input, source)` - Parse from bytes
- `gosysml2/sysml.ParseDirectory(dir)` - Parse all .sysml files
- `gosysml2/sysml.ParseDirectoryParallel(dir, workers)` - Parallel parsing
- `gosysml2/sysml.ParseDirectoryStream(dir, handler)` - Streaming parse

**Validation:**
- `gosysml2/sysml.Validate(input)` - Validate without building model
- `gosysml2/sysml.ValidateFile(filename)` - Validate file

**Low-level parsing:**
- `gosysml2/low.Parse(input)` - Parse to AST
- `gosysml2/low.ParseBytes(input)` - Parse bytes to AST
- `gosysml2/low.Validate(input)` - Validation only

### Application Entry Points

**Example:**
- `gosysml2/examples/main.go` - Library usage examples

**CLI tools:**
- `cmd/verify-parser/main.go` - Parse tree analysis
- `cmd/verify-completeness/main.go` - Model completeness check

## Error Handling

**Strategy:** Error collection with continuation

**Patterns:**
- ANTLR error listeners collect all errors without stopping
- `ParseResult` contains both Model and Errors
- `Success()` method checks if errors exist
- Low-level errors converted to high-level user-friendly errors

**Error types:**
```go
// Low-level
type SyntaxError struct {
    Line, Column int
    Message      string
    Source       string // "lexer" or "parser"
}

// High-level
type Error struct {
    Line, Column int
    Message      string
    Context      string
}
```

## Cross-Cutting Concerns

**Logging:** None - errors returned explicitly

**Validation:** Grammar-based via ANTLR, no additional semantic validation

**Authentication:** Not applicable (library code)

**Memory Management:**
- `WithDiscardTree()` option to free parse tree after model building
- `ParseDirectoryStream()` for memory-efficient batch processing
- Lazy token stream consumption in low-level parser

**Concurrency:**
- `ParseDirectoryParallel()` uses worker pool pattern
- Thread-safe for read operations on Model (no writes after construction)

---

*Architecture analysis: 2026-02-05*
