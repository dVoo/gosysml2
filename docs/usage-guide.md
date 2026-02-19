# Project Usage Guide

This guide covers how to use the gosysml2 project for parsing SysML v2 models.

## Table of Contents

1. [Introduction](#introduction)
2. [Quick Start](#quick-start)
3. [Library Usage](#library-usage)
4. [Command Line Tools](#command-line-tools)
5. [Working with Validation Data](#working-with-validation-data)
6. [Troubleshooting](#troubleshooting)

## Introduction

### What is this Project?

The gosysml2 project is a Go library for parsing SysML v2 models. It provides:

- Full SysML v2 grammar support via ANTLR4-generated parser
- Two-tier API design (low-level and high-level)
- Comprehensive error handling with source locations
- Memory-efficient parsing options
- Parallel and streaming parsing modes

### Who Should Use It?

This project is designed for:

- **Go developers** building tools that need to parse SysML v2 models
- **Systems engineers** working with SysML v2 specifications
- **Tool integrators** needing to extract data from SysML models
- **Researchers** analyzing SysML model structures

### Prerequisites

- Go 1.23 or later
- Basic understanding of SysML v2 syntax
- Familiarity with Go programming

## Quick Start

### Environment Setup (Nix)

This project uses Nix for reproducible development environments:

```bash
# Enter the development shell
nix develop

# Or use direnv (automatically enters shell when cd'ing into directory)
direnv allow
```

### Building Tools

```bash
# Build all command-line tools
cd gosysml2
go build -o ../bin/verify-completeness ./cmd/verify-completeness
go build -o ../bin/verify-parser ./cmd/verify-parser
go build -o ../bin/check-validation ./cmd/check-validation
```

### Running First Parse

```bash
cd gosysml2

# Run the basic example
go run examples/basic/main.go

# Or run a specific example
go run examples/requirements/main.go
go run examples/parallel/main.go
```

## Library Usage

### Reference to docs/USAGE.md

For detailed library usage documentation, see [docs/USAGE.md](./USAGE.md).

### Key Concepts

#### Two-Tier API Design

The library provides two levels of API:

1. **Low-level API** (`low` package): Direct access to lexer, parser, and parse trees
   - Use for: Performance-critical applications, custom parsing logic
   - See: `low.Parse()`, `low.Validate()`, `low.NewLexer()`

2. **High-level API** (`sysml` package): Idiomatic Go model with visitor pattern
   - Use for: Most applications, easier to work with
   - See: `sysml.ParseString()`, `sysml.ParseFile()`, `sysml.ParseDirectory()`

#### Type-Safe References

The model uses **type-safe element references** instead of strings:

```go
// Generic reference type
type Ref[T Element] struct {
    // Name() returns the reference name
    // Resolved() returns the resolved element (nil if unresolved)
    // IsResolved() returns true if resolved
}

// Example usage
for _, req := range sysml.FindAll[*sysml.Requirement](model) {
    if req.Subject.IsResolved() {
        fmt.Printf("Subject: %s\n", req.Subject.Resolved().Name())
    }
}
```

#### Element Hierarchy

All model elements implement the `Element` interface:

```go
type Element interface {
    Name() string
    QualifiedName() string
    Location() Location
    Parent() Element
    Children() []Element
    Documentation() string
    Kind() ElementKind
}
```

Specific types add their own properties:
- `Part` - `IsDefinition`, `TypeRef`
- `Requirement` - `RequirementID`, `DerivedFrom`, `VerifiedBy`
- `Verification` - `Method`, `Subject`

## Command Line Tools

### verify-completeness

Validates that all elements in a model are being handled (not discarded).

```bash
# Basic usage
./bin/verify-completeness ./validationdata/parts-tree.sysml

# Check multiple files
./bin/verify-completeness ./validationdata/*.sysml

# With verbose output
./bin/verify-completeness -v ./validationdata/parts-tree.sysml
```

**What it checks:**
- All parsed elements are present in the model
- No elements are silently discarded
- Element counts match expected values

**Output:**
- List of unhandled elements (if any)
- Element counts by type
- Success/failure status

### verify-parser

Validates parser output against expected results.

```bash
# Basic usage
./bin/verify-parser ./validationdata/parts-tree.sysml

# Compare against expected output
./bin/verify-parser -expected ./expected/parts-tree.json ./validationdata/parts-tree.sysml

# Generate expected output
./bin/verify-parser -generate ./validationdata/parts-tree.sysml
```

**What it checks:**
- Parse errors (should be none for valid files)
- Element structure matches expectations
- Qualified names are correct

### check-validation

Runs the validation test suite against all validation data files.

```bash
# Run full validation suite
cd gosysml2
go run ./cmd/check-validation

# With verbose output
go run ./cmd/check-validation -v

# Check specific category
go run ./cmd/check-validation -category parts-tree
```

**What it checks:**
- All files in `./validationdata/` parse successfully
- Element counts are reasonable
- No critical parse errors

**Categories validated:**
- parts-tree
- function-based-behavior
- state-based-behavior
- requirements
- verification
- (and 14 more categories)

## Working with Validation Data

### Understanding validationdata/ Structure

The `./validationdata/` directory contains 18 categories of SysML validation files:

```
validationdata/
├── parts-tree/
│   ├── basic-parts.sysml
│   ├── nested-parts.sysml
│   └── ...
├── requirements/
│   ├── simple-requirements.sysml
│   ├── derived-requirements.sysml
│   └── ...
├── verification/
│   ├── test-cases.sysml
│   └── ...
└── ... (15 more categories)
```

### Running Validation Tests

```bash
cd gosysml2

# Run all validation tests
go test ./... -run TestValidation

# Run specific category
go test ./... -run TestValidation/PartsTree

# With verbose output
go test ./... -v -run TestValidation
```

### Interpreting Results

**Success:**
```
PASS: parts-tree/basic-parts.sysml
  Parsed 15 elements successfully
```

**Failure:**
```
FAIL: requirements/complex-reqs.sysml
  Parse error at line 42: unexpected token
  Elements parsed: 23/25 expected
```

**Partial Success:**
```
WARN: verification/nested-verification.sysml
  Parsed with 2 warnings
  Model built successfully
```

### Adding New Validation Files

1. Create `.sysml` file in appropriate category directory
2. Ensure it represents valid SysML v2 syntax
3. Run validation to generate expected output
4. Add to test suite if it covers new ground

## Troubleshooting

### Common Issues

#### Issue: "cannot find package" error

**Symptom:**
```
go: module github.com/dVoo/gosysml2: not found
```

**Solution:**
```bash
# Update go.mod references
cd gosysml2
go mod tidy

# Or use local replace
go mod edit -replace github.com/dVoo/gosysml2=./
```

#### Issue: Parse errors on valid-looking SysML

**Symptom:**
```
Parse error at line 10: unexpected token 'part'
```

**Possible causes:**
1. Using SysML v1 syntax (this is a v2 parser)
2. Missing semicolons after declarations
3. Using unsupported grammar elements

**Solution:**
- Check [PERFORMANCE.md](./PERFORMANCE.md) for grammar coverage
- Verify syntax against SysML v2 specification
- Check validation data for examples of supported syntax

#### Issue: Out of memory when parsing large files

**Symptom:**
```
fatal error: runtime: out of memory
```

**Solution:**
```go
// Use streaming mode
err := sysml.ParseDirectoryStream(dir, handler, sysml.WithDiscardTree())

// Or use parallel parsing with memory limits
results, err := sysml.ParseDirectoryParallel(dir, 2) // limit workers
```

#### Issue: Slow parsing performance

**Symptom:**
Parsing takes longer than expected

**Solutions:**
```go
// Use parallel parsing for multiple files
results, _ := sysml.ParseDirectoryParallel(dir, 0)

// Discard parse tree to reduce memory pressure
result := sysml.ParseFile(path, sysml.WithDiscardTree())

// Validate only (no model building) for quick checks
errors := low.Validate(input)
```

### Getting Help

1. **Check the examples:** All examples in `gosysml2/examples/` are working code
2. **Review USAGE.md:** Detailed usage patterns and API reference
3. **Check validation data:** See how different SysML constructs are represented
4. **Run tests:** `go test ./...` will show what's working

### Debug Mode

Enable debug logging:

```go
// Set debug environment variable
os.Setenv("GOSYSML_DEBUG", "1")

// Or use parser option (if available)
result := sysml.ParseString(input, sysml.WithDebug())
```

### Performance Profiling

```bash
# CPU profiling
go run -cpuprofile=cpu.prof examples/parallel/main.go
go tool pprof cpu.prof

# Memory profiling
go run -memprofile=mem.prof examples/parallel/main.go
go tool pprof mem.prof
```

## Additional Resources

- [Library Usage Guide](./USAGE.md) - Detailed API documentation
- [Performance Guide](./PERFORMANCE.md) - Optimization strategies
- [Examples](../examples/) - Working code samples
- [README](../README.md) - Project overview and quick reference
