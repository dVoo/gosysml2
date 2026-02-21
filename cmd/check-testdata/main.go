package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dVoo/gosysml2/sysml"
)

func main() {
	testdataDir := "docs/testdata"

	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		fmt.Printf("Error reading testdata directory: %v\n", err)
		os.Exit(1)
	}

	var sysmlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sysml") {
			sysmlFiles = append(sysmlFiles, entry.Name())
		}
	}

	fmt.Printf("Found %d .sysml files in %s\n\n", len(sysmlFiles), testdataDir)

	successCount := 0
	failCount := 0
	var failures []string

	for _, filename := range sysmlFiles {
		filePath := filepath.Join(testdataDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("✗ %s - Error reading file: %v\n", filename, err)
			failCount++
			failures = append(failures, fmt.Sprintf("%s (read error)", filename))
			continue
		}

		result := sysml.ParseString(string(content))
		if result.Err() == nil {
			if result.Model != nil {
				elementCount := countElements(result.Model)
				fmt.Printf("✓ %s - Parsed successfully (%d elements)\n", filename, elementCount)
			} else {
				fmt.Printf("✓ %s - Parsed successfully (no model)\n", filename)
			}
			successCount++
		} else {
			fmt.Printf("✗ %s - Parse failed: %v\n", filename, result.ParseError)
			failCount++
			failures = append(failures, filename)
		}
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf("Results: %d passed, %d failed\n", successCount, failCount)

	if failCount > 0 {
		fmt.Printf("\nFailed files:\n")
		for _, f := range failures {
			fmt.Printf("  - %s\n", f)
		}
		os.Exit(1)
	} else {
		fmt.Printf("\n✓ All %d testdata files parse correctly!\n", successCount)
	}
}

func countElements(model *sysml.Model) int {
	count := 0
	for range sysml.All(model) {
		count++
	}
	return count
}
