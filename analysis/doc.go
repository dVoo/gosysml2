// Package analysis provides utilities for analyzing and transforming SysML models.
//
// This package includes helpers for D2 diagram generation, including valid shape
// mappings from SysML element types to D2 shapes.
//
// # D2 Shape Mapping
//
// D2 (Declarative Diagramming) is a diagram scripting language. This package provides
// utilities to map SysML element kinds to valid D2 shapes.
//
// Common mistake: Using "folder" shape (INVALID). Use "package" instead.
//
// Example:
//
//	import "github.com/dVoo/gosysml2/analysis"
//
//	// Get valid D2 shape for a SysML element
//	shape := analysis.ShapeForElement("package")
//	// Returns: "package" (not "folder"!)
//
//	// Validate a shape string
//	if !analysis.IsValidShape("folder") {
//	    // Invalid! Use FixShape to get valid equivalent
//	    shape = analysis.FixShape("folder")
//	    // Returns: "package"
//	}
//
// # Available Shapes
//
// Valid D2 shapes (BlockType constants):
//   - BlockRectangle, BlockSquare, BlockCircle, BlockOval
//   - BlockDiamond, BlockHexagon, BlockCloud, BlockCylinder
//   - BlockQueue, BlockPackage, BlockPage, BlockDocument
//   - BlockStoredData, BlockPerson, BlockSQLTable, BlockClass
//   - BlockCode, BlockImage, BlockSequence
//
// See: https://d2-lang.com/tour/shapes
//
// # SysML to D2 Mappings
//
// Recommended mappings for SysML elements:
//   - Package → BlockPackage
//   - Part → BlockRectangle
//   - Requirement → BlockDocument
//   - UseCase → BlockOval
//   - Actor → BlockPerson
//   - Port → BlockCircle
//   - State → BlockHexagon
//   - Action → BlockRectangle
//
package analysis
