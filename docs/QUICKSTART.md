# Quick Start

This page gets you from zero to useful parsing quickly.

## 1. Install

```bash
go get github.com/dVoo/gosysml2
```

## 2. Parse a String (High-Level)

Use `sysml` when you want typed model elements and resolved references.

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
    part def Vehicle {
        part engine : Engine;
    }
}
`

    result := sysml.ParseString(input, sysml.WithSource("memory://VehicleModel.sysml"))
    if err := result.Err(); err != nil {
        fmt.Printf("parse errors: %s\n", err)
        return
    }

    parts := sysml.FindAll[*sysml.Part](result.Model)
    fmt.Printf("parsed %d part elements\n", len(parts))
}
```

## 3. Parse a File

```go
result := sysml.ParseFile("model.sysml")
if err := result.Err(); err != nil {
    panic(err)
}
```

## 4. Parse a Directory

```go
import (
    "context"
    "fmt"
    "slices"

    "github.com/dVoo/gosysml2/sysml"
)

opts := sysml.DirOptions{
    Workers:      0, // 0 = NumCPU
    ParseOptions: []sysml.ParseOption{sysml.WithDiscardTree()},
}

results := slices.Collect(sysml.ParseDir(context.Background(), "./models", opts))
for _, r := range results {
    if r.Err() == nil {
        fmt.Printf("ok: %s\n", r.Source)
    } else {
        fmt.Printf("failed: %s: %v\n", r.Source, r.Err())
    }
}
```

When `Workers > 1`, `ParseDir` results are not ordered.

## 5. Parse a Directory with Cache

For repeated repository parses, keep a cache handle and reuse it:

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

for r := range sysml.ParseDir(context.Background(), "./models", opts) {
    if err := r.Err(); err != nil {
        fmt.Printf("failed: %s: %v\n", r.Source, err)
        continue
    }
    fmt.Printf("ok: %s\n", r.Source)
}
```

## 6. Validate Syntax Only (Low-Level)

Use `low` for pure syntax checks and ANTLR-level work.

```go
errs := low.Validate(`package P { part def A; }`)
if errs.HasErrors() {
    fmt.Println(errs)
}
```

## 7. Next Reading

- [`LOW_LEVEL_GUIDE.md`](LOW_LEVEL_GUIDE.md)
- [`HIGH_LEVEL_GUIDE.md`](HIGH_LEVEL_GUIDE.md)
- [`EXAMPLES.md`](EXAMPLES.md)
