# SysML High-Level Data Model (Currently Implemented)

This document describes the `sysml` package data model as currently implemented in code.

## Foundation

### Core abstractions

- `ElementKind`: enum for element categories (`KindPackage`, `KindPart`, `KindRequirement`, etc.)
- `Location`: source span (`Line`, `Column`, `EndLine`, `EndColumn`)
- `Element` interface:
  - `Kind()`, `Name()`, `QualifiedName()`, `Location()`
  - `Parent()`, `Children()`, `SetParent(...)`
  - `Documentation()`
- `Definition` and `Usage` marker interfaces
- `Ref[T Element]`: typed reference that can be unresolved/resolved
- concrete elements also preserve `DeclaredShortName()` when source uses `<...>`

### Root container

- `Model`
  - typed top-level collections:
    - `Packages`, `Imports`, `Comments`, `Dependencies`, `Docs`
    - `Flows`, `ControlNodes`, `Occurrences`
    - `Aliases`, `Metadata`, `Renderings`, `Messages`, `Filters`
    - `Satisfies`, `Verifies`
  - `Elements []Element` generic top-level list
  - supports `BuildIndex()`, `FindByQualifiedName(...)`, `ResolveReferences()`
  - index behavior:
    - qualified-name lookup
    - short-name lookup (for `declaredShortName`)
    - dotted-chain fallback (`a.b` -> `a::b`) for relationship resolution

## Structural and architecture elements

- `Package`
  - `IsLibrary bool`
  - typed children accessors (`Packages()`, `Parts()`, `Requirements()`, etc.)
- `Part` (definition or usage)
  - `IsDefinition`, `TypeRef Ref[*Part]`
  - `Multiplicity string` (for usages, extracted from `[ ... ]` bounds such as `4`, `0..1`, `*`)
  - specialization via `Specializes Ref[*Part]`
  - typed children: attributes/parts/ports
- `Item` (definition or usage)
  - `IsDefinition`, `TypeRef Ref[*Item]`
  - specialization via `Specializes Ref[*Item]`
  - typed children: attributes/items
- `Attribute`
  - `IsDefinition`, `TypeRef Ref[Element]`
  - `DefaultValue`, `IsReadOnly`, `IsDerived`
  - for usage forms (`attribute x : T = v`), parser populates both `TypeRef` and `DefaultValue`
- `Port`
  - `IsDefinition`, `TypeRef Ref[*Port]`, `Direction`
  - optional `ConjugatedPort`
- `ConjugatedPort`
  - `OriginalPort Ref[*Port]`
- `Interface`
  - `IsDefinition`, `TypeRef Ref[*Interface]`
  - typed child ports
- `Connection`
  - `IsDefinition`, `TypeRef Ref[*Connection]`
  - `Ends []*ConnectionEnd`
- `ConnectionEnd`
  - `EndRef Ref[Element]`
- `Allocation`
  - `IsDefinition`, `TypeRef Ref[*Allocation]`
  - `Source Ref[Element]`, `Target Ref[Element]`
- `Dependency`
  - `Client []Element`, `Supplier []Element`

## Requirements, verification, and cases

- `Constraint`
  - `IsDefinition`, `TypeRef Ref[*Constraint]`, `Expression`
- `RequirementConstraint`
  - `IsAssume`, `Expression`
- `Requirement`
  - `IsDefinition`, `TypeRef Ref[*Requirement]`, `RequirementID`
  - `RequirementID` is populated from `declaredShortName` when `<...>` is present
  - `Subject Ref[Element]`
  - relationships:
    - `DerivedFrom`, `DerivedReqs`
    - `SatisfiedBy []Element`
    - `VerifiedBy []*Verification`
  - `Assumptions`, `Constraints`
- `Verification`
  - `IsDefinition`, `TypeRef Ref[*Verification]`
  - `Subject Ref[Element]`
  - `VerifiedRequirement *Requirement`
  - `Method VerificationMethod`
  - typed child actions
- `Concern`
  - `IsDefinition`, `TypeRef Ref[*Concern]`
  - `Stakeholders []Element`
- `UseCase`
  - `IsDefinition`, `TypeRef Ref[*UseCase]`
  - `Subject Ref[Element]`, `Actors []Element`
  - `IncludedUseCases []*UseCase`
- `IncludeUseCase`
  - usage relationship element
  - `IncludedUseCase Ref[*UseCase]`, `Owner Ref[*UseCase]`
- `Case`
  - `IsDefinition`, `TypeRef Ref[*Case]`
  - `Subject Ref[Element]`, `Actors []Element`, `Objectives []*Requirement`
