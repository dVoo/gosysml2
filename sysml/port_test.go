package sysml

import (
	"testing"
)

func TestNewConjugatedPort(t *testing.T) {
	loc := Location{Line: 10, Column: 5}

	// Test creating conjugated port
	conj := NewConjugatedPort("~TestPort", loc)
	if conj == nil {
		t.Fatal("NewConjugatedPort returned nil")
	}
	if conj.Name() != "~TestPort" {
		t.Errorf("Expected name '~TestPort', got '%s'", conj.Name())
	}
	if conj.Kind() != KindConjugatedPort {
		t.Errorf("Expected kind KindConjugatedPort, got %v", conj.Kind())
	}
	if conj.Location() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, conj.Location())
	}
}

func TestConjugatedPortIsDefinition(t *testing.T) {
	loc := Location{}
	conj := NewConjugatedPort("~TestPort", loc)

	// Test that ConjugatedPort implements Definition interface
	var _ Definition = conj
}

func TestConjugatedPortOriginalPort(t *testing.T) {
	loc := Location{}
	conj := NewConjugatedPort("~TestPort", loc)

	// Initially no original port
	if conj.GetOriginalPort() != nil {
		t.Error("Expected nil original port for new conjugated port")
	}

	// Create an original port
	original := NewPort("TestPort", loc, true)
	conj.OriginalPort.Resolve(original)

	// Now GetOriginalPort should return the original
	if conj.GetOriginalPort() != original {
		t.Error("GetOriginalPort should return resolved port")
	}
}

func TestConjugatedPortSetUnresolvedOriginal(t *testing.T) {
	loc := Location{}
	conj := NewConjugatedPort("~TestPort", loc)

	// Initially OriginalPort should be unresolved
	if conj.OriginalPort.IsResolved() {
		t.Error("Expected OriginalPort to be initially unresolved")
	}

	// Set unresolved original port
	conj.SetUnresolvedOriginalPort("MyPort")
	// After setting, it should still be unresolved (just stored the name)
	if conj.OriginalPort.IsResolved() {
		t.Error("Expected OriginalPort to remain unresolved until resolution")
	}
}

func TestConjugatedPortEffectiveName(t *testing.T) {
	loc := Location{}

	// Test with unresolved original port
	conj := NewConjugatedPort("~TestPort", loc)
	effectiveName := conj.EffectiveName()
	if effectiveName != "~TestPort" {
		t.Errorf("Expected effective name '~TestPort' (unresolved), got '%s'", effectiveName)
	}

	// Test with resolved original port
	original := NewPort("MyPort", loc, true)
	conj.OriginalPort.Resolve(original)
	effectiveName = conj.EffectiveName()
	if effectiveName != "~MyPort" {
		t.Errorf("Expected effective name '~MyPort' (resolved), got '%s'", effectiveName)
	}
}

func TestConjugatedPortExists(t *testing.T) {
	input := `package TestPkg {
		port def MyPort;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Err())
	}

	// Find ports
	ports := FindAll[*Port](result.Model)
	if len(ports) != 1 {
		t.Errorf("Expected 1 port, got %d", len(ports))
		return
	}

	port := ports[0]

	// Check that conjugated port exists
	if port.ConjugatedPort == nil {
		t.Error("Expected ConjugatedPort to be auto-created with PortDefinition")
		return
	}

	// Check conjugated port name
	if port.ConjugatedPort.Name() != "~MyPort" {
		t.Errorf("Expected conjugated port name '~MyPort', got '%s'", port.ConjugatedPort.Name())
	}

	// Check that conjugated port is a child of the port
	foundAsChild := false
	for _, child := range port.Children() {
		if child == port.ConjugatedPort {
			foundAsChild = true
			break
		}
	}
	if !foundAsChild {
		t.Error("Expected ConjugatedPort to be a child of the Port")
	}
}

func TestConjugatedPortOriginalRef(t *testing.T) {
	input := `package TestPkg {
		port def MyPort;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Err())
	}

	// Build index and resolve references
	result.Model.BuildIndex()
	result.Model.ResolveReferences()

	ports := FindAll[*Port](result.Model)
	if len(ports) != 1 {
		t.Fatalf("Expected 1 port, got %d", len(ports))
	}

	port := ports[0]
	if port.ConjugatedPort == nil {
		t.Fatal("Expected ConjugatedPort to exist")
	}

	// Check that OriginalPort reference is resolved
	if !port.ConjugatedPort.OriginalPort.IsResolved() {
		t.Error("Expected OriginalPort to be resolved")
	}

	// Check that it points to the correct port
	if port.ConjugatedPort.GetOriginalPort() != port {
		t.Error("Expected OriginalPort to reference the parent Port")
	}
}

