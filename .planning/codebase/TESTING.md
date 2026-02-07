# Testing

## Overview

Comprehensive test coverage using Go's standard testing framework with table-driven tests and benchmarks.

## Test Framework

- **Standard Go testing** — `testing` package
- **No external test frameworks** — Pure Go standard library

## Test Structure

### File Naming
- `*_test.go` — Unit tests
- `*_bench_test.go` — Benchmarks
- Tests alongside source files

### Test Organization

#### Unit Tests
```go
func TestPartCreation(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
    }{
        {"simple part", "partdef Engine;", 1},
        {"nested parts", "partdef Car { part engine; }", 2},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := sysml.ParseString(tt.input)
            if !result.Success() {
                t.Fatalf("parse failed: %v", result.Errors)
            }
            // assertions
        })
    }
}
```

#### Integration Tests
- `sysml/integration_test.go` — End-to-end parsing tests
- Tests against real validation files
- Verifies element counts and relationships

#### Benchmarks
```go
func BenchmarkParseSingle(b *testing.B) {
    input := loadTestFile("large.sysml")
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sysml.ParseString(input)
    }
}
```

## Test Coverage Areas

### Model Types
- Constructor tests for each element type
- Type assertions and interface compliance
- Parent/child relationship validation

### Parser
- Valid input parsing
- Invalid input error handling
- Edge cases (empty input, comments, etc.)

### Reference Resolution
- Resolution success cases
- Unresolved reference handling
- Circular reference detection

### Visitor Pattern
- Visitor traversal order
- Early termination (return false)
- Depth tracking

## Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./sysml/...

# With coverage
go test -cover ./...

# Benchmarks
go test -bench=. ./...

# Verbose
go test -v ./sysml/...
```

## Validation Testing

### Validation Suite
- 18 categories of test files
- 56 validation files total
- 96.4% pass rate (54/56)

### Running Validation
```bash
go run ./cmd/check_validation
```

## Test Data

- `validationdata/` — Official SysML v2 validation files
- `testdata/` — Project-specific test files
- Test files loaded from disk (not embedded)

## Mocking

No mocking framework used. Tests use:
- Real parse results
- String inputs for unit tests
- File-based integration tests
