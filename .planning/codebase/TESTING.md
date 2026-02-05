# Testing Patterns

**Analysis Date:** 2026-02-05

## Test Framework

**Runner:** Go's built-in `testing` package (standard library)

**No external test framework** - uses standard Go testing idioms

**Run Commands:**
```bash
go test ./...                    # Run all tests
go test ./sysml                  # Run sysml package tests
go test ./low                    # Run low package tests
go test -v ./...                 # Verbose output
go test -race ./...              # Race detection
go test -cover ./...             # Coverage
```

## Test File Organization

**Location:** Co-located with source files (Go standard convention)

**Naming:** `*_test.go` suffix

**Structure:**
```
gosysml2/
├── sysml/
│   ├── parse.go
│   ├── parse_test.go           # Unit tests for parsing
│   ├── visitor_test.go         # Tests for visitor pattern
│   ├── integration_test.go     # Integration tests
│   ├── model.go
│   └── errors.go
├── low/
│   ├── parser.go
│   ├── parser_test.go          # Low-level parser tests
│   ├── lexer.go
│   └── errors.go
└── examples/
    └── main.go                 # Example usage (not tested)
```

## Test Structure

**Function Naming:**
```go
func TestParseString(t *testing.T)           // Basic functionality
func TestParseStringWithErrors(t *testing.T) // Error cases
func TestValidate(t *testing.T)              // Specific feature
func TestCounter(t *testing.T)               // Complex scenario
```

**Pattern:** Table-driven tests with inline input data
```go
func TestParseString(t *testing.T) {
    input := `
package TestPackage {
    part def Vehicle {
        part engine : Engine;
    }
    part def Engine;
}
`
    result := ParseString(input)

    if !result.Success() {
        t.Fatalf("parse failed: %s", result.Errors)
    }

    if result.Model == nil {
        t.Fatal("model is nil")
    }
}
```

**Error Testing Pattern:**
```go
func TestParseStringWithErrors(t *testing.T) {
    input := `
package Broken {
    @@@ invalid syntax
}
`
    result := ParseString(input)

    if result.Success() {
        t.Error("expected parse to fail")
    }

    if result.Errors == nil || !result.Errors.HasErrors() {
        t.Error("expected errors to be present")
    }

    t.Logf("Got expected errors: %s", result.Errors)
}
```

**Panic Testing Pattern:**
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

## Mocking

**Framework:** None - tests use real implementations

**Mocking Strategy:**
- No external dependencies to mock (self-contained parser)
- Tests use real SysML input strings
- File system tests use actual testdata files

**Test Data:**
- Inline string literals for unit tests
- `testdata/` directory for integration tests (45+ .sysml files)

## Fixtures and Factories

**Test Data Pattern:**
```go
// Inline fixture
func TestFindRequirements(t *testing.T) {
    input := `
        package TestPackage {
            requirement def SafetyRequirement {
                doc /* Safety requirements must be satisfied */
            }
            requirement def PerformanceRequirement;
            requirement testReq : SafetyRequirement;
        }
    `
    result := ParseString(input)
    // assertions...
}
```

**External Fixtures:**
- Location: `../../testdata/` (relative to test files)
- Format: `.sysml` files
- Used in: `integration_test.go`

**File Walking Pattern:**
```go
func TestParseTestdataFiles(t *testing.T) {
    testdataDir := "../../testdata"

    if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
        t.Skip("testdata directory not found")
    }

    var files []string
    err := filepath.WalkDir(testdataDir, func(path string, d os.DirEntry, err error) error {
        if !d.IsDir() && strings.HasSuffix(path, ".sysml") {
            files = append(files, path)
        }
        return nil
    })
    // assertions...
}
```

## Coverage

**Requirements:** None enforced

**View Coverage:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Coverage Analysis:**
- Unit tests in `parse_test.go`, `visitor_test.go`
- Integration tests in `integration_test.go`
- Low-level tests in `parser_test.go`
- Test data from 45+ SysML example files

## Test Types

**Unit Tests:**
- File: `gosysml2/sysml/parse_test.go` (176 lines)
- Focus: Individual parsing functions
- Pattern: Single input, specific assertions

**Visitor Tests:**
- File: `gosysml2/sysml/visitor_test.go` (434 lines)
- Focus: Visitor pattern, element finding, counting
- Pattern: Multiple element types tested

**Integration Tests:**
- File: `gosysml2/sysml/integration_test.go` (197 lines)
- Focus: Directory parsing, file walking, parallel processing
- Uses: External testdata files

**Low-Level Tests:**
- File: `gosysml2/low/parser_test.go` (213 lines)
- Focus: Lexer, parser, error collection
- Pattern: Direct API testing

## Common Patterns

**Success/Failure Pattern:**
```go
result := ParseString(input)
if !result.Success() {
    t.Fatalf("Failed to parse: %v", result.Errors)
}
```

**Collection Testing:**
```go
names := make(map[string]bool)
for _, req := range reqs {
    names[req.Name()] = true
}

if !names["SafetyRequirement"] {
    t.Error("Expected to find SafetyRequirement")
}
```

**Count Assertions:**
```go
if counter.Counts[KindPackage] != 1 {
    t.Errorf("Expected 1 package, got %d", counter.Counts[KindPackage])
}
```

**Directory Skipping:**
```go
if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
    t.Skip("testdata directory not found")
}
```

**Deferred Assertions:**
```go
defer func() {
    if r := recover(); r == nil {
        t.Error("expected panic")
    }
}()
```

## Test Categories

**Parser Tests:**
- `TestParseString` - Basic parsing
- `TestParseStringWithErrors` - Error handling
- `TestValidate` - Validation-only mode
- `TestWalk` - Tree traversal
- `TestEmptyInput` - Edge cases

**Visitor Tests:**
- `TestFindRequirements` - Element finding by type
- `TestFindVerifications` - Specific element types
- `TestCounter` - Counting visitor
- `TestWalkWithDepth` - Depth tracking
- `TestFilter` - Predicate filtering

**Integration Tests:**
- `TestParseTestdataFiles` - File-based parsing
- `TestParseDirectory` - Directory traversal
- `TestParseDirectoryParallel` - Concurrent parsing
- `TestParseDirectoryStream` - Streaming mode
- `TestWithDiscardTree` - Memory optimization option

**Low-Level Tests:**
- `TestParse` - Basic low-level parsing
- `TestLexer` - Token generation
- `TestErrorCollector` - Error accumulation
- `TestParseErrors` - Error aggregation

## Testing Best Practices (Observed)

1. **Descriptive test names** - `TestFindRequirementsWithVerification`
2. **Inline test data** - SysML strings in tests
3. **t.Fatalf for setup failures** - Stop early if precondition fails
4. **t.Errorf for assertion failures** - Continue to report multiple failures
5. **t.Logf for debugging** - Verbose logging with results
6. **t.Skip for missing dependencies** - Graceful handling of missing testdata
7. **defer for cleanup** - Panic recovery testing

---

*Testing analysis: 2026-02-05*
