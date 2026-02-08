# API Issues Catalog

**Date:** 2026-02-08  
**Source:** Quick Task 006 - Diagram Relations Investigation  
**Status:** Documented for v0.2 planning

---

## Critical API Issues Identified

### Issue 1: Model.ResolveReferences Missing Dependency Case

**Location:** `sysml/model.go:ResolveReferences()` (line 1874)

**Problem:** The `ResolveReferences()` function has a switch statement that handles reference resolution for many element types, but `*Dependency` is missing.

**Impact:** Dependencies are parsed (client/supplier extracted) but references are never resolved, so `Dependency.Client` and `Dependency.Supplier` slices remain empty.

**Code:**
```go
// In ResolveReferences() - MISSING:
case *Dependency:
    m.resolveDependencyRefs(e)  // This function doesn't exist either!
```

**Suggested Fix:**
1. Create `resolveDependencyRefs()` function
2. Add `case *Dependency:` to the switch

**Priority:** HIGH - Data loss (parsed but unusable)

---

### Issue 2: Connection Created Without Ends

**Location:** `sysml/parse.go:EnterBindingConnectorAsUsage()` (line 2428)

**Problem:** Binding connectors are created as Connection elements, but the parser never extracts the connector ends from the context.

**Impact:** Connection exists in model but `Connection.Ends` is always empty. The binding relationship information is lost.

**Current Code:**
```go
func (b *modelBuilder) EnterBindingConnectorAsUsage(ctx *parser.BindingConnectorAsUsageContext) {
    loc := locationFromContext(ctx)
    conn := NewConnection("", loc, false)
    conn.parent = b.getCurrentParent()
    // ... adds to package and stack ...
    // NEVER extracts ends from ctx!
}
```

**Suggested Fix:**
Extract ends from `ctx.BindingConnectorEnd()` or similar context methods and call `conn.AddUnresolvedEnd()`.

**Priority:** HIGH - Incomplete feature

---

### Issue 3: Transition Created Without Source/Target

**Location:** `sysml/parse.go:EnterSuccessionAsUsage()` (line 2458)

**Problem:** Succession (state transition) elements are created but source/target references are never extracted.

**Impact:** Transition exists but `Transition.Source` and `Transition.Target` are never populated (unresolved strings remain empty).

**Current Code:**
```go
func (b *modelBuilder) EnterSuccessionAsUsage(ctx *parser.SuccessionAsUsageContext) {
    loc := locationFromContext(ctx)
    trans := NewTransition("", loc)
    // ... adds to package and stack ...
    // NEVER extracts source/target from ctx!
}
```

**Suggested Fix:**
Extract source/target from context and call `trans.SetUnresolvedSource()` / `SetUnresolvedTarget()`.

**Priority:** HIGH - Incomplete feature

---

### Issue 4: Flow Created Without FlowEnds

**Location:** `sysml/parse.go:EnterFlowDefinition()` and `EnterFlowUsage()`

**Problem:** Flow elements are created but source/target FlowEnd elements are never created or populated.

**Impact:** Flow exists in model but has no endpoints, making it unusable for diagram generation.

**Current Code:**
```go
func (b *modelBuilder) EnterFlowDefinition(ctx *parser.FlowDefinitionContext) {
    // ... extract name ...
    flow := NewFlow(name, loc, true)
    // ... adds to package ...
    // NEVER creates FlowEnd elements!
}
```

**Suggested Fix:**
Create FlowEnd elements from `ctx.FlowEndMember()` and set `FlowEnd.Reference`.

**Priority:** HIGH - Incomplete feature

---

### Issue 5: Missing Diagram Generation API

**Location:** N/A (entire feature missing)

**Problem:** No diagram generation API exists. The `analysis` package only contains D2 shape mapping documentation.

**Impact:** No way to convert parsed models to visual diagrams. Relations cannot be visualized.

**Current State:**
- `analysis/doc.go` - Only documentation about D2 shapes
- No `GenerateD2()` or similar functions
- No diagram output formatters

**Suggested Improvement:**
Create new package or extend analysis with:
```go
// Proposed API
func GenerateD2(model *Model, options D2Options) (string, error)
type D2Options struct {
    IncludeRelations bool
    PackageFilter []string
    // ...
}
```

**Priority:** MEDIUM - New feature needed

---

### Issue 6: Missing Watch Mode API

**Location:** N/A (entire feature missing)

**Problem:** No file watching functionality exists. No way to auto-regenerate on file changes.

**Impact:** Users must manually re-run parser; no real-time feedback during editing.

**Suggested Improvement:**
Create `cmd/watch` or `sysml/watch` package with:
```go
// Proposed API
type Watcher struct {
    RootPath string
    OnChange func(files []string)
}
func (w *Watcher) Start() error
func (w *Watcher) Stop()
```

**Priority:** LOW - Convenience feature

---

### Issue 7: Inconsistent Reference Extraction Pattern

**Location:** Multiple files

**Problem:** Different relation types use inconsistent patterns for reference extraction.

**Inconsistencies:**
- SuccessionFlow: Extracts in `Enter` handler (GOOD)
- Dependency: Extracts in `Enter` handler (GOOD)
- Connection: Doesn't extract (BAD)
- Transition: Doesn't extract (BAD)
- Flow: Doesn't extract (BAD)

**Impact:** Hard to maintain; inconsistent behavior across relation types.

**Suggested Fix:**
Standardize pattern:
1. Extract raw strings in Enter handler
2. Store in `unresolvedXxx` fields
3. Resolve in `resolveXxxRefs()` function
4. All relations follow same pattern

**Priority:** MEDIUM - Code quality

---

## API Usability Notes

### Positive Findings

1. **Model.Connections()** - Clean API for accessing connections
2. **Model.AllFlows()** - Good iterator for flows
3. **Reference resolution pattern** - Consistent `Ref[T]` generic type
4. **Element interface** - Clean visitor pattern implementation

### Areas for Improvement

1. **Error handling in parsers** - Silent failures when extraction missing
2. **Debuggability** - No way to see what was extracted vs what was expected
3. **Test coverage** - Relation extraction not well tested
4. **Documentation** - No docs on what relations are fully supported

---

## Recommendations for v0.2 API Design

1. **Add validation layer** - Check that required fields are populated after parsing
2. **Add diagnostic API** - Query what references couldn't be resolved
3. **Complete extraction** - Fix all parsers to extract relation endpoints
4. **Add diagram API** - Generate D2 or other formats from model
5. **Document limitations** - Clear docs on what grammar elements are supported

---

*Related: diagram-relations-investigation.md*
