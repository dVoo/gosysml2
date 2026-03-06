package sysml

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestLibraryDiscovery tests the library file discovery functionality.
func TestLibraryDiscovery(t *testing.T) {
	// Create temporary directory with mock library files
	tempDir := t.TempDir()

	// Create mock library files
	mockFiles := []string{
		"Library1.sysml",
		"Library2.sysml",
		"Library3.kerml",
		"subdir/Nested1.sysml",
		"subdir/Nested2.kerml",
		"notalibrary.txt",
		"README.md",
	}

	for _, file := range mockFiles {
		path := filepath.Join(tempDir, file)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte("mock content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Create registry and test discovery
	registry := NewLibraryRegistry()
	files, err := registry.DiscoverLibraries(tempDir)
	if err != nil {
		t.Fatalf("DiscoverLibraries failed: %v", err)
	}

	// Should find only .sysml and .kerml files (5 total)
	if len(files) != 5 {
		t.Errorf("Expected 5 library files, got %d", len(files))
	}

	// Verify all expected files are found
	found := make(map[string]bool)
	for _, f := range files {
		found[filepath.Base(f)] = true
	}

	expected := []string{"Library1.sysml", "Library2.sysml", "Library3.kerml", "Nested1.sysml", "Nested2.kerml"}
	for _, exp := range expected {
		if !found[exp] {
			t.Errorf("Expected to find %s", exp)
		}
	}

	// Verify non-library files are NOT found
	if found["notalibrary.txt"] {
		t.Error("Should not find .txt files")
	}
	if found["README.md"] {
		t.Error("Should not find .md files")
	}
}

// TestLibraryDiscoveryNonExistentPath tests discovery with non-existent path.
func TestLibraryDiscoveryNonExistentPath(t *testing.T) {
	registry := NewLibraryRegistry()
	files, err := registry.DiscoverLibraries("/non/existent/path")
	if err != nil {
		t.Fatalf("Should not error for non-existent path: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected 0 files for non-existent path, got %d", len(files))
	}
}

// TestLibraryRegistry tests the basic registry functionality.
func TestLibraryRegistry(t *testing.T) {
	registry := NewLibraryRegistry()

	// Test initial state
	if registry.GetLibraryCount() != 0 {
		t.Errorf("Expected 0 libraries initially, got %d", registry.GetLibraryCount())
	}
	if registry.GetElementCount() != 0 {
		t.Errorf("Expected 0 elements initially, got %d", registry.GetElementCount())
	}
	if registry.IsLoaded() {
		t.Error("Registry should not be loaded initially")
	}

	// Test FindElement returns nil for unknown elements
	elem := registry.FindElement("NonExistent")
	if elem != nil {
		t.Error("FindElement should return nil for unknown elements")
	}

	// Test ResolveImport returns error for unknown packages
	_, err := registry.ResolveImport("UnknownPackage")
	if err == nil {
		t.Error("ResolveImport should return error for unknown packages")
	}
}

// TestLibraryRegistryConcurrentAccess tests thread-safety.
func TestLibraryRegistryConcurrentAccess(t *testing.T) {
	registry := NewLibraryRegistry()

	// Run concurrent operations
	var wg sync.WaitGroup
	numGoroutines := 10

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = registry.FindElement("Test")
			_ = registry.GetLibraryCount()
			_ = registry.GetElementCount()
			_ = registry.IsLoaded()
		}()
	}

	wg.Wait()
}

