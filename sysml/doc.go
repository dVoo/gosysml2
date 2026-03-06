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
//	    "github.com/dVoo/gosysml2/sysml"
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
//	    if err := result.Err(); err != nil {
//	        fmt.Printf("Parse error: %s\n", err)
//	        return
//	    }
//
//	    // Access the model
//	    for _, pkg := range result.Model.Packages() {
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
//	// Parse from bytes.
//	// This avoids an extra high-level []byte -> string copy.
//	result := sysml.ParseBytes(data, "source.sysml")
//
//	// Parse directory with configurable workers and parse options
//	opts := sysml.DirOptions{
//	    Workers:      0, // 0 = runtime.NumCPU()
//	    ParseOptions: []sysml.ParseOption{sysml.WithDiscardTree()},
//	}
//	for r := range sysml.ParseDir(context.Background(), "./models", opts) {
//	    if err := r.Err(); err != nil {
//	        fmt.Printf("failed: %s: %v\n", r.Source, err)
//	        continue
//	    }
//	}
//	// Note: with Workers > 1, ParseDir does not preserve result order.
//
// # Parse Options
//
// Control parsing behavior with options:
//
//	// Discard parse tree to reduce memory (~30% savings)
//	result := sysml.ParseFile("model.sysml", sysml.WithDiscardTree())
//
//	// Build a model but skip index build and reference resolution
//	result := sysml.ParseFile("model.sysml", sysml.WithoutResolution())
//
//	// Syntax parse only, skip high-level model construction entirely
//	result := sysml.ParseFile("model.sysml",
//	    sysml.WithoutModelBuild(),
//	    sysml.WithDiscardTree())
//
//	// Compute and retain input SHA-256 during parsing
//	result := sysml.ParseFile("model.sysml", sysml.WithContentHash())
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
//	parts := sysml.FindAll[*sysml.Part](model)
//	requirements := sysml.FindAll[*sysml.Requirement](model)
//	verifications := sysml.FindAll[*sysml.Verification](model)
//
//	// Find by qualified name
//	elem := sysml.FindByQualifiedName(model, "Vehicle::Engine")
//
//	// Semantic classification without reflecting on concrete structs
//	role := elem.Role()
//	if role.IsUsage() {
//	    if usage, ok := elem.(sysml.Usage); ok {
//	        fmt.Println(usage.TypeName())
//	    }
//	}
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
// Traversal API guidance:
//   - Prefer iterators (All/OfType/OfKind) for most filtering and collection tasks
//   - Use Walk when you need depth and early exit
//   - Use Visit for larger structured processors with per-type methods
//   - Use WalkAll for callback traversal without early-exit handling
//
// # Position Lookup (Editor/LSP)
//
// Use ElementAt for named elements and ElementAtIncludingUnnamed when anonymous
// usages should still be returned (for example, `interface :T` declarations).
//
//	elem := sysml.ElementAt(model, line, column)
//	_ = elem
//	elem = sysml.ElementAtIncludingUnnamed(model, line, column)
//
// # Performance
//
// For large repositories, use ParseDir with appropriate worker count:
//
//	// Sequential streaming style (lowest memory)
//	opts := sysml.DirOptions{Workers: 1, ParseOptions: []sysml.ParseOption{sysml.WithDiscardTree()}}
//	for r := range sysml.ParseDir(context.Background(), "./models", opts) {
//	    if r.Err() == nil {
//	        // Process each model as it's parsed
//	        parts := sysml.FindAll[*sysml.Part](r.Model)
//	        fmt.Printf("%s: %d parts\n", r.Source, len(parts))
//	    }
//	}
//
// # Error Handling
//
// Parse results include detailed error information:
//
//	if err := result.Err(); err != nil {
//	    // Get first error
//	    first := result.ParseError.First()
//	    fmt.Printf("Error at line %d, column %d: %s\n",
//	        first.Line, first.Column, first.Message)
//	}
//
// # Additional Documentation
//
// For complete API documentation and examples, see:
//   - README.md in the gosysml2 directory
//   - examples/ directory for working code samples
//   - Go doc: https://pkg.go.dev/github.com/dVoo/gosysml2/sysml
//
// # Thread Safety
//
// Parser instances are not safe for concurrent use. Create a new parser
// for each concurrent parsing operation. The resulting Model and its
// elements are read-only after parsing and safe for concurrent access.
//
// # Parent-Child Relationships
//
// Elements in a SysML model form a tree structure. Each element (except root elements)
// has a parent that can be accessed via the Parent() method. The parent reference is
// set automatically when an element is added to a container (Package, Part, etc.).
//
// Important: The Parent() method returns the concrete container type (*Package, *Part, etc.)
// not *baseElement. This allows type assertions to work correctly:
//
//	parent := element.Parent()
//	if pkg, ok := parent.(*Package); ok {
//	    fmt.Printf("Parent is package: %s\n", pkg.Name())
//	} else if part, ok := parent.(*Part); ok {
//	    fmt.Printf("Parent is part: %s\n", part.Name())
//	}
//
// The parent reference is established during parsing and is available immediately
// after the element is added to its container via AddChild().
package sysml
