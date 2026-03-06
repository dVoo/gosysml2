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
- `sysml.ParseStringContext(ctx, input, opts...)`
- `sysml.ParseFile(path, opts...)`
- `sysml.ParseFileContext(ctx, path, opts...)`
- `sysml.ParseBytes(data, source, opts...)`
- `sysml.ParseBytesContext(ctx, data, source, opts...)`
- `sysml.ParseReader(r, source, opts...)`
- `sysml.ParseReaderContext(ctx, r, source, opts...)`
- `sysml.ParseStringModel(input, opts...) (*Model, error)`
- `sysml.ParseFileModel(path, opts...) (*Model, error)`
- `sysml.ParseDir(ctx, dir, opts) iter.Seq[*ParseResult]`

## Data Model Snapshot (Currently Implemented)

Main root model access:

- `Model.Elements` (canonical top-level storage)
- typed accessors derived from `Elements`:
  - `Packages()`, `Imports()`, `Comments()`, `Dependencies()`, `Docs()`
  - `Flows()`, `ControlNodes()`, `Occurrences()`
  - `Aliases()`, `Metadata()`, `Renderings()`, `Messages()`
  - `Filters()`, `Satisfies()`, `Verifies()`

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
if err := result.Err(); err != nil {
    panic(err)
}

parts := sysml.FindAll[*sysml.Part](result.Model)
requirements := sysml.FindAll[*sysml.Requirement](result.Model)
fmt.Println(len(parts), len(requirements))
```

## Example: Work with Type-Safe References

```go
for _, p := range sysml.FindAll[*sysml.Part](result.Model) {
    if p.Role() == sysml.RoleUsage {
        fmt.Printf("%s : %s\n", p.Name(), p.TypeName())
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
- `sysml.WithoutResolution()` to skip index build and reference resolution
- `sysml.WithoutModelBuild()` for syntax-only parsing
- `sysml.WithContentHash()` to compute/store SHA-256 during parse
- `sysml.WithSource(source)` to set source id for in-memory parses
- `sysml.WithoutCompatibilityRewrites()` to disable pre-parse compatibility rewrites
- `sysml.WithStandardLibrary()` to load/resolve against standard libraries
- `sysml.WithLibraryPath(path)` for custom library location
- `sysml.WithLibraryRegistry(reg)` for preloaded registry usage
- `sysml.WithParseCache(cache)` to reuse an explicit parse cache handle
- `sysml.WithDefaultParseCache()` for a lazy temp-backed cache
- `sysml.WithParseCacheDir(dir)` for a caller-managed persistent cache directory

## ParseResult Contract

`ParseResult` now exposes one canonical parse status API:

- `Err() error` (nil means success)
- `Errors() []*Error` (flat list)
- `ParseError *ParseError` (detailed aggregate when present)

And parse metadata:

- `Source string`
- `Hash string` (SHA-256 of original input when `WithContentHash()` is used)
- `Rewrites []string` (compatibility rewrites applied)
- `Tree antlr.Tree` (optional; omitted with `WithDiscardTree`)

## Directory Parsing

Use one API for sequential, parallel, and streaming patterns:

```go
opts := sysml.DirOptions{
    Workers:      0, // 0 => runtime.NumCPU()
    ParseOptions: []sysml.ParseOption{sysml.WithDiscardTree()},
}

for r := range sysml.ParseDir(ctx, "./models", opts) {
    if err := r.Err(); err != nil {
        fmt.Printf("failed %s: %v\n", r.Source, err)
        continue
    }
    fmt.Printf("ok: %s\n", r.Source)
}
```

`ParseDir` ordering rules:

- `Workers: 1` preserves walk order.
- `Workers > 1` does not preserve result order.

For repeated repository parses, the optimized interface is still `ParseDir`, but
with an attached cache option:

```go
cache, err := sysml.NewParseCache(
    sysml.WithCacheDir(".gosysml2-cache"),
    sysml.WithCachePersistence(true),
)
if err != nil {
    panic(err)
}
defer cache.Close()

opts := sysml.DirOptions{
    Workers: 0,
    ParseOptions: []sysml.ParseOption{
        sysml.WithParseCache(cache),
        sysml.WithDiscardTree(),
    },
}

for r := range sysml.ParseDir(ctx, "./models", opts) {
    if err := r.Err(); err != nil {
        continue
    }
    _ = r.Model
}
```

Cache behavior summary:

- no cache is used unless explicitly requested
- `ParseFile` and `ParseDir` reuse compatible cached results automatically
- `ParseBytes` participates only when `source` names a stable filesystem path
- `ParseDir` reparses changed files, rebuilds a shared repo index, and selectively re-resolves affected files

## Traversal APIs

- visitor API: `sysml.Visit(model, visitor)`
- generic walk: `sysml.Walk(model, fn)`
- iterators and filters: `All`, `OfType[T]`, `OfKind`, `Filter`, `FindAll[T]`

## Semantic Roles and Typed Usages

Prefer the method-based semantic API when tooling needs stable
definition-vs-usage classification and unresolved/resolved type names:

- `elem.Role() ElementRole`
- `usage.TypeName() string`

Convenience wrappers (`RoleOf`, `UsageTypeName`, `IsDefinitionElement`,
`IsUsageElement`) remain available but are not the canonical path.

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
