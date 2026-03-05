# gosysml2 Usage Guide

This guide covers the current high-level `sysml` API and the low-level `low` API.

## Getting Started

### Install

```bash
go get github.com/dVoo/gosysml2
```

### Parse In-Memory Text

```go
package main

import (
    "fmt"

    "github.com/dVoo/gosysml2/sysml"
)

func main() {
    input := `
package VehicleModel {
    part def Engine;
    part car : Engine;
}
`

    result := sysml.ParseString(input, sysml.WithSource("memory://VehicleModel.sysml"))
    if err := result.Err(); err != nil {
        fmt.Printf("parse failed: %v\n", err)
        return
    }

    for _, pkg := range result.Model.Packages() {
        fmt.Printf("package: %s\n", pkg.Name())
    }
}
```

## ParseResult Model

`ParseResult` has one canonical status API:

- `Err() error`: `nil` means success
- `Errors() []*sysml.Error`: flat list for iteration
- `ParseError *sysml.ParseError`: aggregate error object when parse failed

Other useful fields:

- `Model *sysml.Model`
- `Source string`
- `Hash string` (SHA-256 of original input)
- `Rewrites []string` (compatibility rewrites applied)
- `Tree antlr.Tree` (unless discarded with `WithDiscardTree`)

Example:

```go
result := sysml.ParseFile("model.sysml")
if err := result.Err(); err != nil {
    fmt.Println("first:", result.ParseError.First())
    for _, e := range result.Errors() {
        fmt.Printf("%d:%d %s\n", e.Line, e.Column, e.Message)
    }
    return
}
```

## Parsing APIs

### Single Input

- `ParseString`, `ParseFile`, `ParseBytes`, `ParseReader`
- Context-aware variants: `ParseStringContext`, `ParseFileContext`, `ParseBytesContext`, `ParseReaderContext`
- Idiomatic helpers: `ParseStringModel`, `ParseFileModel`

### Directory Parsing (Unified)

Use `ParseDir` for sequential, parallel, and streaming patterns.

```go
import (
    "context"
    "slices"

    "github.com/dVoo/gosysml2/sysml"
)

opts := sysml.DirOptions{
    Workers:      0, // 0 => runtime.NumCPU()
    ParseOptions: []sysml.ParseOption{sysml.WithDiscardTree()},
}

// streaming style
for r := range sysml.ParseDir(context.Background(), "./models", opts) {
    if err := r.Err(); err != nil {
        continue
    }
    _ = r.Model
}

// collect style
results := slices.Collect(sysml.ParseDir(context.Background(), "./models", opts))
_ = results
```

## Parse Options

- `WithDiscardTree()`
- `WithSource(source)` for in-memory parse source labels
- `WithoutCompatibilityRewrites()` for strict grammar mode
- `WithStandardLibrary()`
- `WithLibraryPath(path)`
- `WithLibraryRegistry(reg)`

## Working with Models

`Model` uses one canonical top-level store: `Elements []Element`.
Typed top-level views are accessor methods:

- `Packages()`, `Imports()`, `Comments()`, `Dependencies()`, `Docs()`
- `Flows()`, `ControlNodes()`, `Occurrences()`
- `Aliases()`, `Metadata()`, `Renderings()`, `Messages()`
- `Filters()`, `Satisfies()`, `Verifies()`

Example:

```go
for _, pkg := range model.Packages() {
    fmt.Println(pkg.Name())
}

parts := sysml.FindAll[*sysml.Part](model)
requirements := sysml.FindAll[*sysml.Requirement](model)
_ = parts
_ = requirements
```

## Traversal and Queries

- `Visit(model, visitor)` for structured processors
- `Walk(model, fn)` for recursive traversal with early exit
- `All(model)`, `OfType[T](model)`, `OfKind(model, kind)` for iterator-style loops
- `Filter`, `FindByName`, `FindByKind`, `FindByQualifiedName`

Which one to use:

- Prefer iterators (`All`, `OfType[T]`, `OfKind`) for most filtering/collection tasks.
- Use `Walk` when you need depth and short-circuit traversal (`false` stops descent).
- Use `Visit` for larger structured processors that benefit from per-type methods.
- Use `WalkAll` when you want traversal callbacks without early-exit logic.

Example:

```go
for elem := range sysml.OfType[*sysml.Requirement](model) {
    fmt.Println(elem.Name())
}
```

## Position-Based Lookup

Use `ElementAt` for editor/LSP style lookups of named elements.
Use `ElementAtIncludingUnnamed` when anonymous usages should still be discoverable.

```go
elem := sysml.ElementAt(model, line, column) // zero-based line/column
if elem != nil {
    fmt.Println(elem.QualifiedName())
}

elem = sysml.ElementAtIncludingUnnamed(model, line, column)
_ = elem
```

## References

`Ref[T]` supports:

- `Name()`
- `EffectiveName()`
- `IsResolved()`
- `Resolved()`

Example:

```go
for _, p := range sysml.FindAll[*sysml.Part](model) {
    if p.Role() == sysml.RoleUsage {
        fmt.Printf("%s : %s\n", p.Name(), p.TypeName())
    }
}
```

## Semantic Roles (Definition vs Usage)

For tooling code that needs stable semantic classification, prefer `elem.Role()`
and `usage.TypeName()` over reflecting on concrete `IsDefinition` fields:

```go
for elem := range sysml.All(model) {
    switch elem.Role() {
    case sysml.RoleDefinition:
        // ...
    case sysml.RoleUsage:
        // ...
    }
}

if usage, ok := elem.(sysml.Usage); ok {
    fmt.Println(usage.TypeName()) // resolved name or unresolved declared name
}
```

## Low-Level Validation

For syntax-only validation (no high-level model):

```go
import "github.com/dVoo/gosysml2/low"

errs := low.Validate(input)
if errs.HasErrors() {
    fmt.Println(errs)
}
```

## Performance Notes

- Use `WithDiscardTree()` to lower memory usage.
- Use `ParseDir` with `Workers: 0` for parallel throughput.
- Use `ParseDir` with `Workers: 1` and process each result immediately for streaming-like memory behavior.
