// Visitor Pattern Example
//
// This example demonstrates:
// - Implementing a custom visitor struct
// - Counting specific element types
// - Filtering elements by criteria
// - Showing depth tracking during walk
//
// Run with: go run main.go

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dVoo/gosysml2/sysml"
)

// RequirementAuditVisitor is a custom visitor that audits requirements
// It finds requirements without IDs and tracks statistics
type RequirementAuditVisitor struct {
	sysml.BaseVisitor
	RequirementsWithoutID []*sysml.Requirement
	RequirementsWithDoc   []*sysml.Requirement
	TotalRequirements     int
	TotalParts            int
	MaxDepth              int
	CurrentDepth          int
}

// VisitRequirement implements custom logic for requirements
func (v *RequirementAuditVisitor) VisitRequirement(req *sysml.Requirement) bool {
	v.TotalRequirements++

	// Track depth
	if v.CurrentDepth > v.MaxDepth {
		v.MaxDepth = v.CurrentDepth
	}

	// Find requirements without RequirementID
	if req.RequirementID == "" {
		v.RequirementsWithoutID = append(v.RequirementsWithoutID, req)
	}

	// Find requirements with documentation
	if req.Documentation() != "" {
		v.RequirementsWithDoc = append(v.RequirementsWithDoc, req)
	}

	v.CurrentDepth++
	return true // continue visiting children
}

// VisitPart counts parts and tracks depth
func (v *RequirementAuditVisitor) VisitPart(part *sysml.Part) bool {
	v.TotalParts++
	if v.CurrentDepth > v.MaxDepth {
		v.MaxDepth = v.CurrentDepth
	}
	v.CurrentDepth++
	return true
}

// DepthTrackingVisitor tracks the depth of each element type
type DepthTrackingVisitor struct {
	sysml.BaseVisitor
	DepthsByKind map[sysml.ElementKind][]int
	currentDepth int
}

func NewDepthTrackingVisitor() *DepthTrackingVisitor {
	return &DepthTrackingVisitor{
		DepthsByKind: make(map[sysml.ElementKind][]int),
	}
}

func (v *DepthTrackingVisitor) trackDepth(kind sysml.ElementKind) {
	v.DepthsByKind[kind] = append(v.DepthsByKind[kind], v.currentDepth)
	v.currentDepth++
}

func (v *DepthTrackingVisitor) VisitPackage(pkg *sysml.Package) bool {
	v.trackDepth(sysml.KindPackage)
	return true
}

func (v *DepthTrackingVisitor) VisitPart(part *sysml.Part) bool {
	v.trackDepth(sysml.KindPart)
	return true
}

func (v *DepthTrackingVisitor) VisitRequirement(req *sysml.Requirement) bool {
	v.trackDepth(sysml.KindRequirement)
	return true
}

func (v *DepthTrackingVisitor) VisitVerification(ver *sysml.Verification) bool {
	v.trackDepth(sysml.KindVerification)
	return true
}

func (v *DepthTrackingVisitor) VisitElement(elem sysml.Element) bool {
	v.trackDepth(elem.Kind())
	return true
}

// NameFilterVisitor finds elements matching a name pattern
type NameFilterVisitor struct {
	sysml.BaseVisitor
	Pattern       string
	Matches       []sysml.Element
	MatchCount    int
	caseSensitive bool
}

func NewNameFilterVisitor(pattern string, caseSensitive bool) *NameFilterVisitor {
	return &NameFilterVisitor{
		Pattern:       pattern,
		caseSensitive: caseSensitive,
	}
}

func (v *NameFilterVisitor) VisitElement(elem sysml.Element) bool {
	name := elem.Name()
	pattern := v.Pattern

	if !v.caseSensitive {
		name = strings.ToLower(name)
		pattern = strings.ToLower(pattern)
	}

	if strings.Contains(name, pattern) {
		v.Matches = append(v.Matches, elem)
		v.MatchCount++
	}
	return true
}

// HierarchyPrinterVisitor prints the model hierarchy with indentation
type HierarchyPrinterVisitor struct {
	sysml.BaseVisitor
	Output []string
	indent string
	depth  int
}

