package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dVoo/gosysml2/sysml"
)

// CategoryStats holds statistics for a single validation category.
type CategoryStats struct {
	Name            string  `json:"name"`
	Files           int     `json:"files"`
	Parsed          int     `json:"parsed"`
	Failed          int     `json:"failed"`
	Elements        int     `json:"elements"`
	LibraryImports  int     `json:"library_imports"`
	LibraryResolved int     `json:"library_resolved"`
	SuccessRate     float64 `json:"success_rate"`
}

// ValidationReport holds the complete validation results.
type ValidationReport struct {
	Timestamp     time.Time       `json:"timestamp"`
	TotalFiles    int             `json:"total_files"`
	TotalParsed   int             `json:"total_parsed"`
	TotalFailed   int             `json:"total_failed"`
	TotalElements int             `json:"total_elements"`
	LibraryLoaded bool            `json:"library_loaded"`
	Categories    []CategoryStats `json:"categories"`
	OverallRate   float64         `json:"overall_rate"`
	Duration      time.Duration   `json:"duration"`
}

func main() {
	var (
		libraryPath    = flag.String("library-path", "./libraries/sysml.library", "Path to standard library")
		validationPath = flag.String("validation-path", "./validationdata", "Path to validation data")
		verbose        = flag.Bool("verbose", false, "Show detailed per-file output")
		category       = flag.String("category", "", "Test specific category only (e.g., '01-Parts Tree')")
		jsonOutput     = flag.Bool("json", false, "Output results as JSON")
		compareFile    = flag.String("compare", "", "Compare with previous results file")
	)
	flag.Parse()

	start := time.Now()

	// Load standard library
	var registry *sysml.LibraryRegistry
	libraryLoaded := false

	if _, err := os.Stat(*libraryPath); err == nil {
		registry = sysml.NewLibraryRegistry(sysml.WithLibraryPaths(*libraryPath))
		if err := registry.RegisterStandardLibrary(); err == nil {
			libraryLoaded = true
			fmt.Printf("✓ Loaded standard library from %s\n", *libraryPath)
		} else {
			fmt.Fprintf(os.Stderr, "⚠ Warning: Failed to load standard library: %v\n", err)
			// Don't use partially loaded registry
			registry = nil
		}
	} else {
		fmt.Fprintf(os.Stderr, "⚠ Warning: Standard library not found at %s\n", *libraryPath)
	}

	// Collect validation files
	categories := make(map[string][]string)

	entries, err := os.ReadDir(*validationPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading validation data: %v\n", err)
		os.Exit(2)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		categoryName := entry.Name()

		// Filter by category if specified
		if *category != "" && categoryName != *category {
			continue
		}

		categoryPath := filepath.Join(*validationPath, categoryName)
		var files []string

		err := filepath.WalkDir(categoryPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".sysml") {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking category %s: %v\n", categoryName, err)
			continue
		}

		if len(files) > 0 {
			categories[categoryName] = files
		}
	}

	if len(categories) == 0 {
		fmt.Fprintf(os.Stderr, "No validation files found in %s\n", *validationPath)
		os.Exit(2)
	}

	// Process each category
	var report ValidationReport
	report.Timestamp = time.Now()
	report.LibraryLoaded = libraryLoaded

	for categoryName, files := range categories {
		stats := processCategory(categoryName, files, registry, *verbose)
		report.Categories = append(report.Categories, stats)
		report.TotalFiles += stats.Files
		report.TotalParsed += stats.Parsed
		report.TotalFailed += stats.Failed
		report.TotalElements += stats.Elements
	}

	report.Duration = time.Since(start)

	if report.TotalFiles > 0 {
		report.OverallRate = float64(report.TotalParsed) / float64(report.TotalFiles) * 100
	}

	// Handle comparison mode
	if *compareFile != "" {
		compareResults(&report, *compareFile)
	}

	// Output results
	if *jsonOutput {
		outputJSON(&report)
	} else {
		outputTable(&report)
	}

	// Exit with appropriate code
	if report.TotalFailed > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

// processCategory parses all files in a category and returns statistics.
func processCategory(name string, files []string, registry *sysml.LibraryRegistry, verbose bool) CategoryStats {
	stats := CategoryStats{
		Name:  name,
		Files: len(files),
	}

	if verbose {
		fmt.Printf("\n## %s\n", name)
	}

	for _, file := range files {
		filename := filepath.Base(file)

		var result *sysml.ParseResult
		if registry != nil {
			result = sysml.ParseFile(file, sysml.WithLibraryRegistry(registry))
		} else {
			result = sysml.ParseFile(file)
		}

		if result.Err() == nil {
			stats.Parsed++
			if verbose {
				fmt.Printf("  ✓ %s\n", filename)
			}

			if result.Model != nil {
				// Count elements
				count := 0
				sysml.Walk(result.Model, func(elem sysml.Element, depth int) bool {
					count++
					return true
				})
				stats.Elements += count

				// Count library imports
				for _, imp := range result.Model.Imports() {
					if imp.ResolvedPackage != nil {
						stats.LibraryImports++
						stats.LibraryResolved++
					} else if imp.ResolvedElement != nil {
						stats.LibraryImports++
					}
				}
			}
		} else {
			stats.Failed++
			if verbose {
				errMsg := "parse error"
				if errs := result.Errors(); len(errs) > 0 {
					errMsg = errs[0].Message
				}
				fmt.Printf("  ✗ %s - %s\n", filename, errMsg)
			}
		}
	}

	if stats.Files > 0 {
		stats.SuccessRate = float64(stats.Parsed) / float64(stats.Files) * 100
	}

	return stats
}

// outputTable prints results in a formatted table.
func outputTable(report *ValidationReport) {
	fmt.Println()
	fmt.Println("=== SysML Validation Test Results ===")
	fmt.Println()

	if report.LibraryLoaded {
		fmt.Println("✓ Standard library loaded")
	} else {
		fmt.Println("⚠ Standard library not loaded (limited resolution)")
	}

	fmt.Println()

	// Print table header
	fmt.Printf("%-30s | %5s | %6s | %6s | %12s\n", "Category", "Files", "Parsed", "Failed", "Success Rate")
	fmt.Println(strings.Repeat("-", 80))

	// Print each category
	for _, cat := range report.Categories {
		fmt.Printf("%-30s | %5d | %6d | %6d | %11.1f%%\n",
			cat.Name, cat.Files, cat.Parsed, cat.Failed, cat.SuccessRate)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-30s | %5d | %6d | %6d | %11.1f%%\n",
		"TOTAL", report.TotalFiles, report.TotalParsed, report.TotalFailed, report.OverallRate)

	fmt.Println()
	fmt.Printf("Total elements parsed: %d\n", report.TotalElements)

	// Library statistics
	var totalImports, totalResolved int
	for _, cat := range report.Categories {
		totalImports += cat.LibraryImports
		totalResolved += cat.LibraryResolved
	}

	if totalImports > 0 {
		fmt.Printf("Library imports: %d (%d resolved)\n", totalImports, totalResolved)
	}

	fmt.Printf("Duration: %v\n", report.Duration.Round(time.Millisecond))

	// Summary
	fmt.Println()
	if report.TotalFailed == 0 {
		fmt.Println("✓ All files parsed successfully!")
	} else {
		fmt.Printf("⚠ %d/%d files failed to parse\n", report.TotalFailed, report.TotalFiles)
	}
}

// outputJSON prints results as JSON.
func outputJSON(report *ValidationReport) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(report)
}

