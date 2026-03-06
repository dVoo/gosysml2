package sysml

import (
	"testing"
)

// TestPartSpecializationWithColonGT verifies that part definitions using :> syntax
// have their specialization relationships properly captured and resolved.
func TestPartSpecializationWithColonGT(t *testing.T) {
	input := `
package TestPackage {
    part def Vehicle;
    part def Car :> Vehicle;
    part def ElectricCar :> Car;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Err())
	}

	model := result.Model

	// Find all parts
	parts := FindAll[*Part](model)
	if len(parts) != 3 {
		t.Fatalf("Expected 3 parts, got %d", len(parts))
	}

	// Find Vehicle (base definition)
	var vehicle, car, electricCar *Part
	for _, p := range parts {
		switch p.Name() {
		case "Vehicle":
			vehicle = p
		case "Car":
			car = p
		case "ElectricCar":
			electricCar = p
		}
	}

	if vehicle == nil {
		t.Fatal("Vehicle part not found")
	}
	if car == nil {
		t.Fatal("Car part not found")
	}
	if electricCar == nil {
		t.Fatal("ElectricCar part not found")
	}

	// Verify all are definitions
	if !vehicle.IsDefinition {
		t.Error("Vehicle should be a definition")
	}
	if !car.IsDefinition {
		t.Error("Car should be a definition")
	}
	if !electricCar.IsDefinition {
		t.Error("ElectricCar should be a definition")
	}

	// Verify Car specializes Vehicle
	if !car.Specializes.IsResolved() {
		t.Error("Car.Specializes should be resolved to Vehicle")
	} else {
		if car.Specializes.Resolved() != vehicle {
			t.Error("Car should specialize Vehicle")
		}
	}

	// Verify ElectricCar specializes Car
	if !electricCar.Specializes.IsResolved() {
		t.Error("ElectricCar.Specializes should be resolved to Car")
	} else {
		if electricCar.Specializes.Resolved() != car {
			t.Error("ElectricCar should specialize Car")
		}
	}

	// Verify Vehicle has no specialization
	if vehicle.Specializes.IsResolved() {
		t.Error("Vehicle should not have a specialization")
	}
}

// TestItemSpecializationWithColonGT verifies that item definitions using :> syntax
// have their specialization relationships properly captured and resolved.
func TestItemSpecializationWithColonGT(t *testing.T) {
	input := `
package TestPackage {
    item def Product;
    item def Widget :> Product;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Err())
	}

	model := result.Model

	// Find all items
	items := FindAll[*Item](model)
	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	var product, widget *Item
	for _, i := range items {
		switch i.Name() {
		case "Product":
			product = i
		case "Widget":
			widget = i
		}
	}

	if product == nil {
		t.Fatal("Product item not found")
	}
	if widget == nil {
		t.Fatal("Widget item not found")
	}

	// Verify Widget specializes Product
	if !widget.Specializes.IsResolved() {
		t.Error("Widget.Specializes should be resolved to Product")
	} else {
		if widget.Specializes.Resolved() != product {
			t.Error("Widget should specialize Product")
		}
	}
}

