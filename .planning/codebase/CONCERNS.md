# Codebase Concerns

**Analysis Date:** 2026-01-30

## Tech Debt

**Memory Efficiency - Critical Issue:**
- Issue: Parser loads entire input into memory before processing. ANTLR tokens are fully materialized via `tokens.Fill()` (now removed but still affects design). Parse trees are retained in memory indefinitely for large files.
- Files: `gosysml2/low/parser.go:53`, `gosysml2/sysml/parse.go:64-95`
- Impact: Cannot parse repositories >500MB-1GB on typical systems. Each file multiplied by 30-50x memory factor. Parallel parsing amplifies memory pressure.
- Fix approach:
  1. Short term (implemented): Remove `tokens.Fill()` - saves 30-50% memory
  2. Short term (available): Use `WithDiscardTree()` option - saves additional 20-40% memory
  3. Medium term: Implement streaming handler (`ParseDirectoryStream`) - processes files sequentially with GC between files
  4. Long term: SAX-style event streaming for mega-repositories (>1GB) - see PERFORMANCE.md for detailed roadmap

**Reference Resolution - Fragile:**
- Issue: Many element types store unresolved references as string arrays that are resolved only after full model construction. No caching of resolution results. Multiple resolution passes may be needed.
- Files: `gosysml2/sysml/model.go:481-562`, `gosysml2/sysml/model.go:1435-1610`
- Impact: O(n²) reference resolution for large models. Unresolved references silently fail (no warnings). Bidirectional relationships can become inconsistent if resolution fails partially.
- Fix approach:
  1. Cache reference resolution results to avoid re-resolution
  2. Add logging for unresolved references with qualified names attempted
  3. Consider two-pass approach: collect all unresolved refs first, then validate resolvability

**ANTLR-Generated Code Duplication:**
- Issue: Parser files duplicated in two locations with identical content. ANTLR parser is 80KB, lexer 1.3KB per copy.
- Files: `code/parser/sysmlv2_*.go` (identical to `gosysml2/internal/parser/sysmlv2_*.go`)
- Impact: ~175KB of duplicated generated code increases maintenance burden. Regenerating grammar requires updating two locations.
- Fix approach:
  1. Generate once to canonical location
  2. Other modules import/vendor from canonical location
  3. Add pre-commit hook to detect duplicates

**Unused Variable Warnings in Generated Code:**
- Issue: ANTLR-generated parser has intentional TODO comments for unused variables that cause compiler warnings.
- Files: `code/parser/sysmlv2_parser.go:59893`, `gosysml2/internal/parser/sysmlv2_parser.go:59892`, `gosysml2/internal/parser/sysmlv2_lexer.go:21, 1085`
- Impact: Noise in build output. May cause issues with strict CI pipelines that treat warnings as errors.
- Fix approach: These are harmless and generated. Either suppress warnings globally for generated files, or accept as expected from ANTLR generation.

## Known Gaps

**Incomplete Element Type Coverage:**
- Issue: Model supports parts, requirements, verifications, concerns, use cases, and analysis cases. Attributes, ports, connections, allocations, items are defined but incomplete in tree builder.
- Files: `gosysml2/sysml/parse.go:265-549` (modelBuilder only handles subset of element types)
- Impact: Advanced SysML models with ports, connections, allocations will parse but elements won't be captured in model tree. Users see partial models.
- Risk: Medium - only affects specific model types not in primary use case

**Reference Type Resolution:**
- Issue: Type references for Part, Requirement, Verification, etc. use string names initially, resolved via `findElement()` but type assertions can fail silently returning nil.
- Files: `gosysml2/sysml/model.go:1531-1572`
- Impact: Usages with unresolved type references appear as separate elements with nil type. Cross-model references may not resolve.
- Safe fix: Add validation phase after resolution to report unresolved references with source locations

## Scaling Limits

**Single-File Capacity:**
- Current capacity: ~100MB per file (with defaults), ~500MB with `WithDiscardTree()`
- Limit: Hits wall at file size * 30-50 = available RAM
- On 32GB machine: Can handle ~600-1000MB input with current approach
- Scaling path: Implement Phase 2.3 (streaming handler) from PERFORMANCE.md for multi-GB files

**Repository Parsing:**
- Current capacity: Parallel mode good for <500MB total across workers
- Limit: Memory contention when parsing 10+ files concurrently
- Scaling path: Semaphore-bounded worker pool (already implemented in `ParseDirectoryParallel`) with `WithDiscardTree()` option

## Test Coverage Gaps

**Missing Tests for Reference Resolution:**
- What's not tested: Unresolved references, cross-file references, circular dependencies in derived-from relationships
- Files: `gosysml2/sysml/parse_test.go`, `gosysml2/sysml/integration_test.go`
- Risk: Reference resolution bugs won't be caught. False positives/negatives in model relationships.
- Priority: High - affects correctness of requirement traceability

**Missing Negative Test Cases:**
- What's not tested: Malformed references, qualified name lookup failures, type assertion failures during resolution
- Files: `gosysml2/sysml/parse_test.go`
- Risk: Silent failures become production issues. Users won't know if references couldn't be resolved.
- Priority: Medium

