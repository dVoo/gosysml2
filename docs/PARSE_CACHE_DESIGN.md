# Parse Cache

This document describes the current parse cache interface and behavior.

## Status

The parse cache is implemented and optional.

Use it through normal parse entry points:

- `ParseFile`
- `ParseBytes`
- `ParseDir`

by attaching one of:

- `WithParseCache(cache)`
- `WithDefaultParseCache()`
- `WithParseCacheDir(dir)`

## User-Facing API

### Parse options

```go
func WithParseCache(cache *ParseCache) ParseOption
func WithDefaultParseCache() ParseOption
func WithParseCacheDir(dir string) ParseOption
```

Behavior:

- `WithParseCache(cache)` reuses an existing cache handle.
- `WithDefaultParseCache()` lazily creates a temp-backed disposable cache handle.
- `WithParseCacheDir(dir)` lazily creates a cache rooted at `dir`.

### Cache handle

```go
type ParseCache struct { /* internal fields */ }

func NewParseCache(opts ...ParseCacheOption) (*ParseCache, error)
func (c *ParseCache) Dir() string
func (c *ParseCache) Close() error
func (c *ParseCache) Clear() error
```

Cache construction options:

```go
type ParseCacheOption func(*parseCacheConfig)

func WithCacheDir(dir string) ParseCacheOption
func WithTemporaryCacheDir() ParseCacheOption
func WithCachePersistence(enabled bool) ParseCacheOption
```

## Recommended Usage

### Repeated file parses

```go
cache, err := sysml.NewParseCache(
    sysml.WithCacheDir(".gosysml2-cache"),
    sysml.WithCachePersistence(true),
)
if err != nil {
    panic(err)
}
defer cache.Close()

result := sysml.ParseFile("model.sysml",
    sysml.WithParseCache(cache),
    sysml.WithDiscardTree(),
)
```

### Repeated repository parses

```go
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

## Current Behavior

### Default behavior

- cache is disabled unless explicitly requested
- uncached semantics and cached semantics are intended to match
- cache entries are disposable best-effort artifacts

### `ParseFile`

- hashes file content
- looks up a compatible cached entry
- reparses on miss
- stores the result back into the cache

### `ParseBytes`

- uses the cache only when `source` is a stable filesystem path
- skips cache for non-filesystem sources such as `memory://...`

### `ParseDir`

- walks repository files
- reuses compatible cached entries for unchanged files
- reparses cache misses
- removes cache entries for files that disappeared
- rebuilds a shared repository index
- re-resolves changed files and their reverse dependency closure
- refreshes manifest metadata and reverse dependencies

## Cache Keys

Cached file reuse is scoped by:

- absolute source path
- content hash
- parse mode fingerprint
- grammar version
- library fingerprint

`WithContentHash()` affects the returned `ParseResult.Hash`, but it does not
change semantic cache compatibility.

## Internal Invalidation Model

Each cached file tracks:

- content hash
- export signature
- top-level package names
- imports
- parse mode fingerprint
- library fingerprint

`ParseDir` maintains a reverse-dependency graph at file granularity.

The current dependency graph is conservative:

- files with overlapping top-level packages are treated as related
- files importing exported namespaces from another file are treated as dependents

If a file changes, `ParseDir` reparses that file and re-resolves its dependency
closure against the rebuilt repository index.

## Cache Directory Layout

```text
<cache-dir>/
  manifest.json
  lock
  files/
    <cache-entry-key>.json
```

Notes:

- `manifest.json` stores cache metadata and reverse dependencies
- `lock` is a best-effort lock-file placeholder
- `files/` stores per-entry metadata files for inspection/debugging

The in-memory cache handle also keeps live parse results for fast reuse within
the current process.

## Limits

Current implementation intentionally does not do:

- true incremental parsing inside a file
- parse-tree reuse across processes
- compact IR serialization for full model reload across process runs

Persistent cache directories currently preserve manifest and entry metadata.
Process-local cached `ParseResult` values are still the fast path for full reuse
inside one process.
