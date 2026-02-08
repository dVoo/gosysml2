package sysml

import (
	"strings"
	"testing"
)

// TestPartString verifies that Part.String() returns useful debug information
func TestPartString(t *testing.T) {
	// Test definition
	partDef := NewPart("Vehicle", Location{Line: 1, Column: 1}, true)
	str := partDef.String()
	if !strings.Contains(str, "Vehicle") {
		t.Error("Part.String() should contain the name")
	}
	if !strings.Contains(str, "definition") {
		t.Error("Part.String() should indicate it's a definition")
	}

	// Test usage
	partUsage := NewPart("myCar", Location{Line: 2, Column: 1}, false)
	partUsage.TypeRef = NewRef[*Part]("Car")
	str = partUsage.String()
	if !strings.Contains(str, "myCar") {
		t.Error("Part.String() should contain the name")
	}
	if !strings.Contains(str, "usage") {
		t.Error("Part.String() should indicate it's a usage")
	}

	// Test with specialization
	partWithSpec := NewPart("Car", Location{Line: 1, Column: 1}, true)
	partWithSpec.SetUnresolvedSpecializes("Vehicle")
	str = partWithSpec.String()
	if !strings.Contains(str, "Car") {
		t.Error("Part.String() should contain the name")
	}
	if !strings.Contains(str, "Vehicle (unresolved)") {
		t.Error("Part.String() should show unresolved specialization")
	}
}

// TestItemString verifies that Item.String() returns useful debug information
func TestItemString(t *testing.T) {
	item := NewItem("Product", Location{Line: 1, Column: 1}, true)
	str := item.String()
	if !strings.Contains(str, "Product") {
		t.Error("Item.String() should contain the name")
	}
	if !strings.Contains(str, "Item") {
		t.Error("Item.String() should indicate it's an item")
	}
	if !strings.Contains(str, "definition") {
		t.Error("Item.String() should indicate it's a definition")
	}

	// Test with specialization
	itemWithSpec := NewItem("Widget", Location{Line: 1, Column: 1}, true)
	itemWithSpec.SetUnresolvedSpecializes("Product")
	str = itemWithSpec.String()
	if !strings.Contains(str, "Widget") {
		t.Error("Item.String() should contain the name")
	}
	if !strings.Contains(str, "Product (unresolved)") {
		t.Error("Item.String() should show unresolved specialization")
	}
}
