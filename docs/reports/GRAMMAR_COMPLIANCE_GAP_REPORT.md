# Grammar Compliance Gap Report

Date: 2026-02-18

## Scope checked

- Grammar sources:
  - `docs/bnf/SysML-textual-bnf.kebnf`
  - `docs/bnf/KerML-textual-bnf.kebnf`
- Parser entry path:
  - `sysml/parse.go:120` (`low.Parse`)
  - `low/parser.go:89` (`EntryRuleRootNamespace`)
- High-level extraction path:
  - `sysml/parse.go` (`modelBuilder` Enter/Exit handlers)
- High-level model:
  - `sysml/model.go` + element files (`sysml/comment.go`, `sysml/flow.go`, etc.)

## What is happening now

1. ANTLR parser accepts grammar-level syntax and builds parse trees.
2. High-level model extraction is selective and currently captures a subset.
3. Large portions of grammar parse successfully but are not mapped into model elements/fields.

## Coverage snapshot (listener callback parity)

- Total parser listener `Enter*` rules available: **448**
- `modelBuilder` `Enter*` handlers implemented: **61**
- Missing `Enter*` handlers: **387**

- `*Definition` rules:
  - total: 27
  - implemented: 22
  - missing: 5 (`EnterIndividualDefinition`, `EnterRenderingDefinition`, `EnterMetadataDefinition`, `EnterExtendedDefinition`, generic `EnterDefinition`)

- `*Usage` rules:
  - total: 61
  - implemented: 27
  - missing: 34

## Missing grammar features (major)

From SysML textual BNF, the following are defined but not extracted by `modelBuilder`:

- Package/member features:
  - `AliasMember` (`docs/bnf/SysML-textual-bnf.kebnf:142`)
  - `ElementFilterMember` (`docs/bnf/SysML-textual-bnf.kebnf:172`)
  - import subforms (`MembershipImport`, `NamespaceImport`, `FilterPackage`)

- Definition/usage elements:
  - `IndividualDefinition`, `IndividualUsage` (`docs/bnf/SysML-textual-bnf.kebnf:188`, `docs/bnf/SysML-textual-bnf.kebnf:575`)
  - `PortionUsage`, `EventOccurrenceUsage` (`docs/bnf/SysML-textual-bnf.kebnf:579`, `docs/bnf/SysML-textual-bnf.kebnf:588`)
  - `Message` (`docs/bnf/SysML-textual-bnf.kebnf:805`)
  - `PerformActionUsage`, `ExhibitStateUsage` (`docs/bnf/SysML-textual-bnf.kebnf:944`, `docs/bnf/SysML-textual-bnf.kebnf:1268`)
  - `SatisfyRequirementUsage` (`docs/bnf/SysML-textual-bnf.kebnf:1466`)
  - `RenderingDefinition`/`RenderingUsage` (`docs/bnf/SysML-textual-bnf.kebnf:1642`, `docs/bnf/SysML-textual-bnf.kebnf:1646`)
  - `MetadataDefinition`/`MetadataUsage` (`docs/bnf/SysML-textual-bnf.kebnf:1652`, `docs/bnf/SysML-textual-bnf.kebnf:1666`)

## Features parsed but currently discarded or flattened

- Import text is captured as raw `ctx.GetText()` (`sysml/parse.go:613`) and then suffix-parsed.
  - Loses structured import details (visibility/import-kind/filter structure).
- Dependency parsing stores only first client and first supplier (`sysml/parse.go:2143` to `sysml/parse.go:2151`).
  - Grammar allows lists on both sides.
- Many semantic fields in model exist but are not filled by parser:
  - Only `Part` type refs are set from parse (`sysml/parse.go:497`).
  - Other `TypeRef` fields (Requirement, Verification, UseCase, View, etc.) are not set during parse.
  - Most unresolved relationship setters are never called in parse (e.g. satisfaction/verification/stakeholder/view relationships).
- Relationship/membership expressions are often stored as raw text (`GetText`) rather than structured AST:
  - subjects/actors/objectives (`sysml/parse.go:1854`, `sysml/parse.go:1890`, `sysml/parse.go:1922`)
  - requirement constraints (`sysml/parse.go:1967`)
  - transitions/succession endpoints (`sysml/parse.go:1566`, `sysml/parse.go:2316`)

## Missing grammar elements in model (representation gaps)

- `KindMetadata` and `KindAlias` exist as enum constants (`sysml/model.go:34`, `sysml/model.go:36`) but there are no concrete `Metadata` or `Alias` element structs/builders.
- No model element for rendering.
- No model element for message usage.
- No dedicated model representation for:
  - element filters
  - prefix metadata usage/annotations
  - many action node forms (accept/send/assignment/etc.) as typed nodes

Model root collections also show current scope (`sysml/model.go:1759`): packages/imports/comments/dependencies/docs/flows/control nodes/occurrences.

## Resolver gap

- `ResolveReferences()` covers many element types (`sysml/model.go:1874`) but not dependency/comment-specific reference resolution.
- `Dependency` unresolved refs are stored, but there is no dependency resolution step in `ResolveReferences`.
- `Comment.about` unresolved refs are not populated by parser and not resolved.

## Compliance boundary

- Compliance target is **textual SysML v2/KerML grammar only**.
- Graphical notation grammar is out of scope for this parser and this report.

## Plan for 100% grammar compliance

### Phase 0: Compliance harness first

1. Add a generated "grammar rule coverage" test:
   - compare parser listener `Enter*` set against implemented modelBuilder handlers.
   - classify each rule as `implemented`, `intentionally-ignored`, or `missing`.
