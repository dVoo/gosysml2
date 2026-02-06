// Basic Parsing Example
//
// This example demonstrates:
// - Parsing a simple SysML model from string
// - Accessing packages and parts
// - Using FindParts and FindRequirements
// - Error handling
//
// Run with: go run main.go

package main

import (
	"fmt"
	"os"

	"github.com/dVoo/gosysml2_oc/sysml"
)

func main() {
	// Example SysML input - a simple Vehicle model
	input := `
package VehicleModel {
    part def Vehicle {
        part engine : Engine;
        part wheels : Wheel[4];
    }

    part def Engine {
        attribute horsepower : Real;
    }

    part def Wheel;

    requirement def SafetyReq {
        doc /* The vehicle must be safe. */
    }

    part vehicle1 : Vehicle;
}
`

	fmt.Println("=== Basic Parsing Example ===")
	fmt.Println()

	// Step 1: Parse the SysML string
	fmt.Println("1. Parsing SysML input...")
	result := sysml.ParseString(input)

	// Step 2: Check for errors
	if !result.Success() {
		fmt.Printf("   Parse errors:\n%s\n", result.Errors)
		os.Exit(1)
	}
	fmt.Println("   Parsing successful!")

	// Step 3: Access top-level packages
	fmt.Printf("\n2. Found %d top-level package(s):\n", len(result.Model.Packages))
	for _, pkg := range result.Model.Packages {
		fmt.Printf("   - Package: %s\n", pkg.Name())
	}

	// Step 4: Find all parts in the model
	parts := sysml.FindParts(result.Model)
	fmt.Printf("\n3. Found %d part(s):\n", len(parts))
	for _, part := range parts {
		defType := "usage"
		if part.IsDefinition {
			defType = "definition"
		}
		fmt.Printf("   - %s (%s)\n", part.Name(), defType)
	}

	// Step 5: Find all requirements
	requirements := sysml.FindRequirements(result.Model)
	fmt.Printf("\n4. Found %d requirement(s):\n", len(requirements))
	for _, req := range requirements {
		fmt.Printf("   - %s\n", req.Name())
		if doc := req.Documentation(); doc != "" {
			fmt.Printf("     Documentation: %s\n", doc)
		}
	}

	// Step 6: Walk the model to show hierarchy
	fmt.Println("\n5. Model hierarchy:")
	sysml.Walk(result.Model, func(elem sysml.Element, depth int) bool {
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}
		fmt.Printf("%s- %s: %s\n", indent, elem.Kind(), elem.Name())
		return true
	})

	// Step 7: Use a visitor to count elements
	fmt.Println("\n6. Element counts by kind:")
	counter := sysml.NewCounter()
	sysml.Visit(result.Model, counter)
	for kind, count := range counter.Counts {
		if count > 0 {
			fmt.Printf("   %s: %d\n", kind, count)
		}
	}
	fmt.Printf("   Total elements: %d\n", counter.Total())

	fmt.Println("\n=== Example Complete ===")
}
