# Technology Stack

**Analysis Date:** 2026-02-05

## Languages

**Primary:**
- Go 1.22 - Core library implementation and parser wrapper

**Secondary:**
- Java (implicit) - ANTLR4 tool runtime dependency for parser generation
- SysML v2 - The domain-specific language being parsed

## Runtime

**Environment:**
- Go 1.22+ (specified in `go.mod`)

**Package Manager:**
- Go modules (`go.mod`)
- Lockfile: `go.sum` (present)

## Frameworks

**Core:**
- ANTLR4-Go v4.13.1 - Parser generation framework
  - Generates lexer/parser from grammar files
  - Used in `internal/parser/` for generated code
  - Wrapped by `low/` and `sysml/` packages

**Testing:**
- Go standard library `testing` - Unit and integration tests
  - Test files: `*_test.go` throughout codebase
  - No external test framework dependencies

**Build/Dev:**
- Nix Flakes - Development environment management
  - Configuration: `flake.nix`
  - Provides: ANTLR4, OpenJDK, opencode

## Key Dependencies

**Critical:**
- `github.com/antlr4-go/antlr/v4` v4.13.1 - Core parser runtime
  - Used by: `gosysml2/low/parser.go`, `gosysml2/low/lexer.go`
  - Purpose: ANTLR4 runtime for Go

**Infrastructure:**
- `golang.org/x/exp` (indirect) - Extended Go packages
  - Required by ANTLR4 runtime

## Configuration

**Environment:**
- No environment variables required for library operation
- Development environment managed via Nix flake
- No `.env` files or external configuration needed

**Build:**
- `go.mod` - Module definition and dependencies
- `go.sum` - Dependency lock file
- `flake.nix` - Nix development shell configuration
- `flake.lock` - Nix flake lock file

## Platform Requirements

**Development:**
- Go 1.22 or later
- Nix (optional, for reproducible dev environment)
- ANTLR4 tool (for regenerating parser)
- OpenJDK (ANTLR4 runtime dependency)

**Production:**
- Pure Go library - no external runtime dependencies
- Single binary deployment for CLI tools
- Cross-platform: Linux, macOS, Windows (Go-supported platforms)

## Parser Generation

**Grammar Files:** (outside `gosysml2/` module)
- `code/SysMLv2Lexer.g4` - Lexer grammar
- `code/SysMLv2Parser.g4` - Parser grammar
- `code/KerMLParser.g4` - Kernel Modeling Language grammar

**Generated Code:**
- `gosysml2/internal/parser/` - ANTLR-generated Go code
  - `sysmlv2_lexer.go`
  - `sysmlv2_parser.go`
  - `sysmlv2parser_listener.go`
  - `sysmlv2parser_base_listener.go`

**Regenerate Command:**
```bash
cd code
antlr -Dlanguage=Go -o parser SysMLv2Lexer.g4 SysMLv2Parser.g4
```

---

*Stack analysis: 2026-02-05*