2. Add parse-tree vs model extraction assertions for representative files in `validationdata/`.
3. Fail CI when new grammar rules are unclassified.

### Phase 1: Close model representation gaps

1. Add missing model types:
   - `Alias`, `Metadata`, `Rendering`, `Message`.
2. Add missing relationship types:
   - satisfy/verify relationships as first-class model nodes or typed edges.
3. Add model structs for package filters and prefix metadata annotations/usages.

### Phase 2: Fill missing listener mappings (highest value first)

1. Implement missing top-level/member handlers:
   - alias/filter/import variants.
2. Implement missing behavior/occurrence handlers:
   - perform/exhibit/satisfy/message/individual/portion/event usage.
3. Implement rendering/metadata definition+usage extraction.

### Phase 3: Stop semantic loss in existing handlers

1. Replace raw `GetText()` extraction with structured field extraction.
2. Parse full cardinalities/lists where grammar allows many entries:
   - dependency clients/suppliers
   - specialization lists
   - import variants
3. Set all existing unresolved/type fields during parse so `ResolveReferences()` can do real work.

### Phase 4: Resolver parity

1. Add `resolveDependencyRefs`, `resolveCommentRefs`, and resolvers for new element types.
2. Ensure bidirectional links are set consistently for satisfy/verify/include/allocate relationships.
3. Add unresolved-reference diagnostics by element kind.

### Phase 5: Definition of done for "100%"

1. Every listener `Enter*` rule is either:
   - mapped to model semantics, or
   - explicitly documented as lexical/syntactic-only with tests.
2. Every model unresolved/type field is populated by parser where grammar provides input.
3. Validation corpus passes with:
   - parse success,
   - no silent drops for supported rules,
   - traceable coverage report.

## Phase 5 status (2026-02-19)

Phase 5 gates are now enforced in tests:

1. Listener rule classification gate:
   - `sysml/grammar_compliance_phase0_test.go`
   - Fails on any unclassified parser `Enter*` rule.
   - Requires each rule to be either implemented in `modelBuilder` or listed in `sysml/testdata/grammar_coverage_classification.json`.

2. No-silent-drop gate for supported rules:
   - `sysml/phase5_definition_of_done_test.go` (`TestPhase5RepresentativeSupportedRulesNotDropped`)
   - Parses representative standard examples and asserts extraction for supported constructs:
     - alias
     - metadata/filter
     - message
     - individual/snapshot/timeslice occurrences
     - satisfy/verify relationships
     - rendering

3. Unresolved/type-field population gate:
   - `sysml/phase5_definition_of_done_test.go` (`TestPhase5ReferenceAndTypeFieldsPopulatedFromGrammar`)
   - Asserts parser populates extracted unresolved/type fields from grammar text for supported constructs:
     - alias targets
     - message sender/receiver endpoints
     - satisfy/verify unresolved references
     - rendering usage type refs
     - filter expressions

4. Traceable coverage artifacts:
   - Rule classification file: `sysml/testdata/grammar_coverage_classification.json`
   - Gap/status report: `docs/reports/GRAMMAR_COMPLIANCE_GAP_REPORT.md`

## Open TODOs after Phase 5

Yes. There are still open items to maximize standards compliance:

1. Large intentionally-ignored surface remains:
   - `sysml/testdata/grammar_coverage_classification.json` currently lists 369 `Enter*` rules as `intentionally_ignored`.
   - This preserves test traceability but is not semantic compliance.

2. Semantic loss hotspots still exist in parser extraction:
   - Several handlers still rely on raw `GetText()` where grammar structure should be preserved (imports, transition fields, subjects/actors/objectives, requirement constraints, flow endpoints, payload text).

3. Import subform handlers are still no-op callbacks:
   - `EnterMembershipImport`, `EnterNamespaceImport`, `EnterFilterPackage` in `sysml/parse.go` are intentionally empty and rely on parent-level extraction.
   - This is functional for current behavior but not full parse-rule semantic mapping parity.

4. Current Phase 5 "no silent drop" checks are representative, not exhaustive per-rule conformance checks over all validation files and standard libraries.

## Further phases (post-Phase-5)

### Phase 6: Eliminate classification debt

1. Reduce `intentionally_ignored` rules by implementing semantic extraction for highest-value ignored rules first (relationships, behaviors, actions, structured expressions).
2. For rules that are truly syntax-only, document rationale per rule category (not only per rule name).
3. Add a ratchet test: `intentionally_ignored` count must never increase; target a monotonic decrease per release.

### Phase 7: Structured extraction parity (remove `GetText` fallbacks)

1. Replace raw-text extraction with structured field mapping for:
   - import forms and filters,
   - transition trigger/guard/source/target,
   - requirement subject/actor/objective and constraint bodies,
   - flow endpoints and payload typing.
2. Keep text snapshots only as optional diagnostics fields, not primary semantic storage.
3. Add focused tests that verify AST field values instead of string snapshots.

### Phase 8: Full unresolved/type-field population matrix

1. Build a test matrix of all model unresolved/type-reference fields and assert population paths from grammar where applicable.
2. Add per-kind unresolved diagnostics in parse results (counts by element/rule) to detect silent omissions early.
3. Add regression tests for multi-reference cardinalities (lists, multiple suppliers/clients, multiple actors/objectives, etc.).

### Phase 9: Corpus-wide compliance harness

1. Expand from representative examples to full `validationdata/` corpus with per-file expected extracted constructs.
2. Add conformance snapshots for key standard library files (`libraries/sysml.library/`) and verify extracted element inventories.
3. Produce machine-readable compliance report artifacts in CI (rule coverage + extraction coverage + unresolved coverage).
