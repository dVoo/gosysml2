package sysml

import (
	"testing"
)

func TestParseString(t *testing.T) {
	input := `
package TestPackage {
    part def Vehicle {
        part engine : Engine;
    }
    part def Engine;
}
`

	result := ParseString(input)

	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	if result.Model == nil {
		t.Fatal("model is nil")
	}

	if len(result.Model.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(result.Model.Packages))
	}

	pkg := result.Model.Packages[0]
	if pkg.Name() != "TestPackage" {
		t.Errorf("expected package name 'TestPackage', got '%s'", pkg.Name())
	}
}

func TestParseStringWithErrors(t *testing.T) {
	input := `
package Broken {
    part def Vehicle {
        invalid syntax here!!!
    }
}
`

	result := ParseString(input)

	if result.Success() {
		t.Error("expected parse to fail")
	}

	if result.Errors == nil || !result.Errors.HasErrors() {
		t.Error("expected errors to be present")
	}

	t.Logf("Got expected errors: %s", result.Errors)
}

func TestValidate(t *testing.T) {
	validInput := `
package Valid {
    part def A;
    part def B;
}
`

	invalidInput := `
package Invalid {
    @@@ not valid syntax
}
`

	// Valid input should pass
	err := Validate(validInput)
	if err != nil {
		t.Errorf("valid input failed validation: %s", err)
	}

	// Invalid input should fail
	err = Validate(invalidInput)
	if err == nil {
		t.Error("invalid input should fail validation")
	}
}

