# Codebase Concerns

**Analysis Date:** 2026-02-05

## Tech Debt

### 1. Parent Assignment Bug in Model Builder
- **Issue:** Attributes inside requirements/verifications are added to package instead of their containing element due to not checking `elementStack`
- **Files:** `gosysml2/sysml/parse.go` (lines 844-895, 897-963)
- **Impact:** Broken parent-child relationships. Can't navigate `requirement.Children()` to find its attributes
- **Fix approach:** Update `EnterAttributeDefinition()` and `EnterAttributeUsage()` to check `elementStack` before falling back to `currentPkg`
- **Reference:** See `analysis/CUSTOM_ATTRIBUTES_ANALYSIS.md` for detailed analysis

### 2. Memory Efficiency - Critical Issue
- **Issue:** Parser loads entire input into memory before processing. Parse trees retained indefinitely for large files
- **Files:** `gosysml2/low/parser.go:53`, `gosysml2/sysml/parse.go:64-95`
- **Impact:** Cannot parse repositories >500MB-1GB on typical systems. Each file multiplied by 30-50x memory factor
- **Fix approach:**
  1. ✅ Implemented: Remove `tokens.Fill()` - saves 30-50% memory
  2. ✅ Available: Use `WithDiscardTree()` option - saves additional 20-40% memory
  3. ✅ Implemented: `ParseDirectoryStream()` - processes files sequentially with GC between files
  4. Long term: SAX-style event streaming for mega-repositories (>1GB)

### 3. Missing Comment/Documentation Extraction
- **Issue:** Model types `Comment` and `Doc` exist but no listener methods implemented to extract them from parse tree
- **Files:** `gosysml2/sysml/model.go` (lines 1415-1453), `gosysml2/sysml/parse.go` (no EnterComment_/EnterDocumentation methods)
- **Impact:** Comments and documentation are parsed by ANTLR but never appear in high-level model
- **Fix approach:** Implement `EnterComment_()` and `EnterDocumentation()` methods in modelBuilder
- **Reference:** See `analysis/PARSER_COMPLETENESS_REPORT.md`

### 4. ANTLR Generated Code Contains TODOs
- **Issue:** ANTLR-generated parser has intentional TODO comments for unused variables that cause compiler warnings
- **Files:** `gosysml2/internal/parser/sysmlv2_lexer.go` (lines 21, 1085), `gosysml2/internal/parser/sysmlv2_parser.go` (line 59892)
- **Impact:** Compiler warnings about unused variables; EOF string not properly initialized
- **Fix approach:** These are in generated code - regenerate with updated ANTLR or post-process to clean up

### 5. Duplicate Parser Code in Two Locations
- **Issue:** Same ANTLR-generated parser exists in both `/code/parser/` and `/gosysml2/internal/parser/`
- **Files:** `code/parser/*.go`, `gosysml2/internal/parser/*.go`
- **Impact:** Code duplication (~175KB), maintenance burden, risk of version drift
- **Fix approach:** Consolidate to single location; `code/` appears to be standalone, `gosysml2/` is the library

### 6. Reference Resolution - Fragile
- **Issue:** Many element types store unresolved references as string arrays resolved only after full model construction. No caching of results
- **Files:** `gosysml2/sysml/model.go:481-562`, `gosysml2/sysml/model.go:1576-1993`
- **Impact:** O(n²) reference resolution for large models. Unresolved references silently fail
- **Fix approach:** Cache resolution results; add logging for unresolved references; two-pass validation

### 7. Missing Verification Elements
- **Issue:** Verification subjects and methods are not extracted (always shows "unspecified")
- **Files:** `gosysml2/sysml/parse.go` (EnterVerificationCaseDefinition, EnterVerificationCaseUsage)
- **Impact:** Incomplete verification case modeling
- **Fix approach:** Add extraction logic for subject member and objective clause

## Known Bugs

### 1. Attribute Parent Assignment Bug
- **Symptoms:** Attributes appear at package level instead of nested under their containing requirement/verification
- **Files:** `gosysml2/sysml/parse.go` lines 844-963
- **Trigger:** Parse any requirement or verification with custom attributes
- **Workaround:** Use location-based filtering to find attributes near their intended parent
- **Code pattern:**
  ```go
  // Current (buggy) - always adds to currentPkg
  if b.currentPkg != nil {
      attr.parent = b.currentPkg
      b.currentPkg.AddChild(attr)
  }
  
  // Should be - check elementStack first
  if len(b.elementStack) > 0 {
      parent := b.elementStack[len(b.elementStack)-1]
      attr.parent = parent
      parent.AddChild(attr)
  } else if b.currentPkg != nil {
      // ...
  }
  ```

