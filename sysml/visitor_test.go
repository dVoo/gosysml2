package sysml

import (
	"testing"
)

func TestFindRequirements(t *testing.T) {
	input := `
		package TestPackage {
			requirement def SafetyRequirement {
				doc /* Safety requirements must be satisfied */
			}
			requirement def PerformanceRequirement;
			requirement testReq : SafetyRequirement;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	reqs := FindAll[*Requirement](result.Model)
	if len(reqs) != 3 {
		t.Errorf("Expected 3 requirements, got %d", len(reqs))
	}

	// Check that we found the expected requirements
	names := make(map[string]bool)
	for _, req := range reqs {
		names[req.Name()] = true
	}

	if !names["SafetyRequirement"] {
		t.Error("Expected to find SafetyRequirement")
	}
	if !names["PerformanceRequirement"] {
		t.Error("Expected to find PerformanceRequirement")
	}
	if !names["testReq"] {
		t.Error("Expected to find testReq")
	}
}

func TestFindVerifications(t *testing.T) {
	input := `
		package TestPackage {
			verification def TestVerification {
				doc /* Test the system */
			}
			verification def InspectionVerification;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	vers := FindAll[*Verification](result.Model)
	if len(vers) != 2 {
		t.Errorf("Expected 2 verifications, got %d", len(vers))
	}

	names := make(map[string]bool)
	for _, ver := range vers {
		names[ver.Name()] = true
	}

	if !names["TestVerification"] {
		t.Error("Expected to find TestVerification")
	}
	if !names["InspectionVerification"] {
		t.Error("Expected to find InspectionVerification")
	}
}

func TestFindConcerns(t *testing.T) {
	input := `
		package TestPackage {
			concern def SecurityConcern;
			concern def SafetyConcern;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	concerns := FindAll[*Concern](result.Model)
	if len(concerns) != 2 {
		t.Errorf("Expected 2 concerns, got %d", len(concerns))
	}
}

func TestFindUseCases(t *testing.T) {
	input := `
		package TestPackage {
			use case def DriveVehicle;
			use case def ParkVehicle;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	useCases := FindAll[*UseCase](result.Model)
	if len(useCases) != 2 {
		t.Errorf("Expected 2 use cases, got %d", len(useCases))
	}
}

func TestFindAnalysisCases(t *testing.T) {
	input := `
		package TestPackage {
			analysis def ThermalAnalysis;
			analysis def StructuralAnalysis;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	analyses := FindAll[*AnalysisCase](result.Model)
	if len(analyses) != 2 {
		t.Errorf("Expected 2 analysis cases, got %d", len(analyses))
	}
}

func TestCounterWithVerification(t *testing.T) {
	input := `
		package TestPackage {
			part def Vehicle;
			part def Engine;
			requirement def SafetyReq;
			verification def TestVer;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	counter := NewCounter()
	Visit(result.Model, counter)

	if counter.Counts[KindPackage] != 1 {
		t.Errorf("Expected 1 package, got %d", counter.Counts[KindPackage])
	}
	if counter.Counts[KindPart] != 2 {
		t.Errorf("Expected 2 parts, got %d", counter.Counts[KindPart])
	}
	if counter.Counts[KindRequirement] != 1 {
		t.Errorf("Expected 1 requirement, got %d", counter.Counts[KindRequirement])
	}
	if counter.Counts[KindVerification] != 1 {
		t.Errorf("Expected 1 verification, got %d", counter.Counts[KindVerification])
	}

	total := counter.Total()
	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
}

func TestWalkWithDepth(t *testing.T) {
	input := `
		package OuterPackage {
			package InnerPackage {
				part def InnerPart;
			}
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	depths := make(map[string]int)
	Walk(result.Model, func(elem Element, depth int) bool {
		depths[elem.Name()] = depth
		return true
	})

	if depths["OuterPackage"] != 0 {
		t.Errorf("Expected OuterPackage at depth 0, got %d", depths["OuterPackage"])
	}
	if depths["InnerPackage"] != 1 {
		t.Errorf("Expected InnerPackage at depth 1, got %d", depths["InnerPackage"])
	}
	if depths["InnerPart"] != 2 {
		t.Errorf("Expected InnerPart at depth 2, got %d", depths["InnerPart"])
	}
}

func TestFilter(t *testing.T) {
	input := `
		package TestPackage {
			part def PartA;
			part def PartB;
			requirement def ReqA;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	parts := Filter(result.Model, func(elem Element) bool {
		return elem.Kind() == KindPart
	})

	if len(parts) != 2 {
		t.Errorf("Expected 2 parts, got %d", len(parts))
	}
}

func TestFindByKindWithRequirements(t *testing.T) {
	input := `
		package TestPackage {
			part def PartA;
			part def PartB;
			requirement def ReqA;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	parts := FindByKind(result.Model, KindPart)
	if len(parts) != 2 {
		t.Errorf("Expected 2 parts, got %d", len(parts))
	}

	reqs := FindByKind(result.Model, KindRequirement)
	if len(reqs) != 1 {
		t.Errorf("Expected 1 requirement, got %d", len(reqs))
	}
}

func TestFindByName(t *testing.T) {
	input := `
		package TestPackage {
			part def Vehicle;
			requirement def Vehicle;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	vehicles := FindByName(result.Model, "Vehicle")
	if len(vehicles) != 2 {
		t.Errorf("Expected 2 elements named Vehicle, got %d", len(vehicles))
	}
}

func TestRequirementIsDefinition(t *testing.T) {
	input := `
		package TestPackage {
			requirement def SafetyReq;
			requirement instanceReq : SafetyReq;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	reqs := FindAll[*Requirement](result.Model)
	if len(reqs) != 2 {
		t.Fatalf("Expected 2 requirements, got %d", len(reqs))
	}

	defCount := 0
	usageCount := 0
	for _, req := range reqs {
		if req.IsDefinition {
			defCount++
		} else {
			usageCount++
		}
	}

	if defCount != 1 {
		t.Errorf("Expected 1 definition, got %d", defCount)
	}
	if usageCount != 1 {
		t.Errorf("Expected 1 usage, got %d", usageCount)
	}
}

func TestQualifiedName(t *testing.T) {
	input := `
		package Outer {
			package Inner {
				part def MyPart;
			}
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	// Find by qualified name
	elem := result.Model.FindByQualifiedName("Outer::Inner::MyPart")
	if elem == nil {
		t.Fatal("Expected to find Outer::Inner::MyPart")
	}

	if elem.Name() != "MyPart" {
		t.Errorf("Expected name 'MyPart', got '%s'", elem.Name())
	}

	if elem.QualifiedName() != "Outer::Inner::MyPart" {
		t.Errorf("Expected QN 'Outer::Inner::MyPart', got '%s'", elem.QualifiedName())
	}
}

func TestTypedChildAccessors(t *testing.T) {
	input := `
		package TestPackage {
			part def PartA;
			part def PartB;
			requirement def ReqA;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	pkg := result.Model.FindPackage("TestPackage")
	if pkg == nil {
		t.Fatal("Expected to find TestPackage")
	}

	// Use typed accessors
	parts := pkg.Parts()
	if len(parts) != 2 {
		t.Errorf("Expected 2 parts via typed accessor, got %d", len(parts))
	}

	reqs := pkg.Requirements()
	if len(reqs) != 1 {
		t.Errorf("Expected 1 requirement via typed accessor, got %d", len(reqs))
	}
}

func TestDefinitionUsageInterface(t *testing.T) {
	input := `
		package TestPackage {
			part def VehicleDef;
			part vehicle : VehicleDef;
		}
	`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Failed to parse: %v", result.Errors)
	}

	// Note: In Go's structural typing, both Part definitions and usages
	// implement both Definition and Usage interfaces.
	// The IsDefinition field is what actually distinguishes them.

	parts := FindAll[*Part](result.Model)
	if len(parts) != 2 {
		t.Fatalf("Expected 2 parts, got %d", len(parts))
	}

	actualDefs := 0
	actualUsages := 0
	for _, p := range parts {
		if p.IsDefinition {
			actualDefs++
		} else {
			actualUsages++
		}
	}

	if actualDefs != 1 {
		t.Errorf("Expected 1 part with IsDefinition=true, got %d", actualDefs)
	}
	if actualUsages != 1 {
		t.Errorf("Expected 1 part with IsDefinition=false, got %d", actualUsages)
	}

	// Verify FindDefinitions and FindUsages work (they find all types that implement the interface)
	definitions := FindAll[Definition](result.Model)
	usages := FindAll[Usage](result.Model)

	// Both should find the same Part elements (since Part implements both interfaces)
	if len(definitions) < 2 {
		t.Errorf("Expected at least 2 elements implementing Definition, got %d", len(definitions))
	}
	if len(usages) < 2 {
		t.Errorf("Expected at least 2 elements implementing Usage, got %d", len(usages))
	}
}

func TestRefType(t *testing.T) {
	// Test the Ref generic type
	ref := NewRef[Element]("SomeElement")

	if ref.Name() != "SomeElement" {
		t.Errorf("Expected name 'SomeElement', got '%s'", ref.Name())
	}

	if ref.IsResolved() {
		t.Error("Expected ref to be unresolved initially")
	}

	if ref.Resolved() != nil {
		t.Error("Expected Resolved() to return nil for unresolved ref")
	}
}
