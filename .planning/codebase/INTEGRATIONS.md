# External Integrations

**Analysis Date:** 2026-02-05

## APIs & External Services

**None.**

This is a self-contained library and CLI tool suite. No external API calls are made.

## Data Storage

**Databases:**
- Not applicable - No database integration

**File Storage:**
- Local filesystem only
- Input: Reads `.sysml` files from local paths
- Output: Console/stdout only

**Caching:**
- None

## Authentication & Identity

**Auth Provider:**
- Not applicable - No authentication required

## Monitoring & Observability

**Error Tracking:**
- None

**Logs:**
- Console output only
- Tools write analysis results to stdout

## CI/CD & Deployment

**Hosting:**
- Not applicable - Library and CLI tools

**CI Pipeline:**
- None detected
- No `.github/workflows/` or similar CI configuration

**Distribution:**
- Go module: `github.com/dVoo/gosysml2`
- Can be imported as Go library
- CLI tools built from source with `go build`

## Environment Configuration

**Required env vars:**
- None

**Secrets location:**
- Not applicable - No secrets required

## Webhooks & Callbacks

**Incoming:**
- None

**Outgoing:**
- None

## Network Dependencies

**None.**

The codebase is completely offline-capable:
- Parses local SysML files only
- No HTTP client usage
- No network-dependent imports

## File Format Dependencies

**Input Formats:**
- SysML v2 text files (`.sysml` extension)
- Parsed using ANTLR4-generated parsers

**Grammar Sources:**
- Based on OMG SysML v2 specification
- BNF specifications in `docs/bnf/`
- ANTLR4 grammars in `code/*.g4`

---

*Integration audit: 2026-02-05*
