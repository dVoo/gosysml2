# High-Level Guide (`sysml`)

Use the `sysml` package for application code: typed model elements, index building, and reference resolution.

## Parse Flow

1. Parse text/file/directory
2. Build model (`*sysml.Model`)
3. Build index
4. Resolve references (`Ref[T]`)

The parse helpers already perform index build and reference resolution.

## Core Entry Points

- `sysml.ParseString(input, opts...)`
- `sysml.ParseFile(path, opts...)`
- `sysml.ParseBytes(data, source, opts...)`
- `sysml.ParseReader(r, source, opts...)`
- `sysml.ParseStringModel(input, opts...) (*Model, error)`
- `sysml.ParseFileModel(path, opts...) (*Model, error)`
- `sysml.ParseDirectory(dir, opts...)`
- `sysml.ParseDirectoryParallel(dir, workers, opts...)`
- `sysml.ParseDirectoryStream(dir, handler, opts...)`

## Data Model Snapshot (Currently Implemented)

Main root collections:

- `Model.Packages`
- `Model.Imports`
- `Model.Elements`
- relationship slices like `Model.Satisfies` and `Model.Verifies`

Common element types:

- structure: `Package`, `Part`, `Item`, `Attribute`, `Port`, `Interface`
- requirements: `Requirement`, `Verification`, `Case`, `AnalysisCase`, `Concern`
- behavior/flow: `Action`, `State`, `Transition`, `Flow`, `ControlNode`, `Occurrence`
- meta/docs: `Alias`, `Metadata`, `Rendering`, `Comment`, `Doc`, `Import`
- KerML support: `KerMLType`, `KerMLFeature`

For full field-level details, see [`SYSML_DATA_MODEL.md`](SYSML_DATA_MODEL.md).

## Example: Parse and Query Elements

```go
result := sysml.ParseString(input)
if !result.Success() {
    panic(result.Errors)
}

parts := sysml.FindAll[*sysml.Part](result.Model)
requirements := sysml.FindAll[*sysml.Requirement](result.Model)
fmt.Println(len(parts), len(requirements))
```

## Example: Work with Type-Safe References

```go
for _, p := range sysml.FindAll[*sysml.Part](result.Model) {
    if !p.IsDefinition && p.TypeRef.IsResolved() {
        fmt.Printf("%s : %s\n", p.Name(), p.TypeRef.Resolved().Name())
    }
}
```

## Example: Requirement Relationships

```go
for _, req := range sysml.FindAll[*sysml.Requirement](result.Model) {
    for _, parent := range req.DerivedFrom {
        fmt.Printf("%s derived from %s\n", req.Name(), parent.Name())
    }
    for _, ver := range req.VerifiedBy {
        fmt.Printf("%s verified by %s\n", req.Name(), ver.Name())
    }
}
```

## Short Name and Name Semantics

Identification forms map to separate fields:

- `Name()` is the declared name
- `DeclaredShortName()` preserves `<...>` short form when present

Reference resolution uses qualified names and short-name index paths.

## Options

- `sysml.WithDiscardTree()` for lower memory usage
- `sysml.WithStandardLibrary()` to load/resolve against standard libraries
- `sysml.WithLibraryPath(path)` for custom library location
- `sysml.WithLibraryRegistry(reg)` for preloaded registry usage

## Traversal APIs

- visitor API: `sysml.Visit(model, visitor)`
- generic walk: `sysml.Walk(model, fn)`
- iterators and filters: `All`, `OfType[T]`, `OfKind`, `Filter`, `FindAll[T]`

## View/Viewpoint Selection APIs

The high-level API provides direct selectors for model-designer-defined views:

- `FindView(model, nameOrQualifiedName)`
- `FindViewpoint(model, nameOrQualifiedName)`
- `ElementsForView(view)`
- `ElementsByView(model, nameOrQualifiedName)`
- `ViewsByViewpoint(model, nameOrQualifiedName)`
- `ElementsByViewpoint(model, nameOrQualifiedName)`

These APIs use parsed `view`/`viewpoint` elements and `expose` clauses, including
namespace/member wildcard forms and common filter annotations such as
`@SysML::PartUsage`.

## End-to-End Examples

See [`EXAMPLES.md`](EXAMPLES.md) for runnable programs.