// TestLibraryResolver tests import resolution functionality.
func TestLibraryResolver(t *testing.T) {
	// Create temporary directory with mock library
	tempDir := t.TempDir()

	// Create a mock library file with nested packages
	libContent := `standard library package TestLib {
		package NestedPkg {
			attribute def TestAttribute;
		}
		attribute def TopLevelAttribute;
	}
	`

	libPath := filepath.Join(tempDir, "TestLib.sysml")
	if err := os.WriteFile(libPath, []byte(libContent), 0644); err != nil {
		t.Fatalf("Failed to create mock library: %v", err)
	}

	// Create registry with custom path
	registry := NewLibraryRegistry(WithLibraryPaths(tempDir))

	// Register the library
	err := registry.RegisterLibrary(libPath)
	if err != nil {
		t.Fatalf("Failed to register library: %v", err)
	}

	// Test ResolveImport with simple package name
	pkg, err := registry.ResolveImport("TestLib")
	if err != nil {
		t.Fatalf("Failed to resolve TestLib: %v", err)
	}
	if pkg.Name() != "TestLib" {
		t.Errorf("Expected package name 'TestLib', got '%s'", pkg.Name())
	}
	if !pkg.IsLibrary {
		t.Error("Package should be marked as library")
	}

	// Test ResolveImport with nested package
	nestedPkg, err := registry.ResolveImport("TestLib::NestedPkg")
	if err != nil {
		t.Fatalf("Failed to resolve TestLib::NestedPkg: %v", err)
	}
	if nestedPkg.Name() != "NestedPkg" {
		t.Errorf("Expected package name 'NestedPkg', got '%s'", nestedPkg.Name())
	}

	// Test FindElement
	elem := registry.FindElement("TestLib::TopLevelAttribute")
	if elem == nil {
		t.Error("Should find TopLevelAttribute")
	}

	// Test FindElement for nested element
	nestedElem := registry.FindElement("TestLib::NestedPkg::TestAttribute")
	if nestedElem == nil {
		t.Error("Should find TestAttribute in nested package")
	}
}

// TestLibraryResolverInvalidPatterns tests resolution with invalid patterns.
func TestLibraryResolverInvalidPatterns(t *testing.T) {
	registry := NewLibraryRegistry()

	// Test empty namespace
	_, err := registry.ResolveImport("")
	if err == nil {
		t.Error("Should error on empty namespace")
	}

	// Test unknown package
	_, err = registry.ResolveImport("UnknownPackage")
	if err == nil {
		t.Error("Should error on unknown package")
	}

	// Test unknown nested package
	// First create a simple library
	tempDir := t.TempDir()
	libContent := `standard library package TestPkg {}
	`
	libPath := filepath.Join(tempDir, "TestPkg.sysml")
	os.WriteFile(libPath, []byte(libContent), 0644)

	registry.RegisterLibrary(libPath)

	_, err = registry.ResolveImport("TestPkg::NonExistent")
	if err == nil {
		t.Error("Should error on unknown nested package")
	}
}

// TestStandardLibraryLoading tests loading the actual standard library.
func TestStandardLibraryLoading(t *testing.T) {
	registry := NewLibraryRegistry(WithLibraryPaths(standardLibraryDir(t)))

	// Load standard library
	err := registry.RegisterStandardLibrary()
	if err != nil {
		t.Logf("Standard library loaded with warnings: %v", err)
	}

	// Log statistics
	libCount := registry.GetLibraryCount()
	elemCount := registry.GetElementCount()
	t.Logf("Loaded %d libraries with %d elements", libCount, elemCount)

	// Verify we loaded some libraries
	if libCount == 0 {
		t.Error("Should have loaded at least some libraries")
	}

	// Test that registry is marked as loaded
	if !registry.IsLoaded() {
		t.Error("Registry should be marked as loaded")
	}

	// Test GetLibraryNames returns sorted list
	names := registry.GetLibraryNames()
	if len(names) != libCount {
		t.Errorf("GetLibraryNames returned %d names, expected %d", len(names), libCount)
	}

	// Verify names are sorted
	for i := 1; i < len(names); i++ {
		if strings.Compare(names[i-1], names[i]) > 0 {
			t.Error("Library names should be sorted")
			break
		}
	}

	// Test finding common standard library elements
	// Note: These may or may not exist depending on the library content
	testElements := []string{
		"ScalarValues::Real",
		"ScalarValues::Integer",
		"ScalarValues::Boolean",
	}

	foundCount := 0
	for _, elemName := range testElements {
		if elem := registry.FindElement(elemName); elem != nil {
			t.Logf("Found element: %s", elemName)
			foundCount++
		}
	}
	t.Logf("Found %d/%d test elements", foundCount, len(testElements))
}

