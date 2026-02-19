// Example usage of the gosysml2 library.
package main

import (
	"fmt"

	"github.com/dVoo/gosysml2/low"
	"github.com/dVoo/gosysml2/sysml"
)

func main() {
	// Example SysML input
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

	fmt.Println("=== High-Level API Example ===")
	highLevelExample(input)

	fmt.Println("\n=== Low-Level API Example ===")
	lowLevelExample(input)
}

func highLevelExample(input string) {
	// Parse using high-level API
	result := sysml.ParseString(input)

	if !result.Success() {
		fmt.Printf("Parse errors:\n%s\n", result.Errors)
		return
	}

	fmt.Println("Parsing successful!")
	fmt.Printf("Found %d top-level packages\n", len(result.Model.Packages))

	// Walk the model
	fmt.Println("\nElements in model:")
	sysml.Walk(result.Model, func(elem sysml.Element, depth int) bool {
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}
		fmt.Printf("%s- %s: %s\n", indent, elem.Kind(), elem.Name())
		return true
	})

	// Find all parts
	parts := sysml.FindAll[*sysml.Part](result.Model)
	fmt.Printf("\nFound %d parts\n", len(parts))
	for _, part := range parts {
		defType := "usage"
		if part.IsDefinition {
			defType = "definition"
		}
		fmt.Printf("  - %s (%s)\n", part.Name(), defType)
	}

	// Count elements
	counter := sysml.NewCounter()
	sysml.Visit(result.Model, counter)
	fmt.Printf("\nElement counts: %v\n", counter.Counts)
}

func lowLevelExample(input string) {
	// Parse using low-level API
	tree, errors := low.Parse(input)

	if errors.HasErrors() {
		fmt.Printf("Parse errors:\n%s\n", errors)
		return
	}

	fmt.Println("Parsing successful!")

	// Create a parser to get token count
	parser := low.NewParser(input)
	fmt.Printf("Token count: %d\n", parser.TokenCount())

	// Access raw parse tree
	fmt.Printf("Parse tree type: %T\n", tree)

	// Example: Validate without building tree (faster)
	validationErrors := low.Validate(input)
	if validationErrors.HasErrors() {
		fmt.Printf("Validation errors: %s\n", validationErrors)
	} else {
		fmt.Println("Validation passed!")
	}
}