### 2. ItemDefinition/ItemUsage Missing Element Stack Push
- **Symptoms:** Items cannot have nested elements because they don't push themselves to elementStack
- **Files:** `gosysml2/sysml/parse.go` lines 416-464
- **Impact:** Items defined inside items won't have correct parent; deeper nesting broken
- **Safe modification:** Add Exit methods and push/pop from elementStack

### 3. Port Direction Not Extracted
- **Symptoms:** All ports have `PortDirectionNone` regardless of declared direction (`in`, `out`, `inout`)
- **Files:** `gosysml2/sysml/parse.go` (EnterPortDefinition, EnterPortUsage)
- **Impact:** Can't determine port directionality for analysis

### 4. Silent Reference Resolution Failures
- **Problem:** When `m.findElement()` returns nil, element silently remains unresolved. No error collected
- **Files:** `gosysml2/sysml/model.go:1576-1993`
- **Current behavior:** Unresolved references simply don't appear in relationship arrays (e.g., `DerivedFrom` remains empty)
- **Recommendation:** Add optional strict mode that tracks unresolved references and reports them

## Security Considerations

### 1. Panic on Parse Error in Must* Functions
- **Risk:** `MustParseString` and `MustParseFile` panic on error, which could crash server applications
- **Files:** `gosysml2/sysml/parse.go` lines 1686-1702
- **Current mitigation:** Well-documented behavior (Must prefix convention); callers should use non-Must variants for production
- **Recommendations:** Add recovery examples in documentation; ensure production code uses error-returning variants

### 2. No Input Size Limits
- **Risk:** Parsing extremely large SysML files could cause OOM
- **Files:** `gosysml2/sysml/parse.go` (ParseString, ParseFile, ParseDirectory)
- **Current mitigation:** Streaming API (`ParseDirectoryStream`) available for memory-efficient processing
- **Recommendations:** Document size limits; consider adding context.WithTimeout support

### 3. File Path Traversal in ParseFile
- **Risk:** No validation of file path before reading
- **Files:** `gosysml2/sysml/parse.go` (ParseFile)
- **Current mitigation:** Uses standard Go file operations which respect OS permissions
- **Recommendations:** If used in server context, validate paths against allowlist

### 4. Memory Exhaustion from Malicious Input
- **Risk:** Malicious input with deeply nested structures could cause stack overflow
- **Files:** `gosysml2/sysml/visitor.go:196-206` (recursive walkDepth function)
- **Current mitigation:** None; recursion depth unbounded
- **Recommendations:** Add depth tracking with configurable limit; consider iterative traversal

## Performance Bottlenecks

### 1. Large ANTLR Parser Files
- **Problem:** Generated parser files are very large (>80,000 lines)
- **Files:** `gosysml2/internal/parser/sysmlv2_parser.go` (80,592 lines)
- **Impact:** Slow compilation, large binary size
- **Cause:** ANTLR generates full LL(*) parser for complete SysML v2 grammar
- **Improvement path:** Consider splitting grammar or using parser combinators for specific use cases

### 2. Reference Resolution O(n) per Lookup
- **Problem:** `Model.findElement()` does multiple index lookups and parent chain walks
- **Files:** `gosysml2/sysml/model.go` lines 1962-1993
- **Impact:** Could be slow for models with many references
- **Cause:** No caching of failed lookups; repeated string concatenation
- **Improvement path:** Add LRU cache for failed lookups; optimize qualified name construction

### 3. Single-File Capacity Limits
- **Current capacity:** ~100MB per file (with defaults), ~500MB with `WithDiscardTree()`
- **Limit:** Hits wall at file size * 30-50 = available RAM
- **On 32GB machine:** Can handle ~600-1000MB input with current approach
- **Scaling path:** Use streaming APIs for large files

## Fragile Areas

### 1. Element Stack Management
- **Files:** `gosysml2/sysml/parse.go` (multiple Enter/Exit methods)
- **Why fragile:** Complex stack operations; easy to forget Exit handler or push/pop mismatch
- **Safe modification:** Always add Exit method when adding Enter; test with deeply nested input
- **Test coverage:** Basic coverage exists but doesn't test deep nesting (>5 levels)

### 2. Type Assertion Chains in AddChild Methods
- **Files:** `gosysml2/sysml/model.go` (Package.AddChild, Part.AddChild, etc.)
- **Why fragile:** Adding new element types requires updating multiple switch statements
- **Safe modification:** Create generic container interface; use reflection or code generation
- **Test coverage:** Partial - only tested types covered

### 3. Grammar Version Coupling
- **Files:** `code/SysMLv2Parser.g4`, `code/SysMLv2Lexer.g4`
- **Why fragile:** Generated code is tied to specific grammar version; updates require regeneration
- **Safe modification:** Pin ANTLR version; document grammar version compatibility
- **Test coverage:** Grammar tests exist but not comprehensive

