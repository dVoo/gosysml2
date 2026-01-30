# Technology Stack

**Analysis Date:** 2026-01-30

## Languages

**Primary:**
- Go 1.22 - Used in `gosysml2/` module for SysML parser implementation and model manipulation
- ANTLR4 Grammar Files (.g4) - Parser and lexer grammar specifications for SysML v2 and KerML languages

**Secondary:**
- EBNF/BNF (Backus-Naur Form) - Used in `.kebnf` and `.kgbnf` files for formal grammar documentation
- SVG - 284 graphical notation diagram files in `docs/bnf/images/`
- HTML - Rendered documentation of BNF specifications

## Runtime

**Environment:**
- Go 1.22 (from `gosysml2/go.mod`)
- OpenJDK (for ANTLR4 tool execution)

**Package Manager:**
- Go modules (go.mod/go.sum)
- Nix Flakes for development environment management

## Frameworks

**Core:**
- ANTLR4 (v4.13.2) - Parser generator for creating SysML and KerML lexers and parsers
  - `antlr4-go/antlr/v4 v4.13.1` - Go runtime bindings for ANTLR4

**Testing:**
- Go's standard `testing` package - Unit tests in `*_test.go` files

**Build/Dev:**
- Nix Flakes (`flake.nix`, `flake.lock`) - Reproducible development environment
- ANTLR command-line tool (v4.13.2+) - Parser generation tool

## Key Dependencies

**Critical:**
- `github.com/antlr4-go/antlr/v4 v4.13.1` - Runtime for ANTLR-generated parsers and lexers. Used by `gosysml2/sysml` and `gosysml2/low` packages for parsing SysML text.
- `golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842` - Go experimental package (indirect dependency)

**Infrastructure:**
- None detected - No external service SDKs or infrastructure clients

## Configuration

**Environment:**
- No `.env` files detected
- Development environment managed through Nix Flakes
- No configuration files for development parameters

**Build:**
- `flake.nix` - Nix Flake configuration providing development environment with:
  - `claude-code` - Claude Code integration
  - `antlr` - ANTLR4 parser generator
  - `openjdk` - Java runtime for ANTLR

## Platform Requirements

**Development:**
- Nix 2.4+ with flakes support
- x86_64-linux architecture (configured in `flake.nix`)
- OpenJDK for ANTLR4 tool execution
- Go 1.22 runtime (provided by Nix flake)

**Production:**
- Go 1.22+ runtime
- No external dependencies required for `gosysml2` when used as library
- Parser requires ANTLR4 Go runtime (`antlr/v4`)

## Code Generation

**Parser Generation:**
- ANTLR4 grammar files in `code/` are compiled to Go parsers:
  - `code/SysMLv2Lexer.g4` → `code/parser/sysmlv2_lexer.go`
  - `code/KerMLParser.g4` → (KerML parser)
  - `code/SysMLv2Parser.g4` → `code/parser/sysmlv2_parser.go`
- Generated code includes lexer, parser, listeners, and base listeners
- ANTLR generation command: `antlr -Dlanguage=Go -o parser [grammar files]`

## Module Structure

**Main Module:**
- `github.com/dVoo/gosysml2` - Go SysML v2 parser library

**Sub-packages:**
- `gosysml2/sysml` - High-level idiomatic Go API for parsed models
- `gosysml2/low` - Low-level parser API with direct parse tree access
- `gosysml2/internal/parser` - ANTLR4-generated parser code (internal)

---

*Stack analysis: 2026-01-30*
