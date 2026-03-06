// Validation Example
//
// This example demonstrates:
// - Validate SysML syntax without full parsing
// - Show error handling for invalid input
// - Demonstrate error location reporting
//
// Run with: go run main.go

package main

import (
	"fmt"
	"os"

	"github.com/dVoo/gosysml2/low"
	"github.com/dVoo/gosysml2/sysml"
)

func main() {
	fmt.Println("=== Validation Example ===")
	fmt.Println()

	// Example 1: Valid SysML input
	fmt.Println("1. Validating valid SysML input...")
	validInput := `
package ValidModel {
    part def Vehicle {
        part engine : Engine;
    }

    part def Engine {
        attribute horsepower : Real;
    }
}
`

	// Using low-level API for validation only (faster)
	errors := low.Validate(validInput)
	if errors.HasErrors() {
		fmt.Printf("   Validation errors:\n%s\n", errors)
	} else {
		fmt.Println("   Validation passed!")
	}

	// Example 2: Syntax-only validation through the high-level API
	fmt.Println("\n2. Syntax-only parsing without model build...")
	fastResult := sysml.ParseString(validInput,
		sysml.WithoutModelBuild(),
		sysml.WithDiscardTree(),
	)
	if fastResult.Err() != nil {
		fmt.Printf("   Unexpected syntax errors: %v\n", fastResult.Err())
	} else {
		fmt.Println("   Syntax-only parse passed!")
		fmt.Printf("   Model built: %t\n", fastResult.Model != nil)
		fmt.Printf("   Parse tree retained: %t\n", fastResult.Tree != nil)
	}

	// Example 3: Invalid SysML input
	fmt.Println("\n3. Validating invalid SysML input...")
	invalidInput := `
package InvalidModel {
    part def Vehicle {
        part engine : Engine
        // Missing semicolon and closing brace
    
    part def Engine {
        attribute horsepower Real;
        // Missing colon
    }
}
`

	errors = low.Validate(invalidInput)
	if errors.HasErrors() {
		allErrors := errors.All()
		fmt.Printf("   Validation found %d error(s):\n", len(allErrors))
		for i, err := range allErrors {
			fmt.Printf("   Error %d:\n", i+1)
			fmt.Printf("      Line: %d\n", err.Line)
			fmt.Printf("      Column: %d\n", err.Column)
			fmt.Printf("      Message: %s\n", err.Message)
			fmt.Printf("      Source: %s\n", err.Source)
		}
	} else {
		fmt.Println("   Validation passed!")
	}

	// Example 4: High-level API with error handling
	fmt.Println("\n4. High-level API error handling...")

	// Parse with errors
	result := sysml.ParseString(invalidInput)
	if result.Err() != nil {
		fmt.Printf("   Parse failed with %d error(s)\n", len(result.Errors()))

		// Get first error
		if first := result.ParseError.First(); first != nil {
			fmt.Printf("   First error at line %d, column %d:\n",
				first.Line, first.Column)
			fmt.Printf("      %s\n", first.Message)
		}

		// Show all errors
		fmt.Println("   All errors:")
		for _, err := range result.Errors() {
			fmt.Printf("      Line %d:%d - %s\n",
				err.Line, err.Column, err.Message)
		}
	}

	// Example 5: Valid input with full parsing
	fmt.Println("\n5. Parsing valid input...")
	result = sysml.ParseString(validInput)
	if result.Err() == nil {
		fmt.Println("   Parse successful!")
		fmt.Printf("   Model has %d package(s)\n", len(result.Model.Packages()))
		for _, pkg := range result.Model.Packages() {
			fmt.Printf("      - %s\n", pkg.Name())
		}
	}

	// Example 6: File validation
	fmt.Println("\n6. File validation...")

	// Create a temporary file with invalid content
	tmpFile := "/tmp/invalid_test.sysml"
	invalidContent := []byte(`
package Test {
    part def Broken {
        attribute value
    }
}
`)
	if err := os.WriteFile(tmpFile, invalidContent, 0644); err != nil {
		fmt.Printf("   Failed to create test file: %v\n", err)
	} else {
		defer os.Remove(tmpFile)

		// Validate the file
		result := sysml.ParseFile(tmpFile)
		if result.Err() != nil {
			fmt.Printf("   File validation found %d error(s)\n", len(result.Errors()))
			for _, err := range result.Errors() {
				fmt.Printf("      Line %d:%d - %s\n",
					err.Line, err.Column, err.Message)
			}
		}
	}

	// Example 7: Error formatting
	fmt.Println("\n7. Error formatting options...")
	result = sysml.ParseString(invalidInput)
	if result.Err() != nil {
		fmt.Println("   Formatted error string:")
		fmt.Printf("   %s\n", result.ParseError.Error())
	}

	// Example 8: Checking specific error types
	fmt.Println("\n8. Error analysis...")
	result = sysml.ParseString(invalidInput)
	if result.Err() != nil {
		syntaxErrors := 0
		otherErrors := 0

		for _, err := range result.Errors() {
			// Check if it's a syntax error
			if err.Line >= 0 {
				syntaxErrors++
			} else {
				otherErrors++
			}
		}

		fmt.Printf("   Syntax errors: %d\n", syntaxErrors)
		fmt.Printf("   Other errors: %d\n", otherErrors)
	}

	// Example 9: Building a validation report
	fmt.Println("\n9. Validation report...")
	inputs := []struct {
		name  string
		input string
	}{
		{
			name: "Valid model",
			input: `
package Valid {
    part def Component;
}
`,
		},
		{
			name: "Missing semicolon",
			input: `
package Invalid {
    part def Component
}
`,
		},
		{
			name: "Invalid syntax",
			input: `
package Invalid {
    invalid keyword here
}
`,
		},
	}

	fmt.Println("   Running validation report...")
	for _, test := range inputs {
		errors := low.Validate(test.input)
		status := "PASS"
		if errors.HasErrors() {
			status = "FAIL"
		}
		fmt.Printf("   [%s] %s (%d errors)\n", status, test.name, len(errors.All()))
	}

	fmt.Println("\n=== Example Complete ===")
}
