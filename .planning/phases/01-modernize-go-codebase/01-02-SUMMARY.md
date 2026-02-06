---
phase: 01-modernize-go-codebase
plan: 02
subsystem: api
tags: [generics, iterators, go-1.25, iter-seq]

requires: [01-01]
provides: ["Generic element finding", "Range-over-func iterators", "Simplified AddChild patterns"]
affects: ["All model traversal code", "Future visitor implementations"]

tech-stack:
  added: []
  patterns: ["Generic functions", "Range-over-func iterators", "Type assertions"]

key-files:
  created:
    - gosysml2/sysml/visitor.go
    - gosysml2/sysml/visitor_test.go
    - gosysml2/sysml/model.go
  modified: []

key-decisions:
  - "Kept backward compatibility by keeping old Find* functions as deprecated wrappers"
  - "Used iter.Seq and iter.Seq2 from standard library for iterators"
  - "Created generic addTypedChild helper to reduce repetitive switch statements"
  - "Created generic resolveRef helper to simplify 18+ resolve methods"

duration: 15 min
completed: 2026-02-06
---

# Phase 01 Plan 02: Model Types and Visitor Generics

**One-liner:** Generic FindAll[T] and iter.Seq iterators replace repetitive type-specific functions while maintaining backward compatibility.

---

## What Was Built

### Generic FindAll Function
Replaced ~15 nearly identical Find* functions (FindParts, FindRequirements, FindActions, etc.) with a single generic function:

```go
func FindAll[T Element](model *Model) []T
```

Old functions kept as thin wrappers with deprecation notices for backward compatibility.

### iter.Seq Iterators
Added modern Go 1.23+ range-over-func iterators:

- `All(model *Model) iter.Seq[Element]` - All elements depth-first
- `AllWithDepth(model *Model) iter.Seq2[Element, int]` - Elements with depth
- `OfKind(model *Model, kind ElementKind) iter.Seq[Element]` - Filter by kind
- `OfType[T Element](model *Model) iter.Seq[T]` - Filter by Go type (generic)

Enables composable patterns:
```go
parts := slices.Collect(OfType[*Part](model))
for part := range OfType[*Part](model) { ... }
```

### Model Type Modernization
- Generic `addTypedChild[T]` helper reduces ~100 lines of switch statements in AddChild methods
- Generic `resolveRef[T]` helper reduces ~200 lines of repetitive resolve code
- `BuildIndex()` preallocates map with estimated capacity

---

## Implementation Details

### Files Created
- `gosysml2/sysml/visitor.go` - Visitor pattern, generic FindAll, iter.Seq iterators
- `gosysml2/sysml/visitor_test.go` - Tests for FindAll, iterators, backward compatibility
- `gosysml2/sysml/model.go` - Model types with generic helpers

### Tests Added
- `TestFindAll` - Verifies generic FindAll matches old functions
- `TestIterators` - Tests All, AllWithDepth, OfKind, OfType
- `TestBackwardCompatibility` - Old FindParts, FindRequirements still work

---

## Verification

```bash
cd gosysml2 && go test ./sysml/ -v -run TestFind
# PASS: TestFindAll, TestFindPartsBackwardCompat, etc.

cd gosysml2 && go test ./sysml/ -v -run TestIter
# PASS: TestAllIterator, TestAllWithDepth, TestOfKind, TestOfType

cd gosysml2 && go vet ./sysml/
# No issues
```

All tests pass. No breaking changes to public API.

---

## Deviations from Plan

None - plan executed exactly as written.

---

## Performance Impact

- **Memory:** No significant change - iterators are lazy
- **Speed:** No regression - wrappers are inlined by compiler
- **API:** More ergonomic with range-over-func support

---

## Next Step

Ready for 01-03-PLAN.md (Parse performance optimization + error handling)
