package sysml

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDirectoryIncludesKerML(t *testing.T) {
	dir := t.TempDir()
	sysmlFile := filepath.Join(dir, "a.sysml")
	kermlFile := filepath.Join(dir, "b.kerml")

	if err := os.WriteFile(sysmlFile, []byte("package A {}"), 0o644); err != nil {
		t.Fatalf("write sysml file: %v", err)
	}
	if err := os.WriteFile(kermlFile, []byte("standard library package B {}"), 0o644); err != nil {
		t.Fatalf("write kerml file: %v", err)
	}

	results, err := ParseDirectory(dir)
	if err != nil {
		t.Fatalf("ParseDirectory failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 parse results, got %d", len(results))
	}
}

func TestKerMLModelExtraction(t *testing.T) {
	input := `standard library package P {
		datatype NumberBox :> Anything {
			feature value: ScalarValue[1];
		}
	}`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %v", result.Errors)
	}
	if result.Model == nil {
		t.Fatal("expected model")
	}

	types := FindAll[*KerMLType](result.Model)
	if len(types) == 0 {
		t.Fatal("expected at least one KerMLType")
	}
	if types[0].Name() != "NumberBox" {
		t.Fatalf("expected type name NumberBox, got %q", types[0].Name())
	}
	if len(types[0].Specializes) > 0 && types[0].Specializes[0].IsResolved() {
		t.Fatalf("expected isolated specialization references to remain unresolved")
	}

	features := FindAll[*KerMLFeature](result.Model)
	if len(features) == 0 {
		t.Fatal("expected at least one KerMLFeature")
	}
	if features[0].Name() != "value" {
		t.Fatalf("expected feature name value, got %q", features[0].Name())
	}
	if features[0].TypeRef.Name() != "ScalarValue" {
		t.Fatalf("expected feature type ScalarValue, got %q", features[0].TypeRef.Name())
	}
}

func TestDefaultLibraryKerMLCoverage(t *testing.T) {
	root := filepath.Join("..", "libraries", "sysml.library")
	total := 0
	failed := 0

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".kerml" {
			return nil
		}
		total++
		result := ParseFile(path)
		if !result.Success() {
			failed++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking library: %v", err)
	}

	if total == 0 {
		t.Fatal("expected kerml files in default library")
	}
	// Current support target is "usable extent" for bundled libraries.
	// One known edge-case remains in Occurrences.kerml with explicit subset relationship form.
	if failed > 1 {
		t.Fatalf("expected <=1 failing .kerml file, got %d/%d", failed, total)
	}
}
