# Phase 3: Grammar Completion - Context

## User Vision

**Goal**: Enable parsing of large enterprise model repositories using full sysml-core syntax and libraries.

**Critical Requirement**: All elements must land in the model - nothing discarded during parsing.

**Current Pain Point**: 90% of grammar constructs are parsed but discarded (449 parser handlers, only 46 implemented in model builder).

## Locked Decisions (Must Honor)

1. **ANTLR-generated code policy**: Do NOT modify files in `internal/parser/` - work within constraints
2. **Element retention**: Every parsed element must be stored in model (no silent discarding)
3. **Backward compatibility**: Maintain existing API, add new element types as extensions
4. **Test coverage**: New handlers must have corresponding tests

## Claude's Discretion (Implementation Freedom)

1. **Implementation order**: Prioritize by enterprise usage frequency
2. **Model structure**: Design model types to capture BNF semantics accurately
3. **Helper patterns**: Create reusable patterns for similar element types
4. **Performance**: Balance completeness with parsing speed

## Deferred Ideas (Out of Scope)

1. **Graphical notation**: Only textual syntax in this phase
2. **Semantic validation**: Syntax parsing only, not semantic correctness
3. **Code generation**: No transform to other formats
4. **Editing capabilities**: Read-only parsing for now

## Priority Ranking (Based on Enterprise Usage)

### P0 - Critical (Blocks Most Models)
- Dependency (relationships between elements)
- Comment/Documentation (model documentation)
- FlowDefinition/Usage (data flows)

### P1 - High (Common in Enterprise)
- ControlNodes (Fork, Join, Merge, Decision)
- OccurrenceDefinition/Usage (time-based)
- BindingConnectorAsUsage
- SuccessionAsUsage

### P2 - Medium (Regular Usage)
- CaseDefinition/Usage
- IncludeUseCaseUsage
- PortConjugation
- FlowEnd

### P3 - Low (Advanced Features)
- Metadata constructs
- Element filtering
- View exposure
- Textual representations

## Success Criteria

1. **Coverage**: Increase from 10% to 80%+ parser handler coverage
2. **Validation**: 98%+ success rate on validation dataset
3. **Element retention**: Zero silent discarding of parsed elements
4. **Test coverage**: Each new handler has unit tests

## Reference Materials

- BNF Grammar: docs/bnf/SysML-textual-bnf.kebnf (1705 lines)
- KerML Grammar: docs/bnf/KerML-textual-bnf.kebnf (1467 lines)
- Current Parser: gosysml2/internal/parser/ (449 handlers)
- Current Model: gosysml2/sysml/model.go (28 types)
- Gap Analysis: .planning/GRAMMAR_GAP_ANALYSIS.md
