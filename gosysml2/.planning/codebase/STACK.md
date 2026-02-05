# Technology Stack

**Analysis Date:** 2026-02-05

## Languages

**Primary:**
- Go 1.22 - Core language for parser implementation and high-level API

## Runtime

**Environment:**
- Go 1.22 runtime
- Works on x86_64-linux (primary target), cross-compilable

**Package Manager:**
- Go modules (go.mod/go.sum)
- Lockfile: present (`go.sum`)

## Frameworks

**Core:**
- ANTLR4 (antlr4-go) v4.13.1 - Parser generator for SysML v2 grammar
  - Provides lexer and parser generation from grammar files
  - Used for low-level API in `github.com/dVoo/gosysml2/low`
  - Generates parser code at `github.com/dVoo/gosysml2/internal/parser`

**Build/Dev:**
- Nix flakes - Development environment management
  - Located in `/home/daniel/projects/claudecode/flake.nix`
  - Provides isolated dev shell with Go, ANTLR, OpenJDK, and claude-code

## Key Dependencies

**Critical:**
- `github.com/antlr4-go/antlr/v4 v4.13.1` - Core parser infrastructure
  - Used by `low` package for lexer/parser interface
  - Provides token streaming and parse tree walking

**Indirect/Transitive:**
- `golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842` - Standard library experimental features
  - Likely used for runtime features or compatibility

## Configuration

**Environment:**
- No explicit .env configuration required
- Build configuration managed via Go build system
- Development environment: `nix develop` (from parent directory)

**Build:**
- Standard Go build: `go build ./...`
- Test execution: `go test ./...`
- No custom Makefiles or build scripts at project root

## Generated Code

**Parser Generation:**
The ANTLR-generated parser resides in `github.com/dVoo/gosysml2/internal/parser`:
- `sysmlv2_lexer.go` - Lexer rules for tokenization
- `sysmlv2_parser.go` - Parser rules for grammar
- `sysmlv2parser_base_listener.go` - ANTLR listener interface
- `sysmlv2parser_listener.go` - Listener implementation

These files are generated from grammar files (in separate parent repository at `/home/daniel/projects/claudecode/code/`) using ANTLR4 with Java.

## Platform Requirements

**Development:**
- Go 1.22+
- OpenJDK (for ANTLR code generation, if regenerating parser)
- ANTLR 4.13.1 (if regenerating parser)
- Nix (optional, for guaranteed environment reproducibility)

**Production:**
- Go 1.22+ runtime
- Deployed as Go library (consumed via `go get`)
- No external service dependencies
- Standalone parser library

## API Surface

**Public Packages:**
- `github.com/dVoo/gosysml2/sysml` - High-level model API
  - Idiomatic Go types and interfaces
  - Model building and visitor pattern
  - Directory parsing (sequential, parallel, streaming)
- `github.com/dVoo/gosysml2/low` - Low-level parser API
  - Direct ANTLR lexer/parser access
  - Token streaming
  - Parse tree inspection

**Internal Packages:**
- `github.com/dVoo/gosysml2/internal/parser` - ANTLR-generated code (not for direct use)

---

*Stack analysis: 2026-02-05*
