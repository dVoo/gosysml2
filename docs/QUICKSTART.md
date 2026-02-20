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

    result := sysml.ParseString(input)
    if !result.Success() {
        fmt.Printf("parse errors: %s\n", result.Errors)
        return
    }

    parts := sysml.FindAll[*sysml.Part](result.Model)
    fmt.Printf("parsed %d part elements\n", len(parts))
}
```

## 3. Parse a File

```go
result := sysml.ParseFile("model.sysml")
if !result.Success() {
    panic(result.Errors)
}
```

## 4. Parse a Directory

```go
results, err := sysml.ParseDirectory("./models")
if err != nil {
    panic(err)
}
for _, r := range results {
    if r.Success() {
        fmt.Printf("ok: %s\n", r.Source)
    }
}
```

## 5. Validate Syntax Only (Low-Level)

Use `low` for pure syntax checks and ANTLR-level work.

```go
errs := low.Validate(`package P { part def A; }`)
if errs.HasErrors() {
    fmt.Println(errs)
}
```

## 6. Next Reading

- [`LOW_LEVEL_GUIDE.md`](LOW_LEVEL_GUIDE.md)
- [`HIGH_LEVEL_GUIDE.md`](HIGH_LEVEL_GUIDE.md)
- [`EXAMPLES.md`](EXAMPLES.md)
