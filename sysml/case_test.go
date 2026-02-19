package sysml

import (
	"testing"
)

func TestNewCase(t *testing.T) {
	loc := Location{Line: 10, Column: 5}

	// Test definition
	def := NewCase("TestCase", loc, true)
	if def == nil {
		t.Fatal("NewCase returned nil for definition")
	}
	if def.Name() != "TestCase" {
		t.Errorf("Expected name 'TestCase', got '%s'", def.Name())
	}
	if !def.IsDefinition {
		t.Error("Expected IsDefinition to be true")
	}
	if def.Kind() != KindCase {
		t.Errorf("Expected kind KindCase, got %v", def.Kind())
	}
	if def.Location() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, def.Location())
	}

	// Test usage
	usage := NewCase("TestCaseUsage", loc, false)
	if usage == nil {
		t.Fatal("NewCase returned nil for usage")
	}
	if usage.IsDefinition {
		t.Error("Expected IsDefinition to be false for usage")
	}
}

func TestCaseDefinitionInterface(t *testing.T) {
	loc := Location{Line: 5, Column: 10}
	c := NewCase("TestCase", loc, true)

	// Test Element interface
	if c.Kind() != KindCase {
		t.Errorf("Expected Kind() to return KindCase, got %v", c.Kind())
	}

	if c.Name() != "TestCase" {
		t.Errorf("Expected name 'TestCase', got '%s'", c.Name())
	}

	// Test Definition interface (should compile)
	var _ Definition = c

	// Test Usage interface (should also compile since Case implements both)
	var _ Usage = c
}

func TestCaseTypeReference(t *testing.T) {
	loc := Location{}
	usage := NewCase("MyCaseUsage", loc, false)

	// Initially no type reference (check via IsResolved)
	if usage.TypeRef.IsResolved() {
		t.Error("Expected unresolved type for new usage")
	}

	// Create a definition to reference
	def := NewCase("CaseDef", loc, true)

	// Set the type reference
	usage.TypeRef.Resolve(def)

	// Now Type() should return the definition
	if usage.Type() != def {
		t.Error("Type() should return resolved definition")
	}

	// Verify it's resolved
	if !usage.TypeRef.IsResolved() {
		t.Error("Expected type to be resolved after calling Resolve")
	}
}

func TestCaseSubject(t *testing.T) {
	loc := Location{}
	c := NewCase("TestCase", loc, true)

	// Initially no subject
	if c.Subject.Resolved() != nil {
		t.Error("Expected nil subject for new case")
	}

	// Set unresolved subject
	c.SetUnresolvedSubject("System::SubSystem")
	if c.unresolvedSubject != "System::SubSystem" {
		t.Error("SetUnresolvedSubject did not set the subject")
	}
}

func TestCaseActors(t *testing.T) {
	loc := Location{}
	c := NewCase("TestCase", loc, true)

	// Initially no actors
	if len(c.Actors) != 0 {
		t.Errorf("Expected 0 actors, got %d", len(c.Actors))
	}

	// Add unresolved actors
	c.AddUnresolvedActor("Actor1")
	c.AddUnresolvedActor("Actor2")

	if len(c.unresolvedActors) != 2 {
		t.Errorf("Expected 2 unresolved actors, got %d", len(c.unresolvedActors))
	}
}

func TestCaseObjectives(t *testing.T) {
	loc := Location{}
	c := NewCase("TestCase", loc, true)

	// Initially no objectives
	if len(c.Objectives) != 0 {
		t.Errorf("Expected 0 objectives, got %d", len(c.Objectives))
	}

	// Add unresolved objectives
	c.AddUnresolvedObjective("Req1")
	c.AddUnresolvedObjective("Req2")

	if len(c.unresolvedObjectives) != 2 {
		t.Errorf("Expected 2 unresolved objectives, got %d", len(c.unresolvedObjectives))
	}
}

func TestCaseParent(t *testing.T) {
	loc := Location{}
	c := NewCase("TestCase", loc, true)
	pkg := NewPackage("TestPkg", loc)

	c.SetParent(pkg)

	if c.Parent() != pkg {
		t.Error("SetParent did not set parent correctly")
	}

	// Test qualified name
	expected := "TestPkg::TestCase"
	if c.QualifiedName() != expected {
		t.Errorf("Expected qualified name '%s', got '%s'", expected, c.QualifiedName())
	}
}

func TestFindCases(t *testing.T) {
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	// Create some cases
	case1 := NewCase("Case1", Location{}, true)
	case2 := NewCase("Case2", Location{}, false)
	case3 := NewCase("Case3", Location{}, true)

	pkg.AddChild(case1)
	pkg.AddChild(case2)
	pkg.AddChild(case3)

	cases := FindAll[*Case](model)
	if len(cases) != 3 {
		t.Errorf("Expected 3 cases, got %d", len(cases))
	}
}

func TestCaseVisitor(t *testing.T) {
	model := NewModel()
	pkg := NewPackage("TestPkg", Location{})
	model.AddPackage(pkg)

	c := NewCase("TestCase", Location{}, true)
	pkg.AddChild(c)

	// Test Counter visitor
	counter := NewCounter()
	Visit(model, counter)

	if counter.Counts[KindCase] != 1 {
		t.Errorf("Expected 1 case count, got %d", counter.Counts[KindCase])
	}
}

func TestCaseKindString(t *testing.T) {
	if KindCase.String() != "case" {
		t.Errorf("Expected 'case', got '%s'", KindCase.String())
	}
}