// TestLibraryWithOptions tests creating registry with custom options.
func TestLibraryWithOptions(t *testing.T) {
	customPath := "/custom/library/path"

	registry := NewLibraryRegistry(WithLibraryPaths(customPath))

	// The registry should use the custom path
	// We can verify this by checking that discovery fails for the custom path
	// (assuming it doesn't exist)
	files, err := registry.DiscoverLibraries(customPath)
	if err != nil {
		t.Fatalf("Should not error for non-existent path in test: %v", err)
	}
	if len(files) != 0 {
		t.Error("Should return empty slice for non-existent custom path")
	}
}

// TestLibraryLoadError tests error handling for invalid library files.
func TestLibraryLoadError(t *testing.T) {
	// Create temporary directory with invalid library file
	tempDir := t.TempDir()

	// Create an invalid SysML file
	invalidContent := `this is not valid sysml { broken syntax`
	invalidPath := filepath.Join(tempDir, "Invalid.sysml")
	os.WriteFile(invalidPath, []byte(invalidContent), 0644)

	registry := NewLibraryRegistry()

	// Try to load the invalid library
	_, err := registry.LoadLibrary(invalidPath)
	if err == nil {
		t.Error("Should error when loading invalid library file")
	}

	// Try to register the invalid library
	err = registry.RegisterLibrary(invalidPath)
	if err == nil {
		t.Error("Should error when registering invalid library file")
	}
}

// TestLibraryMultiplePaths tests registry with multiple library paths.
func TestLibraryMultiplePaths(t *testing.T) {
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()

	// Create libraries in both directories
	content1 := `standard library package Lib1 {}`
	content2 := `standard library package Lib2 {}`

	os.WriteFile(filepath.Join(tempDir1, "Lib1.sysml"), []byte(content1), 0644)
	os.WriteFile(filepath.Join(tempDir2, "Lib2.sysml"), []byte(content2), 0644)

	// Create registry with multiple paths
	registry := NewLibraryRegistry(WithLibraryPaths(tempDir1, tempDir2))

	// Load standard library (should search both paths)
	err := registry.RegisterStandardLibrary()
	if err != nil {
		t.Logf("Registration completed with warnings: %v", err)
	}

	// Should have loaded libraries from both paths
	if registry.GetLibraryCount() < 2 {
		t.Errorf("Expected at least 2 libraries from multiple paths, got %d",
			registry.GetLibraryCount())
	}

	// Verify both packages can be resolved
	_, err = registry.ResolveImport("Lib1")
	if err != nil {
		t.Errorf("Should resolve Lib1: %v", err)
	}

	_, err = registry.ResolveImport("Lib2")
	if err != nil {
		t.Errorf("Should resolve Lib2: %v", err)
	}
}

// BenchmarkLibraryDiscovery benchmarks the library discovery process.
func BenchmarkLibraryDiscovery(b *testing.B) {
	registry := NewLibraryRegistry()
	libDir := standardLibraryDir(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.DiscoverLibraries(libDir)
		if err != nil {
			b.Fatalf("Discovery failed: %v", err)
		}
	}
}

// BenchmarkLibraryRegistryFindElement benchmarks element lookup.
func BenchmarkLibraryRegistryFindElement(b *testing.B) {
	// Create temporary library
	tempDir := b.TempDir()
	content := `standard library package BenchLib {
		attribute def Attr1;
		attribute def Attr2;
		attribute def Attr3;
	}`
	os.WriteFile(filepath.Join(tempDir, "BenchLib.sysml"), []byte(content), 0644)

	registry := NewLibraryRegistry(WithLibraryPaths(tempDir))
	registry.RegisterStandardLibrary()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.FindElement("BenchLib::Attr1")
	}
}
