# Custom Attributes in Requirements and Verifications - Analysis Report

**Date:** 2026-01-31
**Question:** Are custom attributes from requirements and validations parsed to the high-level model?

## Answer: PARTIALLY ✅ ⚠️

Custom attributes **are parsed and extracted**, but there's a **parent assignment bug** that affects the model structure.

## Test Results

### Test File Created

```sysml
package TestCustomAttributes {
    // Requirement with custom attributes
    requirement def MyRequirement {
        attribute priority : String;
        attribute status : String;
        attribute owner : String;
        attribute effort : Real;

        require constraint { effort > 0 }
    }

    // Verification with attributes
    verification def MyVerification {
        subject testSubject;

        attribute testType : String;
        attribute duration : Real;

        objective {
            verify requirement : MyRequirement;
        }
    }
}
```

### What Gets Extracted

✅ **Attributes are extracted:**
- All 4 attributes from MyRequirement (priority, status, owner, effort)
- All 2 attributes from MyVerification (testType, duration)
- Total: 6 attributes successfully parsed

```
Element Statistics:
  Total elements: 4

  Breakdown by type:
    verification        : 1
    package             : 1
    requirement         : 2
    attribute           : 6    <-- All extracted!
```

### ⚠️ The Bug: Incorrect Parent Assignment

**Problem:** Attributes inside requirements/verifications are being added to the **package** instead of their containing element.

**Expected Structure:**
```
package: TestCustomAttributes
  ├─ requirement: MyRequirement
  │  ├─ attribute: priority        <-- Should be here
  │  ├─ attribute: status           <-- Should be here
  │  ├─ attribute: owner            <-- Should be here
  │  └─ attribute: effort           <-- Should be here
  └─ verification: MyVerification
     ├─ attribute: testType         <-- Should be here
     └─ attribute: duration          <-- Should be here
```

**Actual Structure:**
```
package: TestCustomAttributes
  ├─ requirement: MyRequirement
  ├─ attribute: priority        <-- WRONG PARENT
  ├─ attribute: status           <-- WRONG PARENT
  ├─ attribute: owner            <-- WRONG PARENT
  ├─ attribute: effort           <-- WRONG PARENT
  ├─ verification: MyVerification
  ├─ attribute: testType         <-- WRONG PARENT
  └─ attribute: duration          <-- WRONG PARENT
```

### Root Cause

**File:** `gosysml2/sysml/parse.go`
**Functions:** `EnterAttributeDefinition()` (line 779) and `EnterAttributeUsage()` (line 804)

**Buggy Code:**
```go
func (b *modelBuilder) EnterAttributeDefinition(ctx *parser.AttributeDefinitionContext) {
    // ... extract name and location ...

    attr := NewAttribute(name, loc, true)

    // BUG: Always adds to currentPkg, doesn't check elementStack
    if b.currentPkg != nil {
        attr.parent = b.currentPkg          // <-- WRONG
        b.currentPkg.AddChild(attr)         // <-- WRONG
    }
}
```

**Correct Pattern (used by Requirements):**
```go
func (b *modelBuilder) EnterRequirementDefinition(ctx *parser.RequirementDefinitionContext) {
    // ... extract name and location ...

    req := NewRequirement(name, loc, true)

    if b.currentPkg != nil {
        req.parent = b.currentPkg
        b.currentPkg.AddChild(req)
    }

    // Push requirement onto stack for nested elements
    b.elementStack = append(b.elementStack, req)  // <-- This is the key!
}
```

The model builder uses an `elementStack` to track nested elements:
- Requirements/Verifications/Parts push themselves onto the stack when entered
- They pop themselves when exited
- Child elements should check the stack to find their correct parent

**BUT** attributes don't check the stack - they always add to `currentPkg`.

## Impact Assessment

### What Works ✅

1. **Attributes are parsed:** All custom attributes are recognized by the lexer/parser
2. **Attributes are extracted:** The high-level model contains all attribute elements
3. **Attribute properties captured:**
   - Name
   - Type reference
   - Source location (line/column)
   - IsDefinition flag
   - IsReadOnly, IsDerived flags

### What's Broken ⚠️

1. **Parent-child relationships:** Attributes don't appear as children of their containing elements
2. **Navigation:** Can't do `requirement.Children()` to find its attributes
3. **Scoping:** Attributes appear at package scope instead of requirement/verification scope

### Workaround

You can still find attributes, but you need to search by location:

```go
// Find requirement
req := findRequirementByName(model, "MyRequirement")
reqLine := req.Location().Line

// Find attributes declared right after the requirement
attrs := make([]*Attribute, 0)
for _, elem := range model.FindByKind(sysml.KindAttribute) {
    attr := elem.(*sysml.Attribute)
    // Attributes right after req definition, before next major element
    if attr.Location().Line > reqLine && attr.Location().Line < reqLine + 10 {
        attrs = append(attrs, attr)
    }
}
```

This is fragile and not recommended for production use.

## Constraint Extraction ✅

**Good news:** Requirement constraints (assume/require) **are correctly extracted** and properly linked to their requirements.

```
Requirement: MyRequirement (def=true)
  Constraints (1):
    - <unnamed> (require constraint)
      Expression: constraint{effort>0}
```

The `RequirementConstraint` type includes:
- `IsAssume` flag (true for assume, false for require)
- `Expression` field (captured as raw text)
- Proper parent relationship

## Verification Elements

Similar issues affect verifications:
- Verification subjects are **not extracted**
- Verification method is **not extracted** (always shows "unspecified")
- Attributes have wrong parent

## Recommendations

### For Users (Current State)

**If you need to extract custom attributes:**

1. ⚠️ **Known limitation:** Attributes will be at package level, not nested under requirements
2. ✅ **Workaround:** Use location-based filtering or search all attributes in package
3. ✅ **Type references work:** You can still get attribute types via `attr.TypeRef`
4. ✅ **All metadata available:** Names, locations, flags are all captured

**What you CAN reliably extract today:**
- Requirement structure (definitions and usages)
- Requirement constraints (assume/require)
- Verification existence (but not subjects/methods)
- Attribute existence (but not proper nesting)
- All type references and relationships

### For Developers (Fix Required)

To properly support custom attributes in requirements/verifications, the model builder needs updates:

**Required Changes in `gosysml2/sysml/parse.go`:**

1. Attributes should check element stack before adding to package
2. Verifications should push themselves onto element stack
3. Port definitions should push themselves onto element stack
4. Update all child element handlers to use stack-aware parent assignment

**Example Fix:**
```go
func (b *modelBuilder) EnterAttributeDefinition(ctx *parser.AttributeDefinitionContext) {
    // ... extract name and location ...

    attr := NewAttribute(name, loc, true)

    // Check element stack first
    if len(b.elementStack) > 0 {
        parent := b.elementStack[len(b.elementStack)-1]
        attr.parent = parent
        parent.AddChild(attr)
    } else if b.currentPkg != nil {
        attr.parent = b.currentPkg
        b.currentPkg.AddChild(attr)
    }
}
```

## Conclusion

**Short answer:** Yes, custom attributes are parsed and extracted, but their parent-child relationships are incorrectly assigned.

**For most use cases:** This is a non-critical bug if you only need to know what attributes exist. The attribute data is all there, just not properly nested in the tree structure.

**For hierarchical analysis:** This is a blocking issue if you need to programmatically determine which attributes belong to which requirements.

**Bottom line:** The parser gets all the information from the source code - the bug is in how it's organized in the resulting model structure.
