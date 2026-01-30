# Architecture

**Analysis Date:** 2026-01-30

## Pattern Overview

**Overall:** Multi-layered grammar-driven parser architecture with separation between low-level parsing and high-level model abstraction.

**Key Characteristics:**
- ANTLR4-based lexer/parser generation from grammar specifications
- Two-tier API design: low-level (parse trees) and high-level (idiomatic Go objects)
- Visitor pattern for model traversal
- Type-safe element references with automatic resolution
- Parallel and streaming parsing support for large repositories

## Layers

**Grammar Layer:**
- Purpose: Define SysML v2 and KerML syntax rules in ANTLR format
- Location: `code/SysMLv2Lexer.g4`, `code/KerMLParser.g4`, `code/SysMLv2Parser.g4`
- Contains: Lexical and parser rules defining all SysML v2 and KerML constructs
- Depends on: None (foundational)
- Used by: ANTLR code generation targeting multiple languages (Go, Java, etc.)

**Documentation Layer:**
- Purpose: Provide BNF/EBNF specifications and rendered grammar documentation
- Location: `docs/bnf/` with `.kebnf`, `.kgbnf` files and HTML renderings
- Contains: Formal grammar specifications with OMG specification clause references
- Depends on: Grammar files (for documentation alignment)
- Used by: Specification compliance and reference

**Low-Level Parser Layer (`low` package):**
- Purpose: Direct access to ANTLR-generated lexer, parser, and parse trees
- Location: `gosysml2/low/`
- Contains:
  - `lexer.go`: Wrapper around ANTLR lexer with error collection
  - `parser.go`: Wrapper around ANTLR parser with token stream management
  - `errors.go`: Error collection and reporting for lexer/parser
  - `*_test.go`: Low-level parsing tests
- Depends on: `gosysml2/internal/parser/` (ANTLR-generated code)
- Used by: High-level API layer and applications needing raw parse trees

**ANTLR Generated Layer:**
- Purpose: Lexer and parser implementation generated from grammar files
- Location: `gosysml2/internal/parser/`
- Contains:
  - `sysmlv2_lexer.go`: Generated lexer (~1300 lines)
  - `sysmlv2_parser.go`: Generated parser (~80k lines)
  - `sysmlv2parser_listener.go`: ANTLR listener interface
  - `sysmlv2parser_base_listener.go`: Base listener implementation
- Depends on: ANTLR runtime library (`github.com/antlr4-go/antlr/v4`)
- Used by: `low` package

**High-Level Model Layer (`sysml` package):**
- Purpose: Idiomatic Go abstraction over parse trees with semantic model
- Location: `gosysml2/sysml/`
- Contains:
  - `model.go`: Element types, kinds, references, and model structure (~1672 lines)
  - `parse.go`: Parsing entry points and orchestration (~612 lines)
  - `visitor.go`: Visitor pattern for model traversal (~374 lines)
  - `errors.go`: Error conversion and reporting
  - `*_test.go`: Integration tests
- Depends on: `low` package, ANTLR runtime
- Used by: Applications consuming SysML models

**Test/Example Layer:**
- Purpose: Validate parsing and demonstrate API usage
- Location: `gosysml2/examples/`, `testdata/`
- Contains:
  - `examples/main.go`: Complete API usage examples
  - `testdata/`: 34 `.sysml` test files covering all language features
- Depends on: `sysml` and `low` packages
- Used by: Developers and validation

## Data Flow

**Single File Parsing (High-Level):**

1. Application calls `sysml.ParseFile(path)` or `sysml.ParseString(input)`
2. High-level parser delegates to `low.Parse(input)`
3. Lexer tokenizes input → Parser builds ANTLR parse tree
4. Errors collected and converted to `ParseError` structures
5. Parse tree passed to `buildModel()` function
6. Visitor traverses tree, creating `Element` objects with typed fields
7. References resolved (qualified names → object pointers)
8. `ParseResult` returned with `Model`, `Errors`, and optional `Tree`

**Directory Parsing (Parallel Mode):**

1. Application calls `sysml.ParseDirectoryParallel(dir, workers)`
2. Walker finds all `.sysml` files in directory
3. File list distributed across worker goroutines
4. Each worker calls `ParseFile()` independently
5. Results collected in slice
6. Optional: `WithDiscardTree()` reduces memory per worker by ~30%

**Model Traversal (Visitor Pattern):**

1. Application calls `sysml.Visit(model, visitor)` or `sysml.Walk(model, func)`
2. Traversal iterates `model.Elements` (top-level packages)
3. For each element, dispatcher calls appropriate `Visitor.Visit*()` method
4. Visitor returns `bool`: `true` to recurse into children, `false` to skip
5. Children accessed via `Element.Children()` interface
6. Walk completes depth-first traversal

