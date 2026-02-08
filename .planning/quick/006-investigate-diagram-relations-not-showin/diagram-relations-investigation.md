# Diagram Relations Investigation Report

**Date:** 2026-02-08  
**Investigation:** Quick Task 006 - Why diagram relations are not showing in watch mode  
**Status:** COMPLETE - Root cause identified

---

## Executive Summary

**The investigation revealed that "watch mode" and "diagram rendering" features do not exist in the current codebase.** The plan assumed these features were implemented but not working correctly, when in fact they were never built. The parser does handle relation elements, but several critical gaps prevent relations from being fully usable.

---

## Issue Summary

- **Expected:** Watch mode should generate diagrams showing relations between elements
- **Actual:** Watch mode does not exist; diagram generation does not exist
- **Root Cause Layer:** Missing features (parser extraction gaps + no diagram generation layer)

---

## Investigation Methodology

1. Searched for `watch` command implementation → **Not found**
2. Searched for diagram/D2 generation code → **Not found** (only shape mapping docs)
3. Analyzed parser relation handlers in `sysml/parse.go`
4. Analyzed reference resolution in `sysml/model.go`
5. Identified specific gaps in each relation type

---

## Findings

### 1. Missing Features (Architecture Layer)

| Feature | Status | Evidence |
|---------|--------|----------|
| **Watch Mode** | ❌ **MISSING** | No `cmd/watch/` directory, no file watching code |
| **Diagram Generation** | ❌ **MISSING** | Only `analysis/doc.go` with D2 shape mapping documentation |
| **D2 Output** | ❌ **MISSING** | No code generates D2 syntax from models |

**analysis/doc.go** contains:
- Documentation about valid D2 shapes
- SysML-to-D2 shape mapping recommendations
- **No actual generation code**

### 2. Parser Relation Handling

The parser HAS handlers for relation elements, but implementation is incomplete:

#### ✅ Working Correctly

**SuccessionFlow** (`EnterSuccessionFlowUsage`)
- **File:** `sysml/parse.go:2295-2351`
- **Extracts:** Source and target from flow declaration
- **Resolves:** Yes, via `resolveSuccessionFlowRefs`
- **Status:** ✅ Functional

**Dependency** (`EnterDependency`)
- **File:** `sysml/parse.go:2133-2173`
- **Extracts:** Client and supplier qualified names
- **Resolves:** ❌ **NO** - Missing from `ResolveReferences` switch
- **Status:** ⚠️ Parsed but never resolved

#### ❌ Broken/Missing Extraction

**Connection** (`EnterBindingConnectorAsUsage`)
- **File:** `sysml/parse.go:2426-2454`
- **Creates:** Connection element with `NewConnection("", loc, false)`
- **Extracts:** ❌ **NOTHING** - Ends never populated
- **Result:** `unresolvedEnds` slice always empty
- **Resolution:** N/A (nothing to resolve)
- **Status:** ❌ **NON-FUNCTIONAL**

**Code evidence:**
```go
func (b *modelBuilder) EnterBindingConnectorAsUsage(ctx *parser.BindingConnectorAsUsageContext) {
    loc := locationFromContext(ctx)
    conn := NewConnection("", loc, false)  // Created empty
    // ...
    // NEVER extracts end references from ctx!
    // ...
    b.elementStack = append(b.elementStack, conn)
}
```

**Transition** (`EnterSuccessionAsUsage`)
- **File:** `sysml/parse.go:2456-2481`
- **Creates:** Transition element with `NewTransition("", loc)`
- **Extracts:** ❌ **NOTHING** - Source/target never populated
- **Result:** `unresolvedSource`/`unresolvedTarget` always empty
- **Resolution:** N/A (nothing to resolve)
- **Status:** ❌ **NON-FUNCTIONAL**

**Flow** (`EnterFlowDefinition`, `EnterFlowUsage`)
- **Files:** `sysml/parse.go:2231-2291`
- **Creates:** Flow element
- **Extracts:** ❌ **NOTHING** - Source/target FlowEnd elements never created
- **Result:** Flow exists but has no connections
- **Status:** ❌ **NON-FUNCTIONAL**

### 3. Reference Resolution Gap

In `sysml/model.go:ResolveReferences()` (line 1874), the switch statement handles:
- ✅ SuccessionFlow → `resolveSuccessionFlowRefs`
- ❌ **Dependency** → **NOT HANDLED** (missing case)
- ✅ Transition → `resolveTransitionRefs` (but never has data to resolve)
- ✅ Connection → `resolveConnectionRefs` (but never has data to resolve)
- ❌ Flow → **NOT HANDLED** (no resolve function exists)

**Missing dependency resolution:**
```go
// In model.go:1876-1925 - Dependency is MISSING from this switch
case *Requirement: m.resolveRequirementRefs(e)
case *Verification: m.resolveVerificationRefs(e)
// ... many others ...
case *SuccessionFlow: m.resolveSuccessionFlowRefs(e)
// NO case *Dependency!
```

---

## Root Cause Breakdown

