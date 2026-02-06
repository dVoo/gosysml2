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

	"github.com/dVoo/gosysml2_oc/low"
	"github.com/dVoo/gosysml2_oc/sysml"
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

	// Example 2: Invalid SysML input
	fmt.Println("\n2. Validating invalid SysML input...")
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

	// Example 3: High-level API with error handling
	fmt.Println("\n3. High-level API error handling...")

	// Parse with errors
	result := sysml.ParseString(invalidInput)
	if !result.Success() {
		fmt.Printf("   Parse failed with %d error(s)\n", len(result.Errors.Errors))

		// Get first error
		if first := result.Errors.First(); first != nil {
			fmt.Printf("   First error at line %d, column %d:\n",
				first.Line, first.Column)
			fmt.Printf("      %s\n", first.Message)
		}

		// Show all errors
		fmt.Println("   All errors:")
		for _, err := range result.Errors.Errors {
			fmt.Printf("      Line %d:%d - %s\n",
				err.Line, err.Column, err.Message)
		}
	}

	// Example 4: Valid input with full parsing
	fmt.Println("\n4. Parsing valid input...")
	result = sysml.ParseString(validInput)
	if result.Success() {
		fmt.Println("   Parse successful!")
		fmt.Printf("   Model has %d package(s)\n", len(result.Model.Packages))
		for _, pkg := range result.Model.Packages {
			fmt.Printf("      - %s\n", pkg.Name())
		}
	}

	// Example 5: File validation
	fmt.Println("\n5. File validation...")

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
		if !result.Success() {
			fmt.Printf("   File validation found %d error(s)\n", len(result.Errors.Errors))
			for _, err := range result.Errors.Errors {
				fmt.Printf("      Line %d:%d - %s\n",
					err.Line, err.Column, err.Message)
			}
		}
	}

	// Example 6: Error formatting
	fmt.Println("\n6. Error formatting options...")
	result = sysml.ParseString(invalidInput)
	if !result.Success() {
		fmt.Println("   Formatted error string:")
		fmt.Printf("   %s\n", result.Errors.Error())
	}

	// Example 7: Checking specific error types
	fmt.Println("\n7. Error analysis...")
	result = sysml.ParseString(invalidInput)
	if !result.Success() {
		syntaxErrors := 0
		otherErrors := 0

		for _, err := range result.Errors.Errors {
			// Check if it's a syntax error
			if err.Line > 0 {
				syntaxErrors++
			} else {
				otherErrors++
			}
		}

		fmt.Printf("   Syntax errors: %d\n", syntaxErrors)
		fmt.Printf("   Other errors: %d\n", otherErrors)
	}

	// Example 8: Building a validation report
	fmt.Println("\n8. Validation report...")
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
