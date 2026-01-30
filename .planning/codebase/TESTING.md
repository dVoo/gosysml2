# Testing Patterns

**Analysis Date:** 2026-01-30

## Test Framework

**Runner:**
- Go standard `testing` package (no external test framework)
- Run tests with `go test ./...` across the module

**Assertion Library:**
- Manual assertions using `testing.T` methods: `t.Fatalf()`, `t.Errorf()`, `t.Fatal()`, `t.Error()`
- No external assertion library (testify, assert, etc.)

**Run Commands:**
```bash
# Run all tests in a package
go test ./gosysml2/sysml

# Run tests in low-level parser
go test ./gosysml2/low

# Run specific test by name
go test ./gosysml2/sysml -run TestParseString

# Watch mode (requires external tool, not configured)
# Manual: watch -n 2 'go test ./...'

# Coverage
go test ./gosysml2/sysml -cover
go test ./gosysml2/... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

## Test File Organization

**Location:**
- Co-located with source files in same package directory
- `parse_test.go` alongside `parse.go` in `gosysml2/sysml/`
- `visitor_test.go` alongside `visitor.go` in `gosysml2/sysml/`
- `parser_test.go` alongside `parser.go` in `gosysml2/low/`
- `integration_test.go` in `gosysml2/sysml/` for cross-module scenarios

**Naming:**
- `*_test.go` suffix following Go standard
- Test functions prefixed with `Test` followed by capitalized subject: `TestParseString`, `TestParseStringWithErrors`, `TestWalk`, `TestFindByKind`
- Helper functions may be prefixed with lowercase or Test prefix depending on scope

**Structure:**
```
gosysml2/
├── sysml/
│   ├── parse.go
│   ├── parse_test.go           # Tests for ParseString, ParseFile, ParseBytes, etc.
│   ├── visitor.go
│   ├── visitor_test.go         # Tests for FindRequirements, FindVerifications, etc.
│   ├── integration_test.go     # Tests parsing testdata files
│   ├── model.go
│   └── errors.go
└── low/
    ├── parser.go
    ├── parser_test.go          # Tests for low-level parser
    ├── lexer.go
    └── errors.go
```

## Test Structure

**Suite Organization:**
Tests use Go's standard table-driven test pattern where applicable, individual test functions for integration/behavior tests.

**Basic pattern from `parse_test.go`:**
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

	if len(result.Model.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(result.Model.Packages))
	}

	pkg := result.Model.Packages[0]
	if pkg.Name() != "TestPackage" {
		t.Errorf("expected package name 'TestPackage', got '%s'", pkg.Name())
	}
}
```

**Patterns:**
- **Setup:** Input data defined inline as string literals (SysML syntax embedded in Go strings)
- **Teardown:** Implicit through Go's garbage collection; no explicit cleanup
- **Assertions:** Direct comparisons with `if` statements followed by `t.Errorf()` or `t.Fatalf()`
- **Success cases:** Call function, check result, verify expected values
- **Error cases:** Check `result.Success()` returns `false`, verify `result.Errors` is populated
- **Assertions for nil:** Explicit `if X == nil` or `if X != nil` checks

## Mocking

**Framework:** No external mocking library used; simple stubbing with test fixtures

