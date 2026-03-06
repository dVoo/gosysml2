# ANTLR vs KEBNF Audit Summary

Note: the referenced KEBNF sources are maintained in the upstream SysML release
repository: https://github.com/Systems-Modeling/SysML-v2-Release

This summarizes the structural diff between:

- SysML textual KEBNF from the SysML-v2-Release project
- KerML textual KEBNF from the SysML-v2-Release project
- ANTLR-generated grammar artifacts:
  - `internal/parser/sysmlv2_parser.go`
  - `internal/parser/sysmlv2_lexer.go`
- High-level extraction coverage:
  - `sysml/parse.go` (`modelBuilder` Enter* handlers)

The complete machine-generated list is in:

- `docs/reports/antlr_vs_kebnf_diff_report.md`

## Snapshot counts

- KEBNF named constructs: `557`
- KEBNF parser-like constructs: `525`
- KEBNF lexer-like constructs: `32`
- ANTLR parser rules: `448`
- ANTLR lexer symbolic tokens: `227`
- Remaining open diffs after closure classification: `0`

## Naming changes applied in this iteration

Source-level ANTLR grammar/parser naming alignment was applied for the
underscore-suffixed rules that were pure naming artifacts:

- `comment_` -> `comment`
- `package_` -> `package`
- `class_` -> `class`
- `expression_` -> `expression`

Propagation completed through:

- ANTLR generated parser/listener in `internal/parser/`
- high-level listener hooks in `sysml/parse.go`:
  - `EnterPackage` / `ExitPackage`
  - `EnterComment`
- low/high-level naming resolver tests (`low/*`, `sysml/*`)

Known exception:

- `import_` remains unchanged because `import` is reserved in ANTLR grammar
  syntax and cannot be used directly as a parser rule identifier.

## Closure model

The original raw name-diff over-reported implementation-shape mismatches.
Diffs are now considered closed when they fall into one of these buckets:

1. Source-level rename closure:
- `comment_` -> `comment`
- `package_` -> `package`
- `class_` -> `class`
- `expression_` -> `expression`

2. Reserved-keyword closure:
- `import_` retained because `import` is reserved in ANTLR grammar syntax.
- Naming resolver APIs map KEBNF-style names to concrete ANTLR names.

3. Structural-equivalence closure:
- KEBNF wrapper/nonterminal decomposition vs consolidated ANTLR helper rules.

4. Lexer-macro closure:
- KEBNF lexical convention macros are not expected as ANTLR symbolic tokens.

5. High-level scope closure:
- ANTLR constructs not handled in `modelBuilder` are either implemented or
  explicitly tracked as intentionally ignored in `sysml/testdata/grammar_coverage_classification.json`.

## Verification

- `go test ./low ./sysml` passes.
- `TestGrammarEnterRuleCoverageClassification` passes with zero unclassified enter-rules.
