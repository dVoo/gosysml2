package sysml

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRefEffectiveName(t *testing.T) {
	unresolved := NewRef[*Part]("Vehicle")
	if got := unresolved.EffectiveName(); got != "Vehicle" {
		t.Fatalf("expected unresolved EffectiveName %q, got %q", "Vehicle", got)
	}

	resolvedPart := NewPart("ResolvedVehicle", Location{}, true)
	resolved := NewRef[*Part]("Vehicle")
	resolved.Resolve(resolvedPart)
	if got := resolved.EffectiveName(); got != "ResolvedVehicle" {
		t.Fatalf("expected resolved EffectiveName %q, got %q", "ResolvedVehicle", got)
	}
}

func TestParseStringModel(t *testing.T) {
	model, err := ParseStringModel("package P { part def A; }")
	if err != nil {
		t.Fatalf("expected parse success, got error: %v", err)
	}
	if model == nil {
		t.Fatal("expected non-nil model")
	}

	if _, err := ParseStringModel("package P { !!! }"); err == nil {
		t.Fatal("expected parse error for invalid input")
	}
}

func TestParseFileModel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "model.sysml")
	if err := os.WriteFile(path, []byte("package P { part def A; }"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	model, err := ParseFileModel(path)
	if err != nil {
		t.Fatalf("expected parse success, got error: %v", err)
	}
	if model == nil {
		t.Fatal("expected non-nil model")
	}
}

func TestWalkAllHelpers(t *testing.T) {
	result := ParseString(`
package P {
	part def A;
	part def B;
}
`)
	if !result.Success() {
		t.Fatalf("parse failed: %v", result.Errors)
	}

	var methodCount int
	result.Model.WalkAll(func(Element) {
		methodCount++
	})
	if methodCount == 0 {
		t.Fatal("expected model.WalkAll to visit elements")
	}

	var funcCount int
	WalkAll(result.Model, func(Element, int) {
		funcCount++
	})
	if funcCount != methodCount {
		t.Fatalf("expected WalkAll function count %d to match model.WalkAll count %d", funcCount, methodCount)
	}
}
