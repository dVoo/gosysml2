# API Design (Current + Next Steps)

This document defines the intended high-level API shape for `gosysml2` after validating
integration feedback from `sysml2ls` and `sysml2_d2`.

Goals:

- Straightforward, human-friendly API for common parsing and tooling tasks
- Stable, non-surprising semantics (especially definition-vs-usage and typed usages)
- Maintain parser throughput and low allocation behavior for large model repositories

## Principles

1. One obvious path for common tasks.
2. Non-breaking additions first; breaking interface changes only with strong payoff.
3. Parser performance features remain opt-in and cheap (`ParseDir`, `WithDiscardTree`).
4. API docs should describe the canonical path, while advanced paths remain available.

## Canonical API Surface

### Parsing

- In-memory parse: `ParseString(text, opts...)`
- File parse: `ParseFile(path, opts...)`
- Idiomatic model/error forms: `ParseStringModel`, `ParseFileModel`
- Repository parse/streaming: `ParseDir(ctx, dir, DirOptions)`

Preferred parse-status contract:

- `result.Err() error`
- `result.Errors() []*Error`
- `result.ParseError *ParseError` (aggregate details)

### Traversal

Recommended usage order:

1. `All`, `OfType[T]`, `OfKind` for most query/filter loops
2. `Walk` when depth and early-exit are needed
3. `Visit` for large structured processors

`WalkAll` is a convenience helper for no-early-exit traversal.

### Semantic Queries (Definitions/Usages)

Canonical semantic classification API:

- `elem.Role() ElementRole`
- `usage.TypeName() string` for typed-usage names even when unresolved
- `RoleOf`, `IsDefinitionElement`, `IsUsageElement`, `UsageTypeName` remain convenience wrappers

Rationale:

- Avoids reflection on concrete `IsDefinition` fields
- Preserves compatibility (no new methods on public interfaces yet)
- Gives tooling one stable semantic API to build on

### Position-based Lookup

- `ElementAt(model, line, col)` for named elements
- `ElementAtIncludingUnnamed(model, line, col)` when anonymous usages matter

This split keeps existing behavior stable while allowing richer editor features.

## Parser/Runtime Concerns (Not Primarily API Shape)

These are important but should be handled as implementation work, not API redesign:

- Incremental/partial re-parse
- Rewrite position mapping for compatibility rewrites
- Documentation propagation (`doc /* ... */`) semantics
- Error-recovery guarantees/documentation

## Gap Review Summary (Both Tooling Notes)

### `sysml2ls`-reported gaps

- Many were already resolved in current `gosysml2` (parse options, parse result status,
  model storage, `ParseDir`, rewrite visibility, etc.).
- Remaining API-shape issues are mainly semantic ergonomics:
  definition-vs-usage classification and unresolved type-name access.
- Remaining parser-behavior issue: element documentation population.

### `sysml2_d2`-reported gaps

- Most parser gaps (1–28) are already closed.
- Open gaps 29–31 are explicitly non-SysML-v2-compliant syntax and should remain rejected.

This means the current redesign focus should stay on API clarity and editor/tooling
ergonomics rather than broad grammar changes.

## Performance Notes

- Prefer `ParseDir` + `DirOptions{Workers: N}` for repo-scale parsing.
- Use `WithDiscardTree()` when parse-tree inspection is unnecessary.
- Keep semantic helpers allocation-free and O(1) (type switches, simple field reads).
- Avoid post-parse full-model scans in parser listeners when a local listener-time action
  can preserve performance (for example, doc attachment heuristics).

## Possible Future Breaking Cleanup (Major Version)

If a major-version API cleanup is acceptable, consider:

- Narrowing `Usage` into typed/untyped usage interfaces
- Reducing traversal API surface if real-world usage converges on iterators + `Walk`

Even with the clean method-based API, the helper functions remain useful convenience wrappers.
