# Codebase Structure

**Analysis Date:** 2026-01-30

## Directory Layout

```
claudecode/
├── code/                      # ANTLR4 grammar files and generated parser
│   ├── SysMLv2Lexer.g4        # Shared lexer (keywords, operators, literals)
│   ├── KerMLParser.g4         # KerML parser rules (foundation)
│   ├── SysMLv2Parser.g4       # SysML parser rules (extends KerML)
│   └── parser/                # Generated Go parser code
│       └── *.go               # ANTLR-generated lexer, parser, listeners
├── docs/                      # Specification and reference documentation
│   └── bnf/                   # BNF/EBNF grammar specifications
│       ├── KerML-textual-bnf.kebnf
│       ├── SysML-textual-bnf.kebnf
│       ├── SysML-graphical-bnf.kgbnf
│       ├── *.html             # Rendered grammar documentation
│       └── images/            # SVG diagrams for graphical notation (284 files)
├── gosysml2/                  # Go SysML v2 parser library
│   ├── go.mod                 # Go module definition
│   ├── go.sum                 # Dependency checksums
│   ├── README.md              # Library documentation and API reference
│   ├── PERFORMANCE.md         # Performance tuning and capacity guidelines
│   ├── internal/
│   │   └── parser/            # ANTLR-generated Go code (do not edit)
│   ├── low/                   # Low-level lexer/parser API
│   │   ├── lexer.go           # Lexer wrapper with error collection
│   │   ├── parser.go          # Parser wrapper with token stream mgmt
│   │   ├── errors.go          # Error collector and error types
│   │   └── *_test.go          # Low-level parsing tests
│   ├── sysml/                 # High-level Go model API
│   │   ├── model.go           # Element types, kinds, Model struct
│   │   ├── parse.go           # Parsing entry points (string/file/dir)
│   │   ├── visitor.go         # Visitor pattern and traversal
│   │   ├── errors.go          # High-level error handling
│   │   └── *_test.go          # Model and integration tests
│   └── examples/
│       └── main.go            # Complete usage examples
└── testdata/                  # SysML test files for grammar validation
    ├── simple/                # Basic test files (ActionTest, etc.)
    ├── Vehicle Example/       # Complete vehicle model example
    ├── Import Tests/
    ├── Requirements Examples/
    └── Variability Examples/
```

## Directory Purposes

**code/:**
- Purpose: Grammar specifications in ANTLR4 format
- Contains: Lexer and parser rules, generated code output
- Key files: `SysMLv2Lexer.g4`, `SysMLv2Parser.g4`, `KerMLParser.g4`
- Generated: Yes (parser/ contains ~91k lines of ANTLR output)
- Edit: Only `.g4` files; parser/ is generated

**docs/bnf/:**
- Purpose: Formal BNF/EBNF grammar documentation and visual specifications
- Contains: Grammar source (`.kebnf`, `.kgbnf`), rendered HTML, SVG diagrams
- Key files: `KerML-textual-bnf.kebnf`, `SysML-textual-bnf.kebnf`
- Committed: Yes (reference documentation)
- Generated: HTML and images are generated from source

**gosysml2/:**
- Purpose: Go library for parsing and traversing SysML v2 models
- Contains: Parser library code, tests, examples, documentation
- Module: `github.com/dVoo/gosysml2`
- Requires: Go 1.22+, `github.com/antlr4-go/antlr/v4`

**gosysml2/internal/parser/:**
- Purpose: ANTLR-generated Go parser implementation (do not edit)
- Contains: Lexer, parser, listener interfaces, base listeners (~85k lines)
- Generated: Yes (via `antlr -Dlanguage=Go -o parser ...`)
- Regenerate: Run ANTLR when grammar files change

**gosysml2/low/:**
- Purpose: Low-level parser API for performance-critical applications
- Contains: `Lexer`, `Parser` wrappers with error collection
- Key functions: `NewLexer()`, `NewParser()`, `Parse()`, `Validate()`
- Returns: ANTLR parse trees (raw parse tree structures)
- Usage: Performance tuning, advanced use cases, debugging

**gosysml2/sysml/:**
- Purpose: High-level idiomatic Go API for SysML model consumption
- Contains: Element types (`Package`, `Part`, `Requirement`, etc.), parsing functions, visitor pattern
- Key functions: `ParseString()`, `ParseFile()`, `ParseDirectory()`, `Visit()`, `Walk()`
- Returns: `ParseResult` with `Model` and `Errors`
- Usage: Most applications use this API

**gosysml2/examples/:**
- Purpose: Demonstrate library usage (both APIs)
- File: `main.go` shows high-level and low-level examples
- Run: `go run examples/main.go`

**testdata/:**
- Purpose: SysML test files covering grammar features
- Contains: 34 `.sysml` files exercising all language constructs
- Organization:
  - `simple/` - Individual feature tests (ActionTest.sysml, etc.)
  - `Vehicle Example/` - Complete model with multiple files
  - `Requirements Examples/` - Requirement-focused examples
  - Other categories for imports, variability, etc.

## Key File Locations

