package sysml

import (
	"testing"
)

func TestLibraryPackageSyntax(t *testing.T) {
	// Test parsing standard library package syntax
	content := `standard library package TestLib {
		attribute def TestAttr;
	}`
	
	result := ParseString(content)
	if result.Errors != nil && result.Errors.HasErrors() {
		t.Fatalf("Parse errors: %v", result.Errors)
	}
	
	if result.Model == nil {
		t.Fatal("Model is nil")
	}
	
	t.Logf("Packages found: %d", len(result.Model.Packages))
	for _, pkg := range result.Model.Packages {
		t.Logf("Package: %s (IsLibrary: %v)", pkg.Name(), pkg.IsLibrary)
	}
	
	if len(result.Model.Packages) == 0 {
		t.Error("Expected at least one package, got 0")
	}
}
