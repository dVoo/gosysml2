# Codebase Structure

**Analysis Date:** 2026-02-05

## Directory Layout

```
gosysml2/
├── internal/
│   └── parser/              # ANTLR-generated parser
├── low/                     # Low-level API
├── sysml/                   # High-level API
└── examples/                # Usage examples

code/
└── parser/                  # Standalone parser module

cmd/
├── verify-completeness/     # Model verification tool
├── verify-parser/           # Parser verification tool
├── test-attrs/              # Attribute test tool
├── test-low-level/          # Low-level API test tool
└── test-requirement-attributes/  # Requirement test tool

testdata/
├── simple/                  # Basic feature tests (30+ files)
├── Vehicle Example/         # Vehicle model examples
├── Requirements Examples/   # Requirements modeling
├── Import Tests/            # Import functionality tests
└── Variability Examples/    # Variability modeling

examples/                    # User example files
├── test_example_req_verification.sysml
├── test_corrected.sysml
├── test_user_input.sysml
└── test_req_custom.sysml

docs/
└── bnf/                     # BNF grammar documentation
```

## Directory Purposes

### gosysml2/

**Purpose:** Main library module

**Contains:**
- `internal/parser/` - ANTLR-generated lexer and parser (~2.2MB generated code)
- `low/` - Low-level parsing API with direct ANTLR access
- `sysml/` - High-level model API with visitor pattern
- `examples/` - Library usage examples

**Key files:**
- `go.mod` - Module definition (`github.com/dVoo/gosysml2`)

### gosysml2/internal/parser/

**Purpose:** Generated parser code (should not be edited manually)

**Contains:**
- `sysmlv2_lexer.go` - Lexer (~86KB)
- `sysmlv2_parser.go` - Parser (~2.1MB, main grammar implementation)
- `sysmlv2parser_listener.go` - Listener interface
- `sysmlv2parser_base_listener.go` - Base listener
- `.interp`, `.tokens` files - ANTLR metadata

**Pattern:** Generated from grammar, never edit directly

### gosysml2/low/

**Purpose:** Low-level API with minimal overhead

**Contains:**
- `lexer.go` - Lexer wrapper with error collection
- `parser.go` - Parser wrapper with configuration
- `errors.go` - Error types and collection
- `parser_test.go` - Unit tests

**Design:** Thin wrapper around ANTLR for performance-critical use cases

### gosysml2/sysml/

**Purpose:** High-level developer-friendly API

**Contains:**
- `model.go` - All domain types (Model, Package, Part, Requirement, etc.)
- `parse.go` - Parsing facade and model builder (~1733 lines)
- `visitor.go` - Visitor pattern implementation and traversal utilities
- `errors.go` - User-friendly error types
- `*_test.go` - Unit and integration tests

**Design:** Rich domain model with reference resolution and traversal

### cmd/

**Purpose:** CLI tools for testing and verification

**Pattern:** Each subdirectory is a standalone main package

**Tools:**

| Directory | Purpose | Key Function |
|-----------|---------|--------------|
| `verify-parser/` | Analyze parse tree coverage | Source coverage checking |
| `verify-completeness/` | Verify model extraction | Element statistics, unresolved refs |
| `test-low-level/` | Test low-level API | Token analysis |
| `test-attrs/` | Test attribute parsing | Attribute extraction |
| `test-requirement-attributes/` | Test requirements | Requirement attribute tests |

### testdata/

**Purpose:** SysML test files for validation

**Structure:**
- `simple/` - 30+ files testing individual features
- `Vehicle Example/` - Complex multi-file vehicle model
- `Requirements Examples/` - Requirements derivation, HSUV model
- `Import Tests/` - Import functionality (alias, qualified, circular)
- `Variability Examples/` - Variability modeling

**Pattern:** Each `.sysml` file tests specific language features

### examples/

**Purpose:** User-facing example SysML files

**Contains:**
- Requirement verification examples
- User input test cases
- Corrected/validated models

### docs/

**Purpose:** Documentation assets

**Contains:**
- `bnf/` - BNF grammar specification and images

## Key File Locations

### Entry Points

| Purpose | Path |
|---------|------|
| Library root | `gosysml2/sysml/parse.go` |
| Example main | `gosysml2/examples/main.go` |
| Parser verify | `cmd/verify-parser/main.go` |
| Completeness verify | `cmd/verify-completeness/main.go` |

### Configuration

