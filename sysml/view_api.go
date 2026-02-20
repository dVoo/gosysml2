package sysml

import "fmt"

// FindView returns the first view whose name or qualified name matches.
func FindView(model *Model, nameOrQualifiedName string) *View {
	if model == nil || nameOrQualifiedName == "" {
		return nil
	}
	for _, view := range FindAll[*View](model) {
		if view.Name() == nameOrQualifiedName || view.QualifiedName() == nameOrQualifiedName {
			return view
		}
	}
	return nil
}

// FindViewpoint returns the first viewpoint whose name or qualified name matches.
func FindViewpoint(model *Model, nameOrQualifiedName string) *Viewpoint {
	if model == nil || nameOrQualifiedName == "" {
		return nil
	}
	for _, viewpoint := range FindAll[*Viewpoint](model) {
		if viewpoint.Name() == nameOrQualifiedName || viewpoint.QualifiedName() == nameOrQualifiedName {
			return viewpoint
		}
	}
	return nil
}

// ElementsForView returns a de-duplicated copy of the elements exposed by a view.
func ElementsForView(view *View) []Element {
	if view == nil {
		return nil
	}
	out := make([]Element, len(view.ExposedElements))
	copy(out, view.ExposedElements)
	return dedupeElements(out)
}

// ElementsByView resolves a view by name/qualified-name and returns exposed elements.
func ElementsByView(model *Model, nameOrQualifiedName string) ([]Element, error) {
	view := FindView(model, nameOrQualifiedName)
	if view == nil {
		return nil, fmt.Errorf("view %q not found", nameOrQualifiedName)
	}
	return ElementsForView(view), nil
}

// ViewsByViewpoint returns all views that resolve to the given viewpoint.
func ViewsByViewpoint(model *Model, viewpointNameOrQualifiedName string) ([]*View, error) {
	vp := FindViewpoint(model, viewpointNameOrQualifiedName)
	if vp == nil {
		return nil, fmt.Errorf("viewpoint %q not found", viewpointNameOrQualifiedName)
	}

	views := make([]*View, 0)
	for _, view := range FindAll[*View](model) {
		if view.Viewpoint.IsResolved() && view.Viewpoint.Resolved() == vp {
			views = append(views, view)
		}
	}
	return views, nil
}

// ElementsByViewpoint returns de-duplicated elements exposed by all views that
// conform to the specified viewpoint.
func ElementsByViewpoint(model *Model, viewpointNameOrQualifiedName string) ([]Element, error) {
	views, err := ViewsByViewpoint(model, viewpointNameOrQualifiedName)
	if err != nil {
		return nil, err
	}

	all := make([]Element, 0)
	for _, view := range views {
		all = append(all, ElementsForView(view)...)
	}
	return dedupeElements(all), nil
}
