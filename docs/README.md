# gosysml2

A Go library for parsing SysML v2 models. Provides both a low-level performance-oriented API and a high-level developer-friendly API.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [Performance](#performance)
- [Element Kinds](#element-kinds)
- [Error Handling](#error-handling)
- [Examples](#examples)
- [Documentation](#documentation)
- [License](#license)

## Features

- Full SysML v2 grammar support via ANTLR4-generated parser
- Two-tier API design:
  - **Low-level API** (`low` package): Direct access to lexer, parser, and parse trees
  - **High-level API** (`sysml` package): Idiomatic Go model with visitor pattern
- Comprehensive error handling with source locations
- Memory-efficient parsing options for large repositories
- Parallel and streaming parsing modes
- Element finder functions for requirements, verifications, parts, and more

## Installation

```bash
go get github.com/dVoo/gosysml2
```

## Quick Start

### Parse a SysML String

```go
package main

import (
    "fmt"
    "github.com/dVoo/gosysml2/sysml"
)

func main() {
    input := `
        package Vehicle {
            part def Engine {
                attribute power : Real;
            }
            part def Car {
                part engine : Engine;
            }
        }
    `

    result := sysml.ParseString(input)
    if !result.Success() {
        fmt.Printf("Parse error: %s\n", result.Errors)
        return
    }

    // Access the model
    for _, pkg := range result.Model.Packages {
        fmt.Printf("Package: %s\n", pkg.Name())
    }
}
```

### Parse a File

```go
result := sysml.ParseFile("model.sysml")
if !result.Success() {
    fmt.Printf("Error: %s\n", result.Errors)
    return
}

// Find all parts in the model
parts := sysml.FindParts(result.Model)
for _, part := range parts {
    fmt.Printf("Part: %s (definition: %v)\n", part.Name(), part.IsDefinition)
}
```

### Parse a Directory

```go
// Sequential parsing
results, err := sysml.ParseDirectory("./models")

// Parallel parsing (faster for large repositories)
results, err := sysml.ParseDirectoryParallel("./models", 4) // 4 workers

// Streaming (lowest memory usage)
err := sysml.ParseDirectoryStream("./models", func(r *sysml.ParseResult) error {
    if r.Success() {
        // Process each model as it's parsed
        reqs := sysml.FindRequirements(r.Model)
        fmt.Printf("%s: %d requirements\n", r.Source, len(reqs))
    }
    return nil
})
```

## API Reference

### High-Level API (`sysml` package)

#### Parsing Functions

| Function | Description |
|----------|-------------|
| `ParseString(input, opts...)` | Parse a SysML string |
| `ParseFile(path, opts...)` | Parse a SysML file |
| `ParseBytes(data, source, opts...)` | Parse from byte slice (avoids copy) |
| `ParseReader(r, source, opts...)` | Parse from io.Reader |
| `ParseDirectory(dir, opts...)` | Parse all .sysml files in directory |
| `ParseDirectoryParallel(dir, workers, opts...)` | Parse in parallel |
| `ParseDirectoryStream(dir, handler, opts...)` | Parse with streaming handler |
| `Validate(input)` | Validate syntax without building model |
| `ValidateFile(path)` | Validate a file's syntax |

#### Parse Options

```go
// Discard parse tree after building model (reduces memory ~30%)
sysml.ParseString(input, sysml.WithDiscardTree())
```

#### Model Types

The model uses **type-safe element references** instead of strings. References are automatically resolved after parsing.

```go
// Generic reference type - can be resolved or unresolved
type Ref[T Element] struct {
    // Name() returns the reference name
    // Resolved() returns the resolved element (nil if unresolved)
    // IsResolved() returns true if resolved
}

type Model struct {
    Packages []*Package
    Imports  []*Import
    Elements []Element

    // FindByQualifiedName(qn string) Element - fast lookup by qualified name
    // BuildIndex() - builds element index (called automatically)
    // ResolveReferences() - resolves all references (called automatically)
}

type Package struct {
    // Element interface: Name(), QualifiedName(), Location(), Parent(), Children(), Documentation()
    IsLibrary bool

    // Typed accessors (type-safe, no casting needed)
    // Packages() []*Package
    // Parts() []*Part
    // Requirements() []*Requirement
    // Actions() []*Action
}

type Part struct {
    IsDefinition bool
    TypeRef      Ref[*Part]  // Type-safe reference to part definition
    Multiplicity string      // Usage multiplicity, e.g. "4", "0..1", "*"

    // Typed accessors
    // Attributes() []*Attribute
    // Parts() []*Part
    // Ports() []*Port
}

type Attribute struct {
    TypeRef      Ref[Element]  // Reference to attribute type
    DefaultValue string
    IsReadOnly   bool
    IsDerived    bool
}

type Requirement struct {
    IsDefinition  bool
    TypeRef       Ref[*Requirement]  // Reference to requirement definition
    RequirementID string
    Subject       Ref[Element]       // Reference to subject element

    // Relationships with REAL references (not strings!)
    DerivedFrom  []*Requirement    // Requirements this is derived from
    DerivedReqs  []*Requirement    // Requirements derived from this (inverse)
    SatisfiedBy  []Element         // Elements that satisfy this
    VerifiedBy   []*Verification   // Verification cases that verify this

    // Constraints
    Assumptions []*RequirementConstraint  // assume constraints
    Constraints []*RequirementConstraint  // require constraints

    // Text() string - returns documentation text
    // Requirements() []*Requirement - nested requirements
}

type Verification struct {
    IsDefinition        bool
    TypeRef             Ref[*Verification]
    Subject             Ref[Element]
    VerifiedRequirement *Requirement        // Direct reference to verified requirement
    Method              VerificationMethod  // test, analysis, inspection, demonstration

    // Actions() []*Action - nested actions
}

// Also available: Action, Concern, UseCase, AnalysisCase, Import, Comment, Port, Attribute
```

#### Working with References

```go
result := sysml.ParseString(input)
model := result.Model

// Find a requirement
reqs := sysml.FindRequirements(model)
for _, req := range reqs {
    // Access derived requirements (real pointers, not strings!)
    for _, derived := range req.DerivedFrom {
        fmt.Printf("%s is derived from %s\n", req.Name(), derived.Name())
    }

    // Check verification cases
    for _, ver := range req.VerifiedBy {
        fmt.Printf("%s is verified by %s\n", req.Name(), ver.Name())
    }

    // Access subject if resolved
    if req.Subject.IsResolved() {
        fmt.Printf("Subject: %s\n", req.Subject.Resolved().Name())
    }
}

// For usages, access the type definition
parts := sysml.FindParts(model)
for _, part := range parts {
    if !part.IsDefinition && part.TypeRef.IsResolved() {
        def := part.TypeRef.Resolved()
        fmt.Printf("Part %s is of type %s\n", part.Name(), def.Name())
    }
}
```

#### Qualified Names and Lookups

```go
// Elements have qualified names
elem := model.FindByQualifiedName("Vehicle::Engine::power")
fmt.Println(elem.QualifiedName())  // "Vehicle::Engine::power"

// Or use the finder function
elem = sysml.FindByQualifiedName(model, "Vehicle::Engine")
```

#### Typed Child Accessors

```go
// No more casting! Use typed accessors
pkg := model.FindPackage("Vehicle")
for _, part := range pkg.Parts() {        // Returns []*Part
    for _, attr := range part.Attributes() {  // Returns []*Attribute
        fmt.Printf("%s.%s\n", part.Name(), attr.Name())
    }
}

// Also: Action, Concern, UseCase, AnalysisCase, Import, Comment
```

#### Finder Functions

```go
// Find elements by type (returns typed slices)
packages := sysml.FindPackages(model)         // []*Package
parts := sysml.FindParts(model)               // []*Part
requirements := sysml.FindRequirements(model) // []*Requirement
verifications := sysml.FindVerifications(model)
concerns := sysml.FindConcerns(model)
useCases := sysml.FindUseCases(model)
analysisCases := sysml.FindAnalysisCases(model)
actions := sysml.FindActions(model)
attributes := sysml.FindAttributes(model)
ports := sysml.FindPorts(model)

// Find definitions vs usages
definitions := sysml.FindDefinitions(model)  // []Definition interface
usages := sysml.FindUsages(model)            // []Usage interface

// Find by kind or name
elements := sysml.FindByKind(model, sysml.KindPart)
elements := sysml.FindByName(model, "Engine")

// Find by qualified name (fast indexed lookup)
elem := sysml.FindByQualifiedName(model, "Vehicle::Engine")

// Custom filter
elements := sysml.Filter(model, func(e sysml.Element) bool {
    return e.Kind() == sysml.KindRequirement
})
```

#### Visitor Pattern

```go
// Implement custom visitor
type MyVisitor struct {
    sysml.BaseVisitor
    partCount int
}

func (v *MyVisitor) VisitPart(part *sysml.Part) bool {
    v.partCount++
    return true // continue visiting children
}

visitor := &MyVisitor{}
sysml.Visit(model, visitor)
fmt.Printf("Found %d parts\n", visitor.partCount)

// Or use the built-in counter
counter := sysml.NewCounter()
sysml.Visit(model, counter)
fmt.Printf("Counts: %v\n", counter.Counts)
fmt.Printf("Total: %d\n", counter.Total())
```

#### Walk Function

```go
// Walk with depth tracking
sysml.Walk(model, func(elem sysml.Element, depth int) bool {
    indent := strings.Repeat("  ", depth)
    fmt.Printf("%s%s: %s\n", indent, elem.Kind(), elem.Name())
    return true // continue walking
})
```

### Low-Level API (`low` package)

For performance-critical applications or when you need direct access to the parse tree:

```go
import "github.com/dVoo/gosysml2/low"

// Parse and get raw parse tree
tree, errors := low.Parse(input)
if errors.HasErrors() {
    for _, err := range errors.All() {
        fmt.Printf("Line %d: %s\n", err.Line, err.Message)
    }
}

// Access tokens
lexer := low.NewLexer(input)
for token := lexer.NextToken(); token != nil; token = lexer.NextToken() {
    fmt.Printf("%s: %s\n", low.TokenName(token.GetTokenType()), token.GetText())
}

// Validation only (no parse tree construction)
errors := low.Validate(input)
```

## Performance

### Memory Optimization

For large repositories, use these options to reduce memory usage:

```go
// 1. Discard parse tree (saves ~30% memory)
result := sysml.ParseFile(path, sysml.WithDiscardTree())

// 2. Use streaming for very large repositories
sysml.ParseDirectoryStream(dir, handler, sysml.WithDiscardTree())

// 3. Use ParseBytes when you already have []byte
data, _ := os.ReadFile(path)
result := sysml.ParseBytes(data, path, sysml.WithDiscardTree())
```

### Parallel Parsing

For multi-core machines, parallel parsing significantly improves throughput:

```go
// Use all available CPU cores
results, _ := sysml.ParseDirectoryParallel(dir, 0)

// Or specify worker count
results, _ := sysml.ParseDirectoryParallel(dir, 4)
```

### Capacity Guidelines

| Repository Size | Recommended Approach | 32GB RAM |
|-----------------|---------------------|----------|
| < 100 MB | `ParseDirectory()` | Easy |
| 100-500 MB | `ParseDirectoryParallel()` + `WithDiscardTree()` | OK |
| 500 MB - 1 GB | `ParseDirectoryStream()` | Possible |
| > 1 GB | See PERFORMANCE.md | Challenging |

## Element Kinds

```go
const (
    KindPackage
    KindPart
    KindRequirement
    KindVerification
    KindConcern
    KindUseCase
    KindAnalysis
    KindAction
    KindState
    KindConstraint
    KindConnection
    KindInterface
    KindAllocation
    KindImport
    KindComment
    // ... and more
)
```

## Error Handling

```go
result := sysml.ParseFile("model.sysml")

if !result.Success() {
    // Get first error
    first := result.Errors.First()
    fmt.Printf("Error at line %d, column %d: %s\n",
        first.Line, first.Column, first.Message)

    // Iterate all errors
    for _, err := range result.Errors.Errors {
        fmt.Printf("- %s\n", err)
    }

    // Get formatted error string
    fmt.Println(result.Errors.Error())
}
```

## Examples

### Extract All Requirements with Traceability

```go
result := sysml.ParseFile("requirements.sysml")
if !result.Success() {
    log.Fatal(result.Errors)
}

for _, req := range sysml.FindRequirements(result.Model) {
    fmt.Printf("Requirement: %s\n", req.Name())
    if req.RequirementID != "" {
        fmt.Printf("  ID: %s\n", req.RequirementID)
    }
    if len(req.DerivedFrom) > 0 {
        fmt.Printf("  Derived from: %v\n", req.DerivedFrom)
    }
    if len(req.Verifies) > 0 {
        fmt.Printf("  Verified by: %v\n", req.Verifies)
    }
}
```

### Count Elements by Type

```go
results, _ := sysml.ParseDirectory("./models")

totals := make(map[sysml.ElementKind]int)
for _, r := range results {
    if r.Success() {
        counter := sysml.NewCounter()
        sysml.Visit(r.Model, counter)
        for kind, count := range counter.Counts {
            totals[kind] += count
        }
    }
}

fmt.Println("Element counts across all models:")
for kind, count := range totals {
    fmt.Printf("  %s: %d\n", kind, count)
}
```

### Custom Model Traversal

```go
// Find all parts that are definitions (not usages)
definitions := sysml.Filter(model, func(e sysml.Element) bool {
    if part, ok := e.(*sysml.Part); ok {
        return part.IsDefinition
    }
    return false
})

fmt.Printf("Found %d part definitions\n", len(definitions))
```

## Documentation

- **[USAGE.md](USAGE.md)** - Comprehensive usage guide covering common scenarios, patterns, and best practices
- **[PERFORMANCE.md](PERFORMANCE.md)** - Detailed performance optimization strategies and capacity planning
- **[PARSER_LAYERS.md](PARSER_LAYERS.md)** - Low-level and high-level parser architecture and API boundaries
- **[SYSML_DATA_MODEL.md](SYSML_DATA_MODEL.md)** - Currently implemented high-level SysML data model inventory
- **[Examples](../examples/)** - Working code examples:
  - [basic](../examples/basic/) - Simple parsing example
  - [requirements](../examples/requirements/) - Requirement traceability
  - [validation](../examples/validation/) - Error handling
  - [parallel](../examples/parallel/) - Parallel parsing
  - [visitor](../examples/visitor/) - Custom visitors
- **[Project Usage Guide](../docs/usage-guide.md)** - Project-level documentation covering command-line tools and validation data

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
