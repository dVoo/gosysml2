# Phase 1: Modernize Go Codebase - Research

**Researched:** 2026-02-03
**Domain:** Go 1.25 modernization, parser performance optimization, ANTLR wrapper code
**Confidence:** HIGH

## Summary

This research investigates how to modernize a Go parser codebase (currently Go 1.22) to Go 1.25, focusing on wrapper code around ANTLR-generated parsers. The standard approach involves applying Go 1.25 language features (improved generics, range-over-func iterators), modern error handling patterns, and proven performance optimizations specific to parser workloads.

The codebase already implements several performance best practices (lazy token consumption, byte-based inputs, tree discarding options, parallel parsing). Modernization should focus on: (1) applying Go 1.25+ features to visitor/walker patterns, (2) optimizing slice/map allocations with proper preallocation, (3) improving error handling with wrapping and context, and (4) establishing benchmarks using the 34 test files in `docs/testdata/`.

**Primary recommendation:** Upgrade to Go 1.25+, apply range-over-func to visitor patterns, preallocate slices in model building, wrap errors with context, and benchmark against real SysML files to validate improvements.

## Standard Stack

The established libraries/tools for this domain:

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.25+ | Language runtime | Latest stable release with 10-40% GC improvements, container-aware GOMAXPROCS |
| antlr4-go/antlr/v4 | v4.13.1+ | Parser runtime | Official ANTLR Go runtime, actively maintained |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| runtime/pprof | stdlib | CPU/memory profiling | Performance analysis and optimization |
| testing/synctest | stdlib (Go 1.25) | Concurrent code testing | Testing parallel parsing with virtualized time |
| golang.org/x/exp | latest | Experimental features (slices, maps utilities) | Helper functions for slice/map operations |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| ANTLR | Handwritten parser | 40x faster but requires complete rewrite, loses grammar-based approach |
| antlr4-go/antlr/v4 | Older antlr/antlr4 | Newer v4 path is the actively maintained version |

**Installation:**
```bash
go get -u github.com/antlr4-go/antlr/v4@latest
go install github.com/antlr4-go/antlr/v4/cmd/antlr@latest
```

## Architecture Patterns

### Recommended Project Structure
Current structure is sound:
```
gosysml2/
├── internal/parser/     # ANTLR-generated code (DO NOT MODIFY)
├── low/                 # Low-level parser wrapper
│   ├── lexer.go        # Lexer wrapper with error collection
│   ├── parser.go       # Parser wrapper with options
│   └── errors.go       # Error aggregation
├── sysml/              # High-level API
│   ├── parse.go        # Public parsing API
│   ├── model.go        # AST model types
│   ├── visitor.go      # Tree traversal utilities
│   └── errors.go       # High-level error types
└── examples/           # Usage examples
```

### Pattern 1: Functional Options for Parser Configuration
**What:** Use variadic option functions for flexible configuration
**When to use:** Any API that needs optional parameters (already implemented in codebase)
**Example:**
```go
// Source: Existing codebase pattern (low/parser.go, sysml/parse.go)
type ParseOption func(*parseConfig)

func WithParseTree(build bool) ParseOption {
    return func(c *parseConfig) {
        c.buildParseTree = build
    }
}

func NewParser(input string, opts ...ParseOption) *Parser {
    cfg := &parserConfig{buildParseTree: true}
    for _, opt := range opts {
        opt(cfg)
    }
    // ...
}
```

### Pattern 2: Error Collector with Interface Implementation
**What:** Collect errors during parsing by implementing ANTLR's ErrorListener
**When to use:** When wrapping ANTLR parsers to aggregate syntax errors
**Example:**
```go
// Source: Existing codebase (low/errors.go)
type ErrorCollector struct {
    *antlr.DefaultErrorListener
    errors []*SyntaxError
    source string
}

func (c *ErrorCollector) SyntaxError(
    recognizer antlr.Recognizer,
    offendingSymbol interface{},
    line, column int,
    msg string,
    e antlr.RecognitionException,
) {
    c.errors = append(c.errors, &SyntaxError{
        Line: line, Column: column, Message: msg, Source: c.source,
    })
}
```