func TestConjugatedPortVisitor(t *testing.T) {
	input := `package TestPkg {
		port def MyPort;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Err())
	}

	// Test Counter visitor
	counter := NewCounter()
	Visit(result.Model, counter)

	if counter.Counts[KindConjugatedPort] != 1 {
		t.Errorf("Expected 1 conjugated port count, got %d", counter.Counts[KindConjugatedPort])
	}
}

func TestFindConjugatedPorts(t *testing.T) {
	input := `package TestPkg {
		port def Port1;
		port def Port2;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Err())
	}

	conjugatedPorts := FindAll[*ConjugatedPort](result.Model)
	if len(conjugatedPorts) != 2 {
		t.Errorf("Expected 2 conjugated ports, got %d", len(conjugatedPorts))
	}
}

func TestConjugatedPortKindString(t *testing.T) {
	if KindConjugatedPort.String() != "conjugated port" {
		t.Errorf("Expected 'conjugated port', got '%s'", KindConjugatedPort.String())
	}
}

// Integration test: Parse port with conjugation syntax
func TestParsePortConjugation(t *testing.T) {
	input := `package TestPkg {
		port def MyPort;
		
		part def MyPart {
			port p : ~MyPort;
		}
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Err())
	}

	// Find the port definition
	ports := FindAll[*Port](result.Model)

	// We should have MyPort (definition) and p (usage)
	var myPortDef *Port

	for _, port := range ports {
		if port.Name() == "MyPort" && port.IsDefinition {
			myPortDef = port
		}
	}

	if myPortDef == nil {
		t.Fatal("Could not find MyPort definition")
	}

	if myPortDef.ConjugatedPort == nil {
		t.Fatal("Expected ConjugatedPort to be auto-created")
	}

	// The conjugated port should exist with name ~MyPort
	if myPortDef.ConjugatedPort.Name() != "~MyPort" {
		t.Errorf("Expected conjugated port name '~MyPort', got '%s'", myPortDef.ConjugatedPort.Name())
	}
}

// Test that conjugated port is only created for definitions, not usages
func TestConjugatedPortOnlyForDefinitions(t *testing.T) {
	input := `package TestPkg {
		port def MyPort;
		port myPortUsage : MyPort;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Err())
	}

	ports := FindAll[*Port](result.Model)

	var defCount, usageCount int
	for _, port := range ports {
		if port.IsDefinition {
			defCount++
			// Definitions should have conjugated ports
			if port.ConjugatedPort == nil {
				t.Errorf("Port definition '%s' should have ConjugatedPort", port.Name())
			}
		} else {
			usageCount++
			// Usages should not have conjugated ports
			if port.ConjugatedPort != nil {
				t.Errorf("Port usage '%s' should not have ConjugatedPort", port.Name())
			}
		}
	}

	if defCount != 1 {
		t.Errorf("Expected 1 port definition, got %d", defCount)
	}
	if usageCount != 1 {
		t.Errorf("Expected 1 port usage, got %d", usageCount)
	}
}

// Test multiple port definitions with conjugated ports
func TestMultipleConjugatedPorts(t *testing.T) {
	input := `package TestPkg {
		port def InputPort;
		port def OutputPort;
		port def BidirectionalPort;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Err())
	}

	conjugatedPorts := FindAll[*ConjugatedPort](result.Model)
	if len(conjugatedPorts) != 3 {
		t.Fatalf("Expected 3 conjugated ports, got %d", len(conjugatedPorts))
	}

	// Check each conjugated port has correct name
	expectedNames := map[string]bool{
		"~InputPort":         false,
		"~OutputPort":        false,
		"~BidirectionalPort": false,
	}

	for _, conj := range conjugatedPorts {
		if _, ok := expectedNames[conj.Name()]; ok {
			expectedNames[conj.Name()] = true
		} else {
			t.Errorf("Unexpected conjugated port name: %s", conj.Name())
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Missing conjugated port: %s", name)
		}
	}
}

// Test visitor pattern with BaseVisitor
func TestConjugatedPortBaseVisitor(t *testing.T) {
	input := `package TestPkg {
		port def MyPort;
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("Parsing failed: %v", result.Err())
	}

	// Test that BaseVisitor handles ConjugatedPort without panic
	baseVisitor := BaseVisitor{}
	Visit(result.Model, baseVisitor)
	// If we get here without panic, BaseVisitor works correctly
}