### 4. Model.Walk vs Visitor Pattern
- **Files:** `gosysml2/sysml/model.go` (Walk), `gosysml2/sysml/visitor.go` (Visit)
- **Why fragile:** Two traversal patterns that could diverge in behavior
- **Safe modification:** Consider unifying or clearly documenting when to use each

### 5. Visitor Pattern with Type-Specific Methods
- **Files:** `gosysml2/sysml/visitor.go:64-96` (visitElement switch statement)
- **Why fragile:** 20+ type-specific visit methods + catch-all. Adding new element types requires updating visitor interface + all implementations
- **Safe modification:** Use reflection-based visitor or add new types to catch-all handler first
- **Test coverage:** Base tests only; no tests for custom visitor implementations

## Missing Critical Features

### 1. Metadata Extraction
- **Problem:** Metadata annotations exist in grammar but not extracted to model
- **Files:** `gosysml2/sysml/model.go` (KindMetadata defined but unused)
- **Blocks:** Full SysML v2 compliance; metadata-driven tooling

### 2. Full Expression Parsing
- **Problem:** Constraint/expression bodies captured as raw text only
- **Files:** `gosysml2/sysml/parse.go` (EnterRequirementConstraintMember)
- **Blocks:** Expression analysis, constraint validation, code generation

### 3. Import Resolution
- **Problem:** Imports are parsed but not resolved to actual files
- **Files:** `gosysml2/sysml/model.go` (Import struct)
- **Blocks:** Cross-file reference resolution, multi-file model building

### 4. Library Model Loading
- **Problem:** Standard libraries (ScalarValues, etc.) not automatically available
- **Blocks:** Complete type checking for real-world models

### 5. Model Serialization
- **Problem:** Can only parse SysML, not regenerate it. No way to modify and save models
- **Impact:** Read-only tool; cannot be used in model generation pipelines
- **Solution:** Add serializer that reconstructs SysML text from model

### 6. Semantic Validation
- **Problem:** Parser validates syntax only. No enforcement of SysML rules like "requirement ID must be unique"
- **Impact:** Invalid but parseable models silently accepted
- **Solution:** Add semantic validator phase after model building

## Test Coverage Gaps

### 1. Deep Nesting Tests
- **What's not tested:** Elements nested >5 levels deep
- **Files:** Test files focus on simple structures
- **Risk:** Stack management bugs only visible in complex hierarchies
- **Priority:** Medium

### 2. Error Recovery Tests
- **What's not tested:** Parser behavior with malformed input in various contexts
- **Risk:** Panics or infinite loops on edge case inputs
- **Priority:** High

### 3. Concurrent Access Tests
- **What's not tested:** Thread safety of Model and element access
- **Risk:** Race conditions in multi-threaded use
- **Priority:** Medium (currently single-threaded design)

### 4. Large File Tests
- **What's not tested:** Files >10MB or with >10,000 elements
- **Risk:** Memory exhaustion, performance degradation
- **Priority:** Low (streaming API available)

### 5. Visitor Pattern Exhaustive Tests
- **What's not tested:** All visitor methods for all element types
- **Risk:** Visitor silently skips certain element types
- **Priority:** Low

### 6. Reference Resolution Tests
- **What's not tested:** Unresolved references, cross-file references, circular dependencies
- **Files:** `gosysml2/sysml/parse_test.go`, `gosysml2/sysml/integration_test.go`
- **Risk:** Reference resolution bugs won't be caught
- **Priority:** High - affects correctness of requirement traceability

## Dependencies at Risk

### 1. ANTLR4 Go Runtime
- **Package:** `github.com/antlr4-go/antlr/v4`
- **Risk:** Major version changes could require grammar regeneration
- **Impact:** Parser breaks, build failures
- **Mitigation:** Version pinned in go.mod; test before upgrading

### 2. ANTLR v4 Lock-in
- **Risk:** Grammar is tied to ANTLR4 specifically
- **Impact:** v4.13.1 is current; stuck on this version
- **Alternatives:** Consider parsing combinator library (Parsec) for future, but high migration cost

## Module Boundary Issues

### 1. Root Module vs gosysml2 Submodule
- **Issue:** Root module appears incomplete (no main package, just cmd tools)
- **Impact:** Confusion about which module to use as dependency
- **Recommendation:** Clarify in README that `gosysml2/` is the library module

### 2. Internal Package Usage
- **Issue:** `internal/parser` is accessible within module but not to importers
- **Impact:** External users cannot access low-level parsing details if needed
- **Recommendation:** Consider exposing through `low` package if needed

---

*Concerns audit: 2026-02-05*
