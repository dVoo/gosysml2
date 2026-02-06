// Package sysml provides a high-level, idiomatic Go API for parsing SysML v2 models.
//
// This package wraps the low-level ANTLR-generated parser with a clean, type-safe
// model representation that follows Go conventions. It provides comprehensive
// parsing capabilities, element traversal, and reference resolution.
//
// # Overview
//
// The sysml package offers a two-tier API design:
//   - High-level API (this package): Idiomatic Go model with visitor pattern
//   - Low-level API (low package): Direct access to lexer, parser, and parse trees
//
// Use this package when you need:
//   - Type-safe access to SysML elements
//   - Reference resolution between elements
//   - Visitor pattern for model traversal
//   - Convenient parsing functions
//
// # Quick Start
//
// Parse a SysML string:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/dVoo/gosysml2_oc/sysml"
//	)
//
//	func main() {
//	    input := `
//	        package Vehicle {
//	            part def Engine {
//	                attribute power : Real;
//	            }
//	            part def Car {
//	                part engine : Engine;
//	            }
//	        }
//	    `
//
//	    result := sysml.ParseString(input)
//	    if !result.Success() {
//	        fmt.Printf("Parse error: %s\n", result.Errors)
//	        return
//	    }
//
//	    // Access the model
//	    for _, pkg := range result.Model.Packages {
//	        fmt.Printf("Package: %s\n", pkg.Name())
//	    }
//	}
//
// # Key Types
//
// The following types represent the core SysML v2 model elements:
//
// Model: The root container holding all parsed elements.
// Provides element indexing and reference resolution.
//
// Package: A namespace containing related elements.
// Packages can be nested and support typed accessors for children.
//
// Part: Represents part definitions and usages.
// Parts can contain attributes, nested parts, and ports.
//
// Requirement: Captures requirement definitions and usages
// with support for derivation, satisfaction, and verification.
//
// Attribute: Typed attributes with optional default values.
//
// Port: Interface points with direction (in, out, inout).
//
// # Type-Safe References
//
// The package uses generic Ref[T] types for element references:
//
//	// Reference to a part definition
//	var typeRef Ref[*Part]
//
//	// Check if reference is resolved
//	if typeRef.IsResolved() {
//	    def := typeRef.Resolved()
//	    fmt.Println(def.Name())
//	}
//
// References are automatically resolved after parsing via
// Model.ResolveReferences().
//
// # Parsing Functions
//
// The package provides multiple parsing entry points:
//
//	// Parse from string
//	result := sysml.ParseString(input)
//
//	// Parse from file
//	result := sysml.ParseFile("model.sysml")
//
//	// Parse from bytes (avoids copy)
//	result := sysml.ParseBytes(data, "source.sysml")
//
//	// Parse directory
//	results, err := sysml.ParseDirectory("./models")
//
//	// Parse directory in parallel
//	results, err := sysml.ParseDirectoryParallel("./models", 4)
//
// # Parse Options
//
// Control parsing behavior with options:
//
//	// Discard parse tree to reduce memory (~30% savings)
//	result := sysml.ParseFile("model.sysml", sysml.WithDiscardTree())
//
//	// Use standard library for imports
//	result := sysml.ParseFile("model.sysml", sysml.WithStandardLibrary())
//
//	// Use custom library path
//	result := sysml.ParseFile("model.sysml",
//	    sysml.WithStandardLibrary(),
//	    sysml.WithLibraryPath("/path/to/libraries"))
//
// # Finder Functions
//
// Convenient functions to find elements by type:
//
//	parts := sysml.FindParts(model)
//	requirements := sysml.FindRequirements(model)
//	verifications := sysml.FindVerifications(model)
//
//	// Find by qualified name
//	elem := sysml.FindByQualifiedName(model, "Vehicle::Engine")
//
// # Visitor Pattern
//
// Implement custom visitors for model traversal:
//
//	type PartCounter struct {
//	    sysml.BaseVisitor
//	    Count int
//	}
//
//	func (v *PartCounter) VisitPart(part *sysml.Part) bool {
//	    v.Count++
//	    return true // continue visiting children
//	}
//
//	counter := &PartCounter{}
//	sysml.Visit(model, counter)
//
// # Performance
//
// For large repositories, use streaming or parallel parsing:
//
//	// Streaming (lowest memory)
//	err := sysml.ParseDirectoryStream("./models", func(r *sysml.ParseResult) error {
//	    if r.Success() {
//	        // Process each model as it's parsed
//	        parts := sysml.FindParts(r.Model)
//	        fmt.Printf("%s: %d parts\n", r.Source, len(parts))
//	    }
//	    return nil
//	}, sysml.WithDiscardTree())
//
// # Error Handling
//
// Parse results include detailed error information:
//
//	if !result.Success() {
//	    // Get first error
//	    first := result.Errors.First()
//	    fmt.Printf("Error at line %d, column %d: %s\n",
//	        first.Line, first.Column, first.Message)
//	}
//
// # Additional Documentation
//
// For complete API documentation and examples, see:
//   - README.md in the gosysml2 directory
//   - examples/ directory for working code samples
//   - Go doc: https://pkg.go.dev/github.com/dVoo/gosysml2_oc/sysml
//
// # Thread Safety
//
// Parser instances are not safe for concurrent use. Create a new parser
// for each concurrent parsing operation. The resulting Model and its
// elements are read-only after parsing and safe for concurrent access.
package sysml
