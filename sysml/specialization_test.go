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
		t.Fatalf("Failed to parse: %v", result.Errors)
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
		t.Fatalf("Failed to parse: %v", result.Errors)
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
		t.Fatalf("Failed to parse: %v", result.Errors)
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
		t.Fatalf("Failed to parse: %v", result.Errors)
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
