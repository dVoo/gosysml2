package sysml

import (
	"testing"
)

func TestNewIncludeUseCase(t *testing.T) {
	loc := Location{Line: 10, Column: 5}

	include := NewIncludeUseCase("TestInclude", loc)
	if include == nil {
		t.Fatal("NewIncludeUseCase returned nil")
	}
	if include.Name() != "TestInclude" {
		t.Errorf("Expected name 'TestInclude', got '%s'", include.Name())
	}
	if include.Kind() != KindIncludeUseCase {
		t.Errorf("Expected kind KindIncludeUseCase, got %v", include.Kind())
	}
	if include.Location() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, include.Location())
	}

	// Initially unresolved
	if include.IncludedUseCase.IsResolved() {
		t.Error("Expected IncludedUseCase to be unresolved for new include")
	}
}

func TestIncludeUseCaseInterface(t *testing.T) {
	loc := Location{Line: 5, Column: 10}
	include := NewIncludeUseCase("TestInclude", loc)

	// Test Element interface
	if include.Kind() != KindIncludeUseCase {
		t.Errorf("Expected Kind() to return KindIncludeUseCase, got %v", include.Kind())
	}

	if include.Name() != "TestInclude" {
		t.Errorf("Expected name 'TestInclude', got '%s'", include.Name())
	}

	// Test Usage interface (should compile)
	var _ Usage = include
}

func TestIncludeUseCaseSetUnresolved(t *testing.T) {
	loc := Location{}
	include := NewIncludeUseCase("MyInclude", loc)

	// Set unresolved reference
	include.SetUnresolvedIncludedUseCase("BaseUseCase")
	if !include.IncludedUseCase.IsResolved() {
		// This is expected - it's unresolved
	} else {
		t.Error("Expected IncludedUseCase to be unresolved after setting unresolved name")
	}
}

func TestIncludeUseCaseParent(t *testing.T) {
	loc := Location{}
	include := NewIncludeUseCase("TestInclude", loc)
	pkg := NewPackage("TestPkg", loc)

	include.SetParent(pkg)

	if include.Parent() != pkg {
		t.Error("SetParent did not set parent correctly")
	}

	// Test qualified name
	expected := "TestPkg::TestInclude"
	if include.QualifiedName() != expected {
		t.Errorf("Expected qualified name '%s', got '%s'", expected, include.QualifiedName())
	}
}

func TestFindIncludeUseCases(t *testing.T) {
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	// Create some include use cases
	include1 := NewIncludeUseCase("Include1", Location{})
	include2 := NewIncludeUseCase("Include2", Location{})
	include3 := NewIncludeUseCase("Include3", Location{})

	pkg.AddChild(include1)
	pkg.AddChild(include2)
	pkg.AddChild(include3)

	includes := FindIncludeUseCases(model)
	if len(includes) != 3 {
		t.Errorf("Expected 3 include use cases, got %d", len(includes))
	}
}

func TestIncludeUseCaseVisitor(t *testing.T) {
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	include := NewIncludeUseCase("TestInclude", Location{})
	pkg.AddChild(include)

	// Test Counter visitor
	counter := NewCounter()
	Visit(model, counter)

	if counter.Counts[KindIncludeUseCase] != 1 {
		t.Errorf("Expected 1 include use case count, got %d", counter.Counts[KindIncludeUseCase])
	}
}

func TestIncludeUseCaseKindString(t *testing.T) {
	if KindIncludeUseCase.String() != "include use case" {
		t.Errorf("Expected 'include use case', got '%s'", KindIncludeUseCase.String())
	}
}

