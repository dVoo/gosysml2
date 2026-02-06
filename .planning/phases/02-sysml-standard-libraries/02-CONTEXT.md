# Phase 02: SysML Standard Libraries Support - Context

**Gathered:** 2026-02-06
**Status:** Ready for planning

<domain>
## Phase Boundary

Enable the parser to resolve and use SysML standard library definitions found in `./libraries/*`. Validate implementation against all files in `./validationdata`, which contains 56 validation files across 18 categories including Parts Tree, Function-based Behavior, State-based Behavior, Requirements, Verification, and more.

This phase focuses on library resolution infrastructure, not custom user libraries or remote library fetching.

</domain>

<decisions>
## Implementation Decisions

### Library Discovery & Loading
- **Lazy loading**: Libraries are discovered and loaded on first import reference, not at parser initialization
- **Path resolution priority**:
  1. API option: `WithLibraryPath()` for programmatic control
  2. Environment variable: `SYSML_LIBRARY_PATH`
  3. Standard locations: Check for `sysml.library*` folders in project directory

### Error Handling
- **Fail fast**: Parsing stops immediately when a library element cannot be resolved
- **Helpful error messages**: Include suggestions for similar available elements when reporting unresolved elements
- **Collect all errors**: Continue parsing to find all unresolved elements, then report the complete list (don't stop at first)
- **Differentiated error types**: Distinguish between `LibraryNotFoundError` (library doesn't exist) and `ElementNotFoundError` (library exists but element doesn't)

### Import Resolution
- **Wildcard imports**: Implementation choice - balance memory usage vs lookup speed
- **Qualified names**: Implementation choice - optimize for standard SysML library structure
- **Namespace collisions**: Implementation choice - follow SysML specification semantics
- **Import aliases**: Implementation choice - support if SysML spec allows

### Claude's Discretion
- Cache strategy for loaded libraries (memory vs disk vs none)
- Library file search strategy (flat vs recursive vs index-based)
- Wildcard import expansion strategy (eager vs lazy vs index)
- Qualified name resolution implementation (direct lookup vs registry)
- Namespace collision handling (first wins vs explicit qualification vs error)
- Import alias support (if SysML spec permits)

</decisions>

<specifics>
## Specific Ideas

- Library resolution should integrate cleanly with existing parse pipeline
- Must support the standard SysML library structure found in `./libraries/sysml.library/`
- Validation should test all 56 files across 18 categories with per-category pass/fail reporting
- Consider creating a standalone validation checker tool for CI/CD integration

</specifics>

<deferred>
## Deferred Ideas

- Custom user library paths outside standard locations — could be a future enhancement phase
- Remote library fetching over HTTP — belongs in a separate "Remote Libraries" phase
- Library versioning and version resolution — future phase
- Dynamic library reloading during long-running sessions — future enhancement

</deferred>

---

*Phase: 02-sysml-standard-libraries*
*Context gathered: 2026-02-06*
