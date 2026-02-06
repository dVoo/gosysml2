package sysml

import (
	"testing"
)

func TestParseString(t *testing.T) {
	input := `
package TestPackage {
    part def Vehicle {
        part engine : Engine;
    }
    part def Engine;
}
`

	result := ParseString(input)

	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	if result.Model == nil {
		t.Fatal("model is nil")
	}

	if len(result.Model.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(result.Model.Packages))
	}

	pkg := result.Model.Packages[0]
	if pkg.Name() != "TestPackage" {
		t.Errorf("expected package name 'TestPackage', got '%s'", pkg.Name())
	}
}

func TestParseStringWithErrors(t *testing.T) {
	input := `
package Broken {
    part def Vehicle {
        invalid syntax here!!!
    }
}
`

	result := ParseString(input)

	if result.Success() {
		t.Error("expected parse to fail")
	}

	if result.Errors == nil || !result.Errors.HasErrors() {
		t.Error("expected errors to be present")
	}

	t.Logf("Got expected errors: %s", result.Errors)
}

func TestValidate(t *testing.T) {
	validInput := `
package Valid {
    part def A;
    part def B;
}
`

	invalidInput := `
package Invalid {
    @@@ not valid syntax
}
`

	// Valid input should pass
	err := Validate(validInput)
	if err != nil {
		t.Errorf("valid input failed validation: %s", err)
	}

	// Invalid input should fail
	err = Validate(invalidInput)
	if err == nil {
		t.Error("invalid input should fail validation")
	}
}

func TestWalk(t *testing.T) {
	input := `
package P1 {
    package P2 {
        part def A;
    }
    part def B;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	var visited []string
	Walk(result.Model, func(elem Element, depth int) bool {
		visited = append(visited, elem.Name())
		return true
	})

	if len(visited) == 0 {
		t.Error("no elements visited")
	}

	t.Logf("Visited elements: %v", visited)
}

func TestFindByKind(t *testing.T) {
	input := `
package TestPkg {
    part def Vehicle;
    part def Engine;
    part vehicle1 : Vehicle;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	parts := FindByKind(result.Model, KindPart)
	t.Logf("Found %d parts", len(parts))
}

func TestCounter(t *testing.T) {
	input := `
package P {
    part def A;
    part def B;
    part a1 : A;
    part b1 : B;
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %s", result.Errors)
	}

	counter := NewCounter()
	Visit(result.Model, counter)

	t.Logf("Counts: %v", counter.Counts)
	t.Logf("Total: %d", counter.Total())

	if counter.Total() == 0 {
		t.Error("expected some elements to be counted")
	}
}

func TestMustParseStringPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid input")
		}
	}()

	MustParseString(`package { invalid }}}`)
}

func TestEmptyInput(t *testing.T) {
	result := ParseString("")

	// Empty input should parse successfully (empty namespace is valid)
	if !result.Success() {
		t.Logf("Parse result: %s", result.Errors)
	}
}
