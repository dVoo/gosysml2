# Technology Stack

**Analysis Date:** 2026-02-05

## Languages

**Primary:**
- Go 1.22 - All application code

**Grammar Definition:**
- ANTLR4 Grammar (.g4) - Parser/lexer definitions
- BNF - Language specification documentation in `docs/bnf/`

## Runtime

**Environment:**
- Go 1.22 runtime
- Supports x86_64-linux (per `flake.nix`)

**Package Manager:**
- Go Modules (`go.mod`)
- Lockfile: `go.sum` present

**Development Environment:**
- Nix Flakes (`flake.nix`) for reproducible dev environment

## Frameworks & Tools

**Core Parsing:**
- ANTLR4 v4.13.1 - Parser generator for SysML v2 grammar
- `github.com/antlr4-go/antlr/v4` - Go ANTLR runtime

**Testing:**
- Go standard testing (`testing` package)
- No external test framework dependencies

**Build/Dev:**
- `nix develop` - Enter dev shell with ANTLR, OpenJDK, and tools
- `go build` - Standard Go compilation
- ANTLR CLI (via Nix) - Grammar compilation

## Key Dependencies

**Critical:**
- `github.com/antlr4-go/antlr/v4 v4.13.1` - ANTLR Go runtime for generated parsers
  - Required by all parser-generated code
  - Used in: `gosysml2/internal/parser/`, `code/parser/`

**Infrastructure:**
- `golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842` - Extended Go packages (indirect)

## Module Structure

**Three Go Modules:**

1. **Root Module** (`go.mod`)
   - Path: `github.com/dVoo/sysmlv2-tools`
   - Replaces: `github.com/dVoo/gosysml2` → `./gosysml2`

2. **Library Module** (`gosysml2/go.mod`)
   - Path: `github.com/dVoo/gosysml2`
   - Self-contained parser library
   - Contains: `internal/parser/`, `sysml/`, `low/`

3. **Verify Module** (`code/parser/go.mod`)
   - Path: `verify`
   - Standalone verification tool module
   - Also replaces gosysml2 → local path

## Configuration

**Environment:**
- No environment variables required
- No external configuration files
- Library is self-contained

**Build:**
- `go.mod` - Module definitions
- `flake.nix` - Nix development environment
- `.gitignore` - Excludes build artifacts, ANTLR generated files

**ANTLR Grammar Files:**
- `code/SysMLv2Lexer.g4` - Lexer definition
- `code/SysMLv2Parser.g4` - Parser definition
- `code/KerMLParser.g4` - Kernel Modeling Language parser

## Platform Requirements

**Development:**
- Nix with flakes support, OR
- Go 1.22+ with ANTLR4 and OpenJDK installed

**Production:**
- Pure Go - no CGO dependencies
- Cross-platform: Linux, macOS, Windows (Go supported platforms)
- No external runtime dependencies

**Building ANTLR Parsers:**
```bash
# Requires ANTLR4 and Java
cd code
antlr -Dlanguage=Go -o parser SysMLv2Lexer.g4 SysMLv2Parser.g4
```

---

*Stack analysis: 2026-02-05*
