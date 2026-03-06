package sysml

import (
	"testing"
)

// TestThenSuccessionInDefinitionBody verifies that `then X;` target succession
// parses correctly inside a part definition body (Batmobile pattern).
func TestThenSuccessionInDefinitionBody(t *testing.T) {
	input := `
part def Batmobile {
    timeslice batmanDriving {
        attribute speed : Real;
    }
    then charging;
    timeslice charging {
        attribute speed : Real;
    }
}
`
	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("expected parse to succeed, got: %s", result.Err())
	}
}

func TestNewOccurrence(t *testing.T) {
	loc := Location{Line: 10, Column: 5}
	occ := NewOccurrence("MyOccurrence", loc, true, false)

	if occ == nil {
		t.Fatal("NewOccurrence returned nil")
	}

	if occ.Name() != "MyOccurrence" {
		t.Errorf("Expected name 'MyOccurrence', got '%s'", occ.Name())
	}

	if occ.Kind() != KindOccurrence {
		t.Errorf("Expected kind %v, got %v", KindOccurrence, occ.Kind())
	}

	if !occ.IsDefinition {
		t.Error("Expected IsDefinition to be true")
	}

	if occ.IsIndividual {
		t.Error("Expected IsIndividual to be false")
	}

	if occ.Location() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, occ.Location())
	}

	if occ.Parent() != nil {
		t.Error("Expected nil parent for new occurrence")
	}
}

func TestNewOccurrenceUsage(t *testing.T) {
	loc := Location{Line: 20, Column: 10}
	occ := NewOccurrence("MyOccurrence", loc, false, false)

	if occ.IsDefinition {
		t.Error("Expected IsDefinition to be false for usage")
	}

	if occ.IsIndividual {
		t.Error("Expected IsIndividual to be false")
	}
}

func TestNewIndividualOccurrence(t *testing.T) {
	loc := Location{Line: 30, Column: 15}
	occ := NewOccurrence("IndividualOcc", loc, true, true)

	if !occ.IsDefinition {
		t.Error("Expected IsDefinition to be true for individual definition")
	}

	if !occ.IsIndividual {
		t.Error("Expected IsIndividual to be true")
	}
}

func TestNewEventOccurrence(t *testing.T) {
	loc := Location{Line: 40, Column: 20}
	occ := NewEventOccurrence("EventOcc", loc)

	if occ.IsDefinition {
		t.Error("Expected IsDefinition to be false for event occurrence")
	}

	if occ.IsIndividual {
		t.Error("Expected IsIndividual to be false for event occurrence")
	}

	if !occ.IsEvent {
		t.Error("Expected IsEvent to be true")
	}
}

func TestNewTimeSlice(t *testing.T) {
	loc := Location{Line: 50, Column: 25}
	occ := NewTimeSlice("TimeSlice1", loc)

	if occ.IsDefinition {
		t.Error("Expected IsDefinition to be false for time slice")
	}

	if !occ.IsTimeSlice {
		t.Error("Expected IsTimeSlice to be true")
	}
}

func TestNewSnapshot(t *testing.T) {
	loc := Location{Line: 60, Column: 30}
	occ := NewSnapshot("Snapshot1", loc)

	if occ.IsDefinition {
		t.Error("Expected IsDefinition to be false for snapshot")
	}

	if !occ.IsSnapshot {
		t.Error("Expected IsSnapshot to be true")
	}

	if occ.PortionKind_ != PortionKindSnapshot {
		t.Errorf("Expected PortionKindSnapshot, got %v", occ.PortionKind_)
	}
}

func TestNewLifeStep(t *testing.T) {
	loc := Location{Line: 70, Column: 35}

	// Test start life step
	startOcc := NewLifeStep("StartStep", loc, LifeStepStart)
	if startOcc.LifeStep_ != LifeStepStart {
		t.Errorf("Expected LifeStepStart, got %v", startOcc.LifeStep_)
	}

	// Test end life step
	endOcc := NewLifeStep("EndStep", loc, LifeStepEnd)
	if endOcc.LifeStep_ != LifeStepEnd {
		t.Errorf("Expected LifeStepEnd, got %v", endOcc.LifeStep_)
	}
}

func TestPortionKindString(t *testing.T) {
	testCases := []struct {
		kind     PortionKind
		expected string
	}{
		{PortionKindNone, ""},
		{PortionKindSnapshot, "snapshot"},
		{PortionKindTimeslice, "timeslice"},
	}

	for _, tc := range testCases {
		result := tc.kind.String()
		if result != tc.expected {
			t.Errorf("PortionKind(%v).String() = '%s', expected '%s'", tc.kind, result, tc.expected)
		}
	}
}

