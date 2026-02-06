# gosysml2 Usage Guide

A comprehensive guide for using the gosysml2 library to parse and work with SysML v2 models in Go.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Parsing Strategies](#parsing-strategies)
3. [Working with Models](#working-with-models)
4. [Common Patterns](#common-patterns)
5. [Error Handling](#error-handling)
6. [Performance Tips](#performance-tips)

## Getting Started

### Installation

```bash
go get github.com/dVoo/gosysml2
```

### First Parse Example

Here's a minimal example to get you parsing SysML models:

```go
package main

import (
    "fmt"
    "github.com/dVoo/gosysml2/sysml"
)

func main() {
    // Simple SysML input
    input := `
        package VehicleModel {
            part def Vehicle {
                part engine : Engine;
            }
            part def Engine;
        }
    `

    // Parse the input
    result := sysml.ParseString(input)
    
    // Check for errors
    if !result.Success() {
        fmt.Printf("Parse errors: %s\n", result.Errors)
        return
    }
    
    // Access the model
    for _, pkg := range result.Model.Packages {
        fmt.Printf("Package: %s\n", pkg.Name())
    }
}
```

### Understanding the Result

The `ParseResult` struct provides:

- `Model` - The parsed model containing all elements
- `Errors` - Collection of parse errors (if any)
- `Source` - The source file path or identifier
- `Tree` - The raw parse tree (unless discarded)

```go
result := sysml.ParseString(input)

if result.Success() {
    // Access parsed model
    model := result.Model
    
    // Iterate packages
    for _, pkg := range model.Packages {
        fmt.Printf("Package: %s\n", pkg.Name())
    }
} else {
    // Handle errors
    for _, err := range result.Errors.Errors {
        fmt.Printf("Line %d: %s\n", err.Line, err.Message)
    }
}
```

## Parsing Strategies

### Single File Parsing

Parse a single SysML file:

```go
result := sysml.ParseFile("model.sysml")
if !result.Success() {
    log.Fatal(result.Errors)
}

// Process the model
processModel(result.Model)
```

Parse from bytes (avoids string copy):

```go
data, err := os.ReadFile("model.sysml")
if err != nil {
    log.Fatal(err)
}

result := sysml.ParseBytes(data, "model.sysml")
```

### Directory Parsing (Sequential)

Parse all `.sysml` files in a directory:

```go
results, err := sysml.ParseDirectory("./models")
if err != nil {
    log.Fatal(err)
}

for _, r := range results {
    if r.Success() {
        fmt.Printf("Parsed %s: %d packages\n", 
            r.Source, len(r.Model.Packages))
    } else {
        fmt.Printf("Failed %s: %s\n", r.Source, r.Errors)
    }
}
```

### Parallel Parsing for Large Repositories

For multi-core machines, parallel parsing significantly improves throughput:

```go
// Use all available CPU cores
results, err := sysml.ParseDirectoryParallel("./models", 0)

// Or specify worker count
results, err := sysml.ParseDirectoryParallel("./models", 4)

// Process results
for _, r := range results {
    if r.Success() {
        // Process each model
        reqs := sysml.FindRequirements(r.Model)
        fmt.Printf("%s: %d requirements\n", r.Source, len(reqs))
    }
}
```

See [examples/parallel](examples/parallel) for a complete demonstration.

### Streaming for Memory-Constrained Scenarios

Process files one at a time without loading all into memory:

```go
err := sysml.ParseDirectoryStream("./models", func(r *sysml.ParseResult) error {
    if r.Success() {
        // Process immediately, then discard
        reqs := sysml.FindRequirements(r.Model)
        fmt.Printf("%s: %d requirements\n", r.Source, len(reqs))
    }
    return nil // continue to next file
})
```

## Working with Models

### Accessing Packages

```go
// Get top-level packages
for _, pkg := range model.Packages {
    fmt.Printf("Package: %s\n", pkg.Name())
    
    // Access nested packages
    for _, subPkg := range pkg.Packages() {
        fmt.Printf("  Sub-package: %s\n", subPkg.Name())
    }
}

// Find a specific package by name
vehiclePkg := model.FindPackage("VehicleModel")
if vehiclePkg != nil {
    // Work with the package
}
```

### Finding Elements by Type

The library provides typed finder functions:

```go
// Find all elements of a specific type
parts := sysml.FindParts(model)
requirements := sysml.FindRequirements(model)
verifications := sysml.FindVerifications(model)
concerns := sysml.FindConcerns(model)
useCases := sysml.FindUseCases(model)
analysisCases := sysml.FindAnalysisCases(model)
actions := sysml.FindActions(model)
attributes := sysml.FindAttributes(model)
ports := sysml.FindPorts(model)

// Find definitions vs usages
definitions := sysml.FindDefinitions(model) // []Definition interface
usages := sysml.FindUsages(model)          // []Usage interface
```

### Using Typed Accessors

Instead of casting elements, use typed accessors on packages:

```go
pkg := model.FindPackage("VehicleModel")

// Get typed children (no casting needed!)
for _, part := range pkg.Parts() {
    fmt.Printf("Part: %s\n", part.Name())
    
    // Access nested parts
    for _, subPart := range part.Parts() {
        fmt.Printf("  Sub-part: %s\n", subPart.Name())
    }
    
    // Access attributes
    for _, attr := range part.Attributes() {
        fmt.Printf("  Attribute: %s\n", attr.Name())
    }
}

// Other accessors
for _, req := range pkg.Requirements() {
    fmt.Printf("Requirement: %s\n", req.Name())
}

for _, action := range pkg.Actions() {
    fmt.Printf("Action: %s\n", action.Name())
}
```

### Working with References

The model uses **type-safe element references** instead of strings:

```go
type Ref[T Element] struct {
    // Name() returns the reference name
    // Resolved() returns the resolved element (nil if unresolved)
    // IsResolved() returns true if resolved
}
```

Example usage:

```go
// Find a requirement
reqs := sysml.FindRequirements(model)
for _, req := range reqs {
    // Access derived requirements (real pointers!)
    for _, derived := range req.DerivedFrom {
        fmt.Printf("%s is derived from %s\n", 
            req.Name(), derived.Name())
    }
    
    // Check verification cases
    for _, ver := range req.VerifiedBy {
        fmt.Printf("%s is verified by %s\n", 
            req.Name(), ver.Name())
    }
    
    // Access subject if resolved
    if req.Subject.IsResolved() {
        fmt.Printf("Subject: %s\n", 
            req.Subject.Resolved().Name())
    }
}

// For usages, access the type definition
parts := sysml.FindParts(model)
for _, part := range parts {
    if !part.IsDefinition && part.TypeRef.IsResolved() {
        def := part.TypeRef.Resolved()
        fmt.Printf("Part %s is of type %s\n", 
            part.Name(), def.Name())
    }
}
```

### Qualified Names and Lookups

```go
// Elements have qualified names
elem := model.FindByQualifiedName("Vehicle::Engine::power")
fmt.Println(elem.QualifiedName())  // "Vehicle::Engine::power"

// Or use the finder function
elem = sysml.FindByQualifiedName(model, "Vehicle::Engine")
```

## Common Patterns

### Extracting Requirements

```go
result := sysml.ParseFile("requirements.sysml")
if !result.Success() {
    log.Fatal(result.Errors)
}

for _, req := range sysml.FindRequirements(result.Model) {
    fmt.Printf("Requirement: %s\n", req.Name())
    
    // Access requirement ID
    if req.RequirementID != "" {
        fmt.Printf("  ID: %s\n", req.RequirementID)
    }
    
    // Access documentation
    if doc := req.Documentation(); doc != "" {
        fmt.Printf("  Documentation: %s\n", doc)
    }
    
    // Show traceability
    if len(req.DerivedFrom) > 0 {
        fmt.Printf("  Derived from:\n")
        for _, derived := range req.DerivedFrom {
            fmt.Printf("    - %s\n", derived.Name())
        }
    }
    
    if len(req.VerifiedBy) > 0 {
        fmt.Printf("  Verified by:\n")
        for _, ver := range req.VerifiedBy {
            fmt.Printf("    - %s\n", ver.Name())
        }
    }
}
```

See [examples/requirements](examples/requirements) for a complete example.

### Finding Parts and Their Types

```go
parts := sysml.FindParts(model)

for _, part := range parts {
    if part.IsDefinition {
        fmt.Printf("Part definition: %s\n", part.Name())
    } else {
        // It's a usage - show its type
        fmt.Printf("Part usage: %s", part.Name())
        if part.TypeRef.IsResolved() {
            fmt.Printf(" : %s", part.TypeRef.Resolved().Name())
        }
        fmt.Println()
    }
}
```

### Requirement Traceability

```go
// Build a traceability matrix
requirements := sysml.FindRequirements(model)

for _, req := range requirements {
    fmt.Printf("\n%s:\n", req.Name())
    
    // Show derivation chain
    if len(req.DerivedFrom) > 0 {
        fmt.Println("  Derived from:")
        for _, parent := range req.DerivedFrom {
            fmt.Printf("    - %s\n", parent.Name())
        }
    }
    
    // Show derived requirements (inverse)
    if len(req.DerivedReqs) > 0 {
        fmt.Println("  Derived requirements:")
        for _, child := range req.DerivedReqs {
            fmt.Printf("    - %s\n", child.Name())
        }
    }
    
    // Show satisfaction
    if len(req.SatisfiedBy) > 0 {
        fmt.Println("  Satisfied by:")
        for _, elem := range req.SatisfiedBy {
            fmt.Printf("    - %s (%s)\n", 
                elem.Name(), elem.Kind())
        }
    }
    
    // Show verification
    if len(req.VerifiedBy) > 0 {
        fmt.Println("  Verified by:")
        for _, ver := range req.VerifiedBy {
            fmt.Printf("    - %s\n", ver.Name())
        }
    }
}
```

### Custom Filtering

Use the `Filter` function for custom predicates:

```go
// Find all parts that are definitions
definitions := sysml.Filter(model, func(e sysml.Element) bool {
    if part, ok := e.(*sysml.Part); ok {
        return part.IsDefinition
    }
    return false
})

// Find requirements with specific prefix
reqPrefix := "SYS-"
systemReqs := sysml.Filter(model, func(e sysml.Element) bool {
    if req, ok := e.(*sysml.Requirement); ok {
        return strings.HasPrefix(req.Name(), reqPrefix)
    }
    return false
})

// Find elements by name pattern
matches := sysml.Filter(model, func(e sysml.Element) bool {
    return strings.Contains(e.Name(), "Engine")
})
```

### Using the Visitor Pattern

Implement custom visitors for model analysis:

```go
// Custom visitor that counts specific element types
type MyVisitor struct {
    sysml.BaseVisitor
    partCount int
    reqCount  int
}

func (v *MyVisitor) VisitPart(part *sysml.Part) bool {
    v.partCount++
    return true // continue visiting children
}

func (v *MyVisitor) VisitRequirement(req *sysml.Requirement) bool {
    v.reqCount++
    return true
}

// Use the visitor
visitor := &MyVisitor{}
sysml.Visit(model, visitor)
fmt.Printf("Found %d parts and %d requirements\n", 
    visitor.partCount, visitor.reqCount)

// Or use the built-in counter
counter := sysml.NewCounter()
sysml.Visit(model, counter)
fmt.Printf("Counts: %v\n", counter.Counts)
fmt.Printf("Total: %d\n", counter.Total())
```

See [examples/visitor](examples/visitor) for a complete demonstration.

### Walking the Model

```go
// Walk with depth tracking
sysml.Walk(model, func(elem sysml.Element, depth int) bool {
    indent := strings.Repeat("  ", depth)
    fmt.Printf("%s%s: %s\n", indent, elem.Kind(), elem.Name())
    return true // continue walking
})
```

### Generic FindAll (Go 1.23+)

```go
// Use generics to find all elements of a specific type
requirements := sysml.FindAll[*sysml.Requirement](model)
parts := sysml.FindAll[*sysml.Part](model)
verifications := sysml.FindAll[*sysml.Verification](model)

// Iterate over all elements
for elem := range sysml.All(model) {
    fmt.Printf("%s: %s\n", elem.Kind(), elem.Name())
}

// Iterate over specific kind
for elem := range sysml.OfKind(model, sysml.KindRequirement) {
    fmt.Printf("Requirement: %s\n", elem.Name())
}
```

## Error Handling

### Checking for Errors

```go
result := sysml.ParseFile("model.sysml")

if !result.Success() {
    // Get first error
    first := result.Errors.First()
    fmt.Printf("First error at line %d, column %d: %s\n",
        first.Line, first.Column, first.Message)
    
    // Iterate all errors
    for _, err := range result.Errors.Errors {
        fmt.Printf("Line %d:%d - %s\n",
            err.Line, err.Column, err.Message)
    }
    
    // Get formatted error string
    fmt.Println(result.Errors.Error())
}
```

### Accessing Error Details

```go
result := sysml.ParseString(input)
if !result.Success() {
    for _, err := range result.Errors.Errors {
        fmt.Printf("Error:\n")
        fmt.Printf("  Line: %d\n", err.Line)
        fmt.Printf("  Column: %d\n", err.Column)
        fmt.Printf("  Message: %s\n", err.Message)
        fmt.Printf("  Source: %s\n", err.Source)
    }
}
```

### Handling Partial Success

Some parse operations may succeed partially:

```go
result := sysml.ParseString(input)

if result.Model != nil {
    // Model was built, but there may be errors
    fmt.Printf("Model has %d packages\n", len(result.Model.Packages))
    
    if !result.Success() {
        fmt.Printf("With %d errors\n", len(result.Errors.Errors))
    }
}
```

### Validation Without Full Parse

For quick syntax checking without building a model:

```go
import "github.com/dVoo/gosysml2/low"

// Validate only (faster, no model built)
errors := low.Validate(input)
if errors.HasErrors() {
    fmt.Printf("Validation failed: %s\n", errors)
} else {
    fmt.Println("Validation passed!")
}
```

See [examples/validation](examples/validation) for a complete example.

## Performance Tips

### When to Use WithDiscardTree

The parse tree can consume significant memory. Discard it if you only need the model:

```go
// Saves ~30% memory
result := sysml.ParseFile(path, sysml.WithDiscardTree())

// Also works with directory parsing
results, err := sysml.ParseDirectory(dir, sysml.WithDiscardTree())

// And streaming
err := sysml.ParseDirectoryStream(dir, handler, sysml.WithDiscardTree())
```

### Memory Optimization

```go
// 1. Use ParseBytes when you have []byte (avoids copy)
data, _ := os.ReadFile(path)
result := sysml.ParseBytes(data, path)

// 2. Discard parse tree
result := sysml.ParseFile(path, sysml.WithDiscardTree())

// 3. Use streaming for very large repositories
err := sysml.ParseDirectoryStream(dir, func(r *sysml.ParseResult) error {
    // Process and discard each file immediately
    processModel(r.Model)
    return nil
})
```

### Capacity Guidelines

| Repository Size | Recommended Approach | 32GB RAM |
|-----------------|---------------------|----------|
| < 100 MB | `ParseDirectory()` | Easy |
| 100-500 MB | `ParseDirectoryParallel()` + `WithDiscardTree()` | OK |
| 500 MB - 1 GB | `ParseDirectoryStream()` | Possible |
| > 1 GB | See PERFORMANCE.md | Challenging |

### Parallel Parsing Best Practices

```go
// Use all CPU cores (good default)
results, _ := sysml.ParseDirectoryParallel(dir, 0)

// Or specify based on your workload
// CPU-bound: match CPU cores
// IO-bound: can use more workers
results, _ := sysml.ParseDirectoryParallel(dir, runtime.NumCPU())

// For small repositories, sequential may be faster
// (parallelization overhead exceeds benefit)
if repoSize < 50*1024*1024 { // < 50MB
    results, _ = sysml.ParseDirectory(dir)
}
```

## Additional Resources

- [API Reference](README.md) - Complete API documentation
- [Performance Guide](PERFORMANCE.md) - Detailed optimization strategies
- [Examples](examples/) - Working code examples
  - [Basic](examples/basic) - Simple parsing example
  - [Requirements](examples/requirements) - Requirement traceability
  - [Validation](examples/validation) - Error handling
  - [Parallel](examples/parallel) - Parallel parsing
  - [Visitor](examples/visitor) - Custom visitors
