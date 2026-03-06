package sysml

import "testing"

func mustParsePhase5File(t *testing.T, path string) *Model {
	t.Helper()

	resolvedPath := externalTestPath(t, path)
	result := ParseFile(resolvedPath)
	if !result.Success() {
		t.Fatalf("parse failed for %s: %v", resolvedPath, result.Err())
	}
	if result.Model == nil {
		t.Fatalf("model is nil for %s", resolvedPath)
	}
	return result.Model
}

func TestPhase5RepresentativeSupportedRulesNotDropped(t *testing.T) {
	t.Run("alias", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/13-Model Containment/13a-Model Containment.sysml")
		aliases := FindAll[*Alias](model)
		if len(aliases) == 0 {
			t.Fatal("expected at least one alias element")
		}
	})

	t.Run("metadata and filters", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/13-Model Containment/13b-Safety and Security Features Element Group-1.sysml")
		metadata := FindAll[*Metadata](model)
		if len(metadata) == 0 {
			t.Fatal("expected at least one metadata element")
		}
		filters := FindAll[*ElementFilter](model)
		if len(filters) == 0 {
			t.Fatal("expected at least one element filter")
		}
	})

	t.Run("message", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/17-Sequence Modeling/17a-Sequence-Modeling.sysml")
		messages := FindAll[*Message](model)
		if len(messages) == 0 {
			t.Fatal("expected at least one message element")
		}
	})

	t.Run("individual and portions", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/06-Individual and Snapshots/6-Individual and Snapshots.sysml")
		occurrences := FindAll[*Occurrence](model)
		if len(occurrences) == 0 {
			t.Fatal("expected at least one occurrence element")
		}

		hasIndividual := false
		hasPortion := false
		for _, occ := range occurrences {
			if occ.IsIndividual {
				hasIndividual = true
			}
			if occ.IsSnapshot || occ.IsTimeSlice || occ.PortionKind() != PortionKindNone {
				hasPortion = true
			}
		}

		if !hasIndividual {
			t.Fatal("expected at least one individual occurrence")
		}
		if !hasPortion {
			t.Fatal("expected at least one snapshot/timeslice occurrence")
		}
	})

	t.Run("satisfy", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/08-Requirements/8-Requirements.sysml")
		rels := FindAll[*SatisfyRelationship](model)
		if len(rels) == 0 {
			t.Fatal("expected at least one satisfy relationship")
		}
	})

	t.Run("verify", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/09-Verification/9-Verification-simplified.sysml")
		rels := FindAll[*VerifyRelationship](model)
		if len(rels) == 0 {
			t.Fatal("expected at least one verify relationship")
		}
	})

	t.Run("rendering", func(t *testing.T) {
		model := mustParsePhase5File(t, "../libraries/sysml.library/Systems Library/Views.sysml")
		renderings := FindAll[*Rendering](model)
		if len(renderings) == 0 {
			t.Fatal("expected at least one rendering element")
		}
	})
}

func TestPhase5ReferenceAndTypeFieldsPopulatedFromGrammar(t *testing.T) {
	t.Run("alias target", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/13-Model Containment/13a-Model Containment.sysml")
		aliases := FindAll[*Alias](model)
		for _, alias := range aliases {
			if alias.unresolvedTarget != "" {
				return
			}
		}
		t.Fatal("expected at least one alias unresolved target populated from grammar")
	})

	t.Run("message endpoints", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/17-Sequence Modeling/17a-Sequence-Modeling.sysml")
		messages := FindAll[*Message](model)
		for _, msg := range messages {
			if msg.unresolvedSender != "" && msg.unresolvedReceiver != "" {
				return
			}
		}
		t.Fatal("expected at least one message with sender/receiver unresolved fields populated")
	})

	t.Run("satisfy relationship refs", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/08-Requirements/8-Requirements.sysml")
		rels := FindAll[*SatisfyRelationship](model)
		for _, rel := range rels {
			if rel.unresolvedSatisfier != "" && rel.unresolvedRequired != "" {
				return
			}
		}
		t.Fatal("expected at least one satisfy relationship with unresolved fields populated")
	})

	t.Run("verify relationship refs", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/09-Verification/9-Verification-simplified.sysml")
		rels := FindAll[*VerifyRelationship](model)
		for _, rel := range rels {
			if rel.unresolvedRequired != "" {
				return
			}
		}
		t.Fatal("expected at least one verify relationship with unresolved requirement populated")
	})

	t.Run("rendering usage typing", func(t *testing.T) {
		model := mustParsePhase5File(t, "../libraries/sysml.library/Systems Library/Views.sysml")
		renderings := FindAll[*Rendering](model)
		for _, rendering := range renderings {
			if !rendering.IsDefinition && rendering.TypeRef.name != "" {
				return
			}
		}
		t.Fatal("expected at least one rendering usage with type reference populated")
	})

	t.Run("filter expression", func(t *testing.T) {
		model := mustParsePhase5File(t, "../validationdata/13-Model Containment/13b-Safety and Security Features Element Group-1.sysml")
		filters := FindAll[*ElementFilter](model)
		for _, filter := range filters {
			if filter.Expression != "" {
				return
			}
		}
		t.Fatal("expected at least one element filter expression populated from grammar")
	})
}
