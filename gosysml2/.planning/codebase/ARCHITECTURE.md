# Architecture

**Analysis Date:** 2026-02-05

## Pattern Overview

**Overall:** Layered Parser Architecture with Tree-to-Model Transformation

**Key Characteristics:**
- Two-layer API design: low-level ANTLR access and high-level model representation
- ANTLR4-based lexing and parsing with error collection
- Parse tree visitor pattern for converting ANTLR contexts to semantic model
- Lazy reference resolution: unresolved references stored during parsing, resolved post-parse
- Qualified name indexing for efficient element lookup
- Type-safe child element tracking within parent containers

## Layers

**Low Layer (ANTLR Foundation):**
- Purpose: Direct ANTLR lexer and parser access for performance-sensitive use cases
- Location: `internal/parser/`, `low/`
- Contains: Generated lexer/parser from g4 files, ANTLR wrapper (Lexer, Parser, error collection)
- Depends on: antlr4-go/antlr/v4, generated parser classes
- Used by: High layer parser, direct users needing parse tree access
- Key files: `low/lexer.go`, `low/parser.go`, `low/errors.go`

**High Layer (Semantic Model):**
- Purpose: Developer-friendly API for SysML model representation and traversal
- Location: `sysml/`
- Contains: Model elements (Package, Part, Requirement, etc.), model builders, visitors, helper functions
- Depends on: Low layer, ANTLR runtime
- Used by: User applications, examples
- Key files: `sysml/parse.go`, `sysml/model.go`, `sysml/visitor.go`

## Data Flow

**Parsing Pipeline:**

1. **Lexing** (`low/lexer.go`)
   - Input: SysML text (string or []byte)
   - ANTLR lexer tokenizes input
   - Errors collected via ErrorCollector interface

2. **Parsing** (`low/parser.go`)
   - Input: Token stream from lexer
   - ANTLR parser generates parse tree
   - Two modes: full tree construction or validation-only
   - Errors collected from parser error listener

3. **Tree Walking & Model Building** (`sysml/parse.go`, `sysml/model.go`)
   - Input: Parse tree (IEntryRuleRootNamespaceContext)
   - `modelBuilder` (embedded BaseSysMLv2ParserListener) walks tree
   - Creates high-level elements (Package, Part, Requirement, etc.)
   - Maintains stack of packages and elements for nesting context
   - Stores unresolved references (strings) for post-parse resolution

4. **Indexing** (`sysml/model.go` - BuildIndex)
   - Input: Complete model with all elements
   - Builds qualified-name index map for O(1) lookup
   - Used for reference resolution

5. **Reference Resolution** (`sysml/model.go` - ResolveReferences)
   - Input: Indexed model with unresolved reference strings
   - Walks all elements, resolves unresolved refs to actual elements
   - Sets bidirectional relationships (e.g., DerivedFrom/DerivedReqs)
   - Handles scope-relative lookup (parent chain walking)

**State Management:**

- **Parse Result**: `ParseResult` contains Model, Errors, raw parse tree, and source identifier
- **Element Tree**: Hierarchical with parent pointers and typed child collections
- **References**: Stored as `Ref[T]` generic type with lazy resolution
- **Unresolved State**: Kept as strings (line 572-575 in model.go for Requirements)

## Key Abstractions

**Element Interface** (`sysml/model.go`)
- Purpose: Common contract for all SysML model elements
- Provides: Kind, Name, QualifiedName, Location, Parent, Children, Documentation
- Implementations: Package, Part, Requirement, Port, etc. (25+ concrete types)
- Pattern: baseElement struct with common fields, type-specific fields added per concrete type

**Definition/Usage Duality** (`sysml/model.go`)
- Purpose: Distinguish type definitions from usages (SysML semantic distinction)
- Implementation: IsDefinition bool flag + TypeRef field (Ref to definition for usages)
- Examples: `part def Vehicle { }` (definition) vs `part vehicle1 : Vehicle;` (usage)

**Ref[T] Generic Reference** (`sysml/model.go`)
- Purpose: Represent element references that may be unresolved during parsing
- Pattern: Two-phase resolution (string name → actual element)
- Methods: Name(), Resolved(), IsResolved(), Resolve()
- Usage: TypeRef fields (Part TypeRef Ref[*Part]), cross-references (Requirement.Subject)

