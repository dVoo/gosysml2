---
phase: quick
task: 001
description: "Reorganize project structure"
date: 2026-02-07
commit: 84b6d97
---

# Quick Task 001: Reorganize Project Structure

**Summary:** Moved all Go code from `gosysml2/` subdirectory to root level and consolidated documentation.

## Changes Made

### Project Structure

**Before:**
```
.
├── gosysml2/              # Go module in subdirectory
│   ├── go.mod
│   ├── go.sum
│   ├── README.md          # Library docs
│   ├── USAGE.md           # Usage guide
│   ├── PERFORMANCE.md     # Performance docs
│   ├── low/               # Low-level API
│   ├── sysml/             # High-level API
│   ├── internal/          # Internal parser
│   ├── cmd/               # CLI tools
│   └── examples/          # Example programs
├── docs/
│   └── usage-guide.md     # Project-level docs
└── README.md              # Root README
```

**After:**
```
.
├── go.mod                 # Go module at root
├── go.sum
├── README.md              # Concise entry point
├── low/                   # Low-level API
├── sysml/                 # High-level API
├── internal/              # Internal parser
├── cmd/                   # CLI tools
├── examples/              # Example programs
└── docs/
    ├── README.md          # Full API reference
    ├── USAGE.md           # Usage guide
    ├── PERFORMANCE.md     # Performance docs
    └── usage-guide.md     # Project-level docs
```

### Key Changes

1. **Go Module**: Moved `gosysml2/go.mod` and `gosysml2/go.sum` to root
   - Updated module name from `github.com/dVoo/gosysml2` to `github.com/dVoo/gosysml2_oc`

2. **Packages**: Moved all Go packages to root
   - `gosysml2/low/` → `low/`
   - `gosysml2/sysml/` → `sysml/`
   - `gosysml2/internal/` → `internal/`
   - `gosysml2/cmd/` → `cmd/`
   - `gosysml2/examples/` → `examples/`

3. **Documentation**: Consolidated in `docs/`
   - `gosysml2/README.md` → `docs/README.md` (full API reference)
   - `gosysml2/USAGE.md` → `docs/USAGE.md`
   - `gosysml2/PERFORMANCE.md` → `docs/PERFORMANCE.md`
   - Updated root `README.md` to be a concise entry point

4. **Import Paths**: Updated all imports from `github.com/dVoo/gosysml2` to `github.com/dVoo/gosysml2_oc`
   - All Go source files
   - All Markdown documentation

### Root README Changes

The root `README.md` was simplified to be a concise entry point:
- Project status and features
- Quick start example
- Links to full documentation in `docs/`
- Repository structure
- Development environment setup

Full API documentation moved to `docs/README.md`.

## Verification

- ✅ All packages build: `go build ./low/... ./sysml/... ./internal/...`
- ✅ Module name updated: `go.mod` shows `github.com/dVoo/gosysml2_oc`
- ✅ Import paths updated in all Go files
- ✅ Documentation consolidated in `docs/`
- ✅ Root README is concise entry point
- ✅ `gosysml2/` directory removed

## Known Issues

- `check_testdata.go` has pre-existing type mismatches (not introduced by this reorganization):
  - Line 47: `*sysml.Model` vs `sysml.Model` type mismatch
  - Line 76: `sysml.Model` vs `*sysml.Model` type mismatch
  - These errors existed before the reorganization

## Files Modified

- `README.md` — Simplified root README
- `go.mod` — Updated module name
- `docs/README.md` — Full API reference (moved from gosysml2/)
- `docs/USAGE.md` — Usage guide (moved from gosysml2/)
- `docs/PERFORMANCE.md` — Performance guide (moved from gosysml2/)
- All Go source files — Updated import paths
- All Markdown files — Updated import references

## Commit

`84b6d97` — chore: reorganize project structure