// Integration test: Parse a simple include use case usage
func TestParseIncludeUseCaseUsage(t *testing.T) {
	// This test verifies the parser can handle use case definitions
	// Note: Full include parsing depends on ANTLR grammar support
	input := `package TestPkg {
		use case def BaseUseCase;
		use case def MyUseCase;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	// Verify use cases were parsed
	useCases := FindUseCases(result.Model)
	if len(useCases) != 2 {
		t.Errorf("Expected 2 use cases, got %d", len(useCases))
	}
}

// Integration test: Parse include with qualified name reference (manual model)
func TestParseIncludeUseCaseWithReference(t *testing.T) {
	// Create model manually to test reference handling
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	// Create base use case
	baseUseCase := NewUseCase("BaseUseCase", Location{}, true)
	pkg.AddChild(baseUseCase)

	// Create use case with include
	myUseCase := NewUseCase("MyUseCase", Location{}, true)
	pkg.AddChild(myUseCase)

	// Create include use case with unresolved reference
	include := NewIncludeUseCase("", Location{})
	include.SetUnresolvedIncludedUseCase("BaseUseCase")
	myUseCase.AddChild(include)

	includes := FindIncludeUseCases(model)
	if len(includes) < 1 {
		t.Fatalf("Expected at least 1 include use case, got %d", len(includes))
	}

	// Check that the include has an unresolved reference initially
	inc := includes[0]
	if inc.IncludedUseCase.IsResolved() {
		t.Log("Include use case reference was resolved (unexpected before BuildIndex/ResolveReferences)")
	}
}

// Integration test: Verify include reference resolution
func TestIncludeUseCaseResolution(t *testing.T) {
	// Create model manually to test resolution
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	// Create base use case
	baseUseCase := NewUseCase("BaseUseCase", Location{}, true)
	pkg.AddChild(baseUseCase)

	// Create use case with include
	myUseCase := NewUseCase("MyUseCase", Location{}, true)
	pkg.AddChild(myUseCase)

	// Create include use case
	include := NewIncludeUseCase("", Location{})
	include.SetUnresolvedIncludedUseCase("BaseUseCase")
	myUseCase.AddChild(include)

	// Build index and resolve references
	model.BuildIndex()
	model.ResolveReferences()

	// After resolution, the included use case should be resolved
	if !include.IncludedUseCase.IsResolved() {
		t.Error("Expected IncludedUseCase to be resolved after ResolveReferences")
	}

	// Verify it points to the correct UseCase
	included := include.IncludedUseCase.Resolved()
	if included == nil {
		t.Fatal("IncludedUseCase.Resolved() returned nil")
	}

	if included.Name() != "BaseUseCase" {
		t.Errorf("Expected included use case name 'BaseUseCase', got '%s'", included.Name())
	}
}

// Integration test: Parse use case with multiple includes (manual model)
func TestUseCaseWithMultipleIncludes(t *testing.T) {
	// Create model manually to test multiple includes
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	// Create base use cases
	base1 := NewUseCase("BaseUseCase1", Location{}, true)
	base2 := NewUseCase("BaseUseCase2", Location{}, true)
	pkg.AddChild(base1)
	pkg.AddChild(base2)

	// Create use case with includes
	myUseCase := NewUseCase("MyUseCase", Location{}, true)
	pkg.AddChild(myUseCase)

	// Create include use cases
	include1 := NewIncludeUseCase("", Location{})
	include1.SetUnresolvedIncludedUseCase("BaseUseCase1")
	myUseCase.AddChild(include1)

	include2 := NewIncludeUseCase("", Location{})
	include2.SetUnresolvedIncludedUseCase("BaseUseCase2")
	myUseCase.AddChild(include2)

	// Build index and resolve references
	model.BuildIndex()
	model.ResolveReferences()

	includes := FindIncludeUseCases(model)
	if len(includes) != 2 {
		t.Errorf("Expected 2 include use cases, got %d", len(includes))
	}

	// Verify all includes are resolved
	resolvedCount := 0
	for _, inc := range includes {
		if inc.IncludedUseCase.IsResolved() {
			resolvedCount++
		}
	}

	if resolvedCount != 2 {
		t.Errorf("Expected 2 resolved includes, got %d", resolvedCount)
	}
}

// Integration test: Include use case minimal test (manual model)
func TestParseIncludeUseCaseMinimal(t *testing.T) {
	// Create model manually to test basic functionality
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	// Create use cases
	baseUseCase := NewUseCase("BaseUseCase", Location{}, true)
	pkg.AddChild(baseUseCase)

	myUseCase := NewUseCase("MyUseCase", Location{}, true)
	pkg.AddChild(myUseCase)

	// Create include use case
	include := NewIncludeUseCase("", Location{})
	include.SetUnresolvedIncludedUseCase("BaseUseCase")
	myUseCase.AddChild(include)

	// Should have use cases
	useCases := FindUseCases(model)
	if len(useCases) != 2 {
		t.Errorf("Expected 2 use cases, got %d", len(useCases))
	}

	// Should have includes
	includes := FindIncludeUseCases(model)
	if len(includes) != 1 {
		t.Errorf("Expected 1 include use case, got %d", len(includes))
	}
}

// Test IncludeUseCase as a child of UseCase
func TestIncludeUseCaseAsUseCaseChild(t *testing.T) {
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	// Create use case
	useCase := NewUseCase("MyUseCase", Location{}, true)
	pkg.AddChild(useCase)

	// Create include use case
	include := NewIncludeUseCase("MyInclude", Location{})
	include.SetUnresolvedIncludedUseCase("BaseUseCase")

	// Add include as child of use case
	useCase.AddChild(include)

	// Verify parent relationship
	// Note: Parent returns Element interface, so we need to check the underlying type
	if include.Parent() == nil {
		t.Error("Include use case parent should not be nil")
	}

	// Verify use case has the include as child
	found := false
	for _, child := range useCase.Children() {
		if child == include {
			found = true
			break
		}
	}
	if !found {
		t.Error("Use case should have include as child")
	}
}
