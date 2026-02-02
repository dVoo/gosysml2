# Phase 1: Modernize Go Codebase - Context

**Gathered:** 2026-02-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Update the existing SysML v2 parser codebase to the newest Golang version with latest language features, optimized for highest performance. The focus is on modernizing the wrapper code around ANTLR-generated parsers, not on parser grammar changes or new functionality.

</domain>

<decisions>
## Implementation Decisions

### Go Version and Language Features
- Target latest Go version with no backward compatibility constraints
- Apply all applicable modern features: generics, range-over-func, improved error handling, and other relevant language improvements
- Keep ANTLR-generated code as-is, modernize wrapper code only (don't modify parser output)
- Moderate refactoring approach - improve common patterns (error handling, iteration) without aggressive restructuring

### Performance Optimization
- Balanced optimization target - reasonable speed without excessive memory use (speed/memory trade-offs)
- Basic benchmarks for key operations (lexing, parsing, tree walking)
- Apply all applicable optimization techniques: algorithmic improvements, compiler optimizations, concurrency where appropriate
- Validate performance wins using real-world test files from `docs/testdata/` directory

### Claude's Discretion
- Specific modern features to apply where (generics vs range-over-func in different contexts)
- Exact optimization techniques for each component
- Profiling approach and tooling choices
- Benchmark implementation details

</decisions>

<specifics>
## Specific Ideas

- Use the 34 `.sysml` test files in `docs/testdata/` as the validation corpus for performance improvements
- Focus modernization on wrapper code that uses ANTLR output, not the generated lexer/parser files
- Measure "before and after" performance on real SysML files to validate improvements

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 01-modernize-go-codebase*
*Context gathered: 2026-02-03*
