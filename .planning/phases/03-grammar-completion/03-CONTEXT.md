# Phase 3: Grammar Completion - Context

**Gathered:** 2026-02-06
**Status:** Ready for planning

<domain>
## Phase Boundary

Complete grammar implementation to enable parsing of large enterprise model repositories using full sysml-core syntax and libraries. Critical requirement: All elements must land in the model — nothing discarded during parsing.

This phase focuses on implementing parser handlers for the 90% of grammar constructs that are currently parsed by ANTLR but discarded by the model builder. Scope includes SysML v2 textual syntax and KerML elements required for SysML v2 file imports.

Out of scope: Graphical notation, semantic validation, code generation, editing capabilities.

</domain>

<decisions>
## Implementation Decisions

### ANTLR Policy
- **Do NOT modify files in `internal/parser/`** — Work within ANTLR-generated code constraints
- Accept go vet warnings in parser files as code generation artifacts

### Element Retention
- **Zero silent discarding** — Every parsed element must be stored in model
- All 449 parser handlers must have corresponding model builder implementations or explicit handling
- No element loss during parse-to-model conversion

### Backward Compatibility
- Maintain existing API — no breaking changes
- Add new element types as extensions
- Keep deprecated wrappers for old functions when adding generics

### Test Coverage
- Each new handler must have corresponding unit tests
- No regression: existing 54 passing validation files must continue to pass
- Target: 98%+ validation success rate (currently 96.4%)

### KerML Support
- Include KerML elements as needed to import SysML v2 files
- Support .kerml file extensions (already discovered in libraries)
- Focus on KerML constructs used by standard libraries

### Priority Ranking (Updated)

**P0 - Critical (Must Have)**
- Dependency (relationships between elements) — blocks enterprise models
- Comment/Documentation (model documentation) — essential for model understanding
- FlowDefinition/Usage (data flows) — common in system modeling

**P1 - High (Should Have)**
- ControlNodes (Fork, Join, Merge, Decision) — activity diagrams
- OccurrenceDefinition/Usage (time-based) — temporal modeling
- BindingConnectorAsUsage — feature bindings
- SuccessionAsUsage — predecessor-successor relationships
- **CaseDefinition/Usage** — use case modeling (promoted from P2)
- Other low-effort elements that fit within capacity

**P2 - Medium (Nice to Have)**
- IncludeUseCaseUsage
- PortConjugation
- FlowEnd
- (Cases moved to P1)

**P3 - Low (Future)**
- Metadata constructs
- Element filtering
- View exposure
- Textual representations

### Coverage Target
- **Claude's discretion** — Determine appropriate coverage target (80%, 90%, or other) based on complexity analysis of remaining handlers
- Balance completeness with implementation effort
- Prioritize by enterprise usage frequency

### Performance
- **Claude's discretion** — Add performance criteria if parsing speed or memory usage is at risk
- Balance completeness with parsing speed
- Use Builder preallocation patterns where beneficial

### Claude's Discretion
- Implementation order within priority levels
- Model structure design to capture BNF semantics accurately
- Helper patterns for similar element types
- Exact coverage percentage target (80-90% range)
- Performance thresholds if needed

</decisions>

<specifics>
## Specific Ideas

- Parser currently has 449 handlers available, only 46 implemented (10.2% coverage)
- Critical gap: Dependencies are essential for enterprise models but completely missing
- Comment and Doc model types exist but have no parser handlers (being discarded)
- Two validation files fail (07-Variant Configuration, 15-Properties-Values-Expressions) — acceptable to leave as edge cases for now
- KerML files discovered in standard libraries (36 .kerml files) — need support for these constructs

</specifics>

<deferred>
## Deferred Ideas

- **Graphical notation** — Only textual syntax in this phase; graphical is separate capability
- **Semantic validation** — Syntax parsing only, not semantic correctness checking
- **Code generation** — No transform to other formats (PlantUML, XMI, etc.)
- **Editing capabilities** — Read-only parsing for now; editing is future phase
- **100% validation success** — 98% acceptable, 2 edge case files can fail
- **Metadata/Filtering/View exposure** — P3 priority, defer if effort too high

</deferred>

<reference>
## Reference Materials

- BNF Grammar: docs/bnf/SysML-textual-bnf.kebnf (1705 lines)
- KerML Grammar: docs/bnf/KerML-textual-bnf.kebnf (1467 lines)
- Current Parser: gosysml2/internal/parser/ (449 handlers)
- Current Model: gosysml2/sysml/model.go (28 types)
- Gap Analysis: .planning/GRAMMAR_GAP_ANALYSIS.md
- Validation Dataset: 56 files across 18 categories

</reference>

---

*Phase: 03-grammar-completion*
*Context gathered: 2026-02-06*
*Last updated: 2026-02-06*
