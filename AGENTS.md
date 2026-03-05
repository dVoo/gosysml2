# Repository Guidelines

## Project Structure & Module Organization
- `sysml/`: High-level typed SysML/KerML model API (main integration surface).
- `low/`: Low-level ANTLR-facing parser/lexer access.
- `internal/parser/`: Generated parser artifacts.
- `cmd/`: Utility CLIs for validation/checking repository data.
- `examples/`: Runnable usage examples (`basic`, `parallel`, `requirements`, `validation`, `visitor`).
- `docs/`: User guides, architecture notes, performance notes, and grammar references.
- `validationdata/`: Corpus used by tests and validation tooling.
- `libraries/`: Standard library model files used for import resolution.

## Build, Test, and Development Commands
- `go build ./...`: Build all packages.
- `go test ./...`: Run full unit/integration test suite.
- `go test ./sysml -run TestName`: Run a focused test.
- `go test ./... -count=1`: Disable test caching while debugging.
- `jj st`: Show working-copy status.
- `jj describe -m "scope: short summary"`: Set/update change description.

## Coding Style & Naming Conventions
- Language: Go. Format with `gofmt` before submitting changes.
- Prefer clear, exported API names; use concise comments on non-obvious logic only.
- Keep new public API consistent with current direction:
  - Parse status via `result.Err()`.
  - Directory parsing via `ParseDir(ctx, dir, DirOptions)`.
  - Top-level model access via methods like `model.Packages()`.
- File naming: lowercase with underscores only when needed; tests in `*_test.go`.

## Testing Guidelines
- Framework: Go `testing` package.
- Add/adjust tests in the same package as the code under test (`sysml/*_test.go`, `low/*_test.go`).
- Prefer table-driven tests for parser behavior and edge cases.
- For parser changes, include positive and negative cases (valid parse + diagnostic assertions).

## Commit & Pull Request Guidelines
- Keep change descriptions imperative and scoped (seen style: `parser: ...`, `refactor ...`).
- One logical change per commit/change set.
- PRs should include:
  - What changed and why.
  - API breakages/migration notes (if any).
  - Evidence of validation (`go build ./...`, `go test ./...`).
  - Updated docs/examples when public behavior changes.