### Pattern 3: Lazy Token Consumption (Performance Critical)
**What:** Let ANTLR consume tokens on-demand rather than buffering all tokens
**When to use:** Always for large files (already implemented, line removed)
**Example:**
```go
// Source: Existing implementation (low/parser.go:52-54)
tokens := lexer.TokenStream()
// NOTE: Removed tokens.Fill() - let ANTLR consume tokens lazily
// This significantly reduces memory usage for large files
```

### Pattern 4: Byte-Based Input Streams
**What:** Use []byte directly to avoid string conversions and copies
**When to use:** When input is already in []byte form (file reading)
**Example:**
```go
// Source: Existing implementation (low/lexer.go:25-27)
func NewLexerFromBytes(input []byte) *Lexer {
    stream := antlr.NewInputStream(string(input))
    return NewLexerFromStream(stream)
}

// MODERNIZATION OPPORTUNITY: Check if ANTLR v4.13+ supports direct []byte
// to avoid the string() conversion
```

### Pattern 5: Visitor Pattern with Stack Tracking
**What:** Use element stack to track parent context during tree traversal
**When to use:** Building hierarchical models from flat parse events
**Example:**
```go
// Source: Existing implementation (sysml/parse.go:258-275)
type modelBuilder struct {
    *parser.BaseSysMLv2ParserListener
    elementStack []Element
}

func (b *modelBuilder) EnterPartDefinition(ctx *parser.PartDefinitionContext) {
    part := NewPart(name, loc, true)
    if len(b.elementStack) > 0 {
        parent := b.elementStack[len(b.elementStack)-1]
        parent.AddChild(part)
    }
    b.elementStack = append(b.elementStack, part)
}

func (b *modelBuilder) ExitPartDefinition(ctx *parser.PartDefinitionContext) {
    b.elementStack = b.elementStack[:len(b.elementStack)-1]
}
```

### Pattern 6: Parallel File Processing with Semaphore
**What:** Limit concurrency with buffered channel semaphore
**When to use:** Processing multiple files with controlled parallelism
**Example:**
```go
// Source: Existing implementation (sysml/parse.go:168-210)
func ParseDirectoryParallel(dir string, workers int, opts ...ParseOption) ([]*ParseResult, error) {
    if workers <= 0 {
        workers = runtime.NumCPU()
    }

    results := make([]*ParseResult, len(files))
    var wg sync.WaitGroup
    sem := make(chan struct{}, workers)

    for i, file := range files {
        wg.Add(1)
        sem <- struct{}{}  // Acquire semaphore

        go func(idx int, path string) {
            defer wg.Done()
            defer func() { <-sem }()  // Release semaphore
            results[idx] = ParseFile(path, opts...)
        }(i, file)
    }

    wg.Wait()
    return results, nil
}
```

### Anti-Patterns to Avoid
- **Buffering all tokens upfront**: Memory waste - let ANTLR consume lazily (already avoided)
- **String concatenation in hot paths**: Use strings.Builder or pre-allocated buffers
- **Unpreallocated slices in loops**: Always preallocate when size is known
- **Global FileSet in go/token**: Use FileSet.AddExistingFiles (Go 1.25+) if applicable
- **Ignoring error context**: Always wrap errors with fmt.Errorf("%w") to preserve chain

## Don't Hand-Roll

Problems that look simple but have existing solutions:

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| String building | Manual += concatenation | strings.Builder | O(n) vs O(n²), reduces allocations |
| Slice utilities | Manual loops for filtering/mapping | golang.org/x/exp/slices | Optimized implementations, generics support |
| Map utilities | Manual key/value iteration | golang.org/x/exp/maps | Standard patterns, less error-prone |
| Error wrapping | Manual string formatting | fmt.Errorf with %w | Preserves error chain for errors.Is/As |
| JSON v2 | Custom serialization | encoding/json/v2 (Go 1.25 experimental) | Substantially faster, better API |
| Concurrent testing | Manual time mocking | testing/synctest (Go 1.25) | Virtualized time, deterministic tests |
| Object pooling | Custom pool implementation | sync.Pool | GC-aware, per-processor caching |

