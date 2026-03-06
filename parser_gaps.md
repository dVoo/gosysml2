# Parser Gaps / Known Bugs

## 1. FQDN truncated for anonymous redefinition members (`part :>> name`)

### Symptom
Elements declared with a redefinition shorthand and no explicit own name produce a
qualified name that ends with `::` — the parent path is correct but the element's
own name component is missing.

**Example** (from `Batmobile.sysml`):
```sysml
part def BatmobileConfigurations :> Batmobile {
    part :>> batmobileEngine : EngineChoices;
    part :>> wheels : WheelChoices [4];   // ← broken
}
```
`wheels` is displayed as `Batmobile::BatmobileConfigurations::` instead of
`Batmobile::BatmobileConfigurations::wheels`.

### Root Cause

**`EnterPartUsage`** (`gosysml2/sysml/parse.go`, line 1404):

```go
func (b *modelBuilder) EnterPartUsage(ctx *parser.PartUsageContext) {
    name := ""
    ...
    if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
        usageDecl := ctx.Usage().UsageDeclaration()
        if ident := usageDecl.Identification(); ident != nil {
            name = extractName(ident)   // ← nil when syntax is `part :>> foo`
        }
        ...
    }
    part := NewPart(name, ...)   // name == "" for anonymous redefinitions
```

When the syntax is `part :>> wheels`, `Identification()` returns `nil` (there is
no explicit name token — the name comes from the redefinition target). The element
is created with `name = ""`.

**`baseElement.QualifiedName()`** (`gosysml2/sysml/model.go`, line 420):

```go
e.cachedQN = parentQN + "::" + e.name   // becomes "...::BatmobileConfigurations::"
```

When `e.name == ""` the `"::"` separator is still appended, producing a trailing
`::`.

### Fix Needed

In `EnterPartUsage` (and analogously in `EnterItemUsage`, `EnterAttributeUsage`,
and any other usage visitors that call `extractName(ident)` and may get empty
names), derive the element name from the redefinition chain when `Identification()`
is nil.

The redefinition target name is already extracted into `redefinedRefs` (see
`extractRedefinitionRefs`). The last path component of the first redefinition ref
should become the element's name:

```go
// After the existing name extraction block:
if name == "" && len(redefinedRefs) > 0 {
    // Derive name from first redefinition target (last path component of QN)
    ref := redefinedRefs[0]
    if i := strings.LastIndex(ref, "::"); i >= 0 {
        name = ref[i+2:]
    } else {
        name = ref
    }
}
```

Alternative (simpler, lower risk): fix `QualifiedName()` to not append `"::"` when
`e.name == ""`, so unnamed elements don't pollute the index with broken QNs:

```go
func (e *baseElement) QualifiedName() string {
    ...
    if e.name == "" {
        e.cachedQN = ""   // anonymous — no qualified name
    } else if e.parent == nil {
        e.cachedQN = e.name
    } else {
        parentQN := e.parent.QualifiedName()
        if parentQN == "" {
            e.cachedQN = e.name
        } else {
            e.cachedQN = parentQN + "::" + e.name
        }
    }
    ...
}
```

The parser-side fix (deriving name from redefinition) is preferred because it
correctly assigns the inherited name and allows the element to be found by
qualified name lookup.

### Affected Constructs
- `part :>> foo`
- `item :>> foo`
- `attribute :>> foo`
- Any usage member without an explicit `Identification()` but with a
  FeatureSpecialization redefinition (`:>>`).

### Test Case
```sysml
package BugRepro {
    part def Base { part x; }
    part def Child :> Base {
        part :>> x : Base;   // should have QN "BugRepro::Child::x"
    }
}
```
Currently produces `BugRepro::Child::` for the redefined `x` member.