// compareResults compares current results with a previous run.
func compareResults(current *ValidationReport, compareFile string) {
	data, err := os.ReadFile(compareFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not read comparison file: %v\n", err)
		return
	}

	var previous ValidationReport
	if err := json.Unmarshal(data, &previous); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not parse comparison file: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("=== Comparison with Previous Run ===")
	fmt.Println()

	// Overall comparison
	prevRate := 0.0
	if previous.TotalFiles > 0 {
		prevRate = float64(previous.TotalParsed) / float64(previous.TotalFiles) * 100
	}

	diff := current.OverallRate - prevRate
	sign := "+"
	if diff < 0 {
		sign = ""
	}

	fmt.Printf("Overall success rate: %.1f%% → %.1f%% (%s%.1f%%)\n",
		prevRate, current.OverallRate, sign, diff)

	// Category comparisons
	fmt.Println("\nCategory changes:")
	for _, currCat := range current.Categories {
		found := false
		for _, prevCat := range previous.Categories {
			if prevCat.Name == currCat.Name {
				found = true
				if currCat.SuccessRate != prevCat.SuccessRate {
					diff := currCat.SuccessRate - prevCat.SuccessRate
					sign := "+"
					if diff < 0 {
						sign = ""
					}
					fmt.Printf("  %s: %.1f%% → %.1f%% (%s%.1f%%)\n",
						currCat.Name, prevCat.SuccessRate, currCat.SuccessRate, sign, diff)
				}
				break
			}
		}
		if !found {
			fmt.Printf("  %s: NEW (%.1f%%)\n", currCat.Name, currCat.SuccessRate)
		}
	}

	fmt.Println()
}