**Patterns:**
Embed test input data directly in test functions as SysML string literals:

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

	if result.Success() {
		t.Error("expected parse to fail")
	}

	if result.Errors == nil || !result.Errors.HasErrors() {
		t.Error("expected errors to be present")
	}

	t.Logf("Got expected errors: %s", result.Errors)
}
```

**What to Mock:**
- File system operations tested with actual `filepath.WalkDir` and real test files (see integration tests)
- ANTLR parser behavior tested directly through public API

**What NOT to Mock:**
- Parser internals (tested through public API)
- Tree walking and model building (integration tested)
- Error collection (tested with known error inputs)

## Fixtures and Fixtures and Factories

**Test Data:**
Test data embedded directly in test functions as inline string literals:

```go
func TestWalk(t *testing.T) {
	input := `
package P1 {
    package P2 {
        part def A;
    }
    part def B;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}
	// ... test continues
}
```

**Testdata directory:**
- `testdata/simple/` contains 34+ `.sysml` files covering language features:
  - `ActionTest.sysml`, `AliasTest.sysml`, `AllocationTest.sysml`, `AnalysisTest.sysml`
  - `AssignmentTest.sysml`, `CalculationTest.sysml`, `CommentTest.sysml`, `ConjugationTest.sysml`
  - `ConnectionTest.sysml`, `ConstraintTest.sysml`, etc.
- `testdata/Vehicle Example/` contains realistic example models:
  - `SysML v2 Spec Annex A SimpleVehicleModel.sysml`
  - `VehicleDefinitions.sysml`, `VehicleUsages.sysml`, `VehicleIndividuals.sysml`

**Location:**
- Testdata files at project root: `testdata/simple/*.sysml` and `testdata/Vehicle Example/*.sysml`
- Accessed from tests using relative paths: `../../testdata`
- Integration tests use `filepath.WalkDir()` to load and parse all `.sysml` files

## Coverage

**Requirements:** No explicit code coverage requirements configured

**View Coverage:**
```bash
# Generate coverage report
go test ./gosysml2/sysml -coverprofile=coverage.out

# View in HTML
go tool cover -html=coverage.out

# Show coverage for specific function
go tool cover -html=coverage.out -o coverage.html
```

## Test Types

**Unit Tests:**
- **Scope:** Individual functions and methods
- **Location:** `parse_test.go`, `visitor_test.go`, `parser_test.go`
- **Approach:** Parse simple inline SysML strings, verify result structure, check error handling
- **Examples:**
  - `TestParseString` - basic parsing
  - `TestParseStringWithErrors` - error collection
  - `TestValidate` - validation without model building
  - `TestFindByKind` - element filtering
  - `TestCounter` - visitor pattern usage
  - `TestEmptyInput` - edge case handling

**Integration Tests:**
- **Scope:** Full parsing pipeline with real testdata files
- **Location:** `integration_test.go` in `gosysml2/sysml/`
- **Approach:**
  - Walk testdata directory
  - Parse each `.sysml` file using `ParseFile()`
  - Collect pass/fail statistics
  - Report failures with line/column details
- **Examples:**
  - `TestParseTestdataFiles` - parses 30+ test files, reports pass/fail counts
  - `TestValidateTestdataFiles` - validates files without building full model
- **Coverage:** Tests against real SysML examples from spec and complex vehicle model

**E2E Tests:**
- **Framework:** Not explicitly present; integration tests serve this purpose
- **Behavior:** Tests parse complete realistic models (`Vehicle Example` directory) and verify parsing succeeds

## Common Patterns

**Async Testing:**
Testing does not use goroutines or async patterns. All tests are synchronous.

**Error Testing:**
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

	// Check Success() method
	if result.Success() {
		t.Error("expected parse to fail")
	}

	// Check Errors is populated
	if result.Errors == nil || !result.Errors.HasErrors() {
		t.Error("expected errors to be present")
	}

	// Verify error details
	t.Logf("Got expected errors: %s", result.Errors)
}
```

**Panic Testing:**
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

**Skipable Tests:**
```go
func TestParseTestdataFiles(t *testing.T) {
	testdataDir := "../../testdata"

	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}
	// ... test continues
}
```

**Visitor Pattern Testing:**
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
	if counter.Total() == 0 {
		t.Error("expected some elements to be counted")
	}
}
```

## Test Data Organization

**Embedded test data:**
- SysML syntax embedded directly in test functions as multi-line strings
- Simple structures: 3-20 lines per test case
- Complex structures: 50+ lines with nested packages and elements

**External test data:**
- 34+ SysML files in `testdata/simple/` directory (one per language feature)
- Realistic examples in `testdata/Vehicle Example/` (4 files: definitions, usages, individuals, spec example)
- Files used by `integration_test.go` to validate parser against real specifications

**Test file locations from codebase:**
- `testdata/simple/ActionTest.sysml`
- `testdata/simple/AliasTest.sysml`
- `testdata/simple/AllocationTest.sysml`
- `testdata/simple/AnalysisTest.sysml`
- `testdata/simple/AssignmentTest.sysml`
- `testdata/simple/CalculationTest.sysml`
- `testdata/simple/CommentTest.sysml`
- `testdata/simple/ConjugationTest.sysml`
- `testdata/simple/ConnectionTest.sysml`
- `testdata/simple/ConstraintTest.sysml`
- `testdata/Vehicle Example/SysML v2 Spec Annex A SimpleVehicleModel.sysml`
- `testdata/Vehicle Example/VehicleDefinitions.sysml`
- `testdata/Vehicle Example/VehicleUsages.sysml`
- `testdata/Vehicle Example/VehicleIndividuals.sysml`

---

*Testing analysis: 2026-01-30*