| File | Purpose |
|------|---------|
| `go.mod` | Root module: `github.com/dVoo/sysmlv2-tools` |
| `gosysml2/go.mod` | Library module: `github.com/dVoo/gosysml2` |
| `code/parser/go.mod` | Standalone parser: `module verify` |

### Core Logic

| File | Responsibility |
|------|----------------|
| `gosysml2/sysml/model.go` | All domain types (lines ~2000+) |
| `gosysml2/sysml/parse.go` | Parsing facade, model builder (lines ~1733) |
| `gosysml2/sysml/visitor.go` | Visitor pattern, traversal (lines ~495) |
| `gosysml2/low/parser.go` | Low-level parser wrapper |
| `gosysml2/low/lexer.go` | Low-level lexer wrapper |
| `gosysml2/low/errors.go` | Error types |

### Testing

| File | Type |
|------|------|
| `gosysml2/sysml/parse_test.go` | Unit tests for parsing |
| `gosysml2/sysml/visitor_test.go` | Unit tests for visitor |
| `gosysml2/sysml/integration_test.go` | Integration tests |
| `gosysml2/low/parser_test.go` | Low-level parser tests |

## Naming Conventions

### Files

| Pattern | Example | Purpose |
|---------|---------|---------|
| `*.go` | `model.go` | Implementation files |
| `*_test.go` | `parse_test.go` | Test files |
| `sysmlv2_*.go` | `sysmlv2_parser.go` | Generated files |
| `main.go` | `main.go` | CLI entry points |

### Directories

| Pattern | Example | Purpose |
|---------|---------|---------|
| `cmd/<tool-name>/` | `cmd/verify-parser/` | CLI tools (kebab-case) |
| `<category> Examples/` | `Requirements Examples/` | Test data (spaces allowed) |
| `<feature> Tests/` | `Import Tests/` | Feature-specific tests |

### Go Identifiers

| Pattern | Example | Usage |
|---------|---------|-------|
| PascalCase | `PartDefinition`, `NewRequirement` | Exported types/functions |
| camelCase | `elementStack`, `unresolvedSubject` | Unexported fields |
| Kind prefix | `KindPackage`, `KindRequirement` | Element kind constants |

## Where to Add New Code

### New Element Type

1. **Add to `gosysml2/sysml/model.go`:**
   - Add `KindXXX` constant to `ElementKind` enum
   - Define struct type (embed `baseElement`)
   - Implement `isDefinition()` and/or `isUsage()` if applicable
   - Add `AddChild()` with type tracking
   - Add accessor methods for typed children

2. **Add to `gosysml2/sysml/parse.go`:**
   - Add `EnterXxxDefinition()` method to `modelBuilder`
   - Add `EnterXxxUsage()` method to `modelBuilder`
   - Add `resolveXxxRefs()` method to `Model`
   - Call resolver from `ResolveReferences()`

3. **Add to `gosysml2/sysml/visitor.go`:**
   - Add `VisitXxx()` to `Visitor` interface
   - Add `VisitXxx()` to `BaseVisitor`
   - Add case to `visitElement()`
   - Add `FindXxx()` helper function
   - Add to `Counter` visitor

### New Parse Option

1. **Add to `gosysml2/sysml/parse.go`:**
   - Add field to `parseConfig` struct
   - Create `WithXxx()` option function
   - Use option in `parseWithSource()`

### New CLI Tool

1. **Create `cmd/<tool-name>/main.go`:**
   - Use `package main`
   - Import `github.com/dVoo/gosysml2/sysml`
   - Handle command-line args
   - Print results

2. **Update root `go.mod`:**
   - Add tool as dependency if needed

### New Test Data

1. **Add `.sysml` file to appropriate `testdata/` subdirectory:**
   - `testdata/simple/` for single-feature tests
   - `testdata/<Category> Examples/` for complex scenarios

## Special Directories

### internal/

**Purpose:** Internal implementation details

**Characteristics:**
- Go visibility: accessible only within module
- Contains generated code (should not be edited)
- Breaking changes OK (internal API)

### cmd/

**Purpose:** Standalone executables

**Characteristics:**
- Each subdirectory is independent `package main`
- Built with `go build ./cmd/<name>/`
- Not part of library API

### testdata/

**Purpose:** Test fixtures

**Characteristics:**
- SysML source files for testing
- Organized by feature/semantic area
- Referenced by tests using relative paths
- Committed to repository

---

*Structure analysis: 2026-02-05*
