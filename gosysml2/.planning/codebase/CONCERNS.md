# Codebase Concerns

**Analysis Date:** 2026-02-05

## Tech Debt

**Memory Usage in Large Repositories:**
- Issue: Current implementation loads entire parse tree into memory. PERFORMANCE.md documents that input files require 30-50x multiplier in RAM (e.g., 100MB file = 3-5GB RAM).
- Files: `low/parser.go`, `sysml/parse.go`
- Impact: Cannot parse repositories larger than available RAM. Limits production use to small-to-medium models. Very large enterprise SysML repositories (>1GB) are infeasible.
- Fix approach: Implement Phase 2-4 recommendations from PERFORMANCE.md:
  1. Remove `tokens.Fill()` in `low/parser.go:53` (already done - see comment on line 53)
  2. Implement streaming parser that processes files without building full in-memory model
  3. Add optional parse tree discarding (already implemented via `WithDiscardTree()` option)
  4. Use byte streams throughout to avoid string copies

**Unresolved Reference Handling:**
- Issue: References may remain unresolved if targets cannot be found. No explicit tracking of resolution failures for debugging.
- Files: `sysml/model.go` (lines 1960-1993 in `findElement`), `sysml/parse.go` (lines 1565-1597 in reference capture)
- Impact: Silent failures - unresolved references are indistinguishable from intentional lazy references. Users cannot easily identify broken links in models.
- Fix approach: Add resolution statistics and warning collection to Model. Provide method to list all unresolved references with context.

**String-Based Reference Names:**
- Issue: References stored as strings (`unresolvedDerivedFrom`, `unresolvedSatisfiedBy`, etc. in `sysml/model.go` lines 572-575) without validation against actual element names.
- Files: `sysml/model.go` (Reference type definition and resolution code), `sysml/parse.go` (reference extraction)
- Impact: Typos in reference names silently fail to resolve. No compile-time type checking for cross-references.
- Fix approach: Consider using type-safe reference IDs instead of strings where possible. Add validation pass after parsing to report unresolvable references.

**Scattered Element Type Handling:**
- Issue: modelBuilder in `sysml/parse.go` contains ~1732 lines with extensive switch statements for different element types (lines 344-354, 390-399, 868-889, etc.).
- Files: `sysml/parse.go` (modelBuilder type with ~100 Enter/Exit methods)
- Impact: High risk of missing handlers when new element types are added. Difficult to maintain consistency across all element kinds.
- Fix approach: Extract element creation into a factory/visitor pattern to centralize logic by element type.

**Model Building Assumes Single Pass:**
- Issue: Model building from parse tree happens in single traversal with immediate parent/child tracking. No validation that tree structure matches expected SysML hierarchy.
- Files: `sysml/parse.go` (lines 277-1562 in modelBuilder)
- Impact: Malformed parse trees may produce incorrect models without errors. Scope and nesting rules not enforced.
- Fix approach: Add structural validation pass after model construction to verify valid nesting rules.

## Known Bugs

**Package Stack Initialization Gap:**
- Bug: In `EnterPackage_()` (line 304), current package is pushed to stack AFTER being set, but this can cause issues if nested packages appear immediately.
- Files: `sysml/parse.go` line 304-306
- Symptoms: First child of nested package may have incorrect parent reference
- Trigger: Parse model with nested packages
- Workaround: Ensure top-level package is fully closed before starting nested package

**Reference Type Checking Not Type-Safe:**
- Bug: When resolving references in `sysml/model.go`, type assertions check `ok` flag but continue processing either way in some paths (lines 1806-1809, 1833-1837, etc.).
- Files: `sysml/model.go` (all resolve*Refs methods)
- Symptoms: Type mismatches silently ignored; reference remains unresolved without logging
- Trigger: Reference to wrong element type (e.g., referencing Part when expecting Requirement)
- Workaround: Carefully verify element names match expected types

**Panic on Parse Errors in MustParse Functions:**
- Bug: `MustParseString()` and `MustParseFile()` panic with `ParseError` value (not `error` interface).
- Files: `sysml/parse.go` lines 1690, 1699
- Symptoms: Panic recovered code receives `*ParseError` type, not standard error. Stack trace less useful.
- Trigger: Call MustParseString on invalid SysML with deferred recovery
- Workaround: Convert to `error` before panic or use regular Parse functions with error checks

## Security Considerations

