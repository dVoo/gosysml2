# Codebase Structure

**Analysis Date:** 2026-02-05

## Directory Layout

```
gosysml2/
├── .planning/              # Planning documents (not part of build)
├── code/                   # ANTLR4 grammar files (from OMG SysML spec)
│   ├── SysMLv2Lexer.g4
│   ├── SysMLv2Parser.g4
│   └── parser/             # Generated parser code (checked in)
├── examples/               # Example usage
│   └── main.go
├── internal/
│   └── parser/             # Generated ANTLR code (from antlr tool)
│       ├── sysmlv2_lexer.go
│       ├── sysmlv2_parser.go
│       ├── sysmlv2parser_listener.go
│       └── sysmlv2parser_base_listener.go
├── low/                    # Low-level ANTLR wrapper API
│   ├── lexer.go           # Lexer wrapper with error collection
│   ├── parser.go          # Parser wrapper with error collection
│   ├── errors.go          # Error types (SyntaxError, ErrorCollector)
│   ├── lexer_test.go
│   └── parser_test.go
├── sysml/                  # High-level semantic model API
│   ├── parse.go           # Main entry points (ParseString, ParseFile, etc.)
│   ├── model.go           # Model structure, element types, indexing, resolution
│   ├── visitor.go         # Visitor interface and BaseVisitor implementation
│   ├── errors.go          # User-friendly error types (Error, ParseError)
│   ├── parse_test.go
│   ├── integration_test.go
│   └── visitor_test.go
├── go.mod                  # Module definition
├── go.sum                  # Dependency checksums
└── README.md               # Project documentation
```

## Directory Purposes

