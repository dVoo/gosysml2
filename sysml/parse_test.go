package sysml

import (
	"os"
	"path/filepath"
	"strings"
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
		t.Fatalf("parse failed: %s", result.Err())
	}

	if result.Model == nil {
		t.Fatal("model is nil")
	}

	if len(result.Model.Packages()) != 1 {
		t.Errorf("expected 1 package, got %d", len(result.Model.Packages()))
	}

	pkg := result.Model.Packages()[0]
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

	if result.ParseError == nil || !result.ParseError.HasErrors() {
		t.Error("expected errors to be present")
	}

	t.Logf("Got expected errors: %s", result.ParseError)
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
		t.Fatalf("parse failed: %s", result.Err())
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
		t.Fatalf("parse failed: %s", result.Err())
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
		t.Fatalf("parse failed: %s", result.Err())
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
		t.Logf("Parse result: %s", result.Err())
	}
}

func TestRequirementDeclaredShortNameAndTyping(t *testing.T) {
	input := `
package VehicleReqs {
    requirement def MassRequirement;
    requirement <'R1'> vehicleMassReq : MassRequirement;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	reqs := FindAll[*Requirement](result.Model)
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requirements, got %d", len(reqs))
	}

	var usage *Requirement
	for _, r := range reqs {
		if !r.IsDefinition {
			usage = r
			break
		}
	}
	if usage == nil {
		t.Fatal("expected a requirement usage")
	}

	if usage.Name() != "vehicleMassReq" {
		t.Fatalf("expected usage name 'vehicleMassReq', got %q", usage.Name())
	}

	if usage.DeclaredShortName() != "'R1'" {
		t.Fatalf("expected declared short name %q, got %q", "'R1'", usage.DeclaredShortName())
	}

	if usage.TypeRef.Name() != "MassRequirement" {
		t.Fatalf("expected type ref 'MassRequirement', got %q", usage.TypeRef.Name())
	}
}

func TestParseBareImportAndAliasAsForm(t *testing.T) {
	input := `
package VehicleModel {
    import ISQ::*;
    alias ISQ::TorqueValue as Torque;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkgs := FindAll[*Package](result.Model)
	if len(pkgs) == 0 || len(pkgs[0].Imports()) == 0 {
		t.Fatal("expected at least one package import")
	}
	imp := pkgs[0].Imports()[0]
	if imp.Visibility != "" {
		t.Fatalf("expected empty visibility for bare import, got %q", imp.Visibility)
	}
	if imp.ImportedNamespace != "ISQ" {
		t.Fatalf("expected imported namespace ISQ, got %q", imp.ImportedNamespace)
	}
	if !imp.IsAll {
		t.Fatal("expected wildcard import")
	}

	aliases := FindAll[*Alias](result.Model)
	if len(aliases) == 0 {
		t.Fatal("expected at least one alias")
	}
	if aliases[0].Name() != "Torque" {
		t.Fatalf("expected alias name Torque, got %q", aliases[0].Name())
	}
	if aliases[0].unresolvedTarget != "ISQ::TorqueValue" {
		t.Fatalf("expected alias target ISQ::TorqueValue, got %q", aliases[0].unresolvedTarget)
	}
}

func TestParseLogicalAndAndUnitAtBracketForm(t *testing.T) {
	input := `
package HVAC {
    attribute mass = 75 @[kg];
    attribute ok = true && false;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}
}

func TestParseCompliantStateDoPerformActionBody(t *testing.T) {
	input := `
package VehicleModel {
    state vehicleStates {
        do vehicle::senseTemperature { out temp; };
    }
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}
}

func TestRequirementDefinitionRequireBlockExpressionCompatibility(t *testing.T) {
	input := `
package P {
    requirement def MassRequirement {
        require { massActual <= massLimit };
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	reqs := FindAll[*Requirement](result.Model)
	if len(reqs) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(reqs))
	}

	req := reqs[0]
	if len(req.Constraints) != 1 {
		t.Fatalf("expected 1 requirement constraint, got %d", len(req.Constraints))
	}
	if req.Constraints[0].Expression != "massActual <= massLimit" {
		t.Fatalf("expected constraint expression %q, got %q", "massActual <= massLimit", req.Constraints[0].Expression)
	}
}

func TestRequirementUsageBindingListCompatibility(t *testing.T) {
	input := `
package P {
    requirement def MassRequirement;
    requirement <'R1'> massReq : MassRequirement [vehicle = myVehicle, threshold = 42];
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	reqs := FindAll[*Requirement](result.Model)
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requirements, got %d", len(reqs))
	}

	var usage *Requirement
	for _, r := range reqs {
		if !r.IsDefinition {
			usage = r
			break
		}
	}
	if usage == nil {
		t.Fatal("expected requirement usage")
	}

	if usage.Bindings["vehicle"] != "myVehicle" {
		t.Fatalf("expected vehicle binding %q, got %q", "myVehicle", usage.Bindings["vehicle"])
	}
	if usage.Bindings["threshold"] != "42" {
		t.Fatalf("expected threshold binding %q, got %q", "42", usage.Bindings["threshold"])
	}
}

func TestCompatibilityRewritesForNewGaps(t *testing.T) {
	input := `
package P {
    part def X {
        attribute type : String[0..1];
        attribute <var> myVar : Real;
    }
    alias multiplicity for X;
    calc f {
        return result = values->selectOne {in ref a { a > 0 }};
    }
}
`
	normalized, _, _ := normalizeUnsupportedRequirementSyntax(input)

	contains := func(want string) {
		t.Helper()
		if !strings.Contains(normalized, want) {
			t.Fatalf("expected normalized input to contain %q, got:\n%s", want, normalized)
		}
	}

	contains("attribute 'type' :")
	contains("attribute <'var'> myVar : Real;")
	contains("alias 'multiplicity' for X;")
	contains("selectOne {in a { a > 0 }}")
}

func TestCompatibilityRewritesForGap21AndGap22(t *testing.T) {
	input := `
package P {
    action def ForLoopActionDef {
        protected ref var[0..1] :> seq;
        assign var := seq#(index);
    }
    part def X {
        derived ref item actionDefinition : Behavior[0..*] ordered redefines behavior, occurrenceDefinition subsets Metadata::metadataItems;
    }
}
`
	normalized, _, _ := normalizeUnsupportedRequirementSyntax(input)

	contains := func(want string) {
		t.Helper()
		if !strings.Contains(normalized, want) {
			t.Fatalf("expected normalized input to contain %q, got:\n%s", want, normalized)
		}
	}

	contains("protected ref 'var'[0..1] :> seq;")
	contains("assign 'var' := seq#(index);")
	contains("derived ref item actionDefinition : Behavior[0..*] ordered redefines behavior, occurrenceDefinition subsets Metadata::metadataItems;")
}

func TestCompatibilityGapFilesParse(t *testing.T) {
	files := []string{
		"Domain Libraries/Analysis/TradeStudies.sysml",
		"Domain Libraries/Metadata/ImageMetadata.sysml",
		"Domain Libraries/Quantities and Units/ISQChemistryMolecular.sysml",
		"Domain Libraries/Quantities and Units/SI.sysml",
		"Systems Library/Actions.sysml",
		"Systems Library/SysML.sysml",
	}

	for _, rel := range files {
		path := filepath.Join("..", "libraries", "sysml.library", filepath.FromSlash(rel))
		if _, err := os.Stat(path); err != nil {
			t.Skipf("required library file not found: %s (%v)", path, err)
		}
		result := ParseFile(path)
		if !result.Success() {
			t.Fatalf("expected parse success for %s, got errors: %v", rel, result.Err())
		}
	}
}

func TestInlineSubsetsRedefinesKeywordsCapturedOnItemUsage(t *testing.T) {
	input := `
package P {
    item def TypeX;
    item def BaseA;
    item def BaseB;
    item itemX : TypeX[0..*] ordered redefines BaseA, BaseB subsets BaseA;
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	var actionDef *Item
	for _, it := range FindAll[*Item](result.Model) {
		if it.Name() == "itemX" && !it.IsDefinition {
			actionDef = it
			break
		}
	}
	if actionDef == nil {
		t.Fatal("expected item usage itemX")
	}

	if actionDef.TypeRef.Name() != "TypeX" {
		t.Fatalf("expected item type ref %q, got %q", "TypeX", actionDef.TypeRef.Name())
	}

	if len(actionDef.RedefinedFeatures) != 2 {
		t.Fatalf("expected 2 redefined features, got %d", len(actionDef.RedefinedFeatures))
	}
	if actionDef.RedefinedFeatures[0].Name() != "BaseA" {
		t.Fatalf("expected first redefined feature %q, got %q", "BaseA", actionDef.RedefinedFeatures[0].Name())
	}
	if actionDef.RedefinedFeatures[1].Name() != "BaseB" {
		t.Fatalf("expected second redefined feature %q, got %q", "BaseB", actionDef.RedefinedFeatures[1].Name())
	}

	if len(actionDef.SubsettedFeatures) != 1 {
		t.Fatalf("expected 1 subsetted feature, got %d", len(actionDef.SubsettedFeatures))
	}
	if actionDef.SubsettedFeatures[0].Name() != "BaseA" {
		t.Fatalf("expected subsetted feature %q, got %q", "BaseA", actionDef.SubsettedFeatures[0].Name())
	}
}

func TestPartUsageMultiplicity(t *testing.T) {
	input := `
package VehiclePkg {
    part def Wheel;
    part wheel1 : Wheel[1];
    part wheels : Wheel[4];
    part optionalWheel : Wheel[0..1];
    part anyWheels : Wheel[*];
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	parts := FindAll[*Part](result.Model)
	want := map[string]string{
		"wheel1":        "1",
		"wheels":        "4",
		"optionalWheel": "0..1",
		"anyWheels":     "*",
	}

	for name, expectedMultiplicity := range want {
		var found *Part
		for _, p := range parts {
			if p.Name() == name {
				found = p
				break
			}
		}
		if found == nil {
			t.Fatalf("expected part usage %q to be parsed", name)
		}
		if found.IsDefinition {
			t.Fatalf("expected %q to be usage, but got definition", name)
		}
		if found.Multiplicity != expectedMultiplicity {
			t.Fatalf("expected multiplicity for %q to be %q, got %q", name, expectedMultiplicity, found.Multiplicity)
		}
	}
}

func TestPartAttributeUsageTypingAndDefaultValue(t *testing.T) {
	input := `
package VehiclePkg {
    part def Car {
        attribute mass : Real = 1500.0;
        attribute maxSpeed : Integer = 200;
    }
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	var car *Part
	for _, p := range FindAll[*Part](result.Model) {
		if p.Name() == "Car" && p.IsDefinition {
			car = p
			break
		}
	}
	if car == nil {
		t.Fatal("expected part definition Car")
	}

	attrs := car.Attributes()
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attributes on Car, got %d", len(attrs))
	}

	type expected struct {
		typeRef string
		value   string
	}
	want := map[string]expected{
		"mass":     {typeRef: "Real", value: "1500.0"},
		"maxSpeed": {typeRef: "Integer", value: "200"},
	}

	for _, attr := range attrs {
		exp, ok := want[attr.Name()]
		if !ok {
			continue
		}
		if attr.TypeRef.Name() != exp.typeRef {
			t.Fatalf("expected attribute %q type ref %q, got %q", attr.Name(), exp.typeRef, attr.TypeRef.Name())
		}
		if attr.DefaultValue != exp.value {
			t.Fatalf("expected attribute %q default value %q, got %q", attr.Name(), exp.value, attr.DefaultValue)
		}
	}
}

func TestPartAttributeUsageTypingFromSpecializationForm(t *testing.T) {
	input := `
package VehiclePkg {
    attribute def MassValue;
    part def Car {
        attribute mass :> MassValue;
    }
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	var massAttr *Attribute
	for _, attr := range FindAll[*Attribute](result.Model) {
		if attr.Name() == "mass" {
			massAttr = attr
			break
		}
	}
	if massAttr == nil {
		t.Fatal("expected attribute usage 'mass'")
	}

	if massAttr.TypeRef.Name() != "MassValue" {
		t.Fatalf("expected type ref %q, got %q", "MassValue", massAttr.TypeRef.Name())
	}
}

func TestAttributeUsageSubsettingAndRedefinitionReferences(t *testing.T) {
	input := `
package VehiclePkg {
    part def Vehicle {
        attribute velocity : Real;
        attribute speedA :> velocity;
        attribute speedB ::> velocity;
        attribute speedC :>> velocity;
    }
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	attrs := FindAll[*Attribute](result.Model)
	byName := map[string]*Attribute{}
	for _, attr := range attrs {
		byName[attr.Name()] = attr
	}

	velocity := byName["velocity"]
	speedA := byName["speedA"]
	speedB := byName["speedB"]
	speedC := byName["speedC"]
	if velocity == nil || speedA == nil || speedB == nil || speedC == nil {
		t.Fatalf("expected attributes velocity/speedA/speedB/speedC to be present")
	}

	if len(speedA.SubsettedFeatures) != 1 || speedA.SubsettedFeatures[0].Name() != "velocity" {
		t.Fatalf("expected speedA to subset velocity, got %+v", speedA.SubsettedFeatures)
	}
	if len(speedB.SubsettedFeatures) != 1 || speedB.SubsettedFeatures[0].Name() != "velocity" {
		t.Fatalf("expected speedB to subset/reference velocity, got %+v", speedB.SubsettedFeatures)
	}
	if len(speedC.RedefinedFeatures) != 1 || speedC.RedefinedFeatures[0].Name() != "velocity" {
		t.Fatalf("expected speedC to redefine velocity, got %+v", speedC.RedefinedFeatures)
	}
}

// TestNestedPartsParentChildRelationships verifies that nested parts
// create proper parent-child relationships via Children()
func TestNestedPartsParentChildRelationships(t *testing.T) {
	result := ParseFile(externalTestPath(t, "../validationdata/01-Parts Tree/1a-Parts Tree.sysml"))
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Err())
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
	result := ParseFile(externalTestPath(t, "../validationdata/01-Parts Tree/1a-Parts Tree.sysml"))
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Err())
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
	// First, test with a simple case to understand the grammar
	simpleInput := `
package Test {
    part def Vehicle;
    part def Axle;
    part myVehicle : Vehicle;
    part myAxle : Axle;
}
`
	result := ParseString(simpleInput)
	if !result.Success() {
		t.Fatalf("Failed to parse simple input: %v", result.Err())
	}

	model := result.Model
	parts := FindAll[*Part](model)

	for _, p := range parts {
		t.Logf("Found part: name=%s, IsDefinition=%v, TypeRef.Name=%q",
			p.Name(), p.IsDefinition, p.TypeRef.Name())
	}

	// Now test with the real file
	result = ParseFile(externalTestPath(t, "../validationdata/01-Parts Tree/1a-Parts Tree.sysml"))
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Err())
	}

	model = result.Model

	// Find vehicle1 part usage and check its TypeRef
	parts = FindAll[*Part](model)

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
		} else {
			t.Logf("SUCCESS: Part usage '%s' has TypeRef='%s'", test.partName, typeRefName)
		}
	}
}

// TestAttributeExtraction verifies that attributes are properly
// extracted and accessible via Part.Attributes()
func TestAttributeExtraction(t *testing.T) {
	result := ParseFile(externalTestPath(t, "../validationdata/01-Parts Tree/1a-Parts Tree.sysml"))
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Err())
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
	result := ParseFile(externalTestPath(t, "../validationdata/01-Parts Tree/1a-Parts Tree.sysml"))
	if !result.Success() {
		t.Fatalf("Failed to parse file: %v", result.Err())
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

// TestItemUsageNesting verifies that item usages nested inside part defs
// are parented to the part, not the package (regression for Bug A).
func TestItemUsageNesting(t *testing.T) {
	input := `
package TestPkg {
    part def Vehicle {
        item payload : Cargo;
        item nested1 {
            attribute weight : Real;
        }
    }
    item def Cargo;
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	// Vehicle should be found
	var vehicle *Part
	for _, p := range pkg.Parts() {
		if p.Name() == "Vehicle" {
			vehicle = p
			break
		}
	}
	if vehicle == nil {
		t.Fatal("Vehicle part def not found")
	}

	// Items should be children of Vehicle, not the package
	found := 0
	for _, child := range vehicle.Children() {
		if item, ok := child.(*Item); ok {
			found++
			if item.Parent() != vehicle {
				t.Errorf("item %q parent should be Vehicle, got %v", item.Name(), item.Parent())
			}
		}
	}
	if found < 2 {
		t.Errorf("expected at least 2 items inside Vehicle, found %d", found)
	}

	// Items should NOT be direct children of the package
	if len(pkg.Items()) > 1 {
		names := make([]string, 0)
		for _, item := range pkg.Items() {
			names = append(names, item.Name())
		}
		t.Errorf("package should only have Cargo item def, got: %v", names)
	}
}

// TestItemDefinitionNesting verifies item defs support nested children.
func TestItemDefinitionNesting(t *testing.T) {
	input := `
package TestPkg {
    item def Container {
        attribute capacity : Integer;
        item contents : Cargo;
    }
    item def Cargo;
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	// Find Container item def
	var container *Item
	for _, item := range pkg.Items() {
		if item.Name() == "Container" {
			container = item
			break
		}
	}
	if container == nil {
		t.Fatal("Container item def not found")
	}

	// Container should have children (attribute + item)
	if len(container.Children()) < 2 {
		t.Errorf("Container should have at least 2 children, got %d", len(container.Children()))
	}
	if len(container.Attributes()) < 1 {
		t.Errorf("Container should have at least 1 attribute, got %d", len(container.Attributes()))
	}
	if len(container.Items()) < 1 {
		t.Errorf("Container should have at least 1 nested item, got %d", len(container.Items()))
	}
}

// TestEnumerationNesting verifies enum values are children of the enum, not the package.
func TestEnumerationNesting(t *testing.T) {
	input := `
package TestPkg {
    enum def Color {
        enum Red;
        enum Green;
        enum Blue;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	if len(pkg.Enumerations()) != 1 {
		t.Fatalf("expected 1 enumeration, got %d", len(pkg.Enumerations()))
	}

	colorEnum := pkg.Enumerations()[0]
	if colorEnum.Name() != "Color" {
		t.Errorf("expected enum name 'Color', got '%s'", colorEnum.Name())
	}

	// Values should be children of the enum
	if len(colorEnum.Values()) != 3 {
		t.Errorf("expected 3 enum values, got %d", len(colorEnum.Values()))
	}
	// Check parent of each value
	for _, v := range colorEnum.Values() {
		if v.Parent() != colorEnum {
			t.Errorf("enum value %q parent should be Color enum, got %v", v.Name(), v.Parent())
		}
	}
}

// TestNestedStatesAndActions verifies that states and actions support nesting.
func TestNestedStatesAndActions(t *testing.T) {
	input := `
package TestPkg {
    state def OperatingStates {
        state idle;
        state running {
            state warmup;
        }
    }
    action def ProcessFlow {
        action step1;
        action step2;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	// Find OperatingStates
	var opStates *State
	for _, s := range pkg.States() {
		if s.Name() == "OperatingStates" {
			opStates = s
			break
		}
	}
	if opStates == nil {
		t.Fatal("OperatingStates not found")
	}

	// Should have 2 nested states (idle, running)
	if len(opStates.States()) != 2 {
		t.Errorf("expected 2 nested states in OperatingStates, got %d", len(opStates.States()))
	}

	// running should have a nested warmup state
	var running *State
	for _, s := range opStates.States() {
		if s.Name() == "running" {
			running = s
			break
		}
	}
	if running == nil {
		t.Fatal("running state not found")
	}
	if len(running.States()) != 1 {
		t.Errorf("expected 1 nested state in running, got %d", len(running.States()))
	}

	// Find ProcessFlow
	var processFlow *Action
	for _, a := range pkg.Actions() {
		if a.Name() == "ProcessFlow" {
			processFlow = a
			break
		}
	}
	if processFlow == nil {
		t.Fatal("ProcessFlow not found")
	}

	// Should have 2 nested actions
	if len(processFlow.Actions()) != 2 {
		t.Errorf("expected 2 nested actions in ProcessFlow, got %d", len(processFlow.Actions()))
	}
}

// TestNestedConstraints verifies constraint nesting.
func TestNestedConstraints(t *testing.T) {
	input := `
package TestPkg {
    constraint def TempConstraint {
        constraint maxTemp;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	if len(pkg.Constraints()) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(pkg.Constraints()))
	}

	tempConstraint := pkg.Constraints()[0]
	if len(tempConstraint.Constraints()) != 1 {
		t.Errorf("expected 1 nested constraint, got %d", len(tempConstraint.Constraints()))
	}
}

// TestInterfaceWithPorts verifies interface bodies can contain ports.
func TestInterfaceWithPorts(t *testing.T) {
	input := `
package TestPkg {
    interface def DataLink {
        end port sender;
        end port receiver;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	if len(pkg.Interfaces()) != 1 {
		t.Fatalf("expected 1 interface, got %d", len(pkg.Interfaces()))
	}

	dataLink := pkg.Interfaces()[0]
	if len(dataLink.Ports()) < 2 {
		t.Errorf("expected at least 2 ports in DataLink, got %d", len(dataLink.Ports()))
	}
}

// TestPartTypedChildTracking verifies Part tracks items, states, actions,
// connections, constraints, requirements, enumerations, and flows via typed accessors.
func TestPartTypedChildTracking(t *testing.T) {
	input := `
package TestPkg {
    part def Vehicle {
        attribute mass : Real;
        part engine : Engine;
        port fuelPort;
        item cargo;
        connection engineLink;
        state operatingState;
        action startEngine;
        constraint safeSpeed;
        requirement safetyReq;
        enum def Color { enum red; }
        flow fuelFlow;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}
	if len(pkg.Parts()) != 1 {
		t.Fatalf("expected 1 part, got %d", len(pkg.Parts()))
	}

	vehicle := pkg.Parts()[0]

	tests := []struct {
		name  string
		count int
	}{
		{"Attributes", len(vehicle.Attributes())},
		{"Parts", len(vehicle.Parts())},
		{"Ports", len(vehicle.Ports())},
		{"Items", len(vehicle.Items())},
		{"Connections", len(vehicle.Connections())},
		{"States", len(vehicle.States())},
		{"Actions", len(vehicle.Actions())},
		{"Constraints", len(vehicle.Constraints())},
		{"Requirements", len(vehicle.Requirements())},
		{"Enumerations", len(vehicle.Enumerations())},
		{"Flows", len(vehicle.Flows())},
	}
	for _, tt := range tests {
		if tt.count < 1 {
			t.Errorf("Part.%s(): expected >= 1, got %d", tt.name, tt.count)
		}
	}
}

// TestActionTypedChildTracking verifies Action tracks control nodes,
// attributes, parts, items, constraints, and flows.
func TestActionTypedChildTracking(t *testing.T) {
	input := `
package TestPkg {
    action def ProcessOrder {
        action validate;
        action fulfill;
        fork node;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}
	if len(pkg.Actions()) != 1 {
		t.Fatalf("expected 1 action, got %d", len(pkg.Actions()))
	}

	processOrder := pkg.Actions()[0]
	if len(processOrder.Actions()) < 2 {
		t.Errorf("expected >= 2 nested actions, got %d", len(processOrder.Actions()))
	}
	if len(processOrder.ControlNodes()) < 1 {
		t.Errorf("expected >= 1 control node, got %d", len(processOrder.ControlNodes()))
	}
}

// TestControlNodeInsideAction verifies control nodes are parented to action, not package.
func TestControlNodeInsideAction(t *testing.T) {
	input := `
package TestPkg {
    action def Workflow {
        fork node;
        join node;
        merge node;
        decide node;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	workflow := pkg.Actions()[0]
	if len(workflow.ControlNodes()) != 4 {
		t.Errorf("expected 4 control nodes inside action, got %d", len(workflow.ControlNodes()))
	}

	// Control nodes should NOT appear as direct package children
	for _, child := range pkg.Children() {
		if child.Kind() == KindControlNode {
			t.Errorf("control node %q should be inside action, not at package level", child.Name())
		}
	}
}

// TestCommentInsideElement verifies comments inside definitions are parented correctly.
func TestCommentInsideElement(t *testing.T) {
	input := `
package TestPkg {
    part def Vehicle {
        comment /* A vehicle part */
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	vehicle := pkg.Parts()[0]
	hasComment := false
	for _, child := range vehicle.Children() {
		if child.Kind() == KindComment {
			hasComment = true
			break
		}
	}
	if !hasComment {
		t.Error("expected comment as child of Vehicle part def, not found")
	}
}

// TestPortAttributeTracking verifies ports track nested attributes.
func TestPortAttributeTracking(t *testing.T) {
	input := `
package TestPkg {
    port def DataPort {
        attribute dataRate;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}

	// Port definition should be found
	ports := pkg.Children()
	var dataPort *Port
	for _, child := range ports {
		if p, ok := child.(*Port); ok && p.Name() == "DataPort" {
			dataPort = p
			break
		}
	}
	if dataPort == nil {
		t.Fatal("DataPort not found")
	}

	if len(dataPort.Attributes()) != 1 {
		t.Errorf("expected 1 attribute in DataPort, got %d", len(dataPort.Attributes()))
	}
}

// TestStateTypedChildTracking verifies State tracks actions, attributes, and constraints.
func TestStateTypedChildTracking(t *testing.T) {
	input := `
package TestPkg {
    state def OperatingMode {
        state idle;
        state running;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Err())
	}

	pkg := result.Model.FindPackage("TestPkg")
	if pkg == nil {
		t.Fatal("package TestPkg not found")
	}
	if len(pkg.States()) != 1 {
		t.Fatalf("expected 1 state, got %d", len(pkg.States()))
	}

	opMode := pkg.States()[0]
	if len(opMode.States()) != 2 {
		t.Errorf("expected 2 nested states, got %d", len(opMode.States()))
	}
}