**No Input Size Limits:**
- Risk: Parser will attempt to parse arbitrarily large inputs, enabling DoS via memory exhaustion.
- Files: `low/parser.go`, `sysml/parse.go` (all Parse* functions)
- Current mitigation: None - relies on caller to validate input size
- Recommendations:
  1. Add optional max size parameter to all Parse functions
  2. Enforce per-token limits in lexer
  3. Document memory requirements in API docs

**No Path Validation in File Operations:**
- Risk: `ParseFile()` accepts any path without validation. Could read sensitive files if called with untrusted paths.
- Files: `sysml/parse.go` lines 50-65, `sysml/parse.go` lines 149-165 (ParseDirectory)
- Current mitigation: OS-level permissions
- Recommendations:
  1. Validate paths are within expected directory
  2. Add option to restrict parsing to .sysml files only
  3. Document security implications in README

**ANTLR4 Dependency Risk:**
- Risk: Generated parser (`internal/parser/sysmlv2_parser.go`, etc.) has TODOs indicating incomplete implementation (lines 21, 1085).
- Files: `internal/parser/sysmlv2_lexer.go` (TODOs at lines 21, 1085)
- Current mitigation: None
- Recommendations:
  1. Review and complete TODO items in generated code
  2. Keep antlr4-go updated regularly
  3. Monitor for security advisories

## Performance Bottlenecks

**Parallel Parsing Limited by Goroutines:**
- Problem: `ParseDirectoryParallel()` creates goroutine per file with semaphore limiting concurrency.
- Files: `sysml/parse.go` lines 171-210
- Cause: Unbounded goroutine creation if thousands of files exist before semaphore acquired
- Improvement path:
  1. Pre-allocate file list (already done)
  2. Create bounded work queue
  3. Monitor memory per worker and adjust load

**Reference Resolution O(n²) Behavior:**
- Problem: `ResolveReferences()` calls `findElement()` for every unresolved reference. Each `findElement()` walks element index and parent chain.
- Files: `sysml/model.go` lines 1960-1993, resolution calls throughout
- Cause: No caching of lookups; parent chain walks repeated for similar references
- Improvement path:
  1. Build scope map once during BuildIndex()
  2. Cache common lookups (qualified names)
  3. Use trie for partial name matching

**String Concatenation in QualifiedName:**
- Problem: `QualifiedName()` rebuilds full path by concatenating parent names on every call.
- Files: All element types inherit baseElement
- Cause: No caching of computed qualified names
- Improvement path:
  1. Cache qualified name after first computation
  2. Invalidate on parent changes
  3. Consider immutable model after parsing complete

**Element Index Building Not Optimized:**
- Problem: `BuildIndex()` walks all elements and builds map without pre-sizing.
- Files: `sysml/model.go` (BuildIndex method)
- Cause: Map grows and rehashes repeatedly; unknown number of elements upfront
- Improvement path:
  1. Do initial count pass to pre-allocate map capacity
  2. Use more efficient name keys (avoid string concatenation)

## Fragile Areas

**modelBuilder State Machine:**
- Files: `sysml/parse.go` (lines 257-1597)
- Why fragile: Complex element stack management with push/pop in Enter/Exit methods. Missing Exit handler breaks parent chain for nested elements. Package stack separate from element stack (line 304-305) adds complexity.
- Safe modification:
  1. Add unit tests for each Enter/Exit pair
  2. Use stack helper methods (Push/Pop/Peek) instead of slice operations
  3. Create stack validation function to check invariants
- Test coverage: Limited - only `parse_test.go` has basic happy-path tests

**Reference Resolution Logic:**
- Files: `sysml/model.go` (lines ~1700-1994)
- Why fragile: Each element type has custom resolution method. Changes to element structure or new element types require new resolution methods. No common pattern or base implementation.
- Safe modification:
  1. Extract resolution into central switch statement
  2. Document resolution algorithm
  3. Add regression tests for cross-file references
- Test coverage: `parse_test.go` tests basic cases; no tests for complex reference chains

**Visitor Pattern Implementation:**
- Files: `sysml/visitor.go`, test in `visitor_test.go`
- Why fragile: BaseVisitor struct embeds in custom visitors. Missing Visit method for new element types silently skips them.
- Safe modification:
  1. Use interface-based dispatch instead of embedding
  2. Require explicit method for each type
  3. Add compile-time checks for completeness
- Test coverage: `visitor_test.go` covers basic cases but not all element types

