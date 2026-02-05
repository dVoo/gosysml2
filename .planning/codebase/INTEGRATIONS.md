# External Integrations

**Analysis Date:** 2026-02-05

## Overview

This is a **standalone parsing library** with no external service dependencies. It operates entirely within-process without requiring network access, databases, or external APIs.

## APIs & External Services

**None** - The library is self-contained.

## Data Storage

**Databases:**
- None - In-memory model representation only

**File Storage:**
- Local filesystem only
  - Read: `ParseFile()`, `ParseDirectory()` in `gosysml2/sysml/parse.go`
  - Input: `.sysml` files (SysML v2 textual syntax)

**Caching:**
- None - No external caching layer
- In-memory element index built during parse (`Model.BuildIndex()`)

## Authentication & Identity

**Auth Provider:**
- None - No authentication required for library operation

## Monitoring & Observability

**Error Tracking:**
- Custom error collection via `ErrorCollector` (`gosysml2/low/errors.go`)
- Structured error types: `SyntaxError`, `ParseError` (`gosysml2/sysml/errors.go`)
- No external error tracking services

**Logs:**
- No logging framework - returns errors to caller
- Example programs use `fmt.Println()` for demonstration only

## CI/CD & Deployment

**Hosting:**
- GitHub (inferred from `.git/` structure)

**CI Pipeline:**
- None detected - No GitHub Actions, GitLab CI, or other CI configuration found

**Package Distribution:**
- Go module: `github.com/dVoo/gosysml2`
- Published via Go module proxy (standard Go mechanism)

## Environment Configuration

**Required env vars:**
- None

**Optional env vars:**
- `GOPROXY` - Standard Go module proxy configuration
- `GOPATH` - Standard Go workspace (if not using modules)

**Secrets location:**
- Not applicable - No secrets required

## Webhooks & Callbacks

**Incoming:**
- None

**Outgoing:**
- None

## Integration Points

**Library Consumers:**
- Applications import `github.com/dVoo/gosysml2/sysml` for high-level API
- Applications import `github.com/dVoo/gosysml2/low` for low-level API

**Example Integration:**
```go
import "github.com/dVoo/gosysml2/sysml"

result := sysml.ParseFile("model.sysml")
if !result.Success() {
    // Handle errors locally
}
```

**CLI Tools:** (referenced in README, not in `gosysml2/` module)
- `cmd/verify-completeness` - Model analysis
- `cmd/test-attrs` - Requirement attribute testing
- `cmd/verify-parser` - Parser verification

---

*Integration audit: 2026-02-05*
