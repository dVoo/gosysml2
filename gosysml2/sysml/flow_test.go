package sysml

import (
	"testing"
)

func TestNewFlow(t *testing.T) {
	loc := Location{Line: 10, Column: 5}
	flow := NewFlow("DataFlow", loc, true)

	if flow == nil {
		t.Fatal("NewFlow returned nil")
	}

	if flow.GetName() != "DataFlow" {
		t.Errorf("Expected name 'DataFlow', got '%s'", flow.GetName())
	}

	if flow.GetKind() != KindFlow {
		t.Errorf("Expected kind %v, got %v", KindFlow, flow.GetKind())
	}

	if !flow.IsDefinition {
		t.Error("Expected IsDefinition to be true")
	}

	if flow.GetLocation() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, flow.GetLocation())
	}

	if flow.GetParent() != nil {
		t.Error("Expected nil parent for new flow")
	}
}

func TestNewFlowUsage(t *testing.T) {
	loc := Location{Line: 20, Column: 10}
	flow := NewFlow("myFlow", loc, false)

	if flow.IsDefinition {
		t.Error("Expected IsDefinition to be false for usage")
	}
}

func TestNewFlowEnd(t *testing.T) {
	loc := Location{Line: 15, Column: 8}
	flowEnd := NewFlowEnd(loc)

	if flowEnd == nil {
		t.Fatal("NewFlowEnd returned nil")
	}

	if flowEnd.GetKind() != KindFlowEnd {
		t.Errorf("Expected kind %v, got %v", KindFlowEnd, flowEnd.GetKind())
	}

	if flowEnd.GetName() != "" {
		t.Errorf("Expected empty name, got '%s'", flowEnd.GetName())
	}

	if flowEnd.GetLocation() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, flowEnd.GetLocation())
	}
}

func TestFlowParent(t *testing.T) {
	flow := NewFlow("TestFlow", Location{}, true)
	pkg := NewPackage("TestPkg", Location{})

	flow.SetParent(pkg)

	if flow.GetParent() != pkg {
		t.Error("SetParent did not set parent correctly for flow")
	}
}

func TestFlowEndParent(t *testing.T) {
	flowEnd := NewFlowEnd(Location{})
	flow := NewFlow("TestFlow", Location{}, true)

	flowEnd.SetParent(flow)

	if flowEnd.GetParent() != flow {
		t.Error("SetParent did not set parent correctly for flow end")
	}
}

func TestFlowVisitor(t *testing.T) {
	flow := NewFlow("TestFlow", Location{}, true)

	var visited bool
	visitor := &testFlowVisitor{
		onVisit: func(f *Flow) {
			visited = true
			if f != flow {
				t.Error("Visitor received wrong flow")
			}
		},
	}

	flow.Accept(visitor)

	if !visited {
		t.Error("Visitor was not called for flow")
	}
}

func TestFlowEndVisitor(t *testing.T) {
	flowEnd := NewFlowEnd(Location{})

	var visited bool
	visitor := &testFlowEndVisitor{
		onVisit: func(fe *FlowEnd) {
			visited = true
			if fe != flowEnd {
				t.Error("Visitor received wrong flow end")
			}
		},
	}

	flowEnd.Accept(visitor)

	if !visited {
		t.Error("Visitor was not called for flow end")
	}
}

type testFlowVisitor struct {
	BaseVisitor
	onVisit func(*Flow)
}

func (v *testFlowVisitor) VisitFlow(f *Flow) bool {
	if v.onVisit != nil {
		v.onVisit(f)
	}
	return true
}

type testFlowEndVisitor struct {
	BaseVisitor
	onVisit func(*FlowEnd)
}

func (v *testFlowEndVisitor) VisitFlowEnd(fe *FlowEnd) bool {
	if v.onVisit != nil {
		v.onVisit(fe)
	}
	return true
}

func TestModelAddFlow(t *testing.T) {
	model := NewModel()
	flow := NewFlow("TestFlow", Location{Line: 5, Column: 10}, true)

	model.AddFlow(flow)

	if len(model.Flows) != 1 {
		t.Errorf("Expected 1 flow in model, got %d", len(model.Flows))
	}

	if len(model.Elements) != 1 {
		t.Errorf("Expected 1 element in model, got %d", len(model.Elements))
	}

	if model.Flows[0] != flow {
		t.Error("Flow not added correctly to model")
	}
}

func TestPackageFlows(t *testing.T) {
	pkg := NewPackage("TestPkg", Location{})
	flow1 := NewFlow("Flow1", Location{}, true)
	flow2 := NewFlow("Flow2", Location{}, false)

	pkg.AddChild(flow1)
	pkg.AddChild(flow2)

	flows := pkg.Flows()

	if len(flows) != 2 {
		t.Errorf("Expected 2 flows in package, got %d", len(flows))
	}
}

func TestFlowPayloadFeatures(t *testing.T) {
	flow := NewFlow("TestFlow", Location{}, true)

	// PayloadFeatures should be initialized empty
	if flow.PayloadFeatures == nil {
		t.Error("PayloadFeatures should be initialized")
	}

	if len(flow.PayloadFeatures) != 0 {
		t.Errorf("Expected 0 payload features, got %d", len(flow.PayloadFeatures))
	}
}

func TestFlowEnds(t *testing.T) {
	flow := NewFlow("TestFlow", Location{}, true)

	// Initially no source or target
	if flow.Source != nil {
		t.Error("Expected nil source initially")
	}

	if flow.Target != nil {
		t.Error("Expected nil target initially")
	}

	// Set source and target
	sourceEnd := NewFlowEnd(Location{Line: 1, Column: 1})
	targetEnd := NewFlowEnd(Location{Line: 2, Column: 1})

	flow.Source = sourceEnd
	flow.Target = targetEnd

	if flow.Source != sourceEnd {
		t.Error("Source not set correctly")
	}

	if flow.Target != targetEnd {
		t.Error("Target not set correctly")
	}
}

func TestFlowEndReference(t *testing.T) {
	flowEnd := NewFlowEnd(Location{})

	// Initially no reference
	if flowEnd.Reference != nil {
		t.Error("Expected nil reference initially")
	}

	if flowEnd.Feature != nil {
		t.Error("Expected nil feature initially")
	}

	// Create mock elements for reference and feature
	part := NewPart("TestPart", Location{}, true)
	attr := NewAttribute("TestAttr", Location{}, true)

	flowEnd.Reference = part
	flowEnd.Feature = attr

	if flowEnd.Reference != part {
		t.Error("Reference not set correctly")
	}

	if flowEnd.Feature != attr {
		t.Error("Feature not set correctly")
	}
}