**Entry Points:**
- `gosysml2/examples/main.go`: Complete API usage examples (run with `go run`)
- `gosysml2/sysml/parse.go`: Public parsing API functions
- `gosysml2/low/parser.go`: Low-level parsing wrapper

**Configuration:**
- `gosysml2/go.mod`: Module name and dependencies
- `code/SysMLv2Lexer.g4`: Lexer configuration (keywords, operators)
- `CLAUDE.md`: Project guidance (development environment, structure)

**Core Logic:**
- `gosysml2/sysml/model.go`: All element types and model structure (~1672 lines)
- `gosysml2/sysml/parse.go`: Parse orchestration, tree-to-model conversion (~612 lines)
- `gosysml2/sysml/visitor.go`: Model traversal via visitor pattern (~374 lines)
- `gosysml2/low/parser.go`: Lexer/parser initialization and token management

**Testing:**
- `gosysml2/sysml/parse_test.go`: Parse function tests
- `gosysml2/sysml/visitor_test.go`: Visitor and model traversal tests
- `gosysml2/sysml/integration_test.go`: End-to-end parsing tests
- `gosysml2/low/parser_test.go`: Low-level parser tests
- `testdata/`: Real SysML files for validation

## Naming Conventions

**Files:**
- Grammar files: `[Name]v2Lexer.g4`, `[Name]Parser.g4` (ANTLR standard)
- Generated files: `[name]_lexer.go`, `[name]_parser.go` (ANTLR standard)
- Test files: `*_test.go` (Go convention)
- Model files: `model.go`, `parse.go`, `visitor.go` (function-based organization)
- SysML files: `*.sysml` (test data)

**Directories:**
- Package names lowercase: `sysml`, `low`, `parser`
- Internal packages: `internal/parser/` (generated code, not exported)
- Example code: `examples/` (not in main packages)
- Test data: `testdata/` (repository root)

**Types (in gosysml2/sysml/):**
- Model types: PascalCase (`Package`, `Part`, `Requirement`, `Verification`)
- Type enums: `KindPart`, `KindRequirement` (prefix with `Kind`)
- Interfaces: `Element`, `Visitor` (PascalCase)
- Functions: `ParseString()`, `ParseFile()`, `Visit()` (PascalCase, verb-first)
- Internal helpers: `buildModel()`, `visitElement()` (camelCase)

**Constants:**
- Element kinds: `KindPackage`, `KindPart`, `KindRequirement` (screaming KindCase)
- No magic strings; all use constants

## Where to Add New Code

**New SysML Element Type:**
1. Define struct in `gosysml2/sysml/model.go` (embed `BaseElement`)
2. Add `ElementKind` constant to `model.go` (e.g., `KindMyElement`)
3. Add `String()` case to `ElementKind.String()` method
4. Add visitor method in `gosysml2/sysml/visitor.go` (e.g., `VisitMyElement()`)
5. Add case to `visitElement()` dispatcher
6. Add parsing logic in tree-to-model visitor (in `parse.go`)
7. Add test in `gosysml2/sysml/*_test.go`

**New Parsing Feature:**
1. Add grammar rules to `code/SysMLv2Parser.g4` (or `KerMLParser.g4` if KerML-specific)
2. Regenerate parser: `antlr -Dlanguage=Go -o parser SysMLv2Lexer.g4 SysMLv2Parser.g4`
3. Copy generated files to `gosysml2/internal/parser/`
4. Add element type to `gosysml2/sysml/model.go` (step 1 above)
5. Update tree visitor in `parse.go` to handle new rule contexts
6. Add test in `testdata/` with example SysML
7. Verify with `gosysml2/low/parser_test.go` or integration test

**New Parsing Function:**
- Location: `gosysml2/sysml/parse.go`
- Pattern: Follow existing functions (`ParseString`, `ParseFile`, etc.)
- Return: `*ParseResult` (contains Model, Errors, Tree)
- Handle options: Accept `...ParseOption` parameter

**Error Handling Extension:**
- Location: `gosysml2/sysml/errors.go` (high-level)
- Location: `gosysml2/low/errors.go` (low-level)
- Pattern: Return `*ParseError` containing `[]*Error` slices
- Never use `panic()` for user input errors

**Utilities/Helpers:**
- Shared functions: `gosysml2/sysml/` package (exported as public)
- Internal helpers: Same package with lowercase names
- Reusable across packages: Consider moving to separate file

## Special Directories

**gosysml2/internal/parser/:**
- Purpose: ANTLR-generated code (do NOT edit manually)
- Generated: Yes (from grammar files via ANTLR)
- Committed: Yes (to avoid rebuild requirement on clone)
- Regenerate when: Grammar files change
- Command: `cd code && antlr -Dlanguage=Go -o parser SysMLv2Lexer.g4 SysMLv2Parser.g4`

**code/parser/:**
- Purpose: Alternative output location for ANTLR generation
- Status: May contain stale generated code from development
- Canonical: Use `gosysml2/internal/parser/` instead

**testdata/:**
- Purpose: SysML files for validation and testing
- Committed: Yes (for regression testing)
- Organization: By feature category (simple/, Vehicle Example/, etc.)
- Usage: Loaded by test functions to validate parsing

---

*Structure analysis: 2026-01-30*
