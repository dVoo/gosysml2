# Testing Patterns

**Analysis Date:** 2026-02-05

## Test Framework

**Runner:**
- Standard Go testing (`testing` package)
- No external test framework (no testify, ginkgo, etc.)
- Go version: 1.22

**Run Commands:**

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./gosysml2/sysml/...
go test ./gosysml2/low/...

# Run specific test
go test -run TestParseString ./gosysml2/sysml/...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**No Test Configuration Files:**
- No `jest.config.*`, `vitest.config.*`, or equivalent
- Configuration through standard Go test flags

## Test File Organization

**Location:**
- Tests are co-located with source files (same directory)
- Test files use `*_test.go` naming convention

**Structure:**

```
gosysml2/
├── sysml/
│   ├── parse.go           # Source
│   ├── parse_test.go      # Unit tests
│   ├── visitor.go         # Source
│   ├── visitor_test.go    # Unit tests
│   ├── model.go           # Source (no tests - data structures only)
│   ├── errors.go          # Source (no tests)
│   └── integration_test.go # Integration tests
├── low/
│   ├── parser.go          # Source
│   ├── parser_test.go     # Unit tests
│   ├── lexer.go           # Source
│   └── errors.go          # Source
└── examples/
    └── main.go            # Example code (no tests)
```

**Test Data:**
- SysML test input embedded as raw string literals
- External test data in `../../testdata/` directory
- Test files use `.sysml` extension

## Test Structure

**Basic Test Pattern:**

```go
package sysml

import (
	"testing"
)

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

	if len(result.Model.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(result.Model.Packages))
	}
}
```

**Error Testing Pattern:**

```go
func TestParseStringWithErrors(t *testing.T) {
	input := `
