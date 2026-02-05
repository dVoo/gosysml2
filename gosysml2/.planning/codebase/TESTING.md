# Testing Patterns

**Analysis Date:** 2026-02-05

## Test Framework

**Runner:**
- Go's built-in `testing` package (`testing.T`)
- Run tests with: `go test ./...` or `go test ./sysml` or `go test ./low`
- No external test framework detected (no Jest, Vitest, or third-party libraries)

**Assertion Library:**
- Manual assertions using idiomatic Go patterns (no assertion library)
- Pattern: `if condition != expected { t.Errorf("message: got %v, want %v", got, expected) }`
- Fatal assertions: `t.Fatalf()` stops test immediately
- Non-fatal assertions: `t.Errorf()` continues test execution

**Run Commands:**
```bash
go test ./...                    # Run all tests
go test ./sysml                  # Run tests in sysml package
go test ./low                    # Run tests in low package
go test -v ./...                 # Verbose output
go test -run TestParseString ./  # Run specific test by name pattern
go test -cover ./...             # Show coverage percentage
go test -coverprofile=cov.out ./... && go tool cover -html=cov.out  # HTML coverage
```

## Test File Organization

**Location:**
- Co-located with source code: `parse.go` and `parse_test.go` in same directory
- Tests grouped by package: `/home/daniel/projects/claudecode/gosysml2/sysml/*_test.go`, `/home/daniel/projects/claudecode/gosysml2/low/*_test.go`

**Naming:**
- Test files: `*_test.go` (Go convention)
- Test functions: `TestXxx` where `Xxx` describes what is tested
- Examples: `TestParseString`, `TestParseStringWithErrors`, `TestValidate`, `TestWalk`, `TestFindByKind`, `TestCounter`, `TestMustParseStringPanic`

**Structure:**
```
gosysml2/
├── low/
│   ├── lexer.go
│   ├── parser_test.go        # Tests for Parser, Lexer, ErrorCollector
│   ├── errors.go
│   └── parser.go
├── sysml/
│   ├── parse_test.go         # Tests for ParseString, ParseFile, Validate, Walk
│   ├── integration_test.go    # Tests for ParseDirectory variants, file parsing
│   ├── visitor_test.go        # Tests for Find* functions (Requirements, Parts, etc.)
│   ├── parse.go
│   ├── model.go
│   ├── visitor.go
│   └── errors.go
└── examples/
    └── main.go
```

## Test Structure

**Suite Organization:**
```go
package sysml

import (
    "testing"
)

func TestParseString(t *testing.T) {
    // 1. Arrange: set up test input
    input := `
package TestPackage {
    part def Vehicle {
        part engine : Engine;
    }
    part def Engine;
}
`

    // 2. Act: execute function under test
    result := ParseString(input)

    // 3. Assert: verify results
    if !result.Success() {
        t.Fatalf("parse failed: %s", result.Errors)
    }

    if result.Model == nil {
        t.Fatal("model is nil")
    }

    if len(result.Model.Packages) != 1 {
        t.Errorf("expected 1 package, got %d", len(result.Model.Packages))
    }

    // Specific assertions on parsed structure
    pkg := result.Model.Packages[0]
    if pkg.Name() != "TestPackage" {
        t.Errorf("expected package name 'TestPackage', got '%s'", pkg.Name())
    }
}
```

**Patterns:**
- **Setup**: Test input defined inline as string literals (SysML syntax)
- **Execution**: Functions under test called directly on test data
- **Verification**: Multiple assertions check different aspects (success, structure, values)
- **Error cases**: Separate test function for each failure scenario
- **Cleanup**: Not needed (no file I/O or resources to clean up)

## Test Input Strategy

**Inline SysML Strings:**
- All unit tests use inline SysML syntax strings
- Syntax examples match actual SysML v2 grammar:
  ```
  package TestPackage {
      part def Vehicle;
      requirement def SafetyRequirement;
      verification def TestVer;
  }
  ```

**Testdata Directory:**
- External test files stored in `/home/daniel/projects/claudecode/gosysml2/testdata/` (mentioned in integration tests)
- Contains real SysML v2 files for comprehensive grammar validation
- Referenced in `integration_test.go`: `integration_test.go` uses `filepath.WalkDir` to discover and parse all `.sysml` files

**Invalid Input:**
- Intentionally malformed syntax to test error handling:
  ```
  package Broken {
      @@@ invalid syntax
  }
  ```
- Tests verify errors are captured: `if !errors.HasErrors() { t.Error(...) }`

## Mocking

**Framework:** No mocking library detected; not used

**Patterns:**
- No mocks or stubs; tests use real parsers
- Low-level ANTLR errors captured via real `ErrorCollector` implementation
- Integration tests use actual file system (with skip on missing testdata)
- Visitor pattern allows custom implementations for test-specific behavior

**What to Mock:**
- N/A - direct testing without mocks is preferred in this codebase

**What NOT to Mock:**
- Parser and lexer (always test with real ANTLR)
- Error collection (test error paths with real error handling)
- Model building (test tree-to-model conversion with real builder)

## Fixtures and Factories

**Test Data:**
- Inline SysML strings as fixtures (no separate fixture files)
- Common patterns:
  ```go
  input := `
  package TestPackage {
      part def Vehicle;
      part def Engine;
  }
  `
  result := ParseString(input)
  ```

