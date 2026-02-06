// Requirements Traceability Example
//
// This example demonstrates:
// - Parsing a model with requirements
// - Show requirement relationships (derivedFrom, satisfiedBy, verifiedBy)
// - Demonstrate requirement ID extraction
// - Show how to find requirements by name
//
// Run with: go run main.go

package main

import (
	"fmt"
	"os"

	"github.com/dVoo/gosysml2_oc/sysml"
)

func main() {
	// Example SysML input with requirements and traceability
	// Note: Using simplified syntax supported by the parser
	input := `
package RequirementsModel {
    // System requirements
    requirement def REQ001 {
        doc /* The system shall provide user authentication. */
    }

    requirement def REQ002 {
        doc /* The system shall encrypt all data in transit. */
    }

    requirement def REQ003 {
        doc /* The system shall maintain 99.9% uptime. */
    }

    // Derived requirements
    requirement def REQ101 {
        doc /* The system shall support username/password authentication. */
    }

    requirement def REQ102 {
        doc /* The system shall support multi-factor authentication. */
    }

    requirement def REQ201 {
        doc /* The system shall use TLS 1.3 for all communications. */
    }

    // Parts that satisfy requirements
    part def AuthenticationModule;
    part def EncryptionModule;

    // Verification cases
    verification def VerifyAuth;
    verification def VerifyEncryption;

    // System usage
    part system;
}
`

	fmt.Println("=== Requirements Traceability Example ===")
	fmt.Println()

	// Parse the model
	fmt.Println("1. Parsing requirements model...")
	result := sysml.ParseString(input)

	if !result.Success() {
		fmt.Printf("   Parse errors:\n%s\n", result.Errors)
		os.Exit(1)
	}
	fmt.Println("   Parsing successful!")

	// Find all requirements
	requirements := sysml.FindRequirements(result.Model)
	fmt.Printf("\n2. Found %d requirement(s):\n", len(requirements))

	// Categorize requirements
	var definitions []*sysml.Requirement
	var usages []*sysml.Requirement
	for _, req := range requirements {
		if req.IsDefinition {
			definitions = append(definitions, req)
		} else {
			usages = append(usages, req)
		}
	}

	fmt.Printf("   - %d requirement definitions\n", len(definitions))
	fmt.Printf("   - %d requirement usages\n", len(usages))

	// Show requirement details
	fmt.Println("\n3. Requirement Details:")
	for _, req := range requirements {
		if !req.IsDefinition {
			continue // Skip usages for cleaner output
		}

		fmt.Printf("\n   Requirement: %s\n", req.Name())

		// Show documentation
		if doc := req.Documentation(); doc != "" {
			fmt.Printf("   Documentation: %s\n", doc)
		}

		// Show qualified name
		fmt.Printf("   Qualified name: %s\n", req.QualifiedName())

		// Show derived from relationships (if any)
		if len(req.DerivedFrom) > 0 {
			fmt.Println("   Derived from:")
			for _, derived := range req.DerivedFrom {
				fmt.Printf("      - %s\n", derived.Name())
			}
		}

		// Show derived requirements (inverse)
		if len(req.DerivedReqs) > 0 {
			fmt.Println("   Derived requirements:")
			for _, derived := range req.DerivedReqs {
				fmt.Printf("      - %s\n", derived.Name())
			}
		}

		// Show satisfied by
		if len(req.SatisfiedBy) > 0 {
			fmt.Println("   Satisfied by:")
			for _, elem := range req.SatisfiedBy {
				fmt.Printf("      - %s (%s)\n", elem.Name(), elem.Kind())
			}
		}

		// Show verified by
		if len(req.VerifiedBy) > 0 {
			fmt.Println("   Verified by:")
			for _, ver := range req.VerifiedBy {
				fmt.Printf("      - %s\n", ver.Name())
			}
		}
	}

	// Find requirements by name
	fmt.Println("\n4. Finding requirements by name:")
	targetName := "REQ001"
	found := sysml.Filter(result.Model, func(e sysml.Element) bool {
		return e.Name() == targetName
	})
	if len(found) > 0 {
		fmt.Printf("   Found '%s': %s (%s)\n", targetName, found[0].Name(), found[0].Kind())
		if req, ok := found[0].(*sysml.Requirement); ok {
			if doc := req.Documentation(); doc != "" {
				fmt.Printf("   Documentation: %s\n", doc)
			}
		}
	} else {
		fmt.Printf("   Requirement '%s' not found\n", targetName)
	}

	// Find verification cases
	verifications := sysml.FindVerifications(result.Model)
	fmt.Printf("\n5. Found %d verification case(s):\n", len(verifications))
	for _, ver := range verifications {
		fmt.Printf("   - %s\n", ver.Name())
		if ver.VerifiedRequirement != nil {
			fmt.Printf("     Verifies: %s\n", ver.VerifiedRequirement.Name())
		}
	}

	// Show parts that could satisfy requirements
	parts := sysml.FindParts(result.Model)
	fmt.Printf("\n6. Found %d part(s) in model:\n", len(parts))
	for _, part := range parts {
		defType := "usage"
		if part.IsDefinition {
			defType = "definition"
		}
		fmt.Printf("   - %s (%s)\n", part.Name(), defType)
	}

	// Count elements by type
	fmt.Println("\n7. Model Statistics:")
	counter := sysml.NewCounter()
	sysml.Visit(result.Model, counter)
	fmt.Printf("   Total elements: %d\n", counter.Total())
	fmt.Printf("   Requirements: %d\n", counter.Counts[sysml.KindRequirement])
	fmt.Printf("   Verifications: %d\n", counter.Counts[sysml.KindVerification])
	fmt.Printf("   Parts: %d\n", counter.Counts[sysml.KindPart])
	fmt.Printf("   Packages: %d\n", counter.Counts[sysml.KindPackage])

	// Show model hierarchy
	fmt.Println("\n8. Model Hierarchy:")
	sysml.Walk(result.Model, func(elem sysml.Element, depth int) bool {
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}
		fmt.Printf("%s- %s: %s\n", indent, elem.Kind(), elem.Name())
		return true
	})

	fmt.Println("\n=== Example Complete ===")
}
