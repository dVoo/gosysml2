# Roadmap

**Project:** SysML v2 Parser
**Current Milestone:** v0.1 - Modern Go Implementation

## Phases

### Phase 1: Modernize Go codebase
**Goal:** Update the codebase to the newest Golang version with latest features, optimized for highest performance

**Status:** ✓ Complete - 4/4 plans finished

**Plans:** 4 plans

Plans:
- [x] 01-01-PLAN.md -- Go version upgrade + low-level wrapper modernization
- [x] 01-02-PLAN.md -- Model types + visitor generics and iter.Seq iterators
- [x] 01-03-PLAN.md -- Parse performance optimization + error handling
- [x] 01-04-PLAN.md -- Benchmarks, integration tests, and baseline recording

**Success Criteria:**
- Code uses latest Go language features
- Performance optimized for parsing speed
- Modern Go idioms and patterns applied
- Benchmarks show measurable improvements

### Phase 2: SysML Standard Libraries Support
**Goal:** Allow the usage of SysML standard libraries and validate against comprehensive validation dataset

**Status:** ✓ Complete - 3/3 plans finished

**Depends on:** Phase 1
**Plans:** 3 plans

Plans:
- [x] 02-01-PLAN.md -- Library resolution foundation (registry, discovery, indexing)
- [x] 02-02-PLAN.md -- Import resolution integration (parse pipeline, qualified names)
- [x] 02-03-PLAN.md -- Validation test suite (18 categories, standalone checker)

**Details:**
Enable the parser to resolve and use SysML standard library definitions found in `./libraries/*`. Validate implementation against all files in `./validationdata` (renamed from testdata), which contains 18 categories of SysML validation cases including Parts Tree, Function-based Behavior, State-based Behavior, Requirements, Verification, and more.

**Wave Structure:**
- Wave 1: 02-01 (independent - library foundation)
- Wave 2: 02-02 (depends on 02-01 - integration)
- Wave 3: 02-03 (depends on 02-02 - validation)

### Phase 3: Grammar Completion
**Goal:** Implement full SysML v2 grammar coverage to enable parsing of large enterprise model repositories using complete sysml-core syntax and libraries

**Status:** ✓ Complete - 2/2 plans finished

**Depends on:** Phase 2
**Plans:** 2 plans

Plans:
- [x] 03-01-PLAN.md -- P0 Critical Elements (Dependency, Comment, Documentation, Flows)
- [x] 03-02-PLAN.md -- P1 High Priority (ControlNodes, Occurrences, Bindings, Successions)

**Details:**
Complete grammar implementation to achieve 80%+ parser handler coverage (currently 10%). Critical requirement: all parsed elements must land in the model — nothing discarded. Focus on elements required for enterprise models: dependencies, documentation, flows, control nodes, occurrences.

**Success Criteria:**
- 80%+ grammar coverage (up from 10%)
- 98%+ validation success rate
- Zero elements discarded during parsing
- All new handlers have unit tests

### Phase 4: Advanced Features
**Goal:** Implement P2 medium priority grammar elements for advanced SysML modeling including case modeling, use case relationships, and port conjugation

**Status:** ○ Planned - 3 plans ready

**Depends on:** Phase 3
**Plans:** 3 plans

Plans:
- [ ] 04-01-PLAN.md -- Case modeling (CaseDefinition, CaseUsage)
- [ ] 04-02-PLAN.md -- Use case relationships and Port conjugation (IncludeUseCaseUsage, ConjugatedPortDefinition)
- [ ] 04-03-PLAN.md -- SuccessionFlow usage and Phase 4 completion

**Wave Structure:**
- Wave 1: 04-01, 04-02 (independent - both add new model types)
- Wave 2: 04-03 (depends on 04-01, 04-02 - integration and validation)

**Details:**
Complete P2 grammar elements to reach 85%+ coverage. Focus on enterprise modeling features: case definitions and usages for use case modeling, include/extend relationships between use cases, and port conjugation for interface modeling. These elements are essential for full SysML v2 compliance in enterprise contexts.

**Success Criteria:**
- 85%+ grammar coverage (up from 68%)
- 98%+ validation success rate maintained
- All P2 elements have parser handlers
- Unit tests for all new grammar elements

---

*Created: 2026-02-03*
*Updated: 2026-02-06*
