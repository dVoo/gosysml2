# External Integrations

**Analysis Date:** 2026-01-30

## APIs & External Services

**None detected.**

The codebase does not integrate with external APIs or SaaS platforms. This is a language parser and grammar specification repository focused on SysML v2 and KerML syntax.

## Data Storage

**Databases:**
- None - This is a stateless parser library with no persistent storage

**File Storage:**
- Local filesystem only - Codebase reads `.sysml` files from local file system
  - `sysml.ParseFile()` - Reads SysML files
  - `sysml.ParseDirectory()` - Reads all `.sysml` files from directory
  - `sysml.ParseDirectoryParallel()` - Parallel file reading
  - `sysml.ParseDirectoryStream()` - Streaming file processing
  - Located in `gosysml2/sysml/parse.go`

**Caching:**
- None - No caching layer implemented
- Parse results are ephemeral and returned directly to caller

## Authentication & Identity

**Auth Provider:**
- None - This is an open-source language parser with no authentication requirements

## Monitoring & Observability

**Error Tracking:**
- None - No error tracking service integration

**Logs:**
- Console/stdout logging via Go standard library
- Error collection mechanism in `gosysml2/sysml/errors.go`
- Parse errors include source location (line, column) for debugging
- Located in `gosysml2/sysml/errors.go`

## CI/CD & Deployment

**Hosting:**
- GitHub (public repository, no deployment infrastructure)
- Module hosted at `github.com/dVoo/gosysml2`

**CI Pipeline:**
- Not detected - No GitHub Actions, GitLab CI, or other CI service configuration
- Manual testing via `go test` commands

## Environment Configuration

**Required env vars:**
- None - No environment variables required

**Secrets location:**
- None - No secrets management needed

## Webhooks & Callbacks

**Incoming:**
- None

**Outgoing:**
- None

## Test Data Integration

**Test Files:**
- SysML test files in `testdata/` and `docs/testdata/` (46 `.sysml` files)
- Test coverage files in `gosysml2/sysml/*.go` (integration tests, unit tests)
- Test files use standard Go `testing` package with no external test framework

**Grammar Reference Data:**
- BNF specifications in `docs/bnf/`:
  - `KerML-textual-bnf.kebnf` - Reference grammar
  - `SysML-textual-bnf.kebnf` - Reference grammar
  - `SysML-graphical-bnf.kgbnf` - Graphical notation grammar
- Rendered HTML documentation from BNF files
- SVG diagrams for graphical notation (284 files in `docs/bnf/images/`)

## Standard Library Dependencies

**Go Standard Library:**
- `io`, `os`, `path/filepath` - File I/O operations
- `sync`, `runtime` - Concurrency and threading
- `strings` - String manipulation

**No external Go packages required except:**
- `github.com/antlr4-go/antlr/v4` - ANTLR4 Go runtime (single external dependency)

## Build Tool Integration

**Code Generation:**
- ANTLR4 generates Go parser files from grammar specifications
- No code generation from external APIs
- Grammar files are source-controlled; generated code is in `code/parser/` and `gosysml2/internal/parser/`

---

*Integration audit: 2026-01-30*
