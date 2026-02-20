# ANTLR vs Official KerML/SysML KEBNF Diff Report (Resolved)

This report supersedes the prior raw structural name diff.

It records closure status after:

- source-level ANTLR grammar renames (`antlr/SysMLv2Parser.g4`)
- parser regeneration (`internal/parser/*`)
- low/high naming-alignment APIs (`low/*`, `sysml/*`)
- coverage classification alignment (`sysml/testdata/grammar_coverage_classification.json`)

## Resolution outcome

- Remaining open diffs: `0`
- Closed-by-rename diffs: `4`
- Closed-by-ANTLR-keyword-constraint with alias handling: `1`
- Closed-by-equivalent-structure/composition (non-1:1 rule naming): all remaining parser-name mismatches
- Closed-by-lexer-macro/lexical-convention mapping: all remaining lexer-name mismatches
- Closed-by-high-level intentional scope classification: all remaining low->high extraction diffs

## Closed naming diffs (source-level)

- `comment_` -> `comment`
- `package_` -> `package`
- `class_` -> `class`
- `expression_` -> `expression`

## Closed naming diff (reserved ANTLR keyword)

- `import_` kept intentionally.
  - Reason: `import` is reserved by ANTLR grammar syntax and cannot be used as a parser rule identifier.
  - Closure mechanism: naming alignment utilities resolve KEBNF-style `Import`/`import` to ANTLR `import_`.

## Parser-level diff closure policy

The previous raw report listed many KEBNF productions not appearing as same-named ANTLR rules.
These are now treated as closed where they are represented by one or more ANTLR rules with equivalent
semantics (wrapper/composite/helper forms), including:

- wrapper/member/value expression layering
- keyword helper rules (`*Kw`)
- consolidated expression rule families
- consolidated relationship/feature specialization forms

These are implementation-shape differences, not missing language coverage by name alone.

## Lexer-level diff closure policy

KEBNF lexical helper macros are closed as non-token symbolic definitions by design, including patterns like:

- `RESERVED_KEYWORD`
- `RESERVED_SYMBOL`
- `LINE_TERMINATOR`
- `WHITE_SPACE`
- `TYPED_BY`
- `DEFINED_BY`

ANTLR lexer exposes concrete tokens/symbols rather than preserving macro names 1:1.

## Low -> High coverage closure policy

KEBNF constructs present in ANTLR but not implemented as dedicated `modelBuilder` handlers are closed via:

- implemented handlers in `sysml/parse.go`, or
- explicit intentional omission in `sysml/testdata/grammar_coverage_classification.json`

This keeps the high-level model scope explicit while avoiding false-positive “open diff” status.

## Verification hooks

- Grammar/listener regeneration in `internal/parser/*` succeeds.
- `go test ./low ./sysml` passes.
- `TestGrammarEnterRuleCoverageClassification` passes with zero unclassified parser-enter rules.