## Scaling Limits

**Single-File Parsing Limited by RAM:**
- Current capacity: ~100MB file on 32GB machine (with WithDiscardTree option)
- Limit: Hits RAM ceiling; GC becomes bottleneck
- Scaling path:
  1. Implement streaming lexer (Phase 2.2 in PERFORMANCE.md)
  2. Process tokens incrementally without buffering
  3. Add streaming visitor interface

**Directory Parsing Limited by File Count:**
- Current capacity: ~1000s of files (depends on file size and RAM)
- Limit: All results held in memory before return
- Scaling path:
  1. Use streaming callback pattern (already implemented in ParseDirectoryStream)
  2. Provide cursor-based pagination for repositories
  3. Add incremental indexing

**Model Element Graph Traversal:**
- Current capacity: ~100,000 elements comfortably
- Limit: Reference resolution becomes slow (O(n²)); qualified name computation slow
- Scaling path:
  1. Implement reference caching
  2. Build scope tree instead of flat index
  3. Use lazy evaluation for computed properties

## Dependencies at Risk

**antlr4-go v4.13.1:**
- Risk: Generated parser files have incomplete implementations (TODO comments). Parser may not handle all valid SysML grammar.
- Impact: Some valid SysML syntax may fail to parse
- Migration plan:
  1. Monitor upstream ANTLR4 Go runtime releases
  2. Regenerate parsers from latest grammar files
  3. Test against full SysML v2 test suite

**No Version Constraints on Dependencies:**
- Risk: `go.get` may pull incompatible future versions of antlr4-go
- Impact: Builds may break unexpectedly
- Migration plan:
  1. Add specific version constraints in go.mod
  2. Implement pre-commit test of dependency versions
  3. Use go mod tidy with -compat flag

## Missing Critical Features

**No Schema/Validation Against SysML Specification:**
- Problem: Parser accepts any syntactically valid input without semantic validation (e.g., illegal element nesting, conflicting modifiers).
- Blocks: Comprehensive model validation, EMF schema generation, XSD schema export
- Solution: Implement semantic validation phase after parsing; compare against spec rule set

**No Support for Comment/Documentation Extraction:**
- Problem: Documentation blocks are parsed but not captured in model. Comments lost.
- Blocks: Documentation generation, traceability reporting
- Solution: Extend parse tree visitor to capture comments; associate with elements

**No Serialization (Write) Support:**
- Problem: Can only read SysML, cannot write models back. Model modifications cannot be persisted.
- Blocks: Model transformation, refactoring, round-trip testing
- Solution: Implement Model.Write() methods with proper formatting

**No Cross-File Reference Support:**
- Problem: References only resolve within single model. Imports not processed.
- Blocks: Modular large-scale models, library management
- Solution: Implement import resolution; build cross-file index

## Test Coverage Gaps

**Complex Reference Resolution Chains:**
- What's not tested: Transitive references (A->B->C), circular references, cross-package references
- Files: `sysml/model.go` (resolution code)
- Risk: Undetected failures in complex model structures used by enterprise customers
- Priority: High

**Error Recovery and Partial Models:**
- What's not tested: Behavior when parse errors occur (does Model construction stop completely or continue with valid parts?)
- Files: `sysml/parse.go` (buildModel function)
- Risk: Unclear how to handle partially valid input
- Priority: High

**Nested Element Handling:**
- What's not tested: Deeply nested packages, parts containing parts, verification cases with actions
- Files: `sysml/parse.go` (modelBuilder stack operations)
- Risk: Scope/parent tracking bugs in real-world deep nesting
- Priority: Medium

**Parallel Parsing Race Conditions:**
- What's not tested: Concurrent access to shared resources in ParseDirectoryParallel
- Files: `sysml/parse.go` (ParseDirectoryParallel function)
- Risk: Subtle race conditions that appear only under load
- Priority: Medium

**Memory Behavior with Large Files:**
- What's not tested: Actual memory usage with files >100MB, garbage collection under stress
- Files: `low/parser.go`, `sysml/parse.go`
- Risk: Out-of-memory panics in production; no graceful degradation
- Priority: High

**Edge Cases in Name Extraction:**
- What's not tested: Qualified names with special characters, unicode identifiers, names with whitespace
- Files: `sysml/parse.go` (extractName function lines 1648-1660)
- Risk: Names silently truncated or corrupted
- Priority: Low

---

*Concerns audit: 2026-02-05*