func (v *HierarchyPrinterVisitor) VisitPackage(pkg *sysml.Package) bool {
	v.Output = append(v.Output, fmt.Sprintf("%s📦 %s", v.indent, pkg.Name()))
	v.indent += "  "
	v.depth++
	return true
}

func (v *HierarchyPrinterVisitor) ExitPackage(pkg *sysml.Package) {
	if len(v.indent) >= 2 {
		v.indent = v.indent[:len(v.indent)-2]
	}
	v.depth--
}

func (v *HierarchyPrinterVisitor) VisitPart(part *sysml.Part) bool {
	icon := "🔧"
	if part.IsDefinition {
		icon = "📐"
	}
	v.Output = append(v.Output, fmt.Sprintf("%s%s %s", v.indent, icon, part.Name()))
	v.indent += "  "
	v.depth++
	return true
}

func (v *HierarchyPrinterVisitor) VisitRequirement(req *sysml.Requirement) bool {
	v.Output = append(v.Output, fmt.Sprintf("%s📋 %s", v.indent, req.Name()))
	v.indent += "  "
	v.depth++
	return true
}

func (v *HierarchyPrinterVisitor) VisitVerification(ver *sysml.Verification) bool {
	v.Output = append(v.Output, fmt.Sprintf("%s✅ %s", v.indent, ver.Name()))
	return true
}