**Visitor Pattern** (`sysml/visitor.go`)
- Purpose: Extensible traversal and processing of model elements
- Implementation: Visitor interface with VisitX methods per element type
- BaseVisitor: Default implementation (all return true) for embedding
- Example: Counter visitor counts elements by kind

**Model Index** (`sysml/model.go`)
- Purpose: Fast element lookup by qualified name
- Implementation: map[string]Element built after parsing
- Used for: Reference resolution, external lookups
- Scope handling: Qualified names include package hierarchy (e.g., "P1::P2::Part1")

## Entry Points

**ParseString** (`sysml/parse.go:45`)
- Location: `sysml/parse.go`
- Signature: `ParseString(input string, opts ...ParseOption) *ParseResult`
- Triggers: User code needs to parse SysML from string
- Responsibilities: Delegates to parseWithSource, applies parse options

**ParseFile** (`sysml/parse.go:49`)
- Location: `sysml/parse.go`
- Signature: `ParseFile(filename string, opts ...ParseOption) *ParseResult`
- Triggers: User code needs to parse SysML from file
- Responsibilities: Reads file, delegates to parseWithSource with filename as source

**ParseDirectory** / **ParseDirectoryParallel** (`sysml/parse.go:147,171`)
- Location: `sysml/parse.go`
- Signature: `ParseDirectory(dir string, ...) ([]*ParseResult, error)`, `ParseDirectoryParallel(dir, workers, ...) ([]*ParseResult, error)`
- Triggers: Batch processing of multiple .sysml files
- Responsibilities: Walk directory, filter .sysml files, parse with optional parallelization

**ParseDirectoryStream** (`sysml/parse.go:212`)
- Location: `sysml/parse.go`
- Signature: `ParseDirectoryStream(dir string, handler func(*ParseResult) error, ...) error`
- Triggers: Memory-efficient processing of large file repositories
- Responsibilities: Parse files one-at-a-time, invoke handler, discard model after processing

## Error Handling

**Strategy:** Two-layer error representation with automatic conversion

**Low Layer** (`low/errors.go`):
- `SyntaxError`: Line, column, message, source (lexer/parser)
- `ErrorCollector`: Implements antlr.ErrorListener, accumulates errors during lexing/parsing
- `ParseErrors`: Container for LexerErrors and ParserErrors arrays

**High Layer** (`sysml/errors.go`):
- `Error`: User-friendly error with context field
- `ParseError`: Result object containing array of Errors, source identifier
- `Success()` method on ParseResult: Returns true if no errors
- Error conversion happens at parse boundary (convertFromLowLevel function)

**Patterns:**
- Don't panic on parse errors: return ParseResult with populated Errors
- Parse continues despite errors (error recovery)
- All errors collected: user sees complete error list, not first error only
- Line numbers available: for IDE integration and error location

## Cross-Cutting Concerns

**Logging:** None (library intentionally silent, errors via return values)

**Validation:**
- Pre-parse validation: `Validate(input)` and `ValidateFile(filename)` in sysml package
- Validation uses low-level parser with parse tree construction disabled for speed
- Post-parse validation: Implicit in reference resolution (unresolved refs don't fail, just stay unresolved)

**Authentication:** Not applicable (library parses local files/strings)

**Nesting & Scoping:**
- Packages: Stack-based tracking during tree walk (b.packageStack)
- Elements: Generic stack tracking (b.elementStack) for nested part definitions, requirement definitions
- getCurrentParent() method: Correctly resolves which element should receive new children
- Qualified names: Built via parent chain walking (::  separator per SysML spec)

**Performance Optimizations:**
- LazyError Collection: Tokens filled lazily (ANTLR lazy consumption)
- Memory Options: WithDiscardTree() option discards parse tree after model build
- Parallel Parsing: ParseDirectoryParallel with semaphore-based worker pool (default: runtime.NumCPU())
- Stream Processing: ParseDirectoryStream processes files without loading all into memory

---

*Architecture analysis: 2026-02-05*