// TestSpecializationWithQualifiedName verifies that specializations with
// qualified names (e.g., Parent::Vehicle) are properly resolved.
func TestSpecializationWithQualifiedName(t *testing.T) {
	input := `
package ParentPackage {
    part def Vehicle;
}
package ChildPackage {
    part def Car :> ParentPackage::Vehicle;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Err())
	}

	model := result.Model

	parts := FindAll[*Part](model)

	var vehicle, car *Part
	for _, p := range parts {
		switch p.Name() {
		case "Vehicle":
			vehicle = p
		case "Car":
			car = p
		}
	}

	if vehicle == nil {
		t.Fatal("Vehicle part not found")
	}
	if car == nil {
		t.Fatal("Car part not found")
	}

	// Verify Car specializes Vehicle via qualified name
	if !car.Specializes.IsResolved() {
		t.Error("Car.Specializes should be resolved to Vehicle via qualified name")
	} else {
		if car.Specializes.Resolved() != vehicle {
			t.Error("Car should specialize ParentPackage::Vehicle")
		}
	}
}

// TestRedefinitionNameDerivation verifies that usage members declared with :>>
// and no explicit Identification get their name derived from the redefinition
// target, so QualifiedName does not end with a trailing "::".
func TestRedefinitionNameDerivation(t *testing.T) {
	input := `
package BugRepro {
    part def Base { part x; }
    part def Child :> Base {
        part :>> x : Base;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Err())
	}

	parts := FindAll[*Part](result.Model)

	var child *Part
	for _, p := range parts {
		if p.Name() == "Child" {
			child = p
		}
	}
	if child == nil {
		t.Fatal("Child part not found")
	}

	childParts := child.Parts()
	if len(childParts) != 1 {
		t.Fatalf("Expected 1 nested part in Child, got %d", len(childParts))
	}

	nested := childParts[0]
	if nested.Name() != "x" {
		t.Errorf("Expected nested part name %q, got %q", "x", nested.Name())
	}

	want := "BugRepro::Child::x"
	if got := nested.QualifiedName(); got != want {
		t.Errorf("QualifiedName: want %q, got %q", want, got)
	}
}

// TestRedefinitionNameDerivationItem verifies that item :>> foo inside a part def
// derives its name from the redefinition target.
func TestRedefinitionNameDerivationItem(t *testing.T) {
	// item :>> x inside a part def is parsed as itemUsage with empty identification;
	// the element name should be derived from the redefinition target.
	input := `
package BugRepro {
    part def Child {
        item :>> x : Base;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Err())
	}

	items := FindAll[*Item](result.Model)

	var usage *Item
	for _, i := range items {
		if !i.IsDefinition {
			usage = i
		}
	}
	if usage == nil {
		t.Fatal("item usage not found")
	}

	if usage.Name() != "x" {
		t.Errorf("Expected item name %q, got %q", "x", usage.Name())
	}
}

// TestRedefinitionNameDerivationAttributeWithName verifies that attribute foo :>> bar
// (with explicit name) correctly records the redefinition target and QN.
// Note: `attribute :>> x` without an explicit name is parsed by the ANTLR grammar
// as a defaultReferenceUsage (keyword "attribute" as element name), so the name
// derivation path in EnterAttributeUsage applies to other unnamed attribute usages.
func TestRedefinitionNameDerivationAttributeWithName(t *testing.T) {
	input := `
package BugRepro {
    part def Base {
        attribute x = 0;
    }
    part def Child :> Base {
        attribute y :>> x = 1;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Err())
	}

	parts := FindAll[*Part](result.Model)

	var child *Part
	for _, p := range parts {
		if p.Name() == "Child" {
			child = p
			break
		}
	}
	if child == nil {
		t.Fatal("Child part not found")
	}

	attrs := child.Attributes()
	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute in Child, got %d", len(attrs))
	}

	nested := attrs[0]
	if nested.Name() != "y" {
		t.Errorf("Expected attribute name %q, got %q", "y", nested.Name())
	}

	want := "BugRepro::Child::y"
	if got := nested.QualifiedName(); got != want {
		t.Errorf("QualifiedName: want %q, got %q", want, got)
	}
}

// TestSpecializationUniqueness verifies that each part maintains its own
// specialization reference independently.
func TestSpecializationUniqueness(t *testing.T) {
	input := `
package TestPackage {
    part def A;
    part def B;
    part def C :> A;
    part def D :> A;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Err())
	}

	model := result.Model

	parts := FindAll[*Part](model)

	var a, c, d *Part
	for _, p := range parts {
		switch p.Name() {
		case "A":
			a = p
		case "C":
			c = p
		case "D":
			d = p
		}
	}

	if a == nil {
		t.Fatal("Part A not found")
	}
	if c == nil {
		t.Fatal("Part C not found")
	}
	if d == nil {
		t.Fatal("Part D not found")
	}

	// Verify both C and D specialize A
	if !c.Specializes.IsResolved() || c.Specializes.Resolved() != a {
		t.Error("C should specialize A")
	}
	if !d.Specializes.IsResolved() || d.Specializes.Resolved() != a {
		t.Error("D should specialize A")
	}
}