**Missing Stress Tests:**
- What's not tested: Large model parsing (>100MB), deeply nested structures (>100 levels), wide structures (1000+ siblings)
- Files: None exist
- Risk: Performance degradation and memory leaks only discovered in production
- Priority: Medium

## Error Handling Issues

**Silent Reference Resolution Failures:**
- Problem: When `m.findElement()` returns nil, element silently remains unresolved. No error collected.
- Files: `gosysml2/sysml/model.go:1457-1610`
- Current behavior: Unresolved references simply don't appear in relationship arrays (e.g., `DerivedFrom` remains empty)
- Recommendation: Add optional strict mode that tracks unresolved references and reports them

**Partial Error Collection:**
- Problem: Parser stops after first error batch; doesn't collect all syntax errors in file
- Files: `gosysml2/low/errors.go` (ErrorCollector), `gosysml2/sysml/errors.go`
- Current: Up to 10 errors shown, rest hidden with "... and N more errors"
- Impact: Users must fix errors iteratively rather than seeing all issues at once

## Fragile Areas

**Model Building from Parse Tree:**
- Files: `gosysml2/sysml/parse.go:257-564` (modelBuilder walker)
- Why fragile: Element type checking via pattern matching on parse tree contexts. Grammar changes require parser regeneration + visitor updates. No validation that all expected contexts are handled.
- Safe modification: Add exhaustive switch statements with compile-time checks for new grammar rules
- Test coverage: Integration tests check main element types but not all combinations

**Reference Lookup by Name:**
- Files: `gosysml2/sysml/model.go:1610-1650` (findElement method not shown but referenced extensively)
- Why fragile: Resolves by qualified name comparison. No namespace handling. Global search inefficient.
- Safe modification: Maintain namespace hierarchy; use indexed lookup within each scope level
- Test coverage: No tests for complex qualified names or namespace resolution

**Visitor Pattern with Type-Specific Methods:**
- Files: `gosysml2/sysml/visitor.go:64-96` (visitElement switch statement)
- Why fragile: 11 type-specific visit methods + catch-all. Adding new element types requires updating visitor interface + all implementations
- Safe modification: Use reflection-based visitor or add new types to catch-all handler first, convert to specific when stable
- Test coverage: Base tests only; no tests for adding custom visitor implementations

## Dependencies at Risk

**ANTLR v4 Lock-in:**
- Risk: Grammar is tied to ANTLR4 specifically. Major version update would require regeneration.
- Impact: v4.13.1 is current; no newer versions available for Go. Stuck on this version.
- Alternatives: Consider parsing combinator library (Parsec) for future, but high migration cost

**Golang.org/x/exp Dependency:**
- Risk: `golang.org/x/exp` is unstable/experimental package used for mmap file reading in PERFORMANCE.md suggestions
- Impact: Not in go.mod yet, but referenced in performance docs. Using this for production would tie to unstable API.
- Migration: Replace with stable stdlib or find alternative when performance optimizations needed

## Security Considerations

**Input Validation:**
- Risk: No validation on input size before parsing. Large allocations possible from malformed input.
- Files: `gosysml2/sysml/parse.go:45-115` (ParseString, ParseFile, ParseBytes)
- Current mitigation: ANTLR parser has internal limits but they're not explicitly set
- Recommendations:
  1. Add max input size check before parsing
  2. Set ANTLR parser recursion limit explicitly
  3. Add timeout for parsing operations

**Memory Exhaustion:**
- Risk: Malicious input with deeply nested structures could cause stack overflow or memory exhaustion
- Files: `gosysml2/sysml/visitor.go:112-122` (recursive walkDepth function)
- Current mitigation: None; recursion depth unbounded
- Recommendations:
  1. Add depth tracking with configurable limit
  2. Consider iterative traversal for visitor pattern

**Parse Tree Access:**
- Risk: Parse tree returned directly to users via `ParseResult.Tree`. Changes to tree could corrupt internal state.
- Files: `gosysml2/sysml/parse.go:79-85` (ParseResult.Tree assignment)
- Current mitigation: `WithDiscardTree()` option discards tree
- Recommendations: Make tree private or return immutable interface; only expose through filtered accessors

## Missing Critical Features

**No validation for SysML semantic constraints:**
- Problem: Parser validates syntax only. No enforcement of SysML rules like "requirement ID must be unique", "definition before usage"
- Impact: Invalid but parseable models silently accepted
- Solution: Add semantic validator phase after model building

**No import resolution:**
- Problem: `Import` elements parsed but not used. No mechanism to load imported packages.
- Impact: Multi-file models with imports cannot resolve cross-file references
- Solution: Add import loader that can resolve packages from paths or registry

**No model serialization:**
- Problem: Can only parse SysML, not regenerate it. No way to modify and save models.
- Impact: Read-only tool; cannot be used in model generation pipelines
- Solution: Add serializer that reconstructs SysML text from model

---

*Concerns audit: 2026-01-30*
