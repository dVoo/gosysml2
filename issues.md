# gosysml2 Parser Integration Issues

**Date:** 2026-02-08  
**Component:** internal/sysml/parser.go  
**Library:** github.com/dVoo/gosysml2  

---

## Issue 1: Recursive Package Traversal Not Implemented

### Summary
The parser recognizes `Package` elements but does not recursively traverse their contents to find nested elements like `part def`, `item def`, etc.

### Affected File
`internal/sysml/parser.go`, lines 182-191

### Current Implementation
```go
switch e := elem.(type) {
case *gosysml2.Package:
    // Process package contents by iterating through its typed fields
    // For now, we'll just note the package exists
    _ = e
case *gosysml2.Part:
    // Part is already added as a block above
    // Connections would be processed separately
    _ = e
}
```

### Problem
When processing `*gosysml2.Package`, the code only assigns to `_` (discard), ignoring all nested elements. This means:

- Top-level package is recognized but its contents are ignored
- Nested packages are never traversed
- Parts, items, blocks inside packages are never extracted

### Impact
Files with nested package structures produce empty diagrams:
- `parts_tree.sysml` - All parts inside `Definitions` and `Usages` packages are lost
- Any real-world SysML v2 project using packages shows nothing

### Root Cause
The `gosysml2.Package` type has an `OwnedElement` field (or similar) that contains all nested elements, but the parser doesn't access it.

### Required Fix
```go
case *gosysml2.Package:
    // Recursively process all elements owned by this package
    for _, ownedElem := range e.OwnedElement {
        processElement(ast, ownedElem, fullPath)
    }
```

Note: Need to verify exact field name in gosysml2 (might be `Elements`, `OwnedElements`, `Members`, etc.)

---

## Issue 2: Part Definitions Not Distinguished from Part Usages

### Summary
The parser treats all parts the same way, but SysML v2 distinguishes between:
- `part def TypeName` (type definition)
- `part instanceName : TypeName` (usage/instance)

### Current Behavior
Both create a Block with Type="Part", but we lose the distinction between:
- The type definition (should appear as a block/node)
- The usage/instance (could be shown as connection or nested structure)

### Example from parts_tree.sysml
```sysml
package Definitions {
    part def Vehicle { ... }  // Type definition - should be node
    part def Axle { ... }     // Type definition - should be node
}

package Usages {
    part vehicle1: Vehicle { ... }  // Usage - could be instance or connection
}
```

### Suggested Enhancement
Add a field to Block to distinguish definition vs usage:
```go
type Block struct {
    Name       string
    Type       string
    IsDef      bool       // true for "part def", false for "part name: Type"
    FullPath   string
    Attributes map[string]string
}
```

---

## Issue 3: Nested Parts Not Connected

### Summary
When a part contains nested parts (composition), the parser should create connections between parent and child.

### Example
```sysml
part vehicle1: Vehicle {
    part frontAxleAssembly: AxleAssembly {
        part frontAxle: Axle;
        part frontWheel: Wheel[2];
    }
}
```

### Expected
Should generate connections:
```
vehicle1 -> frontAxleAssembly
frontAxleAssembly -> frontAxle
frontAxleAssembly -> frontWheel
```

### Current
Nested parts are lost because:
1. Package isn't traversed (Issue #1)
2. Even if it were, nested parts aren't extracted as connections

---

## Issue 4: Imports Not Resolved During Parsing

### Summary
The parser doesn't resolve import statements like:
```sysml
private import SI::kg;
private import Definitions::*;
```

This means:
- Types defined in other packages can't be resolved
- Cross-package references are broken
- The `ImportResolver` exists but isn't integrated during initial parsing

### Note
This might be intentional separation of concerns (parse first, resolve later), but needs verification.

---

## Issue 5: Attributes and Properties Not Extracted

### Summary
Part attributes are mentioned in comments but not actually extracted:

```go
// Add attributes as tooltip or label if present
if len(block.Attributes) > 0 {
    // ... tooltip code
}
```

But in `processElement`, the Attributes map is never populated:
```go
block := Block{
    Name:       name,
    Type:       elemKind,
    FullPath:   fullPath,
    Attributes: make(map[string]string),  // Empty map created, never filled
}
```

### Missing Data
- `attribute mass :> ISQ::mass`
- `attribute steeringAngle: ScalarValues::Real`
- Default values, redefinitions, etc.

---

## Recommended Priority Order

1. **Issue #1 (Critical)** - Fix recursive package traversal - without this, most real files won't work
2. **Issue #3 (High)** - Connect nested parts to show hierarchy
3. **Issue #2 (Medium)** - Distinguish definitions from usages for better diagrams
4. **Issue #5 (Medium)** - Extract attributes for tooltips/labels
5. **Issue #4 (Low)** - Import resolution (may already work via separate mechanism)

---

## Testing Recommendations

Create test files covering:
- [ ] Simple flat file (no packages) - should work
- [ ] Single nested package (one level) - Issue #1
- [ ] Multiple nested packages (two+ levels) - Issue #1
- [ ] Parts with attributes - Issue #5
- [ ] Parts with nested parts (composition) - Issue #3
- [ ] Part definitions vs usages - Issue #2

---

## Related Files

- `internal/sysml/parser.go` - Main issue location
- `internal/sysml/ast.go` - Block struct definition
- `internal/transformer/transformer.go` - Consumes parsed AST
- `testdata/*.sysml` - Test cases needed

---

*Document created during Phase 4 UAT debugging*  
*Issue discovered: parts_tree.sysml produces empty output*
