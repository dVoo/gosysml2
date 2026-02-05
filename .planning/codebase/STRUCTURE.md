# Codebase Structure

**Analysis Date:** 2026-02-05

## Directory Layout

```
gosysml2/                          # Main library module
├── examples/                      # Example usage code
│   └── main.go                   # Demo of both APIs
├── internal/                      # Internal implementation
│   └── parser/                   # ANTLR-generated parser
│       ├── sysmlv2_lexer.go     # Generated lexer (~85KB)
│       ├── sysmlv2_parser.go    # Generated parser (~2.1MB)
│       ├── sysmlv2parser_listener.go      # Generated listener interface
│       ├── sysmlv2parser_base_listener.go # Generated base listener
│       ├── SysMLv2Lexer.interp   # ANTLR interpreter data
│       ├── SysMLv2Parser.interp  # ANTLR interpreter data
│       ├── SysMLv2Lexer.tokens   # Token definitions
│       └── SysMLv2Parser.tokens  # Token definitions
├── low/                          # Low-level API package
│   ├── lexer.go                 # Lexer wrapper
│   ├── parser.go                # Parser wrapper
│   ├── errors.go                # Error types and collection
│   └── parser_test.go           # Low-level parser tests
├── sysml/                        # High-level API package
│   ├── model.go                 # Domain model types (~2000 lines)
│   ├── parse.go                 # Parsing functions (~1700 lines)
│   ├── visitor.go               # Visitor pattern implementation
│   ├── errors.go                # High-level error types
│   ├── parse_test.go            # Parse function tests
│   ├── visitor_test.go          # Visitor tests
│   └── integration_test.go      # Integration tests
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── README.md                    # Documentation
└── PERFORMANCE.md               # Performance documentation

testdata/                          # Test data files (outside module)
├── simple/                       # Simple test cases
├── Vehicle Example/              # Vehicle model examples
├── Requirements Examples/        # Requirements test cases
├── Variability Examples/         # Variability test cases
└── Import Tests/                 # Import test cases

.planning/                         # Planning documentation
├── codebase/                     # This directory
│   ├── ARCHITECTURE.md
│   └── STRUCTURE.md
├── phases/                       # Implementation phases
└── ROADMAP.md

analysis/                          # Analysis documents
├── ATTRIBUTE_FIX_SUMMARY.md
├── EXAMPLE_VERIFICATION_ANALYSIS.md
├── CUSTOM_ATTRIBUTES_ANALYSIS.md
└── PARSER_COMPLETENESS_REPORT.md

docs/                             # Documentation
└── bnf/                         # BNF grammar documentation
```

## Directory Purposes

**`gosysml2/` (Library Root):**
- Purpose: Go module root containing the library code
- Contains: Source packages, examples, module files
- Key files: `go.mod`, `go.sum`, `README.md`

**`gosysml2/internal/parser/`:**
- Purpose: Auto-generated ANTLR4 parser code
- Generated: Yes (from SysML v2 grammar)
- Committed: Yes (for distribution without ANTLR dependency)
- Do not edit manually - regenerate from grammar

**`gosysml2/low/`:**
- Purpose: Low-level performance-oriented API
- Audience: Advanced users needing parse tree access
- Pattern: Thin wrappers around ANTLR with error collection

**`gosysml2/sysml/`:**
- Purpose: High-level developer-friendly API
- Audience: Most library users
- Pattern: Domain model with visitor pattern

**`gosysml2/examples/`:**
- Purpose: Example usage demonstrating both APIs
- Usage: Reference for library consumers

**`testdata/`:**
- Purpose: Test input files in SysML format
- Structure: Organized by feature/example type
- Usage: Integration tests, manual testing

## Key File Locations

**Entry Points:**
- `gosysml2/examples/main.go`: Example application
- `gosysml2/sysml/parse.go`: High-level parsing API (line 1)
- `gosysml2/low/parser.go`: Low-level parsing API (line 1)

**Configuration:**
- `gosysml2/go.mod`: Module definition (Go 1.22, antlr4-go/antlr/v4)

**Core Logic:**
- `gosysml2/sysml/model.go`: Domain model (~2000 lines)
- `gosysml2/sysml/parse.go`: Parse functions and model builder (~1700 lines)
- `gosysml2/sysml/visitor.go`: Visitor pattern (495 lines)

**Testing:**
- `gosysml2/sysml/parse_test.go`: Parse function tests
- `gosysml2/sysml/visitor_test.go`: Visitor tests
- `gosysml2/sysml/integration_test.go`: Integration tests
- `gosysml2/low/parser_test.go`: Low-level parser tests

**Error Handling:**
- `gosysml2/sysml/errors.go`: High-level errors
- `gosysml2/low/errors.go`: Low-level errors

## Naming Conventions

**Files:**
- Package name matches directory: `sysml/`, `low/`
- Test files: `*_test.go` suffix
- Generated files: Named by ANTLR generator

**Types:**
- Public types: PascalCase (`ParseResult`, `ElementKind`)
- Private types: camelCase (`modelBuilder`, `baseElement`)
- Interface suffix: None (`Element`, `Visitor`, `Definition`)
- Generic types: `Ref[T]` pattern

**Functions:**
- Constructor pattern: `New` + Type (`NewModel()`, `NewPackage()`)
- Parse functions: `Parse` + Input type (`ParseString()`, `ParseFile()`)
- Find functions: `Find` + Element type (`FindParts()`, `FindRequirements()`)
- Must functions: `Must` + Verb (`MustParseString()`) - panic on error

**Variables:**
- Stack tracking: `elementStack`, `packageStack`
- Unresolved references: `unresolved` + Reference type (`unresolvedSubject`)
- Collectors: plural nouns (`errors`, `tokens`)

## Where to Add New Code

**New Element Type:**
1. Add `ElementKind` constant in `sysml/model.go`
2. Add `String()` case in `sysml/model.go`
3. Create struct type in `sysml/model.go`
4. Implement `Element` interface (embed `baseElement`)
5. Add `isDefinition()` and/or `isUsage()` methods if applicable
6. Add visitor method in `sysml/visitor.go` interface and `BaseVisitor`
7. Add visit dispatch in `visitElement()` function
8. Add find function in `sysml/visitor.go`
9. Add builder methods in `sysml/parse.go` (in `modelBuilder`)
10. Add reference resolution in `model.ResolveReferences()`

**New Parse Function:**
- Add to `sysml/parse.go` near existing `Parse*` functions
- Follow pattern: parse → build model → index → resolve → return
- Add corresponding test in `sysml/parse_test.go`

**New Visitor:**
- Embed `BaseVisitor` to only override methods of interest
- Implement `Visitor` interface
- Use `sysml.Visit()` or `sysml.Walk()` to traverse model

**New Test:**
- Package tests: `sysml/*_test.go` or `low/*_test.go`
- Test data: Add `.sysml` files to appropriate `testdata/` subdirectory
- Follow naming: `Test` + FunctionName

## Special Directories

**`internal/parser/`:**
- Purpose: Auto-generated ANTLR4 code
- Generated: Yes (via ANTLR tool)
- Committed: Yes
- Do not modify directly

**`testdata/`:**
- Purpose: Test fixtures
- Generated: No
- Committed: Yes
- Format: Raw SysML v2 files

**`.planning/`:**
- Purpose: Planning and analysis documentation
- Generated: No
- Committed: Yes
- Usage: GSD command integration

---

*Structure analysis: 2026-02-05*
