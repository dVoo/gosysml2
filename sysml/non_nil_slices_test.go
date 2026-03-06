package sysml

import (
	"reflect"
	"testing"
)

func assertSliceHasNoNilElements(t *testing.T, path string, v any) {
	t.Helper()
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		t.Fatalf("%s: expected slice, got %s", path, rv.Kind())
	}
	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i)
		switch item.Kind() {
		case reflect.Interface, reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
			if item.IsNil() {
				t.Fatalf("%s[%d] is nil", path, i)
			}
		}
	}
}

func assertExportedSliceFieldsNoNil(t *testing.T, path string, v any) {
	t.Helper()
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			t.Fatalf("%s is nil", path)
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fv := rv.Field(i)
		if fv.Kind() != reflect.Slice {
			continue
		}
		assertSliceHasNoNilElements(t, path+"."+field.Name, fv.Interface())
	}
}

func TestModelAddMethodsRejectNil(t *testing.T) {
	m := NewModel()

	m.AddDoc(nil)
	m.AddDependency(nil)
	m.AddPackage(nil)
	m.AddImport(nil)
	m.AddComment(nil)
	m.AddControlNode(nil)
	m.AddOccurrence(nil)
	m.AddAlias(nil)
	m.AddMetadata(nil)
	m.AddRendering(nil)
	m.AddMessage(nil)
	m.AddFilter(nil)
	m.AddSatisfy(nil)
	m.AddVerify(nil)

	if len(m.Elements) != 0 {
		t.Fatalf("expected no elements after nil adds, got %d", len(m.Elements))
	}
}

func TestNoNilElementsInPublicSlices(t *testing.T) {
	result := ParseString(`
package P {
	part def Vehicle {
		attribute mass : Real;
	}
	part vehicle1 : Vehicle;
	requirement def R;
	verification def V;
}
`)
	if !result.Success() {
		t.Fatalf("parse failed: %v", result.Err())
	}

	model := result.Model
	if model == nil {
		t.Fatal("model is nil")
	}

	assertExportedSliceFieldsNoNil(t, "Model", model)
	assertSliceHasNoNilElements(t, "Model.Elements", model.Elements)

	for elem := range All(model) {
		assertExportedSliceFieldsNoNil(t, elem.Kind().String(), elem)

		switch e := elem.(type) {
		case *Package:
			assertSliceHasNoNilElements(t, "Package.Packages()", e.Packages())
			assertSliceHasNoNilElements(t, "Package.Parts()", e.Parts())
			assertSliceHasNoNilElements(t, "Package.Requirements()", e.Requirements())
			assertSliceHasNoNilElements(t, "Package.Actions()", e.Actions())
			assertSliceHasNoNilElements(t, "Package.Imports()", e.Imports())
			assertSliceHasNoNilElements(t, "Package.Items()", e.Items())
			assertSliceHasNoNilElements(t, "Package.States()", e.States())
			assertSliceHasNoNilElements(t, "Package.Connections()", e.Connections())
			assertSliceHasNoNilElements(t, "Package.Interfaces()", e.Interfaces())
			assertSliceHasNoNilElements(t, "Package.Allocations()", e.Allocations())
			assertSliceHasNoNilElements(t, "Package.Views()", e.Views())
			assertSliceHasNoNilElements(t, "Package.Viewpoints()", e.Viewpoints())
			assertSliceHasNoNilElements(t, "Package.Calculations()", e.Calculations())
			assertSliceHasNoNilElements(t, "Package.Enumerations()", e.Enumerations())
			assertSliceHasNoNilElements(t, "Package.Constraints()", e.Constraints())
			assertSliceHasNoNilElements(t, "Package.Dependencies()", e.Dependencies())
		case *Part:
			assertSliceHasNoNilElements(t, "Part.Attributes()", e.Attributes())
			assertSliceHasNoNilElements(t, "Part.Parts()", e.Parts())
			assertSliceHasNoNilElements(t, "Part.Ports()", e.Ports())
		case *Port:
			assertSliceHasNoNilElements(t, "Port.Ports()", e.Ports())
			assertSliceHasNoNilElements(t, "Port.Parts()", e.Parts())
		case *Requirement:
			assertSliceHasNoNilElements(t, "Requirement.Requirements()", e.Requirements())
		case *Action:
			assertSliceHasNoNilElements(t, "Action.Actions()", e.Actions())
		case *Verification:
			assertSliceHasNoNilElements(t, "Verification.Actions()", e.Actions())
		case *State:
			assertSliceHasNoNilElements(t, "State.States()", e.States())
			assertSliceHasNoNilElements(t, "State.Transitions()", e.Transitions())
		case *Item:
			assertSliceHasNoNilElements(t, "Item.Attributes()", e.Attributes())
			assertSliceHasNoNilElements(t, "Item.Items()", e.Items())
		case *Interface:
			assertSliceHasNoNilElements(t, "Interface.Ports()", e.Ports())
		case *Enumeration:
			assertSliceHasNoNilElements(t, "Enumeration.Values()", e.Values())
		}
	}
}