func TestLifeStepString(t *testing.T) {
	testCases := []struct {
		step     LifeStep
		expected string
	}{
		{LifeStepNone, ""},
		{LifeStepStart, "start"},
		{LifeStepEnd, "end"},
	}

	for _, tc := range testCases {
		result := tc.step.String()
		if result != tc.expected {
			t.Errorf("LifeStep(%v).String() = '%s', expected '%s'", tc.step, result, tc.expected)
		}
	}
}

func TestOccurrenceSetPortionKind(t *testing.T) {
	loc := Location{}
	occ := NewOccurrence("Test", loc, false, false)

	// Initially no portion kind
	if occ.PortionKind_ != PortionKindNone {
		t.Errorf("Expected PortionKindNone initially, got %v", occ.PortionKind_)
	}

	// Set portion kind
	occ.SetPortionKind(PortionKindTimeslice)
	if occ.PortionKind_ != PortionKindTimeslice {
		t.Errorf("Expected PortionKindTimeslice, got %v", occ.PortionKind_)
	}
}

func TestOccurrenceSetLifeStep(t *testing.T) {
	loc := Location{}
	occ := NewOccurrence("Test", loc, false, false)

	// Initially no life step
	if occ.LifeStep_ != LifeStepNone {
		t.Errorf("Expected LifeStepNone initially, got %v", occ.LifeStep_)
	}

	// Set life step
	occ.SetLifeStep(LifeStepStart)
	if occ.LifeStep_ != LifeStepStart {
		t.Errorf("Expected LifeStepStart, got %v", occ.LifeStep_)
	}
}

func TestOccurrenceParent(t *testing.T) {
	occ := NewOccurrence("TestOcc", Location{}, false, false)
	pkg := NewPackage("TestPkg", Location{})

	occ.SetParent(pkg)

	if occ.Parent() != pkg {
		t.Error("SetParent did not set parent correctly for occurrence")
	}
}

func TestOccurrenceAddChild(t *testing.T) {
	parent := NewOccurrence("Parent", Location{}, true, false)
	child := NewOccurrence("Child", Location{}, false, false)

	parent.AddChild(child)

	children := parent.Children()
	if len(children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(children))
	}

	if children[0] != child {
		t.Error("Child not added correctly to occurrence")
	}

	// Check Occurrences() method
	occurrences := parent.Occurrences()
	if len(occurrences) != 1 {
		t.Errorf("Expected 1 occurrence, got %d", len(occurrences))
	}
}

func TestOccurrenceControlNodeMethods(t *testing.T) {
	loc := Location{}
	occ := NewOccurrence("Test", loc, false, false)

	// Occurrences should not be control nodes
	if occ.IsFork() {
		t.Error("Occurrence should not be a fork")
	}
	if occ.IsJoin() {
		t.Error("Occurrence should not be a join")
	}
	if occ.IsMerge() {
		t.Error("Occurrence should not be a merge")
	}
	if occ.IsDecision() {
		t.Error("Occurrence should not be a decision")
	}
}

func TestModelAddOccurrence(t *testing.T) {
	model := NewModel()
	occ := NewOccurrence("TestOcc", Location{Line: 5, Column: 10}, true, false)

	model.AddOccurrence(occ)

	if len(model.Occurrences()) != 1 {
		t.Errorf("Expected 1 occurrence in model, got %d", len(model.Occurrences()))
	}

	if len(model.Elements) != 1 {
		t.Errorf("Expected 1 element in model, got %d", len(model.Elements))
	}

	if model.Occurrences()[0] != occ {
		t.Error("Occurrence not added correctly to model")
	}
}

func TestOccurrenceQualifiedName(t *testing.T) {
	pkg := NewPackage("TestPkg", Location{})
	occ := NewOccurrence("MyOcc", Location{}, true, false)

	// Initially no parent, so qualified name is just the name
	if occ.QualifiedName() != "MyOcc" {
		t.Errorf("Expected qualified name 'MyOcc', got '%s'", occ.QualifiedName())
	}

	// Add to package
	pkg.AddChild(occ)

	// Now qualified name should include package
	expected := "TestPkg::MyOcc"
	qn := occ.QualifiedName()
	if qn != expected {
		t.Errorf("Expected qualified name '%s', got '%s'", expected, qn)
	}
}
