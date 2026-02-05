# External Integrations

**Analysis Date:** 2026-02-05

## APIs & External Services

**None detected** - This is a standalone parsing library with no external API integrations.

The library provides interfaces for consuming SysML v2 models but does not call external services.

## Data Storage

**Databases:**
- None - Library is in-memory only

**File Storage:**
- Local filesystem only
  - Input: Reads `.sysml` files from filesystem via `sysml.ParseFile()`, `sysml.ParseDirectory()`, etc.
  - Output: No file writing. Library returns in-memory `*sysml.Model` objects

**Caching:**
- None - All results are computed on-demand

## Authentication & Identity

**Auth Provider:**
- Not applicable - No authentication layer

## Monitoring & Observability

**Error Tracking:**
- None - Error handling is local to caller

**Logs:**
- No logging framework integrated
- Library is silent by default
- Callers receive errors via `ParseResult.Errors` field containing structured `ParseError` objects

## CI/CD & Deployment

**Hosting:**
- Public Go module hosted at `github.com/dVoo/gosysml2`
- Distributed via standard Go module system

**CI Pipeline:**
- Not detected in repository
- Test execution: Standard `go test ./...` pattern

## Environment Configuration

**Required env vars:**
- None - Library has no environment configuration

**Secrets location:**
- Not applicable - No secrets needed

## Webhooks & Callbacks

**Incoming:**
- None

**Outgoing:**
- None - Library is request-response only (synchronous parsing)

## ANTLR Integration Details

**Grammar Source:**
- Reference grammar files located in parent repository at `/home/daniel/projects/claudecode/code/`:
  - `SysMLv2Lexer.g4` - Lexer specification for both KerML and SysML
  - `KerMLParser.g4` - KerML parser rules
  - `SysMLv2Parser.g4` - Full SysML parser rules

**Generated Parser:**
- ANTLR code generation target: Go runtime
- Generated into: `github.com/dVoo/gosysml2/internal/parser`
- Generator: ANTLR4 with `-Dlanguage=Go` option

**Parser Generation Command (from parent repo):**
```bash
cd code
antlr -Dlanguage=Go -o parser SysMLv2Lexer.g4 SysMLv2Parser.g4
```

## Streaming & Parallel Processing

**Streaming Input:**
- `sysml.ParseReader(io.Reader, source)` - Supports reading from any `io.Reader` (HTTP responses, pipes, etc.)
- `sysml.ParseDirectoryStream(dir, handler)` - Streaming directory parsing with callback

**Parallel Processing:**
- `sysml.ParseDirectoryParallel(dir, workers)` - Multi-worker directory parsing
- Uses Go goroutines with semaphore-based concurrency control
- Worker count: `runtime.NumCPU()` by default or custom value

## Performance Integration

No external performance monitoring, but library provides:
- Memory optimization: `WithDiscardTree()` option reduces memory ~30%
- Capacity guidelines for large repositories up to 1+ GB
- See `PERFORMANCE.md` for detailed capacity analysis

---

*Integration audit: 2026-02-05*
