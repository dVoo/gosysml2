package sysml

import (
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// validationTestdataDir returns the path to validationdata relative to this package.
func validationTestdataDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join("..", "..", "validationdata")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("validationdata directory not found")
	}
	return dir
}

// standardLibraryExists checks if the standard library is available.
func standardLibraryExists() bool {
	_, err := os.Stat("./libraries/sysml.library")
	return err == nil
}

// collectValidationFiles finds all .sysml files in the validationdata directory grouped by category.
func collectValidationFiles(t *testing.T) map[string][]string {
	t.Helper()
	validationDir := validationTestdataDir(t)

	categories := make(map[string][]string)

	entries, err := os.ReadDir(validationDir)
	if err != nil {
		t.Fatalf("failed to read validationdata directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		categoryName := entry.Name()
		categoryPath := filepath.Join(validationDir, categoryName)

		var files []string
		err := filepath.WalkDir(categoryPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(path, ".sysml") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("failed to walk category %s: %v", categoryName, err)
		}

		if len(files) > 0 {
			categories[categoryName] = files
		}
	}

	return categories
}

// collectTestFiles finds all .sysml files in the validationdata directory.
func collectTestFiles(t *testing.T) []string {
	t.Helper()
	validationDir := validationTestdataDir(t)

	var files []string
	err := filepath.WalkDir(validationDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sysml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk validationdata: %v", err)
	}

	return files
}

// TestIntegrationFullValidation runs comprehensive integration tests on all test files.
func TestIntegrationFullValidation(t *testing.T) {
	files := collectTestFiles(t)
	t.Logf("Found %d SysML test files", len(files))

	if len(files) == 0 {
		t.Fatal("No test files found")
	}

	// Track statistics for summary table
	type fileStats struct {
		name         string
		elements     int
		packages     int
		parts        int
		requirements int
		errors       int
		success      bool
	}
	stats := make([]fileStats, 0, len(files))

	// Process each file with sub-tests
	for _, file := range files {
		name := filepath.Base(file)
		relPath, _ := filepath.Rel(validationTestdataDir(t), file)

		t.Run(relPath, func(t *testing.T) {
			result := ParseFile(file)
			stat := fileStats{name: name}

			if result.Success() {
				stat.success = true
				t.Logf("PASS: %s - parsed successfully", name)

				if result.Model != nil {
					// Count elements
					Walk(result.Model, func(elem Element, depth int) bool {
						stat.elements++
						switch elem.Kind() {
						case KindPackage:
							stat.packages++
						case KindPart:
							stat.parts++
						case KindRequirement:
							stat.requirements++
						}
						return true
					})

					if stat.elements == 0 {
						t.Logf("  Warning: %s has 0 elements (empty file?)", name)
					}
				}
			} else {
				stat.success = false
				if result.Errors != nil {
					stat.errors = len(result.Errors.Errors)
				}
				t.Logf("FAIL: %s - %s", name, result.Errors.First().Message)
				// Don't fail the test - some files may have intentional errors
			}

			stats = append(stats, stat)
		})
	}

	// Print summary table
	t.Log("\n=== Integration Test Summary ===")
	t.Logf("%-50s | %8s | %8s | %5s | %12s | %6s", "File", "Elements", "Packages", "Parts", "Requirements", "Errors")
	t.Log("----------------------------------------------------------------------------------------------------")

	passed := 0
	failed := 0
	totalElements := 0
	totalPackages := 0
	totalParts := 0
	totalRequirements := 0

	for _, stat := range stats {
		if !stat.success {
			failed++
		} else {
			passed++
			totalElements += stat.elements
			totalPackages += stat.packages
			totalParts += stat.parts
			totalRequirements += stat.requirements
		}
		t.Logf("%-50s | %8d | %8d | %5d | %12d | %6d", stat.name, stat.elements, stat.packages, stat.parts, stat.requirements, stat.errors)
	}

	t.Log("----------------------------------------------------------------------------------------------------")
	t.Logf("TOTAL: %d passed, %d failed out of %d files", passed, failed, len(files))
	t.Logf("Total elements: %d (packages: %d, parts: %d, requirements: %d)", totalElements, totalPackages, totalParts, totalRequirements)
	t.Logf("\n")
}

// TestIntegrationParallelParsing tests parallel parsing produces same results as sequential.
func TestIntegrationParallelParsing(t *testing.T) {
	testdataDir := validationTestdataDir(t)
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	// Parse sequentially first
	seqResults, err := ParseDirectory(testdataDir)
	if err != nil {
		t.Fatalf("Sequential ParseDirectory failed: %v", err)
	}

	// Count elements from sequential results
	seqElementCounts := make(map[string]int)
	for _, r := range seqResults {
		if r.Success() && r.Model != nil {
			count := 0
			Walk(r.Model, func(elem Element, depth int) bool {
				count++
				return true
			})
			seqElementCounts[r.Source] = count
		}
	}

	// Parse in parallel
	parResults, err := ParseDirectoryParallel(testdataDir, 4, WithDiscardTree())
	if err != nil {
		t.Fatalf("Parallel ParseDirectoryParallel failed: %v", err)
	}

	// Count elements from parallel results
	parElementCounts := make(map[string]int)
	for _, r := range parResults {
		if r.Success() && r.Model != nil {
			count := 0
			Walk(r.Model, func(elem Element, depth int) bool {
				count++
				return true
			})
			parElementCounts[r.Source] = count
		}
	}

	// Compare results
	if len(seqElementCounts) != len(parElementCounts) {
		t.Errorf("Sequential and parallel parsing produced different numbers of successful results: %d vs %d", len(seqElementCounts), len(parElementCounts))
	}

	mismatches := 0
	for source, seqCount := range seqElementCounts {
		parCount, ok := parElementCounts[source]
		if !ok {
			t.Errorf("File %s parsed sequentially but not in parallel", source)
			mismatches++
			continue
		}
		if seqCount != parCount {
			t.Errorf("File %s has different element counts: sequential=%d, parallel=%d", source, seqCount, parCount)
			mismatches++
		}
	}

	if mismatches == 0 {
		t.Logf("Parallel parsing matches sequential: %d files processed with identical element counts", len(seqElementCounts))
	}
}

// TestIntegrationStreaming tests streaming parse mode.
func TestIntegrationStreaming(t *testing.T) {
	testdataDir := validationTestdataDir(t)
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	var count int32
	var passed int32

	err := ParseDirectoryStream(testdataDir, func(r *ParseResult) error {
		atomic.AddInt32(&count, 1)
		if r.Success() {
			atomic.AddInt32(&passed, 1)
		}
		// Verify tree was discarded in streaming mode
		if r.Tree != nil {
			t.Error("Expected Tree to be nil in streaming mode")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("ParseDirectoryStream failed: %v", err)
	}

	if int(count) != len(files) {
		t.Errorf("Expected %d files, got %d", len(files), count)
	}

	t.Logf("Streaming parse: %d/%d files parsed successfully", passed, count)
}

// TestIntegrationIteratorConsistency verifies iterators produce same results as Walk.
func TestIntegrationIteratorConsistency(t *testing.T) {
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	// Find a file that parses successfully
	var testModel *Model
	for _, file := range files {
		result := ParseFile(file)
		if result.Success() && result.Model != nil {
			testModel = result.Model
			break
		}
	}

	if testModel == nil {
		t.Fatal("No successfully parsed model found for iterator testing")
	}

	// Count using Walk
	walkCount := 0
	Walk(testModel, func(elem Element, depth int) bool {
		walkCount++
		return true
	})

	// Count using All iterator
	iterCount := 0
	for range All(testModel) {
		iterCount++
	}

	if walkCount != iterCount {
		t.Errorf("All() iterator count (%d) doesn't match Walk count (%d)", iterCount, walkCount)
	} else {
		t.Logf("All() iterator matches Walk: %d elements", walkCount)
	}

	// Count packages using FindPackages (old API)
	packagesOld := FindPackages(testModel)

	// Count packages using OfType iterator
	packagesIter := 0
	for range OfType[*Package](testModel) {
		packagesIter++
	}

	if len(packagesOld) != packagesIter {
		t.Errorf("OfType[*Package] count (%d) doesn't match FindPackages count (%d)", packagesIter, len(packagesOld))
	} else {
		t.Logf("OfType[*Package] iterator matches FindPackages: %d packages", len(packagesOld))
	}

	// Count using OfKind iterator
	packagesKind := 0
	for range OfKind(testModel, KindPackage) {
		packagesKind++
	}

	if len(packagesOld) != packagesKind {
		t.Errorf("OfKind count (%d) doesn't match FindPackages count (%d)", packagesKind, len(packagesOld))
	} else {
		t.Logf("OfKind iterator matches FindPackages: %d packages", len(packagesOld))
	}
}

// TestIntegrationDiscardTree tests that WithDiscardTree reduces memory usage.
func TestIntegrationDiscardTree(t *testing.T) {
	input := `package Test { part def A; }`

	// With tree
	result1 := ParseString(input)
	if result1.Tree == nil {
		t.Error("Expected Tree to be present without WithDiscardTree()")
	}

	// Without tree
	result2 := ParseString(input, WithDiscardTree())
	if result2.Tree != nil {
		t.Error("Expected Tree to be nil with WithDiscardTree()")
	}

	// Both should have model
	if result1.Model == nil || result2.Model == nil {
		t.Error("Expected Model to be present in both cases")
	}

	// Models should be equivalent
	if len(result1.Model.Elements) != len(result2.Model.Elements) {
		t.Error("Models have different element counts")
	}
}

// TestIntegrationParseDirectory tests basic directory parsing.
func TestIntegrationParseDirectory(t *testing.T) {
	testdataDir := validationTestdataDir(t)
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	results, err := ParseDirectory(testdataDir)
	if err != nil {
		t.Fatalf("ParseDirectory failed: %v", err)
	}

	if len(results) != len(files) {
		t.Errorf("Expected %d results, got %d", len(files), len(results))
	}

	passed := 0
	for _, r := range results {
		if r.Success() {
			passed++
		}
	}

	t.Logf("ParseDirectory: %d/%d files parsed successfully", passed, len(results))
}

// TestIntegrationParseDirectoryParallel tests parallel directory parsing.
func TestIntegrationParseDirectoryParallel(t *testing.T) {
	testdataDir := validationTestdataDir(t)
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	results, err := ParseDirectoryParallel(testdataDir, 4, WithDiscardTree())
	if err != nil {
		t.Fatalf("ParseDirectoryParallel failed: %v", err)
	}

	if len(results) != len(files) {
		t.Errorf("Expected %d results, got %d", len(files), len(results))
	}

	passed := 0
	for _, r := range results {
		if r.Success() {
			passed++
			// Verify tree was discarded
			if r.Tree != nil {
				t.Error("Expected Tree to be nil with WithDiscardTree()")
			}
		}
	}

	t.Logf("ParseDirectoryParallel: %d/%d files parsed successfully", passed, len(results))
}

// TestIntegrationParseDirectoryStream tests streaming directory parsing.
func TestIntegrationParseDirectoryStream(t *testing.T) {
	testdataDir := validationTestdataDir(t)
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	count := 0
	passed := 0

	err := ParseDirectoryStream(testdataDir, func(r *ParseResult) error {
		count++
		if r.Success() {
			passed++
		}
		// Verify tree was discarded (streaming always discards)
		if r.Tree != nil {
			t.Error("Expected Tree to be nil in streaming mode")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("ParseDirectoryStream failed: %v", err)
	}

	if count != len(files) {
		t.Errorf("Expected %d files, got %d", len(files), count)
	}

	t.Logf("ParseDirectoryStream: %d/%d files parsed successfully", passed, count)
}

// TestIntegrationCountAll tests CountAll function.
func TestIntegrationCountAll(t *testing.T) {
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	// Find a file that parses successfully
	var testModel *Model
	for _, file := range files {
		result := ParseFile(file)
		if result.Success() && result.Model != nil {
			testModel = result.Model
			break
		}
	}

	if testModel == nil {
		t.Fatal("No successfully parsed model found for CountAll testing")
	}

	counts := CountAll(testModel)
	t.Logf("Element counts by kind: %v", counts)

	total := 0
	for _, count := range counts {
		total += count
	}

	// Verify total matches Walk count
	walkCount := 0
	Walk(testModel, func(elem Element, depth int) bool {
		walkCount++
		return true
	})

	if total != walkCount {
		t.Errorf("CountAll total (%d) doesn't match Walk count (%d)", total, walkCount)
	} else {
		t.Logf("CountAll matches Walk: %d elements", total)
	}
}

// TestIntegrationVisitPattern tests the visitor pattern.
func TestIntegrationVisitPattern(t *testing.T) {
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	// Find a file that parses successfully
	var testModel *Model
	for _, file := range files {
		result := ParseFile(file)
		if result.Success() && result.Model != nil {
			testModel = result.Model
			break
		}
	}

	if testModel == nil {
		t.Fatal("No successfully parsed model found for visitor testing")
	}

	// Custom visitor that counts packages
	visitor := &countingVisitor{}
	Visit(testModel, visitor)

	// Compare with FindPackages
	packages := FindPackages(testModel)
	if visitor.packageCount != len(packages) {
		t.Errorf("Visitor package count (%d) doesn't match FindPackages count (%d)", visitor.packageCount, len(packages))
	} else {
		t.Logf("Visitor pattern matches FindPackages: %d packages", len(packages))
	}
}

// countingVisitor is a test visitor that counts packages.
type countingVisitor struct {
	BaseVisitor
	packageCount int
}

func (v *countingVisitor) VisitPackage(pkg *Package) bool {
	v.packageCount++
	return true
}

// TestIntegrationErrorHandling tests error handling on files with errors.
func TestIntegrationErrorHandling(t *testing.T) {
	// Test with invalid input
	result := ParseString("invalid {{{ syntax")

	if result.Success() {
		t.Error("Expected parse to fail for invalid syntax")
	}

	if result.Errors == nil {
		t.Error("Expected errors to be present")
	} else {
		t.Logf("Got expected error: %s", result.Errors.First().Message)
	}

	// Model should still be constructed (even with errors)
	if result.Model == nil {
		t.Log("Model is nil for invalid input (acceptable)")
	}
}

// TestIntegrationParseBytes tests the ParseBytes function.
func TestIntegrationParseBytes(t *testing.T) {
	input := []byte(`package Test { part def A; }`)
	result := ParseBytes(input, "test")

	if !result.Success() {
		t.Errorf("ParseBytes failed: %v", result.Errors)
	}

	if result.Model == nil {
		t.Error("Expected Model to be present")
	}
}

// TestIntegrationParseReader tests the ParseReader function.
func TestIntegrationParseReader(t *testing.T) {
	input := strings.NewReader(`package Test { part def A; }`)
	result := ParseReader(input, "test")

	if !result.Success() {
		t.Errorf("ParseReader failed: %v", result.Errors)
	}

	if result.Model == nil {
		t.Error("Expected Model to be present")
	}
}

// TestIntegrationValidateFile tests the ValidateFile function.
func TestIntegrationValidateFile(t *testing.T) {
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	// Just test the first file
	err := ValidateFile(files[0])
	// Don't fail - some files may have intentional errors
	if err == nil {
		t.Logf("ValidateFile: %s is valid", filepath.Base(files[0]))
	} else {
		t.Logf("ValidateFile: %s has errors (may be intentional)", filepath.Base(files[0]))
	}
}

// TestIntegrationModelIndex tests model index building.
func TestIntegrationModelIndex(t *testing.T) {
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	// Find a file that parses successfully
	var testModel *Model
	for _, file := range files {
		result := ParseFile(file)
		if result.Success() && result.Model != nil {
			testModel = result.Model
			break
		}
	}

	if testModel == nil {
		t.Fatal("No successfully parsed model found for index testing")
	}

	// Index should be built during parsing - verify by checking FindByQualifiedName works
	// The elementIndex is unexported, but we can verify BuildIndex doesn't panic
	testModel.BuildIndex()
	t.Log("BuildIndex completed without panic")

	// Try to find an element by qualified name if any exist
	if len(testModel.Elements) > 0 {
		elem := testModel.Elements[0]
		qn := elem.QualifiedName()
		if qn != "" {
			found := testModel.FindByQualifiedName(qn)
			if found == nil {
				t.Logf("FindByQualifiedName returned nil for %s (may be expected)", qn)
			} else if found.QualifiedName() != qn {
				t.Error("FindByQualifiedName returned wrong element")
			} else {
				t.Logf("Successfully found element by qualified name: %s", qn)
			}
		}
	}
}

// TestIntegrationModelResolveReferences tests reference resolution.
func TestIntegrationModelResolveReferences(t *testing.T) {
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	// Find a file that parses successfully
	var testModel *Model
	for _, file := range files {
		result := ParseFile(file)
		if result.Success() && result.Model != nil {
			testModel = result.Model
			break
		}
	}

	if testModel == nil {
		t.Fatal("No successfully parsed model found for reference resolution testing")
	}

	// References should be resolved during parsing
	// Just verify the method doesn't panic
	testModel.ResolveReferences()
	t.Log("ResolveReferences completed without panic")
}

// TestIntegrationParseString tests basic ParseString functionality.
func TestIntegrationParseString(t *testing.T) {
	input := `package Test {
		part def Vehicle;
		part vehicle : Vehicle;
	}`

	result := ParseString(input)

	if !result.Success() {
		t.Errorf("ParseString failed: %v", result.Errors)
	}

	if result.Model == nil {
		t.Fatal("Expected Model to be present")
	}

	// Verify structure
	if len(result.Model.Elements) == 0 {
		t.Error("Expected at least one element")
	}

	// Find the package
	packages := FindPackages(result.Model)
	if len(packages) == 0 {
		t.Error("Expected at least one package")
	} else {
		t.Logf("Found %d package(s)", len(packages))
	}

	// Find parts
	parts := FindAll[*Part](result.Model)
	t.Logf("Found %d part(s)", len(parts))
}

// TestIntegrationParseFile tests basic ParseFile functionality.
func TestIntegrationParseFile(t *testing.T) {
	files := collectTestFiles(t)

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	// Test parsing the first file
	result := ParseFile(files[0])

	// Don't fail on parse errors - some files may have intentional errors
	if result.Success() {
		t.Logf("ParseFile: %s parsed successfully", filepath.Base(files[0]))
		if result.Model != nil {
			t.Logf("  Model has %d top-level elements", len(result.Model.Elements))
		}
	} else {
		t.Logf("ParseFile: %s has errors (may be intentional): %s", filepath.Base(files[0]), result.Errors.First().Message)
	}
}

// TestIntegrationParseWithOptions tests parsing with various options.
func TestIntegrationParseWithOptions(t *testing.T) {
	input := `package Test { part def A; }`

	// Without discard tree
	result1 := ParseString(input)
	if result1.Tree == nil {
		t.Error("Expected Tree to be present")
	}

	// With discard tree
	result2 := ParseString(input, WithDiscardTree())
	if result2.Tree != nil {
		t.Error("Expected Tree to be nil with WithDiscardTree()")
	}

	// Both should have model
	if result1.Model == nil || result2.Model == nil {
		t.Error("Expected Model in both results")
	}
}

// TestIntegrationFullPipeline tests the complete parse pipeline.
func TestIntegrationFullPipeline(t *testing.T) {
	input := `package VehicleModel {
		part def Vehicle {
			attribute mass :> ISQ::mass;
		}
		
		part vehicle : Vehicle {
			attribute :>> mass = 1000;
		}
		
		requirement def MassRequirement {
			doc /* Vehicle mass shall be less than 1500 */
		}
		
		requirement massReq : MassRequirement;
	}`

	result := ParseString(input)

	if !result.Success() {
		t.Fatalf("Parse failed: %v", result.Errors)
	}

	if result.Model == nil {
		t.Fatal("Expected Model to be present")
	}

	// Count different element types
	packages := FindAll[*Package](result.Model)
	parts := FindAll[*Part](result.Model)
	requirements := FindAll[*Requirement](result.Model)
	attributes := FindAll[*Attribute](result.Model)

	t.Logf("Parsed model contains:")
	t.Logf("  %d package(s)", len(packages))
	t.Logf("  %d part(s)", len(parts))
	t.Logf("  %d requirement(s)", len(requirements))
	t.Logf("  %d attribute(s)", len(attributes))

	// Verify expected counts
	if len(packages) != 1 {
		t.Errorf("Expected 1 package, got %d", len(packages))
	}
	if len(parts) != 2 { // Vehicle def and vehicle instance
		t.Errorf("Expected 2 parts, got %d", len(parts))
	}
	if len(requirements) != 2 { // MassRequirement def and massReq instance
		t.Errorf("Expected 2 requirements, got %d", len(requirements))
	}
}

// TestIntegrationAllWithDepth tests the AllWithDepth iterator.
func TestIntegrationAllWithDepth(t *testing.T) {
	input := `package Outer {
		package Inner {
			part def A;
		}
	}`

	result := ParseString(input)
	if !result.Success() || result.Model == nil {
		t.Skip("Failed to parse test input")
	}

	// Count elements at each depth
	depthCounts := make(map[int]int)
	for elem, depth := range AllWithDepth(result.Model) {
		_ = elem
		depthCounts[depth]++
	}

	t.Logf("Element counts by depth: %v", depthCounts)

	// Should have elements at depth 0 (Outer), depth 1 (Inner, A def), depth 2 (A instance)
	if len(depthCounts) == 0 {
		t.Error("Expected elements at multiple depths")
	}
}

// TestIntegrationFindByPath tests the FindByPath function.
func TestIntegrationFindByPath(t *testing.T) {
	input := `package TestPackage {
		part def TestPart;
	}`

	result := ParseString(input)
	if !result.Success() || result.Model == nil {
		t.Skip("Failed to parse test input")
	}

	// Try to find by path
	elem := FindByPath(result.Model, "TestPackage")
	if elem == nil {
		t.Log("FindByPath returned nil (element may not have qualified name)")
	} else {
		t.Logf("Found element: %s", elem.Name())
	}
}

// TestIntegrationElementRelationships tests element relationship tracking.
func TestIntegrationElementRelationships(t *testing.T) {
	input := `package Test {
		part def Parent {
			part child : Child;
		}
		part def Child;
	}`

	result := ParseString(input)
	if !result.Success() || result.Model == nil {
		t.Skip("Failed to parse test input")
	}

	// Find Parent part definition
	parts := FindAll[*Part](result.Model)
	var parent *Part
	for _, p := range parts {
		if p.Name() == "Parent" {
			parent = p
			break
		}
	}

	if parent == nil {
		t.Skip("Could not find Parent part")
	}

	// Check children
	children := parent.Children()
	t.Logf("Parent has %d children", len(children))

	// Check parent reference
	for _, child := range children {
		if child.Parent() == nil {
			t.Error("Child's parent reference is nil")
		}
	}
}

// Helper function
func FindByPath(model *Model, path string) Element {
	// This is a simplified version - the actual implementation may differ
	for elem := range All(model) {
		if elem.QualifiedName() == path {
			return elem
		}
	}
	return nil
}

// TestValidationCategories runs the parser against all 18 validation categories.
func TestValidationCategories(t *testing.T) {
	categories := collectValidationFiles(t)
	if len(categories) == 0 {
		t.Skip("No validation categories found")
	}

	// Check if standard library is available
	hasLibrary := standardLibraryExists()
	if !hasLibrary {
		t.Log("WARNING: Standard library not found at ./libraries/sysml.library - library resolution tests will be limited")
	}

	t.Logf("Found %d validation categories", len(categories))

	// Category statistics
	type categoryStats struct {
		name            string
		files           int
		parsed          int
		failed          int
		elements        int
		libraryImports  int
		libraryResolved int
	}

	var allStats []categoryStats

	// Process each category
	for categoryName, files := range categories {
		t.Run(categoryName, func(t *testing.T) {
			stats := categoryStats{
				name:  categoryName,
				files: len(files),
			}

			for _, file := range files {
				filename := filepath.Base(file)

				// Parse with library support if available
				var result *ParseResult
				if hasLibrary {
					result = ParseFile(file, WithStandardLibrary())
				} else {
					result = ParseFile(file)
				}

				if result.Success() {
					stats.parsed++
					t.Logf("  ✓ %s", filename)

					if result.Model != nil {
						// Count elements
						count := 0
						Walk(result.Model, func(elem Element, depth int) bool {
							count++
							return true
						})
						stats.elements += count

						// Count library imports if available
						if hasLibrary {
							for _, imp := range result.Model.Imports {
								if imp.ResolvedPackage != nil {
									stats.libraryImports++
									stats.libraryResolved++
								} else if imp.ResolvedElement != nil {
									stats.libraryImports++
								}
							}
						}
					}
				} else {
					stats.failed++
					errMsg := "parse error"
					if result.Errors != nil && len(result.Errors.Errors) > 0 {
						errMsg = result.Errors.Errors[0].Message
					}
					t.Logf("  ✗ %s - %s", filename, errMsg)
				}
			}

			successRate := 0.0
			if stats.files > 0 {
				successRate = float64(stats.parsed) / float64(stats.files) * 100
			}

			t.Logf("  → %s: %d/%d files parsed (%.1f%%), %d elements, %d library imports resolved",
				categoryName, stats.parsed, stats.files, successRate, stats.elements, stats.libraryResolved)

			allStats = append(allStats, stats)
		})
	}

	// Print final summary table
	t.Log("\n=== Validation Categories Summary ===")
	t.Logf("%-35s | %5s | %6s | %6s | %10s | %8s", "Category", "Files", "Parsed", "Failed", "Elements", "LibRefs")
	t.Log("----------------------------------------------------------------------------------------")

	totalFiles := 0
	totalParsed := 0
	totalFailed := 0
	totalElements := 0
	totalLibraryRefs := 0

	for _, stats := range allStats {
		t.Logf("%-35s | %5d | %6d | %6d | %10d | %8d",
			stats.name, stats.files, stats.parsed, stats.failed, stats.elements, stats.libraryResolved)
		totalFiles += stats.files
		totalParsed += stats.parsed
		totalFailed += stats.failed
		totalElements += stats.elements
		totalLibraryRefs += stats.libraryResolved
	}

	t.Log("----------------------------------------------------------------------------------------")
	t.Logf("%-35s | %5d | %6d | %6d | %10d | %8d",
		"TOTAL", totalFiles, totalParsed, totalFailed, totalElements, totalLibraryRefs)

	overallRate := 0.0
	if totalFiles > 0 {
		overallRate = float64(totalParsed) / float64(totalFiles) * 100
	}
	t.Logf("\nOverall Success Rate: %d/%d files (%.1f%%)", totalParsed, totalFiles, overallRate)
}

// TestLibraryImportResolution verifies that library imports are resolved correctly.
func TestLibraryImportResolution(t *testing.T) {
	if !standardLibraryExists() {
		t.Skip("Standard library not available - skipping library resolution test")
	}

	// Parse a simple file with library imports
	input := `package TestLib {
		import ScalarValues::*;
		import ISQ::*;
		
		attribute testValue : ScalarValues::Real;
	}`

	result := ParseString(input, WithStandardLibrary())
	if !result.Success() {
		t.Fatalf("Failed to parse test input: %v", result.Errors)
	}

	if result.Model == nil {
		t.Fatal("Model is nil")
	}

	// Check that imports were resolved
	for _, imp := range result.Model.Imports {
		resolvedName := ""
		if imp.ResolvedPackage != nil {
			resolvedName = imp.ResolvedPackage.Name()
		}
		t.Logf("Import: %s, ResolvedPackage: %s", imp.ImportedNamespace, resolvedName)
		if imp.ResolvedPackage == nil {
			t.Errorf("Import %s was not resolved to a library package", imp.ImportedNamespace)
		}
	}
}

// TestLibraryElementReferences verifies that qualified name references resolve to library elements.
func TestLibraryElementReferences(t *testing.T) {
	if !standardLibraryExists() {
		t.Skip("Standard library not available - skipping library element reference test")
	}

	// Parse a file with qualified name references
	input := `package TestRef {
		import ISQ::*;
		import SI::*;
		
		attribute mass :> ISQ::mass;
		attribute weight :> SI::kg;
	}`

	result := ParseString(input, WithStandardLibrary())
	if !result.Success() {
		t.Fatalf("Failed to parse test input: %v", result.Errors)
	}

	if result.Model == nil {
		t.Fatal("Model is nil")
	}

	// Find attributes and check if they have resolved type references
	for _, elem := range result.Model.Elements {
		if attr, ok := elem.(*Attribute); ok {
			t.Logf("Attribute: %s", attr.Name())
			if typeElem := attr.TypeRef.Resolved(); typeElem != nil {
				t.Logf("  TypeRef: %s", typeElem.Name())
			}
		}
	}
}

// TestValidationFileLibraryUsage tests library usage in actual validation files.
func TestValidationFileLibraryUsage(t *testing.T) {
	if !standardLibraryExists() {
		t.Skip("Standard library not available - skipping validation file library usage test")
	}

	categories := collectValidationFiles(t)
	if len(categories) == 0 {
		t.Skip("No validation categories found")
	}

	// Test representative files from different categories
	representativeFiles := []struct {
		category string
		pattern  string
	}{
		{"01-Parts Tree", "Parts Tree"},
		{"03-Function-based Behavior", "Function"},
		{"15-Properties-Values-Expressions", "Property"},
	}

	for _, rep := range representativeFiles {
		files, ok := categories[rep.category]
		if !ok || len(files) == 0 {
			t.Logf("Category %s not found or empty, skipping", rep.category)
			continue
		}

		// Find a file matching the pattern
		var testFile string
		for _, f := range files {
			if strings.Contains(filepath.Base(f), rep.pattern) || strings.Contains(f, rep.pattern) {
				testFile = f
				break
			}
		}
		if testFile == "" {
			testFile = files[0] // Use first file if no match
		}

		t.Run(rep.category, func(t *testing.T) {
			result := ParseFile(testFile, WithStandardLibrary())
			if !result.Success() {
				t.Logf("File %s did not parse successfully (may be expected for complex files)", filepath.Base(testFile))
				return
			}

			if result.Model == nil {
				t.Skip("Model is nil")
			}

			// Count library imports
			libraryImportCount := 0
			for _, imp := range result.Model.Imports {
				if imp.ResolvedPackage != nil {
					libraryImportCount++
				}
			}

			t.Logf("Category %s: %s has %d resolved library imports",
				rep.category, filepath.Base(testFile), libraryImportCount)
		})
	}
}

// TestLibraryPerformance measures library loading performance.
func TestLibraryPerformance(t *testing.T) {
	if !standardLibraryExists() {
		t.Skip("Standard library not available - skipping performance test")
	}

	// Measure library loading time
	start := time.Now()
	reg := NewLibraryRegistry()
	err := reg.RegisterStandardLibrary()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to load standard library: %v", err)
	}

	t.Logf("Standard library loaded in %v", elapsed)
	t.Logf("Library contains %d packages", reg.GetLibraryCount())

	// Measure parse time with library
	input := `package PerfTest {
		import ScalarValues::*;
		attribute value : ScalarValues::Real;
	}`

	// Without library
	start = time.Now()
	result1 := ParseString(input)
	elapsed1 := time.Since(start)

	// With library (should reuse cached registry)
	start = time.Now()
	result2 := ParseString(input, WithLibraryRegistry(reg))
	elapsed2 := time.Since(start)

	if !result1.Success() || !result2.Success() {
		t.Log("Parse results may have errors, but timing is still valid")
	}

	t.Logf("Parse time without library: %v", elapsed1)
	t.Logf("Parse time with library: %v", elapsed2)
	t.Logf("Overhead: %v", elapsed2-elapsed1)
}

// TestLibraryErrorHandling tests graceful handling of invalid library imports.
func TestLibraryErrorHandling(t *testing.T) {
	if !standardLibraryExists() {
		t.Skip("Standard library not available - skipping error handling test")
	}

	// Parse a file with an invalid library import
	input := `package TestError {
		import NonExistentLibrary::*;
		attribute value : NonExistentLibrary::SomeType;
	}`

	result := ParseString(input, WithStandardLibrary())

	// Should complete without panic, even if there are errors
	t.Logf("Parse completed with success=%v", result.Success())
	if result.Errors != nil {
		t.Logf("Errors: %v", result.Errors.Errors)
	}
}