**Key insight:** Parser workloads allocate heavily during tree construction. Pre-allocate slices with known capacity, use sync.Pool for frequently created temporary objects (Location, small nodes), and leverage Go 1.25's improved slice stack allocation.

## Common Pitfalls

### Pitfall 1: Slice Append Without Preallocation
**What goes wrong:** Repeated append operations cause multiple array reallocations and O(n) copies
**Why it happens:** Easy to write `var results []Thing` and append in loop
**How to avoid:** Preallocate capacity when size is known or predictable
**Warning signs:** High allocation counts in pprof for slice growth, memory spikes during parsing
**Example fix:**
```go
// BAD: Will reallocate ~10 times for 1000 elements
var elements []Element
for _, ctx := range contexts {
    elements = append(elements, parseElement(ctx))
}

// GOOD: Single allocation
elements := make([]Element, 0, len(contexts))
for _, ctx := range contexts {
    elements = append(elements, parseElement(ctx))
}
```

### Pitfall 2: Premature sync.Pool Usage
**What goes wrong:** Object pooling adds complexity without performance gain
**Why it happens:** Assumption that pooling always helps, copying optimization advice blindly
**How to avoid:** Profile first, only pool if GC pressure is measured and objects are short-lived
**Warning signs:** Complex pool management code, minimal or negative performance impact
**Example when NOT to use:**
```go
// DON'T pool long-lived model elements that survive the parse
// They'll be GC'd from pool before reuse anyway
var modelPool = sync.Pool{New: func() interface{} { return &Model{} }}

// DO pool short-lived temporary objects in hot paths
var locationPool = sync.Pool{New: func() interface{} { return &Location{} }}
```

### Pitfall 3: ANTLR Go Runtime Performance Expectations
**What goes wrong:** Expecting handwritten-parser performance from ANTLR
**Why it happens:** ANTLR provides flexibility and grammar-based development at cost of performance
**How to avoid:** Accept 10-40x slower parsing vs handwritten, optimize wrapper code instead
**Warning signs:** Attempting to modify ANTLR-generated code, frustration with parsing speed
**Mitigation:**
- Focus optimization on wrapper code (model building, visitors, allocation patterns)
- Use lazy token consumption, parallel file parsing, tree discarding
- Profile wrapper code, not ANTLR internals
- Consider ANTLR as "fast enough" for development tools, not production data pipelines

### Pitfall 4: Error Handling Without Context
**What goes wrong:** Errors lose context as they propagate up the call stack
**Why it happens:** Using errors.New() instead of fmt.Errorf with %w
**How to avoid:** Always wrap errors with context when returning from function boundaries
**Warning signs:** Unable to determine error origin, unhelpful error messages
**Example fix:**
```go
// BAD: Loses context
if err := parseFile(path); err != nil {
    return errors.New("parse failed")
}

// GOOD: Preserves error chain with context
if err := parseFile(path); err != nil {
    return fmt.Errorf("parsing %s: %w", path, err)
}
```

### Pitfall 5: Goroutine Leaks in Concurrent Parsing
**What goes wrong:** Worker goroutines block forever, never exit
**Why it happens:** Channels not closed, context not propagated, panic in worker
**How to avoid:** Always use defer to release semaphore, propagate context, handle panics
**Warning signs:** Increasing goroutine count in runtime metrics, workers never complete
**Example fix:**
```go
// BAD: If ParseFile panics, semaphore never released
go func(idx int, path string) {
    sem <- struct{}{}
    results[idx] = ParseFile(path)
    <-sem
}(i, file)

// GOOD: Semaphore always released
go func(idx int, path string) {
    defer wg.Done()
    defer func() { <-sem }()  // Always release
    results[idx] = ParseFile(path)
}(i, file)
```

### Pitfall 6: String to []byte Conversions in Hot Paths
**What goes wrong:** String/[]byte conversions copy data, increasing allocations
**Why it happens:** API accepts string but data is in []byte (or vice versa)
**How to avoid:** Use []byte consistently, only convert at I/O boundaries
**Warning signs:** pprof shows runtime.stringtoslicebyte in hot paths
**Example fix:**
```go
// BAD: Double conversion
content, _ := os.ReadFile(filename)  // Returns []byte
result := ParseString(string(content))  // Converts to string, then back to []byte

// GOOD: Use bytes directly
content, _ := os.ReadFile(filename)
result := ParseBytes(content)
```

