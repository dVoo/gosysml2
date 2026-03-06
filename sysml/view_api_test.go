package sysml

import "testing"

func TestViewpointAndViewExposeExtractionAndAPI(t *testing.T) {
	input := `
package Model {
	part def Vehicle;
	part vehicle : Vehicle {
		part wheel;
	}
	part systemsEngineer;

	concern systemBreakdown {
		stakeholder :>> systemsEngineer;
	}

	viewpoint structurePerspective {
		frame systemBreakdown;
	}

	view structureGeneration {
		satisfy structurePerspective;
		expose Model::vehicle::**[@SysML::PartUsage];
	}
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %v", result.Err())
	}

	views := FindAll[*View](result.Model)
	if len(views) != 1 {
		t.Fatalf("expected 1 view, got %d", len(views))
	}
	view := views[0]

	if len(view.Exposures) != 1 {
		t.Fatalf("expected 1 expose clause, got %d", len(view.Exposures))
	}
	if view.Exposures[0].Namespace != "Model::vehicle" {
		t.Fatalf("expected expose namespace Model::vehicle, got %q", view.Exposures[0].Namespace)
	}

	if !view.Viewpoint.IsResolved() {
		t.Fatalf("expected view viewpoint to resolve")
	}
	if got := view.Viewpoint.Resolved().Name(); got != "structurePerspective" {
		t.Fatalf("expected viewpoint structurePerspective, got %q", got)
	}

	exposed := ElementsForView(view)
	if len(exposed) == 0 {
		t.Fatalf("expected exposed elements")
	}
	foundWheel := false
	for _, elem := range exposed {
		if p, ok := elem.(*Part); ok && p.Name() == "wheel" && !p.IsDefinition {
			foundWheel = true
			break
		}
	}
	if !foundWheel {
		t.Fatalf("expected wheel part usage in exposed elements")
	}

	byView, err := ElementsByView(result.Model, "structureGeneration")
	if err != nil {
		t.Fatalf("ElementsByView error: %v", err)
	}
	if len(byView) == 0 {
		t.Fatalf("expected non-empty elements by view")
	}

	byVP, err := ElementsByViewpoint(result.Model, "structurePerspective")
	if err != nil {
		t.Fatalf("ElementsByViewpoint error: %v", err)
	}
	if len(byVP) == 0 {
		t.Fatalf("expected non-empty elements by viewpoint")
	}
}

func TestViewpointConcernAndStakeholderResolution(t *testing.T) {
	input := `
package P {
	part stakeholder1;
	concern c1 {
		stakeholder :>> stakeholder1;
	}
	viewpoint vp1 {
		frame c1;
	}
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %v", result.Err())
	}

	concerns := FindAll[*Concern](result.Model)
	if len(concerns) != 1 {
		t.Fatalf("expected 1 concern, got %d", len(concerns))
	}
	if len(concerns[0].Stakeholders) != 1 {
		t.Fatalf("expected concern stakeholder resolution, got %d", len(concerns[0].Stakeholders))
	}

	vps := FindAll[*Viewpoint](result.Model)
	if len(vps) != 1 {
		t.Fatalf("expected 1 viewpoint, got %d", len(vps))
	}
	if len(vps[0].Concerns) != 1 {
		t.Fatalf("expected viewpoint concern resolution, got %d", len(vps[0].Concerns))
	}
}

func TestViewExposeResolvesImportedNamespaceReferences(t *testing.T) {
	input := `
package ASPICE_Toolchain_Model {
	package Tool_Cluster {
		package Tools {
			part def Tool;
			part Codebeamer : Tool;
			part Jira : Tool;
		}
	}

	package Views {
		import Tool_Cluster::*;

		view def ToolArchitectureView;

		view toolArchitecture : ToolArchitectureView {
			expose Tools::Codebeamer;
			expose Tools::*;
		}
	}
}
`

	result := ParseString(input)
	if !result.Success() {
		t.Fatalf("parse failed: %v", result.Err())
	}

	elements, err := ElementsByView(result.Model, "toolArchitecture")
	if err != nil {
		t.Fatalf("ElementsByView error: %v", err)
	}
	if len(elements) == 0 {
		t.Fatalf("expected imported-scope expose references to resolve, got 0 elements")
	}

	foundCodebeamer := false
	for _, elem := range elements {
		part, ok := elem.(*Part)
		if ok && part.Name() == "Codebeamer" && !part.IsDefinition {
			foundCodebeamer = true
			break
		}
	}
	if !foundCodebeamer {
		t.Fatalf("expected Codebeamer part usage to be exposed via imported scope")
	}
}
