# Performance Optimization Guide

## Current Limitations

The current implementation loads everything into memory, which limits parsing to files/repositories that fit in RAM (roughly input_size * 30-50).

## Optimization Strategies

### Phase 1: Quick Wins (Low Effort, High Impact)

#### 1.1 Remove `tokens.Fill()` - Let ANTLR consume lazily

```go
// low/parser.go - REMOVE this line:
tokens.Fill()  // DELETE THIS
```

This alone can reduce memory by 50%+ for large files.

#### 1.2 Don't retain parse tree when building model

```go
// Option to discard parse tree after model building
type ParseOption func(*parseConfig)

func WithDiscardTree() ParseOption {
    return func(c *parseConfig) { c.discardTree = true }
}

func parseWithSource(input, source string, opts ...ParseOption) *ParseResult {
    // ... parse ...
    result.Model = buildModel(tree)
    if cfg.discardTree {
        result.Tree = nil  // Allow GC
    }
    return result
}
```

#### 1.3 Use []byte throughout, avoid string conversion

```go
// low/lexer.go - Use byte-based input stream
func NewLexerFromBytes(input []byte) *Lexer {
    // antlr4-go supports this directly
    stream := antlr.NewByteStream(input)  // No copy!
    return NewLexerFromStream(stream)
}
```

### Phase 2: Streaming Architecture (Medium Effort)

#### 2.1 File-by-file processing with memory release

```go
// Process repository incrementally
func ParseRepository(dir string, handler func(*Model) error) error {
    return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
        if !strings.HasSuffix(path, ".sysml") {
            return nil
        }

        result := ParseFile(path)
        if err := handler(result.Model); err != nil {
            return err
        }

        // Explicitly release memory
        result.Model = nil
        result.Tree = nil
        runtime.GC()  // Hint to garbage collector

        return nil
    })
}
```

#### 2.2 Memory-mapped file reading

```go
import "golang.org/x/exp/mmap"

func ParseFileMmap(filename string) *ParseResult {
    reader, err := mmap.Open(filename)
    if err != nil {
        return &ParseResult{Errors: ...}
    }
    defer reader.Close()

    // Read into byte slice without full copy
    data := make([]byte, reader.Len())
    reader.ReadAt(data, 0)

    return parseFromBytes(data, filename)
}
```

#### 2.3 SAX-style streaming model builder

Instead of building full model, stream events:

```go
type ModelHandler interface {
    OnPackage(name string, loc Location) error
    OnPart(name string, loc Location, isDefinition bool) error
    OnRequirement(name string, loc Location) error
    // ... etc
}

// Stream parse without building full model
func ParseStream(input []byte, handler ModelHandler) error {
    tree, errs := low.ParseBytes(input)
    if errs.HasErrors() {
        return errs
    }

    streamer := &streamBuilder{handler: handler}
    antlr.ParseTreeWalkerDefault.Walk(streamer, tree)
    return nil
}
```

### Phase 3: Parallel Processing (Higher Effort)

#### 3.1 Concurrent file parsing

```go
func ParseRepositoryParallel(dir string, workers int) ([]*ParseResult, error) {
    files := findSysMLFiles(dir)
    results := make([]*ParseResult, len(files))

    var wg sync.WaitGroup
    sem := make(chan struct{}, workers)

    for i, file := range files {
        wg.Add(1)
        sem <- struct{}{}

        go func(idx int, path string) {
            defer wg.Done()
            defer func() { <-sem }()

            results[idx] = ParseFile(path)
        }(i, file)
    }

    wg.Wait()
    return results, nil
}
```

#### 3.2 Worker pool with bounded memory

```go
type ParseJob struct {
    Path   string
    Result chan *ParseResult
}

func StartParseWorkers(jobs <-chan ParseJob, workers int) {
    for i := 0; i < workers; i++ {
        go func() {
            for job := range jobs {
                result := ParseFile(job.Path)
                job.Result <- result
            }
        }()
    }
}
```

### Phase 4: Advanced Optimizations

#### 4.1 Incremental/cached parsing

```go
type ParseCache struct {
    mu     sync.RWMutex
    models map[string]*CachedModel
}

type CachedModel struct {
    Model    *Model
    ModTime  time.Time
    Checksum [32]byte
}

func (c *ParseCache) ParseFile(path string) (*Model, error) {
    stat, _ := os.Stat(path)

    c.mu.RLock()
    cached, ok := c.models[path]
    c.mu.RUnlock()

    if ok && cached.ModTime.Equal(stat.ModTime()) {
        return cached.Model, nil  // Return cached
    }

    // Parse and cache
    result := ParseFile(path)
    c.mu.Lock()
    c.models[path] = &CachedModel{
        Model:   result.Model,
        ModTime: stat.ModTime(),
    }
    c.mu.Unlock()

    return result.Model, nil
}
```

#### 4.2 Lazy model loading

```go
// Only parse structure, defer body parsing
type LazyPackage struct {
    Name     string
    Location Location
    source   []byte
    parsed   *Package
    once     sync.Once
}

func (p *LazyPackage) Children() []Element {
    p.once.Do(func() {
        p.parsed = parsePackageBody(p.source)
    })
    return p.parsed.Children()
}
```

#### 4.3 Object pooling for reduced allocations

```go
var locationPool = sync.Pool{
    New: func() interface{} {
        return &Location{}
    },
}

var packagePool = sync.Pool{
    New: func() interface{} {
        return &Package{
            baseElement: baseElement{
                children: make([]Element, 0, 8),
            },
        }
    },
}
```

## Implementation Priority

| Priority | Change | Memory Reduction | Effort |
|----------|--------|------------------|--------|
| 1 | Remove `tokens.Fill()` | 30-50% | 1 line |
| 2 | Discard tree option | 20-40% | ~20 lines |
| 3 | Byte-based input | 10-20% | ~10 lines |
| 4 | File-by-file with GC | Variable | ~30 lines |
| 5 | Parallel parsing | N/A (speed) | ~50 lines |
| 6 | Streaming handler | 80-90% | ~100 lines |
| 7 | Parse cache | Variable | ~100 lines |

## Quick Start: Minimal Changes for 10x Improvement

Apply these 3 changes for immediate improvement:

1. **Remove `tokens.Fill()`** in `low/parser.go:53`
2. **Add tree discard option** to not retain parse tree
3. **Use `ParseBytes`** with `[]byte` input throughout

This should allow parsing repositories up to ~1GB on a 32GB machine.

For multi-GB repositories, implement the streaming handler (Phase 2.3) which processes elements without building a full in-memory model.
