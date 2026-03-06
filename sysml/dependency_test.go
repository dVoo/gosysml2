package sysml

import (
	"testing"
)

func TestNewDependency(t *testing.T) {
	loc := Location{Line: 10, Column: 5}
	dep := NewDependency(loc)

	if dep == nil {
		t.Fatal("NewDependency returned nil")
	}

	if dep.Kind() != KindDependency {
		t.Errorf("Expected kind %v, got %v", KindDependency, dep.Kind())
	}

	if dep.GetKind() != "dependency" {
		t.Errorf("Expected GetKind() = 'dependency', got '%s'", dep.GetKind())
	}

	if dep.Name() != "" {
		t.Errorf("Expected empty name, got '%s'", dep.Name())
	}

	if dep.GetLocation() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, dep.GetLocation())
	}

	if dep.GetParent() != nil {
		t.Error("Expected nil parent for new dependency")
	}
}

func TestDependencyParent(t *testing.T) {
	dep := NewDependency(Location{})
	pkg := NewPackage("TestPkg", Location{})

	dep.SetParent(pkg)

	if dep.GetParent() != pkg {
		t.Error("SetParent did not set parent correctly")
	}
}

func TestDependencyUnresolvedReferences(t *testing.T) {
	dep := NewDependency(Location{})

	// Add unresolved references
	dep.AddUnresolvedClient("ClientA")
	dep.AddUnresolvedClient("ClientB")
	dep.AddUnresolvedSupplier("SupplierX")

	clientRefs, supplierRefs := dep.UnresolvedReferences()

	if len(clientRefs) != 2 {
		t.Errorf("Expected 2 client refs, got %d", len(clientRefs))
	}

	if len(supplierRefs) != 1 {
		t.Errorf("Expected 1 supplier ref, got %d", len(supplierRefs))
	}

	if clientRefs[0] != "ClientA" || clientRefs[1] != "ClientB" {
		t.Errorf("Unexpected client refs: %v", clientRefs)
	}

	if supplierRefs[0] != "SupplierX" {
		t.Errorf("Unexpected supplier ref: %v", supplierRefs)
	}
}

func TestDependencyVisitor(t *testing.T) {
	dep := NewDependency(Location{Line: 1, Column: 1})

	// Test with a counting visitor
	var visited bool
	visitor := &testDependencyVisitor{
		onVisit: func(d *Dependency) {
			visited = true
			if d != dep {
				t.Error("Visitor received wrong dependency")
			}
		},
	}

	dep.Accept(visitor)

	if !visited {
		t.Error("Visitor was not called")
	}
}

type testDependencyVisitor struct {
	BaseVisitor
	onVisit func(*Dependency)
}

func (v *testDependencyVisitor) VisitDependency(d *Dependency) bool {
	if v.onVisit != nil {
		v.onVisit(d)
	}
	return true
}

func TestDependencyTypeMethods(t *testing.T) {
	dep := NewDependency(Location{})

	// Test Type() method returns nil
	if dep.Type() != nil {
		t.Error("Expected Type() to return nil for Dependency")
	}
}

func TestModelAddDependency(t *testing.T) {
	model := NewModel()
	dep := NewDependency(Location{Line: 5, Column: 10})

	model.AddDependency(dep)

	if len(model.Dependencies()) != 1 {
		t.Errorf("Expected 1 dependency in model, got %d", len(model.Dependencies()))
	}

	if len(model.Elements) != 1 {
		t.Errorf("Expected 1 element in model, got %d", len(model.Elements))
	}

	if model.Dependencies()[0] != dep {
		t.Error("Dependency not added correctly to model")
	}
}

func TestDependencyIsDefinitionAndUsage(t *testing.T) {
	dep := NewDependency(Location{})

	// These methods exist to satisfy interfaces but don't change behavior
	dep.isDefinition()
	dep.isUsage()
}
