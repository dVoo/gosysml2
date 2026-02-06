---
phase: quick
task: 001
type: execute
files_modified: [
  "go.mod",
  "go.sum",
  "README.md",
  "docs/USAGE.md",
  "docs/PERFORMANCE.md",
  "docs/README.md",
  "low/",
  "sysml/",
  "internal/",
  "cmd/",
  "examples/"
]
autonomous: true
---

<objective>
Reorganize project structure by moving gosysml2/ contents to root and consolidating documentation.

Purpose: Simplify project structure - all Go code at root, documentation in docs/
Output: Flattened project structure with merged README and consolidated docs
</objective>

<execution_context>
@/home/daniel/.config/opencode/get-shit-done/workflows/execute-plan.md
</execution_context>

<tasks>

<task type="auto">
  <name>Move Go module and packages to root</name>
  <files>
    go.mod
    go.sum
    low/
    sysml/
    internal/
    cmd/
    examples/
    check_testdata.go
  </files>
  <action>
Move all Go code from gosysml2/ to root:
1. Move go.mod and go.sum to root
2. Move low/ directory to root
3. Move sysml/ directory to root
4. Move internal/ directory to root
5. Move cmd/ directory to root
6. Move examples/ directory to root
7. Move check_testdata.go to root

Update go.mod module path from "github.com/dVoo/gosysml2" to "github.com/dVoo/gosysml2_oc" (matching the repo name).

After moving, run `go mod tidy` to verify the module works.
  </action>
  <verify>go build ./... succeeds from root</verify>
  <done>All Go code moved to root and builds successfully</done>
</task>

<task type="auto">
  <name>Consolidate documentation into docs/</name>
  <files>
    README.md
    docs/README.md
    docs/USAGE.md
    docs/PERFORMANCE.md
  </files>
  <action>
Consolidate documentation:

1. Create docs/README.md by merging content from gosysml2/README.md:
   - Keep the comprehensive API documentation
   - Include installation, quick start, features
   - Include API reference section
   - Update import paths to remove "gosysml2/" prefix

2. Move gosysml2/USAGE.md to docs/USAGE.md
   - Update any internal references

3. Move gosysml2/PERFORMANCE.md to docs/PERFORMANCE.md

4. Update root README.md:
   - Keep project overview and status
   - Keep repository structure section (update paths)
   - Keep getting started section
   - Add "Documentation" section linking to docs/README.md
   - Remove detailed API reference (now in docs/README.md)
   - Update all paths to remove "gosysml2/" references

Ensure no duplicate information between root README and docs/README.
Root README should be the entry point, docs/README should have full API reference.
  </action>
  <verify>All markdown files render correctly, no broken internal links</verify>
  <done>Documentation consolidated - root README is entry point, full docs in docs/</done>
</task>

<task type="auto">
  <name>Clean up and remove gosysml2/ directory</name>
  <files>
    gosysml2/ (deleted)
  </files>
  <action>
1. Verify all files have been moved from gosysml2/
2. Remove the empty gosysml2/ directory
3. Update any import paths in Go files that reference "gosysml2/" to use relative paths
4. Run `go build ./...` and `go test ./...` to verify everything works
5. Update .gitignore if needed
  </action>
  <verify>gosysml2/ directory removed, all builds and tests pass</verify>
  <done>Project successfully reorganized, gosysml2/ removed</done>
</task>

</tasks>

<verification>
- All Go code builds from root: `go build ./...`
- All tests pass: `go test ./...`
- Root README.md is clear entry point
- docs/ contains full documentation
- No duplicate information between README and docs
- No broken internal links
</verification>

<success_criteria>
- Go module at root level
- All packages (low/, sysml/, internal/, cmd/, examples/) at root
- Root README is concise entry point
- Full API documentation in docs/README.md
- USAGE.md and PERFORMANCE.md in docs/
- gosysml2/ directory removed
</success_criteria>

<output>
After completion, create `.planning/quick/001-reorganize-project-structure/001-SUMMARY.md`
</output>
