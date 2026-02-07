---
phase: quick
task: 002
description: "Rename library from gosysml2_oc to gosysml2"
date: 2026-02-07
commit: b4a267c
---

# Quick Task 002: Rename Library to gosysml2

**Summary:** Renamed the library from `gosysml2_oc` to the cleaner `gosysml2` name.

## Changes Made

### Module Name
- **go.mod**: Changed from `github.com/dVoo/gosysml2_oc` to `github.com/dVoo/gosysml2`

### Go Source Files (15 files)
All import paths updated:
- `low/lexer.go`
- `low/parser.go`
- `low/doc.go`
- `sysml/errors.go`
- `sysml/bench_test.go`
- `sysml/doc.go`
- `sysml/parse.go`
- `cmd/check_validation.go`
- `examples/main.go`
- `examples/basic/main.go`
- `examples/parallel/main.go`
- `examples/requirements/main.go`
- `examples/validation/main.go`
- `examples/visitor/main.go`
- `check_testdata.go`

### Documentation (7 files)
- `README.md` — Updated go get command and import examples
- `docs/README.md` — Updated all import paths
- `docs/USAGE.md` — Updated import paths
- `docs/usage-guide.md` — Updated references
- `.planning/codebase/STACK.md` — Updated module name
- `.planning/codebase/INTEGRATIONS.md` — Updated references
- `.planning/codebase/CONVENTIONS.md` — Updated import example

## Verification

- ✅ go.mod shows `module github.com/dVoo/gosysml2`
- ✅ All Go files use `github.com/dVoo/gosysml2` imports
- ✅ All documentation references `gosysml2`
- ✅ Core packages build: `go build ./low/... ./sysml/... ./internal/...`
- ✅ No remaining references to `gosysml2_oc`

## Files Modified

23 files changed, 179 insertions(+), 32 deletions(-)

## Commit

`b4a267c` — chore: rename library from gosysml2_oc to gosysml2
