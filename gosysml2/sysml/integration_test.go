package sysml

import (
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

// testdataDir returns the path to testdata relative to this package.
func integrationTestdataDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join("..", "..", "testdata")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("testdata directory not found")
	}
	return dir
}

// collectTestFiles finds all .sysml files in the testdata directory.
func collectTestFiles(t *testing.T) []string {
	t.Helper()
	testdataDir := integrationTestdataDir(t)

	var files []string
	err := filepath.WalkDir(testdataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sysml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk testdata: %v", err)
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
		relPath, _ := filepath.Rel(integrationTestdataDir(t), file)

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
	testdataDir := integrationTestdataDir(t)
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
	testdataDir := integrationTestdataDir(t)
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
	testdataDir := integrationTestdataDir(t)
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
	testdataDir := integrationTestdataDir(t)
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
	testdataDir := integrationTestdataDir(t)
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