**code/**
- Purpose: ANTLR4 grammar source files defining SysML v2 syntax
- Contains: Lexer rules (SysMLv2Lexer.g4), parser rules (SysMLv2Parser.g4)
- Key files: Grammar files reference OMG specification clauses
- Generated into: `internal/parser/`

**internal/parser/**
- Purpose: Generated ANTLR parser classes (checked into version control)
- Contains: Lexer token classes, parser context classes, listener base classes
- Key files:
  - `sysmlv2_lexer.go`: Token definitions and lexing logic
  - `sysmlv2_parser.go`: Parse rules and parse tree construction
  - `sysmlv2parser_base_listener.go`: BaseSysMLv2ParserListener for tree walking
- Generated from: ANTLR tool with `antlr -Dlanguage=Go code/SysMLv2Lexer.g4 code/SysMLv2Parser.g4`
- Committed: Yes (for reproducible builds without ANTLR tool)

**low/**
- Purpose: Low-level performance-oriented API wrapping ANTLR
- Contains: Lexer, Parser, ErrorCollector abstractions around generated code
- Key classes:
  - `Lexer`: Wraps SysMLv2Lexer with error collection
  - `Parser`: Wraps SysMLv2Parser, provides parse tree construction options
  - `ErrorCollector`: Implements antlr.ErrorListener interface
- Used by: High-level sysml package, advanced users needing parse tree access
- Testing: Unit tests verify lexer/parser directly (`lexer_test.go`, `parser_test.go`)

**sysml/**
- Purpose: High-level semantic model and user-facing API
- Contains: 25+ element types, model builder, visitor pattern, helper functions
- Key types:
  - `Element` interface: Base for all model elements
  - `baseElement`: Common fields (kind, name, location, parent, children)
  - Concrete types: Package, Part, Requirement, Port, Action, State, View, etc.
  - `Ref[T]`: Generic reference type with lazy resolution
  - `Model`: Root container with packages, imports, comments
- Key functions:
  - `ParseString`, `ParseFile`, `ParseDirectory`, `ParseDirectoryParallel`, `ParseDirectoryStream`
  - `Validate`, `ValidateFile`: Parse without building full model
  - `MustParseString`, `MustParseFile`: Parse or panic
  - `Walk`: Depth-first traversal with callback
  - `Visit`: Apply Visitor pattern traversal
- Testing: Full test suite with parse, visitor, integration tests

**examples/**
- Purpose: Demonstrates library usage for both APIs
- Contains: Single main.go showing:
  - High-level API: ParseString, Walk, FindParts, Visit
  - Low-level API: low.Parse, parser.TokenCount
  - Error handling patterns
- Not built/deployed: Serves as documentation and manual test

## Key File Locations

**Entry Points:**
- `sysml/parse.go:45`: ParseString - main entry for string input
- `sysml/parse.go:49`: ParseFile - main entry for file input
- `sysml/parse.go:147`: ParseDirectory - batch parsing
- `sysml/parse.go:171`: ParseDirectoryParallel - parallel batch parsing
- `sysml/parse.go:212`: ParseDirectoryStream - streaming batch parsing
- `low/parser.go:110`: low.Parse - low-level parse entry point

**Core Model:**
- `sysml/model.go:1456`: Model struct (root)
- `sysml/model.go:227`: Package element
- `sysml/model.go:361`: Part element
- `sysml/model.go:548`: Requirement element
- `sysml/model.go:421`: Port element
- `sysml/model.go:700`: Verification element
- `sysml/model.go:147`: Element interface definition

**Builder & Resolution:**
- `sysml/parse.go:240`: buildModel function (tree → Model)
- `sysml/parse.go:257`: modelBuilder struct (implements visitor pattern)
- `sysml/model.go:1513`: BuildIndex method (creates qualified name index)
- `sysml/model.go:1526`: ResolveReferences method (resolves unresolved refs)
- `sysml/model.go:1576`: resolveRequirementRefs, resolveVerificationRefs, etc. (per-type resolution)

**Visitors & Traversal:**
- `sysml/visitor.go:5`: Visitor interface (dispatch methods per element type)
- `sysml/visitor.go:83`: BaseVisitor (default implementation for embedding)

**Error Types:**
- `sysml/errors.go:12`: Error struct (user-friendly)
- `sysml/errors.go:27`: ParseError struct (parse result with errors array)
- `low/errors.go:13`: SyntaxError struct (low-level)
- `low/errors.go:26`: ErrorCollector struct (implements antlr.ErrorListener)

**Testing:**
- `sysml/parse_test.go`: Unit tests for parsing functions
- `sysml/integration_test.go`: Integration tests for full workflows
- `sysml/visitor_test.go`: Visitor pattern tests
- `low/parser_test.go`: Low-level parser tests

## Naming Conventions

**Files:**
- `*_test.go`: Go test files (tested with `go test`)
- `*_errors.go`: Error types (separate from main logic)
- `parse.go`, `parser.go`: Parser-related code
- `model.go`: Model structure and element definitions
- `visitor.go`: Visitor pattern implementation
- `lexer.go`: Lexer wrapper

**Directories:**
- `internal/`: Unexported packages (internal/parser), not for direct client use
- `code/`: Source files (grammars), not compiled code
- `examples/`: Documentation-by-example, not deployed
- Package names lowercase: `sysml`, `low` (no underscores)

**Functions:**
- `Parse*`: Public parsing functions (ParseString, ParseFile)
- `parse*`: Private helpers (parseWithSource, parseConfig)
- `New*`: Constructor functions (NewModel, NewPackage, NewPart)
- `Resolve*`: Reference resolution (ResolveReferences, resolveRequirementRefs)
- `Walk*`: Tree traversal (Walk, walkElement)
- `Enter*/Exit*`: ANTLR listener callbacks (EnterPackage_, ExitPackage_)

**Types:**
- `*Element`: Concrete element types (Part, Requirement, Port)
- `*Definition`, `*Usage`: Semantically significant naming per SysML
- `ParseResult`, `ParseError`: Result/error types
- `Visitor`, `BaseVisitor`: Interface and implementation

## Where to Add New Code

**New SysML Element Type:**
1. Define struct in `sysml/model.go` (e.g., `type MyElement struct { baseElement ... }`)
2. Implement Element interface methods (Kind, Name, QualifiedName, etc.)
3. Add constructor `NewMyElement(...)` in same file
4. Add listener methods in `sysml/parse.go`: `EnterMyElementDefinition`, `ExitMyElementDefinition`
5. Add resolution method in `sysml/model.go`: `resolveMyElementRefs`
6. Call resolution in `Model.ResolveReferences()` dispatch
7. Add visitor method in `sysml/visitor.go`: `VisitMyElement(elem *MyElement) bool`
8. Add test in `sysml/*_test.go` or new test file

**New Grammar Rule:**
1. Add rule to `code/SysMLv2Parser.g4` or `code/SysMLv2Lexer.g4`
2. Regenerate parser: `antlr -Dlanguage=Go -o internal/parser code/SysMLv2Lexer.g4 code/SysMLv2Parser.g4`
3. Update generated files in `internal/parser/`
4. Implement listener method in model builder (step 4 above)

**New Public API Function:**
- Place in `sysml/parse.go` if parsing-related
- Place in `sysml/model.go` if model query/traversal-related
- Export uppercase (e.g., `FindParts`, `Walk`)
- Add tests in appropriate test file

**New Visitor Implementation:**
1. Embed `BaseVisitor` in your struct
2. Override only Visit* methods you need
3. Example: Counter visitor in examples or tests

## Special Directories

**internal/parser/**
- Purpose: Hold ANTLR-generated code
- Generated: Yes, from antlr tool
- Committed: Yes (checked into git)
- Modified: Never manually edit (regenerate from grammar changes)

**code/**
- Purpose: Source grammars
- Generated: No (hand-written BNF/EBNF)
- Committed: Yes
- Modified: Only when SysML syntax changes

**.planning/codebase/**
- Purpose: Architecture/structure documentation
- Generated: By gsd mapper (you are here)
- Committed: Yes
- Modified: Updated when architecture changes significantly

---

*Structure analysis: 2026-02-05*