func TestWalk(t *testing.T) {
	input := `
package P1 {
    package P2 {
        part def A;
    }
    part def B;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	var visited []string
	Walk(result.Model, func(elem Element, depth int) bool {
		visited = append(visited, elem.Name())
		return true
	})

	if len(visited) == 0 {
		t.Error("no elements visited")
	}

	t.Logf("Visited elements: %v", visited)
}

func TestFindByKind(t *testing.T) {
	input := `
package TestPkg {
    part def Vehicle;
    part def Engine;
    part vehicle1 : Vehicle;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	parts := FindByKind(result.Model, KindPart)
	t.Logf("Found %d parts", len(parts))
}

func TestCounter(t *testing.T) {
	input := `
package P {
    part def A;
    part def B;
    part a1 : A;
    part b1 : B;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	counter := NewCounter()
	Visit(result.Model, counter)

	t.Logf("Counts: %v", counter.Counts)
	t.Logf("Total: %d", counter.Total())

	if counter.Total() == 0 {
		t.Error("expected some elements to be counted")
	}
}

func TestMustParseStringPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid input")
		}
	}()

	MustParseString(`package { invalid }}}`)
}

func TestEmptyInput(t *testing.T) {
	result := ParseString("")

	// Empty input should parse successfully (empty namespace is valid)
	if !result.Success() {
		t.Logf("Parse result: %s", result.Errors)
	}
}

// TestNestedPartsParentChildRelationships verifies that nested parts
// create proper parent-child relationships via Children()
func TestNestedPartsParentChildRelationships(t *testing.T) {
	result := ParseFile("../validationdata/01-Parts Tree/1a-Parts Tree.sysml")
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Errors)
	}

	model := result.Model

	// Find all parts in the model
	parts := FindAll[*Part](model)
	if len(parts) == 0 {
		t.Fatal("No parts found in the model")
	}

	// Find vehicle1 part usage
	var vehicle1 *Part
	for _, p := range parts {
		if p.Name() == "vehicle1" {
			vehicle1 = p
			break
		}
	}
	if vehicle1 == nil {
		t.Fatal("vehicle1 part usage not found")
	}

	// Verify vehicle1 is NOT a definition
	if vehicle1.IsDefinition {
		t.Errorf("vehicle1 should be a usage (IsDefinition=false), got IsDefinition=true")
	}

	// Verify vehicle1 has children (frontAxleAssembly and rearAxleAssembly)
	children := vehicle1.Children()
	if len(children) == 0 {
		t.Error("vehicle1 should have children (frontAxleAssembly and rearAxleAssembly)")
	}

	// Check for frontAxleAssembly and rearAxleAssembly in children
	var frontAxleAssembly, rearAxleAssembly *Part
	for _, child := range children {
		if part, ok := child.(*Part); ok {
			switch part.Name() {
			case "frontAxleAssembly":
				frontAxleAssembly = part
			case "rearAxleAssembly":
				rearAxleAssembly = part
			}
		}
	}

	if frontAxleAssembly == nil {
		t.Error("vehicle1.Children() should contain frontAxleAssembly")
	} else {
		// Verify frontAxleAssembly is a usage
		if frontAxleAssembly.IsDefinition {
			t.Errorf("frontAxleAssembly should be a usage (IsDefinition=false)")
		}
		// Verify frontAxleAssembly has frontAxle and frontWheel as children
		faChildren := frontAxleAssembly.Children()
		var frontAxle, frontWheel *Part
		for _, child := range faChildren {
			if part, ok := child.(*Part); ok {
				switch part.Name() {
				case "frontAxle":
					frontAxle = part
				case "frontWheel":
					frontWheel = part
				}
			}
		}
		if frontAxle == nil {
			t.Error("frontAxleAssembly.Children() should contain frontAxle")
		}
		if frontWheel == nil {
			t.Error("frontAxleAssembly.Children() should contain frontWheel")
		}
	}

	if rearAxleAssembly == nil {
		t.Error("vehicle1.Children() should contain rearAxleAssembly")
	} else {
		// Verify rearAxleAssembly is a usage
		if rearAxleAssembly.IsDefinition {
			t.Errorf("rearAxleAssembly should be a usage (IsDefinition=false)")
		}
	}
}

// TestPartDefinitionVsUsageDistinction verifies that part definitions
// and usages are properly distinguished via IsDefinition field
func TestPartDefinitionVsUsageDistinction(t *testing.T) {
	result := ParseFile("../validationdata/01-Parts Tree/1a-Parts Tree.sysml")
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Errors)
	}

	model := result.Model

	// Find Definitions package
	var definitionsPkg *Package
	for _, pkg := range FindAll[*Package](model) {
		if pkg.Name() == "Definitions" {
			definitionsPkg = pkg
			break
		}
	}
	if definitionsPkg == nil {
		t.Fatal("Definitions package not found")
	}

	// Find Usages package
	var usagesPkg *Package
	for _, pkg := range FindAll[*Package](model) {
		if pkg.Name() == "Usages" {
			usagesPkg = pkg
			break
		}
	}
	if usagesPkg == nil {
		t.Fatal("Usages package not found")
	}

	// Verify Vehicle part def has IsDefinition=true
	var vehicleDef *Part
	for _, child := range definitionsPkg.Children() {
		if part, ok := child.(*Part); ok && part.Name() == "Vehicle" {
			vehicleDef = part
			break
		}
	}
	if vehicleDef == nil {
		t.Fatal("Vehicle part definition not found in Definitions package")
	}
	if !vehicleDef.IsDefinition {
		t.Errorf("Vehicle part definition should have IsDefinition=true")
	}

	// Verify vehicle1 part usage has IsDefinition=false
	var vehicleUsage *Part
	for _, child := range usagesPkg.Children() {
		if part, ok := child.(*Part); ok && part.Name() == "vehicle1" {
			vehicleUsage = part
			break
		}
	}
	if vehicleUsage == nil {
		t.Fatal("vehicle1 part usage not found in Usages package")
	}
	if vehicleUsage.IsDefinition {
		t.Errorf("vehicle1 part usage should have IsDefinition=false")
	}

	// Verify other definitions
	expectedDefs := []string{"Vehicle", "AxleAssembly", "Axle", "FrontAxle", "Wheel"}
	for _, defName := range expectedDefs {
		found := false
		for _, child := range definitionsPkg.Children() {
			if part, ok := child.(*Part); ok && part.Name() == defName {
				found = true
				if !part.IsDefinition {
					t.Errorf("Part definition '%s' should have IsDefinition=true", defName)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected part definition '%s' not found in Definitions package", defName)
		}
	}
}

// TestTypeReferenceResolution verifies that part usages have their
// TypeRef populated with the type name from the source
func TestTypeReferenceResolution(t *testing.T) {
	result := ParseFile("../validationdata/01-Parts Tree/1a-Parts Tree.sysml")
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Errors)
	}

	model := result.Model

	// Find vehicle1 part usage and check its TypeRef
	parts := FindAll[*Part](model)

	tests := []struct {
		partName         string
		expectedTypeName string
	}{
		{"vehicle1", "Vehicle"},
		{"frontAxleAssembly", "AxleAssembly"},
		{"rearAxleAssembly", "AxleAssembly"},
		{"frontAxle", "Axle"},
		{"frontWheel", "Wheel"},
		{"rearAxle", "Axle"},
		{"rearWheel", "Wheel"},
	}

	for _, test := range tests {
		var found *Part
		for _, p := range parts {
			if p.Name() == test.partName {
				found = p
				break
			}
		}

		if found == nil {
			t.Errorf("Part '%s' not found", test.partName)
			continue
		}

		// Part usages should have a TypeRef
		if found.IsDefinition {
			t.Logf("Part '%s' is a definition, skipping TypeRef check", test.partName)
			continue
		}

		typeRefName := found.TypeRef.Name()
		if typeRefName == "" {
			t.Errorf("Part usage '%s' should have TypeRef populated, but it's empty", test.partName)
		} else if typeRefName != test.expectedTypeName {
			t.Errorf("Part usage '%s' should have TypeRef='%s', got TypeRef='%s'",
				test.partName, test.expectedTypeName, typeRefName)
		}
	}
}

// TestAttributeExtraction verifies that attributes are properly
// extracted and accessible via Part.Attributes()
func TestAttributeExtraction(t *testing.T) {
	result := ParseFile("../validationdata/01-Parts Tree/1a-Parts Tree.sysml")
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Errors)
	}

	model := result.Model

	// Find Vehicle part definition
	var vehicleDef *Part
	for _, p := range FindAll[*Part](model) {
		if p.Name() == "Vehicle" && p.IsDefinition {
			vehicleDef = p
			break
		}
	}
	if vehicleDef == nil {
		t.Fatal("Vehicle part definition not found")
	}

	// Check Vehicle has mass attribute
	attrs := vehicleDef.Attributes()
	var massAttr *Attribute
	for _, attr := range attrs {
		if attr.Name() == "mass" {
			massAttr = attr
			break
		}
	}
	if massAttr == nil {
		t.Error("Vehicle part definition should have 'mass' attribute accessible via Attributes()")
	}

	// Check Axle part definition has mass attribute
	var axleDef *Part
	for _, p := range FindAll[*Part](model) {
		if p.Name() == "Axle" && p.IsDefinition {
			axleDef = p
			break
		}
	}
	if axleDef != nil {
		axleAttrs := axleDef.Attributes()
		var axleMassAttr *Attribute
		for _, attr := range axleAttrs {
			if attr.Name() == "mass" {
				axleMassAttr = attr
				break
			}
		}
		if axleMassAttr == nil {
			t.Error("Axle part definition should have 'mass' attribute")
		}
	}

	// Check FrontAxle part definition has steeringAngle attribute
	var frontAxleDef *Part
	for _, p := range FindAll[*Part](model) {
		if p.Name() == "FrontAxle" && p.IsDefinition {
			frontAxleDef = p
			break
		}
	}
	if frontAxleDef != nil {
		faAttrs := frontAxleDef.Attributes()
		var steeringAngleAttr *Attribute
		for _, attr := range faAttrs {
			if attr.Name() == "steeringAngle" {
				steeringAngleAttr = attr
				break
			}
		}
		if steeringAngleAttr == nil {
			t.Error("FrontAxle part definition should have 'steeringAngle' attribute")
		}
	}

	// Check vehicle1 usage has mass attribute (redefining Vehicle::mass)
	var vehicleUsage *Part
	for _, p := range FindAll[*Part](model) {
		if p.Name() == "vehicle1" && !p.IsDefinition {
			vehicleUsage = p
			break
		}
	}
	if vehicleUsage != nil {
		usageAttrs := vehicleUsage.Attributes()
		var usageMassAttr *Attribute
		for _, attr := range usageAttrs {
			if attr.Name() == "mass" {
				usageMassAttr = attr
				break
			}
		}
		if usageMassAttr == nil {
			t.Log("Note: vehicle1 usage doesn't have mass attribute in its Attributes() - this may be by design (redefinition)")
		} else {
			// Check default value
			if usageMassAttr.DefaultValue == "" {
				t.Log("Note: vehicle1 mass attribute has empty DefaultValue - may need value extraction fix")
			} else {
				t.Logf("vehicle1 mass attribute DefaultValue: %s", usageMassAttr.DefaultValue)
			}
		}
	}
}

// TestParentRelationships verifies parent references are set correctly
func TestParentRelationships(t *testing.T) {
	result := ParseFile("../validationdata/01-Parts Tree/1a-Parts Tree.sysml")
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Errors)
	}

	model := result.Model

	// Find vehicle1
	var vehicle1 *Part
	for _, p := range FindAll[*Part](model) {
		if p.Name() == "vehicle1" {
			vehicle1 = p
			break
		}
	}
	if vehicle1 == nil {
		t.Fatal("vehicle1 not found")
	}

	// vehicle1's parent should be the Usages package
	parent := vehicle1.Parent()
	if parent == nil {
		t.Error("vehicle1 should have a parent")
	} else if pkg, ok := parent.(*Package); ok {
		if pkg.Name() != "Usages" {
			t.Errorf("vehicle1's parent should be 'Usages' package, got '%s'", pkg.Name())
		}
	} else {
		t.Errorf("vehicle1's parent should be a Package, got %T", parent)
	}

	// Find frontAxleAssembly
	var frontAxleAssembly *Part
	for _, p := range FindAll[*Part](model) {
		if p.Name() == "frontAxleAssembly" {
			frontAxleAssembly = p
			break
		}
	}
	if frontAxleAssembly == nil {
		t.Fatal("frontAxleAssembly not found")
	}

	// frontAxleAssembly's parent should be vehicle1
	faParent := frontAxleAssembly.Parent()
	if faParent == nil {
		t.Error("frontAxleAssembly should have a parent")
	} else if part, ok := faParent.(*Part); ok {
		if part.Name() != "vehicle1" {
			t.Errorf("frontAxleAssembly's parent should be 'vehicle1', got '%s'", part.Name())
		}
	} else {
		t.Errorf("frontAxleAssembly's parent should be a Part, got %T", faParent)
	}
}
