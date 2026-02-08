---
phase: quick-003
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - sysml/parse_test.go
  - sysml/model.go
  - sysml/parse.go
autonomous: true

must_haves:
  truths:
    - "Nested parts create proper parent-child relationships via Children()"
    - "Part usages have correct TypeRef resolved from type name"
    - "Attributes are accessible via Part.Attributes() method"
    - "Part definitions vs usages are distinguishable via IsDefinition field"
  artifacts:
    - path: "sysml/parse_test.go"
      provides: "Unit tests for nested parts, type resolution, and attributes"
      min_lines: 100
  key_links:
    - from: "EnterPartUsage"
      to: "Part.TypeRef"
      via: "type reference extraction from UsageDeclaration"
    - from: "PartUsage children"
      to: "Parent.Children()"
      via: "addToParent with elementStack"
---

<objective>
Investigate and fix parser issues related to nested parts connectivity, type reference resolution for part usages, and attribute extraction. The key findings from verification indicate:
- Issue #1 (Package traversal): Already working correctly
- Issue #2 (Part defs vs usages): IsDefinition field exists and is populated
- Issue #3 (Nested parts): Needs investigation - may be missing type reference resolution
- Issue #5 (Attributes): Needs verification if attributes are properly attached to parts

Purpose: Ensure nested parts properly reference their parent and have their type references resolved, and that attributes are accessible from parts.
Output: Working tests demonstrating correct nested part behavior, type resolution, and attribute access.
</objective>

<execution_context>
@/home/daniel/.config/opencode/get-shit-done/workflows/execute-plan.md
@/home/daniel/.config/opencode/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/STATE.md
@issues.md

@sysml/parse.go
@sysml/model.go
@sysml/visitor.go

Test file showing nested parts structure:
@validationdata/01-Parts Tree/1a-Parts Tree.sysml
</context>

<tasks>

<task type="auto">
  <name>Task 1: Create comprehensive test for nested parts and type resolution</name>
  <files>sysml/parse_test.go</files>
  <action>
Create a new test file sysml/parse_test.go with tests to verify:

1. **Nested parts parent-child relationship test**:
   - Parse the "1a-Parts Tree.sysml" test file
   - Find vehicle1 part usage
   - Verify vehicle1.Children() contains frontAxleAssembly and rearAxleAssembly
   - Verify frontAxleAssembly.Children() contains frontAxle and frontWheel

2. **Part definition vs usage distinction test**:
   - Verify Vehicle part def has IsDefinition=true
   - Verify vehicle1 part usage has IsDefinition=false
   - Check that definitions are in Definitions package
   - Check that usages are in Usages package

3. **Type reference resolution test**:
   - Check if vehicle1 (part usage) has TypeRef resolved to Vehicle (part def)
   - Check if frontAxleAssembly has TypeRef resolved to AxleAssembly
   - If TypeRef is not being resolved, note this as the actual bug to fix

4. **Attribute extraction test**:
   - Verify Vehicle part def has mass attribute accessible via Attributes()
   - Verify attributes have proper names and default values

The test should use the existing ParseFile function and model traversal methods (FindAll, Walk, etc.).

Key code patterns to use:
```go
result := ParseFile("validationdata/01-Parts Tree/1a-Parts Tree.sysml")
model := result.Model

// Find all parts
parts := FindAll[*Part](model)

// Walk the model
Walk(model, func(elem Element, depth int) bool {
    if part, ok := elem.(*Part); ok {
        // Check part properties
    }
    return true
})
```
  </action>
  <verify>go test ./sysml -run TestNestedParts -v</verify>
  <done>Test file exists and reveals the actual state of nested parts, type references, and attributes</done>
</task>

<task type="auto">
  <name>Task 2: Fix type reference extraction for part usages</name>
  <files>sysml/parse.go</files>
  <action>
Based on the test results from Task 1, if TypeRef is not being populated for part usages:

1. Look at EnterPartUsage (lines 465-485) - it extracts name but NOT the type reference
2. The type reference comes from the UsageDeclaration which has a FeatureSpecializationPart containing the type name
3. Extract the type name from ctx.Usage().UsageDeclaration().FeatureSpecializationPart() and set it as the unresolved type reference

The fix should look like:
```go
func (b *modelBuilder) EnterPartUsage(ctx *parser.PartUsageContext) {
    name := ""
    typeRef := ""
    
    if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
        usageDecl := ctx.Usage().UsageDeclaration()
        if ident := usageDecl.Identification(); ident != nil {
            name = extractName(ident)
        }
        // Extract type reference from FeatureSpecializationPart
        if featSpec := usageDecl.FeatureSpecializationPart(); featSpec != nil {
            typeRef = extractTypeReference(featSpec)
        }
    }

    part := NewPart(name, locationFromContext(ctx), false)
    part.parent = b.getCurrentParent()
    if typeRef != "" {
        part.TypeRef = NewRef[*Part](typeRef)
    }
    b.addToParent(part)
    b.elementStack = append(b.elementStack, part)
}
```

You'll need to create or use an existing helper function to extract the type reference name from FeatureSpecializationPart context. Look at how extractName works and create a similar extractTypeReference function if needed.
  </action>
  <verify>go test ./sysml -run TestTypeReference -v passes</verify>
  <done>Part usages have their TypeRef.name populated from the type specification in the source</done>
</task>

<task type="auto">
  <name>Task 3: Verify and document attribute extraction</name>
  <files>sysml/parse.go, sysml/model.go</files>
  <action>
1. Review EnterAttributeUsage (lines 1058-1143) to confirm attributes ARE being extracted properly
2. Verify that attributes are being added to Part.children and Part.attributes via AddChild
3. If working correctly, ensure the test from Task 1 passes for attributes
4. If NOT working, debug why attributes aren't being attached to parts

The EnterAttributeUsage already extracts:
- Attribute name from UsageDeclaration
- DefaultValue from UsageCompletion.ValuePart.FeatureValue.OwnedExpression
- Adds to parent via getCurrentParent() and type-specific AddChild calls

If attributes are not appearing on parts, the issue might be:
- The parent element stack isn't correctly set when parsing nested attributes
- The AddChild switch case for Part isn't being hit

Add debug logging or additional test assertions to verify the attribute flow.
  </action>
  <verify>go test ./sysml -run TestAttributes -v passes</verify>
  <done>Attributes are properly extracted and accessible via Part.Attributes() method</done>
</task>

</tasks>

<verification>
After all tasks:
1. Run all new tests: `go test ./sysml -v`
2. Ensure existing tests still pass: `go test ./... -v`
3. Verify the Parts Tree example can be parsed and traversed correctly
4. Confirm part definitions vs usages are distinguishable
5. Confirm nested parts have proper parent relationships
6. Confirm type references are resolved or at least populated for later resolution
</verification>

<success_criteria>
- [ ] Test file exists with comprehensive coverage of nested parts
- [ ] Part usages have TypeRef populated with the type name from source
- [ ] Attributes are accessible from parts via Part.Attributes()
- [ ] Parent-child relationships work via Children() method
- [ ] All tests pass (new and existing)
</success_criteria>

<output>
After completion, create `.planning/quick/003-fix-parser-issues/003-SUMMARY.md`
</output>