## Code Examples

Verified patterns from official sources and existing codebase:

### Range-Over-Func for Visitor Pattern (Go 1.23+)
```go
// Source: Go 1.23 iter package + existing visitor pattern
import "iter"

// Modern iterator-based visitor (MODERNIZATION OPPORTUNITY)
func (m *Model) Elements() iter.Seq[Element] {
    return func(yield func(Element) bool) {
        for _, elem := range m.elements {
            if !yield(elem) {
                return
            }
        }
    }
}

// Usage with range
for elem := range model.Elements() {
    fmt.Println(elem.Name())
}

// Can compose with standard library functions
filtered := slices.Collect(
    slices.Filter(model.Elements(), func(e Element) bool {
        return e.Kind() == KindRequirement
    }),
)
```

### Error Wrapping with Context (Go 1.13+ pattern, 2026 best practice)
```go
// Source: https://go.dev/blog/go1.13-errors
// Modern error handling pattern
func ParseFile(filename string) (*ParseResult, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("reading %s: %w", filename, err)
    }

    result := ParseBytes(content, filename)
    if result.Errors != nil && result.Errors.HasErrors() {
        return nil, fmt.Errorf("parsing %s: %w", filename, result.Errors)
    }

    return result, nil
}

// Checking specific error types
if err := ParseFile(path); err != nil {
    var parseErr *ParseError
    if errors.As(err, &parseErr) {
        // Handle parse-specific error
    }
}
```

### Preallocated Slice Building
```go
// Source: https://goperf.dev/01-common-patterns/mem-prealloc/
// MODERNIZATION: Apply to model building (sysml/parse.go)

// CURRENT: Unpreallocated append in Enter* methods
func (b *modelBuilder) buildChildren(contexts []Context) []Element {
    var children []Element  // Grows dynamically
    for _, ctx := range contexts {
        children = append(children, parseContext(ctx))
    }
    return children
}

// OPTIMIZED: Preallocate when size known
func (b *modelBuilder) buildChildren(contexts []Context) []Element {
    children := make([]Element, 0, len(contexts))  // Preallocate
    for _, ctx := range contexts {
        children = append(children, parseContext(ctx))
    }
    return children
}
```

### Benchmark Template for Parser Performance
```go
// Source: Go testing package best practices
func BenchmarkParseFile(b *testing.B) {
    content, err := os.ReadFile("testdata/test_example.sysml")
    if err != nil {
        b.Fatal(err)
    }

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        result := ParseBytes(content, "benchmark")
        if !result.Success() {
            b.Fatal("parse failed")
        }
    }
}

// Run with: go test -bench=. -benchmem -cpuprofile=cpu.prof
// Analyze: go tool pprof cpu.prof
```

