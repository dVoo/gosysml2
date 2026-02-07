---
phase: quick
task: 002
type: execute
files_modified: [
  "go.mod",
  "go.sum",
  "README.md",
  "docs/README.md",
  "docs/USAGE.md",
  "docs/PERFORMANCE.md",
  "docs/usage-guide.md",
  "low/doc.go",
  "sysml/doc.go",
  "sysml/errors.go",
  "sysml/parse.go",
  "low/lexer.go",
  "low/parser.go",
  "cmd/check_validation.go",
  "examples/main.go",
  "examples/basic/main.go",
  "examples/parallel/main.go",
  "examples/requirements/main.go",
  "examples/validation/main.go",
  "examples/visitor/main.go",
  "check_testdata.go",
  ".planning/codebase/STACK.md",
  ".planning/codebase/INTEGRATIONS.md"
]
autonomous: true
---

<objective>
Rename the library from gosysml2_oc to gosysml2.

Purpose: Use the cleaner library name gosysml2 instead of gosysml2_oc
Output: All imports and references updated to github.com/dVoo/gosysml2
</objective>

<execution_context>
@/home/daniel/.config/opencode/get-shit-done/workflows/execute-plan.md
</execution_context>

<tasks>

<task type="auto">
  <name>Update go.mod module name</name>
  <files>go.mod</files>
  <action>
Update the module name in go.mod from github.com/dVoo/gosysml2 to github.com/dVoo/gosysml2.

Run go mod tidy to update dependencies.
  </action>
  <verify>go.mod shows module github.com/dVoo/gosysml2</verify>
  <done>Module name updated to gosysml2</done>
</task>

<task type="auto">
  <name>Update all Go import paths</name>
  <files>
    low/doc.go
    sysml/doc.go
    sysml/errors.go
    sysml/parse.go
    low/lexer.go
    low/parser.go
    cmd/check_validation.go
    examples/main.go
    examples/basic/main.go
    examples/parallel/main.go
    examples/requirements/main.go
    examples/validation/main.go
    examples/visitor/main.go
    check_testdata.go
  </files>
  <action>
Replace all import paths from github.com/dVoo/gosysml2 to github.com/dVoo/gosysml2 in all Go source files.

Use find and sed to update all files:
- Find all .go files
- Replace gosysml2_oc with gosysml2 in import statements
  </action>
  <verify>No files contain gosysml2_oc in imports</verify>
  <done>All Go imports updated to gosysml2</done>
</task>

<task type="auto">
  <name>Update documentation references</name>
  <files>
    README.md
    docs/README.md
    docs/USAGE.md
    docs/PERFORMANCE.md
    docs/usage-guide.md
    .planning/codebase/STACK.md
    .planning/codebase/INTEGRATIONS.md
  </files>
  <action>
Update all documentation to reference gosysml2 instead of gosysml2_oc:

1. Update README.md - change go get command and import examples
2. Update docs/README.md - change all import paths
3. Update docs/USAGE.md - change import paths
4. Update docs/PERFORMANCE.md - change import paths if any
5. Update docs/usage-guide.md - change references
6. Update .planning/codebase/STACK.md - update module name
7. Update .planning/codebase/INTEGRATIONS.md - update references

Use sed to replace gosysml2_oc with gosysml2 in all markdown files.
  </action>
  <verify>No documentation files contain gosysml2_oc references</verify>
  <done>All documentation updated to gosysml2</done>
</task>

<task type="auto">
  <name>Verify build and commit</name>
  <files>all updated files</files>
  <action>
1. Run go mod tidy to ensure module is clean
2. Build all packages to verify imports work: go build ./...
3. Stage all changes
4. Commit with descriptive message
  </action>
  <verify>go build ./... succeeds with no errors</verify>
  <done>Library successfully renamed to gosysml2</done>
</task>

</tasks>

<verification>
- go.mod shows module github.com/dVoo/gosysml2
- No Go files contain gosysml2_oc in imports
- No documentation contains gosysml2_oc references
- go build ./... succeeds
- All changes committed
</verification>

<success_criteria>
- Module name is github.com/dVoo/gosysml2
- All imports use github.com/dVoo/gosysml2
- All documentation references gosysml2
- Project builds successfully
</success_criteria>

<output>
After completion, create `.planning/quick/002-rename-library-to-gosysml2/002-SUMMARY.md`
</output>
