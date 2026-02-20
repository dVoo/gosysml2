# Low-Level Guide (`low`)

Use the `low` package when you need direct parser/tree behavior and minimal overhead.

## Typical Scenarios

- syntax validation without building a high-level model
- custom parse-tree visitors/listeners over ANTLR contexts
- token-level tooling
- parser-rule/token-name alignment with official grammar docs

## Core Entry Points

- `low.Parse(input string, opts ...ParseOption)`
- `low.ParseBytes(input []byte, opts ...ParseOption)`
- `low.ParseWithContext(ctx, input, opts...)`
- `low.Validate(input string)`
- `low.ValidateBytes(input []byte)`

## Example: Parse and Handle Errors

```go
package main

import (
    "fmt"

    "github.com/dVoo/gosysml2/low"
)

func main() {
    tree, errs := low.Parse(`package P { part def A; }`)
    if errs.HasErrors() {
        for _, e := range errs.All() {
            fmt.Printf("%d:%d %s\n", e.Line, e.Column, e.Message)
        }
        return
    }
    fmt.Printf("tree root: %T\n", tree)
}
```

## Example: Validation-Only Fast Path

```go
errs := low.Validate(input)
if errs.HasErrors() {
    fmt.Println("invalid syntax")
}
```

## Example: Context Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

_, err := low.ParseWithContext(ctx, input)
if err != nil {
    // context timeout/cancel or parser-level error
}
```

## Parse Options

- `low.WithParseTree(false)` for syntax checks with less memory
- `low.WithContext(ctx)` for cancellation support

## Grammar Name Helpers

Use these when mapping official KEBNF names to generated ANTLR names:

- `low.NormalizeGrammarName(name)`
- `low.ResolveParserRuleName(name)`
- `low.ResolveTokenName(name)`
- `low.ParserRuleNameCandidates(name)`
- `low.TokenNameCandidates(name)`

## When to Move Up to `sysml`

If you need resolved references, typed elements, and semantic traversal, use the high-level API: [`HIGH_LEVEL_GUIDE.md`](HIGH_LEVEL_GUIDE.md).