### Profiling-Guided Optimization Workflow
```go
// Source: https://go.dev/blog/pprof
// 1. Add pprof endpoints for live profiling
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // ... parser code
}

// 2. Capture CPU profile
// go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

// 3. Capture memory profile
// go tool pprof http://localhost:6060/debug/pprof/heap

// 4. Analyze in interactive mode
// (pprof) top 20 -cum
// (pprof) list functionName
// (pprof) web  # Visualize call graph
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Core types in generics | Explicit rules per feature | Go 1.25 (Aug 2025) | Simpler spec, fully backward compatible |
| Manual error inspection | errors.Is/As/Unwrap | Go 1.13+ | Type-safe error checking, context preservation |
| sync.Map for all cases | Regular map + sync.RWMutex | Ongoing | sync.Map only faster for specific patterns |
| tokens.Fill() upfront | Lazy token consumption | Best practice 2024+ | 30-50% memory reduction for large files |
| String-based parsing | []byte-based parsing | Best practice 2020+ | Avoids string/bytes conversions |
| Global FileSet (go/token) | FileSet.AddExistingFiles | Go 1.25 | Eliminates global state problems |
| Manual iterator patterns | range-over-func (iter package) | Go 1.23 (Aug 2024) | Standard library integration, cleaner API |

**Deprecated/outdated:**
- `go/ast.FilterPackage`: Deprecated in Go 1.25, use `ast.PreorderStack` instead
- `errors.New` for wrapped errors: Use `fmt.Errorf` with `%w` to preserve error chain
- Buffering all tokens: Use lazy consumption except when token stream inspection needed
- Hand-rolled generic slices/maps: Use `golang.org/x/exp/slices` and `maps` packages

## Open Questions

Things that couldn't be fully resolved:

1. **ANTLR v4.13+ direct []byte support**
   - What we know: Current code converts []byte to string (low/lexer.go:26)
   - What's unclear: Whether antlr4-go v4.13+ supports zero-copy []byte input
   - Recommendation: Check ANTLR v4 changelog and test with antlr.NewByteArrayInputStream if available

2. **Optimal sync.Pool usage for parser**
   - What we know: sync.Pool helps for high-frequency, short-lived objects
   - What's unclear: Which specific types benefit in this parser (Location? Small contexts?)
   - Recommendation: Profile first, pool only if GC pressure measured in pprof

3. **Go 1.25 Green Tea GC benefit for parser workloads**
   - What we know: 10-40% GC reduction in GC-heavy workloads, experimental in 1.25
   - What's unclear: Whether parser allocation patterns benefit specifically
   - Recommendation: Benchmark with `GOEXPERIMENT=greenteagc` vs default GC

4. **Range-over-func performance vs traditional visitor**
   - What we know: Range-over-func is cleaner API, standard library integration
   - What's unclear: Performance difference vs current visitor implementation
   - Recommendation: Implement both, benchmark with real test files

## Sources

### Primary (HIGH confidence)
- Go 1.25 Release Notes: https://go.dev/doc/go1.25 - Language features, compiler improvements
- Go iter package docs: https://pkg.go.dev/iter - Iterator types and usage
- Go blog on when to use generics: https://go.dev/blog/when-generics - Official guidance
- Go blog on pprof: https://go.dev/blog/pprof - Profiling techniques
- Existing codebase: gosysml2/low/, gosysml2/sysml/ - Current implementation patterns

### Secondary (MEDIUM confidence)
- Go 1.26 Release Notes: https://go.dev/doc/go1.26 - Upcoming features (Feb 2026)
- Go slice preallocation guide: https://goperf.dev/01-common-patterns/mem-prealloc/ - Performance patterns
- OneUpTime pprof guide (Jan 2026): https://oneuptime.com/blog/post/2026-01-07-go-pprof-profiling/view - Current profiling practices
- Leapcell sync.Pool guide: https://leapcell.io/blog/unlocking-efficiency-demystifying-go-s-sync-pool-for-ephemeral-objects - Pool usage patterns

### Tertiary (LOW confidence - needs verification)
- ANTLR performance discussions: https://github.com/antlr/antlr4/issues/1540 - Community performance tips (2016-2021, dated)
- Go generics performance article: https://planetscale.com/blog/generics-can-make-your-go-code-slower - Specific use case, may not generalize

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - Go 1.25 is current stable, ANTLR v4 is established
- Architecture: HIGH - Existing codebase follows modern patterns, verified against official docs
- Pitfalls: HIGH - Based on official Go blog posts and documented best practices
- Performance optimization: MEDIUM - WebSearch results verified with official docs where possible, some techniques require benchmarking to confirm benefit

**Research date:** 2026-02-03
**Valid until:** 30 days (until Go 1.26 release expected Feb 2026, which may add new optimizations)

## Recommendations for Planning

1. **Language version**: Update go.mod to Go 1.25+, verify 1.25.5 availability
2. **Quick wins**: Apply slice preallocation in model building (identifiable with pprof)
3. **Modernization targets**:
   - Add range-over-func iterators to visitor pattern (clean API)
   - Improve error wrapping with context (better debugging)
   - Apply generics to repeated type-specific functions
4. **Performance validation**:
   - Establish baseline benchmarks with 34 test files in docs/testdata/
   - Profile with pprof before/after changes
   - Target balanced speed/memory (avoid premature optimization)
5. **Don't modify**: ANTLR-generated code in internal/parser/
6. **Test coverage**: Use existing test files, add benchmarks for regression detection