- `AnalysisCase`
  - `IsDefinition`, `TypeRef Ref[*AnalysisCase]`
  - `Subject Ref[Element]`, `ReturnType Ref[Element]`
- typed relationship edges:
  - `SatisfyRelationship` (`Satisfier -> Requirement`)
  - `VerifyRelationship` (`Verification -> Requirement`)
  - unresolved names are also preserved in `Ref.Name()` for diagnostics
  - verifier resolution prefers the enclosing `Verification` context when parsing `verify` usage forms

## Behavioral, flow, and time model elements

- `Action`
  - `IsDefinition`, `TypeRef Ref[*Action]`
  - typed nested actions
- `State`
  - `IsDefinition`, `TypeRef Ref[*State]`
  - `EntryAction`, `DoAction`, `ExitAction`
  - nested states and transitions
- `Transition`
  - `Source Ref[*State]`, `Target Ref[*State]`
  - `TriggerExpr`, `GuardExpr`, `EffectAction`
- `ControlNode`
  - `NodeKind` (`merge`, `decision`, `join`, `fork`)
  - optional `Condition`
- `Flow`
  - `IsDefinition`
  - `Source *FlowEnd`, `Target *FlowEnd`
  - `PayloadFeatures []Element`
- `FlowEnd`
  - `Reference Element`, `Feature Element`
- `SuccessionFlow`
  - usage element with `Source Ref[Element]`, `Target Ref[Element]`, `Guard`
- `Occurrence`
  - `IsDefinition`, `IsIndividual`, `IsEvent`, `IsTimeSlice`, `IsSnapshot`
  - `PortionKind`, `LifeStep`
  - nested occurrences

## Views, metadata, and annotations

- `Viewpoint`
  - `IsDefinition`, `TypeRef Ref[*Viewpoint]`
  - `Concerns []*Concern`, `Stakeholders []Element`
- `View`
  - `IsDefinition`, `TypeRef Ref[*View]`
  - `ExposedElements []Element`, `Viewpoint Ref[*Viewpoint]`
- `Alias`
  - `Target Ref[Element]`
- `Metadata`
  - `IsDefinition`, `TypeRef Ref[*Metadata]`
  - `Annotations []*PrefixMetadataAnnotation`
- `PrefixMetadataAnnotation`
  - `Metadata Ref[*Metadata]`, `Values []string`
- `Rendering`
  - `IsDefinition`, `TypeRef Ref[*Rendering]`, `Body`
- `ElementFilter`
  - `Expression`
- `Message`
  - usage element
  - `Sender Ref[Element]`, `Receiver Ref[Element]`, `Payload`

## Enumerations and calculations

- `Enumeration`
  - `IsDefinition`, `TypeRef Ref[*Enumeration]`
  - `Values []*EnumerationValue`
- `EnumerationValue`
  - single enum literal
- `Calculation`
  - `IsDefinition`, `TypeRef Ref[*Calculation]`
  - `ReturnType Ref[Element]`, `Expression`

## Imports and documentation elements

- `Import`
  - `ImportedNamespace`
  - import semantics flags:
    - `IsMembership`, `IsNamespace`, `IsRecursive`, `IsAll`
  - `FilterExpressions []string`
  - resolution fields:
    - `ResolvedElement Element`
    - `ResolvedPackage *Package`
    - `IsResolved bool`
- `Comment`
  - `Body`, `Locale`, `About []Element`
- `Doc`
  - `Body`, `Locale`

## Traversal and querying API

- visitor-based traversal: `Visitor`, `BaseVisitor`, `Visit(model, visitor)`
- generic traversal: `Walk(model, fn)`
- iterators:
  - `All(model)`, `OfType[T](model)`, `OfKind(model, kind)`
- finder helpers:
  - `FindPackages`, `FindParts`, `FindRequirements`, `FindVerifications`, etc.
  - `FindByKind`, `FindByName`, `FindByQualifiedName`, `Filter`

## Notes on resolution behavior

- Most relationship fields are parsed as unresolved names first.
- `Model.ResolveReferences()` resolves references to typed pointers/`Ref[T]`.
- `Model.BuildIndex()` creates qualified-name and short-name indices used by lookup.
- nested `SatisfyRelationship` / `VerifyRelationship` nodes are tracked in
  `Model.Satisfies` / `Model.Verifies` without duplicate traversal entries in `Model.Elements`.
- If a `LibraryRegistry` is attached, qualified names can resolve into standard library elements as well.
