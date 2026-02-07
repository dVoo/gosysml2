// Package low provides low-level access to the SysML v2 parser.
//
// This package wraps the ANTLR-generated lexer and parser with a thin
// Go-friendly layer that provides error collection, token access, and
// context support. It offers direct access to parse trees for advanced use
// cases where the high-level sysml package model is not needed.
//
// # When to Use This Package
//
// Use the low package when you need:
//   - Direct access to the ANTLR parse tree
//   - Custom parse tree traversal or transformation
//   - Token-level inspection and manipulation
//   - Validation-only parsing (no model building)
//   - Maximum parsing performance with minimal overhead
//   - Integration with other ANTLR-based tools
//
// For most use cases, the high-level sysml package is recommended as it
// provides a cleaner, more idiomatic Go API with type-safe model elements.
//
// # Quick Start
//
// Parse and get the raw parse tree:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/dVoo/gosysml2/low"
//	)
//
//	func main() {
//	    input := `
//	        package Vehicle {
//	            part def Engine;
//	        }
//	    `
//
//	    tree, errors := low.Parse(input)
//	    if errors.HasErrors() {
//	        for _, err := range errors.All() {
//	            fmt.Printf("Line %d: %s\n", err.Line, err.Message)
//	        }
//	        return
//	    }
//
//	    // Access the parse tree
//	    fmt.Printf("Parse tree root: %v\n", tree)
//	}
//
// # Key Types
//
// Parser: Wraps the ANTLR parser with error collection and context support.
// Provides methods for parsing different entry points.
//
// Lexer: Wraps the ANTLR lexer with error collection.
// Provides token stream access.
//
// ParseErrors: Aggregates lexer and parser errors with line/column info.
//
// ParseOption: Functional options for configuring parser behavior.
//
// # Parse Options
//
// Control parser behavior with options:
//
//	// Disable parse tree construction for validation-only
//	tree, errors := low.Parse(input, low.WithParseTree(false))
//
//	// Set context for cancellation support
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	p := low.NewParser(input, low.WithContext(ctx))
//	tree := p.ParseRootNamespace()
//
// # Token Access
//
// Inspect individual tokens from the input:
//
//	lexer := low.NewLexer(input)
//	for token := lexer.NextToken(); token != nil; token = lexer.NextToken() {
//	    tokenType := low.TokenName(token.GetTokenType())
//	    text := token.GetText()
//	    line := token.GetLine()
//	    column := token.GetColumn()
//	    fmt.Printf("%s: '%s' at %d:%d\n", tokenType, text, line, column)
//	}
//
// # Validation Only
//
// For syntax validation without building a parse tree:
//
//	errors := low.Validate(input)
//	if errors.HasErrors() {
//	    fmt.Printf("Validation failed: %s\n", errors.Error())
//	}
//
// # Performance Notes
//
// The low package provides the fastest parsing path with minimal allocations:
//
//   - WithParseTree(false): Skips parse tree construction, fastest option
//   - Token streaming: Lexer consumes tokens lazily, reducing memory
//   - No model building: Direct tree access avoids high-level object creation
//
// For maximum performance in validation scenarios:
//
//	// Fastest validation path
//	errors := low.Validate(input)
//
// # Error Handling
//
// Errors include detailed source location information:
//
//	if errors.HasErrors() {
//	    // Get first error
//	    first := errors.First()
//	    fmt.Printf("Error at line %d, column %d: %s\n",
//	        first.Line, first.Column, first.Message)
//
//	    // Iterate all errors
//	    for _, err := range errors.All() {
//	        fmt.Printf("- %s\n", err)
//	    }
//	}
//
// # Context Support
//
// The parser respects context cancellation for long-running operations:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	tree, err := low.ParseWithContext(ctx, input)
//	if err != nil {
//	    if errors.Is(err, context.DeadlineExceeded) {
//	        fmt.Println("Parsing timed out")
//	    }
//	}
//
// # Relationship to sysml Package
//
// The sysml package builds on top of this low package:
//
//	// low package: Raw parse tree
//	tree, errors := low.Parse(input)
//
//	// sysml package: High-level model (uses low internally)
//	result := sysml.ParseString(input)
//	model := result.Model
//
// Use low when you need parse tree access; use sysml when you need
// a structured model with resolved references.
//
// # Thread Safety
//
// Parser and Lexer instances are not safe for concurrent use. Create
// separate instances for concurrent parsing operations.
package low
