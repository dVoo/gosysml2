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
	if err := result.Err(); err != nil {
		t.Fatalf("parse failed: %v", err)
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

func TestUsageTypeUnresolvedReturnsNilInterface(t *testing.T) {
	part := NewPart("engine", Location{}, false)
	part.TypeRef = NewRef[*Part]("Engine")

	if got := part.Type(); got != nil {
		t.Fatalf("expected nil unresolved type, got %T", got)
	}
}

func TestUsageTypeNameHelper(t *testing.T) {
	part := NewPart("engine", Location{}, false)
	part.TypeRef = NewRef[*Part]("Engine")

	var usage Usage = part
	if got := usage.TypeName(); got != "Engine" {
		t.Fatalf("expected unresolved type name via method %q, got %q", "Engine", got)
	}
	if got := UsageTypeName(usage); got != "Engine" {
		t.Fatalf("expected unresolved type name %q, got %q", "Engine", got)
	}

	resolved := NewPart("ResolvedEngine", Location{}, true)
	part.TypeRef.Resolve(resolved)
	if got := usage.TypeName(); got != "ResolvedEngine" {
		t.Fatalf("expected resolved type name via method %q, got %q", "ResolvedEngine", got)
	}
	if got := UsageTypeName(usage); got != "ResolvedEngine" {
		t.Fatalf("expected resolved type name %q, got %q", "ResolvedEngine", got)
	}
}

func TestDefinitionUsageHelpers(t *testing.T) {
	partDef := NewPart("Engine", Location{}, true)
	partUse := NewPart("engine", Location{}, false)
	dep := NewDependency(Location{})

	if !IsDefinitionElement(partDef) || IsUsageElement(partDef) {
		t.Fatal("expected part definition classification helpers to report definition only")
	}
	if IsDefinitionElement(partUse) || !IsUsageElement(partUse) {
		t.Fatal("expected part usage classification helpers to report usage only")
	}
	if !IsDefinitionElement(dep) || !IsUsageElement(dep) {
		t.Fatal("expected dependency to retain existing dual-role classification")
	}
	if got := RoleOf(partDef); got != RoleDefinition {
		t.Fatalf("expected RoleDefinition, got %v", got)
	}
	if got := partDef.Role(); got != RoleDefinition {
		t.Fatalf("expected partDef.Role() = RoleDefinition, got %v", got)
	}
	if got := RoleOf(partUse); got != RoleUsage {
		t.Fatalf("expected RoleUsage, got %v", got)
	}
	if got := partUse.Role(); got != RoleUsage {
		t.Fatalf("expected partUse.Role() = RoleUsage, got %v", got)
	}
	if got := RoleOf(dep); got != RoleDefinitionAndUsage {
		t.Fatalf("expected RoleDefinitionAndUsage, got %v", got)
	}
	if got := dep.Role(); got != RoleDefinitionAndUsage {
		t.Fatalf("expected dep.Role() = RoleDefinitionAndUsage, got %v", got)
	}
}

func TestElementAtIncludingUnnamed(t *testing.T) {
	input := `
package P {
	interface def EngineInterface;
	part def Car {
		interface : EngineInterface;
	}
}
`
	result := ParseString(input)
	if err := result.Err(); err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if result.Model == nil {
		t.Fatal("expected model")
	}

	// Position on the anonymous interface usage declaration line.
	line := 4 // zero-based in this raw string
	col := 3
	if got := ElementAt(result.Model, line, col); got != nil {
		t.Fatalf("expected named-only ElementAt to return nil for anonymous usage, got %T", got)
	}
	if got := ElementAtIncludingUnnamed(result.Model, line, col); got == nil {
		t.Fatal("expected unnamed-inclusive lookup to return an element")
	}
}

func TestDocMemberPopulatesDocumentation(t *testing.T) {
	input := `
package P {
	part def A;
	doc /* Primary part doc. */;
}
`
	result := ParseString(input)
	if err := result.Err(); err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	pkgs := result.Model.Packages()
	if len(pkgs) != 1 {
		t.Fatalf("expected one package, got %d", len(pkgs))
	}
	parts := pkgs[0].Parts()
	if len(parts) != 1 {
		t.Fatalf("expected one part, got %d", len(parts))
	}
	if got := parts[0].Documentation(); got == "" {
		t.Fatal("expected part documentation to be populated from doc member")
	}
}