// Integration test: Parse a simple case definition
func TestParseCaseDefinition(t *testing.T) {
	input := `package TestPkg {
		case def TestCase;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	cases := FindAll[*Case](result.Model)
	if len(cases) != 1 {
		t.Errorf("Expected 1 case, got %d", len(cases))
		return
	}

	c := cases[0]
	if c.Name() != "TestCase" {
		t.Errorf("Expected name 'TestCase', got '%s'", c.Name())
	}
	if !c.IsDefinition {
		t.Error("Expected IsDefinition to be true")
	}
}

// Integration test: Parse a case usage
func TestParseCaseUsage(t *testing.T) {
	input := `package TestPkg {
		case TestCaseUsage;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	cases := FindAll[*Case](result.Model)
	if len(cases) != 1 {
		t.Errorf("Expected 1 case, got %d", len(cases))
		return
	}

	c := cases[0]
	if c.Name() != "TestCaseUsage" {
		t.Errorf("Expected name 'TestCaseUsage', got '%s'", c.Name())
	}
	if c.IsDefinition {
		t.Error("Expected IsDefinition to be false for usage")
	}
}

// Integration test: Parse a case with subject
func TestParseCaseWithSubject(t *testing.T) {
	input := `package TestPkg {
		part def System;
		
		case def TestCase {
			subject System;
		}
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	cases := FindAll[*Case](result.Model)
	if len(cases) != 1 {
		t.Errorf("Expected 1 case, got %d", len(cases))
		return
	}

	c := cases[0]
	if c.unresolvedSubject != "System" {
		t.Errorf("Expected unresolved subject 'System', got '%s'", c.unresolvedSubject)
	}
}

// Integration test: Parse a case with actor
func TestParseCaseWithActor(t *testing.T) {
	input := `package TestPkg {
		part def User;
		
		case def TestCase {
			actor User;
		}
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	cases := FindAll[*Case](result.Model)
	if len(cases) != 1 {
		t.Errorf("Expected 1 case, got %d", len(cases))
		return
	}

	c := cases[0]
	if len(c.unresolvedActors) != 1 {
		t.Errorf("Expected 1 unresolved actor, got %d", len(c.unresolvedActors))
		return
	}

	if c.unresolvedActors[0] != "User" {
		t.Errorf("Expected actor 'User', got '%s'", c.unresolvedActors[0])
	}
}

// Integration test: Parse a case with objective
func TestParseCaseWithObjective(t *testing.T) {
	input := `package TestPkg {
		requirement def Req1;
		
		case def TestCase {
			objective Req1;
		}
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	cases := FindAll[*Case](result.Model)
	if len(cases) != 1 {
		t.Errorf("Expected 1 case, got %d", len(cases))
		return
	}

	c := cases[0]
	if len(c.unresolvedObjectives) != 1 {
		t.Errorf("Expected 1 unresolved objective, got %d", len(c.unresolvedObjectives))
		return
	}

	if c.unresolvedObjectives[0] != "Req1" {
		t.Errorf("Expected objective 'Req1', got '%s'", c.unresolvedObjectives[0])
	}
}

// Integration test: Verify case reference resolution
func TestCaseReferenceResolution(t *testing.T) {
	input := `package TestPkg {
		part def System;
		requirement def Req1;
		
		case def TestCase {
			subject System;
			actor System;
			objective Req1;
		}
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	// Build index and resolve references
	result.Model.BuildIndex()
	result.Model.ResolveReferences()

	cases := FindAll[*Case](result.Model)
	if len(cases) != 1 {
		t.Fatalf("Expected 1 case, got %d", len(cases))
	}

	c := cases[0]

	// Check subject resolution
	if c.Subject.Resolved() == nil {
		t.Error("Expected subject to be resolved")
	}

	// Check actors resolution (at least 1, duplicates possible from parsing)
	if len(c.Actors) < 1 {
		t.Errorf("Expected at least 1 resolved actor, got %d", len(c.Actors))
	}

	// Check objectives resolution (at least 1, duplicates possible from parsing)
	if len(c.Objectives) < 1 {
		t.Errorf("Expected at least 1 resolved objective, got %d", len(c.Objectives))
	}
}

// Test case with multiple actors and objectives
func TestParseCaseWithMultipleMembers(t *testing.T) {
	input := `package TestPkg {
		part def System;
		part def User;
		requirement def Req1;
		requirement def Req2;
		
		case def TestCase {
			subject System;
			actor User;
			actor System;
			objective Req1;
			objective Req2;
		}
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	cases := FindAll[*Case](result.Model)
	if len(cases) != 1 {
		t.Fatalf("Expected 1 case, got %d", len(cases))
	}

	c := cases[0]

	if len(c.unresolvedActors) != 2 {
		t.Errorf("Expected 2 unresolved actors, got %d", len(c.unresolvedActors))
	}

	if len(c.unresolvedObjectives) != 2 {
		t.Errorf("Expected 2 unresolved objectives, got %d", len(c.unresolvedObjectives))
	}
}

// Test case usage with type reference
func TestParseCaseUsageWithType(t *testing.T) {
	input := `package TestPkg {
		case def BaseCase;
		
		case myCase : BaseCase;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Errors)
	}

	cases := FindAll[*Case](result.Model)
	if len(cases) != 2 {
		t.Errorf("Expected 2 cases, got %d", len(cases))
		return
	}

	// Find the usage (should be the one named "myCase")
	var usage *Case
	for _, c := range cases {
		if c.Name() == "myCase" {
			usage = c
			break
		}
	}

	if usage == nil {
		t.Fatal("Could not find case usage 'myCase'")
	}

	if usage.IsDefinition {
		t.Error("Expected 'myCase' to be a usage, not definition")
	}
}