package Broken {
    @@@ invalid syntax here!!!
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

**Table-Driven Test Pattern:**

Tests iterate over SysML files in testdata directory:

```go
func TestParseTestdataFiles(t *testing.T) {
	testdataDir := "../../testdata"

	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}

	var files []string
	err := filepath.WalkDir(testdataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sysml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk testdata: %v", err)
	}

	passed := 0
	failed := 0

	for _, file := range files {
		result := ParseFile(file)
		name := filepath.Base(file)

		if result.Success() {
			passed++
			t.Logf("  PASS: %s", name)
		} else {
			failed++
			t.Logf("  FAIL: %s - %s", name, result.Errors.First().Message)
		}
	}

	t.Logf("\nSummary: %d passed, %d failed out of %d files", passed, failed, len(files))

	if failed > 0 {
		t.Errorf("%d files failed to parse", failed)
	}
}
```

**Subtest Pattern:**

```go
func TestParseIndividualFiles(t *testing.T) {
	testdataDir := "../../testdata"

	files, err := findSysMLFiles(testdataDir)
	if err != nil {
		t.Fatalf("failed to find test files: %v", err)
	}

	for _, file := range files {
		// Create test name from relative path
		relPath, _ := filepath.Rel(testdataDir, file)
		testName := strings.TrimSuffix(relPath, ".sysml")
		testName = strings.ReplaceAll(testName, string(filepath.Separator), "_")

		t.Run(testName, func(t *testing.T) {
			result := parseSysMLFile(file)

			if !result.Success() {
				t.Errorf("parsing failed:\n%s", strings.Join(result.Errors, "\n"))
			} else {
				t.Logf("parsed successfully (%d tokens)", result.TokenCount)
			}
		})
	}
}
```

## Assertion Patterns

**Fatal vs Error:**
- Use `t.Fatalf` when the test cannot continue (e.g., parse failed, model is nil)
- Use `t.Errorf` for individual assertion failures (continues testing other cases)

```go
result := ParseString(input)
if !result.Success() {
	t.Fatalf("Failed to parse: %v", result.Errors)  // Fatal - can't continue
}

reqs := FindRequirements(result.Model)
if len(reqs) != 3 {
	t.Errorf("Expected 3 requirements, got %d", len(reqs))  // Error - can continue
}
```

**Common Assertions:**

```go
// Success check
if !result.Success() {
	t.Fatalf("parse failed: %s", result.Errors)
}

// Nil check
if result.Model == nil {
	t.Fatal("model is nil")
}

// Length check
if len(parts) != 2 {
	t.Errorf("Expected 2 parts, got %d", len(parts))
}

// Equality check
if pkg.Name() != "TestPackage" {
	t.Errorf("expected package name 'TestPackage', got '%s'", pkg.Name())
}

// Boolean check
if !ref.IsResolved() {
	t.Error("Expected ref to be unresolved initially")
}

// Map key existence
if !names["SafetyRequirement"] {
	t.Error("Expected to find SafetyRequirement")
}
```

## Mocking

**No Mocking Framework:**
- Tests use real SysML input strings
- No external mocking library (no gomock, mockery, etc.)
- Integration-style testing with actual parser

**Test Doubles:**
- `BaseVisitor` provides a no-op visitor implementation for testing
- Counter visitor used to test visitor pattern:

```go
func TestCounter(t *testing.T) {
	input := `
package P {
    part def A;
    part def B;
    part a1 : A;
    part b1 : B;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	counter := NewCounter()
	Visit(result.Model, counter)

	t.Logf("Counts: %v", counter.Counts)
	t.Logf("Total: %d", counter.Total())

	if counter.Total() == 0 {
		t.Error("expected some elements to be counted")
	}
}
```

## Test Data

**Inline Test Data:**

```go
input := `
package TestPackage {
    part def Vehicle {
        part engine : Engine;
    }
    part def Engine;
}
`
```

**External Test Files:**
- Location: `../../testdata/` (relative to test file)
- Format: `.sysml` files containing valid and complex SysML v2 syntax
- Organized in subdirectories by category:
  - `Import Tests/`
  - `Requirements Examples/`
  - `Variability Examples/`
  - `simple/`

**Test Data Helper:**

```go
func findSysMLFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sysml") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
```

## Test Types

**Unit Tests:**
- Test individual functions in isolation
- Examples: `TestParseString`, `TestValidate`, `TestWalk`
- Location: `*_test.go` files alongside source

**Integration Tests:**
- Test parsing of real SysML files from testdata
- Examples: `TestParseTestdataFiles`, `TestValidateTestdataFiles`
- Location: `integration_test.go`

**Error Tests:**
- Verify error handling and messages
- Examples: `TestParseStringWithErrors`, `TestMustParseStringPanic`

**API Tests:**
- Test both high-level and low-level APIs
- High-level: `TestParseString`, `TestFindRequirements`
- Low-level: `TestParse`, `TestLexer`, `TestParser`

## Coverage

**No Enforced Coverage Target:**
- No coverage configuration files
- Coverage reports generated manually with `go test -cover`

**View Coverage:**

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# View in browser
go tool cover -html=coverage.out
```

## Common Test Patterns

**Setup Pattern:**

```go
func TestCounterWithVerification(t *testing.T) {
	input := `
		package TestPackage {
			part def Vehicle;
			part def Engine;
			requirement def SafetyReq;
			verification def TestVer;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	// Test-specific code...
}
```

**Collection Testing:**

```go
// Build map from results for easier assertions
names := make(map[string]bool)
for _, req := range reqs {
	names[req.Name()] = true
}

if !names["SafetyRequirement"] {
	t.Error("Expected to find SafetyRequirement")
}
```

**Depth Testing:**

```go
func TestWalkWithDepth(t *testing.T) {
	depths := make(map[string]int)
	Walk(result.Model, func(elem Element, depth int) bool {
		depths[elem.Name()] = depth
		return true
	})

	if depths["OuterPackage"] != 0 {
		t.Errorf("Expected OuterPackage at depth 0, got %d", depths["OuterPackage"])
	}
}
```

**Skip Pattern:**

```go
func TestParseTestdataFiles(t *testing.T) {
	testdataDir := "../../testdata"

	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}
	// ...
}
```

## Debug Output

Use `t.Logf` for test debugging information:

```go
t.Logf("Found %d SysML test files", len(files))
t.Logf("  PASS: %s", name)
t.Logf("  FAIL: %s - %s", name, result.Errors.First().Message)
t.Logf("\nSummary: %d passed, %d failed out of %d files", passed, failed, len(files))
```

## Test Utilities

**Custom Error Listener (in tests):**

```go
// ErrorListener collects parsing errors
type ErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []string
}

func NewErrorListener() *ErrorListener {
	return &ErrorListener{
		Errors: make([]string, 0),
	}
}

func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.Errors = append(l.Errors, fmt.Sprintf("line %d:%d %s", line, column, msg))
}
```

---

*Testing analysis: 2026-02-05*