**Location:**
- Test data defined inside each test function (tight coupling by design)
- No shared fixture files or helper functions for input generation
- Test data is minimal and focused on what's being tested

## Coverage

**Requirements:** No coverage enforcement detected (no `.codecov.yml` or build rules)

**View Coverage:**
```bash
go test -cover ./sysml
go test -cover ./low
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # Opens browser with detailed coverage map
```

## Test Types

**Unit Tests:**
- **Scope**: Individual functions (`ParseString`, `FindRequirements`, `NewLexer`, etc.)
- **Approach**: Direct function calls with predetermined input/output
- **Files**: `parse_test.go`, `parser_test.go`, `visitor_test.go`
- **Coverage**: Single responsibility per test (e.g., `TestParseString` tests only basic parsing success)

**Integration Tests:**
- **Scope**: Multi-file parsing, directory traversal, complete workflows
- **Approach**: Real file system access, loading actual SysML files from testdata directory
- **Files**: `integration_test.go` (contains `TestParseTestdataFiles`, `TestParseDirectory`, `TestParseDirectoryParallel`, `TestParseDirectoryStream`)
- **Data**: Uses external testdata files (`.sysml` files in project testdata directory)
- **Skip condition**: Tests skip if testdata directory not found: `if _, err := os.Stat(testdataDir); os.IsNotExist(err) { t.Skip(...) }`

**E2E Tests:**
- **Framework**: Not explicitly present
- **Note**: Integration tests function as end-to-end tests for parsing pipeline

## Common Patterns

**Async Testing:**
- Not applicable (no goroutines in parsing, though parallel parsing tested)
- Parallel parsing tested with real `sync.WaitGroup` and semaphore: `ParseDirectoryParallel` invoked with 4 workers
- Test verifies results are correct across parallel execution

**Error Testing:**
Pattern for testing error conditions:
```go
func TestParseStringWithErrors(t *testing.T) {
    input := `
package Broken {
    part def Vehicle {
        invalid syntax here!!!
    }
}
`
    result := ParseString(input)

    // Assert parsing failed
    if result.Success() {
        t.Error("expected parse to fail")
    }

    // Assert errors are present
    if result.Errors == nil || !result.Errors.HasErrors() {
        t.Error("expected errors to be present")
    }

    // Log errors for inspection
    t.Logf("Got expected errors: %s", result.Errors)
}
```

**Panic Testing:**
Pattern for testing intentional panic (Must* functions):
```go
func TestMustParseStringPanic(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Error("expected panic for invalid input")
        }
    }()

    MustParseString(`package { invalid }}}`)
}
```

**Success Path Testing:**
- Primary test for each function tests happy path
- `TestParseString` verifies successful parsing returns valid model
- `TestFindRequirements` verifies correct element discovery
- `TestCounter` verifies visitor pattern works correctly

**Boundary Testing:**
- `TestEmptyInput` tests parsing of empty string (valid edge case)
- Multiple element types tested: parts, requirements, ports, connections, etc.
- Nested structures tested: packages within packages, parts within parts

**Error Collector Testing:**
- `TestErrorCollector` verifies error collection state management
- Tests `.HasErrors()` transitions from false to true
- Tests `.Clear()` resets collector
- Tests `.Errors()` returns all collected errors

**Streaming Pattern Testing:**
```go
func TestParseDirectoryStream(t *testing.T) {
    testdataDir := "../../testdata"
    if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
        t.Skip("testdata directory not found")
    }

    count := 0
    err := ParseDirectoryStream(testdataDir, func(r *ParseResult) error {
        if r.Success() {
            count++
        }
        return nil
    })

    if err != nil {
        t.Fatalf("ParseDirectoryStream failed: %v", err)
    }

    t.Logf("Streamed %d files successfully", count)
}
```

## Test Utilities

**Visitor Pattern for Testing:**
- `Counter` visitor used to count elements: `counter := NewCounter(); Visit(result.Model, counter)`
- `Walk` function with callback for element traversal testing

**Model Navigation:**
- Tests verify model structure: `result.Model.Packages[0].Name()`
- Tests verify child elements: `parts := FindParts(result.Model)`
- Tests verify type-safe element collections: `reqs := FindRequirements(result.Model)`

## Documentation Comments in Tests

- Test comments describe input/expected output: `// Valid input should pass`
- Assertion messages are descriptive: `"expected 1 package, got %d"`
- Error messages include expected vs actual: `"expected package name 'TestPackage', got '%s'"`

## Coverage Areas

**Well-tested:**
- Core parsing: `ParseString`, `ParseFile`, `ParseBytes`, `ParseReader`
- Error scenarios: invalid syntax, missing files, malformed input
- Model building: package creation, element nesting, reference extraction
- Visitor pattern: Find functions, Walk function, Counter visitor
- Error handling: error collection, error formatting, multiple error types
- Directory operations: sequential, parallel, streaming parsing modes
- Special cases: empty input, nested packages, complex element hierarchies

**Gaps (if any):**
- Reference resolution completeness not explicitly tested
- Qualified name computation tested indirectly through navigation
- Complex cross-references (derived from, satisfied by, verified by) not extensively tested
- Memory efficiency under large inputs (parallel/streaming modes test this implicitly)
