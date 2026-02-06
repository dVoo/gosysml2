// Parallel Parsing Example
//
// This example demonstrates:
// - Parse multiple files concurrently
// - Use ParseDirectoryParallel
// - Show worker configuration
// - Aggregate results from multiple files
//
// Run with: go run main.go
// Note: This example creates temporary files for demonstration

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dVoo/gosysml2_oc/sysml"
)

func main() {
	fmt.Println("=== Parallel Parsing Example ===")
	fmt.Println()

	// Create temporary directory with sample files
	tmpDir, err := os.MkdirTemp("", "sysml-parallel-example-*")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Printf("1. Created temporary directory: %s\n", tmpDir)

	// Create sample SysML files
	files := []struct {
		name    string
		content string
	}{
		{
			name: "vehicle.sysml",
			content: `package VehicleModel {
    part def Vehicle {
        part engine : Engine;
        part wheels : Wheel[4];
    }

    part def Engine {
        attribute horsepower : Real;
    }

    part def Wheel;
}`,
		},
		{
			name: "requirements.sysml",
			content: `package RequirementsModel {
    requirement def REQ001 {
        doc /* The system shall provide user authentication. */
    }

    requirement def REQ002 {
        doc /* The system shall encrypt all data. */
    }

    verification def VerifyREQ001;
}`,
		},
		{
			name: "components.sysml",
			content: `package ComponentLibrary {
    part def Sensor {
        attribute accuracy : Real;
    }

    part def Actuator {
        attribute power : Real;
    }

    part def Controller {
        part sensors : Sensor[2];
        part actuator : Actuator;
    }
}`,
		},
		{
			name: "interfaces.sysml",
			content: `package InterfaceDefinitions {
    port def DataPort;
    port def ControlPort;

    interface def DataInterface {
        end p1 : DataPort;
        end p2 : DataPort;
    }
}`,
		},
	}

	// Write sample files
	fmt.Println("2. Creating sample SysML files...")
	for _, f := range files {
		path := filepath.Join(tmpDir, f.name)
		if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
			fmt.Printf("Failed to write %s: %v\n", f.name, err)
			os.Exit(1)
		}
		fmt.Printf("   Created: %s\n", f.name)
	}

	// Sequential parsing
	fmt.Println("\n3. Sequential parsing (baseline)...")
	start := time.Now()
	seqResults, err := sysml.ParseDirectory(tmpDir)
	if err != nil {
		fmt.Printf("   Sequential parsing failed: %v\n", err)
		os.Exit(1)
	}
	seqDuration := time.Since(start)
	fmt.Printf("   Parsed %d files in %v\n", len(seqResults), seqDuration)

	// Analyze sequential results
	var seqSuccess int
	var seqErrors int
	var seqTotalElements int
	for _, r := range seqResults {
		if r.Success() {
			seqSuccess++
			if r.Model != nil {
				counter := sysml.NewCounter()
				sysml.Visit(r.Model, counter)
				seqTotalElements += counter.Total()
			}
		} else {
			seqErrors++
		}
	}
	fmt.Printf("   Success: %d, Errors: %d, Total elements: %d\n", seqSuccess, seqErrors, seqTotalElements)

	// Parallel parsing with default workers (NumCPU)
	fmt.Println("\n4. Parallel parsing (using all CPU cores)...")
	fmt.Printf("   CPU cores available: %d\n", runtime.NumCPU())
	start = time.Now()
	parResults, err := sysml.ParseDirectoryParallel(tmpDir, 0) // 0 means use NumCPU
	if err != nil {
		fmt.Printf("   Parallel parsing failed: %v\n", err)
		os.Exit(1)
	}
	parDuration := time.Since(start)
	fmt.Printf("   Parsed %d files in %v\n", len(parResults), parDuration)

	// Analyze parallel results
	var parSuccess int
	var parErrors int
	var parTotalElements int
	for _, r := range parResults {
		if r.Success() {
			parSuccess++
			if r.Model != nil {
				counter := sysml.NewCounter()
				sysml.Visit(r.Model, counter)
				parTotalElements += counter.Total()
			}
		} else {
			parErrors++
		}
	}
	fmt.Printf("   Success: %d, Errors: %d, Total elements: %d\n", parSuccess, parErrors, parTotalElements)

	// Compare performance
	fmt.Println("\n5. Performance comparison:")
	if seqDuration > 0 && parDuration > 0 {
		speedup := float64(seqDuration) / float64(parDuration)
		fmt.Printf("   Sequential: %v\n", seqDuration)
		fmt.Printf("   Parallel:   %v\n", parDuration)
		fmt.Printf("   Speedup:    %.2fx\n", speedup)
	}

	// Parallel parsing with specific worker count
	fmt.Println("\n6. Parallel parsing with 2 workers...")
	start = time.Now()
	results2, err := sysml.ParseDirectoryParallel(tmpDir, 2)
	if err != nil {
		fmt.Printf("   Parsing failed: %v\n", err)
	} else {
		duration2 := time.Since(start)
		fmt.Printf("   Parsed %d files in %v\n", len(results2), duration2)
	}

	// Show aggregated statistics
	fmt.Println("\n7. Aggregated statistics from parallel parsing:")
	totalCounts := make(map[sysml.ElementKind]int)
	packageNames := []string{}
	requirementCount := 0
	partCount := 0

	for _, r := range parResults {
		if !r.Success() || r.Model == nil {
			continue
		}

		// Count elements
		counter := sysml.NewCounter()
		sysml.Visit(r.Model, counter)
		for kind, count := range counter.Counts {
			totalCounts[kind] += count
		}

		// Collect package names
		for _, pkg := range r.Model.Packages {
			packageNames = append(packageNames, pkg.Name())
		}

		// Count specific types
		requirementCount += len(sysml.FindRequirements(r.Model))
		partCount += len(sysml.FindParts(r.Model))
	}

	fmt.Printf("   Total packages found: %d\n", len(packageNames))
	fmt.Printf("   Package names: %v\n", packageNames)
	fmt.Printf("   Total requirements: %d\n", requirementCount)
	fmt.Printf("   Total parts: %d\n", partCount)
	fmt.Println("   Element counts by kind:")
	for kind, count := range totalCounts {
		if count > 0 {
			fmt.Printf("      %s: %d\n", kind, count)
		}
	}

	// Show per-file results
	fmt.Println("\n8. Per-file results:")
	for _, r := range parResults {
		status := "✓"
		if !r.Success() {
			status = "✗"
		}
		elemCount := 0
		if r.Model != nil {
			counter := sysml.NewCounter()
			sysml.Visit(r.Model, counter)
			elemCount = counter.Total()
		}
		fmt.Printf("   %s %s: %d elements", status, filepath.Base(r.Source), elemCount)
		if !r.Success() {
			fmt.Printf(" (%d errors)", len(r.Errors.Errors))
		}
		fmt.Println()
	}

	// Demonstrate streaming (most memory efficient)
	fmt.Println("\n9. Streaming parse (memory efficient)...")
	elemCount := 0
	start = time.Now()
	err = sysml.ParseDirectoryStream(tmpDir, func(r *sysml.ParseResult) error {
		if r.Success() && r.Model != nil {
			counter := sysml.NewCounter()
			sysml.Visit(r.Model, counter)
			elemCount += counter.Total()
		}
		return nil
	})
	streamDuration := time.Since(start)
	if err != nil {
		fmt.Printf("   Streaming parse failed: %v\n", err)
	} else {
		fmt.Printf("   Streamed %d files in %v\n", len(files), streamDuration)
		fmt.Printf("   Total elements processed: %d\n", elemCount)
	}

	// Memory optimization note
	fmt.Println("\n10. Memory optimization options:")
	fmt.Println("    - Use ParseDirectory() for small repositories (< 100 MB)")
	fmt.Println("    - Use ParseDirectoryParallel() for multi-core machines")
	fmt.Println("    - Use ParseDirectoryStream() for very large repositories")
	fmt.Println("    - Use WithDiscardTree() option to reduce memory by ~30%")

	fmt.Println("\n=== Example Complete ===")
}