```
Expected Flow:
  SysML File → Parse → Model (with relations) → Generate D2 → Watch Mode Display

Actual State:
  SysML File → Parse → Model (incomplete relations) → ❌ NO DIAGRAM CODE → ❌ NO WATCH MODE
```

### Parser Layer Issues
1. **BindingConnectorAsUsage**: Creates Connection but never extracts ends
2. **SuccessionAsUsage**: Creates Transition but never extracts source/target
3. **FlowDefinition/FlowUsage**: Creates Flow but never extracts source/target FlowEnds
4. **Dependency**: Extracts client/supplier but never resolved (missing in switch)

### Model Layer Issues
1. No `resolveDependencyRefs` function exists
2. `ResolveReferences()` switch missing `case *Dependency:`

### Application Layer Issues
1. No watch mode command exists
2. No diagram generation code exists
3. Only shape mapping documentation exists

---

## Specific Code Locations

### Broken Extraction

**sysml/parse.go:2428-2447** - BindingConnectorAsUsage handler
```go
// Creates connection but NEVER extracts ends from context
conn := NewConnection("", loc, false)
// Missing: Extract connector ends from ctx.BindingConnectorEnd() etc.
```

**sysml/parse.go:2458-2474** - SuccessionAsUsage handler
```go
// Creates transition but NEVER extracts source/target from context
trans := NewTransition("", loc)
// Missing: Extract source/target references from ctx
```

### Missing Resolution

**sysml/model.go:1876-1925** - Missing case
```go
// Add this case to ResolveReferences switch:
case *Dependency:
    m.resolveDependencyRefs(e)  // Function doesn't exist either!
```

---

## Recommended Fixes for v0.2

### Priority 1: Parser Fixes

1. **Fix `EnterBindingConnectorAsUsage`**
   - Extract connector ends from `ctx.BindingConnectorEnd()`
   - Add to `conn.unresolvedEnds`
   - Test with: `bind x.y to a.b`

2. **Fix `EnterSuccessionAsUsage`**
   - Extract source/target from succession context
   - Call `trans.SetUnresolvedSource()` / `SetUnresolvedTarget()`
   - Test with: `succession from State1 to State2`

3. **Fix `EnterFlowDefinition` / `EnterFlowUsage`**
   - Create FlowEnd elements for source/target
   - Extract references and populate FlowEnd.Reference
   - Test with: `flow from Source to Target`

### Priority 2: Model Fixes

4. **Add Dependency Resolution**
   ```go
   func (m *Model) resolveDependencyRefs(d *Dependency) {
       // Resolve clients
       for _, name := range d.unresolvedClient {
           if elem := m.findElement(name, d); elem != nil {
               d.Client = append(d.Client, elem)
           }
       }
       // Resolve suppliers
       for _, name := range d.unresolvedSupplier {
           if elem := m.findElement(name, d); elem != nil {
               d.Supplier = append(d.Supplier, elem)
           }
       }
   }
   ```

5. **Update `ResolveReferences` switch**
   - Add `case *Dependency: m.resolveDependencyRefs(e)`

### Priority 3: New Features

6. **Create `cmd/diagram` command**
   - Generate D2 output from parsed model
   - Render relations from SuccessionFlow, Connection, Dependency, etc.

7. **Create `cmd/watch` command**
   - File watching using `fsnotify` or similar
   - Auto-regenerate diagrams on file change
   - Optional: built-in HTTP server for preview

### Priority 4: Grammar Gaps

8. **Check grammar coverage**
   - Verify `BindingConnectorEnd` rules exist in grammar
   - Verify `Succession` source/target extraction context
   - Add tests for relation extraction

---

## Test Files for Verification

When fixes are implemented, test with these SysML patterns:

```sysml
// Test 1: Binding Connector
package Test {
    part def A;
    part def B;
    part a : A;
    part b : B;
    bind a.x to b.y;  // Binding connector
}

// Test 2: Succession  
package Test {
    state def States;
    state s1 : States;
    state s2 : States;
    succession from s1 to s2;  // Transition
}

// Test 3: Dependency
package Test {
    part def Client;
    part def Supplier;
    dependency from Client to Supplier;
}

// Test 4: Flow
package Test {
    action def Action1;
    action def Action2;
    flow from Action1 to Action2;
}
```

---

## Impact on v0.2 Planning

This investigation directly impacts v0.2 scope:

1. **Parser Completeness**: Grammar coverage is 73%, but relation element extraction is incomplete
2. **Feature Gap**: Diagram generation was assumed to exist but doesn't
3. **API Design**: Current model API supports relations but data is often empty
4. **Priority**: Fix parser extraction before implementing diagram generation

**Recommendation**: Include "Complete relation extraction" and "Basic diagram generation" in v0.2 milestones.

---

## Documentation Created

- `.planning/notes/diagram-relations-investigation.md` (this file)
- `.planning/notes/api-issues.md` (related API findings)

---

*Investigation completed: Task 006, Phase Quick*  
*Investigator: GSD Plan Executor*  
*Date: 2026-02-08*
