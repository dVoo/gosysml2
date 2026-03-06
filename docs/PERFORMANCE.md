# Performance Guide

This document reflects the current parser behavior.

## Current Fast Paths

- `low` already parses lazily without eager token filling.
- `sysml.WithDiscardTree()` keeps the model and drops the parse tree.
- `sysml.ParseBytes(...)` avoids an extra high-level `[]byte` -> `string` copy.
- `sysml.ParseDir(...)` uses a bounded worker pool and streams files directly from `WalkDir`.
- `sysml.WithParseCache(...)` reuses compatible file parse results across repeated parses.

## Important Parse Options

- `WithDiscardTree()`: keep the model, drop the parse tree.
- `WithoutResolution()`: build the model, skip index build and reference resolution.
- `WithoutModelBuild()`: syntax parse only; no high-level model construction.
- `WithoutCompatibilityRewrites()`: skip compatibility rewrite passes.
- `WithContentHash()`: compute and retain SHA-256 content hash during parse.
- `WithParseCache(...)`: reuse an explicit cache handle.
- `WithDefaultParseCache()`: create a temp-backed disposable cache lazily.
- `WithParseCacheDir(dir)`: bind parsing to a caller-managed cache directory.

## What Costs Time

For the Batmobile sample benchmark:

- compatibility rewrite normalization is cheap relative to parse time
- model build is small relative to ANTLR parse time
- index build and reference resolution are also small on small models
- the dominant cost is still lexer/parser work in the ANTLR runtime

That means the highest-value wins are usually:

1. Avoid work you do not need.
2. Keep less data alive after parsing.
3. Keep repository parsing bounded and streaming.

## Recommended Usage

### Syntax validation only

```go
result := sysml.ParseBytes(data, source,
    sysml.WithoutModelBuild(),
    sysml.WithDiscardTree(),
)
```

### Model build without semantic resolution

```go
result := sysml.ParseBytes(data, source,
    sysml.WithoutResolution(),
    sysml.WithDiscardTree(),
)
```

### Full parse with lower memory

```go
result := sysml.ParseFile(path, sysml.WithDiscardTree())
```

### Repository-scale parsing

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
    Workers:      0, // 0 => runtime.NumCPU()
    ParseOptions: []sysml.ParseOption{
        sysml.WithParseCache(cache),
        sysml.WithDiscardTree(),
    },
}

for r := range sysml.ParseDir(ctx, repoDir, opts) {
    if err := r.Err(); err != nil {
        continue
    }
    // process r.Model
}
```

Notes:

- `Workers: 1` gives the most predictable streaming behavior.
- `Workers > 1` increases throughput but result order is not preserved.

## Practical Guidance

- Prefer `ParseBytes` when the caller already has `[]byte`.
- Prefer `ParseDir` plus `WithParseCache(...)` for repeated repository parses.
- Prefer `ParseFile` plus `WithParseCache(...)` for repeated single-file reparses from disk.
- Use `WithoutModelBuild()` for editor syntax checks and parser smoke tests.
- Use `WithoutResolution()` for partial semantic work where cross-reference resolution is not needed yet.
- Use `WithContentHash()` only when the caller actually needs hashes.
- Use `WithDiscardTree()` unless you explicitly need the raw parse tree.
- Use `WithDefaultParseCache()` for disposable session-local acceleration.
- Use `WithParseCacheDir(dir)` when cache persistence across process runs matters.

## Likely Future Work

- deeper profiling inside the ANTLR runtime hot path
- more aggressive cache invalidation heuristics and lighter-weight cache payloads
- optional lighter-weight model builders for specific tooling workflows
