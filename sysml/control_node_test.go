package sysml

import (
	"testing"
)

func TestNewControlNode(t *testing.T) {
	loc := Location{Line: 10, Column: 5}

	// Test all four control node kinds
	testCases := []struct {
		name     string
		kind     ControlNodeKind
		expected string
	}{
		{"ForkNode", ControlNodeKindFork, "fork"},
		{"JoinNode", ControlNodeKindJoin, "join"},
		{"MergeNode", ControlNodeKindMerge, "merge"},
		{"DecisionNode", ControlNodeKindDecision, "decision"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := NewControlNode(tc.kind, loc)

			if node == nil {
				t.Fatal("NewControlNode returned nil")
			}

			if node.NodeKind != tc.kind {
				t.Errorf("Expected kind %v, got %v", tc.kind, node.NodeKind)
			}

			if node.NodeKind.String() != tc.expected {
				t.Errorf("Expected string '%s', got '%s'", tc.expected, node.NodeKind.String())
			}

			if node.Location() != loc {
				t.Errorf("Expected location %+v, got %+v", loc, node.Location())
			}

			if !node.IsControlNode {
				t.Error("Expected IsControlNode to be true")
			}
		})
	}
}

func TestControlNodeKindString(t *testing.T) {
	testCases := []struct {
		kind     ControlNodeKind
		expected string
	}{
		{ControlNodeKindFork, "fork"},
		{ControlNodeKindJoin, "join"},
		{ControlNodeKindMerge, "merge"},
		{ControlNodeKindDecision, "decision"},
		{ControlNodeKind(999), "unknown"},
	}

	for _, tc := range testCases {
		result := tc.kind.String()
		if result != tc.expected {
			t.Errorf("ControlNodeKind(%v).String() = '%s', expected '%s'", tc.kind, result, tc.expected)
		}
	}
}

func TestControlNodeTypeChecks(t *testing.T) {
	loc := Location{}

	fork := NewControlNode(ControlNodeKindFork, loc)
	if !fork.IsFork() {
		t.Error("Fork node should return true for IsFork()")
	}
	if fork.IsJoin() || fork.IsMerge() || fork.IsDecision() {
		t.Error("Fork node should not be join, merge, or decision")
	}

	join := NewControlNode(ControlNodeKindJoin, loc)
	if !join.IsJoin() {
		t.Error("Join node should return true for IsJoin()")
	}

	merge := NewControlNode(ControlNodeKindMerge, loc)
	if !merge.IsMerge() {
		t.Error("Merge node should return true for IsMerge()")
	}

	decision := NewControlNode(ControlNodeKindDecision, loc)
	if !decision.IsDecision() {
		t.Error("Decision node should return true for IsDecision()")
	}
}

func TestControlNodeCondition(t *testing.T) {
	loc := Location{}
	node := NewControlNode(ControlNodeKindDecision, loc)

	// Initially no condition
	if node.Condition != "" {
		t.Errorf("Expected empty condition initially, got '%s'", node.Condition)
	}

	// Set condition
	node.SetCondition("x > 0")
	if node.Condition != "x > 0" {
		t.Errorf("Expected condition 'x > 0', got '%s'", node.Condition)
	}
}

func TestControlNodeParent(t *testing.T) {
	node := NewControlNode(ControlNodeKindFork, Location{})
	pkg := NewPackage("TestPkg", Location{})

	node.SetParent(pkg)

	if node.Parent() != pkg {
		t.Error("SetParent did not set parent correctly for control node")
	}
}

func TestControlNodeAddChild(t *testing.T) {
	node := NewControlNode(ControlNodeKindMerge, Location{})
	child := NewControlNode(ControlNodeKindFork, Location{})

	// Control nodes can have children (for nested control nodes)
	node.AddChild(child)

	children := node.Children()
	if len(children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(children))
	}

	if children[0] != child {
		t.Error("Child not added correctly to control node")
	}
}

func TestModelAddControlNode(t *testing.T) {
	model := NewModel()
	node := NewControlNode(ControlNodeKindFork, Location{Line: 5, Column: 10})

	model.AddControlNode(node)

	if len(model.ControlNodes) != 1 {
		t.Errorf("Expected 1 control node in model, got %d", len(model.ControlNodes))
	}

	if len(model.Elements) != 1 {
		t.Errorf("Expected 1 element in model, got %d", len(model.Elements))
	}

	if model.ControlNodes[0] != node {
		t.Error("Control node not added correctly to model")
	}
}

func TestControlNodeElementInterface(t *testing.T) {
	loc := Location{Line: 10, Column: 20}
	node := NewControlNode(ControlNodeKindDecision, loc)

	// Test Element interface methods
	if node.Kind() != KindControlNode {
		t.Errorf("Expected Kind() to return KindControlNode, got %v", node.Kind())
	}

	if node.Name() != "" {
		t.Errorf("Expected empty name, got '%s'", node.Name())
	}

	if node.Location() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, node.Location())
	}

	if node.Parent() != nil {
		t.Error("Expected nil parent for new node")
	}

	if node.Documentation() != "" {
		t.Errorf("Expected empty documentation, got '%s'", node.Documentation())
	}
}

func TestControlNodeQualifiedName(t *testing.T) {
	pkg := NewPackage("TestPkg", Location{})
	node := NewControlNode(ControlNodeKindFork, Location{})

	// Initially no parent, so qualified name is empty
	if node.QualifiedName() != "" {
		t.Errorf("Expected empty qualified name without parent, got '%s'", node.QualifiedName())
	}

	// Add to package
	pkg.AddChild(node)

	// Now qualified name should include package
	expected := "TestPkg::"
	qn := node.QualifiedName()
	if qn != expected {
		t.Errorf("Expected qualified name '%s', got '%s'", expected, qn)
	}
}