func main() {
	// Example SysML input
	input := `
package SystemModel {
    // Requirements
    requirement def REQ001 {
        doc /* The system shall provide user authentication. */
    }

    requirement def REQ002 {
        doc /* The system shall encrypt all data in transit. */
    }

    requirement def REQ003 {
        // No documentation
    }

    requirement def REQ004;

    // Parts
    part def Vehicle {
        part engine : Engine;
        part wheels : Wheel[4];
    }

    part def Engine {
        attribute horsepower : Real;
    }

    part def Wheel;

    part def Chassis;

    // Verification cases
    verification def VerifyAuth;
    verification def VerifyEncryption;

    // System usage
    part system : Vehicle;
}
`

	fmt.Println("=== Visitor Pattern Example ===")
	fmt.Println()

	// Parse the model
	fmt.Println("1. Parsing model...")
	result := sysml.ParseString(input)
	if result.Err() != nil {
		fmt.Printf("   Parse errors:\n%s\n", result.ParseError)
		os.Exit(1)
	}
	fmt.Println("   Parsing successful!")

	// Example 1: Custom visitor for requirement auditing
	fmt.Println("\n2. Requirement Audit Visitor:")
	auditVisitor := &RequirementAuditVisitor{}
	sysml.Visit(result.Model, auditVisitor)

	fmt.Printf("   Total requirements: %d\n", auditVisitor.TotalRequirements)
	fmt.Printf("   Total parts: %d\n", auditVisitor.TotalParts)
	fmt.Printf("   Max depth reached: %d\n", auditVisitor.MaxDepth)

	fmt.Printf("   Requirements with documentation: %d\n", len(auditVisitor.RequirementsWithDoc))
	for _, req := range auditVisitor.RequirementsWithDoc {
		fmt.Printf("      - %s\n", req.Name())
	}

	fmt.Printf("   Requirements without ID: %d\n", len(auditVisitor.RequirementsWithoutID))
	for _, req := range auditVisitor.RequirementsWithoutID {
		fmt.Printf("      - %s\n", req.Name())
	}

	// Example 2: Depth tracking visitor
	fmt.Println("\n3. Depth Tracking Visitor:")
	depthVisitor := NewDepthTrackingVisitor()
	sysml.Visit(result.Model, depthVisitor)

	fmt.Println("   Depth distribution by element kind:")
	for kind, depths := range depthVisitor.DepthsByKind {
		if len(depths) > 0 {
			min, max, avg := calculateDepthStats(depths)
			fmt.Printf("      %s: count=%d, min=%d, max=%d, avg=%.1f\n",
				kind, len(depths), min, max, avg)
		}
	}

	// Example 3: Name filter visitor
	fmt.Println("\n4. Name Filter Visitor (searching for 'REQ'):")
	filterVisitor := NewNameFilterVisitor("REQ", false)
	sysml.Visit(result.Model, filterVisitor)

	fmt.Printf("   Found %d elements matching 'REQ':\n", filterVisitor.MatchCount)
	for _, elem := range filterVisitor.Matches {
		fmt.Printf("      - %s (%s)\n", elem.Name(), elem.Kind())
	}

	// Example 4: Using Walk for depth tracking
	fmt.Println("\n5. Walk Function with Depth:")
	depths := make(map[sysml.Element]int)
	sysml.Walk(result.Model, func(elem sysml.Element, depth int) bool {
		depths[elem] = depth
		return true
	})

	fmt.Println("   Elements at each depth level:")
	maxDepth := 0
	for _, d := range depths {
		if d > maxDepth {
			maxDepth = d
		}
	}
	for d := 0; d <= maxDepth; d++ {
		count := 0
		var names []string
		for elem, elemDepth := range depths {
			if elemDepth == d {
				count++
				if len(names) < 5 {
					names = append(names, fmt.Sprintf("%s:%s", elem.Kind(), elem.Name()))
				}
			}
		}
		fmt.Printf("      Depth %d: %d elements", d, count)
		if len(names) > 0 {
			fmt.Printf(" (%s...)", strings.Join(names, ", "))
		}
		fmt.Println()
	}

	// Example 5: Filter with custom predicate
	fmt.Println("\n6. Custom Filter (definitions only):")
	definitions := sysml.Filter(result.Model, func(e sysml.Element) bool {
		// Check if element is a definition
		switch elem := e.(type) {
		case *sysml.Part:
			return elem.IsDefinition
		case *sysml.Requirement:
			return elem.IsDefinition
		case *sysml.Verification:
			return elem.IsDefinition
		default:
			return false
		}
	})

	fmt.Printf("   Found %d definition elements:\n", len(definitions))
	for _, def := range definitions {
		fmt.Printf("      - %s (%s)\n", def.Name(), def.Kind())
	}

	// Example 6: Count elements by kind using built-in Counter
	fmt.Println("\n7. Built-in Counter:")
	counter := sysml.NewCounter()
	sysml.Visit(result.Model, counter)

	fmt.Println("   Element counts:")
	for kind, count := range counter.Counts {
		if count > 0 {
			fmt.Printf("      %s: %d\n", kind, count)
		}
	}
	fmt.Printf("   Total: %d\n", counter.Total())

	// Example 7: FindAll with generics
	fmt.Println("\n8. Generic FindAll:")
	requirements := sysml.FindAll[*sysml.Requirement](result.Model)
	parts := sysml.FindAll[*sysml.Part](result.Model)
	verifications := sysml.FindAll[*sysml.Verification](result.Model)

	fmt.Printf("   Requirements: %d\n", len(requirements))
	fmt.Printf("   Parts: %d\n", len(parts))
	fmt.Printf("   Verifications: %d\n", len(verifications))

	// Example 9: Iterator usage (Go 1.23+)
	fmt.Println("\n9. Iterator Pattern:")
	fmt.Println("   Using All() iterator:")
	count := 0
	for elem := range sysml.All(result.Model) {
		count++
		if count <= 5 {
			fmt.Printf("      - %s: %s\n", elem.Kind(), elem.Name())
		}
	}
	fmt.Printf("   (showing first 5 of %d elements)\n", count)

	// Example 10: OfKind iterator
	fmt.Println("\n10. OfKind Iterator:")
	fmt.Println("    All requirement elements:")
	for elem := range sysml.OfKind(result.Model, sysml.KindRequirement) {
		fmt.Printf("       - %s\n", elem.Name())
	}

	fmt.Println("\n=== Example Complete ===")
}

func calculateDepthStats(depths []int) (min, max int, avg float64) {
	if len(depths) == 0 {
		return 0, 0, 0
	}
	min = depths[0]
	max = depths[0]
	sum := 0
	for _, d := range depths {
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
		sum += d
	}
	avg = float64(sum) / float64(len(depths))
	return
}