**Reference Resolution:**

1. During parse tree traversal, qualified names stored as `Ref[T]` with name only
2. After all elements parsed, `ResolveReferences()` called
3. Index built via `BuildIndex()` mapping qualified names → elements
4. For each reference, lookup performs `FindByQualifiedName()` search
5. Reference marked `resolved` if found; remains unresolved otherwise
6. Safe navigation: `ref.IsResolved()` checks before `ref.Resolved()` access

**State Management:**

- Parse tree is immutable after construction
- Model elements are mutable (can add/remove children after parsing)
- References are resolved once; no automatic updates if element moves
- Error collection happens during lexing and parsing phases
- Memory optimization: `WithDiscardTree()` discards tree after model built

## Key Abstractions

**Element Interface:**
- Purpose: Unified interface for all SysML model elements
- Examples: `Package`, `Part`, `Requirement`, `Verification`, `Action`
- Pattern: All elements implement `Element` interface with methods: `Kind()`, `Name()`, `QualifiedName()`, `Location()`, `Parent()`, `Children()`, `Documentation()`

**Ref[T] Generic Type:**
- Purpose: Type-safe, nullable references to elements
- Pattern:
  - Unresolved: `Ref{name: "Engine", resolved: nil}`
  - Resolved: `Ref{name: "Engine", resolved: *Part}`
  - Check before use: `if ref.IsResolved() { elem := ref.Resolved() }`

**ParseResult:**
- Purpose: Container for parsing outcome with model and errors
- Pattern: Always check `result.Success()` before accessing `result.Model`
- Includes: `Model`, `Errors`, `Tree` (optional), `Source` (file path)

**ElementKind:**
- Purpose: Discriminate element types at runtime (instead of reflection)
- Pattern: Use `elem.Kind()` for type checks and dispatch
- Examples: `KindPart`, `KindRequirement`, `KindVerification`, etc.

**Location:**
- Purpose: Track source position for error reporting and debugging
- Fields: `Line`, `Column`, `EndLine`, `EndColumn`
- Pattern: Available on all elements via `elem.Location()`

## Entry Points

**Command-Line Parsing (via examples):**
- Location: `gosysml2/examples/main.go`
- Triggers: `go run examples/main.go`
- Responsibilities:
  - Demonstrate both high and low-level APIs
  - Show model traversal via visitor
  - Show element counting and searching

**Library API Entry Points:**
- `sysml.ParseString(input, opts)` - Parse SysML text
- `sysml.ParseFile(path, opts)` - Parse SysML file
- `sysml.ParseDirectory(dir, opts)` - Sequential directory parsing
- `sysml.ParseDirectoryParallel(dir, workers, opts)` - Parallel parsing
- `sysml.ParseDirectoryStream(dir, handler, opts)` - Streaming parsing
- `low.Parse(input)` - Raw parse tree (performance critical)

## Error Handling

**Strategy:** Error collection during lexing/parsing, conversion to domain types, no exceptions (panic reserved for internal bugs only).

**Patterns:**

**Lexer/Parser Errors:**
```go
result := sysml.ParseFile("model.sysml")
if !result.Success() {
    for _, err := range result.Errors.Errors {
        fmt.Printf("Line %d, Col %d: %s\n", err.Line, err.Column, err.Message)
    }
}
```

**Error Types:**
- `ParseError`: Collection of `Error` structs with line/column/message
- `Error`: Single lexical, syntactic, or semantic error
- Multiple errors collected per parse (not fail-fast)

**Reference Resolution Errors:**
- Unresolved references do NOT generate errors
- Applications check `ref.IsResolved()` and handle gracefully
- No exception thrown for missing elements

## Cross-Cutting Concerns

**Logging:** No structured logging built-in. Applications use standard `fmt` or `log` packages.

**Validation:**
- Syntactic validation automatic during parsing
- Semantic validation (reference resolution) automatic post-parse
- No explicit validation API (all required validation done on parse)

**Authentication:** Not applicable (parser is stateless, no network/IO except file reading).

**Concurrency:**
- `ParseDirectoryParallel()` uses worker goroutines for I/O parallelism
- Safe: Each goroutine parses independently, no shared state during parsing
- Unsafe: Modifying a single model from multiple goroutines

**Memory Optimization:**
- `WithDiscardTree()` option discards ANTLR parse tree after model built (~30% reduction)
- Recommended for large repositories (>100MB)
- Trade-off: Cannot access raw parse tree if discarded

---

*Architecture analysis: 2026-01-30*
