# Architecture

**Analysis Date:** 2026-02-05

## Pattern Overview

**Overall:** Layered Parser Architecture with Visitor Pattern

**Key Characteristics:**
- Two-tier API design (high-level and low-level)
- ANTLR4-generated lexer/parser at the core
- Domain model built via parse tree visitors
- Reference resolution with two-phase parsing
- Generic `Ref[T]` type for type-safe element references

## Layers

**Generated Parser Layer:**
- Purpose: ANTLR4-generated lexer and parser for SysML v2 grammar
- Location: `internal/parser/`
- Contains: `sysmlv2_lexer.go`, `sysmlv2_parser.go`, `sysmlv2parser_listener.go`, `sysmlv2parser_base_listener.go`
- Depends on: `github.com/antlr4-go/antlr/v4`
- Used by: `low` package

**Low-Level API Layer:**
- Purpose: Performance-oriented access to lexer/parser with error collection
- Location: `low/`
- Contains: `lexer.go`, `parser.go`, `errors.go`
- Depends on: `internal/parser`
- Used by: `sysml` package, direct consumers needing parse tree access

**High-Level API Layer:**
- Purpose: Developer-friendly domain model with visitor pattern
- Location: `sysml/`
- Contains: `model.go`, `parse.go`, `visitor.go`, `errors.go`
- Depends on: `low` package, `internal/parser`
- Used by: Applications, examples, tests

**Application Layer:**
- Purpose: Example usage and client code
- Location: `examples/`
- Contains: `main.go`
- Depends on: `sysml` package, `low` package

## Data Flow

**Parse Flow (High-Level API):**

1. Input → `sysml.ParseString()` / `sysml.ParseFile()` / `sysml.ParseBytes()`
2. `parseWithSource()` delegates to `low.Parse()`
3. `low.NewParser()` creates ANTLR lexer and parser
4. ANTLR generates parse tree (`IEntryRuleRootNamespaceContext`)
5. `buildModel()` walks parse tree with `modelBuilder`
6. `modelBuilder` (ANTLR listener) constructs domain objects
7. `model.BuildIndex()` creates qualified name lookup index
8. `model.ResolveReferences()` resolves all `Ref[T]` references
9. `ParseResult` returned with `Model`, `Errors`, optional `Tree`

**Visitor Flow:**

1. Client calls `sysml.Visit(model, visitor)` or `sysml.Walk(model, fn)`
2. `visitElement()` dispatches to type-specific visitor methods
3. Visitor returns `bool` to control whether to visit children
4. `Walk()` provides depth tracking for tree traversal

**Reference Resolution Flow:**

1. During parsing: references stored as strings (e.g., `unresolvedSubject`)
2. Post-parse: `BuildIndex()` creates map of qualified name → element
3. Post-parse: `ResolveReferences()` walks all elements
4. For each element with references: calls `findElement()` with scope
5. Resolved references stored in `Ref[T].resolved` field
6. Unresolved references remain accessible via `Ref[T].Name()`

## Key Abstractions

**Element Interface:**
- Purpose: Base interface for all SysML model elements
- Location: `sysml/model.go` (lines 147-172)
- Pattern: Interface with common properties (Kind, Name, Location, Parent, Children)

**Definition/Usage Distinction:**
- Purpose: Distinguish between type definitions and usages
- Location: `sysml/model.go` (lines 174-186)
- Pattern: Marker interfaces (`isDefinition()`, `isUsage()`)
- Examples: `Part.IsDefinition`, `Requirement.IsDefinition`

**Generic Ref[T]:**
- Purpose: Type-safe references that may be resolved or unresolved
- Location: `sysml/model.go` (lines 113-145)
- Pattern: Generic struct with `name` and `resolved` fields
- Usage: `Ref[*Part]`, `Ref[Element]`, `Ref[*Requirement]`

**Visitor Pattern:**
- Purpose: Traverse and process model elements without modifying them
- Location: `sysml/visitor.go` (lines 1-81)
- Pattern: Interface with Visit* methods for each element type
- Base Implementation: `BaseVisitor` with no-op methods for easy embedding

**Parse Tree Listener:**
- Purpose: Build domain model from ANTLR parse tree
- Location: `sysml/parse.go` (lines 258-1733)
- Pattern: ANTLR listener implementing `BaseSysMLv2ParserListener`
- Key Type: `modelBuilder` struct with stacks for context tracking

## Entry Points

**High-Level Parsing:**
- Location: `sysml/parse.go`
- Functions: `ParseString()`, `ParseFile()`, `ParseBytes()`, `ParseReader()`
- Triggers: Application calls these functions directly
- Responsibilities: Parse input, build model, resolve references, return result

**Directory Parsing:**
- Location: `sysml/parse.go` (lines 149-238)
- Functions: `ParseDirectory()`, `ParseDirectoryParallel()`, `ParseDirectoryStream()`
- Triggers: Application needs to parse multiple files
- Responsibilities: File discovery, parallel processing, streaming with low memory

**Low-Level Parsing:**
- Location: `low/parser.go`
- Functions: `Parse()`, `ParseBytes()`, `Validate()`, `ValidateBytes()`
- Triggers: Direct low-level access needed
- Responsibilities: Return ANTLR parse tree and collected errors

**Validation:**
- Location: `sysml/parse.go` (lines 1704-1732), `low/parser.go` (lines 123-135)
- Functions: `Validate()`, `ValidateFile()`, `ValidateBytes()`
- Triggers: Syntax checking without model building
- Responsibilities: Fast validation with parse tree disabled

## Error Handling

**Strategy:** Error collection at each layer with transformation

**Patterns:**
- Low-level: `ErrorCollector` implements ANTLR `ErrorListener`
- Mid-level: `ParseErrors` aggregates lexer and parser errors
- High-level: `ParseError` provides user-friendly messages with source context
- Error transformation: `convertFromLowLevel()` maps `SyntaxError` → `Error`

**Error Types:**
- `low.SyntaxError`: Raw ANTLR errors with line, column, message
- `sysml.Error`: User-friendly errors with context
- `sysml.ParseError`: Collection of errors with source identification

## Cross-Cutting Concerns

**Logging:** Not used - errors returned explicitly

**Validation:** 
- Syntax validation via ANTLR parser
- Semantic validation during reference resolution
- No separate validation framework

**Memory Management:**
- Option `WithDiscardTree()` removes parse tree after model build
- Streaming parser clears references to help GC
- Lazy token consumption in `low/parser.go` (no `tokens.Fill()`)

**Concurrency:**
- `ParseDirectoryParallel()` uses goroutines with semaphore
- Thread-safe for independent parses
- Model building is single-threaded per file

---

*Architecture analysis: 2026-02-05*
