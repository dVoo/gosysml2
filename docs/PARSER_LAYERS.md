# Parser Layers: Low-Level and High-Level APIs

This project exposes two parser layers:

- `low` package: thin wrapper around ANTLR lexer/parser
- `sysml` package: high-level model builder and reference resolver

Use `low` when you need raw parser behavior or parse-tree control. Use `sysml` for almost all application code.

## Layer 1: Low-Level Parser (`low`)

The `low` package gives direct access to lexing/parsing with minimal abstraction.

### Core types

- `Lexer`: wraps `SysMLv2Lexer`, collects lexer errors
- `Parser`: wraps `SysMLv2Parser`, collects parser errors
- `SyntaxError`: single lexer/parser error with line/column/source
- `ParseErrors`: combined lexer+parser error set

### Main entry points

- `low.Parse(input string, opts ...ParseOption) (tree, *ParseErrors)`
- `low.ParseBytes(input []byte, opts ...ParseOption) (tree, *ParseErrors)`
- `low.ParseWithContext(ctx, input, opts...) (tree, error)`
- `low.Validate(input string) *ParseErrors`
- `low.ValidateBytes(input []byte) *ParseErrors`

### Low-level options

- `low.WithParseTree(bool)`: enable/disable parse tree construction
- `low.WithContext(ctx)`: cancellation support

### Typical use cases

- custom parse-tree walking on ANTLR contexts
- syntax-only validation at high throughput
- tooling that needs token streams and token names

## Layer 2: High-Level SysML Model (`sysml`)

The `sysml` package builds a typed data model from the low-level parse tree.

### Parse pipeline

1. Parse textual input via `low.Parse(...)`
2. Build model with internal listener (`buildModel(...)`)
3. Build qualified-name index (`Model.BuildIndex()`)
4. Resolve references (`Model.ResolveReferences()`)
5. Return `ParseResult` with model, errors, optional parse tree

### Identification and short-name handling

SysML/KerML `Identification` allows both:

- `declaredShortName` via `<...>`
- `declaredName` as the element name

The high-level layer now treats these separately:

- `Name()` maps to `declaredName` (cross-reference/stable model name)
- `DeclaredShortName()` preserves the `<...>` identifier

Reference resolution consults both normal qualified-name lookup and
short-name lookup, so relationships can resolve by long name or short name.

Part usages also preserve multiplicity from `FeatureSpecializationPart`:

- `part wheels : Wheel[4]` -> `Part.Multiplicity == "4"`
- `part optionalWheel : Wheel[0..1]` -> `Part.Multiplicity == "0..1"`
- `part anyWheels : Wheel[*]` -> `Part.Multiplicity == "*"`

Attribute usages preserve both typing and inline default values:

- `attribute mass : Real = 1500.0;` ->
  `Attribute.TypeRef.Name() == "Real"` and `Attribute.DefaultValue == "1500.0"`
- `attribute maxSpeed : Integer = 200;` ->
  `Attribute.TypeRef.Name() == "Integer"` and `Attribute.DefaultValue == "200"`

Requirement compatibility handling includes:

- `require { massActual <= massLimit };` inside a requirement definition
  is accepted and exposed as `RequirementConstraint.Expression == "massActual <= massLimit"`.
- requirement usage argument lists after typing are accepted:
  `requirement <'R1'> massReq : MassRequirement [vehicle = myVehicle]`
  and exposed via `Requirement.Bindings`.

Feature relationships are also exposed on `Attribute`:

- `attribute speedA :> velocity;` -> `speedA.SubsettedFeatures` contains `velocity`
- `attribute speedB ::> velocity;` -> `speedB.SubsettedFeatures` contains `velocity`
- `attribute speedC :>> velocity;` -> `speedC.RedefinedFeatures` contains `velocity`

### Core types

- `ParseResult`
  - `Model *Model`
  - `Errors *ParseError`
  - `Tree antlr.Tree` (optional; can be discarded)
  - `Source string`
- `ParseError` / `Error`: high-level parse error model
- `Model`: root aggregate with typed top-level collections

### Main entry points

- `sysml.ParseString(input, opts...)`
- `sysml.ParseFile(path, opts...)`
- `sysml.ParseBytes(data, source, opts...)`
- `sysml.ParseReader(r, source, opts...)`
- `sysml.ParseDirectory(dir, opts...)`
- `sysml.ParseDirectoryParallel(dir, workers, opts...)`
- `sysml.ParseDirectoryStream(dir, handler, opts...)`
- `sysml.Validate(input)`
- `sysml.ValidateFile(path)`

### High-level options

- `sysml.WithDiscardTree()`: keep model, drop parse tree
- `sysml.WithLibraryRegistry(reg)`: use preloaded library resolver
- `sysml.WithStandardLibrary()`: auto-load standard library
- `sysml.WithLibraryPath(path)`: custom standard library path

### Typical use cases

- domain tooling (requirements traceability, architecture checks)
- model traversal with visitors/finders
- working with resolved typed references (`Ref[T]`)

## Choosing the right layer

- Choose `low` if you need ANTLR-level control, token streams, or pure syntax validation.
- Choose `sysml` if you need semantic model traversal, typed elements, and resolved cross-references.
- Use both when needed: parse with `sysml` for model semantics, fall back to `low` only for parser internals.
